package account_test

import (
	"context"
	"net/http"
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
	"github.com/Southclaws/storyden/internal/utils"
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

func TestAccountAdmin(t *testing.T) {
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

			adminHandle := "tester-admin-" + xid.New().String()
			victimHandle := "tester-victim-" + xid.New().String()
			randomHandle := "tester-random-" + xid.New().String()

			// Sign up for a new account with a password

			admin, err := cl.AuthPasswordSignupWithResponse(root, openapi.AuthPair{Identifier: adminHandle, Token: "password"})
			r.NoError(err)
			r.Equal(200, admin.StatusCode())
			adminID := account.AccountID(utils.Must(xid.FromString(admin.JSON200.Id)))
			adminSession := e2e.WithSession(session.WithAccountID(root, adminID), cj)

			ar.Update(root, adminID, account.SetAdmin(true))

			victim, err := cl.AuthPasswordSignupWithResponse(root, openapi.AuthPair{Identifier: victimHandle, Token: "password"})
			r.NoError(err)
			r.Equal(200, victim.StatusCode())
			victimID := account.AccountID(utils.Must(xid.FromString(victim.JSON200.Id)))
			victimSession := e2e.WithSession(session.WithAccountID(root, victimID), cj)

			random, err := cl.AuthPasswordSignupWithResponse(root, openapi.AuthPair{Identifier: randomHandle, Token: "password"})
			r.NoError(err)
			r.Equal(200, random.StatusCode())
			randomID := account.AccountID(utils.Must(xid.FromString(random.JSON200.Id)))
			randomSession := e2e.WithSession(session.WithAccountID(root, randomID), cj)

			// Try to suspend the account without being logged in - fails

			suspend1, err := cl.AdminAccountBanCreateWithResponse(root, victim.JSON200.Id)
			r.NoError(err)
			r.NotNil(suspend1)
			r.Equal(http.StatusForbidden, suspend1.StatusCode())

			// Try to suspend the account as a non-admin - fails

			suspend2, err := cl.AdminAccountBanCreateWithResponse(root, victim.JSON200.Id, randomSession)
			r.NoError(err)
			r.NotNil(suspend2)
			r.Equal(http.StatusUnauthorized, suspend2.StatusCode())

			// Try to suspend the account as an admin - succeeds

			suspend3, err := cl.AdminAccountBanCreateWithResponse(root, victim.JSON200.Id, adminSession)
			r.NoError(err)
			r.NotNil(suspend3)
			r.Equal(http.StatusOK, suspend3.StatusCode())

			victimsigni1, err := cl.AuthPasswordSigninWithResponse(root, openapi.AuthPair{
				Identifier: victimHandle,
				Token:      "password",
			}, victimSession)
			r.NoError(err)
			r.NotNil(victimsigni1)
			r.Equal(http.StatusUnauthorized, victimsigni1.StatusCode())

			// Try to reinstate the account without being logged in - fails

			reinstate1, err := cl.AdminAccountBanRemoveWithResponse(root, victim.JSON200.Id)
			r.NoError(err)
			r.NotNil(reinstate1)
			r.Equal(http.StatusForbidden, reinstate1.StatusCode())

			// Try to reinstate the account as a non-admin - fails

			reinstate2, err := cl.AdminAccountBanRemoveWithResponse(root, victim.JSON200.Id, randomSession)
			r.NoError(err)
			r.NotNil(reinstate2)
			r.Equal(http.StatusUnauthorized, reinstate2.StatusCode())

			// Try to reinstate the account as an admin - succeeds

			reinstate3, err := cl.AdminAccountBanRemoveWithResponse(root, victim.JSON200.Id, adminSession)
			r.NoError(err)
			r.NotNil(reinstate3)
			r.Equal(http.StatusOK, reinstate3.StatusCode())

			victimsignin2, err := cl.AuthPasswordSigninWithResponse(root, openapi.AuthPair{
				Identifier: victimHandle,
				Token:      "password",
			}, victimSession)
			r.NoError(err)
			r.NotNil(victimsignin2)
			r.Equal(http.StatusOK, victimsignin2.StatusCode())
		}))
	}))
}
