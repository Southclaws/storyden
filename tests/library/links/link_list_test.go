package links_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/link/link_querier"
	"github.com/Southclaws/storyden/app/resources/link/link_writer"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestLinkListPaging(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		lw *link_writer.LinkWriter,
		lq *link_querier.LinkQuerier,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			ctx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			session := sh.WithSession(ctx)

			run := uuid.NewString()[:8]
			other := uuid.NewString()[:8]

			const matching = 4
			for i := range matching {
				_, err := lw.Store(ctx,
					fmt.Sprintf("https://%s.example.com/%d", run, i),
					fmt.Sprintf("%s title %d", run, i),
					"a link that matches the keyword",
				)
				r.NoError(err)
			}

			const unrelated = 3
			for i := range unrelated {
				_, err := lw.Store(ctx,
					fmt.Sprintf("https://%s.example.com/%d", other, i),
					fmt.Sprintf("%s title %d", other, i),
					"no keyword here",
				)
				r.NoError(err)
			}

			t.Run("counts_only_the_filtered_rows", func(t *testing.T) {
				result, err := lq.Search(root, pagination.NewPageParams(1, 2), link_querier.WithKeyword(run))
				r.NoError(err)

				a.Equal(2, result.TotalPages, "four matching links at two per page is two pages")
				a.Len(result.Items, 2)
			})

			t.Run("walks_every_page_exactly_once", func(t *testing.T) {
				seen := map[string]struct{}{}

				for page := uint(1); page <= 10; page++ {
					result, err := lq.Search(root, pagination.NewPageParams(page, 2), link_querier.WithKeyword(run))
					r.NoError(err)
					r.Equal(int(page), result.CurrentPage)

					for _, l := range result.Items {
						_, duplicate := seen[l.URL]
						a.False(duplicate, "link %s was returned on more than one page", l.URL)
						seen[l.URL] = struct{}{}
					}

					next, ok := result.NextPage.Get()
					if !ok {
						break
					}

					a.Equal(int(page)+1, next)
				}

				a.Len(seen, matching)
			})

			t.Run("exact_url_lookup_ignores_longer_urls", func(t *testing.T) {
				result, err := lq.Search(root, pagination.NewPageParams(1, 10),
					link_querier.WithURL(fmt.Sprintf("https://%s.example.com/", run)),
				)
				r.NoError(err)
				a.Empty(result.Items, "a prefix of a stored URL is not that URL")
			})

			t.Run("api_reports_a_next_page_the_client_can_ask_for", func(t *testing.T) {
				first, err := cl.LinkListWithResponse(root, &openapi.LinkListParams{
					Q: &run,
				}, session)
				tests.Ok(t, err, first)

				a.Equal(1, first.JSON200.CurrentPage)
				a.Nil(first.JSON200.NextPage, "every match fits on the default page")
			})
		}))
	}))
}