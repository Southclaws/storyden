package visibility_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestNodesVisibility(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		cj *session.Jar,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			a := assert.New(t)

			ctxAdmin, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			ctxAuthor, accAuthor := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			ctxRando, _ := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			adminSession := e2e.WithSession(ctxAdmin, cj)
			authorSession := e2e.WithSession(ctxAuthor, cj)
			randoSession := e2e.WithSession(ctxRando, cj)

			t.Run("public_only", func(t *testing.T) {
				t.Parallel()

				// Public listing without filters does not contain any of them
				// because they were created without being published.

				node1 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n1", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node2 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n2", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node3 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n3", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node4 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n4", Slug: opt.New(uuid.NewString()).Ptr()}, randoSession))(t, http.StatusOK)

				clist := tests.AssertRequest(cl.NodeListWithResponse(root, &openapi.NodeListParams{}))(t, http.StatusOK)

				ids := dt.Map(clist.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

				// List does not contain any because they have not been published
				// and the request was made without auth from the owner.
				a.NotContains(ids, node1.JSON200.Id)
				a.NotContains(ids, node2.JSON200.Id)
				a.NotContains(ids, node3.JSON200.Id)
				a.NotContains(ids, node4.JSON200.Id)
			})

			t.Run("public_filter_by_author", func(t *testing.T) {
				t.Parallel()

				// Public listing with author filter does not contain any of
				// the nodes because they have not been published.

				node1 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n1", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node2 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n2", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node3 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n3", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node4 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n4", Slug: opt.New(uuid.NewString()).Ptr()}, randoSession))(t, http.StatusOK)

				clist := tests.AssertRequest(cl.NodeListWithResponse(root, &openapi.NodeListParams{
					Author: &accAuthor.Handle,
				}))(t, http.StatusOK)

				ids := dt.Map(clist.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

				a.NotContains(ids, node1.JSON200.Id)
				a.NotContains(ids, node2.JSON200.Id)
				a.NotContains(ids, node3.JSON200.Id)
				a.NotContains(ids, node4.JSON200.Id)
			})

			t.Run("admin_change_visibility", func(t *testing.T) {
				t.Parallel()

				// NOTE: This should actually fail because this node is in draft
				// and not submitted for review. The admin will not be able to
				// list this node because it is not in review or published.

				node1 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n1", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node2 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n2", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node3 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n3", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node4 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n4", Slug: opt.New(uuid.NewString()).Ptr()}, randoSession))(t, http.StatusOK)

				update1 := tests.AssertRequest(
					cl.NodeUpdateVisibilityWithResponse(root, node1.JSON200.Slug, openapi.VisibilityMutationProps{
						Visibility: openapi.Published,
					}, adminSession),
				)(t, http.StatusOK)
				a.Equal(openapi.Published, update1.JSON200.Visibility)

				clist := tests.AssertRequest(cl.NodeListWithResponse(root, &openapi.NodeListParams{
					Author: &accAuthor.Handle,
				}))(t, http.StatusOK)

				ids := dt.Map(clist.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

				a.Contains(ids, node1.JSON200.Id, "admin made this node visible")
				a.NotContains(ids, node2.JSON200.Id)
				a.NotContains(ids, node3.JSON200.Id)
				a.NotContains(ids, node4.JSON200.Id)
			})

			// Author can change visibility
			t.Run("author_change_visibility", func(t *testing.T) {
				t.Parallel()

				node1 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n1", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node2 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n2", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node3 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n3", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node4 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n4", Slug: opt.New(uuid.NewString()).Ptr()}, randoSession))(t, http.StatusOK)

				update2 := tests.AssertRequest(
					cl.NodeUpdateVisibilityWithResponse(root, node2.JSON200.Slug, openapi.VisibilityMutationProps{
						Visibility: openapi.Published,
					}, authorSession),
				)(t, http.StatusOK)
				a.Equal(openapi.Published, update2.JSON200.Visibility)

				clist := tests.AssertRequest(cl.NodeListWithResponse(root, &openapi.NodeListParams{
					Author: &accAuthor.Handle,
				}))(t, http.StatusOK)

				ids := dt.Map(clist.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

				a.NotContains(ids, node1.JSON200.Id)
				a.Contains(ids, node2.JSON200.Id)
				a.NotContains(ids, node3.JSON200.Id)
				a.NotContains(ids, node4.JSON200.Id)
			})

			t.Run("author_can_view_own_drafts", func(t *testing.T) {
				t.Parallel()

				node1 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n1", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node2 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n2", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node3 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n3", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node4 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n4", Slug: opt.New(uuid.NewString()).Ptr()}, randoSession))(t, http.StatusOK)

				clist := tests.AssertRequest(cl.NodeListWithResponse(root, &openapi.NodeListParams{
					Visibility: &[]openapi.Visibility{openapi.Draft},
				}, authorSession))(t, http.StatusOK)

				ids := dt.Map(clist.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

				a.Contains(ids, node1.JSON200.Id, "owned by author")
				a.Contains(ids, node2.JSON200.Id, "owned by author")
				a.Contains(ids, node3.JSON200.Id, "owned by author")
				a.NotContains(ids, node4.JSON200.Id, "owned by someone else, should not be visible")
			})

			t.Run("admin_lists_in_review_but_not_drafts", func(t *testing.T) {
				t.Parallel()

				node1 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n1", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node2 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n2", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node3 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n3", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node4 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n4", Slug: opt.New(uuid.NewString()).Ptr()}, randoSession))(t, http.StatusOK)

				clist := tests.AssertRequest(cl.NodeListWithResponse(root, &openapi.NodeListParams{
					Visibility: &[]openapi.Visibility{openapi.Review},
				}, adminSession))(t, http.StatusOK)

				ids := dt.Map(clist.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

				a.NotContains(ids, node1.JSON200.Id)
				a.NotContains(ids, node2.JSON200.Id)
				a.NotContains(ids, node3.JSON200.Id)
				a.NotContains(ids, node4.JSON200.Id)
			})

			t.Run("author_submits_for_review", func(t *testing.T) {
				t.Parallel()

				node1 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n1", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node2 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n2", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node3 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n3", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node4 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n4", Slug: opt.New(uuid.NewString()).Ptr()}, randoSession))(t, http.StatusOK)

				update3 := tests.AssertRequest(
					cl.NodeUpdateVisibilityWithResponse(root, node3.JSON200.Slug, openapi.VisibilityMutationProps{
						Visibility: openapi.Review,
					}, authorSession),
				)(t, http.StatusOK)
				a.Equal(openapi.Review, update3.JSON200.Visibility)

				clist := tests.AssertRequest(cl.NodeListWithResponse(root, &openapi.NodeListParams{
					Visibility: &[]openapi.Visibility{openapi.Review},
				}, adminSession))(t, http.StatusOK)

				ids := dt.Map(clist.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

				a.NotContains(ids, node1.JSON200.Id)
				a.NotContains(ids, node2.JSON200.Id)
				a.Contains(ids, node3.JSON200.Id, "in review so is now visible to admins")
				a.NotContains(ids, node4.JSON200.Id, "")
			})

			t.Run("author_submmits_unlisted", func(t *testing.T) {
				t.Parallel()

				node1 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n1", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node2 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n2", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node3 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n3", Slug: opt.New(uuid.NewString()).Ptr()}, authorSession))(t, http.StatusOK)
				node4 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n4", Slug: opt.New(uuid.NewString()).Ptr()}, randoSession))(t, http.StatusOK)

				update3 := tests.AssertRequest(
					cl.NodeUpdateVisibilityWithResponse(root, node3.JSON200.Slug, openapi.VisibilityMutationProps{
						Visibility: openapi.Unlisted,
					}, authorSession),
				)(t, http.StatusOK)
				a.Equal(openapi.Unlisted, update3.JSON200.Visibility)

				clist := tests.AssertRequest(cl.NodeListWithResponse(root, &openapi.NodeListParams{
					// Visibility: &[]openapi.Visibility{openapi.Unlisted},
				}, adminSession))(t, http.StatusOK)

				ids := dt.Map(clist.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

				a.NotContains(ids, node1.JSON200.Id)
				a.NotContains(ids, node2.JSON200.Id)
				a.NotContains(ids, node3.JSON200.Id, "")
				a.NotContains(ids, node4.JSON200.Id, "")
			})

			t.Run("visibility_affects_children", func(t *testing.T) {
				t.Parallel()

				published := openapi.Published
				draft := openapi.Draft

				node1 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n1", Slug: opt.New(uuid.NewString()).Ptr(), Visibility: &published}, adminSession))(t, http.StatusOK)
				node2 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n2", Slug: opt.New(uuid.NewString()).Ptr(), Visibility: &draft}, authorSession))(t, http.StatusOK)
				node3 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n3", Slug: opt.New(uuid.NewString()).Ptr(), Visibility: &published, Parent: &node1.JSON200.Slug}, adminSession))(t, http.StatusOK)
				node4 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n4", Slug: opt.New(uuid.NewString()).Ptr(), Visibility: &published, Parent: &node1.JSON200.Slug}, adminSession))(t, http.StatusOK)

				clist := tests.AssertRequest(cl.NodeListWithResponse(root, &openapi.NodeListParams{}, adminSession))(t, http.StatusOK)

				ids := dt.Map(clist.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

				a.Contains(ids, node1.JSON200.Id)
				a.NotContains(ids, node2.JSON200.Id, "visibility is not published")
				a.NotContains(ids, node3.JSON200.Id, "visibility is published, but is a child of a node that is not")
				a.NotContains(ids, node4.JSON200.Id, "visibility is published, but is a child of a node that is not")
			})

			t.Run("only_author_sees_non_published_children", func(t *testing.T) {
				t.Parallel()

				published := openapi.Published
				draft := openapi.Draft

				node1 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n1", Slug: opt.New(uuid.NewString()).Ptr(), Visibility: &published}, adminSession))(t, http.StatusOK)
				node2 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n2", Slug: opt.New(uuid.NewString()).Ptr(), Visibility: &draft, Parent: &node1.JSON200.Slug}, authorSession))(t, http.StatusOK)
				node3 := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{Name: "n2", Slug: opt.New(uuid.NewString()).Ptr(), Visibility: &draft, Parent: &node1.JSON200.Slug}, randoSession))(t, http.StatusOK)

				get1asAuthor := tests.AssertRequest(cl.NodeGetWithResponse(root, node1.JSON200.Slug, authorSession))(t, http.StatusOK)
				ids := dt.Map(get1asAuthor.JSON200.Children, func(c openapi.NodeWithChildren) string { return c.Id })
				a.Contains(ids, node2.JSON200.Id, "author can see child of node1 because they own it")
				a.NotContains(ids, node3.JSON200.Id, "cannot see node3 as it is not owned by the author")

				get1asRando := tests.AssertRequest(cl.NodeGetWithResponse(root, node1.JSON200.Slug, randoSession))(t, http.StatusOK)
				ids = dt.Map(get1asRando.JSON200.Children, func(c openapi.NodeWithChildren) string { return c.Id })
				a.NotContains(ids, node2.JSON200.Id, "cannot see node2 as it is not published and owned by another member")
				a.Contains(ids, node3.JSON200.Id, "can see node3 as it is owned by this member")

				get1asGuest := tests.AssertRequest(cl.NodeGetWithResponse(root, node1.JSON200.Slug))(t, http.StatusOK)
				ids = dt.Map(get1asGuest.JSON200.Children, func(c openapi.NodeWithChildren) string { return c.Id })
				a.NotContains(ids, node2.JSON200.Id, "guest cannot see node2 as it is not published")
				a.NotContains(ids, node3.JSON200.Id, "guest cannot see node3 as it is not published")
			})
		}))
	}))
}
