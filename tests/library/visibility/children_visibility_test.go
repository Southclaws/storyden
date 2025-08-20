package visibility_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
					Name: "Published Parent", 
					Slug: opt.New(uuid.NewString()).Ptr(), 
					Visibility: &published,
				}, adminSession))(t, http.StatusOK)

				// Author creates draft child
				draftChild := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Name: "Draft Child", 
					Slug: opt.New(uuid.NewString()).Ptr(), 
					Visibility: &draft,
					Parent: &parentNode.JSON200.Slug,
				}, authorSession))(t, http.StatusOK)

				// Author creates unlisted child  
				unlistedChild := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Name: "Unlisted Child", 
					Slug: opt.New(uuid.NewString()).Ptr(), 
					Visibility: &unlisted,
					Parent: &parentNode.JSON200.Slug,
				}, authorSession))(t, http.StatusOK)

				// Admin creates review child
				reviewChild := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Name: "Review Child", 
					Slug: opt.New(uuid.NewString()).Ptr(), 
					Visibility: &review,
					Parent: &parentNode.JSON200.Slug,
				}, adminSession))(t, http.StatusOK)

				// Admin creates published child
				publishedChild := tests.AssertRequest(cl.NodeCreateWithResponse(ctx, openapi.NodeInitialProps{
					Name: "Published Child", 
					Slug: opt.New(uuid.NewString()).Ptr(), 
					Visibility: &published,
					Parent: &parentNode.JSON200.Slug,
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