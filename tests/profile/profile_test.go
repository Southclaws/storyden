package account_test

import (
	"context"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/transports/openapi/bindings"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/openapi"
	"github.com/Southclaws/storyden/internal/utils"
)

func TestPublicProfiles(t *testing.T) {
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

			// Create 5 fresh accounts
			newAccount(t, root, cl, ar, "odin")
			newAccount(t, root, cl, ar, "frigg")
			newAccount(t, root, cl, ar, "baldur")
			newAccount(t, root, cl, ar, "odin2")
			newAccount(t, root, cl, ar, "þórr")

			// Get them all, default params.

			list1, err := cl.ProfileListWithResponse(root, nil)
			r.NoError(err)
			r.NotNil(list1)
			r.Equal(200, list1.StatusCode())

			a.Equal(1, list1.JSON200.CurrentPage)
			a.GreaterOrEqual(len(list1.JSON200.Profiles), 5)

			// Get one specific name - there are two odins

			list2, err := cl.ProfileListWithResponse(root, &openapi.ProfileListParams{
				Q: opt.New("odin").Ptr(),
			})
			r.NoError(err)
			r.NotNil(list2)
			r.Equal(200, list2.StatusCode())

			a.Equal(1, list1.JSON200.CurrentPage)
			a.Equal(50, list1.JSON200.PageSize)
			a.GreaterOrEqual(len(list2.JSON200.Profiles), 2)

			// Query an invalid page.

			list3, err := cl.ProfileListWithResponse(root, &openapi.ProfileListParams{
				Page: opt.New("2147483647").Ptr(),
			})
			r.NoError(err)
			r.NotNil(list3)
			r.Equal(200, list3.StatusCode())

			a.Equal(2147483647, list3.JSON200.CurrentPage)
			a.Nil(list3.JSON200.NextPage)
			a.Empty(list3.JSON200.Profiles)
		}))
	}))
}

func newAccount(t *testing.T, ctx context.Context, cl *openapi.ClientWithResponses, ar account.Repository, handle string) account.Account {
	r := require.New(t)

	hand1 := handle + "-" + xid.New().String()

	response, err := cl.AuthPasswordSignupWithResponse(ctx, openapi.AuthPair{
		Identifier: hand1,
		Token:      "password",
	})
	r.NoError(err)
	r.NotNil(response)
	r.Equal(200, response.StatusCode())

	acc, err := ar.GetByID(ctx, account.AccountID(utils.Must(xid.FromString(response.JSON200.Id))))
	r.NoError(err)
	r.NotNil(acc)

	return *acc
}
