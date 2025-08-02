package username_password_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestUsernamePasswordAuth(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		accountQuery *account_querier.Querier,
	) {
		lc.Append(fx.StartHook(func() {
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
				ctx1 := e2e.WithAccountID(root, accountID)
				session := sh.WithSession(ctx1)

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

			t.Run("change_password", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				hand1 := "tester1-" + xid.New().String()

				// Sign up for a new account with a password
				a1, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{
					Identifier: hand1,
					Token:      "password",
				})
				r.NoError(err)
				r.NotNil(a1)
				r.Equal(http.StatusOK, a1.StatusCode())
				id1, err := xid.FromString(a1.JSON200.Id)
				r.NoError(err)

				ctx1 := e2e.WithAccountID(root, account.AccountID(id1))

				// Get the new account
				get1, err := cl.AccountGetWithResponse(root, sh.WithSession(ctx1))
				r.NoError(err)
				r.NotNil(get1)
				r.Equal(http.StatusOK, get1.StatusCode())
				a.Equal(hand1, get1.JSON200.Handle)

				// Change the password
				change1, err := cl.AuthPasswordUpdateWithResponse(root, openapi.AuthPasswordMutableProps{
					Old: "password",
					New: "wordpass",
				}, sh.WithSession(ctx1))
				r.NoError(err)
				r.NotNil(change1)
				r.Equal(http.StatusOK, change1.StatusCode())

				// Log in to the new account with the old password - fails
				signin1, err := cl.AuthPasswordSigninWithResponse(root, openapi.AuthPair{
					Identifier: hand1,
					Token:      "password",
				})
				r.NoError(err)
				r.NotNil(signin1)
				r.Equal(http.StatusUnauthorized, signin1.StatusCode())

				// Log in to the new account with the new password - succeeds
				signin2, err := cl.AuthPasswordSigninWithResponse(root, openapi.AuthPair{
					Identifier: hand1,
					Token:      "wordpass",
				})
				r.NoError(err)
				r.NotNil(signin2)
				r.Equal(http.StatusOK, signin2.StatusCode())

				// Sign out
				signout, err := cl.AuthProviderLogoutWithResponse(root, sh.WithSession(ctx1))
				tests.Ok(t, err, signout)
				a.Contains(signout.HTTPResponse.Header.Get("Set-Cookie"), "storyden-session=;")
			})

			t.Run("register_fail_invalid_password", func(t *testing.T) {
				handle := xid.New().String()

				signup, err := cl.AuthPasswordSignupWithResponse(root, nil,
					openapi.AuthPasswordSignupJSONRequestBody{
						Identifier: handle,
						Token:      "weak", // too short password
					})
				tests.Status(t, err, signup, http.StatusBadRequest)
			})
		}))
	}))
}

func TestUsernamePasswordAuthMultiMethod(t *testing.T) {
	t.Parallel()

	// NOTE: A bit of complexity here because Storyden supports both username &
	// email based registration and login, however registration with an email
	// is only possible when an email client has been enabled. This opens a few
	// edge cases where some users may register with only a username and then an
	// administrator enables email-based auth. Those users need to still be able
	// to log in with their username so the authentication APIs both work at the
	// same time if email is enabled. This also means that we need to ensure the
	// password change mechanism works for both username and email based auth.
	//
	// As a result, despite there being two independent auth providers for these
	// (username_password and email_password) the operation AuthPasswordUpdate
	// will work with either. The reason this works is because the auth records
	// are set up with uniqueness constraints across the service name and the
	// token type. If an account registers with a username and password, then
	// their auth record will hold a "token_type" of "password" and if they use
	// email, the "token_type" is still "password" so this allows the service to
	// operate on either type of auth record because the "token_type" is shared
	// across both auth providers. This also means that only a single record of
	// "token_type": "password" may exist, so there's no way to end up with two.

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		accountQuery *account_querier.Querier,
	) {
		lc.Append(fx.StartHook(func() {
			t.Run("signup_with_email_but_login_with_username", func(t *testing.T) {
				r := require.New(t)

				email := xid.New().String() + "@storyden.org"
				handle := xid.New().String()
				password := "password"

				// Sign up with email + password
				signup, err := cl.AuthEmailPasswordSignupWithResponse(root, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{
					Email:    email,
					Handle:   &handle,
					Password: password,
				})
				tests.Ok(t, err, signup)

				// Sign in using just a username, no email
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

			t.Run("signup_with_email_then_change_password", func(t *testing.T) {
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

				// Change password
				newPassword := "wordpass"
				change, err := cl.AuthPasswordUpdateWithResponse(root, openapi.AuthPasswordMutableProps{
					Old: "password",
					New: newPassword,
				}, session)
				tests.Ok(t, err, change)

				// Old password fails
				signin2, err := cl.AuthPasswordSigninWithResponse(root, openapi.AuthPasswordSigninJSONRequestBody{
					Identifier: handle,
					Token:      password,
				})
				tests.Status(t, err, signin2, http.StatusUnauthorized)

				// New password succeeds
				signin3, err := cl.AuthPasswordSigninWithResponse(root, openapi.AuthPasswordSigninJSONRequestBody{
					Identifier: handle,
					Token:      newPassword,
				})
				tests.Ok(t, err, signin3)
				session2 := e2e.WithSessionFromHeader(t, root, signin3.HTTPResponse.Header)

				// Get own account
				get, err := cl.AccountGetWithResponse(root, session2)
				tests.Ok(t, err, get)
				r.Equal(openapi.AccountVerifiedStatusNone, get.JSON200.VerifiedStatus)
				r.Len(get.JSON200.EmailAddresses, 1)
			})
		}))
	}))
}
