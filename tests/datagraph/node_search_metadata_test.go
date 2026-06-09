package search_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/search/search_indexer"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestSearchNodeMetadata(t *testing.T) {
	t.Parallel()

	runNodeMetadataSearchTest(t, &config.Config{}, false, nil)
}

func TestSearchNodeMetadataBleve(t *testing.T) {
	t.Parallel()

	bleveName := time.Now().Format(time.RFC3339) + t.Name()
	runNodeMetadataSearchTest(t, &config.Config{
		SearchProvider: "bleve",
		BlevePath:      fmt.Sprintf("data/%s.bleve", bleveName),
	}, true, func(ctx context.Context, idx *search_indexer.Indexer) {
		require.NoError(t, idx.ReindexAll(ctx))
	})
}

func runNodeMetadataSearchTest(t *testing.T, cfg *config.Config, includeTagQuery bool, reindex func(context.Context, *search_indexer.Indexer)) {
	t.Helper()

	integration.Test(t, cfg, e2e.Setup(), fx.Invoke(func(
		root context.Context,
		lc fx.Lifecycle,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		idx *search_indexer.Indexer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			published := openapi.Published
			nameNeedle := "tidal-" + uuid.NewString()
			tagNeedle := openapi.TagName("aurora-" + uuid.NewString())
			content := "<p>body text without searchable metadata terms</p>"

			nodeResp, err := cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
				Name:        "Getting Blocked from Using " + nameNeedle,
				Description: ptr(openapi.NodeDescription("Description mentions " + uuid.NewString())),
				Content:     &content,
				Tags:        &[]openapi.TagName{tagNeedle},
				Visibility:  &published,
			}, adminSession)
			tests.Ok(t, err, nodeResp)

			if reindex != nil {
				reindex(root, idx)
			}

			nodeKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindNode}
			cases := []struct {
				name  string
				query string
			}{
				{name: "search_by_name", query: nameNeedle},
			}
			if includeTagQuery {
				cases = append(cases, struct {
					name  string
					query string
				}{name: "search_by_tag", query: string(tagNeedle)})
			}

			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					r := require.New(t)

					resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
						Q:    tc.query,
						Kind: &nodeKind,
					}, adminSession)
					tests.Ok(t, err, resp)

					found := findItem(resp.JSON200.Items, nodeResp.JSON200.Id)
					r.NotNil(found)
				})
			}
		}))
	}))
}

func ptr[T any](v T) *T {
	return &v
}
