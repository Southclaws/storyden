package node_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/openapi"
	"github.com/Southclaws/storyden/app/transports/openapi/bindings"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests/testutils"
)

func TestNodesVisibility(t *testing.T) {
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

			ctxAdmin, _ := e2e.WithAccount(ctx, ar, seed.Account_001_Odin)
			ctxAuthor, accAuthor := e2e.WithAccount(ctx, ar, seed.Account_003_Baldur)
			ctxRando, _ := e2e.WithAccount(ctx, ar, seed.Account_004_Loki)

			t.Run("public_only", func(t *testing.T) {
				t.Parallel()

				// Public listing without filters does not contain any of them
				// because they were created without being published.

				node1 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n1", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node2 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n2", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node3 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n3", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node4 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n4", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxRando, cj)))(t, http.StatusOK)

				clist := testutils.AssertRequest(cl.NodeListWithResponse(ctx, &openapi.NodeListParams{}))(t, http.StatusOK)

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

				node1 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n1", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node2 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n2", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node3 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n3", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node4 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n4", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxRando, cj)))(t, http.StatusOK)

				clist := testutils.AssertRequest(cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
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

				node1 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n1", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node2 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n2", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node3 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n3", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node4 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n4", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxRando, cj)))(t, http.StatusOK)

				update1 := testutils.AssertRequest(
					cl.NodeUpdateVisibilityWithResponse(ctx, node1.JSON200.Slug, openapi.VisibilityMutationProps{
						Visibility: openapi.Published,
					}, e2e.WithSession(ctxAdmin, cj)),
				)(t, http.StatusOK)
				a.Equal(openapi.Published, update1.JSON200.Visibility)

				clist := testutils.AssertRequest(cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
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

				node1 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n1", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node2 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n2", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node3 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n3", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node4 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n4", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxRando, cj)))(t, http.StatusOK)

				update2 := testutils.AssertRequest(
					cl.NodeUpdateVisibilityWithResponse(ctx, node2.JSON200.Slug, openapi.VisibilityMutationProps{
						Visibility: openapi.Published,
					}, e2e.WithSession(ctxAuthor, cj)),
				)(t, http.StatusOK)
				a.Equal(openapi.Published, update2.JSON200.Visibility)

				clist := testutils.AssertRequest(cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
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

				node1 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n1", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node2 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n2", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node3 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n3", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node4 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n4", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxRando, cj)))(t, http.StatusOK)

				clist := testutils.AssertRequest(cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
					Visibility: &[]openapi.Visibility{openapi.Draft},
				}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)

				ids := dt.Map(clist.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

				a.Contains(ids, node1.JSON200.Id, "owned by author")
				a.Contains(ids, node2.JSON200.Id, "owned by author")
				a.Contains(ids, node3.JSON200.Id, "owned by author")
				a.NotContains(ids, node4.JSON200.Id, "owned by someone else, should not be visible")
			})

			t.Run("admin_lists_in_review_but_not_drafts", func(t *testing.T) {
				t.Parallel()

				node1 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n1", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node2 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n2", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node3 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n3", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node4 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n4", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxRando, cj)))(t, http.StatusOK)

				clist := testutils.AssertRequest(cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
					Visibility: &[]openapi.Visibility{openapi.Review},
				}, e2e.WithSession(ctxAdmin, cj)))(t, http.StatusOK)

				ids := dt.Map(clist.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

				a.NotContains(ids, node1.JSON200.Id)
				a.NotContains(ids, node2.JSON200.Id)
				a.NotContains(ids, node3.JSON200.Id)
				a.NotContains(ids, node4.JSON200.Id)
			})

			t.Run("author_submits_for_review", func(t *testing.T) {
				t.Parallel()

				node1 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n1", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node2 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n2", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node3 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n3", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node4 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n4", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxRando, cj)))(t, http.StatusOK)

				update3 := testutils.AssertRequest(
					cl.NodeUpdateVisibilityWithResponse(ctx, node3.JSON200.Slug, openapi.VisibilityMutationProps{
						Visibility: openapi.Review,
					}, e2e.WithSession(ctxAuthor, cj)),
				)(t, http.StatusOK)
				a.Equal(openapi.Review, update3.JSON200.Visibility)

				clist := testutils.AssertRequest(cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
					Visibility: &[]openapi.Visibility{openapi.Review},
				}, e2e.WithSession(ctxAdmin, cj)))(t, http.StatusOK)

				ids := dt.Map(clist.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

				a.NotContains(ids, node1.JSON200.Id)
				a.NotContains(ids, node2.JSON200.Id)
				a.Contains(ids, node3.JSON200.Id, "in review so is now visible to admins")
				a.NotContains(ids, node4.JSON200.Id, "")
			})

			t.Run("author_submmits_unlisted", func(t *testing.T) {
				t.Parallel()

				node1 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n1", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node2 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n2", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node3 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n3", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node4 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n4", Slug: opt.New(uuid.NewString()).Ptr()}, e2e.WithSession(ctxRando, cj)))(t, http.StatusOK)

				update3 := testutils.AssertRequest(
					cl.NodeUpdateVisibilityWithResponse(ctx, node3.JSON200.Slug, openapi.VisibilityMutationProps{
						Visibility: openapi.Unlisted,
					}, e2e.WithSession(ctxAuthor, cj)),
				)(t, http.StatusOK)
				a.Equal(openapi.Unlisted, update3.JSON200.Visibility)

				clist := testutils.AssertRequest(cl.NodeListWithResponse(ctx, &openapi.NodeListParams{
					// Visibility: &[]openapi.Visibility{openapi.Unlisted},
				}, e2e.WithSession(ctxAdmin, cj)))(t, http.StatusOK)

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

				node1 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n1", Slug: opt.New(uuid.NewString()).Ptr(), Visibility: &published}, e2e.WithSession(ctxAdmin, cj)))(t, http.StatusOK)
				node2 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n2", Slug: opt.New(uuid.NewString()).Ptr(), Visibility: &draft}, e2e.WithSession(ctxAuthor, cj)))(t, http.StatusOK)
				node3 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n3", Slug: opt.New(uuid.NewString()).Ptr(), Visibility: &published, Parent: &node1.JSON200.Slug}, e2e.WithSession(ctxAdmin, cj)))(t, http.StatusOK)
				node4 := testutils.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{Name: "n4", Slug: opt.New(uuid.NewString()).Ptr(), Visibility: &published, Parent: &node1.JSON200.Slug}, e2e.WithSession(ctxAdmin, cj)))(t, http.StatusOK)

				clist := testutils.AssertRequest(cl.NodeListWithResponse(ctx, &openapi.NodeListParams{}, e2e.WithSession(ctxAdmin, cj)))(t, http.StatusOK)

				ids := dt.Map(clist.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

				a.Contains(ids, node1.JSON200.Id)
				a.NotContains(ids, node2.JSON200.Id, "visibility is not published")
				a.NotContains(ids, node3.JSON200.Id, "visibility is published, but is a child of a node that is not")
				a.NotContains(ids, node4.JSON200.Id, "visibility is published, but is a child of a node that is not")
			})
		}))
	}))
}
