package account_email_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	session1 "github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestAccountEmails(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		cj *session1.Jar,
		db *ent.Client,
		accountWrite *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			t.Run("sign_up_with_email_add_another_email_login_with_both", func(t *testing.T) {
				email1 := xid.New().String() + "first@example.com"
				password := "password"

				signup, err := cl.AuthEmailPasswordSignupWithResponse(root, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{
					Email:    email1,
					Password: password,
				})
				tests.Ok(t, err, signup)
				session := e2e.WithSessionFromHeader(t, root, signup.HTTPResponse.Header)

				email2 := xid.New().String() + "second@example.com"

				add1, err := cl.AccountEmailAddWithResponse(root, openapi.AccountEmailInitialProps{
					EmailAddress: email2,
				}, session)
				tests.Ok(t, err, add1)

				loginWithEmail1, err := cl.AuthEmailPasswordSigninWithResponse(root, openapi.AuthEmailPasswordSigninJSONRequestBody{
					Email:    email1,
					Password: password,
				})
				tests.Ok(t, err, loginWithEmail1)

				loginWithEmail2, err := cl.AuthEmailPasswordSigninWithResponse(root, openapi.AuthEmailPasswordSigninJSONRequestBody{
					Email:    email2,
					Password: password,
				})
				tests.Ok(t, err, loginWithEmail2)
			})

			t.Run("cannot_add_someone_elses_email_to_my_account", func(t *testing.T) {
				emailAttacker := xid.New().String() + "attacker@example.com"
				emailVictim := xid.New().String() + "victim@example.com"
				passwordAttacker := "passwordattacker"
				passwordVictim := "passwordvictim"

				signupVictim, err := cl.AuthEmailPasswordSignupWithResponse(root, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{
					Email:    emailVictim,
					Password: passwordVictim,
				})
				tests.Ok(t, err, signupVictim)

				signupAttacker, err := cl.AuthEmailPasswordSignupWithResponse(root, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{
					Email:    emailAttacker,
					Password: passwordAttacker,
				})
				tests.Ok(t, err, signupAttacker)
				session := e2e.WithSessionFromHeader(t, root, signupAttacker.HTTPResponse.Header)

				addVictims, err := cl.AccountEmailAddWithResponse(root, openapi.AccountEmailInitialProps{
					EmailAddress: emailVictim,
				}, session)
				tests.Status(t, err, addVictims, http.StatusConflict)
			})

			t.Run("sign_up_with_username_add_email_login_with_email", func(t *testing.T) {
				handle := xid.New().String()
				password := "password"

				signup, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPasswordSignupJSONRequestBody{
					Identifier: handle,
					Token:      password,
				})
				tests.Ok(t, err, signup)
				session := e2e.WithSessionFromHeader(t, root, signup.HTTPResponse.Header)

				email := xid.New().String() + "@example.com"

				add, err := cl.AccountEmailAddWithResponse(root, openapi.AccountEmailInitialProps{
					EmailAddress: email,
				}, session)
				tests.Ok(t, err, add)

				loginWithEmail1, err := cl.AuthEmailPasswordSigninWithResponse(root, openapi.AuthEmailPasswordSigninJSONRequestBody{
					Email:    email,
					Password: password,
				})
				tests.Ok(t, err, loginWithEmail1)
			})

			t.Run("delete_email", func(t *testing.T) {
				email1 := xid.New().String() + "first@example.com"
				password := "password"

				signup, err := cl.AuthEmailPasswordSignupWithResponse(root, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{
					Email:    email1,
					Password: password,
				})
				tests.Ok(t, err, signup)
				session := e2e.WithSessionFromHeader(t, root, signup.HTTPResponse.Header)

				email2 := xid.New().String() + "second@example.com"

				addEmail2, err := cl.AccountEmailAddWithResponse(root, openapi.AccountEmailInitialProps{
					EmailAddress: email2,
				}, session)
				tests.Ok(t, err, addEmail2)

				remove, err := cl.AccountEmailRemoveWithResponse(root, addEmail2.JSON200.Id, session)
				tests.Ok(t, err, remove)

				loginWithEmail, err := cl.AuthEmailPasswordSigninWithResponse(root, openapi.AuthEmailPasswordSigninJSONRequestBody{
					Email:    email2,
					Password: password,
				})
				tests.Status(t, err, loginWithEmail, http.StatusNotFound)
			})

			t.Run("add_unclaimed_email", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				email := xid.New().String() + "pre-existing@example.com"
				password := "password"

				// Add an email to the database, not associated with any account
				em := db.Email.Create().SetEmailAddress(email).SetVerificationCode("456123").SaveX(root)

				// Sign up with this email, succeeds because the email is not
				// claimed already. Requires verification still.
				signup, err := cl.AuthEmailPasswordSignupWithResponse(root, nil, openapi.AuthEmailPasswordSignupJSONRequestBody{
					Email:    email,
					Password: password,
				})
				tests.Ok(t, err, signup)
				session := e2e.WithSessionFromHeader(t, root, signup.HTTPResponse.Header)

				acc, err := cl.AccountGetWithResponse(root, session)
				tests.Ok(t, err, acc)

				r.Len(acc.JSON200.EmailAddresses, 1)
				a.Equal(email, acc.JSON200.EmailAddresses[0].EmailAddress)
				a.Equal(em.EmailAddress, acc.JSON200.EmailAddresses[0].EmailAddress)
			})
		}))
	}))
}
