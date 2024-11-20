package like_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestLikeThreads(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		cj *session_cookie.Jar,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			user1Ctx, user1Acc := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			user2Ctx, user2Acc := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			adminSession := e2e.WithSession(adminCtx, cj)
			user1Session := e2e.WithSession(user1Ctx, cj)
			user2Session := e2e.WithSession(user2Ctx, cj)

			cat1create, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{Admin: false, Colour: "#fe4efd", Description: "category testing", Name: "Category " + uuid.NewString()}, adminSession)
			tests.Ok(t, err, cat1create)

			thread1create, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Body:       "<p>this is a thread</p>",
				Category:   cat1create.JSON200.Id,
				Visibility: openapi.Published,
				Title:      "Thread testing",
			}, adminSession)
			tests.Ok(t, err, thread1create)

			like1, err := cl.LikePostAddWithResponse(root, thread1create.JSON200.Id, user1Session)
			tests.Ok(t, err, like1)

			like2, err := cl.LikePostAddWithResponse(root, thread1create.JSON200.Id, user2Session)
			tests.Ok(t, err, like2)

			// Assert both likes are present
			likepostget, err := cl.LikePostGetWithResponse(root, thread1create.JSON200.Id, adminSession)
			tests.Ok(t, err, likepostget)
			r.Len(likepostget.JSON200.Likes, 2)
			a.Equal(user1Acc.ID.String(), likepostget.JSON200.Likes[0].Owner.Id)
			a.Equal(user1Acc.Handle, likepostget.JSON200.Likes[0].Owner.Handle)
			a.Equal(user2Acc.ID.String(), likepostget.JSON200.Likes[1].Owner.Id)
			a.Equal(user2Acc.Handle, likepostget.JSON200.Likes[1].Owner.Handle)

			// Assert LikeStatus is correct
			tget1, err := cl.ThreadGetWithResponse(root, thread1create.JSON200.Id, user1Session)
			tests.Ok(t, err, tget1)
			a.Equal(2, tget1.JSON200.Likes.Likes)
			a.True(tget1.JSON200.Likes.Liked)

			tget2, err := cl.ThreadGetWithResponse(root, thread1create.JSON200.Id, user2Session)
			tests.Ok(t, err, tget2)
			a.Equal(2, tget2.JSON200.Likes.Likes)
			a.True(tget2.JSON200.Likes.Liked)

			tget3, err := cl.ThreadGetWithResponse(root, thread1create.JSON200.Id, adminSession)
			tests.Ok(t, err, tget3)
			a.Equal(2, tget3.JSON200.Likes.Likes)
			a.False(tget3.JSON200.Likes.Liked)

			// Assert profile likes contains thread
			profilelikeget, err := cl.LikeProfileGetWithResponse(root, user1Acc.Handle, &openapi.LikeProfileGetParams{}, user2Session)
			tests.Ok(t, err, profilelikeget)
			r.Equal(profilelikeget.JSON200.Results, 1)
			r.Equal(profilelikeget.JSON200.Results, len(profilelikeget.JSON200.Likes))
			itemPost, err := profilelikeget.JSON200.Likes[0].Item.AsDatagraphItemPost()
			r.NoError(err)
			a.Equal(thread1create.JSON200.Id, itemPost.Ref.Id)
			a.Equal(thread1create.JSON200.Slug, itemPost.Ref.Slug)

			// Assert thread list includes all like statuses
			threadlist1, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{Categories: &[]string{cat1create.JSON200.Slug}}, user1Session)
			tests.Ok(t, err, threadlist1)
			r.Len(threadlist1.JSON200.Threads, 1)
			a.Equal(2, threadlist1.JSON200.Threads[0].Likes.Likes)
			a.True(threadlist1.JSON200.Threads[0].Likes.Liked)

			threadlist2, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{Categories: &[]string{cat1create.JSON200.Slug}}, user2Session)
			tests.Ok(t, err, threadlist2)
			r.Len(threadlist2.JSON200.Threads, 1)
			a.Equal(2, threadlist2.JSON200.Threads[0].Likes.Likes)
			a.True(threadlist2.JSON200.Threads[0].Likes.Liked)

			threadlist3, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{Categories: &[]string{cat1create.JSON200.Slug}}, adminSession)
			tests.Ok(t, err, threadlist3)
			r.Len(threadlist3.JSON200.Threads, 1)
			a.Equal(2, threadlist3.JSON200.Threads[0].Likes.Likes)
			a.False(threadlist3.JSON200.Threads[0].Likes.Liked, "admin has not liked the thread")
		}))
	}))
}

func TestLikeReplies(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		cj *session_cookie.Jar,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			user1Ctx, user1Acc := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			user2Ctx, user2Acc := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			adminSession := e2e.WithSession(adminCtx, cj)
			user1Session := e2e.WithSession(user1Ctx, cj)
			user2Session := e2e.WithSession(user2Ctx, cj)

			cat1create, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{Admin: false, Colour: "#fe4efd", Description: "category testing", Name: "Category " + uuid.NewString()}, adminSession)
			tests.Ok(t, err, cat1create)

			thread1create, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Body:       "<p>this is a thread</p>",
				Category:   cat1create.JSON200.Id,
				Visibility: openapi.Published,
				Title:      "Thread testing",
			}, adminSession)
			tests.Ok(t, err, thread1create)

			reply1, err := cl.ReplyCreateWithResponse(root, thread1create.JSON200.Slug, openapi.ReplyInitialProps{Body: "<p>this is a reply</p>"}, user1Session)
			tests.Ok(t, err, reply1)

			like1, err := cl.LikePostAddWithResponse(root, reply1.JSON200.Id, user1Session)
			tests.Ok(t, err, like1)
			like2, err := cl.LikePostAddWithResponse(root, reply1.JSON200.Id, user2Session)
			tests.Ok(t, err, like2)

			// Assert no likes on thread
			likepostget, err := cl.LikePostGetWithResponse(root, thread1create.JSON200.Id, adminSession)
			tests.Ok(t, err, likepostget)
			r.Len(likepostget.JSON200.Likes, 0)

			// Assert like on reply
			likereplyget, err := cl.LikePostGetWithResponse(root, reply1.JSON200.Id, adminSession)
			tests.Ok(t, err, likereplyget)
			r.Len(likereplyget.JSON200.Likes, 2)
			a.Equal(user1Acc.ID.String(), likereplyget.JSON200.Likes[0].Owner.Id)
			a.Equal(user2Acc.ID.String(), likereplyget.JSON200.Likes[1].Owner.Id)

			// nobody liked the thread
			// user1 did like the reply
			// user2 did like the reply
			// admin user did not like the reply

			// get the thread as user1
			tget1, err := cl.ThreadGetWithResponse(root, thread1create.JSON200.Id, user1Session)
			tests.Ok(t, err, tget1)
			a.Equal(0, tget1.JSON200.Likes.Likes, "the thread has 0 likes")
			a.False(tget1.JSON200.Likes.Liked, "user1 did not like the thread")
			a.Equal(2, tget1.JSON200.Replies[0].Likes.Likes, "the reply has 2 likes from user1+2")
			a.True(tget1.JSON200.Replies[0].Likes.Liked, "user1 liked the reply")

			// get the thread as user2
			tget2, err := cl.ThreadGetWithResponse(root, thread1create.JSON200.Id, user2Session)
			tests.Ok(t, err, tget2)
			a.Equal(0, tget2.JSON200.Likes.Likes, "the thread has 0 likes")
			a.False(tget2.JSON200.Likes.Liked, "user2 did not like the thread")
			a.Equal(2, tget2.JSON200.Replies[0].Likes.Likes, "the reply has 1 likes from user1+2")
			a.True(tget2.JSON200.Replies[0].Likes.Liked, "user2 liked the reply")

			// get the thread as admin
			tget3, err := cl.ThreadGetWithResponse(root, thread1create.JSON200.Id, adminSession)
			tests.Ok(t, err, tget3)
			a.Equal(0, tget3.JSON200.Likes.Likes, "the thread has 0 likes")
			a.False(tget3.JSON200.Likes.Liked, "admin user did not like the thread")
			a.Equal(2, tget3.JSON200.Replies[0].Likes.Likes, "the reply has 2 likes from user1+2")
			a.False(tget3.JSON200.Replies[0].Likes.Liked, "admin user did not like the reply")

			// Assert profile likes contains reply
			profilelikeget, err := cl.LikeProfileGetWithResponse(root, user1Acc.Handle, &openapi.LikeProfileGetParams{}, user2Session)
			tests.Ok(t, err, profilelikeget)
			r.Equal(profilelikeget.JSON200.Results, 1)
			r.Equal(profilelikeget.JSON200.Results, len(profilelikeget.JSON200.Likes))
			itemPost, err := profilelikeget.JSON200.Likes[0].Item.AsDatagraphItemPost()
			r.NoError(err)
			a.Equal(reply1.JSON200.Id, itemPost.Ref.Id)
			a.Equal(thread1create.JSON200.Slug, itemPost.Ref.Slug)
		}))
	}))
}
