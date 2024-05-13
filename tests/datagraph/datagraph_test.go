package datagraph_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/openapi"
	"github.com/Southclaws/storyden/app/transports/openapi/bindings"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

func TestDatagraphHappyPath(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		cj *bindings.CookieJar,
		ar account.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			ctx, acc := e2e.WithAccount(ctx, ar, seed.Account_001_Odin)

			// iurl := "https://picsum.photos/500/500"

			name1 := "test-node-1"
			slug1 := name1 + uuid.NewString()
			clus1, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:        name1,
				Slug:        slug1,
				Description: "testing nodes api",
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(clus1)
			r.Equal(http.StatusOK, clus1.StatusCode())

			a.Equal(name1, clus1.JSON200.Name)
			a.Equal(slug1, clus1.JSON200.Slug)
			a.Equal("testing nodes api", clus1.JSON200.Description)
			a.Equal(acc.ID.String(), string(clus1.JSON200.Owner.Id))

			// Add a child node

			name2 := "test-node-2"
			slug2 := name2 + uuid.NewString()
			clus2, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:        name2,
				Slug:        slug2,
				Description: "testing nodes children",
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(clus2)
			r.Equal(http.StatusOK, clus2.StatusCode())

			cadd, err := cl.NodeAddNodeWithResponse(ctx, slug1, slug2, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(cadd)
			r.Equal(http.StatusOK, cadd.StatusCode())
			r.Equal(clus1.JSON200.Id, cadd.JSON200.Id)

			// Add another child to this child
			// clus1
			// |- clus2
			//    |- clus3

			name3 := "test-node-3"
			slug3 := name3 + uuid.NewString()
			clus3, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:        name3,
				Slug:        slug3,
				Description: "testing nodes children",
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(clus3)
			r.Equal(http.StatusOK, clus3.StatusCode())

			cadd, err = cl.NodeAddNodeWithResponse(ctx, slug2, slug3, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(cadd)
			r.Equal(http.StatusOK, cadd.StatusCode())
			r.Equal(clus2.JSON200.Id, cadd.JSON200.Id)
		}))
	}))
}

func TestDatagraphDeletions(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		cj *bindings.CookieJar,
		ar account.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			ctx, _ := e2e.WithAccount(ctx, ar, seed.Account_001_Odin)

			// Create three nodes in a tree
			// clus1
			// |- clus2
			//    |- clus3

			clus1, err := cl.NodeCreateWithResponse(ctx, uniqueNode("deletions1"), e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.Equal(http.StatusOK, clus1.StatusCode())

			clus2, err := cl.NodeCreateWithResponse(ctx, uniqueNode("deletions2"), e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.Equal(http.StatusOK, clus2.StatusCode())

			clus3, err := cl.NodeCreateWithResponse(ctx, uniqueNode("deletions3"), e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.Equal(http.StatusOK, clus3.StatusCode())

			cadd, err := cl.NodeAddNode(ctx, clus1.JSON200.Slug, clus2.JSON200.Slug, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.Equal(http.StatusOK, cadd.StatusCode)

			cadd, err = cl.NodeAddNode(ctx, clus2.JSON200.Slug, clus3.JSON200.Slug, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.Equal(http.StatusOK, cadd.StatusCode)

			cdel, err := cl.NodeDeleteWithResponse(ctx, clus3.JSON200.Slug, nil, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.Equal(http.StatusOK, cdel.StatusCode())
			a.Nil(cdel.JSON200.Destination)

			clus2clus, err := cl.NodeCreateWithResponse(ctx, uniqueNode("deletions2child"), e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.Equal(http.StatusOK, clus2clus.StatusCode())

			cadd, err = cl.NodeAddNode(ctx, clus2.JSON200.Slug, clus2clus.JSON200.Slug, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.Equal(http.StatusOK, cadd.StatusCode)

			cdel, err = cl.NodeDeleteWithResponse(ctx, clus2.JSON200.Slug, &openapi.NodeDeleteParams{
				TargetNode:     &clus1.JSON200.Slug,
				MoveChildNodes: opt.New(true).Ptr(),
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.Equal(http.StatusOK, cdel.StatusCode())
			a.NotNil(cdel.JSON200.Destination)
			a.Equal(clus1.JSON200.Id, cdel.JSON200.Destination.Id)

			clus1get, err := cl.NodeGetWithResponse(ctx, clus1.JSON200.Slug)
			r.NoError(err)
			r.Equal(http.StatusOK, clus1get.StatusCode())

			a.Len(clus1get.JSON200.Children, 1)
		}))
	}))
}

func uniqueNode(name string) openapi.NodeInitialProps {
	return openapi.NodeInitialProps{
		Name:        name,
		Slug:        name + uuid.NewString(),
		Description: name,
	}
}
