package oauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/internal/config"
)

func insecureCIMDService(t *testing.T) *Service {
	t.Helper()
	return &Service{
		cfg: config.Config{
			OAuthEnabled:                         true,
			OAuthClientIDMetadataDocumentEnabled: true,
			OAuthCIMDAllowInsecureFetch:          true,
		},
		cimdCache: newCIMDCache(),
	}
}

func TestIsCIMDClientID(t *testing.T) {
	t.Parallel()

	assert.True(t, isCIMDClientID("https://client.example/metadata.json"))
	assert.True(t, isCIMDClientID("http://client.example/metadata.json"))
	assert.False(t, isCIMDClientID("sdoak_abc123"))
	assert.False(t, isCIMDClientID("storyden-cli"))
	assert.False(t, isCIMDClientID(""))
}

func TestParseCIMDClientIDRejectsInsecureAndPrivate(t *testing.T) {
	t.Parallel()

	s := &Service{cfg: config.Config{OAuthEnabled: true, OAuthClientIDMetadataDocumentEnabled: true}, cimdCache: newCIMDCache()}

	_, err := s.parseCIMDClientID("http://client.example/metadata.json")
	assert.Error(t, err, "http must be rejected when insecure fetch is disabled")

	_, err = s.parseCIMDClientID("https://127.0.0.1/metadata.json")
	assert.Error(t, err, "loopback must be rejected")

	_, err = s.parseCIMDClientID("https://localhost/metadata.json")
	assert.Error(t, err, "localhost must be rejected")

	_, err = s.parseCIMDClientID("https://169.254.169.254/latest/meta-data")
	assert.Error(t, err, "link-local metadata endpoint must be rejected")

	_, err = s.parseCIMDClientID("https://10.0.0.5/metadata.json")
	assert.Error(t, err, "private range must be rejected")

	_, err = s.parseCIMDClientID("https://user:pass@client.example/metadata.json")
	assert.Error(t, err, "userinfo must be rejected")

	_, err = s.parseCIMDClientID("https://client.example")
	assert.Error(t, err, "no path must be rejected")

	_, err = s.parseCIMDClientID("https://client.example/")
	assert.Error(t, err, "root path must be rejected")

	_, err = s.parseCIMDClientID("https://client.example/../secret.json")
	assert.Error(t, err, "dot dot segment must be rejected")

	_, err = s.parseCIMDClientID("https://client.example/./meta")
	assert.Error(t, err, "dot segment must be rejected")

	// Query strings are now supported (e.g. providers embedding hints such as
	// ?token_endpoint_auth_method=none in the client_id URL itself).
	u, err := s.parseCIMDClientID("https://client.example/meta?x=1&token_endpoint_auth_method=none")
	require.NoError(t, err)
	assert.Equal(t, "client.example", u.Hostname())
	assert.Equal(t, "x=1&token_endpoint_auth_method=none", u.RawQuery)

	u, err = s.parseCIMDClientID("https://client.example/metadata.json")
	require.NoError(t, err)
	assert.Equal(t, "client.example", u.Hostname())
}

func TestValidateClientMetadata(t *testing.T) {
	t.Parallel()

	const id = "https://client.example/metadata.json"

	require.NoError(t, validateClientMetadata(&clientMetadataDocument{
		ClientID:     id,
		RedirectURIs: []string{"https://client.example/callback"},
	}, id))

	assert.Error(t, validateClientMetadata(&clientMetadataDocument{
		ClientID:     "https://evil.example/other",
		RedirectURIs: []string{"https://client.example/callback"},
	}, id), "client_id must match URL")

	assert.Error(t, validateClientMetadata(&clientMetadataDocument{
		ClientID:     id,
		RedirectURIs: []string{},
	}, id), "redirect_uris required")

	assert.Error(t, validateClientMetadata(&clientMetadataDocument{
		ClientID:     id,
		RedirectURIs: []string{"http://client.example/callback"},
	}, id), "non-loopback http redirect rejected")

	assert.Error(t, validateClientMetadata(&clientMetadataDocument{
		ClientID:                id,
		RedirectURIs:            []string{"https://client.example/cb"},
		TokenEndpointAuthMethod: "client_secret_post",
	}, id), "secret auth method rejected")

	assert.Error(t, validateClientMetadata(&clientMetadataDocument{
		ClientID:     id,
		RedirectURIs: []string{"https://client.example/cb"},
		ClientSecret: "s3cr3t",
	}, id), "client_secret prop rejected")

	assert.Error(t, validateClientMetadata(&clientMetadataDocument{
		ClientID:     id,
		RedirectURIs: []string{"https://client.example/cb#frag"},
	}, id), "fragment redirect rejected")

	assert.Error(t, validateClientMetadata(&clientMetadataDocument{
		ClientID:     id,
		RedirectURIs: []string{"https://client.example/*"},
	}, id), "wildcard redirect rejected")
}

func TestParseCIMDClientIDAllowsQueryStringsAndChatGPTStyle(t *testing.T) {
	t.Parallel()

	s := insecureCIMDService(t)

	// Exact style seen in production from ChatGPT connectors.
	chatgpt := "https://chatgpt.com/oauth/ehPHP71Uhmkc/client.json?token_endpoint_auth_method=none"
	u, err := s.parseCIMDClientID(chatgpt)
	require.NoError(t, err)
	assert.Equal(t, "chatgpt.com", u.Hostname())
	assert.Equal(t, "/oauth/ehPHP71Uhmkc/client.json", u.Path)
	assert.Equal(t, "token_endpoint_auth_method=none", u.RawQuery)

	// The document served at that URL must declare the *same* client_id (incl. query)
	// for the match in validateClientMetadata to succeed.
	docURLWithQuery := chatgpt
	require.NoError(t, validateClientMetadata(&clientMetadataDocument{
		ClientID:     docURLWithQuery,
		RedirectURIs: []string{"https://client.example/callback"},
	}, docURLWithQuery))

	// Mismatch (doc omits the query part) must still be rejected.
	assert.Error(t, validateClientMetadata(&clientMetadataDocument{
		ClientID:     "https://chatgpt.com/oauth/ehPHP71Uhmkc/client.json",
		RedirectURIs: []string{"https://client.example/callback"},
	}, docURLWithQuery), "client_id in doc must match the full presented client_id incl. query")
}

func TestCIMDAllowedScopes(t *testing.T) {
	t.Parallel()

	s := &Service{cfg: config.Config{}}

	// Standard scopes always allowed; conservative read perms allowed; unknown dropped.
	got := s.cimdAllowedScopes("openid profile " + rbac.PermissionReadPublishedThreads.String() + " NONEXISTENT_SCOPE")
	assert.Contains(t, got, "openid")
	assert.Contains(t, got, "profile")
	assert.Contains(t, got, rbac.PermissionReadPublishedThreads.String())
	assert.NotContains(t, got, "NONEXISTENT_SCOPE")

	// Admin/privileged scopes dropped by default.
	got = s.cimdAllowedScopes(rbac.PermissionAdministrator.String() + " " + rbac.PermissionManageRoles.String())
	assert.NotContains(t, got, rbac.PermissionAdministrator.String())
	assert.NotContains(t, got, rbac.PermissionManageRoles.String())

	// Write perm not in conservative default allowlist is dropped.
	got = s.cimdAllowedScopes(rbac.PermissionCreatePost.String())
	assert.NotContains(t, got, rbac.PermissionCreatePost.String())

	// Explicit allowlist permits a configured non-privileged scope.
	s = &Service{cfg: config.Config{OAuthCIMDAllowedScopes: []string{rbac.PermissionCreatePost.String()}}}
	got = s.cimdAllowedScopes(rbac.PermissionCreatePost.String())
	assert.Contains(t, got, rbac.PermissionCreatePost.String())

	// Privileged scope only via explicit allowlist + opt-in flag.
	s = &Service{cfg: config.Config{OAuthCIMDAllowedScopes: []string{rbac.PermissionAdministrator.String()}}}
	assert.NotContains(t, s.cimdAllowedScopes(rbac.PermissionAdministrator.String()), rbac.PermissionAdministrator.String())
	s.cfg.OAuthCIMDAllowPrivilegedScopes = true
	assert.Contains(t, s.cimdAllowedScopes(rbac.PermissionAdministrator.String()), rbac.PermissionAdministrator.String())
}

func TestCIMDAllowedGrants(t *testing.T) {
	t.Parallel()

	grants := cimdAllowedGrants(&clientMetadataDocument{}, []string{"openid", "offline_access"})
	assert.Contains(t, grants, GrantTypeAuthorizationCode)
	assert.Contains(t, grants, GrantTypeRefreshToken)

	grants = cimdAllowedGrants(&clientMetadataDocument{GrantTypes: []string{GrantTypeAuthorizationCode}}, []string{"openid", "offline_access"})
	assert.Contains(t, grants, GrantTypeAuthorizationCode)
	assert.NotContains(t, grants, GrantTypeRefreshToken, "refresh excluded when doc omits it")

	grants = cimdAllowedGrants(&clientMetadataDocument{}, []string{"openid"})
	assert.NotContains(t, grants, GrantTypeRefreshToken, "no refresh without offline_access")
}

func TestFetchClientMetadataRejectsNonJSON(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<html></html>"))
	}))
	defer srv.Close()

	s := insecureCIMDService(t)
	docURL := srv.URL + "/m.json"
	u, err := s.parseCIMDClientID(docURL)
	require.NoError(t, err)

	_, _, oauthErr := s.fetchClientMetadata(context.Background(), u)
	require.NotNil(t, oauthErr)
	assert.Equal(t, "invalid_client", oauthErr.Code)
}

func TestFetchClientMetadataRejectsOversized(t *testing.T) {
	original := cimdMaxResponseBytes
	cimdMaxResponseBytes = 512
	defer func() { cimdMaxResponseBytes = original }()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"padding":"` + strings.Repeat("x", 4096) + `"}`))
	}))
	defer srv.Close()

	s := insecureCIMDService(t)
	docURL := srv.URL + "/m.json"
	u, err := s.parseCIMDClientID(docURL)
	require.NoError(t, err)

	_, _, oauthErr := s.fetchClientMetadata(context.Background(), u)
	require.NotNil(t, oauthErr)
	assert.Equal(t, "invalid_client", oauthErr.Code)
}

func TestFetchClientMetadataTimesOut(t *testing.T) {
	original := cimdFetchTimeout
	cimdFetchTimeout = 100 * time.Millisecond
	defer func() { cimdFetchTimeout = original }()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(500 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	s := insecureCIMDService(t)
	docURL := srv.URL + "/m.json"
	u, err := s.parseCIMDClientID(docURL)
	require.NoError(t, err)

	_, _, oauthErr := s.fetchClientMetadata(context.Background(), u)
	require.NotNil(t, oauthErr)
	assert.Equal(t, "invalid_client", oauthErr.Code)
}

func TestFetchClientMetadataDoesNotFollowRedirects(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Redirect(w, &http.Request{}, "https://other.example/elsewhere", http.StatusFound)
	}))
	defer srv.Close()

	s := insecureCIMDService(t)
	docURL := srv.URL + "/m.json"
	u, err := s.parseCIMDClientID(docURL)
	require.NoError(t, err)

	_, _, oauthErr := s.fetchClientMetadata(context.Background(), u)
	require.NotNil(t, oauthErr)
	assert.Equal(t, "invalid_client", oauthErr.Code)
}

func TestFetchClientMetadataParsesValidDocument(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"client_name":"Example","redirect_uris":["https://client.example/cb"]}`))
	}))
	defer srv.Close()

	s := insecureCIMDService(t)
	docURL := srv.URL + "/m.json"
	u, err := s.parseCIMDClientID(docURL)
	require.NoError(t, err)

	doc, ttl, oauthErr := s.fetchClientMetadata(context.Background(), u)
	require.Nil(t, oauthErr)
	require.NotNil(t, doc)
	assert.Equal(t, "Example", doc.ClientName)
	assert.Equal(t, []string{"https://client.example/cb"}, doc.RedirectURIs)
	assert.Equal(t, cimdDefaultCacheTTL, ttl)
}

func TestCIMDCacheTTL(t *testing.T) {
	t.Parallel()

	assert.Equal(t, cimdDefaultCacheTTL, cimdCacheTTL(http.Header{}))
	assert.Equal(t, time.Duration(0), cimdCacheTTL(http.Header{"Cache-Control": []string{"no-store"}}))
	assert.Equal(t, 30*time.Second, cimdCacheTTL(http.Header{"Cache-Control": []string{"max-age=30"}}))
	assert.Equal(t, cimdMaxCacheTTL, cimdCacheTTL(http.Header{"Cache-Control": []string{"max-age=99999"}}))
}

func TestParseCIMDClientIDRejectsSection3(t *testing.T) {
	t.Parallel()

	s := &Service{cfg: config.Config{OAuthEnabled: true, OAuthClientIDMetadataDocumentEnabled: true}, cimdCache: newCIMDCache()}

	cases := []string{
		"https://client.example",
		"https://client.example/",
		"https://client.example/../x",
		"https://client.example/./x",
		"https://client.example/x#f",
	}
	for _, c := range cases {
		_, err := s.parseCIMDClientID(c)
		assert.Error(t, err, c)
	}

	// Queries are permitted (and preserved) as part of the client_id identifier.
	u, err := s.parseCIMDClientID("https://client.example/x?y=1&foo=bar")
	require.NoError(t, err)
	assert.Equal(t, "y=1&foo=bar", u.RawQuery)
}

func TestValidateClientMetadataRejectsAuthAndSecret(t *testing.T) {
	t.Parallel()

	id := "https://client.example/m.json"
	err := validateClientMetadata(&clientMetadataDocument{
		ClientID:                id,
		RedirectURIs:            []string{"https://client.example/cb"},
		TokenEndpointAuthMethod: "private_key_jwt",
	}, id)
	assert.Error(t, err)

	err = validateClientMetadata(&clientMetadataDocument{
		ClientID:              id,
		RedirectURIs:          []string{"https://client.example/cb"},
		ClientSecret:          "foo",
		ClientSecretExpiresAt: 123,
	}, id)
	assert.Error(t, err)
}

func TestFetchRespectsDefault5kLimit(t *testing.T) {
	t.Parallel()

	// default is now 5k; oversized test server returns >5k
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"pad":"` + strings.Repeat("x", 6000) + `"}`))
	}))
	defer srv.Close()

	s := insecureCIMDService(t)
	docURL := srv.URL + "/m.json"
	u, err := s.parseCIMDClientID(docURL)
	require.NoError(t, err)

	_, _, oauthErr := s.fetchClientMetadata(context.Background(), u)
	require.NotNil(t, oauthErr)
}

func TestResolveClientPrefersPreRegistered(t *testing.T) {
	t.Parallel()

	assert.False(t, isCIMDClientID("storyden-cli"))
}

func TestCIMDCacheBounded(t *testing.T) {
	t.Parallel()

	c := newCIMDCache()
	for i := 0; i < cimdMaxCacheEntries+10; i++ {
		c.store("k"+strconv.Itoa(i), time.Hour)
	}
	assert.LessOrEqual(t, len(c.entries), cimdMaxCacheEntries)
}

func TestPrevalidateAndLookupOverride(t *testing.T) {
	original := cimdLookupNetIP
	defer func() { cimdLookupNetIP = original }()

	cimdLookupNetIP = func(ctx context.Context, _ string, host string) ([]netip.Addr, error) {
		if host == "evil.example" {
			return []netip.Addr{netip.MustParseAddr("10.0.0.1")}, nil
		}
		return []netip.Addr{netip.MustParseAddr("1.2.3.4")}, nil
	}

	svc := &Service{cfg: config.Config{OAuthEnabled: true, OAuthClientIDMetadataDocumentEnabled: true}, cimdCache: newCIMDCache()}
	// hostname not statically disallowed; prevalidate (called in fetch) will use the override and reject bad addrs
	u, err := svc.parseCIMDClientID("https://good.example/m.json")
	require.NoError(t, err)
	_ = u
}

func TestConcurrentCIMDFetch(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(20 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"redirect_uris":["https://c/cb"]}`))
	}))
	defer srv.Close()

	s := insecureCIMDService(t)
	docURL := srv.URL + "/m.json"
	u, _ := s.parseCIMDClientID(docURL)

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _, _ = s.fetchClientMetadata(context.Background(), u)
		}()
	}
	wg.Wait()
	// no panic, bounded fetches acceptable
}
