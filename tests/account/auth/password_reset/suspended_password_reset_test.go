package password_reset_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestPasswordResetRejectsSuspendedAccount(t *testing.T) {
	t.Parallel()

	integration.Test(t, &config.Config{
		JWTSecret: []byte("07d422e512b23a056ccc953994d1593f"),
	}, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			email := xid.New().String() + "@storyden.org"
			password := "originalsecretpassword"

			signup, err := cl.AuthEmailPasswordSignupWithResponse(root, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{
				Email:    email,
				Password: password,
			})
			tests.Ok(t, err, signup)
			userSession := e2e.WithSessionFromHeader(t, root, signup.HTTPResponse.Header)

			userAccount, err := cl.AccountGetWithResponse(root, userSession)
			tests.Ok(t, err, userAccount)

			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			tokenResp, err := cl.AccountPasswordResetTokenGetWithResponse(root, userAccount.JSON200.Id, adminSession)
			tests.Ok(t, err, tokenResp)

			suspend, err := cl.AdminAccountBanCreateWithResponse(root, userAccount.JSON200.Handle, adminSession)
			tests.Ok(t, err, suspend)

			resetResp, err := cl.AuthPasswordResetWithResponse(root, openapi.AuthPasswordResetJSONRequestBody{
				Token: tokenResp.JSON200.Token,
				New:   "newsuperpassword",
			})
			r.NoError(err)
			r.Equal(http.StatusForbidden, resetResp.StatusCode())
			a.Empty(resetResp.HTTPResponse.Header.Values("Set-Cookie"))
		}))
	}))
}
