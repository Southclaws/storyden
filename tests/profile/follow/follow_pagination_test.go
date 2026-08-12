package follow_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/profile/follow_querier"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestFollowListsPageWithoutLosingRows(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		ar *account_querier.Querier,
		fq *follow_querier.Querier,
	) {
		lc.Append(fx.StartHook(func() {
			t.Run("followers", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				target := newAccount(t, root, cl, ar, "tgt")

				const followerCount = 3
				for i := range followerCount {
					follower := newAccount(t, root, cl, ar, fmt.Sprintf("fwr%d", i))
					session := sh.WithSession(e2e.WithAccountID(root, follower.ID))

					add, err := cl.ProfileFollowersAddWithResponse(root, target.Handle, session)
					tests.Ok(t, err, add)
				}

				list, err := cl.ProfileFollowersGetWithResponse(root, target.Handle, &openapi.ProfileFollowersGetParams{})
				tests.Ok(t, err, list)
				r.Len(list.JSON200.Followers, followerCount)
				a.Equal(followerCount, list.JSON200.Results)
				a.Nil(list.JSON200.NextPage)

				seen := walkPages(t, func(page uint) (*pagination.Result[*profile.Ref], error) {
					return fq.GetFollowers(root, target.ID, pagination.NewPageParams(page, 2))
				})
				a.Len(seen, followerCount)
			})

			t.Run("following", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				follower := newAccount(t, root, cl, ar, "fwr")
				session := sh.WithSession(e2e.WithAccountID(root, follower.ID))

				const followingCount = 3
				for i := range followingCount {
					target := newAccount(t, root, cl, ar, fmt.Sprintf("tgt%d", i))

					add, err := cl.ProfileFollowersAddWithResponse(root, target.Handle, session)
					tests.Ok(t, err, add)
				}

				list, err := cl.ProfileFollowingGetWithResponse(root, follower.Handle, &openapi.ProfileFollowingGetParams{})
				tests.Ok(t, err, list)
				r.Len(list.JSON200.Following, followingCount)
				a.Equal(followingCount, list.JSON200.Results)
				a.Nil(list.JSON200.NextPage)

				seen := walkPages(t, func(page uint) (*pagination.Result[*profile.Ref], error) {
					return fq.GetFollowing(root, follower.ID, pagination.NewPageParams(page, 2))
				})
				a.Len(seen, followingCount)
			})
		}))
	}))
}

func walkPages(t *testing.T, get func(page uint) (*pagination.Result[*profile.Ref], error)) map[account.AccountID]struct{} {
	r := require.New(t)
	a := assert.New(t)

	seen := map[account.AccountID]struct{}{}

	for page := uint(1); page <= 10; page++ {
		result, err := get(page)
		r.NoError(err)
		r.Equal(int(page), result.CurrentPage)

		for _, p := range result.Items {
			_, duplicate := seen[p.ID]
			a.False(duplicate, "profile %s was returned on more than one page", p.ID)
			seen[p.ID] = struct{}{}
		}

		next, ok := result.NextPage.Get()
		if !ok {
			return seen
		}

		a.Equal(int(page)+1, next)
	}

	t.Fatal("pagination did not terminate")

	return seen
}