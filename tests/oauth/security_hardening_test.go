package oauth_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/role/role_assign"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	oauthresource "github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/app/resources/oauth/oauth_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestOAuthSecurityHardeningDeviceClientImpersonationMitigations(t *testing.T) {
	t.Parallel()

	integration.Test(t, oauthConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			t.Run("built_in_client_rejects_explicit_storyden_permission_scopes", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				resp := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(root, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
					ClientId: "storyden-cli",
					Scope:    ptr("openid profile offline_access ADMINISTRATOR"),
				}))(t, http.StatusBadRequest)
				r.NotNil(resp.JSON400)
				a.Equal("invalid_scope", resp.JSON400.Error)
			})

			t.Run("built_in_client_is_branded_storyden_on_consent", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				start := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(root, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
					ClientId: "storyden-cli",
					Scope:    ptr("openid profile offline_access"),
				}))(t, http.StatusOK)
				r.NotNil(start.JSON200)
				r.NotNil(start.JSON200.UserCode)

				consent := tests.AssertRequest(cl.OAuthDeviceConsentWithResponse(root, &openapi.OAuthDeviceConsentParams{
					UserCode: start.JSON200.UserCode,
				}, adminSession))(t, http.StatusOK)
				r.NotNil(consent.JSON200)
				a.Equal("Storyden", consent.JSON200.ClientName)
				a.True(consent.JSON200.InheritsUserPermissions)
			})
		}))
	}))
}

func TestOAuthSecurityHardeningGrantAllowListEnforcement(t *testing.T) {
	t.Parallel()

	integration.Test(t, oauthConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		ow *oauth_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			t.Run("confidential_client_cannot_start_device_authorization", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				clientID := "confidential-device-" + uuid.NewString()
				clientSecret := "secret-" + uuid.NewString()
				createClient(t, root, ow, admin.ID, clientID, oauthresource.ClientTypeConfidential, oauthresource.ScopePolicyExplicit, opt.New(clientSecretHash(t, clientSecret)), standardScopes(), []string{oauthGrantDeviceCode})

				resp := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(root, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
					ClientId: clientID,
					Scope:    ptr("openid profile offline_access"),
				}))(t, http.StatusBadRequest)
				r.NotNil(resp.JSON400)
				a.Equal("unauthorized_client", resp.JSON400.Error)
			})

			t.Run("offline_access_without_refresh_grant_does_not_issue_refresh_token", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				clientID := "device-no-refresh-" + uuid.NewString()
				createClient(t, root, ow, admin.ID, clientID, oauthresource.ClientTypePublic, oauthresource.ScopePolicyInheritUserPermissions, opt.NewEmpty[string](), standardScopes(), []string{oauthGrantDeviceCode})

				start := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(root, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
					ClientId: clientID,
					Scope:    ptr("openid profile offline_access"),
				}))(t, http.StatusOK)
				r.NotNil(start.JSON200)
				r.NotNil(start.JSON200.DeviceCode)
				r.NotNil(start.JSON200.UserCode)

				tests.AssertRequest(cl.OAuthDeviceConsentWithResponse(root, &openapi.OAuthDeviceConsentParams{
					UserCode: start.JSON200.UserCode,
				}, adminSession))(t, http.StatusOK)

				tests.AssertRequest(cl.OAuthDeviceConsentSubmitWithResponse(root, openapi.OAuthDeviceConsentSubmitJSONRequestBody{
					UserCode: *start.JSON200.UserCode,
					Decision: openapi.OAuthDeviceDecisionApprove,
				}, adminSession))(t, http.StatusOK)

				token := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:  oauthGrantDeviceCode,
					ClientId:   clientID,
					DeviceCode: start.JSON200.DeviceCode,
				}))(t, http.StatusOK)
				r.NotNil(token.JSON200)
				r.NotNil(token.JSON200.AccessToken)
				a.Nil(token.JSON200.RefreshToken)
			})
		}))
	}))
}

func TestOAuthSecurityHardeningDiscoveryURLs(t *testing.T) {
	t.Parallel()

	cfg := oauthConfig(t)
	publicAPIAddress, err := url.Parse("http://localhost:8000/")
	require.NoError(t, err)
	cfg.PublicAPIAddress = *publicAPIAddress

	integration.Test(t, cfg, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			a := assert.New(t)
			r := require.New(t)

			req, err := http.NewRequestWithContext(root, http.MethodGet, ts.URL+"/.well-known/openid-configuration", nil)
			r.NoError(err)
			resp, err := http.DefaultClient.Do(req)
			r.NoError(err)
			defer resp.Body.Close()
			r.Equal(http.StatusOK, resp.StatusCode)
			a.Equal("public, max-age=3600", resp.Header.Get("Cache-Control"))

			var discovery struct {
				Issuer                            string   `json:"issuer"`
				DeviceAuthorizationEndpoint       string   `json:"device_authorization_endpoint"`
				TokenEndpoint                     string   `json:"token_endpoint"`
				UserinfoEndpoint                  string   `json:"userinfo_endpoint"`
				JWKSURI                           string   `json:"jwks_uri"`
				TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
			}
			r.NoError(json.NewDecoder(resp.Body).Decode(&discovery))
			a.Equal("http://localhost:8000", discovery.Issuer)
			a.Equal("http://localhost:8000/api/oauth/device_authorization", discovery.DeviceAuthorizationEndpoint)
			a.Equal("http://localhost:8000/api/oauth/token", discovery.TokenEndpoint)
			a.Equal("http://localhost:8000/api/oauth/userinfo", discovery.UserinfoEndpoint)
			a.Equal("http://localhost:8000/api/oauth/jwks", discovery.JWKSURI)
			a.ElementsMatch([]string{"none", "client_secret_basic", "client_secret_post"}, discovery.TokenEndpointAuthMethodsSupported)

			req2, err := http.NewRequestWithContext(root, http.MethodGet, ts.URL+"/.well-known/oauth-authorization-server", nil)
			r.NoError(err)
			resp2, err := http.DefaultClient.Do(req2)
			r.NoError(err)
			defer resp2.Body.Close()
			r.Equal(http.StatusOK, resp2.StatusCode)
			a.Equal("public, max-age=3600", resp2.Header.Get("Cache-Control"))

			var metadata struct {
				Issuer                            string   `json:"issuer"`
				AuthorizationEndpoint             string   `json:"authorization_endpoint"`
				TokenEndpoint                     string   `json:"token_endpoint"`
				JWKSURI                           string   `json:"jwks_uri"`
				DeviceAuthorizationEndpoint       string   `json:"device_authorization_endpoint"`
				TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
			}
			body2, err := io.ReadAll(resp2.Body)
			r.NoError(err)
			r.NoError(json.Unmarshal(body2, &metadata))
			a.Equal("http://localhost:8000", metadata.Issuer)
			a.Equal("http://localhost:8000/api/oauth/authorize", metadata.AuthorizationEndpoint)
			a.Equal("http://localhost:8000/api/oauth/token", metadata.TokenEndpoint)
			a.ElementsMatch([]string{"none", "client_secret_basic", "client_secret_post"}, metadata.TokenEndpointAuthMethodsSupported)
			a.Equal("http://localhost:8000/api/oauth/jwks", metadata.JWKSURI)
			a.Equal("http://localhost:8000/api/oauth/device_authorization", metadata.DeviceAuthorizationEndpoint)

			var rawMetadata map[string]any
			r.NoError(json.Unmarshal(body2, &rawMetadata))
			a.NotContains(rawMetadata, "userinfo_endpoint")
			a.NotContains(rawMetadata, "id_token_signing_alg_values_supported")

			// RFC 9728 OAuth Protected Resource Metadata
			req3, err := http.NewRequestWithContext(root, http.MethodGet, ts.URL+"/.well-known/oauth-protected-resource", nil)
			r.NoError(err)
			resp3, err := http.DefaultClient.Do(req3)
			r.NoError(err)
			defer resp3.Body.Close()
			r.Equal(http.StatusOK, resp3.StatusCode)
			a.Equal("public, max-age=3600", resp3.Header.Get("Cache-Control"))

			var prm struct {
				Resource               string   `json:"resource"`
				AuthorizationServers   []string `json:"authorization_servers"`
				BearerMethodsSupported []string `json:"bearer_methods_supported"`
			}
			body3, err := io.ReadAll(resp3.Body)
			r.NoError(err)
			r.NoError(json.Unmarshal(body3, &prm))
			a.Equal("http://localhost:8000", prm.Resource)
			r.Len(prm.AuthorizationServers, 1)
			a.Equal("http://localhost:8000", prm.AuthorizationServers[0])
			r.Contains(prm.BearerMethodsSupported, "header")

			// API resource
			req4, err := http.NewRequestWithContext(root, http.MethodGet, ts.URL+"/.well-known/oauth-protected-resource/api", nil)
			r.NoError(err)
			resp4, err := http.DefaultClient.Do(req4)
			r.NoError(err)
			defer resp4.Body.Close()
			r.Equal(http.StatusOK, resp4.StatusCode)

			var prmAPI struct {
				Resource             string   `json:"resource"`
				AuthorizationServers []string `json:"authorization_servers"`
				ScopesSupported      []string `json:"scopes_supported"`
			}
			body4, err := io.ReadAll(resp4.Body)
			r.NoError(err)
			r.NoError(json.Unmarshal(body4, &prmAPI))
			a.Equal("http://localhost:8000/api", prmAPI.Resource)
			r.Len(prmAPI.AuthorizationServers, 1)
			a.Equal("http://localhost:8000", prmAPI.AuthorizationServers[0])
			r.Contains(prmAPI.ScopesSupported, "openid")
			r.Contains(prmAPI.ScopesSupported, "profile")
			r.Contains(prmAPI.ScopesSupported, "email")
			r.Contains(prmAPI.ScopesSupported, "offline_access")

			// MCP SSE resource. MCP is disabled in this config, so the resource
			// must not advertise scopes (it is not a live protected resource).
			req5, err := http.NewRequestWithContext(root, http.MethodGet, ts.URL+"/.well-known/oauth-protected-resource/mcp/sse", nil)
			r.NoError(err)
			resp5, err := http.DefaultClient.Do(req5)
			r.NoError(err)
			defer resp5.Body.Close()
			r.Equal(http.StatusOK, resp5.StatusCode)

			var prmMCP struct {
				Resource             string   `json:"resource"`
				AuthorizationServers []string `json:"authorization_servers"`
				ScopesSupported      []string `json:"scopes_supported"`
			}
			body5, err := io.ReadAll(resp5.Body)
			r.NoError(err)
			r.NoError(json.Unmarshal(body5, &prmMCP))
			a.Equal("http://localhost:8000/mcp/sse", prmMCP.Resource)
			r.Len(prmMCP.AuthorizationServers, 1)
			a.Equal("http://localhost:8000", prmMCP.AuthorizationServers[0])
			a.Empty(prmMCP.ScopesSupported)
		}))
	}))
}

func TestOAuthProtectedResourceMCPScopesWhenMCPEnabled(t *testing.T) {
	t.Parallel()

	cfg := oauthConfig(t)
	cfg.MCPEnabled = true

	integration.Test(t, cfg, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			a := assert.New(t)
			r := require.New(t)

			req, err := http.NewRequestWithContext(root, http.MethodGet, ts.URL+"/.well-known/oauth-protected-resource/mcp/sse", nil)
			r.NoError(err)
			resp, err := http.DefaultClient.Do(req)
			r.NoError(err)
			defer resp.Body.Close()
			r.Equal(http.StatusOK, resp.StatusCode)

			var prmMCP struct {
				Resource        string   `json:"resource"`
				ScopesSupported []string `json:"scopes_supported"`
			}
			r.NoError(json.NewDecoder(resp.Body).Decode(&prmMCP))
			a.Equal("http://localhost:8000/mcp/sse", prmMCP.Resource)
			r.Contains(prmMCP.ScopesSupported, "openid")
		}))
	}))
}

func TestOAuthUserInfoUnauthorisedChallenge(t *testing.T) {
	t.Parallel()

	integration.Test(t, oauthConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			a := assert.New(t)
			r := require.New(t)

			req, err := http.NewRequestWithContext(root, http.MethodGet, ts.URL+"/api/oauth/userinfo", nil)
			r.NoError(err)
			resp, err := http.DefaultClient.Do(req)
			r.NoError(err)
			defer resp.Body.Close()

			r.Equal(http.StatusUnauthorized, resp.StatusCode)
			a.Equal(
				`Bearer resource_metadata="http://localhost:8000/.well-known/oauth-protected-resource/api"`,
				resp.Header.Get("WWW-Authenticate"),
			)
		}))
	}))
}

func TestOAuthSecurityHardeningUserInfoScopes(t *testing.T) {
	t.Parallel()

	integration.Test(t, oauthConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		ow *oauth_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			t.Run("openid_only_userinfo_returns_subject_only", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				token := issueDeviceToken(t, root, cl, ow, admin.ID, adminSession, "openid")

				resp := tests.AssertRequest(cl.OAuthUserInfoWithResponse(root, bearer(*token.JSON200.AccessToken)))(t, http.StatusOK)
				r.NotNil(resp.JSON200)
				r.NotNil(resp.JSON200.Sub)
				a.Equal(admin.ID.String(), *resp.JSON200.Sub)
				a.Nil(resp.JSON200.Name)
				a.Nil(resp.JSON200.Email)
				a.Nil(resp.JSON200.EmailVerified)
				a.Nil(resp.JSON200.PreferredUsername)
			})

			t.Run("profile_and_email_userinfo_returns_scope_granted_claims", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				token := issueDeviceToken(t, root, cl, ow, admin.ID, adminSession, "openid profile email")

				resp := tests.AssertRequest(cl.OAuthUserInfoWithResponse(root, bearer(*token.JSON200.AccessToken)))(t, http.StatusOK)
				r.NotNil(resp.JSON200)
				r.NotNil(resp.JSON200.Sub)
				r.NotNil(resp.JSON200.Name)
				r.NotNil(resp.JSON200.PreferredUsername)
				r.NotNil(resp.JSON200.Email)
				r.NotNil(resp.JSON200.EmailVerified)
				a.Equal(admin.ID.String(), *resp.JSON200.Sub)
				a.Equal(admin.Name, *resp.JSON200.Name)
				a.Equal(admin.Handle, *resp.JSON200.PreferredUsername)
			})

			t.Run("missing_userinfo_bearer_is_unauthorised", func(t *testing.T) {
				a := assert.New(t)

				resp := tests.AssertRequest(cl.OAuthUserInfoWithResponse(root))(t, http.StatusUnauthorized)
				a.NotNil(resp)
			})
		}))
	}))
}

func issueDeviceToken(
	t *testing.T,
	ctx context.Context,
	cl *openapi.ClientWithResponses,
	ow *oauth_writer.Writer,
	accountID account.AccountID,
	session openapi.RequestEditorFn,
	scope string,
) *openapi.OAuthTokenResponse {
	t.Helper()

	clientID := "userinfo-" + uuid.NewString()
	createClient(t, ctx, ow, accountID, clientID, oauthresource.ClientTypePublic, oauthresource.ScopePolicyExplicit, opt.NewEmpty[string](), standardScopes(), []string{oauthGrantDeviceCode})

	start := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(ctx, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
		ClientId: clientID,
		Scope:    ptr(scope),
	}))(t, http.StatusOK)
	require.NotNil(t, start.JSON200)
	require.NotNil(t, start.JSON200.DeviceCode)
	require.NotNil(t, start.JSON200.UserCode)

	tests.AssertRequest(cl.OAuthDeviceConsentWithResponse(ctx, &openapi.OAuthDeviceConsentParams{
		UserCode: start.JSON200.UserCode,
	}, session))(t, http.StatusOK)

	tests.AssertRequest(cl.OAuthDeviceConsentSubmitWithResponse(ctx, openapi.OAuthDeviceConsentSubmitJSONRequestBody{
		UserCode: *start.JSON200.UserCode,
		Decision: openapi.OAuthDeviceDecisionApprove,
	}, session))(t, http.StatusOK)

	token := tests.AssertRequest(oauthToken(t, ctx, cl, oauthTokenRequest{
		GrantType:  oauthGrantDeviceCode,
		ClientId:   clientID,
		DeviceCode: start.JSON200.DeviceCode,
	}))(t, http.StatusOK)
	require.NotNil(t, token.JSON200)
	require.NotNil(t, token.JSON200.AccessToken)

	return token
}

func TestOAuthSecurityHardeningDeviceFlowCrossAccountClaim(t *testing.T) {
	t.Parallel()

	integration.Test(t, oauthConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		assignments *role_assign.Assignment,
		ow *oauth_writer.Writer,
		roles *role_repo.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			_, owner := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			attackerCtx, attacker := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			victimCtx, victim := e2e.WithAccount(root, aw, seed.Account_002_Frigg)
			grantOAuthClientUse(t, root, roles, assignments, attacker.ID)
			grantOAuthClientUse(t, root, roles, assignments, victim.ID)
			attackerSession := sh.WithSession(attackerCtx)
			victimSession := sh.WithSession(victimCtx)

			t.Run("attacker_cannot_approve_device_flow_claimed_by_victim", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				clientID := "cross-account-claim-" + uuid.NewString()
				createClient(t, root, ow, owner.ID, clientID, oauthresource.ClientTypePublic, oauthresource.ScopePolicyInheritUserPermissions, opt.NewEmpty[string](), standardScopes(), []string{oauthGrantDeviceCode})

				start := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(root, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
					ClientId: clientID,
					Scope:    ptr("openid profile offline_access"),
				}))(t, http.StatusOK)
				r.NotNil(start.JSON200)
				r.NotNil(start.JSON200.UserCode)

				tests.AssertRequest(cl.OAuthDeviceConsentWithResponse(root, &openapi.OAuthDeviceConsentParams{
					UserCode: start.JSON200.UserCode,
				}, victimSession))(t, http.StatusOK)

				blocked := tests.AssertRequest(cl.OAuthDeviceConsentSubmitWithResponse(root, openapi.OAuthDeviceConsentSubmitJSONRequestBody{
					UserCode: *start.JSON200.UserCode,
					Decision: openapi.OAuthDeviceDecisionApprove,
				}, attackerSession))(t, http.StatusBadRequest)
				r.NotNil(blocked.JSON400)
				a.Equal("access_denied", blocked.JSON400.Error)
			})

			t.Run("attacker_cannot_claim_already_claimed_device_flow", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				clientID := "cross-account-claim-second-" + uuid.NewString()
				createClient(t, root, ow, owner.ID, clientID, oauthresource.ClientTypePublic, oauthresource.ScopePolicyInheritUserPermissions, opt.NewEmpty[string](), standardScopes(), []string{oauthGrantDeviceCode})

				start := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(root, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
					ClientId: clientID,
					Scope:    ptr("openid profile offline_access"),
				}))(t, http.StatusOK)
				r.NotNil(start.JSON200)
				r.NotNil(start.JSON200.UserCode)

				tests.AssertRequest(cl.OAuthDeviceConsentWithResponse(root, &openapi.OAuthDeviceConsentParams{
					UserCode: start.JSON200.UserCode,
				}, victimSession))(t, http.StatusOK)

				blocked := tests.AssertRequest(cl.OAuthDeviceConsentWithResponse(root, &openapi.OAuthDeviceConsentParams{
					UserCode: start.JSON200.UserCode,
				}, attackerSession))(t, http.StatusBadRequest)
				r.NotNil(blocked.JSON400)
				a.Equal("access_denied", blocked.JSON400.Error)
			})
		}))
	}))
}

func TestOAuthSecurityHardeningScopeCapOnRefresh(t *testing.T) {
	t.Parallel()

	integration.Test(t, oauthConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		assignments *role_assign.Assignment,
		ow *oauth_writer.Writer,
		roles *role_repo.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			_, owner := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			memberCtx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			grantOAuthClientUse(t, root, roles, assignments, member.ID)
			memberSession := sh.WithSession(memberCtx)

			t.Run("inheriting_client_refresh_does_not_escalate_after_promotion", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				clientID := "scope-cap-refresh-" + uuid.NewString()
				createClient(t, root, ow, owner.ID, clientID, oauthresource.ClientTypePublic, oauthresource.ScopePolicyInheritUserPermissions, opt.NewEmpty[string](), standardScopes(), []string{oauthGrantDeviceCode, oauthGrantRefreshToken})

				start := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(root, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
					ClientId: clientID,
					Scope:    ptr("openid profile offline_access"),
				}))(t, http.StatusOK)
				r.NotNil(start.JSON200)
				r.NotNil(start.JSON200.DeviceCode)
				r.NotNil(start.JSON200.UserCode)

				tests.AssertRequest(cl.OAuthDeviceConsentWithResponse(root, &openapi.OAuthDeviceConsentParams{
					UserCode: start.JSON200.UserCode,
				}, memberSession))(t, http.StatusOK)

				tests.AssertRequest(cl.OAuthDeviceConsentSubmitWithResponse(root, openapi.OAuthDeviceConsentSubmitJSONRequestBody{
					UserCode: *start.JSON200.UserCode,
					Decision: openapi.OAuthDeviceDecisionApprove,
				}, memberSession))(t, http.StatusOK)

				initial := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:  oauthGrantDeviceCode,
					ClientId:   clientID,
					DeviceCode: start.JSON200.DeviceCode,
				}))(t, http.StatusOK)
				r.NotNil(initial.JSON200)
				r.NotNil(initial.JSON200.RefreshToken)
				a.NotContains(*initial.JSON200.Scope, "ADMINISTRATOR")

				aw.Update(root, member.ID, account_writer.SetAdmin(true))

				refreshed := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:    oauthGrantRefreshToken,
					ClientId:     clientID,
					RefreshToken: initial.JSON200.RefreshToken,
				}))(t, http.StatusOK)
				r.NotNil(refreshed.JSON200)
				r.NotNil(refreshed.JSON200.AccessToken)
				a.NotContains(*refreshed.JSON200.Scope, "ADMINISTRATOR")
			})
		}))
	}))
}
