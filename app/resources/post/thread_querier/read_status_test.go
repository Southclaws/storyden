package thread_querier_test

import (
	"context"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/post_read_state"
	"github.com/Southclaws/storyden/app/resources/post/thread_querier"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

func TestReadStatusFilter(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		aw *account_writer.Writer,
		rw *post_read_state.Writer,
		tq *thread_querier.Querier,
		db *ent.Client,
	) {
		lc.Append(fx.StartHook(func() {
			_, member := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			_, other := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			baseline := time.Now().UTC().Add(-time.Hour).Truncate(time.Second)

			createThread := func(t *testing.T, title string) *ent.Post {
				t.Helper()
				thread, err := db.Post.Create().
					SetTitle(title).
					SetSlug(xid.New().String()).
					SetAuthorID(xid.ID(other.ID)).
					SetVisibility(ent_post.VisibilityPublished).
					SetBody("<p>A discussion</p>").
					SetShort(title).
					SetCreatedAt(baseline.Add(-time.Hour)).
					SetUpdatedAt(baseline).
					SetLastReplyAt(baseline).
					Save(root)
				require.NoError(t, err)
				return thread
			}

			createReply := func(t *testing.T, thread *ent.Post, at time.Time, deleted bool) {
				t.Helper()
				create := db.Post.Create().
					SetRootID(thread.ID).
					SetAuthorID(xid.ID(other.ID)).
					SetVisibility(ent_post.VisibilityPublished).
					SetBody("<p>A reply</p>").
					SetShort("A reply").
					SetCreatedAt(at).
					SetUpdatedAt(at).
					SetLastReplyAt(at)
				if deleted {
					create.SetDeletedAt(at.Add(time.Minute))
				}

				_, err := create.Save(root)
				require.NoError(t, err)
				require.NoError(t, db.Post.UpdateOneID(thread.ID).SetLastReplyAt(at).Exec(root))
			}

			list := func(t *testing.T, prefix string, reader account.AccountID, unread bool, page, size int) *thread_querier.Result {
				t.Helper()
				result, err := tq.List(root, page, size, opt.New(reader),
					thread_querier.HasKeyword(prefix),
					thread_querier.HasNotBeenDeleted(),
					thread_querier.HasStatus(visibility.VisibilityPublished),
					thread_querier.HasNoPinnedOrdering(true),
					thread_querier.HasReadStatus(reader, unread),
				)
				require.NoError(t, err)
				return result
			}

			t.Run("matches_reply_counts_and_timestamp_boundaries", func(t *testing.T) {
				r := require.New(t)
				prefix := "readquery" + xid.New().String()
				unvisited := createThread(t, prefix+" unvisited")
				caughtUp := createThread(t, prefix+" caught up")
				marked := createThread(t, prefix+" marked")
				newReplies := createThread(t, prefix+" new replies")
				deletedReply := createThread(t, prefix+" deleted reply")
				oldReply := createThread(t, prefix+" old reply")
				equalReply := createThread(t, prefix+" equal reply")
				for _, thread := range []*ent.Post{caughtUp, marked, newReplies, deletedReply, oldReply, equalReply} {
					r.NoError(rw.MarkRead(root, member.ID, post.ID(thread.ID), baseline))
				}

				r.NoError(rw.MarkRead(root, other.ID, post.ID(unvisited.ID), baseline))
				r.NoError(rw.MarkUnread(root, other.ID, post.ID(caughtUp.ID), baseline))
				r.NoError(rw.MarkUnread(root, member.ID, post.ID(marked.ID), baseline.Add(time.Minute)))
				for range 3 {
					createReply(t, newReplies, baseline.Add(time.Minute), false)
				}

				createReply(t, deletedReply, baseline.Add(time.Minute), true)
				createReply(t, oldReply, baseline.Add(-time.Minute), false)
				createReply(t, equalReply, baseline, false)

				unread := list(t, prefix, member.ID, true, 0, 10)
				var unreadIDs []post.ID
				for _, thread := range unread.Threads {
					unreadIDs = append(unreadIDs, thread.ID)
					if thread.ID == post.ID(newReplies.ID) {
						r.Equal(3, thread.ReadStatus.OrZero().Count)
					}
				}

				r.ElementsMatch([]post.ID{post.ID(unvisited.ID), post.ID(marked.ID), post.ID(newReplies.ID)}, unreadIDs)
				r.Equal(3, unread.TotalResults)
				read := list(t, prefix, member.ID, false, 0, 10)
				var readIDs []post.ID
				for _, thread := range read.Threads {
					readIDs = append(readIDs, thread.ID)
					r.Zero(thread.ReadStatus.OrZero().Count)
				}

				r.ElementsMatch([]post.ID{post.ID(caughtUp.ID), post.ID(deletedReply.ID), post.ID(oldReply.ID), post.ID(equalReply.ID)}, readIDs)
				r.Equal(4, read.TotalResults)
				otherRead := list(t, prefix, other.ID, false, 0, 10)
				r.Len(otherRead.Threads, 1)
				r.Equal(post.ID(unvisited.ID), otherRead.Threads[0].ID)
			})

			t.Run("counts_and_pages_filtered_rows_with_equal_activity", func(t *testing.T) {
				r := require.New(t)
				prefix := "readquerypages" + xid.New().String()
				var expected []post.ID
				for index := range 7 {
					thread := createThread(t, prefix)
					expected = append([]post.ID{post.ID(thread.ID)}, expected...)
					if index == 0 {
						r.NoError(db.Post.UpdateOneID(thread.ID).SetPinnedRank(100).Exec(root))
					}
				}

				for range 5 {
					thread := createThread(t, prefix+" read")
					r.NoError(db.Post.UpdateOneID(thread.ID).SetLastReplyAt(baseline.Add(time.Minute)).Exec(root))
					r.NoError(rw.MarkRead(root, member.ID, post.ID(thread.ID), baseline))
				}

				var actual []post.ID
				for page := range 4 {
					result := list(t, prefix, member.ID, true, page, 2)
					r.Equal(7, result.TotalResults)
					r.Equal(4, result.TotalPages)
					r.Equal(page, result.CurrentPage)
					if page < 3 {
						r.Len(result.Threads, 2)
						r.Equal(page+1, result.NextPage.OrZero())
					} else {
						r.Len(result.Threads, 1)
						r.False(result.NextPage.Ok())
					}

					for _, thread := range result.Threads {
						actual = append(actual, thread.ID)
					}
				}

				r.Equal(expected, actual)
				r.Equal(5, list(t, prefix, member.ID, false, 0, 2).TotalResults)
			})
		}))
	}))
}
