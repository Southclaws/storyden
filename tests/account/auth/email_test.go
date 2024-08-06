package auth_test

import (
	"context"
	"net/http"
	"regexp"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	session1 "github.com/Southclaws/storyden/app/transports/http/middleware/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/infrastructure/mailer"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestEmailOnlyAuth(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		cj *session1.Jar,
		accountQuery account_querier.Querier,
		mail mailer.Sender,
	) {
		inbox := mail.(*mailer.Mock)

		lc.Append(fx.StartHook(func() {
			t.Run("verify_success", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				address := xid.New().String() + "@storyden.org"

				// Sign up with email
				signup, err := cl.AuthEmailSignupWithResponse(root, openapi.AuthEmailSignupJSONRequestBody{Email: address})
				tests.Ok(t, err, signup)

				accountID := account.AccountID(openapi.GetAccountID(signup.JSON200.Id))
				ctx1 := session.WithAccountID(root, accountID)
				session := e2e.WithSession(ctx1, cj)

				// Get own account, currently unverified
				unverified, err := cl.AccountGetWithResponse(root, session)
				tests.Ok(t, err, unverified)
				r.Equal(openapi.AccountVerifiedStatusNone, unverified.JSON200.VerifiedStatus)
				r.Len(unverified.JSON200.EmailAddresses, 1)
				a.Equal(address, (unverified.JSON200.EmailAddresses)[0].EmailAddress)
				a.True(unverified.JSON200.EmailAddresses[0].IsAuth)
				a.False(unverified.JSON200.EmailAddresses[0].Verified)

				// Get code from email, verify account
				verification := inbox.GetLast()
				code := regexp.MustCompile(`verify your account: ([0-9]{6})`).FindStringSubmatch(verification.Plain)[1]
				verify, err := cl.AuthEmailVerifyWithResponse(root, openapi.AuthEmailVerifyJSONRequestBody{Email: address, Code: code}, session)
				tests.Ok(t, err, verify)
				a.Equal(accountID.String(), verify.JSON200.Id)

				// Get own account, now verified
				verified, err := cl.AccountGetWithResponse(root, session)
				tests.Ok(t, err, verified)
				a.Equal(openapi.AccountVerifiedStatusVerifiedEmail, verified.JSON200.VerifiedStatus)
				r.NotNil(verified.JSON200.EmailAddresses)
				a.Equal(address, verified.JSON200.EmailAddresses[0].EmailAddress)
				a.True(verified.JSON200.EmailAddresses[0].IsAuth)
				a.True(verified.JSON200.EmailAddresses[0].Verified)
			})

			t.Run("verify_resend", func(t *testing.T) {
				// r := require.New(t)
				a := assert.New(t)

				address := xid.New().String() + "@storyden.org"

				// Sign up with email
				signup, err := cl.AuthEmailSignupWithResponse(root, openapi.AuthEmailSignupJSONRequestBody{Email: address})
				tests.Ok(t, err, signup)

				// Sign up with email, again, resulting in a 202 Accepted and no cookie session
				signup2, err := cl.AuthEmailSignupWithResponse(root, openapi.AuthEmailSignupJSONRequestBody{Email: address})
				tests.Status(t, err, signup2, http.StatusUnprocessableEntity)

				a.Empty(signup2.HTTPResponse.Header.Get("Set-Cookie"))
			})

			t.Run("verify_wrong_code", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				address := xid.New().String() + "@storyden.org"

				// Sign up with email
				signup, err := cl.AuthEmailSignupWithResponse(root, openapi.AuthEmailSignupJSONRequestBody{Email: address})
				tests.Ok(t, err, signup)

				accountID := account.AccountID(openapi.GetAccountID(signup.JSON200.Id))
				ctx1 := session.WithAccountID(root, accountID)
				session := e2e.WithSession(ctx1, cj)

				// Get own account, currently unverified
				unverified, err := cl.AccountGetWithResponse(root, session)
				tests.Ok(t, err, unverified)
				r.Equal(openapi.AccountVerifiedStatusNone, unverified.JSON200.VerifiedStatus)

				incorrectCode := "999999" // one day, this test will fail...
				verify, err := cl.AuthEmailVerifyWithResponse(root, openapi.AuthEmailVerifyJSONRequestBody{Email: address, Code: incorrectCode}, session)
				tests.Status(t, err, verify, http.StatusForbidden)

				// Get own account, still not verified
				verified, err := cl.AccountGetWithResponse(root, session)
				tests.Ok(t, err, verified)
				a.Equal(openapi.AccountVerifiedStatusNone, verified.JSON200.VerifiedStatus)
			})
		}))
	}))
}
