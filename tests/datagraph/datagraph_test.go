package datagraph_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

func TestDatagraphHappyPath(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		cj *session.Jar,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			ctx, acc := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)

			// iurl := "https://picsum.photos/500/500"

			name1 := "test-node-1"
			slug1 := name1 + uuid.NewString()
			node1, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name: name1,
				Slug: &slug1,
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(node1)
			r.Equal(http.StatusOK, node1.StatusCode())

			a.Equal(name1, node1.JSON200.Name)
			a.Equal(slug1, node1.JSON200.Slug)
			a.Equal("", node1.JSON200.Description)
			a.Equal(acc.ID.String(), string(node1.JSON200.Owner.Id))

			// Add a child node

			name2 := "test-node-2"
			slug2 := name2 + uuid.NewString()
			node2, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name: name2,
				Slug: &slug2,
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(node2)
			r.Equal(http.StatusOK, node2.StatusCode())

			cadd, err := cl.NodeAddNodeWithResponse(ctx, slug1, slug2, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(cadd)
			r.Equal(http.StatusOK, cadd.StatusCode())
			r.Equal(node1.JSON200.Id, cadd.JSON200.Id)

			// Add another child to this child
			// node1
			// |- node2
			//    |- node3

			name3 := "test-node-3"
			slug3 := name3 + uuid.NewString()
			node3, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name: name3,
				Slug: &slug3,
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(node3)
			r.Equal(http.StatusOK, node3.StatusCode())

			cadd, err = cl.NodeAddNodeWithResponse(ctx, slug2, slug3, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(cadd)
			r.Equal(http.StatusOK, cadd.StatusCode())
			r.Equal(node2.JSON200.Id, cadd.JSON200.Id)
		}))
	}))
}

func TestDatagraphDeletions(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		cj *session.Jar,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			ctx, _ := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)

			// Create three nodes in a tree
			// node1
			// |- node2
			//    |- node3

			node1, err := cl.NodeCreateWithResponse(ctx, uniqueNode("deletions1"), e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.Equal(http.StatusOK, node1.StatusCode())

			node2, err := cl.NodeCreateWithResponse(ctx, uniqueNode("deletions2"), e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.Equal(http.StatusOK, node2.StatusCode())

			node3, err := cl.NodeCreateWithResponse(ctx, uniqueNode("deletions3"), e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.Equal(http.StatusOK, node3.StatusCode())

			cadd, err := cl.NodeAddNode(ctx, node1.JSON200.Slug, node2.JSON200.Slug, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.Equal(http.StatusOK, cadd.StatusCode)

			cadd, err = cl.NodeAddNode(ctx, node2.JSON200.Slug, node3.JSON200.Slug, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.Equal(http.StatusOK, cadd.StatusCode)

			cdel, err := cl.NodeDeleteWithResponse(ctx, node3.JSON200.Slug, nil, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.Equal(http.StatusOK, cdel.StatusCode())
			a.Nil(cdel.JSON200.Destination)

			node2node, err := cl.NodeCreateWithResponse(ctx, uniqueNode("deletions2child"), e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.Equal(http.StatusOK, node2node.StatusCode())

			cadd, err = cl.NodeAddNode(ctx, node2.JSON200.Slug, node2node.JSON200.Slug, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.Equal(http.StatusOK, cadd.StatusCode)

			cdel, err = cl.NodeDeleteWithResponse(ctx, node2.JSON200.Slug, &openapi.NodeDeleteParams{
				TargetNode: &node1.JSON200.Slug,
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.Equal(http.StatusOK, cdel.StatusCode())
			a.NotNil(cdel.JSON200.Destination)
			a.Equal(node1.JSON200.Id, cdel.JSON200.Destination.Id)

			node1get, err := cl.NodeGetWithResponse(ctx, node1.JSON200.Slug)
			r.NoError(err)
			r.Equal(http.StatusOK, node1get.StatusCode())

			a.Len(node1get.JSON200.Children, 1)
		}))
	}))
}

func uniqueNode(name string) openapi.NodeInitialProps {
	slug := name + uuid.NewString()
	return openapi.NodeInitialProps{
		Name: name,
		Slug: &slug,
	}
}
