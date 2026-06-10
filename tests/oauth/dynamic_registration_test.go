package oauth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/role/role_assign"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	oauthresource "github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestOAuthDynamicClientRegistration(t *testing.T) {
	t.Parallel()

	integration.Test(t, oauthConfig(t), e2e.Setup(), fx.Invoke(func(
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
			memberCtx, member := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			grantOAuthClientUse(t, root, roles, assignments, member.ID)
			memberSession := sh.WithSession(memberCtx)

			t.Run("public_pkce_client_registration_succeeds", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				redirectURI := "https://chatgpt.com/connector_platform_oauth_redirect"
				resp := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					ClientName:              ptr("ChatGPT MCP Connector"),
					RedirectUris:            &[]string{redirectURI},
					TokenEndpointAuthMethod: ptr("none"),
					Scope:                   ptr("openid profile email offline_access"),
				}, memberSession))(t, http.StatusCreated)
				r.NotNil(resp.JSON201)

				a.True(strings.HasPrefix(resp.JSON201.ClientId, oauthresource.OAuthAccessKeyPrefix))
				a.Nil(resp.JSON201.ClientSecret, "public clients must not receive a client secret")
				a.Positive(resp.JSON201.ClientIdIssuedAt)
				a.Equal(int64(0), resp.JSON201.ClientSecretExpiresAt)
				a.Equal("none", resp.JSON201.TokenEndpointAuthMethod)
				a.Equal([]string{redirectURI}, resp.JSON201.RedirectUris)
				a.Contains(resp.JSON201.GrantTypes, oauthGrantAuthorizationCode)
				a.Contains(resp.JSON201.GrantTypes, oauthGrantRefreshToken)
				a.Equal([]string{"code"}, resp.JSON201.ResponseTypes)
				r.NotNil(resp.JSON201.ClientName)
				a.Equal("ChatGPT MCP Connector", *resp.JSON201.ClientName)
			})

			t.Run("public_pkce_client_registration_unauthenticated_succeeds", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				// RFC 7591 requires DCR to be publicly accessible (unauthenticated)
				// This test verifies registration works without a session/auth token
				redirectURI := "https://example.com/oauth/callback"
				resp := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					ClientName:              ptr("Unauthenticated Public Client"),
					RedirectUris:            &[]string{redirectURI},
					TokenEndpointAuthMethod: ptr("none"),
					Scope:                   ptr("openid profile"),
				}))(t, http.StatusCreated) // No session passed - unauthenticated
				r.NotNil(resp.JSON201)

				a.True(strings.HasPrefix(resp.JSON201.ClientId, oauthresource.OAuthAccessKeyPrefix))
				a.Nil(resp.JSON201.ClientSecret, "public clients must not receive a client secret")
				a.Positive(resp.JSON201.ClientIdIssuedAt)
				a.Equal(int64(0), resp.JSON201.ClientSecretExpiresAt)
				a.Equal("none", resp.JSON201.TokenEndpointAuthMethod)
				a.Equal([]string{redirectURI}, resp.JSON201.RedirectUris)
				a.Contains(resp.JSON201.GrantTypes, oauthGrantAuthorizationCode)
				a.Contains(resp.JSON201.GrantTypes, oauthGrantRefreshToken)
				a.Equal([]string{"code"}, resp.JSON201.ResponseTypes)
				r.NotNil(resp.JSON201.ClientName)
				a.Equal("Unauthenticated Public Client", *resp.JSON201.ClientName)
			})

			t.Run("confidential_client_registration_returns_secret", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				resp := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					ClientName:              ptr("Confidential App"),
					RedirectUris:            &[]string{"https://app.example/callback"},
					TokenEndpointAuthMethod: ptr("client_secret_post"),
				}, memberSession))(t, http.StatusCreated)
				r.NotNil(resp.JSON201)

				r.NotNil(resp.JSON201.ClientSecret)
				a.True(strings.HasPrefix(*resp.JSON201.ClientSecret, oauthresource.OAuthAccessSecretPrefix))
				a.Equal("client_secret_post", resp.JSON201.TokenEndpointAuthMethod)
				a.Equal(int64(0), resp.JSON201.ClientSecretExpiresAt)
			})

			t.Run("rejects_client_secret_basic_unsupported", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				// client_secret_basic requires HTTP Basic auth which the token endpoint doesn't support
				resp := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					ClientName:              ptr("Basic Auth Client"),
					RedirectUris:            &[]string{"https://app.example/callback"},
					TokenEndpointAuthMethod: ptr("client_secret_basic"),
				}, memberSession))(t, http.StatusBadRequest)
				r.NotNil(resp.JSON400)
				a.Equal("invalid_client_metadata", resp.JSON400.Error)
			})

			t.Run("rejects_invalid_redirect_uri", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				for _, bad := range []string{
					"http://not-loopback.example/callback",
					"https://app.example/callback#fragment",
					"https://app.example/*",
					"/relative/path",
				} {
					resp := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
						ClientName:              ptr("Bad Redirect"),
						RedirectUris:            &[]string{bad},
						TokenEndpointAuthMethod: ptr("none"),
					}, memberSession))(t, http.StatusBadRequest)
					r.NotNil(resp.JSON400, "expected rejection for %q", bad)
					a.Equal("invalid_redirect_uri", resp.JSON400.Error, "redirect uri %q", bad)
				}
			})

			t.Run("allows_loopback_http_redirect_uri", func(t *testing.T) {
				r := require.New(t)

				resp := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					ClientName:              ptr("Native Loopback"),
					RedirectUris:            &[]string{"http://127.0.0.1:51000/callback"},
					TokenEndpointAuthMethod: ptr("none"),
				}, memberSession))(t, http.StatusCreated)
				r.NotNil(resp.JSON201)
			})

			t.Run("rejects_unknown_scopes", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				resp := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					ClientName:              ptr("Invalid Scope"),
					RedirectUris:            &[]string{"https://app.example/callback"},
					TokenEndpointAuthMethod: ptr("none"),
					Scope:                   ptr("openid TOTALLY_INVALID_SCOPE_12345"),
				}, memberSession))(t, http.StatusBadRequest)
				r.NotNil(resp.JSON400)
				a.Equal("invalid_client_metadata", resp.JSON400.Error)
			})

			t.Run("rejects_administrator_scope", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				// ADMINISTRATOR scope is too powerful for DCR clients
				resp := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					ClientName:              ptr("Admin Seeking Client"),
					RedirectUris:            &[]string{"https://app.example/callback"},
					TokenEndpointAuthMethod: ptr("none"),
					Scope:                   ptr("openid ADMINISTRATOR"),
				}, memberSession))(t, http.StatusBadRequest)
				r.NotNil(resp.JSON400)
				a.Equal("invalid_client_metadata", resp.JSON400.Error, "must reject ADMINISTRATOR scope for DCR")
			})

			t.Run("rejects_unsupported_grant_types", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				resp := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					ClientName:              ptr("Service"),
					RedirectUris:            &[]string{"https://app.example/callback"},
					TokenEndpointAuthMethod: ptr("client_secret_post"),
					GrantTypes:              &[]string{oauthGrantClientCredentials},
				}, memberSession))(t, http.StatusBadRequest)
				r.NotNil(resp.JSON400)
				a.Equal("invalid_client_metadata", resp.JSON400.Error)
			})

			t.Run("metadata_advertises_registration_endpoint", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				for _, path := range []string{
					"/.well-known/oauth-authorization-server",
					"/.well-known/openid-configuration",
				} {
					req, err := http.NewRequestWithContext(root, http.MethodGet, ts.URL+path, nil)
					r.NoError(err)

					httpResp, err := http.DefaultClient.Do(req)
					r.NoError(err)

					var body map[string]any
					r.NoError(json.NewDecoder(httpResp.Body).Decode(&body))
					httpResp.Body.Close()

					a.Equal(http.StatusOK, httpResp.StatusCode, path)
					endpoint, ok := body["registration_endpoint"].(string)
					r.True(ok, "registration_endpoint missing from %s", path)
					a.True(strings.HasSuffix(endpoint, "/oauth/register"), "%s: %q", path, endpoint)
				}
			})

			t.Run("registered_client_can_start_authorize_without_unauthorized_client", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				redirectURI := "https://client.example/callback"
				registered := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					ClientName:              ptr("Authorize Flow Client"),
					RedirectUris:            &[]string{redirectURI},
					TokenEndpointAuthMethod: ptr("none"),
					Scope:                   ptr("openid profile"),
				}, memberSession))(t, http.StatusCreated)
				r.NotNil(registered.JSON201)

				verifier := strings.Repeat("z", 43)
				location := authorizeRedirect(t, root, ts, memberSession, authorizeRequest{
					ClientID:            registered.JSON201.ClientId,
					RedirectURI:         redirectURI,
					Scope:               "openid profile",
					State:               "state-" + uuid.NewString(),
					CodeChallenge:       codeChallenge(verifier),
					CodeChallengeMethod: "S256",
				})

				consentURL, err := url.Parse(location)
				r.NoError(err)
				a.Equal("/oauth/authorize/consent", consentURL.Path)
				a.NotEmpty(consentURL.Query().Get("request_id"))
			})

			// Adversarial security tests
			t.Run("hostile_public_client_credentials", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				// Attempt to register public client with client_credentials grant
				resp := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					ClientName:              ptr("Malicious Public Client"),
					TokenEndpointAuthMethod: ptr("none"),
					GrantTypes:              &[]string{oauthGrantClientCredentials},
				}, memberSession))(t, http.StatusBadRequest)
				r.NotNil(resp.JSON400)
				a.Equal("invalid_client_metadata", resp.JSON400.Error, "must reject public client with client_credentials")
			})

			t.Run("hostile_refresh_token_without_authorization_code", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				// Attempt to register with refresh_token but no authorization_code
				resp := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					ClientName:              ptr("Refresh Only Client"),
					TokenEndpointAuthMethod: ptr("client_secret_post"),
					GrantTypes:              &[]string{oauthGrantRefreshToken},
				}, memberSession))(t, http.StatusBadRequest)
				r.NotNil(resp.JSON400)
				a.Equal("invalid_client_metadata", resp.JSON400.Error, "must reject refresh_token without authorization_code")
			})

			t.Run("hostile_response_type_grant_type_mismatch", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				// Attempt to register with "code" response type but no authorization_code grant
				resp := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					ClientName:              ptr("Inconsistent Client"),
					RedirectUris:            &[]string{"https://app.example/callback"},
					TokenEndpointAuthMethod: ptr("client_secret_post"),
					GrantTypes:              &[]string{oauthGrantClientCredentials},
					ResponseTypes:           &[]string{"code"},
				}, memberSession))(t, http.StatusBadRequest)
				r.NotNil(resp.JSON400)
				a.Equal("invalid_client_metadata", resp.JSON400.Error, "must reject code response type without authorization_code grant")
			})

			t.Run("hostile_authorization_code_without_code_response_type", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				// Attempt to register with authorization_code grant but "token" response type (invalid combo)
				resp := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					ClientName:              ptr("Invalid Response Type"),
					RedirectUris:            &[]string{"https://app.example/callback"},
					TokenEndpointAuthMethod: ptr("none"),
					GrantTypes:              &[]string{oauthGrantAuthorizationCode},
					ResponseTypes:           &[]string{"token"}, // invalid: token is for implicit grant
				}, memberSession))(t, http.StatusBadRequest)
				r.NotNil(resp.JSON400)
				a.Equal("invalid_client_metadata", resp.JSON400.Error, "must reject authorization_code grant with invalid response type")
			})

			t.Run("hostile_http_redirect_uri", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				// Attempt plain HTTP for non-loopback
				resp := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					ClientName:              ptr("Insecure Client"),
					RedirectUris:            &[]string{"http://evil.example.com/callback"},
					TokenEndpointAuthMethod: ptr("none"),
				}, memberSession))(t, http.StatusBadRequest)
				r.NotNil(resp.JSON400)
				a.Equal("invalid_redirect_uri", resp.JSON400.Error, "must reject non-HTTPS redirect URI")
			})

			t.Run("hostile_wildcard_redirect_uri", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				// Attempt wildcard redirect URI
				resp := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					ClientName:              ptr("Wildcard Client"),
					RedirectUris:            &[]string{"https://*.evil.com/callback"},
					TokenEndpointAuthMethod: ptr("none"),
				}, memberSession))(t, http.StatusBadRequest)
				r.NotNil(resp.JSON400)
				a.Equal("invalid_redirect_uri", resp.JSON400.Error, "must reject wildcard redirect URI")
			})

			t.Run("hostile_fragment_in_redirect_uri", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				// Attempt redirect URI with fragment
				resp := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					ClientName:              ptr("Fragment Client"),
					RedirectUris:            &[]string{"https://app.example/callback#evil"},
					TokenEndpointAuthMethod: ptr("none"),
				}, memberSession))(t, http.StatusBadRequest)
				r.NotNil(resp.JSON400)
				a.Equal("invalid_redirect_uri", resp.JSON400.Error, "must reject redirect URI with fragment")
			})

			t.Run("hostile_http_metadata_uri", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				// Attempt HTTP logo_uri (should be rejected)
				resp := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					ClientName:              ptr("HTTP Logo"),
					RedirectUris:            &[]string{"https://app.example/callback"},
					TokenEndpointAuthMethod: ptr("none"),
					LogoUri:                 ptr("http://evil.example/malware.png"),
				}, memberSession))(t, http.StatusBadRequest)
				r.NotNil(resp.JSON400)
				a.Equal("invalid_client_metadata", resp.JSON400.Error, "must reject HTTP metadata URI")
			})

			t.Run("hostile_relative_metadata_uri", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				// Attempt relative client_uri
				resp := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					ClientName:              ptr("Relative URI"),
					RedirectUris:            &[]string{"https://app.example/callback"},
					TokenEndpointAuthMethod: ptr("none"),
					ClientUri:               ptr("/relative/path"),
				}, memberSession))(t, http.StatusBadRequest)
				r.NotNil(resp.JSON400)
				a.Equal("invalid_client_metadata", resp.JSON400.Error, "must reject relative metadata URI")
			})

			t.Run("redirect_uri_deduplication", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				// Register with duplicate redirect URIs
				resp := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					ClientName: ptr("Dedupe Test"),
					RedirectUris: &[]string{
						"https://app.example/callback",
						"https://app.example/callback2",
						"https://app.example/callback", // duplicate
					},
					TokenEndpointAuthMethod: ptr("none"),
				}, memberSession))(t, http.StatusCreated)
				r.NotNil(resp.JSON201)

				// Should only have 2 unique URIs
				a.Len(resp.JSON201.RedirectUris, 2, "should deduplicate redirect URIs")
				a.Contains(resp.JSON201.RedirectUris, "https://app.example/callback")
				a.Contains(resp.JSON201.RedirectUris, "https://app.example/callback2")
			})

			t.Run("hostile_opaque_redirect_uri", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				// Attempt opaque URI like "https:callback" which has no hostname
				resp := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					ClientName:              ptr("Opaque URI Attack"),
					RedirectUris:            &[]string{"https:callback"},
					TokenEndpointAuthMethod: ptr("none"),
				}, memberSession))(t, http.StatusBadRequest)
				r.NotNil(resp.JSON400)
				a.Equal("invalid_redirect_uri", resp.JSON400.Error, "must reject opaque URI with no hostname")
			})

			t.Run("allows_ipv4_loopback_range", func(t *testing.T) {
				r := require.New(t)

				// Test 127.0.0.2 (should be allowed as loopback)
				resp := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					ClientName:              ptr("IPv4 Loopback"),
					RedirectUris:            &[]string{"http://127.0.0.2:8080/callback"},
					TokenEndpointAuthMethod: ptr("none"),
				}, memberSession))(t, http.StatusCreated)
				r.NotNil(resp.JSON201, "should allow 127.0.0.2 as loopback")
			})

			t.Run("allows_ipv6_loopback", func(t *testing.T) {
				r := require.New(t)

				// Test ::1 with brackets (standard IPv6 URL format)
				resp := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					ClientName:              ptr("IPv6 Loopback"),
					RedirectUris:            &[]string{"http://[::1]:8080/callback"},
					TokenEndpointAuthMethod: ptr("none"),
				}, memberSession))(t, http.StatusCreated)
				r.NotNil(resp.JSON201, "should allow ::1 as loopback")
			})

			t.Run("token_endpoint_auth_method_persisted", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				// Register client with client_secret_post
				resp := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					ClientName:              ptr("Auth Method Test"),
					RedirectUris:            &[]string{"https://app.example/callback"},
					TokenEndpointAuthMethod: ptr("client_secret_post"),
				}, memberSession))(t, http.StatusCreated)
				r.NotNil(resp.JSON201)

				// Verify it's returned correctly
				a.Equal("client_secret_post", resp.JSON201.TokenEndpointAuthMethod)
			})
		}))
	}))
}
