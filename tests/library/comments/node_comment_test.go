package comments_test

import (
	"context"
	"testing"

	"github.com/Southclaws/dt"
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

func TestNodeCommentsCRUD(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			a := assert.New(t)

			ctx1, acc1 := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)
			ctx2, acc2 := e2e.WithAccount(ctx, aw, seed.Account_002_Frigg)

			published := openapi.Published
			name := "test-node-comments-" + uuid.NewString()
			slug := name
			content := "<h1>Node for Comments</h1><p>Testing comment threads.</p>"

			node, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:       name,
				Slug:       &slug,
				Content:    &content,
				Visibility: &published,
			}, sh.WithSession(ctx1))
			tests.Ok(t, err, node)
			a.Equal(acc1.ID.String(), string(node.JSON200.Owner.Id))

			t.Run("create_and_list", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				commentBody := "<p>First comment on this node!</p>"
				commentResp, err := cl.NodeCommentCreateWithResponse(ctx, slug, openapi.ThreadInitialProps{
					Body: opt.New(commentBody).Ptr(),
				}, sh.WithSession(ctx1))
				tests.Ok(t, err, commentResp)

				r.NotNil(commentResp.JSON200)
				a.Equal("<body><p>First comment on this node!</p></body>", commentResp.JSON200.Body)
				a.Equal(acc1.ID.String(), string(commentResp.JSON200.Author.Id))

				listResp, err := cl.NodeCommentListWithResponse(ctx, slug, &openapi.NodeCommentListParams{})
				tests.Ok(t, err, listResp)

				r.NotNil(listResp.JSON200)
				a.Equal(1, listResp.JSON200.Results)
				a.Len(listResp.JSON200.Threads, 1)
				a.Equal(commentResp.JSON200.Id, listResp.JSON200.Threads[0].Id)
				a.Equal(openapi.Unlisted, commentResp.JSON200.Visibility)
			})

			t.Run("multiple_comments", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				body1 := "<p>Second comment</p>"
				comment1, err := cl.NodeCommentCreateWithResponse(ctx, slug, openapi.ThreadInitialProps{
					Body: opt.New(body1).Ptr(),
				}, sh.WithSession(ctx2))
				tests.Ok(t, err, comment1)

				body2 := "<p>Third comment</p>"
				comment2, err := cl.NodeCommentCreateWithResponse(ctx, slug, openapi.ThreadInitialProps{
					Body: opt.New(body2).Ptr(),
				}, sh.WithSession(ctx1))
				tests.Ok(t, err, comment2)

				listResp, err := cl.NodeCommentListWithResponse(ctx, slug, &openapi.NodeCommentListParams{})
				tests.Ok(t, err, listResp)

				r.NotNil(listResp.JSON200)
				a.GreaterOrEqual(listResp.JSON200.Results, 3)

				ids := dt.Map(listResp.JSON200.Threads, func(t openapi.ThreadReference) string { return string(t.Id) })
				a.Contains(ids, string(comment1.JSON200.Id))
				a.Contains(ids, string(comment2.JSON200.Id))
			})

			t.Run("reply_to_comment", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				commentBody := "<p>Original comment for replies</p>"
				comment, err := cl.NodeCommentCreateWithResponse(ctx, slug, openapi.ThreadInitialProps{
					Body:       opt.New(commentBody).Ptr(),
					Visibility: &published,
				}, sh.WithSession(ctx1))
				tests.Ok(t, err, comment)

				replyBody := "Reply to the comment"
				reply, err := cl.ReplyCreateWithResponse(ctx, comment.JSON200.Slug, openapi.ReplyInitialProps{
					Body: replyBody,
				}, sh.WithSession(ctx2))
				tests.Ok(t, err, reply)

				r.NotNil(reply.JSON200)
				a.Equal("<body>Reply to the comment</body>", reply.JSON200.Body)
				a.Equal(acc2.ID.String(), string(reply.JSON200.Author.Id))

				threadGet, err := cl.ThreadGetWithResponse(ctx, comment.JSON200.Slug, &openapi.ThreadGetParams{})
				tests.Ok(t, err, threadGet)

				r.NotNil(threadGet.JSON200)
				a.Equal(1, threadGet.JSON200.Replies.Results)
			})

			t.Run("update_comment_via_thread_api", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				commentBody := "<p>Comment to update</p>"
				comment, err := cl.NodeCommentCreateWithResponse(ctx, slug, openapi.ThreadInitialProps{
					Body: opt.New(commentBody).Ptr(),
				}, sh.WithSession(ctx1))
				tests.Ok(t, err, comment)

				newBody := "<p>Updated comment body</p>"
				update, err := cl.ThreadUpdateWithResponse(ctx, comment.JSON200.Slug, openapi.ThreadMutableProps{
					Body: opt.New(newBody).Ptr(),
				}, sh.WithSession(ctx1))
				tests.Ok(t, err, update)

				r.NotNil(update.JSON200)
				a.Equal("<body><p>Updated comment body</p></body>", update.JSON200.Body)
			})
		}))
	}))
}

func TestNodeCommentsPermissions(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			ctx1, acc1 := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)

			t.Run("cannot_comment_on_unpublished_node", func(t *testing.T) {
				r := require.New(t)

				draft := openapi.Draft
				name := "unpublished-node-" + uuid.NewString()
				slug := name
				content := "<p>Unpublished node</p>"

				node, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Name:       name,
					Slug:       &slug,
					Content:    &content,
					Visibility: &draft,
				}, sh.WithSession(ctx1))
				tests.Ok(t, err, node)
				r.Equal(acc1.ID.String(), string(node.JSON200.Owner.Id))

				commentBody := "<p>Should not be allowed</p>"
				commentResp, err := cl.NodeCommentCreateWithResponse(ctx, slug, openapi.ThreadInitialProps{
					Body: opt.New(commentBody).Ptr(),
				}, sh.WithSession(ctx1))
				r.NoError(err)
				r.NotNil(commentResp.JSONDefault)
			})
		}))
	}))
}
