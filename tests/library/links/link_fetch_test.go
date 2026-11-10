package links_test

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"github.com/google/uuid"
	"github.com/oapi-codegen/nullable"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/link/link_scrape"
	"github.com/Southclaws/storyden/app/resources/link/link_url"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/link/fetcher"
	"github.com/Southclaws/storyden/app/services/link/scrape"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	entjob "github.com/Southclaws/storyden/internal/ent/backgroundjob"
	entlink "github.com/Southclaws/storyden/internal/ent/link"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

type scriptedScraper struct {
	mu    sync.Mutex
	fn    func(context.Context, link_url.URL) (*scrape.WebContent, error)
	calls int
}

func (s *scriptedScraper) Scrape(ctx context.Context, u link_url.URL) (*scrape.WebContent, error) {
	s.mu.Lock()
	s.calls++
	fn := s.fn
	s.mu.Unlock()
	return fn(ctx, u)
}

func (s *scriptedScraper) set(fn func(context.Context, link_url.URL) (*scrape.WebContent, error)) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.fn = fn
	return s.calls
}

func (s *scriptedScraper) count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.calls
}

func TestLinkFetchLifecycle(t *testing.T) {
	t.Parallel()
	script := &scriptedScraper{}
	var failQueries atomic.Bool
	var (
		root     context.Context
		client   *openapi.ClientWithResponses
		sessions *e2e.SessionHelper
		accounts *account_writer.Writer
		db       *ent.Client
		fetch    *fetcher.Fetcher
	)
	integration.Test(t, nil, e2e.Setup(),
		fx.Decorate(func(scrape.Scraper) scrape.Scraper { return script }),
		fx.Decorate(func(db *ent.Client) *ent.Client {
			db.Link.Intercept(ent.InterceptFunc(func(next ent.Querier) ent.Querier {
				return ent.QuerierFunc(func(ctx context.Context, q ent.Query) (ent.Value, error) {
					_, isLink := q.(*ent.LinkQuery)
					if isLink && failQueries.Load() {
						return nil, errors.New("link database unavailable")
					}

					return next.Query(ctx, q)
				})
			}))
			return db
		}),
		fx.Populate(&root, &client, &sessions, &accounts, &db, &fetch),
	)

	ctx, owner := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
	session := sessions.WithSession(ctx)

	t.Run("storage_errors_remain_server_errors", func(t *testing.T) {
		failQueries.Store(true)
		defer failQueries.Store(false)
		response, err := client.LinkCreateWithResponse(ctx, openapi.LinkCreateJSONRequestBody{
			Url: "https://" + uuid.NewString() + ".example.com/",
		}, session)
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, response.StatusCode(), string(response.Body))
	})

	t.Run("node_and_thread_reject_invalid_links_instead_of_silently_dropping_them", func(t *testing.T) {
		raw := "http://127.0.0.1/"
		node, err := client.NodeCreateWithResponse(ctx, openapi.NodeCreateJSONRequestBody{
			Name: "invalid-link-" + uuid.NewString(),
			Url:  &raw,
		}, session)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, node.StatusCode(), string(node.Body))
		thread, err := client.ThreadCreateWithResponse(ctx, openapi.ThreadCreateJSONRequestBody{
			Title: "invalid-link-" + uuid.NewString(),
			Url:   &raw,
		}, session)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, thread.StatusCode(), string(thread.Body))
	})

	t.Run("legacy_policy_violations_do_not_break_node_or_link_reads", func(t *testing.T) {
		keyword := "legacy-" + uuid.NewString()
		legacy, err := db.Link.Create().
			SetURL("http://127.0.0.1/" + keyword).
			SetSlug(keyword).
			SetDomain("127.0.0.1").
			SetTitle(keyword).
			SetDescription("").
			Save(ctx)
		require.NoError(t, err)

		created, err := client.NodeCreateWithResponse(ctx, openapi.NodeCreateJSONRequestBody{Name: keyword}, session)
		tests.Ok(t, err, created)
		nodeID, err := xid.FromString(created.JSON200.Id)
		require.NoError(t, err)
		require.NoError(t, db.Node.UpdateOneID(nodeID).SetLinkID(legacy.ID).Exec(ctx))

		stored, err := client.NodeGetWithResponse(ctx, created.JSON200.Slug, nil, session)
		tests.Ok(t, err, stored)
		require.Nil(t, stored.JSON200.Link)

		valid, err := db.Link.Create().
			SetURL("https://example.com/" + keyword).
			SetSlug(keyword + "-valid").
			SetDomain("example.com").
			SetTitle(keyword).
			SetDescription("").
			Save(ctx)
		require.NoError(t, err)
		list, err := client.LinkListWithResponse(ctx, &openapi.LinkListParams{Q: &keyword}, session)
		tests.Ok(t, err, list)
		require.Len(t, list.JSON200.Links, 1)
		require.Equal(t, valid.ID.String(), list.JSON200.Links[0].Id)

		missing, err := client.LinkGetWithResponse(ctx, legacy.Slug, session)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, missing.StatusCode())
	})

	t.Run("nodes_keep_dead_links_on_create_and_update", func(t *testing.T) {
		for _, scenario := range []struct {
			name string
			err  error
		}{
			{name: "not_found", err: errors.New("HTTP 404 Not Found")},
			{name: "dns_failure", err: &net.DNSError{Err: "no such host", Name: "dead.example.com", IsNotFound: true}},
			{name: "timeout", err: context.DeadlineExceeded},
			{name: "connection_refused", err: &net.OpError{Op: "dial", Net: "tcp", Err: errors.New("connection refused")}},
		} {
			t.Run(scenario.name, func(t *testing.T) {
				script.set(func(context.Context, link_url.URL) (*scrape.WebContent, error) {
					return nil, scenario.err
				})
				raw := "https://" + uuid.NewString() + ".example.com/dead"
				created, err := client.NodeCreateWithResponse(ctx, openapi.NodeCreateJSONRequestBody{
					Name: "dead-link-" + uuid.NewString(),
					Url:  &raw,
				}, session)
				tests.Ok(t, err, created)
				require.NotNil(t, created.JSON200.Link)
				require.Equal(t, raw, created.JSON200.Link.Url)
				require.Equal(t, openapi.LinkScrapeStatusFailed, *created.JSON200.Link.ScrapeStatus)

				replacement := "https://" + uuid.NewString() + ".example.com/also-dead"
				updated, err := client.NodeUpdateWithResponse(ctx, created.JSON200.Slug, openapi.NodeMutableProps{
					Url: nullable.NewNullableWithValue(replacement),
				}, session)
				tests.Ok(t, err, updated)
				require.NotNil(t, updated.JSON200.Link)
				require.Equal(t, replacement, updated.JSON200.Link.Url)
				require.Equal(t, openapi.LinkScrapeStatusFailed, *updated.JSON200.Link.ScrapeStatus)

				stored, err := client.NodeGetWithResponse(ctx, created.JSON200.Slug, nil, session)
				tests.Ok(t, err, stored)
				require.NotNil(t, stored.JSON200.Link)
				require.Equal(t, replacement, stored.JSON200.Link.Url)
			})
		}
	})

	t.Run("invalid_input_is_rejected_before_scraping_or_saving", func(t *testing.T) {
		calls := script.set(func(context.Context, link_url.URL) (*scrape.WebContent, error) {
			return nil, errors.New("scraper should not be called")
		})
		for _, raw := range []string{
			"", "/relative", "ftp://" + uuid.NewString() + ".example.com/", "http://127.0.0.1/",
			"http://localhost./", "https://user:secret@example.com/", "https://example.com:99999/",
			"https://example.com:443/", "http://example.com:80/", "https://example.com:8443/",
			"https://8.8.8.8/", "https://[2606:4700:4700::1111]/", "https://intranet/", "https://bad_host.example.com/",
		} {
			response, err := client.LinkCreateWithResponse(ctx, openapi.LinkCreateJSONRequestBody{Url: raw}, session)
			require.NoError(t, err)
			require.Equal(t, http.StatusBadRequest, response.StatusCode(), string(response.Body))
			var problem openapi.APIError
			require.NoError(t, json.Unmarshal(response.Body, &problem))
			require.NotNil(t, problem.Title)
			require.NotEmpty(t, *problem.Title)
			require.NotEqual(t, "Bad Request", *problem.Title, "URL policy errors should explain why the URL is rejected")
			exists, err := db.Link.Query().Where(entlink.URLEQ(raw)).Exist(ctx)
			require.NoError(t, err)
			require.False(t, exists, raw)
		}
		require.Equal(t, calls, script.count())
	})

	t.Run("blocked_resolved_target_is_not_saved_as_a_fallback", func(t *testing.T) {
		script.set(func(context.Context, link_url.URL) (*scrape.WebContent, error) {
			return nil, fault.New("resolved address is blocked", ftag.With(ftag.InvalidArgument))
		})
		raw := "https://" + uuid.NewString() + ".example.com/"
		response, err := client.LinkCreateWithResponse(ctx, openapi.LinkCreateJSONRequestBody{Url: raw}, session)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, response.StatusCode())
		exists, err := db.Link.Query().Where(entlink.URLEQ(raw)).Exist(ctx)
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("unavailable_site_is_saved_with_failed_status_and_retry_cooldown", func(t *testing.T) {
		calls := script.set(func(context.Context, link_url.URL) (*scrape.WebContent, error) {
			return nil, errors.New("site unavailable")
		})
		raw := "https://" + uuid.NewString() + ".example.com/"
		response, err := client.LinkCreateWithResponse(ctx, openapi.LinkCreateJSONRequestBody{Url: raw}, session)
		tests.Ok(t, err, response)
		require.Equal(t, raw, response.JSON200.Url)
		require.Equal(t, openapi.LinkScrapeStatusFailed, *response.JSON200.ScrapeStatus)
		require.NotNil(t, response.JSON200.ScrapeAttemptedAt)
		require.Nil(t, response.JSON200.ScrapedAt)
		again, err := client.LinkCreateWithResponse(ctx, openapi.LinkCreateJSONRequestBody{Url: raw}, session)
		tests.Ok(t, err, again)
		require.Equal(t, response.JSON200.Id, again.JSON200.Id)
		require.Equal(t, calls+1, script.count())
	})

	t.Run("redirect_destination_is_saved_deduplicated_and_survives_failed_refresh", func(t *testing.T) {
		host := uuid.NewString() + ".example.com"
		submitted := "https://www." + host + "/old"
		canonical := "https://" + host + "/new"
		finalURL, err := link_url.Parse(canonical)
		require.NoError(t, err)
		calls := script.set(func(context.Context, link_url.URL) (*scrape.WebContent, error) {
			return &scrape.WebContent{URL: finalURL, Title: "Original title", Description: "Original description"}, nil
		})
		first, err := client.LinkCreateWithResponse(ctx, openapi.LinkCreateJSONRequestBody{Url: submitted}, session)
		tests.Ok(t, err, first)
		require.Equal(t, canonical, first.JSON200.Url)
		require.Equal(t, host, first.JSON200.Domain)
		require.Equal(t, openapi.LinkScrapeStatusSuccess, *first.JSON200.ScrapeStatus)
		require.NotNil(t, first.JSON200.ScrapedAt)
		id, err := xid.FromString(first.JSON200.Id)
		require.NoError(t, err)

		for _, raw := range []string{submitted, canonical} {
			again, err := client.LinkCreateWithResponse(ctx, openapi.LinkCreateJSONRequestBody{Url: raw}, session)
			tests.Ok(t, err, again)
			require.Equal(t, first.JSON200.Id, again.JSON200.Id)
		}
		require.Equal(t, calls+2, script.count(), "the original URL must resolve again; the saved destination stays fresh")

		list, err := client.LinkListWithResponse(ctx, &openapi.LinkListParams{Q: &canonical}, session)
		tests.Ok(t, err, list)
		require.Len(t, list.JSON200.Links, 1)
		require.Equal(t, first.JSON200.Id, list.JSON200.Links[0].Id)

		list, err = client.LinkListWithResponse(ctx, &openapi.LinkListParams{Q: &submitted}, session)
		tests.Ok(t, err, list)
		require.Empty(t, list.JSON200.Links, "redirect history is not stored or searched")

		image, err := db.Asset.Create().SetAccountID(xid.ID(owner.ID)).SetFilename("retained.png").SetSize(1).SetMimeType("image/png").Save(ctx)
		require.NoError(t, err)
		lastSuccess := time.Now().Add(-25 * time.Hour).UTC().Truncate(time.Microsecond)
		_, err = db.Link.UpdateOneID(id).SetScrapedAt(lastSuccess).SetPrimaryImageID(image.ID).SetFaviconImageID(image.ID).Save(ctx)
		require.NoError(t, err)
		calls = script.set(func(context.Context, link_url.URL) (*scrape.WebContent, error) {
			return nil, errors.New("temporary outage")
		})
		_, err = fetch.Fetch(ctx, finalURL)
		require.NoError(t, err)
		require.Eventually(t, func() bool {
			stored, err := db.Link.Get(ctx, id)
			return err == nil && stored.ScrapeStatus == entlink.ScrapeStatusFailed
		}, 5*time.Second, 10*time.Millisecond)
		stored, err := db.Link.Get(ctx, id)
		require.NoError(t, err)
		require.Equal(t, entlink.ScrapeStatusFailed, stored.ScrapeStatus)
		jobs, err := db.BackgroundJob.Query().Where(entjob.KindEQ(link_scrape.Kind)).All(ctx)
		require.NoError(t, err)
		var payloads []string
		for _, job := range jobs {
			payloads = append(payloads, string(job.Payload))
		}
		require.Contains(t, payloads, `{"url":"`+canonical+`"}`, "refreshes are persisted as background jobs")
		require.Equal(t, "Original title", stored.Title)
		require.Equal(t, "Original description", stored.Description)
		require.WithinDuration(t, lastSuccess, *stored.ScrapedAt, time.Microsecond)
		require.Equal(t, image.ID, *stored.PrimaryAssetID)
		require.Equal(t, image.ID, *stored.FaviconAssetID)

		for range 3 {
			_, err := fetch.Fetch(ctx, finalURL)
			require.NoError(t, err)
		}
		require.Equal(t, calls+1, script.count(), "failed refreshes respect the retry cooldown")
		_, err = db.Link.UpdateOneID(id).SetDomain("www." + host).Save(ctx)
		require.NoError(t, err)

		calls = script.set(func(context.Context, link_url.URL) (*scrape.WebContent, error) {
			return &scrape.WebContent{Title: "Recovered title"}, nil
		})
		_, err = db.Link.UpdateOneID(id).SetScrapeAttemptedAt(time.Now().Add(-16 * time.Minute)).Save(ctx)
		require.NoError(t, err)
		_, err = fetch.Fetch(ctx, finalURL)
		require.NoError(t, err)
		require.Eventually(t, func() bool {
			stored, err := db.Link.Get(ctx, id)
			return err == nil && stored.ScrapeStatus == entlink.ScrapeStatusSuccess
		}, 5*time.Second, 10*time.Millisecond)
		stored, err = db.Link.Get(ctx, id)
		require.NoError(t, err)
		require.Equal(t, entlink.ScrapeStatusSuccess, stored.ScrapeStatus)
		require.Equal(t, "Recovered title", stored.Title)
		require.Equal(t, host, stored.Domain)
		require.Equal(t, "Original description", stored.Description, "missing refreshed fields retain good metadata")
		require.Equal(t, image.ID, *stored.PrimaryAssetID)
		require.Equal(t, calls+1, script.count())
	})
}
