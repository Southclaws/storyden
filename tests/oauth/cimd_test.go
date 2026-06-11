package oauth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/role/role_assign"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

type cimdMetadataServer struct {
	server   *httptest.Server
	fetches  int64
	document map[string]any
}

func newCIMDMetadataServer(t *testing.T) *cimdMetadataServer {
	t.Helper()

	cm := &cimdMetadataServer{}
	cm.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt64(&cm.fetches, 1)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(cm.document)
	}))
	t.Cleanup(cm.server.Close)

	// The CIMD client_id is the metadata document URL itself (must have path per spec).
	docURL := cm.server.URL + "/metadata.json"
	cm.document = map[string]any{
		"client_id":     docURL,
		"client_name":   "ChatGPT MCP Connector",
		"redirect_uris": []string{"https://client.example/callback"},
		"grant_types":   []string{"authorization_code", "refresh_token"},
		"scope":         "openid profile " + rbac.PermissionReadPublishedThreads.String(),
	}

	return cm
}

func (cm *cimdMetadataServer) clientID() string { return cm.server.URL + "/metadata.json" }

func (cm *cimdMetadataServer) fetchCount() int64 { return atomic.LoadInt64(&cm.fetches) }

func cimdConfig(t *testing.T) *config.Config {
	t.Helper()
	cfg := oauthConfig(t)
	cfg.OAuthClientIDMetadataDocumentEnabled = true
	cfg.OAuthCIMDAllowInsecureFetch = true // allow httptest loopback servers
	return cfg
}

func TestOAuthCIMDDiscoveryAdvertisesSupport(t *testing.T) {
	t.Parallel()

	integration.Test(t, cimdConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			for _, path := range []string{"/.well-known/openid-configuration", "/.well-known/oauth-authorization-server"} {
				t.Run(path, func(t *testing.T) {
					a := assert.New(t)
					r := require.New(t)

					req, err := http.NewRequestWithContext(root, http.MethodGet, ts.URL+path, nil)
					r.NoError(err)
					resp, err := http.DefaultClient.Do(req)
					r.NoError(err)
					defer resp.Body.Close()

					a.Equal(http.StatusOK, resp.StatusCode)

					var body map[string]any
					r.NoError(json.NewDecoder(resp.Body).Decode(&body))
					a.Equal(true, body["client_id_metadata_document_supported"])
				})
			}
		}))
	}))
}

func TestOAuthCIMDAuthorizationFlow(t *testing.T) {
	t.Parallel()

	integration.Test(t, cimdConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		assignments *role_assign.Assignment,
		roles *role_repo.Repository,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			memberCtx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			grantOAuthClientUse(t, root, roles, assignments, member.ID, rbac.PermissionReadPublishedThreads)
			memberSession := sh.WithSession(memberCtx)

			t.Run("valid_cimd_client_reaches_consent", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				cm := newCIMDMetadataServer(t)

				location := authorizeRedirect(t, root, ts, memberSession, authorizeRequest{
					ClientID:            cm.clientID(),
					RedirectURI:         "https://client.example/callback",
					Scope:               "openid profile " + rbac.PermissionReadPublishedThreads.String(),
					State:               "state-" + uuid.NewString(),
					CodeChallenge:       codeChallenge(strings.Repeat("a", 43)),
					CodeChallengeMethod: "S256",
				})

				consentURL, err := url.Parse(location)
				r.NoError(err)
				a.Equal("/oauth/authorize/consent", consentURL.Path)
				requestID := consentURL.Query().Get("request_id")
				r.NotEmpty(requestID)

				consent := tests.AssertRequest(cl.OAuthAuthoriseConsentWithResponse(root, &openapi.OAuthAuthoriseConsentParams{
					RequestId: (*openapi.OAuthAuthorizationRequestIDQuery)(&requestID),
				}, memberSession))(t, http.StatusOK)
				r.NotNil(consent.JSON200)
				a.Equal(cm.clientID(), consent.JSON200.ClientId)
				a.Contains(consent.JSON200.GrantedScopes, rbac.PermissionReadPublishedThreads.String())
			})

			t.Run("redirect_uri_not_in_metadata_is_rejected", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				cm := newCIMDMetadataServer(t)

				resp := authorizeHTTPResponse(t, root, ts, memberSession, authorizeRequest{
					ClientID:            cm.clientID(),
					RedirectURI:         "https://attacker.example/callback",
					Scope:               "openid",
					State:               "state-" + uuid.NewString(),
					CodeChallenge:       codeChallenge(strings.Repeat("b", 43)),
					CodeChallengeMethod: "S256",
				})
				defer resp.Body.Close()

				a.Equal(http.StatusBadRequest, resp.StatusCode)
				a.Empty(resp.Header.Get("Location"))
				r.NotNil(resp.Body)
			})

			t.Run("unknown_scope_is_rejected", func(t *testing.T) {
				a := assert.New(t)

				cm := newCIMDMetadataServer(t)

				resp := authorizeHTTPResponse(t, root, ts, memberSession, authorizeRequest{
					ClientID:            cm.clientID(),
					RedirectURI:         "https://client.example/callback",
					Scope:               "openid NONEXISTENT_SCOPE",
					State:               "state-" + uuid.NewString(),
					CodeChallenge:       codeChallenge(strings.Repeat("c", 43)),
					CodeChallengeMethod: "S256",
				})
				defer resp.Body.Close()

				a.Equal(http.StatusBadRequest, resp.StatusCode)
				a.Empty(resp.Header.Get("Location"))
			})

			t.Run("admin_scope_is_rejected_by_default", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				grantOAuthClientUse(t, root, roles, assignments, admin.ID, rbac.PermissionAdministrator)
				adminSession := sh.WithSession(adminCtx)

				cm := newCIMDMetadataServer(t)
				cm.document["scope"] = "openid " + rbac.PermissionAdministrator.String()

				location := authorizeRedirect(t, root, ts, adminSession, authorizeRequest{
					ClientID:            cm.clientID(),
					RedirectURI:         "https://client.example/callback",
					Scope:               "openid " + rbac.PermissionAdministrator.String(),
					State:               "state-" + uuid.NewString(),
					CodeChallenge:       codeChallenge(strings.Repeat("d", 43)),
					CodeChallengeMethod: "S256",
				})
				consentURL, err := url.Parse(location)
				r.NoError(err)
				requestID := consentURL.Query().Get("request_id")
				r.NotEmpty(requestID)

				consent := tests.AssertRequest(cl.OAuthAuthoriseConsentWithResponse(root, &openapi.OAuthAuthoriseConsentParams{
					RequestId: (*openapi.OAuthAuthorizationRequestIDQuery)(&requestID),
				}, adminSession))(t, http.StatusOK)
				r.NotNil(consent.JSON200)
				a.NotContains(consent.JSON200.GrantedScopes, rbac.PermissionAdministrator.String())
			})

			t.Run("cached_metadata_is_reused", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				cm := newCIMDMetadataServer(t)

				for i := 0; i < 3; i++ {
					location := authorizeRedirect(t, root, ts, memberSession, authorizeRequest{
						ClientID:            cm.clientID(),
						RedirectURI:         "https://client.example/callback",
						Scope:               "openid",
						State:               "state-" + uuid.NewString(),
						CodeChallenge:       codeChallenge(strings.Repeat("e", 43)),
						CodeChallengeMethod: "S256",
					})
					consentURL, err := url.Parse(location)
					r.NoError(err)
					r.NotEmpty(consentURL.Query().Get("request_id"))
				}

				a.Equal(int64(1), cm.fetchCount(), "metadata document should be fetched once and cached")
			})

			t.Run("cimd_client_id_with_query_string_and_none_auth_hint_full_flow_to_token", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				cm := newCIMDMetadataServer(t)
				// Simulate production usage (ChatGPT-style) where the client_id URL
				// itself carries a query hint for the auth method the client will use.
				queriedID := cm.server.URL + "/metadata.json?token_endpoint_auth_method=none"
				cm.document["client_id"] = queriedID // doc must declare the exact (full) client_id for match
				cm.document["scope"] = "openid offline_access " + rbac.PermissionReadPublishedThreads.String()

				// Use a code_challenge that we can compute verifier for.
				verifier := strings.Repeat("q", 43)
				challenge := codeChallenge(verifier)

				location := authorizeRedirect(t, root, ts, memberSession, authorizeRequest{
					ClientID:            queriedID,
					RedirectURI:         "https://client.example/callback",
					Scope:               "openid offline_access " + rbac.PermissionReadPublishedThreads.String(),
					State:               "state-q-" + uuid.NewString(),
					CodeChallenge:       challenge,
					CodeChallengeMethod: "S256",
				})

				consentURL, err := url.Parse(location)
				r.NoError(err)
				requestID := consentURL.Query().Get("request_id")
				r.NotEmpty(requestID)

				// Get consent form.
				consent := tests.AssertRequest(cl.OAuthAuthoriseConsentWithResponse(root, &openapi.OAuthAuthoriseConsentParams{
					RequestId: (*openapi.OAuthAuthorizationRequestIDQuery)(&requestID),
				}, memberSession))(t, http.StatusOK)
				r.NotNil(consent.JSON200)

				// Approve.
				submit := tests.AssertRequest(cl.OAuthAuthoriseConsentSubmitWithResponse(root, openapi.OAuthAuthoriseConsentSubmitJSONRequestBody{
					RequestId: requestID,
					Decision:  openapi.OAuthAuthoriseDecisionApprove,
				}, memberSession))(t, http.StatusOK)
				r.NotNil(submit.JSON200)
				a.Equal(openapi.OAuthAuthoriseConsentResultStatusApproved, submit.JSON200.Status)

				codeRedirect, err := url.Parse(submit.JSON200.Location)
				r.NoError(err)
				code := codeRedirect.Query().Get("code")
				r.NotEmpty(code)

				// Full token exchange for this public CIMD client (PKCE only, no secret).
				tok, err := oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:    "authorization_code",
					ClientId:     queriedID,
					Code:         &code,
					RedirectUri:  ptr("https://client.example/callback"),
					CodeVerifier: &verifier,
				})
				r.NoError(err)
				r.NotNil(tok)
				a.Equal(http.StatusOK, tok.StatusCode())
				r.NotNil(tok.JSON200)
				r.NotEmpty(tok.JSON200.AccessToken)
				r.NotEmpty(tok.JSON200.RefreshToken) // because doc requests offline_access + refresh grant
			})

			t.Run("cimd_metadata_updates_are_observed_after_cache_expiry", func(t *testing.T) {
				a := assert.New(t)

				cm := newCIMDMetadataServer(t)
				// Force re-fetch on next resolve by making the document response indicate no caching.
				// (Our cimdCacheTTL treats no-cache / max-age=0 as ttl=0 which invalidates.)
				origHandler := cm.server.Config.Handler
				cm.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					atomic.AddInt64(&cm.fetches, 1)
					w.Header().Set("Content-Type", "application/json")
					w.Header().Set("Cache-Control", "no-cache, no-store")
					_ = json.NewEncoder(w).Encode(cm.document)
				})
				t.Cleanup(func() { cm.server.Config.Handler = origHandler })

				id := cm.clientID()

				// First use: initial fetch + create snapshot.
				_ = authorizeHTTPResponse(t, root, ts, memberSession, authorizeRequest{
					ClientID:            id,
					RedirectURI:         "https://client.example/callback",
					Scope:               "openid",
					State:               "state-f1-" + uuid.NewString(),
					CodeChallenge:       codeChallenge(strings.Repeat("1", 43)),
					CodeChallengeMethod: "S256",
				}).Body.Close()

				initialFetches := cm.fetchCount()
				a.GreaterOrEqual(initialFetches, int64(1))

				// Mutate the published metadata: add an extra redirect_uri.
				// (Also keep the original so existing tests don't break.)
				cm.document["redirect_uris"] = []string{
					"https://client.example/callback",
					"https://client.example/extra-cb",
				}

				// Second authorize using the *new* redirect should succeed and trigger a re-fetch
				// because previous response forced ttl=0 / no-cache.
				resp2 := authorizeHTTPResponse(t, root, ts, memberSession, authorizeRequest{
					ClientID:            id,
					RedirectURI:         "https://client.example/extra-cb",
					Scope:               "openid",
					State:               "state-f2-" + uuid.NewString(),
					CodeChallenge:       codeChallenge(strings.Repeat("2", 43)),
					CodeChallengeMethod: "S256",
				})
				defer resp2.Body.Close()
				// It should have reached consent (200 would be after? but for redirect case we get 302 with location to consent or error).
				// A 302 to the consent page (or success) means the client was accepted with the *updated* redirect.
				a.True(resp2.StatusCode == http.StatusFound || resp2.StatusCode == http.StatusOK,
					"updated redirect from mutated CIMD doc should be accepted after re-fetch")

				a.Greater(cm.fetchCount(), initialFetches, "a re-fetch should have occurred due to cache expiry / no-cache header")
			})
		}))
	}))
}

func TestOAuthCIMDRejectsPrivateMetadataHostByDefault(t *testing.T) {
	t.Parallel()

	// Default config: OAuthCIMDAllowInsecureFetch is false, so a loopback
	// httptest server must be rejected by the SSRF protections.
	integration.Test(t, oauthConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		assignments *role_assign.Assignment,
		roles *role_repo.Repository,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			a := assert.New(t)

			memberCtx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			grantOAuthClientUse(t, root, roles, assignments, member.ID, rbac.PermissionReadPublishedThreads)
			memberSession := sh.WithSession(memberCtx)

			cm := newCIMDMetadataServer(t)

			resp := authorizeHTTPResponse(t, root, ts, memberSession, authorizeRequest{
				ClientID:            cm.clientID(), // http loopback URL
				RedirectURI:         "https://client.example/callback",
				Scope:               "openid",
				State:               "state-" + uuid.NewString(),
				CodeChallenge:       codeChallenge(strings.Repeat("f", 43)),
				CodeChallengeMethod: "S256",
			})
			defer resp.Body.Close()

			a.Equal(http.StatusBadRequest, resp.StatusCode)
			a.Empty(resp.Header.Get("Location"))
			a.Equal(int64(0), cm.fetchCount(), "metadata host must be rejected before any fetch")
		}))
	}))
}

func TestOAuthCIMDRejectsSection3ClientID(t *testing.T) {
	t.Parallel()

	integration.Test(t, cimdConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		assignments *role_assign.Assignment,
		roles *role_repo.Repository,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			a := assert.New(t)

			memberCtx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			grantOAuthClientUse(t, root, roles, assignments, member.ID, rbac.PermissionReadPublishedThreads)
			memberSession := sh.WithSession(memberCtx)

			badIDs := []string{
				"https://client.example",
				"https://client.example/",
				"https://client.example/../x.json",
				"https://client.example/x?evil=1",
			}
			for _, bid := range badIDs {
				resp := authorizeHTTPResponse(t, root, ts, memberSession, authorizeRequest{
					ClientID:            bid,
					RedirectURI:         "https://client.example/callback",
					Scope:               "openid",
					State:               "state-" + uuid.NewString(),
					CodeChallenge:       codeChallenge(strings.Repeat("1", 43)),
					CodeChallengeMethod: "S256",
				})
				defer resp.Body.Close()
				a.Equal(http.StatusBadRequest, resp.StatusCode, bid)
			}
		}))
	}))
}

func TestOAuthCIMDRejectsBadMetadataDoc(t *testing.T) {
	t.Parallel()

	integration.Test(t, cimdConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		assignments *role_assign.Assignment,
		roles *role_repo.Repository,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			a := assert.New(t)

			memberCtx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			grantOAuthClientUse(t, root, roles, assignments, member.ID, rbac.PermissionReadPublishedThreads)
			memberSession := sh.WithSession(memberCtx)

			// doc with forbidden auth method
			cm := newCIMDMetadataServer(t)
			cm.document["token_endpoint_auth_method"] = "client_secret_basic"

			resp := authorizeHTTPResponse(t, root, ts, memberSession, authorizeRequest{
				ClientID:            cm.clientID(),
				RedirectURI:         "https://client.example/callback",
				Scope:               "openid",
				State:               "state-" + uuid.NewString(),
				CodeChallenge:       codeChallenge(strings.Repeat("2", 43)),
				CodeChallengeMethod: "S256",
			})
			defer resp.Body.Close()
			a.Equal(http.StatusBadRequest, resp.StatusCode)

			// doc with fragment in redirect
			cm2 := newCIMDMetadataServer(t)
			cm2.document["redirect_uris"] = []string{"https://client.example/cb#f"}

			resp2 := authorizeHTTPResponse(t, root, ts, memberSession, authorizeRequest{
				ClientID:            cm2.clientID(),
				RedirectURI:         "https://client.example/cb#f",
				Scope:               "openid",
				State:               "state-" + uuid.NewString(),
				CodeChallenge:       codeChallenge(strings.Repeat("3", 43)),
				CodeChallengeMethod: "S256",
			})
			defer resp2.Body.Close()
			a.Equal(http.StatusBadRequest, resp2.StatusCode)
		}))
	}))
}
