package thread_test

import (
	"context"
	"testing"

	"github.com/Southclaws/opt"
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

func TestThreadVisibilityWithReviewReplies(t *testing.T) {
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

			memberCtx, _ := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			sessionMember := sh.WithSession(memberCtx)
			sessionAdmin := sh.WithSession(adminCtx)

			threadCreate, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Body:       opt.New("<p>This is a test thread</p>").Ptr(),
				Title:      "Test Thread for Visibility",
				Visibility: opt.New(openapi.Published).Ptr(),
			}, sessionMember)
			tests.Ok(t, err, threadCreate)

			reply1Create, err := cl.ReplyCreateWithResponse(root, threadCreate.JSON200.Slug, openapi.ReplyInitialProps{
				Body: "This is a published reply",
			}, sessionMember)
			tests.Ok(t, err, reply1Create)

			reply2Create, err := cl.ReplyCreateWithResponse(root, threadCreate.JSON200.Slug, openapi.ReplyInitialProps{
				Body: "This is another reply",
			}, sessionMember)
			tests.Ok(t, err, reply2Create)

			_, err = cl.PostUpdateWithResponse(root, reply2Create.JSON200.Id, openapi.PostMutableProps{
				Visibility: opt.New(openapi.Review).Ptr(),
			}, sessionAdmin)
			r.NoError(err)

			t.Run("member_sees_published_and_own_review_replies", func(t *testing.T) {
				a := assert.New(t)

				threadGet, err := cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, nil, sessionMember)
				tests.Ok(t, err, threadGet)

				r.Len(threadGet.JSON200.Replies.Replies, 2, "member should see published reply and their own in-review reply")

				// Find both replies
				var publishedReply, reviewReply *openapi.Reply
				for i := range threadGet.JSON200.Replies.Replies {
					if threadGet.JSON200.Replies.Replies[i].Id == reply1Create.JSON200.Id {
						publishedReply = &threadGet.JSON200.Replies.Replies[i]
					}
					if threadGet.JSON200.Replies.Replies[i].Id == reply2Create.JSON200.Id {
						reviewReply = &threadGet.JSON200.Replies.Replies[i]
					}
				}

				r.NotNil(publishedReply, "should find published reply")
				r.NotNil(reviewReply, "should find review reply (own)")
				a.Equal(openapi.Published, publishedReply.Visibility)
				a.Equal(openapi.Review, reviewReply.Visibility)
			})

			t.Run("admin_sees_both_replies", func(t *testing.T) {
				a := assert.New(t)

				threadGet, err := cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, nil, sessionAdmin)
				tests.Ok(t, err, threadGet)

				r.Len(threadGet.JSON200.Replies.Replies, 2, "admin with ManagePosts should see both replies")

				replyIDs := []string{
					threadGet.JSON200.Replies.Replies[0].Id,
					threadGet.JSON200.Replies.Replies[1].Id,
				}
				a.Contains(replyIDs, reply1Create.JSON200.Id, "admin should see published reply")
				a.Contains(replyIDs, reply2Create.JSON200.Id, "admin should see review reply")

				var reviewReply *openapi.Reply
				for i := range threadGet.JSON200.Replies.Replies {
					if threadGet.JSON200.Replies.Replies[i].Id == reply2Create.JSON200.Id {
						reviewReply = &threadGet.JSON200.Replies.Replies[i]
						break
					}
				}
				r.NotNil(reviewReply, "should find the review reply")
				a.Equal(openapi.Review, reviewReply.Visibility, "reply should be marked as review")
			})

			t.Run("admin_can_accept_review_reply", func(t *testing.T) {
				a := assert.New(t)

				replyUpdate, err := cl.PostUpdateWithResponse(root, reply2Create.JSON200.Id, openapi.PostMutableProps{
					Visibility: opt.New(openapi.Published).Ptr(),
				}, sessionAdmin)
				tests.Ok(t, err, replyUpdate)
				a.Equal(openapi.Published, replyUpdate.JSON200.Visibility, "reply should now be published")

				threadGet, err := cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, nil, sessionMember)
				tests.Ok(t, err, threadGet)
				r.Len(threadGet.JSON200.Replies.Replies, 2, "member should now see both replies after acceptance")
			})

			t.Run("member_cannot_change_reply_visibility", func(t *testing.T) {
				_, err := cl.PostUpdateWithResponse(root, reply2Create.JSON200.Id, openapi.PostMutableProps{
					Visibility: opt.New(openapi.Review).Ptr(),
				}, sessionAdmin)
				r.NoError(err)

				replyUpdate, err := cl.PostUpdateWithResponse(root, reply2Create.JSON200.Id, openapi.PostMutableProps{
					Visibility: opt.New(openapi.Published).Ptr(),
				}, sessionMember)
				r.NoError(err)
				r.Equal(403, replyUpdate.StatusCode(), "member should not be able to change visibility")
			})
		}))
	}))
}
