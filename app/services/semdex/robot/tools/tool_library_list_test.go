package tools_test

import (
	"context"
	"encoding/json"
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
	entnode "github.com/Southclaws/storyden/internal/ent/node"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/lib/mcp"
)

func TestLibraryPageListRecentContent(t *testing.T) {
	integration.Test(t, nil, fx.Invoke(func(lc fx.Lifecycle, ctx context.Context, registry *robottools.Registry, aw *account_writer.Writer, db *ent.Client) {
		lc.Append(fx.StartHook(func() {
			_, other := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)
			ctx, member := e2e.WithAccount(ctx, aw, seed.Account_004_Loki)
			ctx = session.WithAccountPermissions(ctx, *member, rbac.NewList(rbac.PermissionReadPublishedLibrary))
			list, err := registry.GetTool(ctx, "library_page_list")
			require.NoError(t, err)
			since := time.Date(2040, 10, 9, 14, 0, 0, 0, time.UTC)
			before := since.Add(7 * 24 * time.Hour)
			create := func(name string, created time.Time, visibility entnode.Visibility) *ent.Node {
				n, err := db.Node.Create().SetAccountID(xid.ID(other.ID)).SetName(name).SetSlug(xid.New().String()).SetVisibility(visibility).SetCreatedAt(created).SetUpdatedAt(created).Save(ctx)
				require.NoError(t, err)
				return n
			}
			old := create("Old directory with new children", since.Add(-time.Second), entnode.VisibilityPublished)
			_, err = db.Node.UpdateOne(old).SetHideChildTree(true).SetUpdatedAt(before).Save(ctx)
			require.NoError(t, err)
			lower := create("At lower boundary", since, entnode.VisibilityPublished)
			middle := create("Child of an old directory", since.Add(time.Hour), entnode.VisibilityPublished)
			_, err = db.Node.UpdateOne(middle).SetParentNodeID(old.ID).Save(ctx)
			require.NoError(t, err)
			tieA := create("Newest A", before.Add(-time.Hour), entnode.VisibilityPublished)
			tieB := create("Newest B", before.Add(-time.Hour).In(time.FixedZone("UTC-8", -8*60*60)), entnode.VisibilityPublished)
			upper := create("At upper boundary", before, entnode.VisibilityPublished)
			create("Someone else's draft", before.Add(-time.Minute), entnode.VisibilityDraft)
			create("Review submission", before.Add(-time.Minute), entnode.VisibilityReview)
			create("Unlisted", before.Add(-time.Minute), entnode.VisibilityUnlisted)
			own := create("Own draft", since.Add(2*time.Hour), entnode.VisibilityDraft)
			_, err = db.Node.UpdateOne(own).SetAccountID(xid.ID(member.ID)).Save(ctx)
			require.NoError(t, err)
			deleted := create("Deleted page", before.Add(-time.Minute), entnode.VisibilityPublished)
			_, err = db.Node.UpdateOne(deleted).SetDeletedAt(before).Save(ctx)
			require.NoError(t, err)

			call := func(input string) mcp.ToolLibraryPageTreeOutput {
				raw, err := list.Handler(ctx, json.RawMessage(input))
				require.NoError(t, err)
				var out mcp.ToolLibraryPageTreeOutput
				require.NoError(t, json.Unmarshal(raw, &out))
				return out
			}
			ids := func(out mcp.ToolLibraryPageTreeOutput) []string {
				result := make([]string, len(out.Pages))
				for i, page := range out.Pages {
					result[i] = page.Id
				}
				return result
			}
			bounds := `"created_since":"2040-10-09T15:00:00+01:00","created_before":"2040-10-16T16:00:00+02:00"`
			base := bounds + `,"visibility":"published"`
			t.Run("dates_and_sort_before_pagination", func(t *testing.T) {
				first := call(`{` + base + `,"sort":"newest","page_size":2}`)
				assert.Equal(t, []string{tieB.ID.String(), tieA.ID.String()}, ids(first))
				assert.Equal(t, 4, first.TotalResults)
				assert.Equal(t, 2, first.Results)
				assert.Equal(t, 2, first.TotalPages)
				assert.Equal(t, 1, first.CurrentPage)
				require.NotNil(t, first.NextPage)
				assert.Equal(t, 2, *first.NextPage)
				assert.True(t, before.Add(-time.Hour).Equal(first.Pages[0].CreatedAt))
				assert.True(t, before.Add(-time.Hour).Equal(first.Pages[0].UpdatedAt))
				require.NotNil(t, first.CreatedSince)
				require.NotNil(t, first.CreatedBefore)
				assert.True(t, since.Equal(*first.CreatedSince))
				assert.True(t, before.Equal(*first.CreatedBefore))
				second := call(`{` + base + `,"sort":"newest","page_size":2,"page":2}`)
				assert.Equal(t, []string{middle.ID.String(), lower.ID.String()}, ids(second))
				assert.Equal(t, 4, second.TotalResults)
				assert.Nil(t, second.NextPage)
				require.NotNil(t, second.Pages[0].Parent)
				assert.Equal(t, old.Slug, *second.Pages[0].Parent)
				beyond := call(`{` + base + `,"sort":"newest","page_size":2,"page":3}`)
				assert.Empty(t, beyond.Pages)
				assert.Equal(t, 4, beyond.TotalResults)
				assert.Nil(t, beyond.NextPage)
			})
			t.Run("latest_without_dates", func(t *testing.T) {
				latest := call(`{"visibility":"published","sort":"newest","page_size":3}`)
				assert.Equal(t, []string{upper.ID.String(), tieB.ID.String(), tieA.ID.String()}, ids(latest))
				assert.Nil(t, latest.CreatedSince)
				assert.Nil(t, latest.CreatedBefore)
			})
			t.Run("oldest_hierarchy_and_depth", func(t *testing.T) {
				oldest := call(`{` + base + `,"sort":"oldest"}`)
				assert.Equal(t, []string{lower.ID.String(), middle.ID.String(), tieA.ID.String(), tieB.ID.String()}, ids(oldest))
				hierarchy := call(`{` + base + `}`)
				assert.Equal(t, middle.ID.String(), hierarchy.Pages[len(hierarchy.Pages)-1].Id)
				roots := call(`{` + base + `,"depth":0,"sort":"oldest"}`)
				assert.Equal(t, []string{lower.ID.String(), tieA.ID.String(), tieB.ID.String()}, ids(roots))
				assert.Equal(t, 3, roots.TotalResults)
				unbounded := call(`{` + base + `,"depth":-1,"sort":"oldest"}`)
				assert.Equal(t, ids(oldest), ids(unbounded))
			})
			t.Run("visibility_before_counts", func(t *testing.T) {
				drafts := call(`{` + bounds + `,"visibility":"draft","sort":"newest"}`)
				assert.Equal(t, []string{own.ID.String()}, ids(drafts))
				assert.Equal(t, 1, drafts.TotalResults)
				all := call(`{` + bounds + `,"sort":"oldest"}`)
				assert.Equal(t, []string{lower.ID.String(), middle.ID.String(), own.ID.String(), tieA.ID.String(), tieB.ID.String()}, ids(all))
				assert.Equal(t, 5, all.TotalResults)
			})
			t.Run("empty_range", func(t *testing.T) {
				empty := call(`{"sort":"newest","created_since":"2041-01-01T00:00:00Z"}`)
				assert.Empty(t, empty.Pages)
				assert.Zero(t, empty.TotalResults)
				assert.Zero(t, empty.TotalPages)
				assert.Nil(t, empty.NextPage)
			})
			t.Run("invalid_inputs", func(t *testing.T) {
				for _, input := range []string{
					`{"created_since":"2040-10-16T00:00:00Z","created_before":"2040-10-16T00:00:00Z"}`,
					`{"created_since":"2040-10-17T00:00:00Z","created_before":"2040-10-16T00:00:00Z"}`,
					`{"created_since":"yesterday"}`, `{"page":0}`, `{"page_size":0}`, `{"page_size":101}`, `{"sort":"activity"}`, `{"depth":-2}`,
				} {
					_, err := list.Handler(ctx, json.RawMessage(input))
					require.Error(t, err, input)
				}
			})
		}))
	}))
}
