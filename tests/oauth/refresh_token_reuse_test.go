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
	oauthresource "github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/app/resources/oauth/oauth_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestOAuthRefreshTokenReuseRevokesFamily(t *testing.T) {
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
			r := require.New(t)
			a := assert.New(t)

			adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			clientID := "refresh-reuse-" + uuid.NewString()
			createClient(t, root, ow, admin.ID, clientID, oauthresource.ClientTypePublic, oauthresource.ScopePolicyInheritUserPermissions, opt.NewEmpty[string](), standardScopes(), []string{oauthGrantDeviceCode, oauthGrantRefreshToken})

			r1 := mintRefreshTokenViaDeviceFlow(t, root, cl, adminSession, clientID)

			// Rotate r1 -> r2. r1 is now consumed and has r2 as its replacement.
			rotate := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
				GrantType:    oauthGrantRefreshToken,
				ClientId:     clientID,
				RefreshToken: &r1,
			}))(t, http.StatusOK)
			r.NotNil(rotate.JSON200)
			r.NotNil(rotate.JSON200.RefreshToken)
			r2 := *rotate.JSON200.RefreshToken

			// Reuse the already-rotated r1: this is the theft signal.
			reuse := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
				GrantType:    oauthGrantRefreshToken,
				ClientId:     clientID,
				RefreshToken: &r1,
			}))(t, http.StatusBadRequest)
			r.NotNil(reuse.JSON400)
			a.Equal("invalid_grant", reuse.JSON400.Error)

			// The active descendant r2 must now be revoked as part of the family.
			cascaded := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
				GrantType:    oauthGrantRefreshToken,
				ClientId:     clientID,
				RefreshToken: &r2,
			}))(t, http.StatusBadRequest)
			r.NotNil(cascaded.JSON400)
			a.Equal("invalid_grant", cascaded.JSON400.Error)
		}))
	}))
}

func mintRefreshTokenViaDeviceFlow(t *testing.T, ctx context.Context, cl *openapi.ClientWithResponses, session openapi.RequestEditorFn, clientID string) string {
	t.Helper()
	r := require.New(t)

	start := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(ctx, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
		ClientId: clientID,
		Scope:    ptr("openid profile offline_access"),
	}))(t, http.StatusOK)
	r.NotNil(start.JSON200)

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
	r.NotNil(token.JSON200)
	r.NotNil(token.JSON200.RefreshToken)

	return *token.JSON200.RefreshToken
}
