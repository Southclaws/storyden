package oauthremote

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Southclaws/storyden/app/services/authentication/oauthremote/oauth_http_client"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProtectedResourceMetadataURLPreservesResourcePath(t *testing.T) {
	t.Parallel()

	resource, err := url.Parse("https://mcp.example.com/mcp")
	require.NoError(t, err)

	assert.Equal(
		t,
		"https://mcp.example.com/.well-known/oauth-protected-resource/mcp",
		protectedResourceMetadataURL(resource),
	)
}

func TestProtectedResourceMetadataURLUsesOriginForRootResource(t *testing.T) {
	t.Parallel()

	resource, err := url.Parse("https://mcp.example.com/")
	require.NoError(t, err)

	assert.Equal(
		t,
		"https://mcp.example.com/.well-known/oauth-protected-resource",
		protectedResourceMetadataURL(resource),
	)
}

func TestAuthorizationServerMetadataURLPreservesIssuerPath(t *testing.T) {
	t.Parallel()

	assert.Equal(
		t,
		"https://auth.example.com/.well-known/oauth-authorization-server/realms/storyden",
		authorizationServerMetadataURL("https://auth.example.com/realms/storyden"),
	)
}

func TestValidateRemoteOAuthURLRejectsPlainHTTPNonLoopback(t *testing.T) {
	t.Parallel()

	err := oauth_http_client.ValidateRemoteOAuthURL("http://mcp.example.com/token", "token endpoint")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must use https")
}

func TestValidateRemoteOAuthURLAllowsPlainHTTPLoopback(t *testing.T) {
	t.Parallel()

	require.NoError(t, oauth_http_client.ValidateRemoteOAuthURL("http://127.0.0.1:8080/token", "token endpoint"))
	require.NoError(t, oauth_http_client.ValidateRemoteOAuthURL("http://localhost:8080/token", "token endpoint"))
}

func TestHTTPClientRejectsInsecureRedirect(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://example.com/.well-known/oauth-authorization-server", http.StatusFound)
	}))
	t.Cleanup(server.Close)

	_, err := fetchJSON[AuthorizationServerMetadata](context.Background(), oauth_http_client.NewHTTPClient(time.Second), server.URL)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "redirect URL must use https")
}

func TestValidateAuthorizationServerMetadataRequiresS256WhenAdvertised(t *testing.T) {
	t.Parallel()

	err := validateAuthorizationServerMetadata(AuthorizationServerMetadata{
		Issuer:                        "https://auth.example.com",
		AuthorizationEndpoint:         "https://auth.example.com/authorize",
		TokenEndpoint:                 "https://auth.example.com/token",
		ResponseTypesSupported:        []string{"code"},
		GrantTypesSupported:           []string{"authorization_code"},
		CodeChallengeMethodsSupported: []string{"plain"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "S256")
}

func TestDefaultRedirectURIUsesFrontendRemoteCallback(t *testing.T) {
	t.Parallel()

	apiURL, err := url.Parse("https://api.example.com")
	require.NoError(t, err)
	webURL, err := url.Parse("https://storyden.example.com")
	require.NoError(t, err)

	service := Service{config: config.Config{
		PublicAPIAddress: *apiURL,
		PublicWebAddress: *webURL,
	}}

	redirectURI, err := url.Parse(service.DefaultRedirectURI())
	require.NoError(t, err)

	assert.Equal(t, "https", redirectURI.Scheme)
	assert.Equal(t, "storyden.example.com", redirectURI.Host)
	assert.Equal(t, "/oauth/remote/callback", redirectURI.Path)
	assert.Empty(t, redirectURI.RawQuery)

	authURL := url.URL{
		Scheme: "https",
		Host:   "auth.example.com",
		Path:   "/authorize",
	}
	query := authURL.Query()
	query.Set("redirect_uri", service.DefaultRedirectURI())
	authURL.RawQuery = query.Encode()

	assert.Contains(t, authURL.String(), "redirect_uri=https%3A%2F%2Fstoryden.example.com%2Foauth%2Fremote%2Fcallback")
}
