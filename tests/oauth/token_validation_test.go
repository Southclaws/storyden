package oauth_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	oauthresource "github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/app/resources/oauth/oauth_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestOAuthBearerTokenValidation(t *testing.T) {
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
			a := assert.New(t)
			r := require.New(t)

			adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			clientID := "bearer-cli-" + uuid.NewString()
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
			r.NotNil(token.JSON200.AccessToken)

			t.Run("valid_oauth_token_authenticates_request", func(t *testing.T) {
				r := require.New(t)

				adminSettings := tests.AssertRequest(cl.AdminSettingsGetWithResponse(root, bearer(*token.JSON200.AccessToken)))(t, http.StatusOK)
				r.NotNil(adminSettings.JSON200)
			})

			t.Run("tampered_oauth_token_is_rejected", func(t *testing.T) {
				tampered := *token.JSON200.AccessToken
				if strings.HasSuffix(tampered, "a") {
					tampered = strings.TrimSuffix(tampered, "a") + "b"
				} else {
					tampered += "a"
				}

				resp := tests.AssertRequest(cl.AdminSettingsGetWithResponse(root, bearer(tampered)))(t, http.StatusUnauthorized)
				a.NotNil(resp)
			})

			t.Run("malformed_oauth_token_is_rejected", func(t *testing.T) {
				resp := tests.AssertRequest(cl.AdminSettingsGetWithResponse(root, bearer("ey.malformed.token")))(t, http.StatusUnauthorized)
				a.NotNil(resp)
			})
		}))
	}))
}

func TestOAuthExpiredBearerTokenIsRejected(t *testing.T) {
	t.Parallel()

	integration.Test(t, oauthConfigWithAccessTTL(t, -1*time.Minute), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		ow *oauth_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			a := assert.New(t)
			r := require.New(t)

			adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			clientID := "expired-cli-" + uuid.NewString()
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
			r.NotNil(token.JSON200.AccessToken)

			resp := tests.AssertRequest(cl.AdminSettingsGetWithResponse(root, bearer(*token.JSON200.AccessToken)))(t, http.StatusUnauthorized)
			a.NotNil(resp)
		}))
	}))
}
