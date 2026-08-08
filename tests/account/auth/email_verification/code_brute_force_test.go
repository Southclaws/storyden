package email_verification_test

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/email"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	email_ent "github.com/Southclaws/storyden/internal/ent/email"
	"github.com/Southclaws/storyden/internal/infrastructure/mailer"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

var codePattern = regexp.MustCompile(`verify your account: ([0-9]{6})`)

func signupAndReadCode(t *testing.T, ctx context.Context, cl *openapi.ClientWithResponses, inbox *mailer.Mock, address string) string {
	t.Helper()

	before := inbox.Count()

	// signup doubles as resend, which answers 422 rather than 200
	_, err := cl.AuthEmailSignupWithResponse(ctx, nil, openapi.AuthEmailSignupJSONRequestBody{Email: address})
	require.NoError(t, err)

	verification := tests.WaitForNextEmail(t, inbox, before)

	return codePattern.FindStringSubmatch(verification.Plain)[1]
}

func wrongCode(code string) string {
	if code == "000000" {
		return "111111"
	}
	return "000000"
}

func TestVerificationCodeLocksOutAfterRepeatedGuesses(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		mail mailer.Sender,
	) {
		inbox := mail.(*mailer.Mock)

		lc.Append(fx.StartHook(func() {
			a := assert.New(t)

			address := xid.New().String() + "@storyden.org"
			code := signupAndReadCode(t, root, cl, inbox, address)

			for i := range email.MaxVerificationAttempts {
				resp, err := cl.AuthEmailVerifyWithResponse(root, openapi.AuthEmailVerifyJSONRequestBody{
					Email: address,
					Code:  wrongCode(code),
				})
				require.NoError(t, err)
				a.NotEqual(http.StatusOK, resp.StatusCode(), fmt.Sprintf("guess %d must not succeed", i+1))
			}

			resp, err := cl.AuthEmailVerifyWithResponse(root, openapi.AuthEmailVerifyJSONRequestBody{
				Email: address,
				Code:  code,
			})
			require.NoError(t, err)
			a.NotEqual(http.StatusOK, resp.StatusCode(), "the correct code must be refused once the attempt budget is spent")
		}))
	}))
}

func TestVerificationCodeExpires(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		mail mailer.Sender,
		db *ent.Client,
	) {
		inbox := mail.(*mailer.Mock)

		lc.Append(fx.StartHook(func() {
			a := assert.New(t)
			r := require.New(t)

			address := xid.New().String() + "@storyden.org"
			code := signupAndReadCode(t, root, cl, inbox, address)

			_, err := db.Email.Update().
				Where(email_ent.EmailAddress(address)).
				SetVerificationCodeExpiresAt(time.Now().Add(-time.Minute)).
				Save(root)
			r.NoError(err)

			resp, err := cl.AuthEmailVerifyWithResponse(root, openapi.AuthEmailVerifyJSONRequestBody{
				Email: address,
				Code:  code,
			})
			r.NoError(err)
			a.NotEqual(http.StatusOK, resp.StatusCode(), "an expired code must not verify")
		}))
	}))
}

func TestResendRotatesTheCodeAndClearsAttempts(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		mail mailer.Sender,
	) {
		inbox := mail.(*mailer.Mock)

		lc.Append(fx.StartHook(func() {
			a := assert.New(t)
			r := require.New(t)

			address := xid.New().String() + "@storyden.org"
			first := signupAndReadCode(t, root, cl, inbox, address)

			for range email.MaxVerificationAttempts {
				_, err := cl.AuthEmailVerifyWithResponse(root, openapi.AuthEmailVerifyJSONRequestBody{
					Email: address,
					Code:  wrongCode(first),
				})
				r.NoError(err)
			}

			second := signupAndReadCode(t, root, cl, inbox, address)
			a.NotEqual(first, second, "resending must issue a new code, not repeat the burned one")

			stale, err := cl.AuthEmailVerifyWithResponse(root, openapi.AuthEmailVerifyJSONRequestBody{
				Email: address,
				Code:  first,
			})
			r.NoError(err)
			a.NotEqual(http.StatusOK, stale.StatusCode(), "the replaced code must stop working")

			resp, err := cl.AuthEmailVerifyWithResponse(root, openapi.AuthEmailVerifyJSONRequestBody{
				Email: address,
				Code:  second,
			})
			r.NoError(err)
			a.Equal(http.StatusOK, resp.StatusCode(), "a fresh code must work again after a lockout")
		}))
	}))
}
