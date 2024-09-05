package thread_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestThreads(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		cj *session.Jar,
		aw account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			ctx, acc := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)
			ctx2, acc2 := e2e.WithAccount(ctx, aw, seed.Account_003_Baldur)

			cat1name := "Category " + uuid.NewString()

			cat1create, err := cl.CategoryCreateWithResponse(ctx, openapi.CategoryInitialProps{
				Admin:       false,
				Colour:      "#fe4efd",
				Description: "category testing",
				Name:        cat1name,
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(cat1create)
			r.Equal(http.StatusOK, cat1create.StatusCode())

			a.Equal(cat1name, cat1create.JSON200.Name)
			a.Equal("category testing", cat1create.JSON200.Description)
			a.Equal(slug.Make(cat1name), cat1create.JSON200.Slug)

			thread1create, err := cl.ThreadCreateWithResponse(ctx, openapi.ThreadInitialProps{
				Body:       "<p>this is a thread</p>",
				Category:   cat1create.JSON200.Id,
				Visibility: openapi.Published,
				Title:      "Thread testing",
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(thread1create)
			r.Equal(http.StatusOK, thread1create.StatusCode())

			a.Nil(thread1create.JSON200.DeletedAt)
			a.Equal(acc.ID.String(), thread1create.JSON200.Author.Id)
			a.Equal("Thread testing", thread1create.JSON200.Title)
			a.Contains(thread1create.JSON200.Slug, "thread-testing")
			a.Equal("<body><p>this is a thread</p></body>", thread1create.JSON200.Body)
			a.Equal("this is a thread", *thread1create.JSON200.Description)
			a.Equal(false, thread1create.JSON200.Pinned)
			a.Equal(cat1create.JSON200.Name, thread1create.JSON200.Category.Name)
			a.Len(thread1create.JSON200.Replies, 0, "a newly created thread has zero replies")
			a.Len(thread1create.JSON200.Reacts, 0)
			a.Equal(thread1create.JSON200.ReplyStatus.Replies, 0)
			a.Equal(thread1create.JSON200.ReplyStatus.Replied, 0)

			// Get list of all threads

			threadlist, err := cl.ThreadListWithResponse(ctx, &openapi.ThreadListParams{})
			r.NoError(err)
			r.NotNil(threadlist)
			r.Equal(http.StatusOK, threadlist.StatusCode())

			ids := dt.Map(threadlist.JSON200.Threads, func(th openapi.ThreadReference) string { return th.Id })
			a.Contains(ids, thread1create.JSON200.Id)

			// Reply to the thread

			reply1create, err := cl.ReplyCreateWithResponse(ctx, thread1create.JSON200.Slug, openapi.ReplyInitialProps{
				Body: "this is a reply",
			}, e2e.WithSession(ctx2, cj))
			r.NoError(err)
			r.NotNil(reply1create)
			r.Equal(http.StatusOK, reply1create.StatusCode())

			a.Equal(acc2.ID.String(), reply1create.JSON200.Author.Id)

			thread1get, err := cl.ThreadGetWithResponse(ctx, thread1create.JSON200.Slug)
			tests.Ok(t, err, thread1get)

			a.Len(thread1get.JSON200.Replies, 1)
			replyids := dt.Map(thread1get.JSON200.Replies, func(p openapi.Reply) string { return p.Id })
			a.Contains(replyids, reply1create.JSON200.Id)
			a.Equal(thread1get.JSON200.ReplyStatus.Replies, 1)
			a.Equal(thread1get.JSON200.ReplyStatus.Replied, 0)

			thread2get, err := cl.ThreadGetWithResponse(ctx, thread1create.JSON200.Slug, e2e.WithSession(ctx2, cj))
			tests.Ok(t, err, thread2get)
			a.Equal(thread2get.JSON200.ReplyStatus.Replies, 1)
			a.Equal(thread2get.JSON200.ReplyStatus.Replied, 1, "ctx2 replied")
		}))
	}))
}

func TestThreadLinkAggregation(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		cj *session.Jar,
		aw account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			sessionCtx, _ := e2e.WithAccount(root, aw, seed.Account_002_Frigg)
			session := e2e.WithSession(sessionCtx, cj)

			catname := "Category " + uuid.NewString()
			cat, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Admin:       false,
				Colour:      "#fe4efd",
				Description: "category testing",
				Name:        catname,
			}, session)
			tests.Ok(t, err, cat)

			url := "https://ogp.me"
			thread1create, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Body:       "<p>this is a thread</p>",
				Category:   cat.JSON200.Id,
				Visibility: openapi.Published,
				Title:      "Thread URL link aggregation",
				Url:        &url,
			}, session)
			tests.Ok(t, err, thread1create)

			r.NotNil(thread1create.JSON200.Link)
			a.Equal(url, thread1create.JSON200.Link.Url)
			a.Equal("ogp-me", thread1create.JSON200.Link.Slug)
			a.Equal("ogp.me", thread1create.JSON200.Link.Domain)
			a.Equal("Open Graph protocol", *thread1create.JSON200.Link.Title)
			a.Equal("The Open Graph protocol enables any web page to become a rich object in a social graph.", *thread1create.JSON200.Link.Description)
			r.NotNil(thread1create.JSON200.Link.FaviconImage)
			a.NotEmpty(thread1create.JSON200.Link.FaviconImage.Id)
			r.NotNil(thread1create.JSON200.Link.PrimaryImage)
			a.NotEmpty(thread1create.JSON200.Link.PrimaryImage.Id)

			// Get the thread just created, ensure link is present.

			thread1get, err := cl.ThreadGetWithResponse(root, thread1create.JSON200.Slug, session)
			tests.Ok(t, err, thread1get)

			r.NotNil(thread1get.JSON200.Link)
			a.Equal(url, thread1get.JSON200.Link.Url)
			a.Equal("ogp-me", thread1get.JSON200.Link.Slug)
			a.Equal("ogp.me", thread1get.JSON200.Link.Domain)
			a.Equal("Open Graph protocol", *thread1get.JSON200.Link.Title)
			a.Equal("The Open Graph protocol enables any web page to become a rich object in a social graph.", *thread1get.JSON200.Link.Description)
			r.NotNil(thread1get.JSON200.Link.FaviconImage)
			a.NotEmpty(thread1get.JSON200.Link.FaviconImage.Id)
			r.NotNil(thread1get.JSON200.Link.PrimaryImage)
			a.NotEmpty(thread1get.JSON200.Link.PrimaryImage.Id)

			// List threads, ensure link is present.

			threadlist, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{})
			tests.Ok(t, err, threadlist)

			listThread, found := lo.Find(threadlist.JSON200.Threads, func(th openapi.ThreadReference) bool {
				return th.Id == thread1create.JSON200.Id
			})
			r.True(found)

			r.NotNil(listThread.Link)
			a.Equal(url, listThread.Link.Url)
			a.Equal("ogp-me", listThread.Link.Slug)
			a.Equal("ogp.me", listThread.Link.Domain)
			a.Equal("Open Graph protocol", *listThread.Link.Title)
			a.Equal("The Open Graph protocol enables any web page to become a rich object in a social graph.", *listThread.Link.Description)
			r.NotNil(listThread.Link.FaviconImage)
			a.NotEmpty(listThread.Link.FaviconImage.Id)
			r.NotNil(listThread.Link.PrimaryImage)
			a.NotEmpty(listThread.Link.PrimaryImage.Id)
		}))
	}))
}
