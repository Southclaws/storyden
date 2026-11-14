package links_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Southclaws/dt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/link/link_url"
	"github.com/Southclaws/storyden/app/resources/link/link_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/link/scrape"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestLinkGetHidesUnpublishedContent(t *testing.T) {
	t.Parallel()

	var (
		root     context.Context
		client   *openapi.ClientWithResponses
		sessions *e2e.SessionHelper
		accounts *account_writer.Writer
		links    *link_writer.LinkWriter
	)
	script := &scriptedScraper{}
	script.set(func(context.Context, link_url.URL) (*scrape.WebContent, error) {
		return &scrape.WebContent{Title: "Linked"}, nil
	})
	integration.Test(t, nil, e2e.Setup(),
		fx.Decorate(func(scrape.Scraper) scrape.Scraper { return script }),
		fx.Populate(&root, &client, &sessions, &accounts, &links),
	)

	ctx, _ := e2e.WithAccount(root, accounts, seed.Account_004_Loki)
	member := sessions.WithSession(ctx)
	adminCtx, _ := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
	admin := sessions.WithSession(adminCtx)

	raw := "https://" + uuid.NewString() + ".example.com/"
	bookmark, err := links.Store(ctx, parseLinkURL(t, raw), "Linked", "")
	require.NoError(t, err)

	draft := openapi.VisibilityDraft
	secret := "<p>draft secret " + uuid.NewString() + "</p>"

	published := openapi.VisibilityPublished
	publishedPage := tests.AssertRequest(client.NodeCreateWithResponse(adminCtx, openapi.NodeCreateJSONRequestBody{Name: uuid.NewString(), Url: &raw, Visibility: &published}, admin))(t, http.StatusOK)
	draftPage := tests.AssertRequest(client.NodeCreateWithResponse(ctx, openapi.NodeCreateJSONRequestBody{Name: uuid.NewString(), Url: &raw, Content: &secret, Visibility: &draft}, member))(t, http.StatusOK)

	publishedThread := tests.AssertRequest(client.ThreadCreateWithResponse(ctx, openapi.ThreadCreateJSONRequestBody{Title: uuid.NewString(), Body: &secret, Url: &raw, Visibility: &published}, member))(t, http.StatusOK)
	draftThread := tests.AssertRequest(client.ThreadCreateWithResponse(ctx, openapi.ThreadCreateJSONRequestBody{Title: uuid.NewString(), Body: &secret, Url: &raw, Visibility: &draft}, member))(t, http.StatusOK)

	for name, editors := range map[string][]openapi.RequestEditorFn{"guest": nil, "member": {member}} {
		t.Run(name, func(t *testing.T) {
			got := tests.AssertRequest(client.LinkGetWithResponse(root, bookmark.Slug, editors...))(t, http.StatusOK)

			nodeIDs := dt.Map(got.JSON200.Nodes, func(n openapi.Node) string { return n.Id })
			assert.Contains(t, nodeIDs, publishedPage.JSON200.Id)
			assert.NotContains(t, nodeIDs, draftPage.JSON200.Id)

			postIDs := dt.Map(got.JSON200.Posts, func(p openapi.PostReference) string { return p.Id })
			assert.Contains(t, postIDs, publishedThread.JSON200.Id)
			assert.NotContains(t, postIDs, draftThread.JSON200.Id)

			for _, n := range got.JSON200.Nodes {
				if n.Content != nil {
					assert.NotContains(t, *n.Content, "draft secret")
				}
			}
		})
	}
}

