package cluster_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/openapi/bindings"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/openapi"
)

func TestClustersHappyPath(t *testing.T) {
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

			name1 := "test-cluster-1"
			slug1 := name1 + uuid.NewString()
			clus1, err := cl.ClusterCreateWithResponse(ctx, openapi.ClusterInitialProps{
				Name:        name1,
				Slug:        slug1,
				Description: "testing clusters api",
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(clus1)
			r.Equal(200, clus1.StatusCode())

			a.Equal(name1, clus1.JSON200.Name)
			a.Equal(slug1, clus1.JSON200.Slug)
			a.Equal("testing clusters api", clus1.JSON200.Description)
			a.Equal(acc.ID.String(), string(clus1.JSON200.Owner.Id))

			// Get the one just created

			clus1get, err := cl.ClusterGetWithResponse(ctx, slug1)
			r.NoError(err)
			r.NotNil(clus1get)
			r.Equal(200, clus1get.StatusCode())

			a.Equal(name1, clus1get.JSON200.Name)
			a.Equal(slug1, clus1get.JSON200.Slug)
			a.Equal("testing clusters api", clus1get.JSON200.Description)
			a.Equal(acc.ID.String(), string(clus1get.JSON200.Owner.Id))

			// Update the one just created

			name1 = "test-cluster-1-UPDATED"
			slug1 = name1 + uuid.NewString()
			desc1 := "a new description"
			iurl1 := "https://niceme.me"
			prop1 := any(map[string]any{
				"key": "value",
			})
			clus1update, err := cl.ClusterUpdateWithResponse(ctx, clus1.JSON200.Slug, openapi.ClusterMutableProps{
				Name:        &name1,
				Slug:        &slug1,
				Description: &desc1,
				ImageUrl:    &iurl1,
				Properties:  &prop1,
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(clus1update)
			r.Equal(200, clus1update.StatusCode())

			a.Equal(name1, clus1update.JSON200.Name)
			a.Equal(slug1, clus1update.JSON200.Slug)
			a.Equal(desc1, clus1update.JSON200.Description)
			a.Equal(iurl1, *clus1update.JSON200.ImageUrl)
			a.Equal(prop1, clus1update.JSON200.Properties)

			// List all root level clusters

			clist, err := cl.ClusterListWithResponse(ctx, &openapi.ClusterListParams{})
			r.NoError(err)
			r.NotNil(clist)
			r.Equal(200, clist.StatusCode())

			ids := dt.Map(clist.JSON200.Clusters, func(c openapi.Cluster) string { return c.Id })

			a.Contains(ids, clus1.JSON200.Id)

			// Add a child cluster

			name2 := "test-cluster-2"
			slug2 := name2 + uuid.NewString()
			clus2, err := cl.ClusterCreateWithResponse(ctx, openapi.ClusterInitialProps{
				Name:        name2,
				Slug:        slug2,
				Description: "testing clusters children",
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(clus2)
			r.Equal(200, clus2.StatusCode())

			cadd, err := cl.ClusterAddClusterWithResponse(ctx, slug1, slug2, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(cadd)
			r.Equal(200, cadd.StatusCode())
			r.Equal(clus1.JSON200.Id, cadd.JSON200.Id)

			// List all root level clusters

			clist2, err := cl.ClusterListWithResponse(ctx, &openapi.ClusterListParams{})
			r.NoError(err)
			r.NotNil(clist2)
			r.Equal(200, clist2.StatusCode())

			ids = dt.Map(clist2.JSON200.Clusters, func(c openapi.Cluster) string { return c.Id })
			a.Contains(ids, clus1.JSON200.Id)
			a.NotContains(ids, clus2.JSON200.Id, "must not contain clus2 because it's a child of clus1 and thus not considered root level")

			// Add another child to this child
			// clus1
			// |- clus2
			//    |- clus3
			// then query children of clus2, expect clus2+clus3 only

			name3 := "test-cluster-3"
			slug3 := name3 + uuid.NewString()
			clus3, err := cl.ClusterCreateWithResponse(ctx, openapi.ClusterInitialProps{
				Name:        name3,
				Slug:        slug3,
				Description: "testing clusters children",
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(clus3)
			r.Equal(200, clus3.StatusCode())

			cadd, err = cl.ClusterAddClusterWithResponse(ctx, slug2, slug3, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(cadd)
			r.Equal(200, cadd.StatusCode())
			r.Equal(clus2.JSON200.Id, cadd.JSON200.Id)

			// List all root level clusters

			clist3, err := cl.ClusterListWithResponse(ctx, &openapi.ClusterListParams{})
			r.NoError(err)
			r.NotNil(clist3)
			r.Equal(200, clist3.StatusCode())

			ids = dt.Map(clist3.JSON200.Clusters, func(c openapi.Cluster) string { return c.Id })
			a.Contains(ids, clus1.JSON200.Id)
			a.NotContains(ids, clus3.JSON200.Id, "must not contain clus3 because it's a child of clus1 and thus not considered root level")

			// List clus2 + descendants

			clist4, err := cl.ClusterListWithResponse(ctx, &openapi.ClusterListParams{
				ClusterId: &clus2.JSON200.Id,
			})
			r.NoError(err)
			r.NotNil(clist4)
			r.Equal(200, clist4.StatusCode())

			ids = dt.Map(clist4.JSON200.Clusters, func(c openapi.Cluster) string { return c.Id })
			a.NotContains(ids, clus1.JSON200.Id, "must not contain clus1 as it's not a descendant of clus2")
			a.Contains(ids, clus2.JSON200.Id)
			a.Contains(ids, clus3.JSON200.Id)
		}))
	}))
}

func TestClustersErrors(t *testing.T) {
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

			get404, err := cl.ClusterGetWithResponse(ctx, "nonexistent")
			r.NoError(err)
			r.NotNil(get404)
			a.Equal(404, get404.StatusCode())

			update401, err := cl.ClusterUpdateWithResponse(ctx, "nonexistent", openapi.ClusterMutableProps{})
			r.NoError(err)
			r.NotNil(update401)
			a.Equal(401, update401.StatusCode())

			update404, err := cl.ClusterUpdateWithResponse(ctx, "nonexistent", openapi.ClusterMutableProps{}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(update404)
			a.Equal(404, update404.StatusCode())
		}))
	}))
}
