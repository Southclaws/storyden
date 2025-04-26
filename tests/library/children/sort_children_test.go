package children_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/lexorank"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestNodeSorting(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			ctx, _ := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)

			visibility := openapi.Published

			makenode := func(name string, parent *string) *openapi.Node {
				slug := name + uuid.NewString()
				n := tests.AssertRequest(
					cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
						Name: name, Slug: &slug, Visibility: &visibility,
						Parent: parent,
					}, sh.WithSession(ctx)),
				)(t, http.StatusOK)

				return n.JSON200
			}

			getList := func(ids ...string) []openapi.NodeWithChildren {
				listResponse := tests.AssertRequest(
					cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
						Depth:  opt.New("2").Ptr(),
						Format: opt.New(openapi.NodeListParamsFormatFlat).Ptr(),
					}),
				)(t, http.StatusOK)
				list := dt.Filter(listResponse.JSON200.Nodes, func(n openapi.NodeWithChildren) bool {
					return lo.Contains(ids, n.Id)
				})
				return list
			}

			getNode := func(parent string) openapi.NodeWithChildren {
				listResponse := tests.AssertRequest(
					cl.NodeGetWithResponse(ctx, parent, &openapi.NodeGetParams{}),
				)(t, http.StatusOK)
				return *listResponse.JSON200
			}

			getChildren := func(parent string) []openapi.NodeWithChildren {
				listResponse := tests.AssertRequest(
					cl.NodeListChildrenWithResponse(ctx, parent, &openapi.NodeListChildrenParams{}),
				)(t, http.StatusOK)
				return listResponse.JSON200.Nodes
			}

			listIDs := func(list []openapi.NodeWithChildren) []string {
				return dt.Map(list, func(n openapi.NodeWithChildren) string { return n.Id })
			}

			// NOTE:
			// These tests don't work well in a shared database because root
			// level nodes are shared by all tests and another test that runs
			// parallel could insert a new node and change the order if a sort
			// key gets normalised during insertion (shifted around by maybe
			// hundreds or thousands of places) resulting in a different result.
			// However, the child-based sort tests do cover all sorting cases
			// and are not affected by this as they live in isolated sorting
			// contexts. This test can run in SQLite mode where each test runs
			// with a separate isolated database. The only case not covered now
			// is when parentID = nil, but I have tested those cases by removing
			// the comments from the below tests. Source: trust me bro™.
			//
			// t.Run("root_move_node", func(t *testing.T) {
			// 	a := assert.New(t)
			// 	r := require.New(t)

			// 	// create 3 root level nodes
			// 	n1 := makenode("1", nil)
			// 	n2 := makenode("2", nil)
			// 	n3 := makenode("3", nil)

			// 	list := getList(n1.Id, n2.Id, n3.Id)
			// 	r.Len(list, 3)
			// 	a.Equal([]string{n1.Id, n2.Id, n3.Id}, listIDs(list), "newly inserted nodes should be at the end")

			// 	// move n2 to middle, after n1
			// 	resp := tests.AssertRequest(
			// 		cl.NodeUpdatePositionWithResponse(ctx, n1.Slug, openapi.NodeUpdatePositionJSONRequestBody{
			// 			After: &n2.Id,
			// 		}),
			// 	)(t, http.StatusOK)
			// 	r.NotNil(resp.JSON200)

			// 	ids := listIDs(getList(n1.Id, n2.Id, n3.Id))
			// 	r.Len(ids, 3)
			// 	a.Equal([]string{n2.Id, n1.Id, n3.Id}, ids, "node 1 has been moved to after node 3")

			// 	// move n1 to middle, before n2
			// 	resp = tests.AssertRequest(
			// 		cl.NodeUpdatePositionWithResponse(ctx, n1.Slug, openapi.NodeUpdatePositionJSONRequestBody{
			// 			Before: &n3.Id,
			// 		}),
			// 	)(t, http.StatusOK)
			// 	r.NotNil(resp.JSON200)

			// 	ids = listIDs(getList(n1.Id, n2.Id, n3.Id))
			// 	r.Len(ids, 3)
			// 	a.Equal([]string{n2.Id, n1.Id, n3.Id}, ids, "node 1 has been moved to after node 3")

			// 	// move n1 to bottom, after n3
			// 	resp = tests.AssertRequest(
			// 		cl.NodeUpdatePositionWithResponse(ctx, n1.Slug, openapi.NodeUpdatePositionJSONRequestBody{
			// 			After: &n3.Id,
			// 		}),
			// 	)(t, http.StatusOK)
			// 	r.NotNil(resp.JSON200)

			// 	ids = listIDs(getList(n1.Id, n2.Id, n3.Id))
			// 	r.Len(ids, 3)
			// 	a.Equal([]string{n2.Id, n3.Id, n1.Id}, ids, "node 1 has been moved to after node 3")

			// 	// move n3 to top, before n2
			// 	resp = tests.AssertRequest(
			// 		cl.NodeUpdatePositionWithResponse(ctx, n3.Slug, openapi.NodeUpdatePositionJSONRequestBody{
			// 			Before: &n2.Id,
			// 		}),
			// 	)(t, http.StatusOK)
			// 	r.NotNil(resp.JSON200)

			// 	ids = listIDs(getList(n1.Id, n2.Id, n3.Id))
			// 	r.Len(ids, 3)
			// 	a.Equal([]string{n3.Id, n2.Id, n1.Id}, ids, "node 3 has been moved to before node 1")
			// })

			t.Run("new_child_node_insertion_order", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				// create 3 nodes under p
				p := makenode("parent", nil)
				n1 := makenode("1", &p.Slug)
				n2 := makenode("2", &p.Slug)
				n3 := makenode("3", &p.Slug)

				list := getList(n1.Id, n2.Id, n3.Id)
				r.Len(list, 3)
				a.Equal([]string{n1.Id, n2.Id, n3.Id}, listIDs(list), "newly inserted nodes should be at the end")
			})

			t.Run("child_move_child_to_top", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				// create 3 child level nodes under p
				p := makenode("parent", nil)
				n1 := makenode("1", &p.Slug)
				n2 := makenode("2", &p.Slug)
				n3 := makenode("3", &p.Slug)

				resp := tests.AssertRequest(
					cl.NodeUpdatePositionWithResponse(ctx, n3.Slug, openapi.NodeUpdatePositionJSONRequestBody{
						Before: &n1.Id,
					}),
				)(t, http.StatusOK)
				r.NotNil(resp.JSON200)

				wantOrder := []string{n3.Id, n1.Id, n2.Id}

				ids := listIDs(getList(n1.Id, n2.Id, n3.Id))
				r.Len(ids, 3)
				a.Equal(wantOrder, ids, "node 3 has been moved to before node 1")

				ids = listIDs(getNode(p.Id).Children)
				r.Len(ids, 3)
				a.Equal(wantOrder, ids, "node 3 has been moved to before node 1")

				ids = listIDs(getChildren(p.Id))
				r.Len(ids, 3)
				a.Equal(wantOrder, ids, "node 3 has been moved to before node 1")
			})

			t.Run("child_move_child_to_bottom", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				// create 3 child level nodes under p
				p := makenode("parent", nil)
				n1 := makenode("1", &p.Slug)
				n2 := makenode("2", &p.Slug)
				n3 := makenode("3", &p.Slug)

				resp := tests.AssertRequest(
					cl.NodeUpdatePositionWithResponse(ctx, n1.Slug, openapi.NodeUpdatePositionJSONRequestBody{
						After: &n3.Id,
					}),
				)(t, http.StatusOK)
				r.NotNil(resp.JSON200)

				wantOrder := []string{n2.Id, n3.Id, n1.Id}

				ids := listIDs(getList(n1.Id, n2.Id, n3.Id))
				r.Len(ids, 3)
				a.Equal(wantOrder, ids, "node 1 has been moved to after node 3")

				ids = listIDs(getNode(p.Id).Children)
				r.Len(ids, 3)
				a.Equal(wantOrder, ids, "node 1 has been moved to after node 3")

				ids = listIDs(getChildren(p.Id))
				r.Len(ids, 3)
				a.Equal(wantOrder, ids, "node 1 has been moved to after node 3")
			})

			t.Run("child_move_child_to_middle_before", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				// create 3 child level nodes under p
				p := makenode("parent", nil)
				n1 := makenode("1", &p.Slug)
				n2 := makenode("2", &p.Slug)
				n3 := makenode("3", &p.Slug)

				resp := tests.AssertRequest(
					cl.NodeUpdatePositionWithResponse(ctx, n1.Slug, openapi.NodeUpdatePositionJSONRequestBody{
						Before: &n3.Id,
					}),
				)(t, http.StatusOK)
				r.NotNil(resp.JSON200)

				wantOrder := []string{n2.Id, n1.Id, n3.Id}

				ids := listIDs(getList(n1.Id, n2.Id, n3.Id))
				r.Len(ids, 3)
				a.Equal(wantOrder, ids, "node 1 has been moved to before node 3")

				ids = listIDs(getNode(p.Id).Children)
				r.Len(ids, 3)
				a.Equal(wantOrder, ids, "node 1 has been moved to before node 3")

				ids = listIDs(getChildren(p.Id))
				r.Len(ids, 3)
				a.Equal(wantOrder, ids, "node 1 has been moved to before node 3")
			})

			t.Run("child_move_child_to_middle_after", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				// create 3 child level nodes under p
				p := makenode("parent", nil)
				n1 := makenode("1", &p.Slug)
				n2 := makenode("2", &p.Slug)
				n3 := makenode("3", &p.Slug)

				resp := tests.AssertRequest(
					cl.NodeUpdatePositionWithResponse(ctx, n1.Slug, openapi.NodeUpdatePositionJSONRequestBody{
						After: &n2.Id,
					}),
				)(t, http.StatusOK)
				r.NotNil(resp.JSON200)

				wantOrder := []string{n2.Id, n1.Id, n3.Id}

				ids := listIDs(getList(n1.Id, n2.Id, n3.Id))
				r.Len(ids, 3)
				a.Equal(wantOrder, ids, "node 1 has been moved to after node 2")

				ids = listIDs(getNode(p.Id).Children)
				r.Len(ids, 3)
				a.Equal(wantOrder, ids, "node 1 has been moved to before node 2")

				ids = listIDs(getChildren(p.Id))
				r.Len(ids, 3)
				a.Equal(wantOrder, ids, "node 1 has been moved to before node 2")
			})
		}))
	}))
}

func TestNodeSortKeyNormalise(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		db *ent.Client,
	) {
		lc.Append(fx.StartHook(func() {
			ctx, _ := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)

			visibility := openapi.Published

			makenode := func(name string, parent *string) *openapi.Node {
				slug := name + uuid.NewString()
				n := tests.AssertRequest(
					cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
						Name: name, Slug: &slug, Visibility: &visibility,
						Parent: parent,
					}, sh.WithSession(ctx)),
				)(t, http.StatusOK)

				return n.JSON200
			}

			setNodeSort := func(id openapi.Identifier, k lexorank.Key) {
				nid, _ := xid.FromString(id)
				if err := db.Node.UpdateOneID(nid).SetSort(k).Exec(ctx); err != nil {
					t.Fatal(err)
				}
			}

			t.Run("trigger_top_key_normalise_1", func(t *testing.T) {
				t.Parallel()

				r := require.New(t)

				n1 := makenode("1", nil)

				k, _ := lexorank.ParseKey("0|zzzzzz")

				setNodeSort(n1.Id, *k)

				n2 := makenode("2", nil)

				n2id, _ := xid.FromString(n2.Id)
				n2e, err := db.Node.Get(ctx, n2id)
				r.NoError(err)
				r.NotNil(n2e)

				fmt.Println(n2e.Sort)
			})
			t.Run("trigger_top_key_normalise_2", func(t *testing.T) {
				t.Parallel()

				r := require.New(t)

				n1 := makenode("1", nil)

				k, _ := lexorank.ParseKey("0|zzzzzz")

				setNodeSort(n1.Id, *k)

				n2 := makenode("2", nil)

				n2id, _ := xid.FromString(n2.Id)
				n2e, err := db.Node.Get(ctx, n2id)
				r.NoError(err)
				r.NotNil(n2e)

				fmt.Println(n2e.Sort)
			})
		}))
	}))
}
