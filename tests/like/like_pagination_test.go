package like_test

import (
	"context"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/like/like_querier"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestProfileLikesPageWithoutLosingRows(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		lq *like_querier.LikeQuerier,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			likerCtx, likerAcc := e2e.WithAccount(root, aw, seed.Account_006_Freyja)
			adminSession := sh.WithSession(adminCtx)
			likerSession := sh.WithSession(likerCtx)

			const likeCount = 3
			for range likeCount {
				thread, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       opt.New("<p>likeable</p>").Ptr(),
					Visibility: opt.New(openapi.VisibilityPublished).Ptr(),
					Title:      "like paging " + uuid.NewString(),
				}, adminSession)
				tests.Ok(t, err, thread)

				like, err := cl.LikePostAddWithResponse(root, thread.JSON200.Id, likerSession)
				tests.Ok(t, err, like)
			}

			list, err := cl.LikeProfileGetWithResponse(root, likerAcc.Handle, &openapi.LikeProfileGetParams{}, adminSession)
			tests.Ok(t, err, list)
			r.Len(list.JSON200.Likes, likeCount)
			a.Equal(likeCount, list.JSON200.Results)
			a.Nil(list.JSON200.NextPage)

			own, err := cl.LikeProfileGetWithResponse(root, likerAcc.Handle, &openapi.LikeProfileGetParams{}, likerSession)
			tests.Ok(t, err, own)
			a.Len(own.JSON200.Likes, likeCount)

			seen := map[xid.ID]struct{}{}

			for page := uint(1); page <= 10; page++ {
				result, err := lq.GetProfileLikes(root, likerAcc.ID, pagination.NewPageParams(page, 2))
				r.NoError(err)
				r.Equal(int(page), result.CurrentPage)

				for _, like := range result.Items {
					_, duplicate := seen[like.ID]
					a.False(duplicate, "like %s was returned on more than one page", like.ID)
					seen[like.ID] = struct{}{}
				}

				next, ok := result.NextPage.Get()
				if !ok {
					break
				}

				a.Equal(int(page)+1, next)
			}

			a.Len(seen, likeCount)
		}))
	}))
}
