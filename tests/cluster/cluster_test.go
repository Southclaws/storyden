package cluster_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Southclaws/dt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/openapi/bindings"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/openapi"
	"github.com/Southclaws/storyden/tests/testutils"
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
			content1 := "# Clusters\n\nRich text content."
			// iurl1 := "https://picsum.photos/200/200"
			url1 := "https://southcla.ws"
			clus1, err := cl.ClusterCreateWithResponse(ctx, openapi.ClusterInitialProps{
				Name:        name1,
				Slug:        slug1,
				Description: "testing clusters api",
				Content:     &content1,
				Url:         &url1,
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(clus1)
			r.Equal(http.StatusOK, clus1.StatusCode())

			a.Equal(name1, clus1.JSON200.Name)
			a.Equal(slug1, clus1.JSON200.Slug)
			a.Equal("testing clusters api", clus1.JSON200.Description)
			a.Equal(content1, *clus1.JSON200.Content)
			a.Equal(url1, clus1.JSON200.Link.Url)
			a.Equal(acc.ID.String(), string(clus1.JSON200.Owner.Id))

			// Get the one just created

			clus1get, err := cl.ClusterGetWithResponse(ctx, slug1)
			r.NoError(err)
			r.NotNil(clus1get)
			r.Equal(http.StatusOK, clus1get.StatusCode())

			a.Equal(name1, clus1get.JSON200.Name)
			a.Equal(slug1, clus1get.JSON200.Slug)
			a.Equal("testing clusters api", clus1get.JSON200.Description)
			a.Equal(acc.ID.String(), string(clus1get.JSON200.Owner.Id))

			// Update the one just created

			name1 = "test-cluster-1-UPDATED"
			slug1 = name1 + uuid.NewString()
			desc1 := "a new description"
			cont1 := "# New content"
			// iurl1 = "https://picsum.photos/500/500"
			url1 = "https://cla.ws"
			prop1 := any(map[string]any{
				"key": "value",
			})
			clus1update, err := cl.ClusterUpdateWithResponse(ctx, clus1.JSON200.Slug, openapi.ClusterMutableProps{
				Name:        &name1,
				Slug:        &slug1,
				Description: &desc1,
				Content:     &cont1,
				Url:         &url1,
				Properties:  &prop1,
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(clus1update)
			r.Equal(http.StatusOK, clus1update.StatusCode())

			a.Equal(name1, clus1update.JSON200.Name)
			a.Equal(slug1, clus1update.JSON200.Slug)
			a.Equal(desc1, clus1update.JSON200.Description)
			a.Equal(cont1, *clus1update.JSON200.Content)
			a.Equal(url1, clus1update.JSON200.Link.Url)
			a.Equal(prop1, clus1update.JSON200.Properties)

			// List all root level clusters

			clist, err := cl.ClusterListWithResponse(ctx, &openapi.ClusterListParams{})
			r.NoError(err)
			r.NotNil(clist)
			r.Equal(http.StatusOK, clist.StatusCode())

			ids := dt.Map(clist.JSON200.Clusters, func(c openapi.Cluster) string { return c.Id })

			a.Contains(ids, clus1.JSON200.Id)

			// Add a child cluster

			name2 := "test-cluster-2"
			slug2 := name2 + uuid.NewString()
			clus2, err := cl.ClusterCreateWithResponse(ctx, openapi.ClusterInitialProps{
				Name:        name2,
				Slug:        slug2,
				Description: "testing clusters children",
				Parent:      &slug1,
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(clus2)
			r.Equal(http.StatusOK, clus2.StatusCode())

			// List all root level clusters

			clist2, err := cl.ClusterListWithResponse(ctx, &openapi.ClusterListParams{})
			r.NoError(err)
			r.NotNil(clist2)
			r.Equal(http.StatusOK, clist2.StatusCode())

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
			r.Equal(http.StatusOK, clus3.StatusCode())

			// This time, use the initial create method instead of using `parent`.

			cadd, err := cl.ClusterAddClusterWithResponse(ctx, slug2, slug3, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(cadd)
			r.Equal(http.StatusOK, cadd.StatusCode())
			r.Equal(clus2.JSON200.Id, cadd.JSON200.Id)

			// List all root level clusters

			clist3, err := cl.ClusterListWithResponse(ctx, &openapi.ClusterListParams{})
			r.NoError(err)
			r.NotNil(clist3)
			r.Equal(http.StatusOK, clist3.StatusCode())

			ids = dt.Map(clist3.JSON200.Clusters, func(c openapi.Cluster) string { return c.Id })
			a.Contains(ids, clus1.JSON200.Id)
			a.NotContains(ids, clus3.JSON200.Id, "must not contain clus3 because it's a child of clus1 and thus not considered root level")

			// List clus2 + descendants

			clist4, err := cl.ClusterListWithResponse(ctx, &openapi.ClusterListParams{
				ClusterId: &clus2.JSON200.Id,
			})
			r.NoError(err)
			r.NotNil(clist4)
			r.Equal(http.StatusOK, clist4.StatusCode())

			ids = dt.Map(clist4.JSON200.Clusters, func(c openapi.Cluster) string { return c.Id })
			a.NotContains(ids, clus1.JSON200.Id, "must not contain clus1 as it's not a descendant of clus2")
			a.Contains(ids, clus2.JSON200.Id)
			a.Contains(ids, clus3.JSON200.Id)

			// Sever clus3 from clus2 so it's root level again

			cremove, err := cl.ClusterRemoveClusterWithResponse(ctx, slug2, slug3, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(cremove)
			r.Equal(http.StatusOK, cremove.StatusCode())
			r.Equal(clus2.JSON200.Id, cremove.JSON200.Id)

			clist5, err := cl.ClusterListWithResponse(ctx, &openapi.ClusterListParams{
				ClusterId: &clus2.JSON200.Id,
			})
			r.NoError(err)
			r.NotNil(clist5)
			r.Equal(http.StatusOK, clist5.StatusCode())

			ids = dt.Map(clist5.JSON200.Clusters, func(c openapi.Cluster) string { return c.Id })
			a.NotContains(ids, clus1.JSON200.Id, "must not contain clus1 as it's not a descendant of clus2")
			a.Contains(ids, clus2.JSON200.Id)
			a.NotContains(ids, clus3.JSON200.Id, "must not contain clus3 as it was severed from clus2")
		}))
	}))
}

func TestClustersFiltering(t *testing.T) {
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

			ctx1, acc1 := e2e.WithAccount(ctx, ar, seed.Account_001_Odin)
			ctx2, acc2 := e2e.WithAccount(ctx, ar, seed.Account_002_Frigg)

			name1 := "test-cluster-owned-by-1"
			slug1 := name1 + uuid.NewString()
			content1 := "# Clusters\n\nOwned by Odin."
			clus1, err := cl.ClusterCreateWithResponse(ctx, openapi.ClusterInitialProps{
				Name:        name1,
				Slug:        slug1,
				Description: "testing clusters api",
				Content:     &content1,
			}, e2e.WithSession(ctx1, cj))
			r.NoError(err)
			r.NotNil(clus1)
			r.Equal(http.StatusOK, clus1.StatusCode())

			name2 := "test-cluster-owned-by-2"
			slug2 := name2 + uuid.NewString()
			content2 := "# Clusters\n\nOwned by Frigg."
			clus2, err := cl.ClusterCreateWithResponse(ctx, openapi.ClusterInitialProps{
				Name:        name2,
				Slug:        slug2,
				Description: "testing clusters api",
				Content:     &content2,
			}, e2e.WithSession(ctx2, cj))
			r.NoError(err)
			r.NotNil(clus1)
			r.Equal(http.StatusOK, clus1.StatusCode())

			clist, err := cl.ClusterListWithResponse(ctx, &openapi.ClusterListParams{
				Author: &acc1.Handle,
			})
			r.NoError(err)
			r.NotNil(clist)
			r.Equal(http.StatusOK, clist.StatusCode())

			ids := dt.Map(clist.JSON200.Clusters, func(c openapi.Cluster) string { return c.Id })

			a.Contains(ids, clus1.JSON200.Id)
			a.NotContains(ids, clus2.JSON200.Id)

			clist2, err := cl.ClusterListWithResponse(ctx, &openapi.ClusterListParams{
				Author: &acc2.Handle,
			})
			r.NoError(err)
			r.NotNil(clist2)
			r.Equal(http.StatusOK, clist2.StatusCode())

			ids2 := dt.Map(clist2.JSON200.Clusters, func(c openapi.Cluster) string { return c.Id })

			a.NotContains(ids2, clus1.JSON200.Id)
			a.Contains(ids2, clus2.JSON200.Id)
		}))
	}))
}

func TestClustersVisibility(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		cj *bindings.CookieJar,
		ar account.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			a := assert.New(t)

			// Tests:
			// - Admin can change visibility of anyone's cluster
			// - Admin can list non-published clusters
			// - Non-admin can not list non-published clusters
			// - Non-admin cannot update visibility of any clusters
			// - Author can list their own hidden clusters
			// - Author can update visibility of their own clusters

			ctxAdmin, _ := e2e.WithAccount(ctx, ar, seed.Account_001_Odin)
			ctxAuthor, accAuthor := e2e.WithAccount(ctx, ar, seed.Account_002_Frigg)
			ctxRando, _ := e2e.WithAccount(ctx, ar, seed.Account_003_Baldur)

			// Author creates 3 clusters

			name1 := "TestClustersFiltering1"
			slug1 := name1 + uuid.NewString()
			clus1 := testutils.AssertRequest(cl.ClusterCreateWithResponse(ctx, openapi.ClusterInitialProps{Name: name1, Slug: slug1, Description: ""}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)

			name2 := "TestClustersFiltering2"
			slug2 := name2 + uuid.NewString()
			clus2 := testutils.AssertRequest(cl.ClusterCreateWithResponse(ctx, openapi.ClusterInitialProps{Name: name2, Slug: slug2, Description: ""}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)

			name3 := "TestClustersFiltering3"
			slug3 := name3 + uuid.NewString()
			clus3 := testutils.AssertRequest(cl.ClusterCreateWithResponse(ctx, openapi.ClusterInitialProps{Name: name3, Slug: slug3, Description: ""}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)

			name4 := "TestClustersFiltering4"
			slug4 := name4 + uuid.NewString()
			clus4 := testutils.AssertRequest(cl.ClusterCreateWithResponse(ctx, openapi.ClusterInitialProps{Name: name4, Slug: slug4, Description: ""}, e2e.WithSession(ctxRando, cj)))(t, http.StatusOK)

			// Public listing without filters does not contain any of them
			// because they were created without being published.

			clist := testutils.AssertRequest(cl.ClusterListWithResponse(ctx, &openapi.ClusterListParams{}))(t, http.StatusOK)

			ids := dt.Map(clist.JSON200.Clusters, func(c openapi.Cluster) string { return c.Id })

			// List does not contain any because they have not been published
			// and the request was made without auth from the owner.
			a.NotContains(ids, clus1.JSON200.Id)
			a.NotContains(ids, clus2.JSON200.Id)
			a.NotContains(ids, clus3.JSON200.Id)
			a.NotContains(ids, clus4.JSON200.Id)

			clist = testutils.AssertRequest(cl.ClusterListWithResponse(ctx, &openapi.ClusterListParams{
				Author: &accAuthor.Handle,
			}))(t, http.StatusOK)

			ids = dt.Map(clist.JSON200.Clusters, func(c openapi.Cluster) string { return c.Id })

			a.NotContains(ids, clus1.JSON200.Id)
			a.NotContains(ids, clus2.JSON200.Id)
			a.NotContains(ids, clus3.JSON200.Id)
			a.NotContains(ids, clus4.JSON200.Id)

			// Admin can change visibility

			update1 := testutils.AssertRequest(
				cl.ClusterUpdateVisibilityWithResponse(ctx, clus1.JSON200.Slug, openapi.VisibilityMutationProps{
					Visibility: openapi.Published,
				}, e2e.WithSession(ctxAdmin, cj)),
			)(t, http.StatusOK)
			a.Equal(openapi.Published, update1.JSON200.Visibility)

			clist = testutils.AssertRequest(cl.ClusterListWithResponse(ctx, &openapi.ClusterListParams{
				Author: &accAuthor.Handle,
			}))(t, http.StatusOK)

			ids = dt.Map(clist.JSON200.Clusters, func(c openapi.Cluster) string { return c.Id })

			a.Contains(ids, clus1.JSON200.Id, "admin made this cluster visible")
			a.NotContains(ids, clus2.JSON200.Id)
			a.NotContains(ids, clus3.JSON200.Id)
			a.NotContains(ids, clus4.JSON200.Id)

			// Author can change visibility

			update2 := testutils.AssertRequest(
				cl.ClusterUpdateVisibilityWithResponse(ctx, clus2.JSON200.Slug, openapi.VisibilityMutationProps{
					Visibility: openapi.Published,
				}, e2e.WithSession(ctxAuthor, cj)),
			)(t, http.StatusOK)
			a.Equal(openapi.Published, update2.JSON200.Visibility)

			clist = testutils.AssertRequest(cl.ClusterListWithResponse(ctx, &openapi.ClusterListParams{
				Author: &accAuthor.Handle,
			}))(t, http.StatusOK)

			ids = dt.Map(clist.JSON200.Clusters, func(c openapi.Cluster) string { return c.Id })

			a.Contains(ids, clus1.JSON200.Id, "admin made this cluster visible")
			a.Contains(ids, clus2.JSON200.Id)
			a.NotContains(ids, clus3.JSON200.Id)
			a.NotContains(ids, clus4.JSON200.Id)

			// Author can list their own hidden clusters, but not others.

			clist = testutils.AssertRequest(cl.ClusterListWithResponse(ctx, &openapi.ClusterListParams{
				Visibility: &[]openapi.Visibility{openapi.Draft},
			}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)

			ids = dt.Map(clist.JSON200.Clusters, func(c openapi.Cluster) string { return c.Id })

			a.NotContains(ids, clus1.JSON200.Id)
			a.NotContains(ids, clus2.JSON200.Id)
			a.Contains(ids, clus3.JSON200.Id, "this is the only cluster not published above")
			a.NotContains(ids, clus4.JSON200.Id, "owned by someone else, should not be visible")

			// Admin can only list in-review clusters, but not drafts.

			clist = testutils.AssertRequest(cl.ClusterListWithResponse(ctx, &openapi.ClusterListParams{
				Visibility: &[]openapi.Visibility{openapi.Review},
			}, e2e.WithSession(ctxAdmin, cj)))(t, http.StatusOK)

			ids = dt.Map(clist.JSON200.Clusters, func(c openapi.Cluster) string { return c.Id })

			a.NotContains(ids, clus1.JSON200.Id)
			a.NotContains(ids, clus2.JSON200.Id)
			a.NotContains(ids, clus3.JSON200.Id)
			a.NotContains(ids, clus4.JSON200.Id)

			// Author moves clus3 to in-review

			update3 := testutils.AssertRequest(
				cl.ClusterUpdateVisibilityWithResponse(ctx, clus3.JSON200.Slug, openapi.VisibilityMutationProps{
					Visibility: openapi.Review,
				}, e2e.WithSession(ctxAuthor, cj)),
			)(t, http.StatusOK)
			a.Equal(openapi.Review, update3.JSON200.Visibility)

			clist = testutils.AssertRequest(cl.ClusterListWithResponse(ctx, &openapi.ClusterListParams{
				Visibility: &[]openapi.Visibility{openapi.Review},
			}, e2e.WithSession(ctxAdmin, cj)))(t, http.StatusOK)

			ids = dt.Map(clist.JSON200.Clusters, func(c openapi.Cluster) string { return c.Id })

			a.NotContains(ids, clus1.JSON200.Id)
			a.NotContains(ids, clus2.JSON200.Id)
			a.Contains(ids, clus3.JSON200.Id, "in review so is now visible to admins")
			a.NotContains(ids, clus4.JSON200.Id, "")
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
			a.Equal(http.StatusNotFound, get404.StatusCode())

			update403, err := cl.ClusterUpdateWithResponse(ctx, "nonexistent", openapi.ClusterMutableProps{})
			r.NoError(err)
			r.NotNil(update403)
			a.Equal(http.StatusForbidden, update403.StatusCode())

			update404, err := cl.ClusterUpdateWithResponse(ctx, "nonexistent", openapi.ClusterMutableProps{}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(update404)
			a.Equal(http.StatusNotFound, update404.StatusCode())
		}))
	}))
}
