package thread_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/thread_querier"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestThreadListHasReplies(t *testing.T) {
	if tests.IsSharedPostgresDatabase() {
		t.Skip("skipping test on shared postgres")
	}

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
			authorCtx, author := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			otherCtx, _ := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			authorSession := sh.WithSession(authorCtx)
			otherSession := sh.WithSession(otherCtx)
			adminSession := sh.WithSession(adminCtx)
			prefix := "replyfilter" + xid.New().String()

			category := tests.AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Name:        prefix,
				Description: "Reply filter fixtures",
				Colour:      "#abcdef",
			}, adminSession))(t, http.StatusOK).JSON200

			createThread := func(name string, categoryID *string, state openapi.Visibility, session openapi.RequestEditorFn) *openapi.Thread {
				t.Helper()
				return tests.AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      prefix + " " + name,
					Body:       new("<p>A discussion</p>"),
					Category:   categoryID,
					Visibility: new(state),
				}, session))(t, http.StatusOK).JSON200
			}

			createReply := func(thread *openapi.Thread, replyTo *string, session openapi.RequestEditorFn) *openapi.Reply {
				t.Helper()
				return tests.AssertRequest(cl.ReplyCreateWithResponse(root, thread.Slug, openapi.ReplyInitialProps{
					Body:    "<p>A response</p>",
					ReplyTo: replyTo,
				}, session))(t, http.StatusOK).JSON200
			}

			withoutReplies := createThread("no replies", new(category.Id), openapi.VisibilityPublished, authorSession)
			resolved := createThread("resolved without replies", new(category.Id), openapi.VisibilityPublished, authorSession)
			deletedOnly := createThread("deleted replies", new(category.Id), openapi.VisibilityPublished, authorSession)
			outside := createThread("uncategorised", nil, openapi.VisibilityPublished, authorSession)
			selfReply := createThread("self reply", new(category.Id), openapi.VisibilityPublished, authorSession)
			otherReply := createThread("other reply", new(category.Id), openapi.VisibilityPublished, authorSession)
			nestedReply := createThread("nested reply", new(category.Id), openapi.VisibilityPublished, authorSession)
			hiddenReply := createThread("hidden reply", new(category.Id), openapi.VisibilityPublished, authorSession)
			privateThread := createThread("private root", new(category.Id), openapi.VisibilityReview, otherSession)

			createReply(selfReply, nil, authorSession)
			createReply(otherReply, nil, otherSession)
			deleted := createReply(deletedOnly, nil, authorSession)
			parent := createReply(nestedReply, nil, authorSession)
			createReply(nestedReply, new(parent.Id), otherSession)
			hidden := createReply(hiddenReply, nil, otherSession)
			for _, reply := range []*openapi.Reply{deleted, parent} {
				r.NoError(db.Post.UpdateOneID(openapi.ParseID(reply.Id)).SetDeletedAt(time.Now()).Exec(root))
			}
			r.NoError(db.Post.UpdateOneID(openapi.ParseID(hidden.Id)).SetVisibility(ent_post.VisibilityReview).Exec(root))
			tests.AssertRequest(cl.ThreadUpdateWithResponse(root, resolved.Slug, openapi.ThreadMutableProps{
				LifecycleState: new(openapi.ThreadLifecycleStateResolved),
			}, authorSession))(t, http.StatusOK)
			tests.AssertRequest(cl.ThreadReadStatusUpdateWithResponse(root, withoutReplies.Slug, openapi.ThreadReadStatusMutableProps{
				Status: openapi.ThreadReadStatusMutablePropsStatusRead,
			}, authorSession))(t, http.StatusOK)

			list := func(t *testing.T, params openapi.ThreadListParams, session ...openapi.RequestEditorFn) *openapi.ThreadListResult {
				t.Helper()
				params.Q = new(prefix)
				return tests.AssertRequest(cl.ThreadListWithResponse(root, &params, session...))(t, http.StatusOK).JSON200
			}
			ids := func(result *openapi.ThreadListResult) []string {
				values := make([]string, 0, len(result.Threads))
				for _, thread := range result.Threads {
					values = append(values, thread.Id)
				}
				return values
			}

			t.Run("omitted_false_true_and_counts", func(t *testing.T) {
				r := require.New(t)
				all := list(t, openapi.ThreadListParams{}, authorSession)
				r.ElementsMatch([]string{withoutReplies.Id, resolved.Id, deletedOnly.Id, outside.Id, selfReply.Id, otherReply.Id, nestedReply.Id, hiddenReply.Id}, ids(all))
				for _, thread := range all.Threads {
					if thread.Id == hiddenReply.Id || thread.Id == nestedReply.Id {
						r.Equal(1, thread.ReplyStatus.Replies)
					}
				}

				unanswered := list(t, openapi.ThreadListParams{HasReplies: new(false)}, authorSession)
				r.ElementsMatch([]string{withoutReplies.Id, resolved.Id, deletedOnly.Id, outside.Id}, ids(unanswered))
				r.Equal(4, unanswered.Results)
				for _, thread := range unanswered.Threads {
					r.Zero(thread.ReplyStatus.Replies)
				}

				answered := list(t, openapi.ThreadListParams{HasReplies: new(true)}, authorSession)
				r.ElementsMatch([]string{selfReply.Id, otherReply.Id, nestedReply.Id, hiddenReply.Id}, ids(answered))
				r.Equal(4, answered.Results)
			})

			t.Run("combines_with_category_read_and_lifecycle_filters", func(t *testing.T) {
				r := require.New(t)
				params := openapi.ThreadListParams{
					HasReplies: new(false),
					Categories: new([]string{category.Slug}),
				}
				r.ElementsMatch([]string{withoutReplies.Id, resolved.Id, deletedOnly.Id}, ids(list(t, params, authorSession)))
				params.ReadStatus = new(openapi.ThreadReadStatusFilterRead)
				r.Equal([]string{withoutReplies.Id}, ids(list(t, params, authorSession)))
				params.ReadStatus = nil
				params.LifecycleState = new(openapi.ThreadLifecycleStateResolved)
				r.Equal([]string{resolved.Id}, ids(list(t, params, authorSession)))
			})

			t.Run("preserves_thread_visibility", func(t *testing.T) {
				r := require.New(t)
				params := openapi.ThreadListParams{HasReplies: new(false)}
				r.NotContains(ids(list(t, params)), privateThread.Id)
				r.NotContains(ids(list(t, params, authorSession)), privateThread.Id)
				r.Contains(ids(list(t, params, adminSession)), privateThread.Id)
			})

			t.Run("counts_matching_threads_before_pagination", func(t *testing.T) {
				r := require.New(t)
				for _, hasReplies := range []bool{false, true} {
					var actual []post.ID
					for page := range 2 {
						result, err := tq.List(root, page, 2, opt.New(author.ID),
							thread_querier.HasKeyword(prefix),
							thread_querier.HasNotBeenDeleted(),
							thread_querier.HasStatus(visibility.VisibilityPublished),
							thread_querier.HasReplies(hasReplies),
							thread_querier.HasNoPinnedOrdering(true),
						)
						r.NoError(err)
						r.Equal(4, result.TotalResults)
						r.Equal(2, result.TotalPages)
						r.Len(result.Threads, 2)
						r.Equal(page == 0, result.NextPage.Ok())
						for _, thread := range result.Threads {
							r.Equal(hasReplies, thread.ReplyStatus.Count > 0)
							actual = append(actual, thread.ID)
						}
					}
					expected := ids(list(t, openapi.ThreadListParams{HasReplies: new(hasReplies)}, authorSession))
					actualIDs := make([]string, 0, len(actual))
					for _, id := range actual {
						actualIDs = append(actualIDs, xid.ID(id).String())
					}
					r.ElementsMatch(expected, actualIDs)
				}
			})
		}))
	}))
}
