package nodeter_test

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
	"github.com/Southclaws/storyden/app/transports/openapi"
	"github.com/Southclaws/storyden/app/transports/openapi/bindings"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
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

			visibility := openapi.Published

			name1 := "test-nodeter-1"
			slug1 := name1 + uuid.NewString()
			content1 := "# Clusters\n\nRich text content."
			// iurl1 := "https://picsum.photos/200/200"
			url1 := "https://southcla.ws"
			node1, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:        name1,
				Slug:        slug1,
				Description: "testing nodeters api",
				Content:     &content1,
				Url:         &url1,
				Visibility:  &visibility, // Admin account can post directly to published
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(node1)
			r.Equal(http.StatusOK, node1.StatusCode())

			a.Equal(name1, node1.JSON200.Name)
			a.Equal(slug1, node1.JSON200.Slug)
			a.Equal("testing nodeters api", node1.JSON200.Description)
			a.Equal(content1, *node1.JSON200.Content)
			a.Equal(url1, node1.JSON200.Link.Url)
			a.Equal(acc.ID.String(), string(node1.JSON200.Owner.Id))

			// Get the one just created

			node1get, err := cl.NodeGetWithResponse(ctx, slug1)
			r.NoError(err)
			r.NotNil(node1get)
			r.Equal(http.StatusOK, node1get.StatusCode())

			a.Equal(name1, node1get.JSON200.Name)
			a.Equal(slug1, node1get.JSON200.Slug)
			a.Equal("testing nodeters api", node1get.JSON200.Description)
			a.Equal(acc.ID.String(), string(node1get.JSON200.Owner.Id))

			// Update the one just created

			name1 = "test-nodeter-1-UPDATED"
			slug1 = name1 + uuid.NewString()
			desc1 := "a new description"
			cont1 := "# New content"
			// iurl1 = "https://picsum.photos/500/500"
			url1 = "https://cla.ws"
			prop1 := any(map[string]any{
				"key": "value",
			})
			node1update, err := cl.NodeUpdateWithResponse(ctx, node1.JSON200.Slug, openapi.NodeMutableProps{
				Name:        &name1,
				Slug:        &slug1,
				Description: &desc1,
				Content:     &cont1,
				Url:         &url1,
				Properties:  &prop1,
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(node1update)
			r.Equal(http.StatusOK, node1update.StatusCode())

			a.Equal(name1, node1update.JSON200.Name)
			a.Equal(slug1, node1update.JSON200.Slug)
			a.Equal(desc1, node1update.JSON200.Description)
			a.Equal(cont1, *node1update.JSON200.Content)
			a.Equal(url1, node1update.JSON200.Link.Url)
			a.Equal(prop1, node1update.JSON200.Properties)

			// List all root level nodeters

			clist, err := cl.NodeListWithResponse(ctx, &openapi.NodeListParams{})
			r.NoError(err)
			r.NotNil(clist)
			r.Equal(http.StatusOK, clist.StatusCode())

			ids := dt.Map(clist.JSON200.Nodes, func(c openapi.Node) string { return c.Id })

			a.Contains(ids, node1.JSON200.Id)

			// Add a child nodeter

			name2 := "test-nodeter-2"
			slug2 := name2 + uuid.NewString()
			node2, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:        name2,
				Slug:        slug2,
				Description: "testing nodeters children",
				Parent:      &slug1,
				Visibility:  &visibility,
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(node2)
			r.Equal(http.StatusOK, node2.StatusCode())

			// List all root level nodeters

			clist2, err := cl.NodeListWithResponse(ctx, &openapi.NodeListParams{})
			r.NoError(err)
			r.NotNil(clist2)
			r.Equal(http.StatusOK, clist2.StatusCode())

			ids = dt.Map(clist2.JSON200.Nodes, func(c openapi.Node) string { return c.Id })
			a.Contains(ids, node1.JSON200.Id)
			a.NotContains(ids, node2.JSON200.Id, "must not contain node2 because it's a child of node1 and thus not considered root level")

			// Add another child to this child
			// node1
			// |- node2
			//    |- node3
			// then query children of node2, expect node2+node3 only

			name3 := "test-nodeter-3"
			slug3 := name3 + uuid.NewString()
			node3, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:        name3,
				Slug:        slug3,
				Description: "testing nodeters children",
				Visibility:  &visibility,
			}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(node3)
			r.Equal(http.StatusOK, node3.StatusCode())

			// This time, use the initial create method instead of using `parent`.

			cadd, err := cl.NodeAddNodeWithResponse(ctx, slug2, slug3, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(cadd)
			r.Equal(http.StatusOK, cadd.StatusCode())
			r.Equal(node2.JSON200.Id, cadd.JSON200.Id)

			// List all root level nodeters

			clist3, err := cl.NodeListWithResponse(ctx, &openapi.NodeListParams{})
			r.NoError(err)
			r.NotNil(clist3)
			r.Equal(http.StatusOK, clist3.StatusCode())

			ids = dt.Map(clist3.JSON200.Nodes, func(c openapi.Node) string { return c.Id })
			a.Contains(ids, node1.JSON200.Id)
			a.NotContains(ids, node3.JSON200.Id, "must not contain node3 because it's a child of node1 and thus not considered root level")

			// List node2 + descendants

			clist4, err := cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
				NodeId: &node2.JSON200.Id,
			})
			r.NoError(err)
			r.NotNil(clist4)
			r.Equal(http.StatusOK, clist4.StatusCode())

			ids = dt.Map(clist4.JSON200.Nodes, func(c openapi.Node) string { return c.Id })
			a.NotContains(ids, node1.JSON200.Id, "must not contain node1 as it's not a descendant of node2")
			a.Contains(ids, node2.JSON200.Id)
			a.Contains(ids, node3.JSON200.Id)

			// Sever node3 from node2 so it's root level again

			cremove, err := cl.NodeRemoveNodeWithResponse(ctx, slug2, slug3, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(cremove)
			r.Equal(http.StatusOK, cremove.StatusCode())
			r.Equal(node2.JSON200.Id, cremove.JSON200.Id)

			clist5, err := cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
				NodeId: &node2.JSON200.Id,
			})
			r.NoError(err)
			r.NotNil(clist5)
			r.Equal(http.StatusOK, clist5.StatusCode())

			ids = dt.Map(clist5.JSON200.Nodes, func(c openapi.Node) string { return c.Id })
			a.NotContains(ids, node1.JSON200.Id, "must not contain node1 as it's not a descendant of node2")
			a.Contains(ids, node2.JSON200.Id)
			a.NotContains(ids, node3.JSON200.Id, "must not contain node3 as it was severed from node2")
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

			visibility := openapi.Published

			name1 := "test-nodeter-owned-by-1"
			slug1 := name1 + uuid.NewString()
			content1 := "# Clusters\n\nOwned by Odin."
			node1, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:        name1,
				Slug:        slug1,
				Description: "testing nodeters api",
				Content:     &content1,
				Visibility:  &visibility,
			}, e2e.WithSession(ctx1, cj))
			r.NoError(err)
			r.NotNil(node1)
			r.Equal(http.StatusOK, node1.StatusCode())

			name2 := "test-nodeter-owned-by-2"
			slug2 := name2 + uuid.NewString()
			content2 := "# Clusters\n\nOwned by Frigg."
			node2, err := cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
				Name:        name2,
				Slug:        slug2,
				Description: "testing nodeters api",
				Content:     &content2,
				Visibility:  &visibility,
			}, e2e.WithSession(ctx2, cj))
			r.NoError(err)
			r.NotNil(node1)
			r.Equal(http.StatusOK, node1.StatusCode())

			clist, err := cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
				Author: &acc1.Handle,
			})
			r.NoError(err)
			r.NotNil(clist)
			r.Equal(http.StatusOK, clist.StatusCode())

			ids := dt.Map(clist.JSON200.Nodes, func(c openapi.Node) string { return c.Id })

			a.Contains(ids, node1.JSON200.Id)
			a.NotContains(ids, node2.JSON200.Id)

			clist2, err := cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
				Author: &acc2.Handle,
			})
			r.NoError(err)
			r.NotNil(clist2)
			r.Equal(http.StatusOK, clist2.StatusCode())

			ids2 := dt.Map(clist2.JSON200.Nodes, func(c openapi.Node) string { return c.Id })

			a.NotContains(ids2, node1.JSON200.Id)
			a.Contains(ids2, node2.JSON200.Id)
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
			// - Admin can change visibility of anyone's nodeter
			// - Admin can list non-published nodeters
			// - Non-admin can not list non-published nodeters
			// - Non-admin cannot update visibility of any nodeters
			// - Author can list their own hidden nodeters
			// - Author can update visibility of their own nodeters

			ctxAdmin, _ := e2e.WithAccount(ctx, ar, seed.Account_001_Odin)
			ctxAuthor, accAuthor := e2e.WithAccount(ctx, ar, seed.Account_002_Frigg)
			ctxRando, _ := e2e.WithAccount(ctx, ar, seed.Account_003_Baldur)

			// Author creates 3 nodeters

			name1 := "TestClustersFiltering1"
			slug1 := name1 + uuid.NewString()
			node1 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: name1, Slug: slug1, Description: ""}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)

			name2 := "TestClustersFiltering2"
			slug2 := name2 + uuid.NewString()
			node2 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: name2, Slug: slug2, Description: ""}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)

			name3 := "TestClustersFiltering3"
			slug3 := name3 + uuid.NewString()
			node3 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: name3, Slug: slug3, Description: ""}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)

			name4 := "TestClustersFiltering4"
			slug4 := name4 + uuid.NewString()
			node4 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: name4, Slug: slug4, Description: ""}, e2e.WithSession(ctxRando, cj)))(t, http.StatusOK)

			// Public listing without filters does not contain any of them
			// because they were created without being published.

			clist := testutils.AssertRequest(cl.NodeListWithResponse(ctx, &openapi.NodeListParams{}))(t, http.StatusOK)

			ids := dt.Map(clist.JSON200.Nodes, func(c openapi.Node) string { return c.Id })

			// List does not contain any because they have not been published
			// and the request was made without auth from the owner.
			a.NotContains(ids, node1.JSON200.Id)
			a.NotContains(ids, node2.JSON200.Id)
			a.NotContains(ids, node3.JSON200.Id)
			a.NotContains(ids, node4.JSON200.Id)

			clist = testutils.AssertRequest(cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
				Author: &accAuthor.Handle,
			}))(t, http.StatusOK)

			ids = dt.Map(clist.JSON200.Nodes, func(c openapi.Node) string { return c.Id })

			a.NotContains(ids, node1.JSON200.Id)
			a.NotContains(ids, node2.JSON200.Id)
			a.NotContains(ids, node3.JSON200.Id)
			a.NotContains(ids, node4.JSON200.Id)

			// Admin can change visibility

			update1 := testutils.AssertRequest(
				cl.NodeUpdateVisibilityWithResponse(ctx, node1.JSON200.Slug, openapi.VisibilityMutationProps{
					Visibility: openapi.Published,
				}, e2e.WithSession(ctxAdmin, cj)),
			)(t, http.StatusOK)
			a.Equal(openapi.Published, update1.JSON200.Visibility)

			clist = testutils.AssertRequest(cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
				Author: &accAuthor.Handle,
			}))(t, http.StatusOK)

			ids = dt.Map(clist.JSON200.Nodes, func(c openapi.Node) string { return c.Id })

			a.Contains(ids, node1.JSON200.Id, "admin made this nodeter visible")
			a.NotContains(ids, node2.JSON200.Id)
			a.NotContains(ids, node3.JSON200.Id)
			a.NotContains(ids, node4.JSON200.Id)

			// Author can change visibility

			update2 := testutils.AssertRequest(
				cl.NodeUpdateVisibilityWithResponse(ctx, node2.JSON200.Slug, openapi.VisibilityMutationProps{
					Visibility: openapi.Published,
				}, e2e.WithSession(ctxAuthor, cj)),
			)(t, http.StatusOK)
			a.Equal(openapi.Published, update2.JSON200.Visibility)

			clist = testutils.AssertRequest(cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
				Author: &accAuthor.Handle,
			}))(t, http.StatusOK)

			ids = dt.Map(clist.JSON200.Nodes, func(c openapi.Node) string { return c.Id })

			a.Contains(ids, node1.JSON200.Id, "admin made this nodeter visible")
			a.Contains(ids, node2.JSON200.Id)
			a.NotContains(ids, node3.JSON200.Id)
			a.NotContains(ids, node4.JSON200.Id)

			// Author can list their own hidden nodeters, but not others.

			clist = testutils.AssertRequest(cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
				Visibility: &[]openapi.Visibility{openapi.Draft},
			}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)

			ids = dt.Map(clist.JSON200.Nodes, func(c openapi.Node) string { return c.Id })

			a.NotContains(ids, node1.JSON200.Id)
			a.NotContains(ids, node2.JSON200.Id)
			a.Contains(ids, node3.JSON200.Id, "this is the only nodeter not published above")
			a.NotContains(ids, node4.JSON200.Id, "owned by someone else, should not be visible")

			// Admin can only list in-review nodeters, but not drafts.

			clist = testutils.AssertRequest(cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
				Visibility: &[]openapi.Visibility{openapi.Review},
			}, e2e.WithSession(ctxAdmin, cj)))(t, http.StatusOK)

			ids = dt.Map(clist.JSON200.Nodes, func(c openapi.Node) string { return c.Id })

			a.NotContains(ids, node1.JSON200.Id)
			a.NotContains(ids, node2.JSON200.Id)
			a.NotContains(ids, node3.JSON200.Id)
			a.NotContains(ids, node4.JSON200.Id)

			// Author moves node3 to in-review

			update3 := testutils.AssertRequest(
				cl.NodeUpdateVisibilityWithResponse(ctx, node3.JSON200.Slug, openapi.VisibilityMutationProps{
					Visibility: openapi.Review,
				}, e2e.WithSession(ctxAuthor, cj)),
			)(t, http.StatusOK)
			a.Equal(openapi.Review, update3.JSON200.Visibility)

			clist = testutils.AssertRequest(cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
				Visibility: &[]openapi.Visibility{openapi.Review},
			}, e2e.WithSession(ctxAdmin, cj)))(t, http.StatusOK)

			ids = dt.Map(clist.JSON200.Nodes, func(c openapi.Node) string { return c.Id })

			a.NotContains(ids, node1.JSON200.Id)
			a.NotContains(ids, node2.JSON200.Id)
			a.Contains(ids, node3.JSON200.Id, "in review so is now visible to admins")
			a.NotContains(ids, node4.JSON200.Id, "")
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

			get404, err := cl.NodeGetWithResponse(ctx, "nonexistent")
			r.NoError(err)
			r.NotNil(get404)
			a.Equal(http.StatusNotFound, get404.StatusCode())

			update403, err := cl.NodeUpdateWithResponse(ctx, "nonexistent", openapi.NodeMutableProps{})
			r.NoError(err)
			r.NotNil(update403)
			a.Equal(http.StatusForbidden, update403.StatusCode())

			update404, err := cl.NodeUpdateWithResponse(ctx, "nonexistent", openapi.NodeMutableProps{}, e2e.WithSession(ctx, cj))
			r.NoError(err)
			r.NotNil(update404)
			a.Equal(http.StatusNotFound, update404.StatusCode())
		}))
	}))
}
