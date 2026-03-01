package thread_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/google/uuid"
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

func TestThreadCacheWithReactions(t *testing.T) {
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

			acc1ctx, acc1 := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			acc2ctx, _ := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			session1 := sh.WithSession(acc1ctx)
			session2 := sh.WithSession(acc2ctx)

			catName := "Category " + uuid.NewString()

			catCreate := tests.AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Colour:      "#fe4efd",
				Description: "category testing",
				Name:        catName,
			}, session1))(t, http.StatusOK)

			threadCreate := tests.AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Body:       opt.New("<p>test thread for cache</p>").Ptr(),
				Category:   opt.New(catCreate.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
				Title:      "Thread cache test",
			}, session1))(t, http.StatusOK)
			a.Equal(acc1.ID.String(), threadCreate.JSON200.Author.Id)
			a.Len(threadCreate.JSON200.Reacts, 0, "newly created thread should have no reactions")

			threadGet1 := tests.AssertRequest(cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, nil))(t, http.StatusOK)
			a.Len(threadGet1.JSON200.Reacts, 0, "thread should have no reactions")

			etag1 := threadGet1.HTTPResponse.Header.Get("ETag")
			r.NotEmpty(etag1, "ETag header should be present")

			threadGet304 := tests.AssertRequest(cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, &openapi.ThreadGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-None-Match", etag1)
				return nil
			}))(t, http.StatusNotModified)
			a.Nil(threadGet304.JSON200, "304 response should have no body")

			reactAdd := tests.AssertRequest(cl.PostReactAddWithResponse(root, threadCreate.JSON200.Id, openapi.PostReactAddJSONRequestBody{
				Emoji: "👍",
			}, session2))(t, http.StatusOK)
			a.Equal("👍", reactAdd.JSON200.Emoji)

			threadGet2 := tests.AssertRequest(cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, nil))(t, http.StatusOK)
			r.Len(threadGet2.JSON200.Reacts, 1, "thread should now have one reaction")
			a.Equal("👍", threadGet2.JSON200.Reacts[0].Emoji)

			etag2 := threadGet2.HTTPResponse.Header.Get("ETag")
			r.NotEmpty(etag2, "ETag header should be present")

			threadGetAfterReact := tests.AssertRequest(cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, &openapi.ThreadGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-None-Match", etag1)
				return nil
			}))(t, http.StatusOK)
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
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			acc1ctx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			acc2ctx, _ := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			session1 := sh.WithSession(acc1ctx)
			session2 := sh.WithSession(acc2ctx)

			catName := "Category " + uuid.NewString()

			catCreate := tests.AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Colour:      "#abcdef",
				Description: "reply cache test",
				Name:        catName,
			}, session1))(t, http.StatusOK)

			threadCreate := tests.AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Body:       opt.New("<p>test thread for replies</p>").Ptr(),
				Category:   opt.New(catCreate.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
				Title:      "Thread cache test - replies",
			}, session1))(t, http.StatusOK)
			a.Len(threadCreate.JSON200.Replies.Replies, 0, "newly created thread should have no replies")

			threadGet1 := tests.AssertRequest(cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, nil))(t, http.StatusOK)
			a.Len(threadGet1.JSON200.Replies.Replies, 0, "thread should have no replies")

			etag1 := threadGet1.HTTPResponse.Header.Get("ETag")
			r.NotEmpty(etag1, "ETag header should be present")
			lastModified1Header := threadGet1.HTTPResponse.Header.Get("Last-Modified")
			r.NotEmpty(lastModified1Header, "Last-Modified header should be present for backward compatibility")

			threadGet304 := tests.AssertRequest(cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, &openapi.ThreadGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-None-Match", etag1)
				return nil
			}))(t, http.StatusNotModified)
			a.Nil(threadGet304.JSON200, "304 response should have no body")

			tests.AssertRequest(cl.ReplyCreateWithResponse(root, threadCreate.JSON200.Slug, openapi.ReplyInitialProps{
				Body: "<p>This is a test reply</p>",
			}, session2))(t, http.StatusOK)

			threadGet200 := tests.AssertRequest(cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, &openapi.ThreadGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-None-Match", etag1)
				return nil
			}))(t, http.StatusOK)
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
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			acc1ctx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			acc2ctx, _ := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			session1 := sh.WithSession(acc1ctx)
			session2 := sh.WithSession(acc2ctx)

			catName := "Category " + uuid.NewString()

			catCreate := tests.AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Colour:      "#123456",
				Description: "reply update cache test",
				Name:        catName,
			}, session1))(t, http.StatusOK)

			threadCreate := tests.AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Body:       opt.New("<p>test thread for reply update</p>").Ptr(),
				Category:   opt.New(catCreate.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
				Title:      "Thread cache test - reply update",
			}, session1))(t, http.StatusOK)

			replyCreate := tests.AssertRequest(cl.ReplyCreateWithResponse(root, threadCreate.JSON200.Slug, openapi.ReplyInitialProps{
				Body: "<p>Original reply content</p>",
			}, session2))(t, http.StatusOK)
			replyID := replyCreate.JSON200.Id

			threadGet1 := tests.AssertRequest(cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, nil))(t, http.StatusOK)
			r.Len(threadGet1.JSON200.Replies.Replies, 1)
			a.Equal("<body><p>Original reply content</p></body>", threadGet1.JSON200.Replies.Replies[0].Body)

			etag1 := threadGet1.HTTPResponse.Header.Get("ETag")
			r.NotEmpty(etag1, "ETag header should be present")
			lastModified1Header := threadGet1.HTTPResponse.Header.Get("Last-Modified")
			r.NotEmpty(lastModified1Header, "Last-Modified header should be present for backward compatibility")

			threadGet304 := tests.AssertRequest(cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, &openapi.ThreadGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-None-Match", etag1)
				return nil
			}))(t, http.StatusNotModified)
			r.NotNil(threadGet304)

			updatedBody := "<p>Updated reply content</p>"
			tests.AssertRequest(cl.PostUpdateWithResponse(root, replyID, openapi.PostUpdateJSONRequestBody{
				Body: &updatedBody,
			}, session2))(t, http.StatusOK)

			threadGet200 := tests.AssertRequest(cl.ThreadGetWithResponse(root, threadCreate.JSON200.Slug, &openapi.ThreadGetParams{}, func(ctx context.Context, req *http.Request) error {
				req.Header.Set("If-None-Match", etag1)
				return nil
			}))(t, http.StatusOK)
			r.NotNil(threadGet200.JSON200, "should return 200 with body after cache invalidation from reply update")
			r.Len(threadGet200.JSON200.Replies.Replies, 1, "thread should have the reply")
			a.Equal("<body><p>Updated reply content</p></body>", threadGet200.JSON200.Replies.Replies[0].Body)
		}))
	}))
}
