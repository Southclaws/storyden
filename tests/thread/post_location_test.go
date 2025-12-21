package thread_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestPostLocation(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			ctx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			session := sh.WithSession(ctx)

			catResp, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Name:   "PostLocation" + uuid.NewString(),
				Colour: "#123456",
			}, session)
			tests.Ok(t, err, catResp)

			createThread := func(currentT *testing.T, title string) *openapi.Thread {
				currentT.Helper()
				resp, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      title,
					Body:       opt.New(fmt.Sprintf("<p>%s</p>", title)).Ptr(),
					Category:   opt.New(catResp.JSON200.Id).Ptr(),
					Visibility: opt.New(openapi.Published).Ptr(),
				}, session)
				tests.Ok(currentT, err, resp)
				return resp.JSON200
			}

			createReply := func(currentT *testing.T, slug string, index int) openapi.Identifier {
				currentT.Helper()
				resp, err := cl.ReplyCreateWithResponse(root, slug, openapi.ReplyInitialProps{
					Body: fmt.Sprintf("reply-%d", index),
				}, session)
				tests.Ok(currentT, err, resp)
				return resp.JSON200.Id
			}

			t.Run("thread", func(t *testing.T) {
				t.Parallel()
				a := assert.New(t)
				r := require.New(t)

				thr := createThread(t, "Thread location")

				resp, err := cl.PostLocationGetWithResponse(root, &openapi.PostLocationGetParams{
					Id: thr.Id,
				}, session)
				tests.Ok(t, err, resp)

				r.NotNil(resp.JSON200)
				a.Equal(openapi.PostLocationKindThread, resp.JSON200.Kind)
				a.Equal(thr.Slug, resp.JSON200.Slug)
				r.Nil(resp.JSON200.Index)
				r.Nil(resp.JSON200.Page)
				r.Nil(resp.JSON200.Position)
			})

			t.Run("reply_single_page", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				thr := createThread(t, "Thread replies page1")
				var replies []openapi.Identifier
				for i := 0; i < 4; i++ {
					replies = append(replies, createReply(t, thr.Slug, i))
				}

				resp, err := cl.PostLocationGetWithResponse(root, &openapi.PostLocationGetParams{
					Id: replies[2],
				}, session)
				tests.Ok(t, err, resp)

				subslug := fmt.Sprintf("%s#%s", thr.Slug, string(replies[2]))

				r.NotNil(resp.JSON200)
				a.Equal(openapi.PostLocationKindReply, resp.JSON200.Kind)
				a.Equal(subslug, resp.JSON200.Slug)
				r.NotNil(resp.JSON200.Index)
				r.NotNil(resp.JSON200.Page)
				r.NotNil(resp.JSON200.Position)
				a.Equal(2, *resp.JSON200.Index)
				a.Equal(1, *resp.JSON200.Page)
				a.Equal(2, *resp.JSON200.Position)
			})

			t.Run("reply_second_page", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				thr := createThread(t, "Thread replies page2")
				targetIndex := reply.RepliesPerPage + 1
				limit := reply.RepliesPerPage + 3
				var target openapi.Identifier
				for i := 0; i < limit; i++ {
					id := createReply(t, thr.Slug, i)
					if i == targetIndex {
						target = id
					}
				}

				r.NotEmpty(target)

				resp, err := cl.PostLocationGetWithResponse(root, &openapi.PostLocationGetParams{
					Id: target,
				}, session)
				tests.Ok(t, err, resp)

				subslug := fmt.Sprintf("%s#%s", thr.Slug, string(target))

				r.NotNil(resp.JSON200)
				a.Equal(openapi.PostLocationKindReply, resp.JSON200.Kind)
				a.Equal(subslug, resp.JSON200.Slug)
				r.NotNil(resp.JSON200.Index)
				r.NotNil(resp.JSON200.Page)
				r.NotNil(resp.JSON200.Position)
				a.Equal(51, *resp.JSON200.Index)
				a.Equal(2, *resp.JSON200.Page)
				a.Equal(1, *resp.JSON200.Position)
			})

			t.Run("not_found", func(t *testing.T) {
				randomID := openapi.Identifier(xid.New().String())
				resp, err := cl.PostLocationGetWithResponse(root, &openapi.PostLocationGetParams{Id: randomID}, session)
				tests.Status(t, err, resp, http.StatusNotFound)
			})
		}))
	}))
}
