package oauth_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

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

func TestOAuthSuspendedAccountCannotUseExistingGrants(t *testing.T) {
	t.Parallel()

	integration.Test(t, oauthConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		ow *oauth_writer.Writer,
		roles *role_repo.Repository,
		assignments *role_assign.Assignment,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			memberCtx, member := e2e.WithAccount(root, aw, seed.Account_002_Frigg)
			grantOAuthClientUse(t, root, roles, assignments, member.ID)
			memberSession := sh.WithSession(memberCtx)

			clientID := "suspended-account-" + uuid.NewString()
			createClient(t, root, ow, member.ID, clientID, oauthresource.ClientTypePublic, oauthresource.ScopePolicyExplicit, opt.NewEmpty[string](), standardScopes(), []string{oauthGrantDeviceCode, oauthGrantRefreshToken})

			refreshToken := mintRefreshTokenViaDeviceFlow(t, root, cl, memberSession, clientID)

			rotate := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
				GrantType:    oauthGrantRefreshToken,
				ClientId:     clientID,
				RefreshToken: &refreshToken,
			}))(t, http.StatusOK)
			r.NotNil(rotate.JSON200)
			r.NotNil(rotate.JSON200.AccessToken)
			r.NotNil(rotate.JSON200.RefreshToken)

			access := bearer(*rotate.JSON200.AccessToken)
			refreshToken = *rotate.JSON200.RefreshToken

			// A device authorisation approved before suspension but not yet
			// exchanged for tokens.
			pending := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(root, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
				ClientId: clientID,
				Scope:    new("openid profile"),
			}))(t, http.StatusOK)
			r.NotNil(pending.JSON200)

			tests.AssertRequest(cl.OAuthDeviceConsentWithResponse(root, &openapi.OAuthDeviceConsentParams{
				UserCode: pending.JSON200.UserCode,
			}, memberSession))(t, http.StatusOK)

			tests.AssertRequest(cl.OAuthDeviceConsentSubmitWithResponse(root, openapi.OAuthDeviceConsentSubmitJSONRequestBody{
				UserCode: *pending.JSON200.UserCode,
				Decision: openapi.OAuthDeviceDecisionApprove,
			}, memberSession))(t, http.StatusOK)

			tests.AssertRequest(cl.OAuthUserInfoWithResponse(root, access))(t, http.StatusOK)

			tests.AssertRequest(cl.AdminAccountBanCreateWithResponse(root, member.Handle, adminSession))(t, http.StatusOK)

			t.Run("refresh_token_is_rejected", func(t *testing.T) {
				refresh := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:    oauthGrantRefreshToken,
					ClientId:     clientID,
					RefreshToken: &refreshToken,
				}))(t, http.StatusBadRequest)
				r.NotNil(refresh.JSON400)
				a.Equal("invalid_grant", refresh.JSON400.Error)
			})

			t.Run("approved_device_code_is_rejected", func(t *testing.T) {
				device := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:  oauthGrantDeviceCode,
					ClientId:   clientID,
					DeviceCode: pending.JSON200.DeviceCode,
				}))(t, http.StatusBadRequest)
				r.NotNil(device.JSON400)
				a.Equal("invalid_grant", device.JSON400.Error)
			})

			t.Run("userinfo_is_rejected", func(t *testing.T) {
				tests.AssertRequest(cl.OAuthUserInfoWithResponse(root, access))(t, http.StatusForbidden)
			})
		}))
	}))
}
