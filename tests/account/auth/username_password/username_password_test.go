package username_password_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	session1 "github.com/Southclaws/storyden/app/transports/http/middleware/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/tests"
)

func TestUsernamePasswordAuth(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		cj *session1.Jar,
		set *settings.SettingsRepository,
		accountQuery *account_querier.Querier,
	) {
		lc.Append(fx.StartHook(func() {
			utils.Must(set.Set(root, settings.Settings{
				AuthenticationMode: opt.New(authentication.ModeHandle),
			}))

			t.Run("register_success", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				handle := xid.New().String()
				password := "password"

				// Sign up with username + password
				signup, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPasswordSignupJSONRequestBody{Identifier: handle, Token: password})
				tests.Ok(t, err, signup)

				// Sign in with username + password
				signin, err := cl.AuthPasswordSigninWithResponse(root, openapi.AuthPasswordSigninJSONRequestBody{Identifier: handle, Token: password})
				tests.Ok(t, err, signin)
				a.NotEmpty(signin.HTTPResponse.Header.Get("Set-Cookie"))

				accountID := account.AccountID(openapi.GetAccountID(signup.JSON200.Id))
				ctx1 := session.WithAccountID(root, accountID)
				session := e2e.WithSession(ctx1, cj)

				// Get own account
				get, err := cl.AccountGetWithResponse(root, session)
				tests.Ok(t, err, get)
				r.Equal(openapi.AccountVerifiedStatusNone, get.JSON200.VerifiedStatus)
				r.Len(get.JSON200.EmailAddresses, 0)
			})

			t.Run("register_fail_duplicate", func(t *testing.T) {
				handle := xid.New().String()
				password := "password"
				password2 := "password2"

				// Sign up with username + password
				signup, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPasswordSignupJSONRequestBody{Identifier: handle, Token: password})
				tests.Ok(t, err, signup)

				// Sign up again with the same username
				signupAgain, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPasswordSignupJSONRequestBody{Identifier: handle, Token: password2})
				tests.Status(t, err, signupAgain, http.StatusConflict)
			})
		}))
	}))
}

func TestUsernamePasswordAuthFailsInEmailMode(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		cj *session1.Jar,
		set *settings.SettingsRepository,
		accountQuery *account_querier.Querier,
	) {
		lc.Append(fx.StartHook(func() {
			utils.Must(set.Set(root, settings.Settings{
				AuthenticationMode: opt.New(authentication.ModeEmail),
			}))

			t.Run("register_with_username_only_fails", func(t *testing.T) {
				handle := xid.New().String()
				password := "password"

				// Sign up with username + password
				signup, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPasswordSignupJSONRequestBody{Identifier: handle, Token: password})
				tests.Status(t, err, signup, http.StatusBadRequest)
			})

			t.Run("register_with_email_login_with_username_success", func(t *testing.T) {
				r := require.New(t)

				email := xid.New().String() + "@storyden.org"
				handle := xid.New().String()
				password := "password"

				// Sign up with username + password
				signup, err := cl.AuthEmailPasswordSignupWithResponse(root, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{
					Email:    email,
					Handle:   &handle,
					Password: password,
				})
				tests.Ok(t, err, signup)

				// Sign in using just a username
				signin, err := cl.AuthPasswordSigninWithResponse(root, openapi.AuthPasswordSigninJSONRequestBody{
					Identifier: handle,
					Token:      password,
				})
				tests.Ok(t, err, signin)
				session := e2e.WithSessionFromHeader(t, root, signin.HTTPResponse.Header)

				// Get own account
				get, err := cl.AccountGetWithResponse(root, session)
				tests.Ok(t, err, get)
				r.Equal(openapi.AccountVerifiedStatusNone, get.JSON200.VerifiedStatus)
				r.Len(get.JSON200.EmailAddresses, 1)
			})
		}))
	}))
}
