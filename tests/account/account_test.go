package account_test

import (
	"context"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/openapi/bindings"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/openapi"
)

func TestAccountAuth(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		cj *bindings.CookieJar,
		ar account.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			hand1 := "tester1-" + xid.New().String()

			// Sign up for a new account with a password

			a1, err := cl.AuthPasswordSignupWithResponse(root, openapi.AuthPair{
				Identifier: hand1,
				Token:      "password",
			})
			r.NoError(err)
			r.NotNil(a1)
			r.Equal(200, a1.StatusCode())

			id1, err := xid.FromString(a1.JSON200.Id)
			r.NoError(err)

			ctx1 := session.WithAccountID(root, account.AccountID(id1))

			// Get the new account

			get1, err := cl.AccountGetWithResponse(root, e2e.WithSession(ctx1, cj))
			r.NoError(err)
			r.NotNil(get1)
			r.Equal(200, get1.StatusCode())

			a.Equal(hand1, get1.JSON200.Handle)

			// Change the password

			change1, err := cl.AuthPasswordUpdateWithResponse(root, openapi.AuthPasswordMutableProps{
				Old: "password",
				New: "wordpass",
			}, e2e.WithSession(ctx1, cj))
			r.NoError(err)
			r.NotNil(change1)
			r.Equal(200, change1.StatusCode())

			// Log in to the new account with the old password - fails

			signin1, err := cl.AuthPasswordSigninWithResponse(root, openapi.AuthPair{
				Identifier: hand1,
				Token:      "password",
			})
			r.NoError(err)
			r.NotNil(signin1)
			r.Equal(403, signin1.StatusCode())

			// Log in to the new account with the new password - succeeds

			signin2, err := cl.AuthPasswordSigninWithResponse(root, openapi.AuthPair{
				Identifier: hand1,
				Token:      "wordpass",
			})
			r.NoError(err)
			r.NotNil(signin2)
			r.Equal(200, signin2.StatusCode())
		}))
	}))
}
