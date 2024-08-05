package library_test

import (
	"context"
	"testing"

	"github.com/Southclaws/dt"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/cookie"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestNodesTreeMutations(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		cj *cookie.Jar,
		aw account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			ctx, _ := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)

			visibility := openapi.Published

			// SETUP
			//
			// node1       <- root               has no children
			// node2       <- root               has 1 child: node3
			// |- node3    <- child of node1     has 1 child: node4
			//    |- node4 <- child of node3     has no children

			name1 := "test-node-1"
			slug1 := name1 + uuid.NewString()
			node1, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:       name1,
				Slug:       &slug1,
				Visibility: &visibility,
			}, e2e.WithSession(ctx, cj))
			tests.Ok(t, err, node1)

			name2 := "test-node-2"
			slug2 := name2 + uuid.NewString()
			node2, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:       name2,
				Slug:       &slug2,
				Visibility: &visibility,
			}, e2e.WithSession(ctx, cj))
			tests.Ok(t, err, node2)

			name3 := "test-node-3"
			slug3 := name3 + uuid.NewString()
			node3, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:       name3,
				Slug:       &slug3,
				Parent:     &slug2,
				Visibility: &visibility,
			}, e2e.WithSession(ctx, cj))
			tests.Ok(t, err, node3)

			name4 := "test-node-4"
			slug4 := name4 + uuid.NewString()
			node4, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:       name4,
				Slug:       &slug4,
				Parent:     &slug3,
				Visibility: &visibility,
			}, e2e.WithSession(ctx, cj))
			tests.Ok(t, err, node4)

			t.Run("change_tree_structure", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				name5 := "test-node-5"
				slug5 := name5 + uuid.NewString()
				node5, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Name:       name5,
					Slug:       &slug5,
					Visibility: &visibility,
				}, e2e.WithSession(ctx, cj))
				tests.Ok(t, err, node5)

				// List nodes and check that the new node is in the root nodes

				listresponse1, err := cl.NodeListWithResponse(ctx, &openapi.NodeListParams{})
				tests.Ok(t, err, listresponse1)

				ids1 := dt.Map(listresponse1.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })
				a.Contains(ids1, node5.JSON200.Id, "node5 should be in the root nodes as it currently has no parent")

				// Add node5 as a child of node1

				addresponse, err := cl.NodeAddNodeWithResponse(ctx, slug1, slug5, e2e.WithSession(ctx, cj))
				tests.Ok(t, err, addresponse)
				a.Equal(node1.JSON200.Id, addresponse.JSON200.Id)

				// Current situation:
				// node1       <- root               has 1 child: node5
				// |- node5    <- child of node1     has no children
				// node2       <- root               has 1 child: node3
				// |- node3    <- child of node1     has 1 child: node4
				//    |- node4 <- child of node3     has no children

				listresponse2, err := cl.NodeListWithResponse(ctx, &openapi.NodeListParams{})
				tests.Ok(t, err, listresponse2)

				ids2 := dt.Map(listresponse2.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })
				a.NotContains(ids2, node5.JSON200.Id, "node5 should not be in the root nodes as it's been moved under node1")

				n1, n1found := lo.Find(listresponse2.JSON200.Nodes, func(c openapi.NodeWithChildren) bool { return c.Id == node1.JSON200.Id })
				r.True(n1found, "node1 must be in the list")
				r.Len(n1.Children, 1, "node1 has one child: node5, which was just added")

				n5 := n1.Children[0]
				a.Equal(node5.JSON200.Id, n5.Id)

				// Remove node5	from node1

				removeresponse, err := cl.NodeRemoveNodeWithResponse(ctx, slug1, slug5, e2e.WithSession(ctx, cj))
				tests.Ok(t, err, removeresponse)

				// Current situation:
				// node1       <- root               has 1 child: node5
				// node2       <- root               has 1 child: node3
				// |- node3    <- child of node1     has 1 child: node4
				//    |- node4 <- child of node3     has no children
				// node5       <- root               has no children

				listresponse3, err := cl.NodeListWithResponse(ctx, &openapi.NodeListParams{})
				tests.Ok(t, err, listresponse3)

				n1, n1found = lo.Find(listresponse3.JSON200.Nodes, func(c openapi.NodeWithChildren) bool { return c.Id == node1.JSON200.Id })
				r.True(n1found, "node1 must be in the list")
				r.Len(n1.Children, 0, "node1 has no children after removing node5")

				n5, n5found := lo.Find(listresponse3.JSON200.Nodes, func(c openapi.NodeWithChildren) bool { return c.Id == node5.JSON200.Id })
				r.True(n5found, "node5 must be in the list")
				r.Len(n5.Children, 0, "node5 has no childen")
			})

			t.Run("move_to_self", func(t *testing.T) {
				a := assert.New(t)
				// r := require.New(t)

				addresponse, err := cl.NodeAddNodeWithResponse(ctx, slug1, slug1, e2e.WithSession(ctx, cj))
				tests.Status(t, err, addresponse, 400)

				listresponse1, err := cl.NodeListWithResponse(ctx, &openapi.NodeListParams{})
				tests.Ok(t, err, listresponse1)

				ids1 := dt.Map(listresponse1.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })
				a.Contains(ids1, node1.JSON200.Id, "node1 should be unaffected")
			})

			t.Run("move_parent_to_child", func(t *testing.T) {
				a := assert.New(t)
				// r := require.New(t)

				nameP1 := "test-node-P1"
				slugP1 := nameP1 + uuid.NewString()
				nodeP1, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Name:       nameP1,
					Slug:       &slugP1,
					Visibility: &visibility,
				}, e2e.WithSession(ctx, cj))
				tests.Ok(t, err, nodeP1)

				nameC1 := "test-node-C1"
				slugC1 := nameC1 + uuid.NewString()
				nodeC1, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Name:       nameC1,
					Slug:       &slugC1,
					Visibility: &visibility,
				}, e2e.WithSession(ctx, cj))
				tests.Ok(t, err, nodeC1)

				// First make C1 a child of P1
				addresponse, err := cl.NodeAddNodeWithResponse(ctx, slugP1, slugC1, e2e.WithSession(ctx, cj))
				tests.Ok(t, err, addresponse)

				// Then move P1 to be a child of C1
				addresponse2, err := cl.NodeAddNodeWithResponse(ctx, slugC1, slugP1, e2e.WithSession(ctx, cj))
				tests.Ok(t, err, addresponse2)

				// C1 must now be a root node without a parent because there
				// cannot be any circular references in the node tree.
				listresponse1, err := cl.NodeListWithResponse(ctx, &openapi.NodeListParams{})
				tests.Ok(t, err, listresponse1)

				ids := dt.Map(listresponse1.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })
				a.Contains(ids, nodeC1.JSON200.Id, "C1 must appear as a root node")
			})
		}))
	}))
}
