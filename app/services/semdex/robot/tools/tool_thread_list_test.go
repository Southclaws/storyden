package tools_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	robottools "github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	"github.com/Southclaws/storyden/internal/ent"
	entpost "github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/lib/mcp"
)

func TestThreadListRecentContent(t *testing.T) {
	integration.Test(t, nil, fx.Invoke(func(lc fx.Lifecycle, ctx context.Context, registry *robottools.Registry, aw *account_writer.Writer, db *ent.Client) {
		lc.Append(fx.StartHook(func() {
			ctx, account := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)
			ctx = session.WithAccountPermissions(ctx, *account, rbac.NewList(rbac.PermissionReadPublishedThreads))
			_, other := e2e.WithAccount(ctx, aw, seed.Account_002_Frigg)
			list, err := registry.GetTool(ctx, "thread_list")
			require.NoError(t, err)
			since := time.Date(2026, 10, 9, 14, 0, 0, 0, time.UTC)
			before := time.Date(2026, 10, 16, 14, 0, 0, 0, time.UTC)
			keyword := "recent-" + xid.New().String()
			create := func(name string, created time.Time, visibility entpost.Visibility) *ent.Post {
				p, err := db.Post.Create().SetAuthorID(xid.ID(other.ID)).SetTitle(keyword + " " + name).SetSlug(xid.New().String()).SetBody("<p>Discussion content</p>").SetShort("Discussion content").SetVisibility(visibility).SetCreatedAt(created).SetUpdatedAt(created).SetLastReplyAt(created).Save(ctx)
				require.NoError(t, err)
				return p
			}
			old := create("old but active", since.Add(-time.Second), entpost.VisibilityPublished)
			lower := create("at lower boundary", since, entpost.VisibilityPublished)
			middle := create("middle pinned", since.Add(24*time.Hour), entpost.VisibilityPublished)
			tieA := create("newest A", before.Add(-time.Hour), entpost.VisibilityPublished)
			tieB := create("newest B", before.Add(-time.Hour).In(time.FixedZone("UTC-8", -8*60*60)), entpost.VisibilityPublished)
			upper := create("at upper boundary", before, entpost.VisibilityPublished)
			future := create("future", before.Add(time.Hour), entpost.VisibilityPublished)
			create("private draft", before.Add(-time.Minute), entpost.VisibilityDraft)
			create("review", before.Add(-time.Minute), entpost.VisibilityReview)
			deleted := create("deleted", before.Add(-time.Minute), entpost.VisibilityPublished)
			_, err = db.Post.UpdateOneID(deleted.ID).SetDeletedAt(before).Save(ctx)
			require.NoError(t, err)
			_, err = db.Post.UpdateOneID(old.ID).SetLastReplyAt(before.Add(time.Hour)).SetPinnedRank(10).Save(ctx)
			require.NoError(t, err)
			_, err = db.Post.UpdateOneID(middle.ID).SetPinnedRank(5).SetLastReplyAt(before).Save(ctx)
			require.NoError(t, err)
			_, err = db.Post.Create().SetAuthorID(xid.ID(other.ID)).SetRootPostID(old.ID).SetBody(keyword + " recent reply").SetShort(keyword).SetVisibility(entpost.VisibilityPublished).SetCreatedAt(before.Add(-time.Minute)).SetUpdatedAt(before.Add(-time.Minute)).SetLastReplyAt(before.Add(-time.Minute)).Save(ctx)
			require.NoError(t, err)

			call := func(input string) mcp.ToolThreadListOutput {
				raw, err := list.Handler(ctx, json.RawMessage(input))
				require.NoError(t, err)
				var result mcp.ToolThreadListOutput
				require.NoError(t, json.Unmarshal(raw, &result))
				return result
			}
			ids := func(result mcp.ToolThreadListOutput) []string {
				out := make([]string, len(result.Threads))
				for i, thread := range result.Threads {
					out[i] = thread.Id
				}
				return out
			}
			base := fmt.Sprintf(`"query":%q,"visibility":"published","created_since":"2026-10-09T15:00:00+01:00","created_before":"2026-10-16T16:00:00+02:00"`, keyword)
			t.Run("date_filter_before_pagination_and_pins", func(t *testing.T) {
				first := call(`{` + base + `,"sort":"newest","page_size":2}`)
				assert.Equal(t, []string{tieB.ID.String(), tieA.ID.String()}, ids(first))
				assert.Equal(t, 4, first.TotalResults)
				assert.Equal(t, 2, first.TotalPages)
				require.NotNil(t, first.NextPage)
				assert.Equal(t, 2, *first.NextPage)
				assert.True(t, first.Threads[0].CreatedAt.Equal(before.Add(-time.Hour)))
				assert.True(t, first.Threads[0].UpdatedAt.Equal(before.Add(-time.Hour)))
				assert.Contains(t, first.Threads[0].BrowserUrl, tieB.Slug)
				second := call(`{` + base + `,"sort":"newest","page_size":2,"page":2}`)
				assert.Equal(t, []string{middle.ID.String(), lower.ID.String()}, ids(second))
				assert.Equal(t, 4, second.TotalResults)
				assert.Nil(t, second.NextPage)
				beyond := call(`{` + base + `,"sort":"newest","page_size":2,"page":3}`)
				assert.Empty(t, beyond.Threads)
				assert.Equal(t, 4, beyond.TotalResults)
				assert.Nil(t, beyond.NextPage)
			})
			t.Run("newest_without_date_bounds", func(t *testing.T) {
				result := call(fmt.Sprintf(`{"query":%q,"visibility":"published","sort":"newest","page_size":5}`, keyword))
				assert.Equal(t, []string{future.ID.String(), upper.ID.String(), tieB.ID.String(), tieA.ID.String(), middle.ID.String()}, ids(result))
				assert.Equal(t, 7, result.TotalResults)
				assert.Nil(t, result.CreatedSince)
				assert.Nil(t, result.CreatedBefore)
			})
			t.Run("oldest_and_default_feed_order", func(t *testing.T) {
				oldest := call(`{` + base + `,"sort":"oldest"}`)
				assert.Equal(t, []string{lower.ID.String(), middle.ID.String(), tieA.ID.String(), tieB.ID.String()}, ids(oldest))
				feed := call(`{` + base + `}`)
				assert.Equal(t, middle.ID.String(), feed.Threads[0].Id)
			})
			t.Run("effective_bounds_and_visibility", func(t *testing.T) {
				bounded := call(fmt.Sprintf(`{"query":%q,"created_since":"2026-10-09T14:00:00Z","created_before":"2026-10-16T14:00:00Z","sort":"newest"}`, keyword))
				assert.Equal(t, []string{tieB.ID.String(), tieA.ID.String(), middle.ID.String(), lower.ID.String()}, ids(bounded))
				require.NotNil(t, bounded.CreatedSince)
				require.NotNil(t, bounded.CreatedBefore)
				assert.True(t, since.Equal(*bounded.CreatedSince))
				assert.True(t, before.Equal(*bounded.CreatedBefore))
				hidden := call(fmt.Sprintf(`{"query":%q,"visibility":"draft","created_since":"2026-10-09T14:00:00Z","created_before":"2026-10-16T14:00:00Z","sort":"newest"}`, keyword))
				assert.Equal(t, ids(bounded), ids(hidden))
			})
			t.Run("empty_window", func(t *testing.T) {
				result := call(fmt.Sprintf(`{"query":%q,"sort":"newest","created_since":"2026-10-25T00:00:00Z","created_before":"2026-11-01T00:00:00Z"}`, keyword))
				assert.Empty(t, result.Threads)
				assert.Zero(t, result.TotalResults)
				assert.Zero(t, result.TotalPages)
				assert.Nil(t, result.NextPage)
			})
			t.Run("invalid_inputs", func(t *testing.T) {
				for _, input := range []string{
					`{"created_since":"2026-10-16T14:00:00Z","created_before":"2026-10-16T14:00:00Z"}`,
					`{"created_since":"2026-10-17T00:00:00Z","created_before":"2026-10-16T00:00:00Z"}`,
					`{"created_since":"yesterday"}`,
					`{"page_size":0}`, `{"page_size":101}`, `{"page":0}`, `{"sort":"popular"}`,
				} {
					_, err := list.Handler(ctx, json.RawMessage(input))
					require.Error(t, err, input)
				}
			})
		}))
	}))
}
