package search_test

import (
	"context"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestSimpleSearchThreadReplyFiltering(t *testing.T) {
	cfg := &config.Config{}

	integration.Test(t, cfg, e2e.Setup(), fx.Invoke(func(
		root context.Context,
		lc fx.Lifecycle,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			catResp, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Name:   "test-category-" + uuid.NewString(),
				Colour: "#123456",
			}, adminSession)
			tests.Ok(t, err, catResp)

			threadResp, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "Test Thread About Programming",
				Body:       opt.New("<p>A thread discussing programming topics</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadResp)

			reply1Resp, err := cl.ReplyCreateWithResponse(root, threadResp.JSON200.Slug, openapi.ReplyInitialProps{
				Body: "<p>This is a reply about programming languages</p>",
			}, adminSession)
			tests.Ok(t, err, reply1Resp)

			reply2Resp, err := cl.ReplyCreateWithResponse(root, threadResp.JSON200.Slug, openapi.ReplyInitialProps{
				Body: "<p>Another reply discussing programming paradigms</p>",
			}, adminSession)
			tests.Ok(t, err, reply2Resp)

			t.Run("search_only_threads", func(t *testing.T) {
				r := require.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "programming",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find at least one thread")

				for _, item := range resp.JSON200.Items {
					threadItem, err := item.AsDatagraphItemThread()
					r.NoError(err, "all items should be threads when filtering by thread kind")
					r.NotEmpty(threadItem.Ref.Id, "thread should have an ID")
				}

				foundThread := findThreadItem(resp.JSON200.Items, threadResp.JSON200.Id)
				r.NotNil(foundThread, "should find the created thread")
			})

			t.Run("search_only_replies", func(t *testing.T) {
				r := require.New(t)

				replyKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindReply}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "programming",
					Kind: &replyKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 2, "should find at least two replies")

				for _, item := range resp.JSON200.Items {
					replyItem, err := item.AsDatagraphItemReply()
					r.NoError(err, "all items should be replies when filtering by reply kind")
					r.NotEmpty(replyItem.Ref.Id, "reply should have an ID")
				}

				foundReply1 := findReplyItem(resp.JSON200.Items, reply1Resp.JSON200.Id)
				foundReply2 := findReplyItem(resp.JSON200.Items, reply2Resp.JSON200.Id)

				r.NotNil(foundReply1, "should find first reply")
				r.NotNil(foundReply2, "should find second reply")
			})

			t.Run("search_both_threads_and_replies", func(t *testing.T) {
				r := require.New(t)

				bothKinds := []openapi.DatagraphItemKind{
					openapi.DatagraphItemKindThread,
					openapi.DatagraphItemKindReply,
				}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "programming",
					Kind: &bothKinds,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 3, "should find thread and replies")

				foundThread := findThreadItem(resp.JSON200.Items, threadResp.JSON200.Id)
				foundReply1 := findReplyItem(resp.JSON200.Items, reply1Resp.JSON200.Id)
				foundReply2 := findReplyItem(resp.JSON200.Items, reply2Resp.JSON200.Id)

				r.NotNil(foundThread, "should find the thread")
				r.NotNil(foundReply1, "should find first reply")
				r.NotNil(foundReply2, "should find second reply")
			})
		}))
	}))
}

func findReplyItem(items []openapi.DatagraphItem, id openapi.Identifier) *openapi.DatagraphItemReply {
	for _, item := range items {
		if replyItem, err := item.AsDatagraphItemReply(); err == nil {
			if replyItem.Ref.Id == id {
				return &replyItem
			}
		}
	}
	return nil
}
