package library_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_traversal"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	entnode "github.com/Southclaws/storyden/internal/ent/node"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestNodeListFlatPagination(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		aw *account_writer.Writer,
		sh *e2e.SessionHelper,
		db *ent.Client,
		traversal node_traversal.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			ownerCtx, owner := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)
			memberCtx, member := e2e.WithAccount(ctx, aw, seed.Account_003_Baldur)
			created := time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
			root, err := db.Node.Create().
				SetAccountID(xid.ID(owner.ID)).SetName("paginated tree").
				SetSlug(xid.New().String()).SetVisibility(entnode.VisibilityPublished).
				SetCreatedAt(created).Save(ctx)
			require.NoError(t, err)

			create := func(parent *ent.Node, offset int, vis entnode.Visibility) *ent.Node {
				n, err := db.Node.Create().
					SetAccountID(xid.ID(owner.ID)).SetName("page").
					SetSlug(xid.New().String()).SetParentNodeID(parent.ID).
					SetVisibility(vis).SetCreatedAt(created.Add(time.Duration(offset) * time.Hour)).Save(ctx)
				require.NoError(t, err)
				return n
			}

			expected := []string{root.ID.String()}
			for i := 1; i <= 101; i++ {
				n := create(root, i, entnode.VisibilityPublished)
				expected = append(expected, n.ID.String())
			}
			privateIDs := []string{}
			for _, vis := range []entnode.Visibility{entnode.VisibilityDraft, entnode.VisibilityReview, entnode.VisibilityUnlisted} {
				n := create(root, -1, vis)
				privateIDs = append(privateIDs, n.ID.String())
			}
			memberDraft := create(root, -1, entnode.VisibilityDraft)
			memberReview := create(root, -1, entnode.VisibilityReview)
			_, err = db.Node.Update().Where(entnode.IDIn(memberDraft.ID, memberReview.ID)).SetAccountID(xid.ID(member.ID)).Save(ctx)
			require.NoError(t, err)
			hidden := create(root, 102, entnode.VisibilityPublished)
			_, err = db.Node.UpdateOne(hidden).SetHideChildTree(true).Save(ctx)
			require.NoError(t, err)
			expected = append(expected, hidden.ID.String())
			hiddenChild := create(hidden, -1, entnode.VisibilityPublished)
			deleted := create(root, -1, entnode.VisibilityPublished)
			create(deleted, -1, entnode.VisibilityPublished)
			_, err = db.Node.UpdateOne(deleted).SetDeletedAt(created).Save(ctx)
			require.NoError(t, err)

			params := openapi.NodeListParams{
				NodeId:     opt.New(root.ID.String()).Ptr(),
				Format:     opt.New(openapi.NodeListParamsFormatFlat).Ptr(),
				Visibility: &openapi.VisibilityParam{openapi.VisibilityPublished},
			}
			ids := func(result *openapi.NodeListResult) []string {
				return dt.Map(result.Nodes, func(n openapi.NodeWithChildren) string { return n.Id })
			}

			t.Run("visibility_and_branch_pruning_before_pagination", func(t *testing.T) {
				params := params
				params.Page = opt.New("1").Ptr()
				first := tests.AssertRequest(cl.NodeListWithResponse(ctx, &params))(t, http.StatusOK).JSON200
				assert.Equal(t, expected[:100], ids(first))
				assert.Equal(t, 100, first.PageSize)
				assert.Equal(t, 100, first.Results)
				assert.Equal(t, 2, first.TotalPages)
				assert.Equal(t, 1, first.CurrentPage)
				require.NotNil(t, first.NextPage)
				assert.Equal(t, 2, *first.NextPage)

				params.Page = opt.New("2").Ptr()
				second := tests.AssertRequest(cl.NodeListWithResponse(ctx, &params))(t, http.StatusOK).JSON200
				assert.Equal(t, expected[100:], ids(second))
				assert.Equal(t, 3, second.Results)
				assert.Equal(t, 2, second.TotalPages)
				assert.Nil(t, second.NextPage)

				params.Page = opt.New("3").Ptr()
				beyond := tests.AssertRequest(cl.NodeListWithResponse(ctx, &params))(t, http.StatusOK).JSON200
				assert.Empty(t, beyond.Nodes)
				assert.Zero(t, beyond.Results)
				assert.Equal(t, 2, beyond.TotalPages)
				assert.Nil(t, beyond.NextPage)
			})

			t.Run("unpaginated_flat_and_tree_requests_retain_behaviour", func(t *testing.T) {
				params := params
				flat := tests.AssertRequest(cl.NodeListWithResponse(ctx, &params))(t, http.StatusOK).JSON200
				assert.ElementsMatch(t, expected, ids(flat))
				assert.Zero(t, flat.PageSize)

				params.Format = opt.New(openapi.NodeListParamsFormatTree).Ptr()
				tree := tests.AssertRequest(cl.NodeListWithResponse(ctx, &params))(t, http.StatusOK).JSON200
				params.Page = opt.New("2").Ptr()
				withPage := tests.AssertRequest(cl.NodeListWithResponse(ctx, &params))(t, http.StatusOK).JSON200
				assert.Equal(t, tree, withPage)
			})

			t.Run("visibility_filters_respect_guest_member_and_owner_permissions", func(t *testing.T) {
				params := params
				params.Page = opt.New("1").Ptr()
				params.Visibility = &openapi.VisibilityParam{
					openapi.VisibilityPublished, openapi.VisibilityDraft,
					openapi.VisibilityReview, openapi.VisibilityUnlisted,
				}
				guest := tests.AssertRequest(cl.NodeListWithResponse(ctx, &params))(t, http.StatusOK).JSON200
				assert.Equal(t, expected[:100], ids(guest))
				member := tests.AssertRequest(cl.NodeListWithResponse(ctx, &params, sh.WithSession(memberCtx)))(t, http.StatusOK).JSON200
				assert.Equal(t, append([]string{memberDraft.ID.String(), memberReview.ID.String()}, expected[:98]...), ids(member))
				owner := tests.AssertRequest(cl.NodeListWithResponse(ctx, &params, sh.WithSession(ownerCtx)))(t, http.StatusOK).JSON200
				ownerIDs := append(privateIDs, memberReview.ID.String())
				assert.Equal(t, append(ownerIDs, expected[:96]...), ids(owner))
			})

			t.Run("depth_filter_applies_before_totals", func(t *testing.T) {
				params := params
				params.Depth = opt.New("0").Ptr()
				params.Page = opt.New("1").Ptr()
				result := tests.AssertRequest(cl.NodeListWithResponse(ctx, &params))(t, http.StatusOK).JSON200
				assert.Equal(t, []string{root.ID.String()}, ids(result))
				assert.Equal(t, 1, result.Results)
				assert.Equal(t, 1, result.TotalPages)
			})

			t.Run("deleted_subtree_is_empty", func(t *testing.T) {
				params := params
				params.NodeId = opt.New(deleted.ID.String()).Ptr()
				params.Page = opt.New("1").Ptr()
				result := tests.AssertRequest(cl.NodeListWithResponse(ctx, &params))(t, http.StatusOK).JSON200
				assert.Empty(t, result.Nodes)
				assert.Zero(t, result.Results)
				assert.Zero(t, result.TotalPages)
				assert.Nil(t, result.NextPage)
			})

			t.Run("repository_list_still_includes_hidden_children", func(t *testing.T) {
				result, err := traversal.List(ctx, pagination.NewPageParams(1, 200), node_traversal.ListParams{
					Root: opt.New(library.NodeID(root.ID)),
				})
				require.NoError(t, err)
				assert.Equal(t, 104, result.TotalResults)
				assert.Contains(t, dt.Map(result.Items, func(n *library.Node) string { return n.GetID().String() }), hiddenChild.ID.String())
			})

			t.Run("hostile_authors_are_literal_values", func(t *testing.T) {
				payloads := []string{
					"' OR '1'='1",
					"' OR 1=1 --",
					"') OR 1=1 --",
					"' /* comment */ OR 1=1 --",
					"? OR 1=1 --",
					"$1 OR 1=1 --",
					"' UNION SELECT NULL, NULL, NULL, NULL, NULL --",
					`" OR "1"="1`,
					"`; SELECT 1; --",
					"'; SELECT pg_sleep(0.01); --",
					"'; UPDATE nodes SET name = 'sql-injected' WHERE id = '" + root.ID.String() + "'; --",
				}
				for _, payload := range payloads {
					t.Run(payload, func(t *testing.T) {
						for _, format := range []openapi.NodeListParamsFormat{openapi.NodeListParamsFormatTree, openapi.NodeListParamsFormatFlat} {
							params := params
							params.Format = &format
							params.Author = &payload
							result := tests.AssertRequest(cl.NodeListWithResponse(ctx, &params))(t, http.StatusOK).JSON200
							assert.Zero(t, len(result.Nodes))
						}
						params := params
						params.Page = opt.New("1").Ptr()
						params.Author = &payload
						result := tests.AssertRequest(cl.NodeListWithResponse(ctx, &params))(t, http.StatusOK).JSON200
						assert.Zero(t, len(result.Nodes))
						assert.Zero(t, result.TotalPages)
					})
				}
			})

			t.Run("repository_binds_authors_without_transport_validation", func(t *testing.T) {
				for _, payload := range []string{"' OR 1=1 --", "'; DROP TABLE nodes; --", "' UNION SELECT NULL --"} {
					t.Run(payload, func(t *testing.T) {
						result, err := traversal.List(ctx, pagination.NewPageParams(1, 100), node_traversal.ListParams{
							Root: opt.New(library.NodeID(root.ID)),
						}, node_traversal.WithRootOwner(payload))
						require.NoError(t, err)
						assert.Zero(t, len(result.Items))
						assert.Zero(t, result.TotalResults)
					})
				}
			})

			t.Run("hostile_pages_cannot_change_the_query", func(t *testing.T) {
				for _, payload := range []string{"1 OR 1=1", "1; SELECT pg_sleep(0.01); --", "1 UNION SELECT NULL", "0); DROP TABLE nodes; --", "-1", "18446744073709551616"} {
					t.Run(payload, func(t *testing.T) {
						params := params
						params.Page = &payload
						result := tests.AssertRequest(cl.NodeListWithResponse(ctx, &params))(t, http.StatusOK).JSON200
						assert.Equal(t, expected[:100], ids(result))
						assert.Equal(t, 1, result.CurrentPage)
						assert.Equal(t, 2, result.TotalPages)
					})
				}
			})

			t.Run("hostile_typed_filters_are_rejected", func(t *testing.T) {
				for _, field := range []string{"node_id", "depth", "visibility"} {
					for _, payload := range []string{"' OR 1=1 --", "0); DROP TABLE nodes; --", "' UNION SELECT NULL --"} {
						t.Run(field+"/"+payload, func(t *testing.T) {
							params := params
							params.Page = opt.New("1").Ptr()
							switch field {
							case "node_id":
								params.NodeId = &payload
							case "depth":
								params.Depth = &payload
							case "visibility":
								params.Visibility = &openapi.VisibilityParam{openapi.Visibility(payload)}
							}
							result, err := cl.NodeListWithResponse(ctx, &params)
							require.NoError(t, err)
							assert.Contains(t, []int{http.StatusBadRequest, http.StatusInternalServerError}, result.StatusCode())
							assert.Nil(t, result.JSON200)
						})
					}
				}
			})

			t.Run("injection_attempts_leave_the_fixture_intact", func(t *testing.T) {
				canary, err := db.Node.Get(ctx, root.ID)
				require.NoError(t, err)
				assert.Equal(t, "paginated tree", canary.Name)
				draft, err := db.Node.Get(ctx, memberDraft.ID)
				require.NoError(t, err)
				assert.Equal(t, entnode.VisibilityDraft, draft.Visibility)
				params := params
				params.Page = opt.New("1").Ptr()
				result := tests.AssertRequest(cl.NodeListWithResponse(ctx, &params))(t, http.StatusOK).JSON200
				assert.Equal(t, expected[:100], ids(result))
				assert.Equal(t, 2, result.TotalPages)
			})
		}))
	}))
}
