package follow_test

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
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/tests"
)

func TestFollows(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		ar *account_querier.Querier,
	) {
		lc.Append(fx.StartHook(func() {
			t.Run("one_follows_another", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				acc1 := newAccount(t, root, cl, ar, "acc1")
				acc2 := newAccount(t, root, cl, ar, "acc2")
				acc1session := sh.WithSession(session.WithAccountID(root, acc1.ID))

				// Follow acc2 from acc1
				f1, err := cl.ProfileFollowersAddWithResponse(root, acc2.Handle, acc1session)
				tests.Ok(t, err, f1)

				// Get following for acc1
				acc1following, err := cl.ProfileFollowingGetWithResponse(root, acc1.Handle, &openapi.ProfileFollowingGetParams{})
				tests.Ok(t, err, acc1following)
				r.Len(acc1following.JSON200.Following, 1, "acc1 is following acc2")
				r.Equal(acc1following.JSON200.Results, len(acc1following.JSON200.Following))
				a.Equal(acc2.ID.String(), acc1following.JSON200.Following[0].Id)

				// Get followers for acc1
				acc1followers, err := cl.ProfileFollowersGetWithResponse(root, acc1.Handle, &openapi.ProfileFollowersGetParams{})
				tests.Ok(t, err, acc1followers)
				r.Len(acc1followers.JSON200.Followers, 0)
				r.Equal(acc1followers.JSON200.Results, len(acc1followers.JSON200.Followers))

				// Get following for acc2
				acc2following, err := cl.ProfileFollowingGetWithResponse(root, acc2.Handle, &openapi.ProfileFollowingGetParams{})
				tests.Ok(t, err, acc2following)
				r.Len(acc2following.JSON200.Following, 0)
				r.Equal(acc2following.JSON200.Results, len(acc2following.JSON200.Following))

				// Get followers for acc2
				acc2followers, err := cl.ProfileFollowersGetWithResponse(root, acc2.Handle, &openapi.ProfileFollowersGetParams{})
				tests.Ok(t, err, acc2followers)
				r.Len(acc2followers.JSON200.Followers, 1)
				r.Equal(acc2followers.JSON200.Results, len(acc2followers.JSON200.Followers))
				a.Equal(acc1.ID.String(), acc2followers.JSON200.Followers[0].Id)

				profile1, err := cl.ProfileGetWithResponse(root, acc1.Handle)
				tests.Ok(t, err, profile1)
				r.Equal(0, profile1.JSON200.Followers)
				r.Equal(1, profile1.JSON200.Following)

				profile2, err := cl.ProfileGetWithResponse(root, acc2.Handle)
				tests.Ok(t, err, profile2)
				r.Equal(1, profile2.JSON200.Followers)
				r.Equal(0, profile2.JSON200.Following)
			})

			t.Run("both_follow_eachother", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				acc1 := newAccount(t, root, cl, ar, "acc1")
				acc2 := newAccount(t, root, cl, ar, "acc2")
				acc1session := sh.WithSession(session.WithAccountID(root, acc1.ID))
				acc2session := sh.WithSession(session.WithAccountID(root, acc2.ID))

				// Follow acc2 from acc1
				f1, err := cl.ProfileFollowersAddWithResponse(root, acc2.Handle, acc1session)
				tests.Ok(t, err, f1)

				// Follow acc2 from acc1
				f2, err := cl.ProfileFollowersAddWithResponse(root, acc1.Handle, acc2session)
				tests.Ok(t, err, f2)

				// Get following for acc1
				acc1following, err := cl.ProfileFollowingGetWithResponse(root, acc1.Handle, &openapi.ProfileFollowingGetParams{})
				tests.Ok(t, err, acc1following)
				r.Len(acc1following.JSON200.Following, 1)
				r.Equal(acc1following.JSON200.Results, len(acc1following.JSON200.Following))
				a.Equal(acc2.ID.String(), acc1following.JSON200.Following[0].Id)

				// Get followers for acc1
				acc1followers, err := cl.ProfileFollowersGetWithResponse(root, acc1.Handle, &openapi.ProfileFollowersGetParams{})
				tests.Ok(t, err, acc1followers)
				r.Len(acc1followers.JSON200.Followers, 1)
				r.Equal(acc1followers.JSON200.Results, len(acc1followers.JSON200.Followers))

				// Get following for acc2
				acc2following, err := cl.ProfileFollowingGetWithResponse(root, acc2.Handle, &openapi.ProfileFollowingGetParams{})
				tests.Ok(t, err, acc2following)
				r.Len(acc2following.JSON200.Following, 1)
				r.Equal(acc2following.JSON200.Results, len(acc2following.JSON200.Following))

				// Get followers for acc2
				acc2followers, err := cl.ProfileFollowersGetWithResponse(root, acc2.Handle, &openapi.ProfileFollowersGetParams{})
				tests.Ok(t, err, acc2followers)
				r.Len(acc2followers.JSON200.Followers, 1)
				r.Equal(acc2followers.JSON200.Results, len(acc2followers.JSON200.Followers))
				a.Equal(acc1.ID.String(), acc2followers.JSON200.Followers[0].Id)
			})

			t.Run("follow_unfollow", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				acc1 := newAccount(t, root, cl, ar, "acc1")
				acc2 := newAccount(t, root, cl, ar, "acc2")
				acc1session := sh.WithSession(session.WithAccountID(root, acc1.ID))

				// Follow acc2 from acc1
				f1, err := cl.ProfileFollowersAddWithResponse(root, acc2.Handle, acc1session)
				tests.Ok(t, err, f1)

				// Get following for acc1
				acc1following, err := cl.ProfileFollowingGetWithResponse(root, acc1.Handle, &openapi.ProfileFollowingGetParams{})
				tests.Ok(t, err, acc1following)
				r.Len(acc1following.JSON200.Following, 1, "acc1 is following acc2")
				r.Equal(acc1following.JSON200.Results, len(acc1following.JSON200.Following))
				a.Equal(acc2.ID.String(), acc1following.JSON200.Following[0].Id)

				// Get followers for acc2
				acc2followers, err := cl.ProfileFollowersGetWithResponse(root, acc2.Handle, &openapi.ProfileFollowersGetParams{})
				tests.Ok(t, err, acc2followers)
				r.Len(acc2followers.JSON200.Followers, 1)
				r.Equal(acc2followers.JSON200.Results, len(acc2followers.JSON200.Followers))
				a.Equal(acc1.ID.String(), acc2followers.JSON200.Followers[0].Id)

				u1, err := cl.ProfileFollowersRemoveWithResponse(root, acc2.Handle, acc1session)
				tests.Ok(t, err, u1)

				// Get following for acc1
				acc1following, err = cl.ProfileFollowingGetWithResponse(root, acc1.Handle, &openapi.ProfileFollowingGetParams{})
				tests.Ok(t, err, acc1following)
				r.Len(acc1following.JSON200.Following, 0)
				r.Equal(acc1following.JSON200.Results, len(acc1following.JSON200.Following))

				// Get followers for acc2
				acc2followers, err = cl.ProfileFollowersGetWithResponse(root, acc2.Handle, &openapi.ProfileFollowersGetParams{})
				tests.Ok(t, err, acc2followers)
				r.Len(acc2followers.JSON200.Followers, 0)
				r.Equal(acc2followers.JSON200.Results, len(acc2followers.JSON200.Followers))
			})
		}))
	}))
}

func newAccount(t *testing.T, ctx context.Context, cl *openapi.ClientWithResponses, ar *account_querier.Querier, handle string) account.Account {
	r := require.New(t)

	hand1 := handle + "-" + xid.New().String()

	response, err := cl.AuthPasswordSignupWithResponse(ctx, nil, openapi.AuthPair{
		Identifier: hand1,
		Token:      "password",
	})
	r.NoError(err)
	r.NotNil(response)
	r.Equal(http.StatusOK, response.StatusCode())

	acc, err := ar.GetByID(ctx, account.AccountID(utils.Must(xid.FromString(response.JSON200.Id))))
	r.NoError(err)
	r.NotNil(acc)

	return acc.Account
}
