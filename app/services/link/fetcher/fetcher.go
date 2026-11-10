// package fetcher simply grabs URLs, scrapes them and stores them and assets.
package fetcher

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/background_job"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/link"
	"github.com/Southclaws/storyden/app/resources/link/link_querier"
	"github.com/Southclaws/storyden/app/resources/link/link_ref"
	"github.com/Southclaws/storyden/app/resources/link/link_scrape"
	"github.com/Southclaws/storyden/app/resources/link/link_url"
	"github.com/Southclaws/storyden/app/resources/link/link_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/services/asset/asset_upload"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/background_job/job_dispatcher"
	"github.com/Southclaws/storyden/app/services/link/scrape"
	"github.com/Southclaws/storyden/internal/infrastructure/httpsafe"
)

var (
	errAssetTooLarge = fault.New("linked asset exceeds maximum size")
)

const (
	assetFetchTimeout  = 30 * time.Second
	assetFetchMaxBytes = 8 * 1024 * 1024
)

type Fetcher struct {
	logger     *slog.Logger
	uploader   *asset_upload.Uploader
	lq         *link_querier.LinkQuerier
	lr         *link_writer.LinkWriter
	sc         scrape.Scraper
	jobs       *background_job.Repository
	dispatcher job_dispatcher.Dispatcher
	client     *http.Client
}

func New(
	logger *slog.Logger,
	uploader *asset_upload.Uploader,
	lq *link_querier.LinkQuerier,
	lr *link_writer.LinkWriter,
	sc scrape.Scraper,
	jobs *background_job.Repository,
	dispatcher job_dispatcher.Dispatcher,
) (*Fetcher, error) {
	client, err := httpsafe.NewSurfClient(httpsafe.Config{Timeout: assetFetchTimeout})
	if err != nil {
		return nil, err
	}

	return &Fetcher{
		logger:     logger,
		uploader:   uploader,
		lq:         lq,
		lr:         lr,
		sc:         sc,
		jobs:       jobs,
		dispatcher: dispatcher,
		client:     client,
	}, nil
}

const (
	metadataFreshness  = 24 * time.Hour
	metadataRetryDelay = 15 * time.Minute
)

func (s *Fetcher) Fetch(ctx context.Context, u link_url.URL) (*link_ref.LinkRef, error) {
	if u.IsZero() {
		return nil, link_url.ErrEmpty
	}

	existing, err := s.lq.FindByURL(ctx, u)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		ln, _, err := s.ScrapeAndStore(ctx, u)
		return ln, err
	}

	if needsRefresh(existing, time.Now()) {
		if err := s.enqueue(ctx, link_scrape.Payload{URL: existing.URL}); err != nil {
			return nil, err
		}
	}

	return existing, nil
}

func needsRefresh(existing *link_ref.LinkRef, now time.Time) bool {
	if existing.ScrapeAttemptedAt == nil {
		return true
	}

	if existing.ScrapeStatus == "failed" {
		return existing.ScrapeAttemptedAt.Before(now.Add(-metadataRetryDelay))
	}

	return existing.ScrapedAt == nil || existing.ScrapedAt.Before(now.Add(-metadataFreshness))
}

func (s *Fetcher) enqueue(ctx context.Context, payload link_scrape.Payload) error {
	id := xid.New()
	if err := s.jobs.Ensure(ctx, id, link_scrape.Kind, id.String(), payload); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	s.dispatcher.Dispatch(ctx, id)
	return nil
}

// HydrateContentURLs takes all the URLs mentioned in the content of an item and
// queues them for hydration. This process visits each URL and fetches metadata.
// Then, stores that metadata in the database and relates them back to the item.
func (s *Fetcher) HydrateContentURLs(ctx context.Context, item datagraph.Item) {
	urls := item.GetContent().Links()

	for _, l := range urls {
		parsed, err := link_url.Parse(l)
		if err != nil {
			continue
		}

		err = s.queueForItem(ctx, parsed, item)
		if err != nil {
			continue
		}
	}
}

// queueForItem queues a scrape request for a URL that is linked to an item.
// When the scrape job is done, the scraped link will be related to the item.
func (s *Fetcher) queueForItem(ctx context.Context, u link_url.URL, item datagraph.Item) error {
	return s.enqueue(ctx, link_scrape.Payload{URL: u, Item: datagraph.NewRef(item)})
}

func (s *Fetcher) ScrapeAndStore(ctx context.Context, u link_url.URL) (*link_ref.LinkRef, *scrape.WebContent, error) {
	if u.IsZero() {
		return nil, nil, link_url.ErrEmpty
	}

	existing, err := s.lq.FindByURL(ctx, u)
	if err != nil {
		return nil, nil, err
	}

	wc, scrapeErr := s.sc.Scrape(ctx, u)
	if scrapeErr != nil {
		if ftag.Get(scrapeErr) == ftag.InvalidArgument || ctx.Err() != nil {
			return nil, nil, scrapeErr
		}

		s.logger.Warn("failed to scrape URL, retaining existing metadata", slog.String("error", scrapeErr.Error()), slog.String("url", u.String()))
	}

	if wc == nil {
		wc = &scrape.WebContent{}
	}

	target := u
	if !wc.URL.IsZero() {
		target = wc.URL
	}

	if existing == nil {
		existing, err = s.lr.Store(ctx, target, "", "")
		if err != nil {
			return nil, nil, err
		}
	}

	result := link_writer.ScrapeResult{Success: scrapeErr == nil, Title: wc.Title, Description: wc.Description}
	if result.Success {
		result.Domain = existing.URL.Hostname()
		result.Options = s.copyMetadataAssets(ctx, wc)
	}

	if err := s.lr.FinishScrape(ctx, existing.ID, time.Now(), result); err != nil {
		return nil, nil, err
	}

	ln, err := s.lq.GetByID(ctx, link.LinkID(existing.ID))
	return ln, wc, err
}

func (s *Fetcher) copyMetadataAssets(ctx context.Context, wc *scrape.WebContent) []link_writer.Option {
	var opts []link_writer.Option
	for _, image := range []struct {
		url    string
		option func(asset.AssetID) link_writer.Option
	}{
		{wc.Favicon, link_writer.WithFaviconImage},
		{wc.Image, link_writer.WithPrimaryImage},
	} {
		if image.url == "" {
			continue
		}

		a, err := s.CopyAsset(ctx, image.url)
		if err != nil {
			s.logger.Warn("failed to scrape metadata image", slog.String("url", image.url), slog.String("error", err.Error()))
			continue
		}

		opts = append(opts, image.option(a.ID))
	}

	return opts
}

func (s *Fetcher) CopyAsset(ctx context.Context, url string) (*asset.Asset, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	body, err := fetchBounded(ctx, s.client, url, assetFetchMaxBytes)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// TODO: Better naming???
	name := mark.Slugify(url)

	a, err := s.uploader.Upload(
		ctx,
		accountID,
		bytes.NewReader(body),
		int64(len(body)),
		asset.NewFilename(name),
		asset_upload.Options{},
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return a, nil
}

func fetchBounded(ctx context.Context, client *http.Client, url string, maxBytes int64) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ctx = fctx.WithMeta(ctx, "status", resp.Status)
		return nil, fault.Wrap(fault.New("failed to get"), fctx.With(ctx))
	}

	if resp.ContentLength > maxBytes {
		return nil, fault.Wrap(errAssetTooLarge, fctx.With(ctx))
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBytes+1))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if int64(len(body)) > maxBytes {
		return nil, fault.Wrap(errAssetTooLarge, fctx.With(ctx))
	}

	return body, nil
}
