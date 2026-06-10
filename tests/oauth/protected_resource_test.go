package oauth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

func TestOAuthProtectedResourceMetadata(t *testing.T) {
	t.Parallel()

	cfg := oauthConfig(t)
	publicAPIAddress, err := url.Parse("http://localhost:8000")
	require.NoError(t, err)
	cfg.PublicAPIAddress = *publicAPIAddress

	integration.Test(t, cfg, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			t.Run("root_resource", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				req, err := http.NewRequestWithContext(root, http.MethodGet, ts.URL+"/.well-known/oauth-protected-resource", nil)
				r.NoError(err)

				resp, err := http.DefaultClient.Do(req)
				r.NoError(err)
				defer resp.Body.Close()

				r.Equal(http.StatusOK, resp.StatusCode)
				a.Equal("application/json", resp.Header.Get("Content-Type"))

				var meta struct {
					Resource               string   `json:"resource"`
					AuthorizationServers   []string `json:"authorization_servers"`
					BearerMethodsSupported []string `json:"bearer_methods_supported"`
					ScopesSupported        []string `json:"scopes_supported"`
				}
				r.NoError(json.NewDecoder(resp.Body).Decode(&meta))

				a.Equal("http://localhost:8000", meta.Resource)
				r.Len(meta.AuthorizationServers, 1)
				a.Equal("http://localhost:8000", meta.AuthorizationServers[0])
				a.Equal([]string{"header"}, meta.BearerMethodsSupported)
				a.Empty(meta.ScopesSupported, "root resource should not advertise scopes")
			})

			t.Run("api_resource", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				req, err := http.NewRequestWithContext(root, http.MethodGet, ts.URL+"/.well-known/oauth-protected-resource/api", nil)
				r.NoError(err)

				resp, err := http.DefaultClient.Do(req)
				r.NoError(err)
				defer resp.Body.Close()

				r.Equal(http.StatusOK, resp.StatusCode)

				var meta struct {
					Resource               string   `json:"resource"`
					AuthorizationServers   []string `json:"authorization_servers"`
					BearerMethodsSupported []string `json:"bearer_methods_supported"`
					ScopesSupported        []string `json:"scopes_supported"`
				}
				r.NoError(json.NewDecoder(resp.Body).Decode(&meta))

				a.Equal("http://localhost:8000/api", meta.Resource)
				r.Len(meta.AuthorizationServers, 1)
				a.Equal("http://localhost:8000", meta.AuthorizationServers[0])
				a.Equal([]string{"header"}, meta.BearerMethodsSupported)
				a.Contains(meta.ScopesSupported, "openid")
				a.Contains(meta.ScopesSupported, "profile")
				a.Contains(meta.ScopesSupported, "email")
				a.Contains(meta.ScopesSupported, "offline_access")
				a.NotContains(meta.ScopesSupported, "ADMINISTRATOR", "admin scope must not be exposed")
				a.NotContains(meta.ScopesSupported, "MANAGE_SETTINGS", "admin scope must not be exposed")
				a.NotContains(meta.ScopesSupported, "MANAGE_ACCOUNTS", "admin scope must not be exposed")
			})

			t.Run("mcp_sse_resource", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				req, err := http.NewRequestWithContext(root, http.MethodGet, ts.URL+"/.well-known/oauth-protected-resource/mcp/sse", nil)
				r.NoError(err)

				resp, err := http.DefaultClient.Do(req)
				r.NoError(err)
				defer resp.Body.Close()

				r.Equal(http.StatusOK, resp.StatusCode)

				var meta struct {
					Resource               string   `json:"resource"`
					AuthorizationServers   []string `json:"authorization_servers"`
					BearerMethodsSupported []string `json:"bearer_methods_supported"`
					ScopesSupported        []string `json:"scopes_supported"`
				}
				r.NoError(json.NewDecoder(resp.Body).Decode(&meta))

				a.Equal("http://localhost:8000/mcp/sse", meta.Resource)
				r.Len(meta.AuthorizationServers, 1)
				a.Equal("http://localhost:8000", meta.AuthorizationServers[0])
				a.Equal([]string{"header"}, meta.BearerMethodsSupported)
				a.Contains(meta.ScopesSupported, "openid")
				a.NotContains(meta.ScopesSupported, "ADMINISTRATOR", "admin scope must not be exposed")
			})

			t.Run("no_auth_required", func(t *testing.T) {
				r := require.New(t)

				for _, path := range []string{
					"/.well-known/oauth-protected-resource",
					"/.well-known/oauth-protected-resource/api",
					"/.well-known/oauth-protected-resource/mcp/sse",
				} {
					req, err := http.NewRequestWithContext(root, http.MethodGet, ts.URL+path, nil)
					r.NoError(err)

					resp, err := http.DefaultClient.Do(req)
					r.NoError(err)
					resp.Body.Close()

					r.Equal(http.StatusOK, resp.StatusCode, "path %s must be publicly accessible", path)
				}
			})
		}))
	}))
}

func TestOAuthProtectedResourceMetadataWithAPIPath(t *testing.T) {
	t.Parallel()

	cfg := oauthConfig(t)
	publicAPIAddress, err := url.Parse("http://localhost:8000/api")
	require.NoError(t, err)
	cfg.PublicAPIAddress = *publicAPIAddress

	integration.Test(t, cfg, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			t.Run("issuer_is_normalised_without_api_suffix", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				req, err := http.NewRequestWithContext(root, http.MethodGet, ts.URL+"/.well-known/oauth-protected-resource", nil)
				r.NoError(err)

				resp, err := http.DefaultClient.Do(req)
				r.NoError(err)
				defer resp.Body.Close()

				r.Equal(http.StatusOK, resp.StatusCode)

				var meta struct {
					Resource             string   `json:"resource"`
					AuthorizationServers []string `json:"authorization_servers"`
				}
				r.NoError(json.NewDecoder(resp.Body).Decode(&meta))

				a.Equal("http://localhost:8000", meta.Resource)
				r.Len(meta.AuthorizationServers, 1)
				a.Equal("http://localhost:8000", meta.AuthorizationServers[0])
			})
		}))
	}))
}
