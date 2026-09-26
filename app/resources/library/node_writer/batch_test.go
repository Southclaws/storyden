package node_writer_test

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/Southclaws/lexorank"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/node"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

func TestBatchCreateAndUpdate(t *testing.T) {
	integration.Test(t, nil, fx.Invoke(func(lc fx.Lifecycle, ctx context.Context, writer *node_writer.Writer, aw *account_writer.Writer, db *ent.Client) {
		lc.Append(fx.StartHook(func() {
			ctx, owner := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)

			for _, size := range []int{1, 20, 100} {
				t.Run(fmt.Sprint(size), func(t *testing.T) {
					items := make([]node_writer.CreateInput, size)
					ids := make([]xid.ID, size)
					for i := range items {
						ids[i] = xid.New()
						items[i] = node_writer.CreateInput{
							ID:   library.NodeID(ids[i]),
							Name: fmt.Sprintf("Page %d", i),
							Slug: mark.NewSlugFromName(fmt.Sprintf("page-%d-%d", size, i)),
							Tags: opt.New(tag_ref.Names{tag_ref.NewName("shared")}),
						}
					}

					require.NoError(t, writer.CreateMany(ctx, owner.ID, items))
					rows, err := db.Node.Query().Where(node.IDIn(ids...)).All(ctx)
					require.NoError(t, err)
					require.Len(t, rows, size)

					rowsByID := make(map[xid.ID]*ent.Node, len(rows))
					for _, row := range rows {
						rowsByID[row.ID] = row
					}

					updates := make([]node_writer.UpdateInput, size)
					for i, id := range ids {
						row := rowsByID[id]
						require.NotNil(t, row)
						assert.Equal(t, fmt.Sprintf("Page %d", i), row.Name)

						updates[i] = node_writer.UpdateInput{
							ID:        library.NodeID(row.ID),
							UpdatedAt: row.UpdatedAt,
							Versioned: true,
							Options:   []node_writer.Option{node_writer.WithName(fmt.Sprintf("Changed %d", i))},
							Tags:      opt.New(tag_ref.Names{tag_ref.NewName("replacement")}),
						}
					}

					require.NoError(t, writer.UpdateMany(ctx, updates))
					rows, err = db.Node.Query().Where(node.IDIn(ids...)).WithTags().All(ctx)
					require.NoError(t, err)
					rowsByID = make(map[xid.ID]*ent.Node, len(rows))
					for _, row := range rows {
						rowsByID[row.ID] = row
					}

					for i, id := range ids {
						row := rowsByID[id]
						require.NotNil(t, row)
						assert.Equal(t, fmt.Sprintf("Changed %d", i), row.Name)
						assert.Equal(t, fmt.Sprintf("page-%d-%d", size, i), row.Slug)
						require.Len(t, row.Edges.Tags, 1)
						assert.Equal(t, "replacement", row.Edges.Tags[0].Name)
					}
				})
			}
		}))
	}))
}

func TestBatchRejectsStaleUpdatesAndNewDrafts(t *testing.T) {
	integration.Test(t, nil, fx.Invoke(func(lc fx.Lifecycle, ctx context.Context, writer *node_writer.Writer, aw *account_writer.Writer, db *ent.Client) {
		lc.Append(fx.StartHook(func() {
			ctx, owner := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)
			first := db.Node.Create().SetOwnerID(xid.ID(owner.ID)).SetName("First").SetSlug("first").SaveX(ctx)
			second := db.Node.Create().SetOwnerID(xid.ID(owner.ID)).SetName("Second").SetSlug("second").SaveX(ctx)
			updates := []node_writer.UpdateInput{
				{ID: library.NodeID(first.ID), UpdatedAt: first.UpdatedAt, Versioned: true, Options: []node_writer.Option{node_writer.WithName("Lost edit")}},
				{ID: library.NodeID(second.ID), UpdatedAt: second.UpdatedAt, Versioned: true, Options: []node_writer.Option{node_writer.WithName("Must roll back")}},
			}

			db.Node.UpdateOneID(first.ID).SetName("Concurrent edit").ExecX(ctx)
			require.ErrorContains(t, writer.UpdateMany(ctx, updates), "pages changed")
			assert.Equal(t, "Concurrent edit", db.Node.GetX(ctx, first.ID).Name)
			assert.Equal(t, "Second", db.Node.GetX(ctx, second.ID).Name)

			updates[0].UpdatedAt = db.Node.GetX(ctx, first.ID).UpdatedAt
			db.NodeVersion.Create().SetNodeID(first.ID).SetAuthorID(xid.ID(owner.ID)).SetName("Draft").SetSlug("draft").SaveX(ctx)
			require.ErrorContains(t, writer.UpdateMany(ctx, updates), "pages changed")
			assert.Equal(t, "Concurrent edit", db.Node.GetX(ctx, first.ID).Name)
			assert.Equal(t, "Second", db.Node.GetX(ctx, second.ID).Name)
		}))
	}))
}

func TestBatchRebalancesExhaustedSiblings(t *testing.T) {
	integration.Test(t, nil, fx.Invoke(func(lc fx.Lifecycle, ctx context.Context, writer *node_writer.Writer, aw *account_writer.Writer, db *ent.Client) {
		lc.Append(fx.StartHook(func() {
			ctx, owner := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)
			parent := db.Node.Create().SetOwnerID(xid.ID(owner.ID)).SetName("Parent").SetSlug("parent").SaveX(ctx)
			old := db.Node.Create().SetOwnerID(xid.ID(owner.ID)).SetParentID(parent.ID).SetName("Existing").SetSlug("existing").SetSort(lexorank.Top).SaveX(ctx)
			items := make([]node_writer.CreateInput, 3)
			for i := range items {
				items[i] = node_writer.CreateInput{
					ID:      library.NodeID(xid.New()),
					Name:    fmt.Sprintf("Sibling %d", i),
					Slug:    mark.NewSlugFromName(fmt.Sprintf("sibling-%d", i)),
					Options: []node_writer.Option{node_writer.WithParent(library.NodeID(parent.ID))},
				}
			}

			require.NoError(t, writer.CreateMany(ctx, owner.ID, items))
			rows := db.Node.Query().Where(node.ParentNodeID(parent.ID)).AllX(ctx)
			require.Len(t, rows, 4)
			slices.SortFunc(rows, func(a, b *ent.Node) int { return a.Sort.Compare(b.Sort) })
			assert.Equal(t, old.ID, rows[0].ID)
			for i := 1; i < len(rows); i++ {
				assert.Positive(t, rows[i].Sort.Compare(rows[i-1].Sort))
			}
		}))
	}))
}
