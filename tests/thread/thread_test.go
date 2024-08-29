package thread_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
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
			r.NoError(err)
			r.NotNil(thread1get)
			r.Equal(http.StatusOK, thread1get.StatusCode())

			a.Len(thread1get.JSON200.Replies, 1)
			replyids := dt.Map(thread1get.JSON200.Replies, func(p openapi.Reply) string { return p.Id })
			a.Contains(replyids, reply1create.JSON200.Id)
		}))
	}))
}

func TestThreadLinkAggregation(t *testing.T) {
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

			ctx, _ := e2e.WithAccount(ctx, aw, seed.Account_002_Frigg)

			catname := "Category " + uuid.NewString()
			cat, err := cl.CategoryCreateWithResponse(ctx, openapi.CategoryInitialProps{
				Admin:       false,
				Colour:      "#fe4efd",
				Description: "category testing",
				Name:        catname,
			}, e2e.WithSession(ctx, cj))
			tests.Ok(t, err, cat)

			url := "https://ogp.me"
			thread1create, err := cl.ThreadCreateWithResponse(ctx, openapi.ThreadInitialProps{
				Body:       "<p>this is a thread</p>",
				Category:   cat.JSON200.Id,
				Visibility: openapi.Published,
				Title:      "Thread URL link aggregation",
				Url:        &url,
			}, e2e.WithSession(ctx, cj))
			tests.Ok(t, err, thread1create)

			r.NotNil(thread1create.JSON200.Link)
			a.Equal(url, thread1create.JSON200.Link.Url)
			a.Equal("Open Graph protocol", *thread1create.JSON200.Link.Title)
			a.Equal("The Open Graph protocol enables any web page to become a rich object in a social graph.", *thread1create.JSON200.Link.Description)
			a.Equal("ogp.me", thread1create.JSON200.Link.Domain)
			a.Equal("ogp-me", thread1create.JSON200.Link.Slug)
		}))
	}))
}
