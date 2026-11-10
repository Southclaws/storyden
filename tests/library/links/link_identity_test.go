package links_test

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/link/link_querier"
	"github.com/Southclaws/storyden/app/resources/link/link_url"
	"github.com/Southclaws/storyden/app/resources/link/link_writer"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/internal/ent"
	entlink "github.com/Southclaws/storyden/internal/ent/link"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestLinkIdentity(t *testing.T) {
	t.Parallel()
	integration.Test(t, nil, fx.Invoke(func(lc fx.Lifecycle, ctx context.Context, writer *link_writer.LinkWriter, query *link_querier.LinkQuerier, db *ent.Client) {
		lc.Append(fx.StartHook(func() {
			host := uuid.NewString() + ".example.com"
			t.Run("distinct_resources_cannot_overwrite_each_other", func(t *testing.T) {
				urls := []string{
					"https://" + host + "/watch?v=1", "https://" + host + "/watch?v=2",
					"https://" + host + "/watch#one", "https://" + host + "/watch#two",
					"https://" + host + "/a-b", "https://" + host + "/a/b",
					"https://" + host + "/Case", "https://" + host + "/case",
					"https://" + host + "/a%2Fb",
					"http://" + host + "/a-b", "https://www." + host + "/a-b",
				}
				ids := map[xid.ID]bool{}
				slugs := map[string]bool{}
				for i, u := range urls {
					ln, err := writer.Store(ctx, parseLinkURL(t, u), fmt.Sprintf("title %d", i), fmt.Sprintf("description %d", i))
					require.NoError(t, err)
					require.False(t, ids[ln.ID], u)
					require.False(t, slugs[ln.Slug], u)
					ids[ln.ID] = true
					slugs[ln.Slug] = true
				}

				for i, u := range urls {
					ln, err := writer.Store(ctx, parseLinkURL(t, u), "", "")
					require.NoError(t, err)
					require.Equal(t, u, ln.URL.String())
					require.Equal(t, fmt.Sprintf("title %d", i), ln.Title.String())
					require.Equal(t, fmt.Sprintf("description %d", i), ln.Description.String())
					require.True(t, ids[ln.ID])
					require.True(t, slugs[ln.Slug])
				}
			})

			t.Run("normalizes_safe_equivalents", func(t *testing.T) {
				first, err := writer.Store(ctx, parseLinkURL(t, "https://"+host), "root", "root description")
				require.NoError(t, err)
				second, err := writer.Store(ctx, parseLinkURL(t, "HTTPS://"+host+"/"), "", "")
				require.NoError(t, err)
				require.Equal(t, first.ID, second.ID)
				require.Equal(t, first.Slug, second.Slug)
				require.Equal(t, "https://"+host+"/", second.URL.String())
				require.Equal(t, host, second.Domain)
				require.False(t, second.CreatedAt.IsZero())
			})

			t.Run("preserves_existing_slugs_and_repairs_the_domain", func(t *testing.T) {
				legacyHost := uuid.NewString() + ".example.com"
				raw := "https://" + legacyHost + "/legacy"
				legacySlug := "legacy-" + uuid.NewString()
				legacy, err := db.Link.Create().SetURL(raw).SetSlug(legacySlug).
					SetDomain("incorrect.example.com").SetTitle("Legacy title").SetDescription("Legacy description").Save(ctx)
				require.NoError(t, err)
				updated, err := writer.Store(ctx, parseLinkURL(t, raw), "", "")
				require.NoError(t, err)
				require.Equal(t, legacy.ID, updated.ID)
				require.Equal(t, legacySlug, updated.Slug)
				require.Equal(t, "Legacy title", updated.Title.String())
				require.Equal(t, legacyHost, updated.Domain)
				found, err := query.Get(ctx, legacySlug)
				require.NoError(t, err)
				require.Equal(t, legacy.ID, found.ID)
				folded, err := query.Get(ctx, strings.ToUpper(legacySlug))
				require.NoError(t, err)
				require.Equal(t, legacy.ID, folded.ID)
			})

			t.Run("URL_search_respects_hosts_and_finds_www_variants", func(t *testing.T) {
				prefix := uuid.NewString()
				co := "https://" + prefix + ".co/"
				want, err := writer.Store(ctx, parseLinkURL(t, co), "co title", "co description")
				require.NoError(t, err)
				_, err = writer.Store(ctx, parseLinkURL(t, "https://"+prefix+".com/"), co, co)
				require.NoError(t, err)
				for _, q := range []string{co, "https://www." + prefix + ".co/", "https://www." + prefix + ".co", "http://" + prefix + ".co", prefix + ".co", "www." + prefix + ".co/"} {
					result, err := query.Search(ctx, pagination.NewPageParams(1, 10), link_querier.WithKeyword(q))
					require.NoError(t, err)
					require.Len(t, result.Items, 1, q)
					require.Equal(t, want.ID, result.Items[0].ID, q)
					require.Equal(t, 1, result.Results)
				}
			})

			t.Run("scrape_updates_preserve_good_metadata_on_failure_or_missing_fields", func(t *testing.T) {
				ln, err := writer.Store(ctx, parseLinkURL(t, "https://"+uuid.NewString()+".example.com/"), "good title", "good description")
				require.NoError(t, err)
				now := time.Now().UTC().Truncate(time.Microsecond)
				require.NoError(t, writer.FinishScrape(ctx, ln.ID, now, link_writer.ScrapeResult{Success: true, Title: "updated title"}))
				require.NoError(t, writer.FinishScrape(ctx, ln.ID, now.Add(time.Minute), link_writer.ScrapeResult{}))
				stored, err := db.Link.Get(ctx, ln.ID)
				require.NoError(t, err)
				require.Equal(t, "updated title", stored.Title)
				require.Equal(t, "good description", stored.Description)
				require.Equal(t, entlink.ScrapeStatusFailed, stored.ScrapeStatus)
				require.WithinDuration(t, now, *stored.ScrapedAt, time.Microsecond)
			})

			t.Run("concurrent_saves_of_one_URL_share_one_record", func(t *testing.T) {
				raw := "https://" + uuid.NewString() + ".example.com/concurrent"
				var wg sync.WaitGroup
				ids := make(chan xid.ID, 8)
				errors := make(chan error, 8)
				for range 8 {
					wg.Go(func() {
						ln, err := writer.Store(ctx, parseLinkURL(t, raw), "shared title", "shared description")
						errors <- err
						if err == nil {
							ids <- ln.ID
						}
					})
				}

				wg.Wait()
				close(errors)
				for err := range errors {
					require.NoError(t, err)
				}

				close(ids)
				unique := map[xid.ID]bool{}
				for id := range ids {
					unique[id] = true
				}

				require.Len(t, unique, 1)
				count, err := db.Link.Query().Where(entlink.URLEQ(raw)).Count(ctx)
				require.NoError(t, err)
				require.Equal(t, 1, count)
			})
		}))
	}))
}

func parseLinkURL(t *testing.T, raw string) link_url.URL {
	t.Helper()
	u, err := link_url.Parse(raw)
	require.NoError(t, err)
	return u
}
