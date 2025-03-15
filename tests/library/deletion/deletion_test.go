package deletion_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
	"github.com/Southclaws/storyden/tests/library"
)

func TestLibraryNodeDeletion(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		cj *session_cookie.Jar,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			ctx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			session := e2e.WithSession(ctx, cj)

			// Each test creates three nodes in a tree
			//   node1
			//   |- node2
			//      |- node3
			// and then deletes a child or parent node to test tree logic.

			t.Run("delete_leaf_node", func(t *testing.T) {
				a := assert.New(t)

				node1, err := cl.NodeCreateWithResponse(ctx, library.UniqueNode("node1"), session)
				tests.Ok(t, err, node1)
				node2, err := cl.NodeCreateWithResponse(ctx, library.UniqueNode("node2"), session)
				tests.Ok(t, err, node2)
				node3, err := cl.NodeCreateWithResponse(ctx, library.UniqueNode("node3-deleteme"), session)
				tests.Ok(t, err, node3)
				cadd, err := cl.NodeAddNodeWithResponse(ctx, node1.JSON200.Slug, node2.JSON200.Slug, session)
				tests.Ok(t, err, cadd)
				cadd, err = cl.NodeAddNodeWithResponse(ctx, node2.JSON200.Slug, node3.JSON200.Slug, session)
				tests.Ok(t, err, cadd)

				cdel, err := cl.NodeDeleteWithResponse(ctx, node3.JSON200.Slug, nil, session)
				tests.Ok(t, err, cdel)
				a.Nil(cdel.JSON200.Destination)

				node3get, err := cl.NodeGetWithResponse(ctx, node3.JSON200.Slug, &openapi.NodeGetParams{}, session)
				tests.Status(t, err, node3get, http.StatusNotFound)

				node2get, err := cl.NodeGetWithResponse(ctx, node2.JSON200.Slug, &openapi.NodeGetParams{}, session)
				tests.Ok(t, err, node2get)
				a.Len(node2get.JSON200.Children, 0)
			})

			t.Run("delete_parent_node_with_target", func(t *testing.T) {
				a := assert.New(t)

				node1, err := cl.NodeCreateWithResponse(ctx, library.UniqueNode("node1"), session)
				tests.Ok(t, err, node1)
				node2, err := cl.NodeCreateWithResponse(ctx, library.UniqueNode("node2-deleteme"), session)
				tests.Ok(t, err, node2)
				node3, err := cl.NodeCreateWithResponse(ctx, library.UniqueNode("node3"), session)
				tests.Ok(t, err, node3)
				cadd, err := cl.NodeAddNodeWithResponse(ctx, node1.JSON200.Slug, node2.JSON200.Slug, session)
				tests.Ok(t, err, cadd)
				cadd, err = cl.NodeAddNodeWithResponse(ctx, node2.JSON200.Slug, node3.JSON200.Slug, session)
				tests.Ok(t, err, cadd)

				cdel, err := cl.NodeDeleteWithResponse(ctx, node2.JSON200.Slug, &openapi.NodeDeleteParams{
					TargetNode: &node1.JSON200.Slug,
				}, session)
				tests.Ok(t, err, cdel)
				a.NotNil(cdel.JSON200.Destination)
				a.Equal(node1.JSON200.Id, cdel.JSON200.Destination.Id)

				node1get, err := cl.NodeGetWithResponse(ctx, node1.JSON200.Slug, &openapi.NodeGetParams{}, session)
				tests.Ok(t, err, node1get)
				a.Len(node1get.JSON200.Children, 1, "node3 is moved under node1 when node2 is deleted")
			})

			t.Run("delete_parent_node_without_target", func(t *testing.T) {
				a := assert.New(t)

				node1, err := cl.NodeCreateWithResponse(ctx, library.UniqueNode("node1"), session)
				tests.Ok(t, err, node1)
				node2, err := cl.NodeCreateWithResponse(ctx, library.UniqueNode("node2-deleteme"), session)
				tests.Ok(t, err, node2)
				node3, err := cl.NodeCreateWithResponse(ctx, library.UniqueNode("node3"), session)
				tests.Ok(t, err, node3)
				cadd, err := cl.NodeAddNodeWithResponse(ctx, node1.JSON200.Slug, node2.JSON200.Slug, session)
				tests.Ok(t, err, cadd)
				cadd, err = cl.NodeAddNodeWithResponse(ctx, node2.JSON200.Slug, node3.JSON200.Slug, session)
				tests.Ok(t, err, cadd)

				ndelete, err := cl.NodeDeleteWithResponse(ctx, node2.JSON200.Slug, &openapi.NodeDeleteParams{}, session)
				tests.Ok(t, err, ndelete)
				a.Nil(ndelete.JSON200.Destination)

				node1get, err := cl.NodeGetWithResponse(ctx, node1.JSON200.Slug, &openapi.NodeGetParams{}, session)
				tests.Ok(t, err, node1get)
				a.Len(node1get.JSON200.Children, 0)

				// node3 is moved to root, with no parent
				node3get, err := cl.NodeGetWithResponse(ctx, node3.JSON200.Slug, &openapi.NodeGetParams{}, session)
				tests.Ok(t, err, node3get)
				a.Nil(node3get.JSON200.Parent)
			})
		}))
	}))
}
