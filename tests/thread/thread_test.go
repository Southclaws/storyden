package thread_test

import (
	"context"
	"testing"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
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
				Colour:      "#fe4efd",
				Description: "category testing",
				Name:        cat1name,
			}, session1)
			tests.Ok(t, err, cat1create)
			r.Equal(cat1name, cat1create.JSON200.Name)
			r.Equal("category testing", cat1create.JSON200.Description)
			r.Equal(mark.Slugify(cat1name), cat1create.JSON200.Slug)

			t.Run("thread_replies", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				thread1create, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       opt.New("<p>this is a thread</p>").Ptr(),
					Category:   opt.New(cat1create.JSON200.Id).Ptr(),
					Visibility: opt.New(openapi.Published).Ptr(),
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

			t.Run("reply_to_reply", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				thread1create, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       opt.New("<p>thread for nested replies</p>").Ptr(),
					Category:   opt.New(cat1create.JSON200.Id).Ptr(),
					Visibility: opt.New(openapi.Published).Ptr(),
					Title:      "Reply-to-reply testing",
				}, session1)
				tests.Ok(t, err, thread1create)

				// acc2 creates first reply to thread
				reply1create, err := cl.ReplyCreateWithResponse(root, thread1create.JSON200.Slug, openapi.ReplyInitialProps{
					Body: "first reply",
				}, session2)
				tests.Ok(t, err, reply1create)
				a.Equal(acc2.ID.String(), reply1create.JSON200.Author.Id)
				a.Nil(reply1create.JSON200.ReplyTo, "first reply should not have reply_to")

				// acc1 creates a reply to reply1
				reply2create, err := cl.ReplyCreateWithResponse(root, thread1create.JSON200.Slug, openapi.ReplyInitialProps{
					Body:    "nested reply",
					ReplyTo: &reply1create.JSON200.Id,
				}, session1)
				tests.Ok(t, err, reply2create)
				a.Equal(acc1.ID.String(), reply2create.JSON200.Author.Id)
				r.NotNil(reply2create.JSON200.ReplyTo, "nested reply should have reply_to set")
				a.Equal(reply1create.JSON200.Id, reply2create.JSON200.ReplyTo.Id, "nested reply should reference first reply")

				// Get thread and verify both replies are present with correct relationships
				thread1get, err := cl.ThreadGetWithResponse(root, thread1create.JSON200.Slug, nil)
				tests.Ok(t, err, thread1get)
				r.Len(thread1get.JSON200.Replies.Replies, 2)

				replies := thread1get.JSON200.Replies.Replies
				reply1get, found := lo.Find(replies, func(r openapi.Reply) bool { return r.Id == reply1create.JSON200.Id })
				r.True(found)
				a.Nil(reply1get.ReplyTo, "first reply should not have reply_to")

				reply2get, found := lo.Find(replies, func(r openapi.Reply) bool { return r.Id == reply2create.JSON200.Id })
				r.True(found)
				r.NotNil(reply2get.ReplyTo, "nested reply should have reply_to")
				a.Equal(reply1create.JSON200.Id, reply2get.ReplyTo.Id, "nested reply should reference first reply")
			})

			t.Run("threads_ordered_by_last_reply", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				thread1create, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       opt.New("<p>thread one</p>").Ptr(),
					Category:   opt.New(cat1create.JSON200.Id).Ptr(),
					Visibility: opt.New(openapi.Published).Ptr(),
					Title:      "1",
				}, session1)
				tests.Ok(t, err, thread1create)

				thread2create, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       opt.New("<p>thread two</p>").Ptr(),
					Category:   opt.New(cat1create.JSON200.Id).Ptr(),
					Visibility: opt.New(openapi.Published).Ptr(),
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

			t.Run("threads_ordered_by_last_reply_nulls_first", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				// create thread a and reply to it (so it gets last_reply_at set)
				threadA, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       opt.New("<p>thread A</p>").Ptr(),
					Category:   opt.New(cat1create.JSON200.Id).Ptr(),
					Visibility: opt.New(openapi.Published).Ptr(),
					Title:      "Thread A",
				}, session1)
				tests.Ok(t, err, threadA)

				replyA, err := cl.ReplyCreateWithResponse(root, threadA.JSON200.Slug, openapi.ReplyInitialProps{
					Body: "reply to thread A",
				}, session2)
				tests.Ok(t, err, replyA)

				// create thread b (no replies, so last_reply_at will be null)
				threadB, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       opt.New("<p>thread B</p>").Ptr(),
					Category:   opt.New(cat1create.JSON200.Id).Ptr(),
					Visibility: opt.New(openapi.Published).Ptr(),
					Title:      "Thread B",
				}, session1)
				tests.Ok(t, err, threadB)

				idA := threadA.JSON200.Id
				idB := threadB.JSON200.Id

				// check ordering: threadb (no last_reply_at) should appear first
				threadList, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{})
				tests.Ok(t, err, threadList)

				threads := filterThreads(threadList.JSON200.Threads, idA, idB)
				r.Len(threads, 2)

				gotIDs := getIDs(threads)
				wantIDs := []openapi.Identifier{idB, idA}
				a.Equal(wantIDs, gotIDs)

				gotThreadA := threads[1]
				gotThreadB := threads[0]

				// Assertions for clarity
				a.NotNil(gotThreadA.LastReplyAt, "Thread A should have a last_reply_at because it has a reply")
				a.Nil(gotThreadB.LastReplyAt, "Thread B should have no last_reply_at because it has no replies")
			})

			t.Run("delete_replies", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				t1, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       opt.New("<p>t1</p>").Ptr(),
					Category:   opt.New(cat1create.JSON200.Id).Ptr(),
					Visibility: opt.New(openapi.Published).Ptr(),
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
				r := require.New(t)
				a := assert.New(t)

				ctx1, _ := e2e.WithAccount(root, aw, seed.Account_002_Frigg)
				session := sh.WithSession(ctx1)

				catname := "Category " + uuid.NewString()
				cat, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
					Colour:      "#fe4efd",
					Description: "category testing",
					Name:        catname,
				}, session)
				tests.Ok(t, err, cat)

				url := "https://ogp.me"
				thread1create, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       opt.New("<p>this is a thread</p>").Ptr(),
					Category:   opt.New(cat.JSON200.Id).Ptr(),
					Visibility: opt.New(openapi.Published).Ptr(),
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

			t.Run("thread_without_category", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				// Create a thread without a category as a draft
				threadCreate, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       opt.New("<p>this is a thread without category</p>").Ptr(),
					Visibility: opt.New(openapi.Draft).Ptr(),
					Title:      "Thread without category",
				}, session1)
				tests.Ok(t, err, threadCreate)

				{
					// Verify the thread was created successfully
					a.Nil(threadCreate.JSON200.DeletedAt)
					a.Equal(acc1.ID.String(), threadCreate.JSON200.Author.Id)
					a.Equal("Thread without category", threadCreate.JSON200.Title)
					a.Contains(threadCreate.JSON200.Slug, "thread-without-category")
					a.Equal("<body><p>this is a thread without category</p></body>", threadCreate.JSON200.Body)
					a.Equal("this is a thread without category", *threadCreate.JSON200.Description)
					a.Equal(false, threadCreate.JSON200.Pinned)
					a.Nil(threadCreate.JSON200.Category, "thread should not have a category")
					r.Len(threadCreate.JSON200.Replies.Replies, 0, "a newly created thread has zero replies")
					a.Equal(threadCreate.JSON200.ReplyStatus.Replies, 0)
					a.Equal(threadCreate.JSON200.ReplyStatus.Replied, 0)
					// a.Equal(string(openapi.Draft), threadCreate.JSON200.Visibility)
				}

				// Update the thread to add a category and change visibility to published
				threadUpdate, err := cl.ThreadUpdateWithResponse(root, threadCreate.JSON200.Slug, openapi.ThreadMutableProps{
					Category:   &cat1create.JSON200.Id,
					Visibility: opt.New(openapi.Published).Ptr(),
				}, session1)
				tests.Ok(t, err, threadUpdate)

				{
					// Verify the thread was updated successfully
					a.Equal("Thread without category", threadUpdate.JSON200.Title)
					a.NotNil(threadUpdate.JSON200.Category, "thread should now have a category")
					a.Equal(cat1create.JSON200.Id, threadUpdate.JSON200.Category.Id)
					a.Equal(cat1create.JSON200.Name, threadUpdate.JSON200.Category.Name)
					a.Equal(openapi.Published, threadUpdate.JSON200.Visibility)
				}

				{
					// Get the thread to verify the changes persisted
					threadGet, err := cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, nil)
					tests.Ok(t, err, threadGet)
					a.Equal("Thread without category", threadGet.JSON200.Title)
					a.NotNil(threadGet.JSON200.Category, "thread should have a category after update")
					a.Equal(cat1create.JSON200.Id, threadGet.JSON200.Category.Id)
					a.Equal(cat1create.JSON200.Name, threadGet.JSON200.Category.Name)
					a.Equal(openapi.Published, threadGet.JSON200.Visibility)
				}

				{
					// Verify the thread appears in the thread list since it's now published
					threadList, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{})
					tests.Ok(t, err, threadList)
					ids := dt.Map(threadList.JSON200.Threads, func(th openapi.ThreadReference) string { return th.Id })
					a.Contains(ids, threadCreate.JSON200.Id, "published thread should appear in list")
				}

				{
					// Verify the thread does not appear in the uncategorised list
					threadList, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{
						Categories: &[]string{"null"},
					})
					tests.Ok(t, err, threadList)
					ids := dt.Map(threadList.JSON200.Threads, func(th openapi.ThreadReference) string { return th.Id })
					a.NotContains(ids, threadCreate.JSON200.Id, "categorised thread should not appear in query for uncategorised threads")
				}
			})

			t.Run("uncategorised_threads_in_thread_list", func(t *testing.T) {
				a := assert.New(t)

				cat1create, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
					Colour:      "#fe4efd",
					Description: "categorised threads",
					Name:        "categorised " + uuid.NewString(),
				}, session1)
				tests.Ok(t, err, cat1create)

				categorisedThread, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       opt.New("<p>this is a thread without category</p>").Ptr(),
					Visibility: opt.New(openapi.Published).Ptr(),
					Title:      "Thread without category",
					Category:   opt.New(cat1create.JSON200.Id).Ptr(),
				}, session1)
				tests.Ok(t, err, categorisedThread)
				uncategorisedThread, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       opt.New("<p>this is a thread without category</p>").Ptr(),
					Visibility: opt.New(openapi.Published).Ptr(),
					Title:      "Thread without category",
				}, session1)
				tests.Ok(t, err, uncategorisedThread)

				{
					threadList, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{
						Categories: &[]string{"null"},
					})
					tests.Ok(t, err, threadList)
					ids := dt.Map(threadList.JSON200.Threads, func(th openapi.ThreadReference) string { return th.Id })
					a.Contains(ids, uncategorisedThread.JSON200.Id, "uncategorised thread should appear in list")
					a.NotContains(ids, categorisedThread.JSON200.Id, "categorised thread should not appear in uncategorised list")
				}

				{
					threadList, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{})
					tests.Ok(t, err, threadList)
					ids := dt.Map(threadList.JSON200.Threads, func(th openapi.ThreadReference) string { return th.Id })
					a.Contains(ids, uncategorisedThread.JSON200.Id, "uncategorised thread should appear in list")
					a.Contains(ids, categorisedThread.JSON200.Id, "categorised thread should appear in list")
				}

				{
					threadList, err := cl.ThreadListWithResponse(root, &openapi.ThreadListParams{
						Categories: &[]string{cat1create.JSON200.Slug},
					})
					tests.Ok(t, err, threadList)
					ids := dt.Map(threadList.JSON200.Threads, func(th openapi.ThreadReference) string { return th.Id })
					a.NotContains(ids, uncategorisedThread.JSON200.Id, "uncategorised thread should appear in list")
					a.Contains(ids, categorisedThread.JSON200.Id, "categorised thread should appear in list")
				}
			})

			t.Run("thread_with_russian_slug", func(t *testing.T) {
				a := assert.New(t)

				threadCreate, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       opt.New("<p>test content</p>").Ptr(),
					Category:   opt.New(cat1create.JSON200.Id).Ptr(),
					Visibility: opt.New(openapi.Published).Ptr(),
					Title:      "Бабочки",
				}, session1)
				tests.Ok(t, err, threadCreate)
				a.Contains(threadCreate.JSON200.Slug, "бабочки")

				threadGet, err := cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, nil)
				tests.Ok(t, err, threadGet)
				a.Equal("Бабочки", threadGet.JSON200.Title)
				a.Contains(threadGet.JSON200.Slug, "бабочки")
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
