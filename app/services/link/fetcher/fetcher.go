// package fetcher simply grabs URLs, scrapes them and stores them and assets.
package fetcher

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/link/link_querier"
	"github.com/Southclaws/storyden/app/resources/link/link_ref"
	"github.com/Southclaws/storyden/app/resources/link/link_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/services/asset/asset_upload"
	"github.com/Southclaws/storyden/app/services/link/scrape"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

var errEmptyLink = fault.New("empty link")

type Fetcher struct {
	logger   *slog.Logger
	uploader *asset_upload.Uploader
	lq       *link_querier.LinkQuerier
	lr       *link_writer.LinkWriter
	sc       scrape.Scraper
	bus      *pubsub.Bus
}

func New(
	logger *slog.Logger,
	uploader *asset_upload.Uploader,
	lq *link_querier.LinkQuerier,
	lr *link_writer.LinkWriter,
	sc scrape.Scraper,
	bus *pubsub.Bus,
) *Fetcher {
	return &Fetcher{
		logger:   logger,
		uploader: uploader,
		lq:       lq,
		lr:       lr,
		sc:       sc,
		bus:      bus,
	}
}

type Options struct {
	ContentFill opt.Optional[asset.ContentFillCommand]
}

func (s *Fetcher) Fetch(ctx context.Context, u url.URL, opts Options) (*link_ref.LinkRef, error) {
	if u.String() == "" {
		return nil, fault.Wrap(errEmptyLink, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	r, err := s.lq.Search(ctx, 0, 1, link_querier.WithURL(u.String()))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if len(r.Links) > 0 {
		if err := s.bus.SendCommand(ctx, &message.CommandScrapeLink{URL: u}); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		return r.Links[0], nil
	}

	lr, _, err := s.ScrapeAndStore(ctx, u)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return lr, nil
}

// HydrateContentURLs takes all the URLs mentioned in the content of an item and
// queues them for hydration. This process visits each URL and fetches metadata.
// Then, stores that metadata in the database and relates them back to the item.
func (s *Fetcher) HydrateContentURLs(ctx context.Context, item datagraph.Item) {
	urls := item.GetContent().Links()

	for _, l := range urls {
		parsed, err := url.Parse(l)
		if err != nil {
			continue
		}

		err = s.queueForItem(ctx, *parsed, item)
		if err != nil {
			continue
		}
	}
}

// queueForItem queues a scrape request for a URL that is linked to an item.
// When the scrape job is done, the scraped link will be related to the item.
func (s *Fetcher) queueForItem(ctx context.Context, u url.URL, item datagraph.Item) error {
	if err := s.bus.SendCommand(ctx, &message.CommandScrapeLink{
		URL:  u,
		Item: datagraph.NewRef(item),
	}); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (s *Fetcher) ScrapeAndStore(ctx context.Context, u url.URL) (*link_ref.LinkRef, *scrape.WebContent, error) {
	wc, err := s.sc.Scrape(ctx, u)
	if err != nil {
		s.logger.Warn("failed to scrape URL, storing link with basic information only",
			slog.String("error", err.Error()),
			slog.String("url", u.String()))

		ln, err := s.lr.Store(ctx, u.String(), "", "", []link_writer.Option{}...)
		if err != nil {
			return nil, nil, fault.Wrap(err, fctx.With(ctx))
		}

		return ln, &scrape.WebContent{}, nil
	}

	opts := []link_writer.Option{}

	if wc.Favicon != "" {
		a, err := s.CopyAsset(ctx, wc.Favicon)
		if err != nil {
			s.logger.Warn("failed to scrape web content favicon image", slog.String("error", err.Error()), slog.String("url", u.String()))
		} else {
			opts = append(opts, link_writer.WithFaviconImage(a.ID))
		}
	}

	if wc.Image != "" {
		a, err := s.CopyAsset(ctx, wc.Image)
		if err != nil {
			s.logger.Warn("failed to scrape web content primary image", slog.String("error", err.Error()), slog.String("url", u.String()))
		} else {
			opts = append(opts, link_writer.WithPrimaryImage(a.ID))
		}
	}

	ln, err := s.lr.Store(ctx, u.String(), wc.Title, wc.Description, opts...)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	return ln, wc, nil
}

func (s *Fetcher) CopyAsset(ctx context.Context, url string) (*asset.Asset, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if resp.StatusCode != http.StatusOK {
		ctx = fctx.WithMeta(ctx, "status", resp.Status)
		return nil, fault.Wrap(fault.New("failed to get"), fctx.With(ctx))
	}

	// TODO: Better naming???
	name := mark.Slugify(url)

	a, err := s.uploader.Upload(ctx, resp.Body, resp.ContentLength, asset.NewFilename(name), asset_upload.Options{})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return a, nil
}
