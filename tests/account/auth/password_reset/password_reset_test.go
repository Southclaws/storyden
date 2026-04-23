package password_reset_test

import (
	"context"
	"net/http"
	"regexp"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/mailer"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestPasswordReset(t *testing.T) {
	t.Parallel()

	integration.Test(t, &config.Config{
		JWTSecret: []byte("07d422e512b23a056ccc953994d1593f"),
	}, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		mail mailer.Sender,
	) {
		inbox := mail.(*mailer.Mock)

		lc.Append(fx.StartHook(func() {
			t.Run("nonexistent", func(t *testing.T) {
				r := require.New(t)

				email := xid.New().String() + "@storyden.org"
				request, err := cl.AuthPasswordResetRequestEmailWithResponse(root, openapi.AuthEmailPasswordReset{
					Email: email,
					TokenUrl: struct {
						Query string `json:"query"`
						Url   string `json:"url"`
					}{
						Url:   "http://localhost:3000/reset",
						Query: "token",
					},
				})
				r.NoError(err)
				r.Equal(http.StatusNotFound, request.StatusCode())
			})

			t.Run("reset_password", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				email := xid.New().String() + "@storyden.org"
				password := "mysupersecretpasswordwhichissosecretiforgotwhatitwas"
				signupEmailCount := inbox.Count()

				// Sign up with username + password
				signup, err := cl.AuthEmailPasswordSignupWithResponse(root, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{Email: email, Password: password})
				tests.Ok(t, err, signup)
				signupSession := e2e.WithSessionFromHeader(t, root, signup.HTTPResponse.Header)

				accountGet, err := cl.AccountGetWithResponse(root, signupSession)
				tests.Ok(t, err, accountGet)

				tests.WaitForNextEmail(t, inbox, signupEmailCount)
				resetEmailCount := inbox.Count()

				// oh no! I forgot my password :( let's reset it
				request, err := cl.AuthPasswordResetRequestEmailWithResponse(root, openapi.AuthEmailPasswordReset{
					Email: email,
					TokenUrl: struct {
						Query string `json:"query"`
						Url   string `json:"url"`
					}{
						Url:   "http://localhost:3000/reset",
						Query: "token",
					},
				})
				tests.Ok(t, err, request)

				resetEmail := tests.WaitForNextEmail(t, inbox, resetEmailCount)
				a.Equal(accountGet.JSON200.Name, resetEmail.Name)
				a.Equal(email, resetEmail.Address.Address)
				token := regexp.MustCompile(`\?token=(.+)`).FindStringSubmatch(resetEmail.Plain)[1]

				reset, err := cl.AuthPasswordResetWithResponse(root, openapi.AuthPasswordResetJSONRequestBody{
					Token: token,
					New:   "newpassword",
				})
				tests.Ok(t, err, reset)

				r.Equal(signup.JSON200.Id, reset.JSON200.Id)

				get, err := cl.AccountGetWithResponse(root, signupSession)
				tests.Ok(t, err, get)
				a.Equal(signup.JSON200.Id, get.JSON200.Id)
			})
		}))
	}))
}
