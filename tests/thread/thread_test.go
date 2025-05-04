package thread_test

import (
	"context"
	"testing"

	"github.com/Southclaws/dt"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
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

func TestThreads(t *testing.T) {
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

			acc1ctx, acc1 := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			acc2ctx, acc2 := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			session1 := sh.WithSession(acc1ctx)
			session2 := sh.WithSession(acc2ctx)

			cat1name := "Category " + uuid.NewString()

			cat1create, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Admin:       false,
				Colour:      "#fe4efd",
				Description: "category testing",
				Name:        cat1name,
			}, session1)
			tests.Ok(t, err, cat1create)
			r.Equal(cat1name, cat1create.JSON200.Name)
			r.Equal("category testing", cat1create.JSON200.Description)
			r.Equal(slug.Make(cat1name), cat1create.JSON200.Slug)

			t.Run("thread_replies", func(t *testing.T) {
				t.Parallel()

				r := require.New(t)
				a := assert.New(t)

				thread1create, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       "<p>this is a thread</p>",
					Category:   cat1create.JSON200.Id,
					Visibility: openapi.Published,
					Title:      "Thread testing",
				}, session1)
				tests.Ok(t, err, thread1create)
				a.Nil(thread1create.JSON200.DeletedAt)
				a.Equal(acc1.ID.String(), thread1create.JSON200.Author.Id)
				a.Equal("Thread testing", thread1create.JSON200.Title)
				a.Contains(thread1create.JSON200.Slug, "thread-testing")
				a.Equal("<body><p>this is a thread</p></body>", thread1create.JSON200.Body)
				a.Equal("this is a thread", *thread1create.JSON200.Description)
				a.Equal(false, thread1create.JSON200.Pinned)
				a.Equal(cat1create.JSON200.Name, thread1create.JSON200.Category.Name)
				a.Len(thread1create.JSON200.Replies.Replies, 0, "a newly created thread has zero replies")
				a.Len(thread1create.JSON200.Reacts, 0)
				a.Equal(thread1create.JSON200.ReplyStatus.Replies, 0)
				a.Equal(thread1create.JSON200.ReplyStatus.Replied, 0)

				// Get list of all threads

				threadlist, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{})
				tests.Ok(t, err, threadlist)
				ids := dt.Map(threadlist.JSON200.Threads, func(th openapi.ThreadReference) string { return th.Id })
				a.Contains(ids, thread1create.JSON200.Id)

				// Reply to the thread

				reply1create, err := cl.ReplyCreateWithResponse(root, thread1create.JSON200.Slug, openapi.ReplyInitialProps{
					Body: "this is a reply",
				}, session2)
				tests.Ok(t, err, reply1create)
				a.Equal(acc2.ID.String(), reply1create.JSON200.Author.Id)

				thread1get, err := cl.ThreadGetWithResponse(root, thread1create.JSON200.Slug, nil)
				tests.Ok(t, err, thread1get)

				r.Len(thread1get.JSON200.Replies.Replies, 1)
				replyids := dt.Map(thread1get.JSON200.Replies.Replies, func(p openapi.Reply) string { return p.Id })
				a.Contains(replyids, reply1create.JSON200.Id)
				a.Equal(thread1get.JSON200.ReplyStatus.Replies, 1)
				a.Equal(thread1get.JSON200.ReplyStatus.Replied, 0)

				thread2get, err := cl.ThreadGetWithResponse(root, thread1create.JSON200.Slug, nil, session2)
				tests.Ok(t, err, thread2get)
				a.Equal(thread2get.JSON200.ReplyStatus.Replies, 1)
				a.Equal(thread2get.JSON200.ReplyStatus.Replied, 1, "ctx2 replied")
			})

			t.Run("threads_ordered_by_last_reply", func(t *testing.T) {
				t.Parallel()

				r := require.New(t)
				a := assert.New(t)

				thread1create, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       "<p>thread one</p>",
					Category:   cat1create.JSON200.Id,
					Visibility: openapi.Published,
					Title:      "1",
				}, session1)
				tests.Ok(t, err, thread1create)

				thread2create, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       "<p>thread two</p>",
					Category:   cat1create.JSON200.Id,
					Visibility: openapi.Published,
					Title:      "2",
				}, session1)
				tests.Ok(t, err, thread2create)

				id1, id2 := thread1create.JSON200.Id, thread2create.JSON200.Id

				{
					threadlist, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{})
					tests.Ok(t, err, threadlist)
					threads := filterThreads(threadlist.JSON200.Threads, id1, id2)
					ids := getIDs(threads)
					r.Len(ids, 2)
					wantIDs := []openapi.Identifier{id2, id1}
					a.Equal(wantIDs, ids)
					gotThread1 := threads[1]
					gotThread2 := threads[0]
					r.Nil(gotThread1.LastReplyAt)
					r.Nil(gotThread2.LastReplyAt)
				}

				// Reply to thread 1, bumping it to the top
				reply1create, err := cl.ReplyCreateWithResponse(root, thread1create.JSON200.Slug, openapi.ReplyInitialProps{
					Body: "this is a reply",
				}, session2)
				tests.Ok(t, err, reply1create)

				{
					threadlist, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{})
					tests.Ok(t, err, threadlist)
					threads := filterThreads(threadlist.JSON200.Threads, id1, id2)
					ids2 := getIDs(threads)
					r.Len(ids2, 2)
					wantIDs2 := []openapi.Identifier{id1, id2}
					r.Equal(wantIDs2, ids2)
					gotThread1 := threads[0]
					gotThread2 := threads[1]
					r.NotNil(gotThread1.LastReplyAt)
					r.Nil(gotThread2.LastReplyAt)
				}

				// Reply to thread 2, bumping it to the top
				reply2create, err := cl.ReplyCreateWithResponse(root, thread2create.JSON200.Slug, openapi.ReplyInitialProps{
					Body: "this is a reply",
				}, session2)
				tests.Ok(t, err, reply2create)

				{
					threadlist, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{})
					tests.Ok(t, err, threadlist)
					threads := filterThreads(threadlist.JSON200.Threads, id1, id2)
					ids2 := getIDs(threads)
					r.Len(ids2, 2)
					wantIDs2 := []openapi.Identifier{id2, id1}
					a.Equal(wantIDs2, ids2)
					gotThread1 := threads[1]
					gotThread2 := threads[0]
					a.Less(*gotThread1.LastReplyAt, *gotThread2.LastReplyAt)
				}
			})

			t.Run("delete_replies", func(t *testing.T) {
				t.Parallel()

				r := require.New(t)
				a := assert.New(t)

				t1, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       "<p>t1</p>",
					Category:   cat1create.JSON200.Id,
					Visibility: openapi.Published,
					Title:      "t1",
				}, session1)
				tests.Ok(t, err, t1)

				r1, err := cl.ReplyCreateWithResponse(root, t1.JSON200.Slug, openapi.ReplyInitialProps{Body: "r1"}, session2)
				tests.Ok(t, err, r1)
				r2, err := cl.ReplyCreateWithResponse(root, t1.JSON200.Slug, openapi.ReplyInitialProps{Body: "r2"}, session2)
				tests.Ok(t, err, r2)
				r3, err := cl.ReplyCreateWithResponse(root, t1.JSON200.Slug, openapi.ReplyInitialProps{Body: "r3"}, session2)
				tests.Ok(t, err, r3)

				r1del, err := cl.PostDeleteWithResponse(root, r1.JSON200.Id, session2)
				tests.Ok(t, err, r1del)

				t1get, err := cl.ThreadGetWithResponse(root, t1.JSON200.Slug, nil)
				tests.Ok(t, err, t1get)
				a.Equal(2, t1get.JSON200.ReplyStatus.Replies)

				tlist, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{}, session2)
				tests.Ok(t, err, tlist)
				tlist1, found := lo.Find(tlist.JSON200.Threads, func(t openapi.ThreadReference) bool { return t.Id == t1.JSON200.Id })
				r.True(found)
				a.Equal(2, tlist1.ReplyStatus.Replies)
			})

			t.Run("link_aggregation", func(t *testing.T) {
				t.Parallel()

				r := require.New(t)
				a := assert.New(t)

				ctx1, _ := e2e.WithAccount(root, aw, seed.Account_002_Frigg)
				session := sh.WithSession(ctx1)

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

				thread1get, err := cl.ThreadGetWithResponse(root, thread1create.JSON200.Slug, nil, session)
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
			})
		}))
	}))
}

func getIDs(t []openapi.ThreadReference) []openapi.Identifier {
	return dt.Map(t, func(th openapi.ThreadReference) string { return th.Id })
}

func filterThreads(ts []openapi.ThreadReference, ids ...openapi.Identifier) []openapi.ThreadReference {
	filtered := dt.Filter(ts, func(t openapi.ThreadReference) bool {
		return lo.Contains(ids, t.Id)
	})

	return filtered
}
