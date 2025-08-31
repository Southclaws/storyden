package visibility_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestChildrenEndpointVisibilityFiltering(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			ctxAdmin, _ := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)
			ctxAuthor, _ := e2e.WithAccount(ctx, aw, seed.Account_003_Baldur)
			ctxRando, _ := e2e.WithAccount(ctx, aw, seed.Account_004_Loki)

			adminSession := sh.WithSession(ctxAdmin)
			authorSession := sh.WithSession(ctxAuthor)
			randoSession := sh.WithSession(ctxRando)

			published := openapi.Published
			draft := openapi.Draft
			unlisted := openapi.Unlisted
			review := openapi.Review

			t.Run("children_endpoint_respects_visibility_rules", func(t *testing.T) {
				a := assert.New(t)

				// Create a published parent node
				parentNode := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Name:       "Published Parent",
					Slug:       opt.New(uuid.NewString()).Ptr(),
					Visibility: &published,
				}, adminSession))(t, http.StatusOK)

				// Author creates draft child
				draftChild := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Name:       "Draft Child",
					Slug:       opt.New(uuid.NewString()).Ptr(),
					Visibility: &draft,
					Parent:     &parentNode.JSON200.Slug,
				}, authorSession))(t, http.StatusOK)

				// Author creates unlisted child
				unlistedChild := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Name:       "Unlisted Child",
					Slug:       opt.New(uuid.NewString()).Ptr(),
					Visibility: &unlisted,
					Parent:     &parentNode.JSON200.Slug,
				}, authorSession))(t, http.StatusOK)

				// Admin creates review child
				reviewChild := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Name:       "Review Child",
					Slug:       opt.New(uuid.NewString()).Ptr(),
					Visibility: &review,
					Parent:     &parentNode.JSON200.Slug,
				}, adminSession))(t, http.StatusOK)

				// Admin creates published child
				publishedChild := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Name:       "Published Child",
					Slug:       opt.New(uuid.NewString()).Ptr(),
					Visibility: &published,
					Parent:     &parentNode.JSON200.Slug,
				}, adminSession))(t, http.StatusOK)

				// Test 1: Unauthenticated user calls /children - should only see published child
				childrenAsGuest := tests.AssertRequest(cl.NodeListChildrenWithResponse(ctx, parentNode.JSON200.Slug, &openapi.NodeListChildrenParams{}))(t, http.StatusOK)
				guestChildrenIDs := dt.Map(childrenAsGuest.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

				a.Contains(guestChildrenIDs, publishedChild.JSON200.Id, "guest should see published child")
				a.NotContains(guestChildrenIDs, draftChild.JSON200.Id, "guest should NOT see draft child")
				a.NotContains(guestChildrenIDs, unlistedChild.JSON200.Id, "guest should NOT see unlisted child")
				a.NotContains(guestChildrenIDs, reviewChild.JSON200.Id, "guest should NOT see review child")

				// Test 2: Random user calls /children - should only see published child
				childrenAsRando := tests.AssertRequest(cl.NodeListChildrenWithResponse(ctx, parentNode.JSON200.Slug, &openapi.NodeListChildrenParams{}, randoSession))(t, http.StatusOK)
				randoChildrenIDs := dt.Map(childrenAsRando.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

				a.Contains(randoChildrenIDs, publishedChild.JSON200.Id, "random user should see published child")
				a.NotContains(randoChildrenIDs, draftChild.JSON200.Id, "random user should NOT see draft child")
				a.NotContains(randoChildrenIDs, unlistedChild.JSON200.Id, "random user should NOT see unlisted child")
				a.NotContains(randoChildrenIDs, reviewChild.JSON200.Id, "random user should NOT see review child")

				// Test 3: Author calls /children - should see their own draft + published
				childrenAsAuthor := tests.AssertRequest(cl.NodeListChildrenWithResponse(ctx, parentNode.JSON200.Slug, &openapi.NodeListChildrenParams{}, authorSession))(t, http.StatusOK)
				authorChildrenIDs := dt.Map(childrenAsAuthor.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

				a.Contains(authorChildrenIDs, publishedChild.JSON200.Id, "author should see published child")
				a.Contains(authorChildrenIDs, draftChild.JSON200.Id, "author should see their own draft child")
				a.Contains(authorChildrenIDs, unlistedChild.JSON200.Id, "author should see their own unlisted child")
				a.NotContains(authorChildrenIDs, reviewChild.JSON200.Id, "author should NOT see admin's review child")

				// Test 4: Admin calls /children - should see published + review (but NOT author's draft/unlisted)
				childrenAsAdmin := tests.AssertRequest(cl.NodeListChildrenWithResponse(ctx, parentNode.JSON200.Slug, &openapi.NodeListChildrenParams{}, adminSession))(t, http.StatusOK)
				adminChildrenIDs := dt.Map(childrenAsAdmin.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

				a.Contains(adminChildrenIDs, publishedChild.JSON200.Id, "admin should see published child")
				a.NotContains(adminChildrenIDs, draftChild.JSON200.Id, "admin should NOT see author's draft child")
				a.Contains(adminChildrenIDs, reviewChild.JSON200.Id, "admin should see review child")
				a.NotContains(adminChildrenIDs, unlistedChild.JSON200.Id, "admin should NOT see unlisted child (it's personal to author)")
			})
		}))
	}))
}

func TestChildrenEndpointVsTreeConsistency(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			ctxAdmin, _ := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)
			ctxAuthor, _ := e2e.WithAccount(ctx, aw, seed.Account_003_Baldur)
			ctxRando, _ := e2e.WithAccount(ctx, aw, seed.Account_004_Loki)

			adminSession := sh.WithSession(ctxAdmin)
			authorSession := sh.WithSession(ctxAuthor)
			randoSession := sh.WithSession(ctxRando)

			published := openapi.Published
			draft := openapi.Draft

			t.Run("children_endpoint_matches_tree_visibility", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				// Create a published parent node
				parentNode := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Name:       "Tree Parent",
					Slug:       opt.New(uuid.NewString()).Ptr(),
					Visibility: &published,
				}, adminSession))(t, http.StatusOK)

				// Author creates a draft child
				draftChild := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Name:       "Draft Child",
					Slug:       opt.New(uuid.NewString()).Ptr(),
					Visibility: &draft,
					Parent:     &parentNode.JSON200.Slug,
				}, authorSession))(t, http.StatusOK)

				// Admin creates a published child
				publishedChild := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Name:       "Published Child",
					Slug:       opt.New(uuid.NewString()).Ptr(),
					Visibility: &published,
					Parent:     &parentNode.JSON200.Slug,
				}, adminSession))(t, http.StatusOK)

				// Test 1: Compare /children endpoint vs /nodes tree for unauthenticated user
				childrenResponse := tests.AssertRequest(cl.NodeListChildrenWithResponse(ctx, parentNode.JSON200.Slug, &openapi.NodeListChildrenParams{}))(t, http.StatusOK)
				childrenIDs := dt.Map(childrenResponse.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

				treeResponse := tests.AssertRequest(cl.NodeGetWithResponse(ctx, parentNode.JSON200.Slug, &openapi.NodeGetParams{}))(t, http.StatusOK)
				treeChildrenIDs := dt.Map(treeResponse.JSON200.Children, func(c openapi.NodeWithChildren) string { return c.Id })

				r.Len(childrenIDs, 1, "children endpoint should only show published child to unauthenticated user")
				r.Len(treeChildrenIDs, 1, "tree endpoint should only show published child to unauthenticated user")
				a.Equal(childrenIDs, treeChildrenIDs, "children and tree endpoints should return the same results for unauthenticated user")
				a.Contains(childrenIDs, publishedChild.JSON200.Id, "both should contain published child")
				a.NotContains(childrenIDs, draftChild.JSON200.Id, "neither should contain draft child")

				// Test 2: Compare /children endpoint vs /nodes tree for author
				childrenAsAuthor := tests.AssertRequest(cl.NodeListChildrenWithResponse(ctx, parentNode.JSON200.Slug, &openapi.NodeListChildrenParams{}, authorSession))(t, http.StatusOK)
				childrenAsAuthorIDs := dt.Map(childrenAsAuthor.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

				treeAsAuthor := tests.AssertRequest(cl.NodeGetWithResponse(ctx, parentNode.JSON200.Slug, &openapi.NodeGetParams{}, authorSession))(t, http.StatusOK)
				treeAsAuthorIDs := dt.Map(treeAsAuthor.JSON200.Children, func(c openapi.NodeWithChildren) string { return c.Id })

				r.Len(childrenAsAuthorIDs, 2, "children endpoint should show both children to author")
				r.Len(treeAsAuthorIDs, 2, "tree endpoint should show both children to author")
				a.ElementsMatch(childrenAsAuthorIDs, treeAsAuthorIDs, "children and tree endpoints should return the same results for author")
				a.Contains(childrenAsAuthorIDs, publishedChild.JSON200.Id, "both should contain published child")
				a.Contains(childrenAsAuthorIDs, draftChild.JSON200.Id, "both should contain author's draft child")

				// Test 3: Compare /children endpoint vs /nodes tree for random user
				childrenAsRando := tests.AssertRequest(cl.NodeListChildrenWithResponse(ctx, parentNode.JSON200.Slug, &openapi.NodeListChildrenParams{}, randoSession))(t, http.StatusOK)
				childrenAsRandoIDs := dt.Map(childrenAsRando.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

				treeAsRando := tests.AssertRequest(cl.NodeGetWithResponse(ctx, parentNode.JSON200.Slug, &openapi.NodeGetParams{}, randoSession))(t, http.StatusOK)
				treeAsRandoIDs := dt.Map(treeAsRando.JSON200.Children, func(c openapi.NodeWithChildren) string { return c.Id })

				r.Len(childrenAsRandoIDs, 1, "children endpoint should only show published child to random user")
				r.Len(treeAsRandoIDs, 1, "tree endpoint should only show published child to random user")
				a.Equal(childrenAsRandoIDs, treeAsRandoIDs, "children and tree endpoints should return the same results for random user")
				a.Contains(childrenAsRandoIDs, publishedChild.JSON200.Id, "both should contain published child")
				a.NotContains(childrenAsRandoIDs, draftChild.JSON200.Id, "neither should contain draft child")
			})

			t.Run("ensure_no_bypass_of_visibility_rules", func(t *testing.T) {
				a := assert.New(t)

				// This test ensures that the /children endpoint cannot be used to bypass
				// visibility rules by creating a scenario where draft content would be
				// exposed if visibility rules were not properly applied.

				// Create a published parent
				parentNode := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Name:       "Security Test Parent",
					Slug:       opt.New(uuid.NewString()).Ptr(),
					Visibility: &published,
				}, adminSession))(t, http.StatusOK)

				// Different authors create draft children
				author1DraftChild := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Name:       "Author1 Draft",
					Slug:       opt.New(uuid.NewString()).Ptr(),
					Visibility: &draft,
					Parent:     &parentNode.JSON200.Slug,
				}, authorSession))(t, http.StatusOK)

				author2DraftChild := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Name:       "Author2 Draft",
					Slug:       opt.New(uuid.NewString()).Ptr(),
					Visibility: &draft,
					Parent:     &parentNode.JSON200.Slug,
				}, randoSession))(t, http.StatusOK)

				// Unauthenticated user should see no draft children
				unauthChildren := tests.AssertRequest(cl.NodeListChildrenWithResponse(ctx, parentNode.JSON200.Slug, &openapi.NodeListChildrenParams{}))(t, http.StatusOK)
				unauthChildrenIDs := dt.Map(unauthChildren.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

				a.NotContains(unauthChildrenIDs, author1DraftChild.JSON200.Id, "unauthenticated user should not see author1's draft")
				a.NotContains(unauthChildrenIDs, author2DraftChild.JSON200.Id, "unauthenticated user should not see author2's draft")

				// Author1 should only see their own draft
				author1Children := tests.AssertRequest(cl.NodeListChildrenWithResponse(ctx, parentNode.JSON200.Slug, &openapi.NodeListChildrenParams{}, authorSession))(t, http.StatusOK)
				author1ChildrenIDs := dt.Map(author1Children.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

				a.Contains(author1ChildrenIDs, author1DraftChild.JSON200.Id, "author1 should see their own draft")
				a.NotContains(author1ChildrenIDs, author2DraftChild.JSON200.Id, "author1 should not see author2's draft")

				// Author2 should only see their own draft
				author2Children := tests.AssertRequest(cl.NodeListChildrenWithResponse(ctx, parentNode.JSON200.Slug, &openapi.NodeListChildrenParams{}, randoSession))(t, http.StatusOK)
				author2ChildrenIDs := dt.Map(author2Children.JSON200.Nodes, func(c openapi.NodeWithChildren) string { return c.Id })

				a.NotContains(author2ChildrenIDs, author1DraftChild.JSON200.Id, "author2 should not see author1's draft")
				a.Contains(author2ChildrenIDs, author2DraftChild.JSON200.Id, "author2 should see their own draft")
			})
		}))
	}))
}
