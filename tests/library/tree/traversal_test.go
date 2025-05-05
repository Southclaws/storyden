package library_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestNodesTreeQuerying(t *testing.T) {
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

			// SETUP
			//
			// node1       <- root               has 2 children: node2 and node3
			// |- node2    <- child of node1     has no children
			// |- node3    <- child of node1     has 1 child: node4
			//    |- node4 <- child of node3     has no children

			name1 := "test-node-1"
			slug1 := name1 + uuid.NewString()
			node1, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:       name1,
				Slug:       &slug1,
				Visibility: &visibility,
			}, sh.WithSession(ctx))
			tests.Ok(t, err, node1)

			name2 := "test-node-2"
			slug2 := name2 + uuid.NewString()
			node2, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:       name2,
				Slug:       &slug2,
				Parent:     &slug1,
				Visibility: &visibility,
			}, sh.WithSession(ctx))
			tests.Ok(t, err, node2)

			name3 := "test-node-3"
			slug3 := name3 + uuid.NewString()
			node3, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:       name3,
				Slug:       &slug3,
				Parent:     &slug1,
				Visibility: &visibility,
			}, sh.WithSession(ctx))
			tests.Ok(t, err, node3)

			name4 := "test-node-4"
			slug4 := name4 + uuid.NewString()
			node4, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:       name4,
				Slug:       &slug4,
				Parent:     &slug3,
				Visibility: &visibility,
			}, sh.WithSession(ctx))
			tests.Ok(t, err, node4)

			t.Run("query_all_top_level", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				depth := "0"
				response, err := cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
					Depth: &depth,
				})
				tests.Ok(t, err, response)

				ids := dt.Map(response.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })
				a.Contains(ids, node1.JSON200.Id)
				a.NotContains(ids, node2.JSON200.Id, "must not contain node2 because it's a child of node1 and thus not considered root level")

				node1resp, found := lo.Find(response.JSON200.Nodes, func(c openapi.NodeWithChildren) bool { return c.Id == node1.JSON200.Id })
				a.True(found)
				r.Len(node1resp.Children, 0, "depth is 0 so we should get no children")
			})

			t.Run("query_depth_1", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				depth := "1"
				response, err := cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
					Depth: &depth,
				})
				tests.Ok(t, err, response)

				ids := dt.Map(response.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })
				a.Contains(ids, node1.JSON200.Id)
				a.NotContains(ids, node2.JSON200.Id, "must not contain node2 because it's a child of node1 and thus not considered root level")

				node1resp, found := lo.Find(response.JSON200.Nodes, func(c openapi.NodeWithChildren) bool { return c.Id == node1.JSON200.Id })
				a.True(found)
				r.Len(node1resp.Children, 2, "node1 has two children: node2 and node3")

				node3resp, found := lo.Find(node1resp.Children, func(c openapi.NodeWithChildren) bool { return c.Id == node3.JSON200.Id })
				a.True(found)
				r.Len(node3resp.Children, 0, "node3 has one child but depth 1 does not include children of children")
			})

			t.Run("query_all_with_children", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				depth := "2"
				response, err := cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
					Depth: &depth,
				})
				tests.Ok(t, err, response)

				ids := dt.Map(response.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })
				a.Contains(ids, node1.JSON200.Id)
				a.NotContains(ids, node2.JSON200.Id, "must not contain node2 because it's a child of node1 and thus not considered root level")

				node1resp, found := lo.Find(response.JSON200.Nodes, func(c openapi.NodeWithChildren) bool { return c.Id == node1.JSON200.Id })
				a.True(found)
				r.Len(node1resp.Children, 2, "node1 has two children: node2 and node3")

				node3resp, found := lo.Find(node1resp.Children, func(c openapi.NodeWithChildren) bool { return c.Id == node3.JSON200.Id })
				a.True(found)
				r.Len(node3resp.Children, 1, "node3 has one child: node4")

				node4resp := node3resp.Children[0]
				a.Equal(node4.JSON200.Id, node4resp.Id)

				for _, n := range response.JSON200.Nodes {
					a.Nil(n.Parent, "root nodes must not have a parent")
				}
			})

			t.Run("query_node1_with_children", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				depth := "2"
				response, err := cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
					Depth:  &depth,
					NodeId: &node1.JSON200.Id,
				})
				tests.Ok(t, err, response)

				r.Len(response.JSON200.Nodes, 1, "the top level of any node list query with a specific node ID should only contain the node itself")
				n1 := response.JSON200.Nodes[0]
				a.Equal(n1.Id, node1.JSON200.Id)
				a.Nil(n1.Parent, "node1 is a root node")

				r.Len(n1.Children, 2, "node1 has two children: node2 and node3")
				ids := dt.Map(n1.Children, func(c openapi.NodeWithChildren) string { return c.Id })

				a.Contains(ids, node2.JSON200.Id, "node2 is a child of node1")
				a.Contains(ids, node3.JSON200.Id, "node3 is a child of node1")

				n3, n3found := lo.Find(n1.Children, func(c openapi.NodeWithChildren) bool { return c.Id == node3.JSON200.Id })
				r.True(n3found, "node3 is a child of node1")

				r.Len(n3.Children, 1, "node3 has one child: node4")
			})

			t.Run("query_node2_with_children", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				depth := "2"
				response, err := cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
					Depth:  &depth,
					NodeId: &node2.JSON200.Id,
				})
				tests.Ok(t, err, response)

				r.Len(response.JSON200.Nodes, 1, "must return node2 itself")
				n2 := response.JSON200.Nodes[0]
				a.Equal(n2.Id, node2.JSON200.Id)

				r.NotNil(n2.Parent, "the query is a subtree under node2 so it must have parent information")
				a.Equal(n2.Parent.Id, node1.JSON200.Id, "node2's parent is node1")

				r.Len(n2.Children, 0, "node2 has no children")
			})

			t.Run("query_node3_with_children", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				depth := "2"
				response, err := cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
					Depth:  &depth,
					NodeId: &node3.JSON200.Id,
				})
				tests.Ok(t, err, response)

				r.Len(response.JSON200.Nodes, 1, "must return node3 itself")
				n3 := response.JSON200.Nodes[0]
				a.Equal(n3.Id, node3.JSON200.Id)

				r.NotNil(n3.Parent)
				a.Equal(n3.Parent.Id, node1.JSON200.Id)

				r.Len(n3.Children, 1, "node3 has one child")
				n4 := n3.Children[0]
				a.Equal(n4.Id, node4.JSON200.Id)
				a.Nil(n4.Parent, "node4 appears in the children list of node3 so it must not have a parent field set since this would just be duplicated information")

				r.Len(n4.Children, 0, "node4 has no children")
			})
		}))
	}))
}

func TestNodesTreeQuerying_WithHiddenChildNodes(t *testing.T) {
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
			session := sh.WithSession(ctx)

			visibility := openapi.Published

			// SETUP
			//
			// node1       <- root               has 2 children: node2 and node3
			// |- node2    <- child of node1     has no children
			// |- node3    <- child of node1     has 2 children: node4 and node5 - is set to hide child nodes in tree queries so 4 and 5 should not appear
			//    |- node4 <- child of node3     has no children
			//    |- node5 <- child of node3     has no children

			name1 := "test-node-1"
			slug1 := name1 + uuid.NewString()
			node1, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:       name1,
				Slug:       &slug1,
				Visibility: &visibility,
			}, sh.WithSession(ctx))
			tests.Ok(t, err, node1)

			name2 := "test-node-2"
			slug2 := name2 + uuid.NewString()
			node2, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:       name2,
				Slug:       &slug2,
				Parent:     &slug1,
				Visibility: &visibility,
			}, sh.WithSession(ctx))
			tests.Ok(t, err, node2)

			name3 := "test-node-3"
			slug3 := name3 + uuid.NewString()
			node3, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:       name3,
				Slug:       &slug3,
				Parent:     &slug1,
				Visibility: &visibility,
			}, sh.WithSession(ctx))
			tests.Ok(t, err, node3)

			name4 := "test-node-4"
			slug4 := name4 + uuid.NewString()
			node4, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:       name4,
				Slug:       &slug4,
				Parent:     &slug3,
				Visibility: &visibility,
			}, sh.WithSession(ctx))
			tests.Ok(t, err, node4)

			name5 := "test-node-5"
			slug5 := name5 + uuid.NewString()
			node5, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:       name5,
				Slug:       &slug5,
				Parent:     &slug3,
				Visibility: &visibility,
			}, sh.WithSession(ctx))
			tests.Ok(t, err, node5)

			t.Run("query_tree_without_hidden_children", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				depth := "3"
				{
					list := tests.AssertRequest(
						cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
							Depth: &depth,
						}),
					)(t, http.StatusOK)
					node1, found := lo.Find(list.JSON200.Nodes, func(c openapi.NodeWithChildren) bool { return c.Id == node1.JSON200.Id })
					r.True(found)
					node1children := node1.Children
					r.Len(node1children, 2)
					node3resp, found := lo.Find(node1children, func(c openapi.NodeWithChildren) bool { return c.Id == node3.JSON200.Id })
					r.True(found)
					a.Len(node3resp.Children, 2, "has 2 non-hidden children: node4 and node5")
				}

				// Mark Node 3 as hidden children in tree traversal
				updated := tests.AssertRequest(
					cl.NodeUpdateWithResponse(ctx, node3.JSON200.Slug, &openapi.NodeUpdateParams{}, openapi.NodeMutableProps{
						HideChildTree: opt.New(true).Ptr(),
					}, session),
				)(t, http.StatusOK)
				r.True(updated.JSON200.HideChildTree)

				get := tests.AssertRequest(
					cl.NodeGetWithResponse(ctx, node3.JSON200.Slug, &openapi.NodeGetParams{}),
				)(t, http.StatusOK)
				r.True(get.JSON200.HideChildTree, "node3 should have HideChildTree set to true")

				{
					list := tests.AssertRequest(
						cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
							Depth: &depth,
						}),
					)(t, http.StatusOK)
					node1, found := lo.Find(list.JSON200.Nodes, func(c openapi.NodeWithChildren) bool { return c.Id == node1.JSON200.Id })
					r.True(found)
					node1children := node1.Children
					r.Len(node1children, 2)
					node3resp, found := lo.Find(node1children, func(c openapi.NodeWithChildren) bool { return c.Id == node3.JSON200.Id })
					r.True(found)
					a.True(node3resp.HideChildTree)
					a.Len(node3resp.Children, 0, "has 2 now-hidden children: node4 and node5")
				}
			})
		}))
	}))
}
