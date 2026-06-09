package oauth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestOAuthDisabledConfiguration(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			memberCtx, _ := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			memberSession := sh.WithSession(memberCtx)

			t.Run("info_does_not_advertise_oauth_capability", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				info := tests.AssertRequest(cl.GetInfoWithResponse(root))(t, http.StatusOK)
				r.NotNil(info.JSON200)
				a.NotContains(info.JSON200.Capabilities, openapi.InstanceCapability("oauth"))
				a.Equal("http://localhost", info.JSON200.WebAddress)
				a.Equal("http://localhost", info.JSON200.ApiAddress)
			})

			t.Run("discovery_returns_clear_disabled_response", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				req, err := http.NewRequestWithContext(root, http.MethodGet, ts.URL+"/.well-known/openid-configuration", nil)
				r.NoError(err)

				resp, err := http.DefaultClient.Do(req)
				r.NoError(err)
				defer resp.Body.Close()

				var body openapi.APIError
				r.NoError(json.NewDecoder(resp.Body).Decode(&body))
				a.Equal(http.StatusNotFound, resp.StatusCode)
				r.NotNil(body.Type)
				a.Equal("urn:storyden:problem:not-found", *body.Type)
				r.NotNil(body.Title)
				a.Contains(*body.Title, "not enabled")
				r.NotNil(body.Metadata)
				a.Equal("oauth_disabled", (*body.Metadata)["code"])
				a.Contains((*body.Metadata)["suggested"], "administrator")
			})

			t.Run("oauth_authorization_server_metadata_returns_clear_disabled_response", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				req, err := http.NewRequestWithContext(root, http.MethodGet, ts.URL+"/.well-known/oauth-authorization-server", nil)
				r.NoError(err)

				resp, err := http.DefaultClient.Do(req)
				r.NoError(err)
				defer resp.Body.Close()

				var body openapi.APIError
				r.NoError(json.NewDecoder(resp.Body).Decode(&body))
				a.Equal(http.StatusNotFound, resp.StatusCode)
				r.NotNil(body.Type)
				a.Equal("urn:storyden:problem:not-found", *body.Type)
				r.NotNil(body.Title)
				a.Contains(*body.Title, "not enabled")
				r.NotNil(body.Metadata)
				a.Equal("oauth_disabled", (*body.Metadata)["code"])
				a.Contains((*body.Metadata)["suggested"], "administrator")
			})

			t.Run("jwks_returns_clear_disabled_response", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				resp, err := cl.OAuthJWKSWithResponse(root)
				r.NoError(err)
				r.NotNil(resp)

				var body openapi.APIError
				r.NoError(json.Unmarshal(resp.Body, &body))
				a.Equal(http.StatusNotFound, resp.StatusCode())
				r.NotNil(body.Type)
				a.Equal("urn:storyden:problem:not-found", *body.Type)
				r.NotNil(body.Title)
				a.Contains(*body.Title, "not enabled")
			})

			t.Run("device_authorization_returns_oauth_error", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				resp := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(root, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
					ClientId: "storyden-cli",
					Scope:    ptr("openid profile offline_access"),
				}))(t, http.StatusBadRequest)
				r.NotNil(resp.JSON400)
				a.Equal("temporarily_unavailable", resp.JSON400.Error)
			})

			t.Run("token_endpoint_returns_oauth_error_before_signing_is_needed", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				deviceCode := "device-code-" + uuid.NewString()
				resp := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:  oauthGrantDeviceCode,
					ClientId:   "storyden-cli",
					DeviceCode: &deviceCode,
				}))(t, http.StatusBadRequest)
				r.NotNil(resp.JSON400)
				a.Equal("temporarily_unavailable", resp.JSON400.Error)
			})

			t.Run("authorization_code_endpoint_returns_oauth_error_for_signed_in_users", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				resp := authorizeHTTPResponse(t, root, ts, memberSession, authorizeRequest{
					ClientID:            "client-" + uuid.NewString(),
					RedirectURI:         "https://client.example/callback",
					Scope:               "openid profile",
					State:               "state-" + uuid.NewString(),
					CodeChallenge:       codeChallenge(strings.Repeat("i", 43)),
					CodeChallengeMethod: "S256",
				})
				defer resp.Body.Close()

				a.Equal(http.StatusBadRequest, resp.StatusCode)
				a.Empty(resp.Header.Get("Location"))
				r.NotNil(resp.Body)
			})

			t.Run("consent_endpoints_return_oauth_error_for_signed_in_users", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				userCode := openapi.OAuthUserCodeQuery("ABCD-EFGH")
				deviceConsent := tests.AssertRequest(cl.OAuthDeviceConsentWithResponse(root, &openapi.OAuthDeviceConsentParams{
					UserCode: &userCode,
				}, memberSession))(t, http.StatusBadRequest)
				r.NotNil(deviceConsent.JSON400)
				a.Equal("temporarily_unavailable", deviceConsent.JSON400.Error)

				deviceSubmit := tests.AssertRequest(cl.OAuthDeviceConsentSubmitWithResponse(root, openapi.OAuthDeviceConsentSubmitJSONRequestBody{
					UserCode: "ABCD-EFGH",
					Decision: openapi.OAuthDeviceDecisionApprove,
				}, memberSession))(t, http.StatusBadRequest)
				r.NotNil(deviceSubmit.JSON400)
				a.Equal("temporarily_unavailable", deviceSubmit.JSON400.Error)

				requestID := openapi.OAuthAuthorizationRequestIDQuery("request-" + uuid.NewString())
				authConsent := tests.AssertRequest(cl.OAuthAuthoriseConsentWithResponse(root, &openapi.OAuthAuthoriseConsentParams{
					RequestId: &requestID,
				}, memberSession))(t, http.StatusBadRequest)
				r.NotNil(authConsent.JSON400)
				a.Equal("temporarily_unavailable", authConsent.JSON400.Error)

				authSubmit := tests.AssertRequest(cl.OAuthAuthoriseConsentSubmitWithResponse(root, openapi.OAuthAuthoriseConsentSubmitJSONRequestBody{
					RequestId: string(requestID),
					Decision:  openapi.OAuthAuthoriseDecisionApprove,
				}, memberSession))(t, http.StatusBadRequest)
				r.NotNil(authSubmit.JSON400)
				a.Equal("temporarily_unavailable", authSubmit.JSON400.Error)
			})
		}))
	}))
}

func TestOAuthEnabledConfigurationAdvertisesCapability(t *testing.T) {
	t.Parallel()

	integration.Test(t, oauthConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
	) {
		lc.Append(fx.StartHook(func() {
			a := assert.New(t)
			r := require.New(t)

			info := tests.AssertRequest(cl.GetInfoWithResponse(root))(t, http.StatusOK)
			r.NotNil(info.JSON200)
			a.Contains(info.JSON200.Capabilities, openapi.InstanceCapability("oauth"))
		}))
	}))
}
