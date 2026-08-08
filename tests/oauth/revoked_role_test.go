package oauth_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/role/role_assign"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	oauthresource "github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/app/resources/oauth/oauth_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestOAuthTokenLosesPermissionWhenRoleIsRevoked(t *testing.T) {
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

			memberCtx, member := e2e.WithAccount(root, aw, seed.Account_002_Frigg)
			roleID := grantOAuthClientUse(t, root, roles, assignments, member.ID, rbac.PermissionManageCategories)
			memberSession := sh.WithSession(memberCtx)

			clientID := "revoked-role-cli-" + uuid.NewString()
			allowed := append(standardScopes(), rbac.PermissionManageCategories.String())
			createClient(t, root, ow, member.ID, clientID, oauthresource.ClientTypePublic, oauthresource.ScopePolicyExplicit, opt.NewEmpty[string](), allowed, []string{oauthGrantDeviceCode})

			start := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(root, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
				ClientId: clientID,
				Scope:    ptr("openid profile " + rbac.PermissionManageCategories.String()),
			}))(t, http.StatusOK)
			r.NotNil(start.JSON200)

			tests.AssertRequest(cl.OAuthDeviceConsentWithResponse(root, &openapi.OAuthDeviceConsentParams{
				UserCode: start.JSON200.UserCode,
			}, memberSession))(t, http.StatusOK)

			tests.AssertRequest(cl.OAuthDeviceConsentSubmitWithResponse(root, openapi.OAuthDeviceConsentSubmitJSONRequestBody{
				UserCode: *start.JSON200.UserCode,
				Decision: openapi.OAuthDeviceDecisionApprove,
			}, memberSession))(t, http.StatusOK)

			token := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
				GrantType:  oauthGrantDeviceCode,
				ClientId:   clientID,
				DeviceCode: start.JSON200.DeviceCode,
			}))(t, http.StatusOK)
			r.NotNil(token.JSON200.AccessToken)

			access := bearer(*token.JSON200.AccessToken)

			t.Run("granted_while_the_role_is_held", func(t *testing.T) {
				tests.AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
					Name: "before-revocation-" + uuid.NewString(),
				}, access))(t, http.StatusOK)
			})

			revokeOAuthClientUse(t, root, assignments, member.ID, roleID)

			t.Run("denied_once_the_role_is_gone", func(t *testing.T) {
				tests.AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
					Name: "after-revocation-" + uuid.NewString(),
				}, access))(t, http.StatusForbidden)
			})
		}))
	}))
}
