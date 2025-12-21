package thread_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/post/thread_cache"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/infrastructure/cache"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestThreadCacheWithReactions(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		cacheStore cache.Store,
		threadCache *thread_cache.Cache,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			acc1ctx, acc1 := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			acc2ctx, _ := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			session1 := sh.WithSession(acc1ctx)
			session2 := sh.WithSession(acc2ctx)

			catName := "Category " + uuid.NewString()

			catCreate, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Colour:      "#fe4efd",
				Description: "category testing",
				Name:        catName,
			}, session1)
			tests.Ok(t, err, catCreate)

			// Create a thread
			threadCreate, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Body:       opt.New("<p>test thread for cache</p>").Ptr(),
				Category:   opt.New(catCreate.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
				Title:      "Thread cache test",
			}, session1)
			tests.Ok(t, err, threadCreate)
			a.Equal(acc1.ID.String(), threadCreate.JSON200.Author.Id)
			a.Len(threadCreate.JSON200.Reacts, 0, "newly created thread should have no reactions")

			threadID := xid.ID(openapi.ParseID(threadCreate.JSON200.Id))

			// Get the thread for the first time and capture Last-Modified header
			threadGet1, err := cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, nil)
			tests.Ok(t, err, threadGet1)
			a.Len(threadGet1.JSON200.Reacts, 0, "thread should have no reactions")

			// Get the Last-Modified header from the response
			lastModified1Header := threadGet1.HTTPResponse.Header.Get("Last-Modified")
			r.NotEmpty(lastModified1Header, "Last-Modified header should be present")

			// Parse the Last-Modified header
			lastModified1, err := time.Parse(time.RFC1123, lastModified1Header)
			r.NoError(err, "Last-Modified header should be parseable")

			// Get the cached value directly from the cache store
			cacheKey := "thread:last-modified:" + threadID.String()
			cachedValue1, err := cacheStore.Get(root, cacheKey)
			r.NoError(err, "cache value should exist")
			cachedTime1, err := time.Parse(time.RFC3339Nano, cachedValue1)
			r.NoError(err, "cached time should be parseable")

			// Wait a small amount to ensure time difference
			time.Sleep(10 * time.Millisecond)

			// Make a conditional request with If-Modified-Since - should return 304 Not Modified
			threadGet304, err := cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, &openapi.ThreadGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-Modified-Since", lastModified1Header)
				return nil
			})
			tests.Status(t, err, threadGet304, 304)
			a.Nil(threadGet304.JSON200, "304 response should have no body")

			// Add a reaction to the thread
			reactAdd, err := cl.PostReactAddWithResponse(root, threadCreate.JSON200.Id, openapi.PostReactAddJSONRequestBody{
				Emoji: "👍",
			}, session2)
			tests.Ok(t, err, reactAdd)
			a.Equal("👍", reactAdd.JSON200.Emoji)

			// SYNCHRONOUS: Cache is invalidated immediately, no wait needed
			// Get the thread again and assert the reaction is present
			threadGet2, err := cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, nil)
			tests.Ok(t, err, threadGet2)
			r.Len(threadGet2.JSON200.Reacts, 1, "thread should now have one reaction")
			a.Equal("👍", threadGet2.JSON200.Reacts[0].Emoji)

			// Get the Last-Modified header from the second response
			lastModified2Header := threadGet2.HTTPResponse.Header.Get("Last-Modified")
			r.NotEmpty(lastModified2Header, "Last-Modified header should be present")

			// Parse the second Last-Modified header
			lastModified2, err := time.Parse(time.RFC1123, lastModified2Header)
			r.NoError(err, "Last-Modified header should be parseable")

			// Verify the cache was updated by checking the stored value FIRST
			// The cache stores timestamps with nanosecond precision, so we can verify
			// that the timestamp was actually updated even if HTTP headers use second precision
			cachedValue2, err := cacheStore.Get(root, cacheKey)
			r.NoError(err, "cache value should still exist after reaction")
			cachedTime2, err := time.Parse(time.RFC3339Nano, cachedValue2)
			r.NoError(err, "cached time should be parseable")

			// Assert that the cached time was updated (with nanosecond precision)
			a.True(cachedTime2.After(cachedTime1), "cached time should be updated after reaction (cachedTime1: %v, cachedTime2: %v)", cachedTime1, cachedTime2)

			// The Last-Modified header should be at least equal or after (limited by RFC1123 second precision)
			a.True(lastModified2.After(lastModified1) || lastModified2.Equal(lastModified1),
				"Last-Modified header should be updated or equal after reaction")

			// Additional verification: the LastModified from threadCache should return the updated value
			lastModifiedFromCache := threadCache.LastModified(root, threadID)
			r.NotNil(lastModifiedFromCache, "thread cache should return last modified time")
			// This should be equal to or after cachedTime2 since we just fetched it
			a.True(lastModifiedFromCache.Equal(cachedTime2) || lastModifiedFromCache.After(cachedTime2),
				"cache last modified should match the cached value")

			// Now test with a CONDITIONAL REQUEST using old timestamp - should get 200 not 304
			threadGetAfterReact, err := cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, &openapi.ThreadGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-Modified-Since", lastModified1Header)
				return nil
			})
			// Debug: Check what status we got
			if threadGetAfterReact.HTTPResponse.StatusCode == 304 {
				t.Logf("Got 304 - cache timestamps: before=%v, after=%v, last-modified header=%s",
					cachedTime1, cachedTime2, lastModified1Header)
				t.Logf("New Last-Modified header: %s", threadGetAfterReact.HTTPResponse.Header.Get("Last-Modified"))
			}
			tests.Ok(t, err, threadGetAfterReact)
			r.NotNil(threadGetAfterReact.JSON200, "should return 200 with body after cache invalidation")
			r.Len(threadGetAfterReact.JSON200.Reacts, 1, "should have the reaction in response")
		}))
	}))
}

func TestThreadCacheWithReplies(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		cacheStore cache.Store,
		threadCache *thread_cache.Cache,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			acc1ctx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			acc2ctx, _ := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			session1 := sh.WithSession(acc1ctx)
			session2 := sh.WithSession(acc2ctx)

			catName := "Category " + uuid.NewString()

			catCreate, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Colour:      "#abcdef",
				Description: "reply cache test",
				Name:        catName,
			}, session1)
			tests.Ok(t, err, catCreate)

			// Create a thread
			threadCreate, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Body:       opt.New("<p>test thread for replies</p>").Ptr(),
				Category:   opt.New(catCreate.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
				Title:      "Thread cache test - replies",
			}, session1)
			tests.Ok(t, err, threadCreate)
			a.Len(threadCreate.JSON200.Replies.Replies, 0, "newly created thread should have no replies")

			threadID := xid.ID(openapi.ParseID(threadCreate.JSON200.Id))

			// Get the thread for the first time
			threadGet1, err := cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, nil)
			tests.Ok(t, err, threadGet1)
			a.Len(threadGet1.JSON200.Replies.Replies, 0, "thread should have no replies")
			lastModified1Header := threadGet1.HTTPResponse.Header.Get("Last-Modified")
			r.NotEmpty(lastModified1Header, "Last-Modified header should be present")

			// Get cache timestamp before reply
			cacheKey := "thread:last-modified:" + threadID.String()
			cachedValue1, err := cacheStore.Get(root, cacheKey)
			r.NoError(err, "cache value should exist")
			cachedTime1, err := time.Parse(time.RFC3339Nano, cachedValue1)
			r.NoError(err, "cached time should be parseable")

			// Verify 304 works with conditional request
			threadGet304, err := cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, &openapi.ThreadGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-Modified-Since", lastModified1Header)
				return nil
			})
			tests.Status(t, err, threadGet304, 304)

			// Add a reply to the thread
			replyCreate, err := cl.ReplyCreateWithResponse(root, threadCreate.JSON200.Slug, openapi.ReplyInitialProps{
				Body: "<p>This is a test reply</p>",
			}, session2)
			tests.Ok(t, err, replyCreate)

			// SYNCHRONOUS: Cache should be invalidated IMMEDIATELY, no wait needed
			cachedValue2, err := cacheStore.Get(root, cacheKey)
			r.NoError(err, "cache value should still exist after reply")
			cachedTime2, err := time.Parse(time.RFC3339Nano, cachedValue2)
			r.NoError(err, "cached time should be parseable")

			// Assert cache was updated immediately
			a.True(cachedTime2.After(cachedTime1),
				"cache should be updated IMMEDIATELY after reply (no async delay needed)")

			// Conditional GET with old timestamp should now return 200 (not 304)
			threadGet200, err := cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, &openapi.ThreadGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-Modified-Since", lastModified1Header)
				return nil
			})
			tests.Ok(t, err, threadGet200)
			r.NotNil(threadGet200.JSON200, "should return 200 with body after cache invalidation")
			r.Len(threadGet200.JSON200.Replies.Replies, 1, "thread should have the reply")
			a.Equal("<body><p>This is a test reply</p></body>", threadGet200.JSON200.Replies.Replies[0].Body)
		}))
	}))
}

func TestThreadCacheWithReplyUpdate(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		cacheStore cache.Store,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			acc1ctx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			acc2ctx, _ := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			session1 := sh.WithSession(acc1ctx)
			session2 := sh.WithSession(acc2ctx)

			catName := "Category " + uuid.NewString()

			catCreate, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Colour:      "#123456",
				Description: "reply update cache test",
				Name:        catName,
			}, session1)
			tests.Ok(t, err, catCreate)

			threadCreate, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Body:       opt.New("<p>test thread for reply update</p>").Ptr(),
				Category:   opt.New(catCreate.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
				Title:      "Thread cache test - reply update",
			}, session1)
			tests.Ok(t, err, threadCreate)

			threadID := xid.ID(openapi.ParseID(threadCreate.JSON200.Id))
			cacheKey := "thread:last-modified:" + threadID.String()

			replyCreate, err := cl.ReplyCreateWithResponse(root, threadCreate.JSON200.Slug, openapi.ReplyInitialProps{
				Body: "<p>Original reply content</p>",
			}, session2)
			tests.Ok(t, err, replyCreate)
			replyID := replyCreate.JSON200.Id

			threadGet1, err := cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, nil)
			tests.Ok(t, err, threadGet1)
			r.Len(threadGet1.JSON200.Replies.Replies, 1)
			a.Equal("<body><p>Original reply content</p></body>", threadGet1.JSON200.Replies.Replies[0].Body)
			lastModified1Header := threadGet1.HTTPResponse.Header.Get("Last-Modified")
			r.NotEmpty(lastModified1Header, "Last-Modified header should be present")

			cachedValue1, err := cacheStore.Get(root, cacheKey)
			r.NoError(err, "cache value should exist")
			cachedTime1, err := time.Parse(time.RFC3339Nano, cachedValue1)
			r.NoError(err, "cached time should be parseable")

			threadGet304, err := cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, &openapi.ThreadGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-Modified-Since", lastModified1Header)
				return nil
			})
			tests.Status(t, err, threadGet304, 304)

			updatedBody := "<p>Updated reply content</p>"
			postUpdate, err := cl.PostUpdateWithResponse(root, replyID, openapi.PostUpdateJSONRequestBody{
				Body: &updatedBody,
			}, session2)
			tests.Ok(t, err, postUpdate)

			cachedValue2, err := cacheStore.Get(root, cacheKey)
			r.NoError(err, "cache value should still exist after reply update")
			cachedTime2, err := time.Parse(time.RFC3339Nano, cachedValue2)
			r.NoError(err, "cached time should be parseable")

			a.True(cachedTime2.After(cachedTime1),
				"cache should be updated IMMEDIATELY after reply update (cachedTime1: %v, cachedTime2: %v)", cachedTime1, cachedTime2)

			threadGet200, err := cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, &openapi.ThreadGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-Modified-Since", lastModified1Header)
				return nil
			})
			tests.Ok(t, err, threadGet200)
			r.NotNil(threadGet200.JSON200, "should return 200 with body after cache invalidation from reply update")
			r.Len(threadGet200.JSON200.Replies.Replies, 1, "thread should have the reply")
			a.Equal("<body><p>Updated reply content</p></body>", threadGet200.JSON200.Replies.Replies[0].Body)
		}))
	}))
}
