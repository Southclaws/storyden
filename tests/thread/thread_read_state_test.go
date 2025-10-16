package thread_test

import (
	"context"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/post_read_state"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/tests"
)

func TestThreadReadState(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		rw *post_read_state.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			acc1ctx, acc1 := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			acc2ctx, _ := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			session1 := sh.WithSession(acc1ctx)
			session2 := sh.WithSession(acc2ctx)

			catName := "Category " + uuid.NewString()

			catCreate, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Colour:      "#fe4efd",
				Description: "read state testing",
				Name:        catName,
			}, session1)
			tests.Ok(t, err, catCreate)

			t.Run("read_state", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadCreate, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       opt.New("<p>original thread</p>").Ptr(),
					Category:   opt.New(catCreate.JSON200.Id).Ptr(),
					Visibility: opt.New(openapi.Published).Ptr(),
					Title:      "Read State Test Thread",
				}, session1)
				tests.Ok(t, err, threadCreate)

				threadGet1, err := cl.ThreadGetWithResponse(acc1ctx, threadCreate.JSON200.Slug, nil, session1)
				tests.Ok(t, err, threadGet1)

				// Verify LIST before beacon - should have no ReadStatus
				threadListBefore, err := cl.ThreadListWithResponse(acc1ctx, &openapi.ThreadListParams{}, session1)
				tests.Ok(t, err, threadListBefore)
				listThreadBefore, found := lo.Find(threadListBefore.JSON200.Threads, func(t openapi.ThreadReference) bool {
					return t.Id == threadCreate.JSON200.Id
				})
				a.True(found)
				a.Nil(listThreadBefore.ReadStatus)

				// Verify GET before beacon - should have no ReadStatus
				threadGetBefore, err := cl.ThreadGetWithResponse(acc1ctx, threadCreate.JSON200.Slug, nil, session1)
				tests.Ok(t, err, threadGetBefore)
				a.Nil(threadGetBefore.JSON200.ReadStatus)

				// Simulate a page-load, send a beacon (navigator.sendBeacon) to mark the thread as read in its current state.
				err = rw.UpsertReadState(root, acc1.ID, post.ID(utils.Must(xid.FromString(threadCreate.JSON200.Id))))
				r.NoError(err)

				lastRead := time.Now().UTC()

				// Sleep for >1 second to ensure unixepoch() will show different values
				// we don't store millisecond level precision in the database.
				time.Sleep(1100 * time.Millisecond)

				replyCreate, err := cl.ReplyCreateWithResponse(acc2ctx, threadCreate.JSON200.Slug, openapi.ReplyInitialProps{
					Body: "this is a reply from acc2",
				}, session2)
				tests.Ok(t, err, replyCreate)

				threadList, err := cl.ThreadListWithResponse(acc1ctx, &openapi.ThreadListParams{}, session1)
				tests.Ok(t, err, threadList)
				listThread, found := lo.Find(threadList.JSON200.Threads, func(t openapi.ThreadReference) bool {
					return t.Id == threadCreate.JSON200.Id
				})
				a.True(found)
				r.NotNil(listThread.ReadStatus)
				a.Equal(1, listThread.ReadStatus.RepliesSince)

				threadGet2, err := cl.ThreadGetWithResponse(acc1ctx, threadCreate.JSON200.Slug, nil, session1)
				tests.Ok(t, err, threadGet2)
				a.Equal(1, threadGet2.JSON200.ReadStatus.RepliesSince)
				a.WithinDuration(lastRead, threadGet2.JSON200.ReadStatus.LastReadAt, time.Second)
			})
		}))
	}))
}
