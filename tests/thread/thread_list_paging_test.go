package thread_test

import (
	"context"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/post/thread_querier"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestThreadListPagesWhenThreadsShareATimestamp(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		db *ent.Client,
		tq *thread_querier.Querier,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			ctx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			session := sh.WithSession(ctx)

			const threadCount = 6
			created := map[xid.ID]struct{}{}

			for range threadCount {
				thread, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       opt.New("<p>paged</p>").Ptr(),
					Visibility: opt.New(openapi.VisibilityPublished).Ptr(),
					Title:      "thread paging " + uuid.NewString(),
				}, session)
				tests.Ok(t, err, thread)

				created[openapi.ParseID(thread.JSON200.Id)] = struct{}{}
			}

			// every row on one timestamp, which is what an import or a busy
			// second looks like, and what makes an offset page unstable
			ids := make([]xid.ID, 0, threadCount)
			for id := range created {
				ids = append(ids, id)
			}

			shared := time.Now().Add(-time.Hour).UTC()
			_, err := db.Post.Update().Where(ent_post.IDIn(ids...)).SetLastReplyAt(shared).Save(ctx)
			r.NoError(err)

			seen := map[xid.ID]int{}

			for page := 0; page < 10; page++ {
				result, err := tq.List(ctx, page, 2, opt.NewEmpty[account.AccountID]())
				r.NoError(err)
				r.Equal(page, result.CurrentPage)

				for _, thread := range result.Threads {
					seen[xid.ID(thread.ID)]++
				}

				if _, ok := result.NextPage.Get(); !ok {
					break
				}
			}

			for id := range created {
				a.Equal(1, seen[id], "thread %s should appear on exactly one page", id)
			}
		}))
	}))
}
