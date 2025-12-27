package thread_test

import (
	"context"
	"testing"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestThreadListReviewVisibility(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			adminCtx, adminAcc := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			memberCtx, memberAcc := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			otherCtx, _ := e2e.WithAccount(root, aw, seed.Account_002_Frigg)

			sessionAdmin := sh.WithSession(adminCtx)
			sessionMember := sh.WithSession(memberCtx)
			sessionOther := sh.WithSession(otherCtx)

			publishedThread, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Body:       opt.New("<p>Published thread</p>").Ptr(),
				Title:      "Published Thread",
				Visibility: opt.New(openapi.Published).Ptr(),
			}, sessionMember)
			tests.Ok(t, err, publishedThread)

			reviewThread, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Body:       opt.New("<p>Review thread</p>").Ptr(),
				Title:      "Review Thread",
				Visibility: opt.New(openapi.Review).Ptr(),
			}, sessionMember)
			tests.Ok(t, err, reviewThread)

			otherReviewThread, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Body:       opt.New("<p>Other review thread</p>").Ptr(),
				Title:      "Other Review Thread",
				Visibility: opt.New(openapi.Review).Ptr(),
			}, sessionOther)
			tests.Ok(t, err, otherReviewThread)

			t.Run("moderator_sees_all_review_threads_with_visibility_filter", func(t *testing.T) {
				threadList, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{
					Visibility: &[]openapi.Visibility{openapi.Published, openapi.Review},
				}, sessionAdmin)
				tests.Ok(t, err, threadList)

				ids := dt.Map(threadList.JSON200.Threads, func(th openapi.ThreadReference) string { return th.Id })
				a.Contains(ids, publishedThread.JSON200.Id, "moderator should see published thread")
				a.Contains(ids, reviewThread.JSON200.Id, "moderator should see member's review thread")
				a.Contains(ids, otherReviewThread.JSON200.Id, "moderator should see other's review thread")
			})

			t.Run("moderator_sees_all_review_threads_in_default_feed", func(t *testing.T) {
				threadList, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{}, sessionAdmin)
				tests.Ok(t, err, threadList)

				ids := dt.Map(threadList.JSON200.Threads, func(th openapi.ThreadReference) string { return th.Id })
				a.Contains(ids, publishedThread.JSON200.Id, "moderator should see published thread in default feed")
				a.Contains(ids, reviewThread.JSON200.Id, "moderator should see member's review thread in default feed")
				a.Contains(ids, otherReviewThread.JSON200.Id, "moderator should see other's review thread in default feed")
			})

			t.Run("member_sees_only_own_review_threads_in_default_listing", func(t *testing.T) {
				threadList, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{}, sessionMember)
				tests.Ok(t, err, threadList)

				ids := dt.Map(threadList.JSON200.Threads, func(th openapi.ThreadReference) string { return th.Id })
				a.Contains(ids, publishedThread.JSON200.Id, "member should see published thread")
				a.Contains(ids, reviewThread.JSON200.Id, "member should see their own review thread in default listing")
				a.NotContains(ids, otherReviewThread.JSON200.Id, "member should NOT see other's review thread")
			})

			t.Run("member_sees_own_review_threads_when_filtering_by_own_account", func(t *testing.T) {
				threadList, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{
					Visibility: &[]openapi.Visibility{openapi.Published, openapi.Review},
					Author:     &memberAcc.Handle,
				}, sessionMember)
				tests.Ok(t, err, threadList)

				ids := dt.Map(threadList.JSON200.Threads, func(th openapi.ThreadReference) string { return th.Id })
				a.Contains(ids, publishedThread.JSON200.Id, "member should see their own published thread")
				a.Contains(ids, reviewThread.JSON200.Id, "member should see their own review thread")
				a.NotContains(ids, otherReviewThread.JSON200.Id, "member should NOT see other's review thread even when filtering by own account")
			})

			t.Run("member_cannot_see_others_review_threads_when_filtering_by_other_account", func(t *testing.T) {
				threadList, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{
					Visibility: &[]openapi.Visibility{openapi.Published, openapi.Review},
					Author:     &adminAcc.Handle,
				}, sessionMember)
				tests.Ok(t, err, threadList)

				adminThreads := dt.Filter(threadList.JSON200.Threads, func(th openapi.ThreadReference) bool {
					return th.Author.Id == adminAcc.ID.String()
				})

				for _, th := range adminThreads {
					a.NotEqual(openapi.Review, th.Visibility, "non-moderator member should not see admin's review threads")
				}
			})

			t.Run("unauthenticated_user_only_sees_published_threads", func(t *testing.T) {
				threadList, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{
					Visibility: &[]openapi.Visibility{openapi.Published, openapi.Review},
				})
				tests.Ok(t, err, threadList)

				ids := dt.Map(threadList.JSON200.Threads, func(th openapi.ThreadReference) string { return th.Id })
				a.Contains(ids, publishedThread.JSON200.Id, "unauthenticated user should see published thread")
				a.NotContains(ids, reviewThread.JSON200.Id, "unauthenticated user should NOT see review thread")
				a.NotContains(ids, otherReviewThread.JSON200.Id, "unauthenticated user should NOT see other review thread")
			})

			t.Run("member_filtering_by_account_with_only_review_visibility", func(t *testing.T) {
				threadList, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{
					Visibility: &[]openapi.Visibility{openapi.Review},
					Author:     &memberAcc.Handle,
				}, sessionMember)
				tests.Ok(t, err, threadList)

				threads := dt.Filter(threadList.JSON200.Threads, func(th openapi.ThreadReference) bool {
					return lo.Contains([]string{reviewThread.JSON200.Id, publishedThread.JSON200.Id, otherReviewThread.JSON200.Id}, th.Id)
				})

				r.Len(threads, 1, "should only return the member's review thread")
				a.Equal(reviewThread.JSON200.Id, threads[0].Id)
				a.Equal(openapi.Review, threads[0].Visibility)
			})
		}))
	}))
}
