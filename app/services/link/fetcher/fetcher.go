// package fetcher simply grabs URLs, scrapes them and stores them and assets.
package fetcher

import (
	"context"
	"net/http"
	"net/url"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/gosimple/slug"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/link/link_querier"
	"github.com/Southclaws/storyden/app/resources/link/link_ref"
	"github.com/Southclaws/storyden/app/resources/link/link_writer"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/services/asset/asset_upload"
	"github.com/Southclaws/storyden/app/services/link/scrape"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

var errEmptyLink = fault.New("empty link")

type Fetcher struct {
	l        *zap.Logger
	uploader *asset_upload.Uploader
	lq       *link_querier.LinkQuerier
	lr       *link_writer.LinkWriter
	sc       scrape.Scraper
	queue    pubsub.Topic[mq.ScrapeLink]
}

func New(
	l *zap.Logger,
	uploader *asset_upload.Uploader,
	nr library.Repository,
	lq *link_querier.LinkQuerier,
	lr *link_writer.LinkWriter,
	sc scrape.Scraper,
	queue pubsub.Topic[mq.ScrapeLink],
) *Fetcher {
	return &Fetcher{
		l:        l.With(zap.String("service", "hydrator")),
		uploader: uploader,
		lq:       lq,
		lr:       lr,
		sc:       sc,
		queue:    queue,
	}
}

func (s *Fetcher) Fetch(ctx context.Context, u url.URL) (*link_ref.LinkRef, error) {
	if u.String() == "" {
		return nil, fault.Wrap(errEmptyLink, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	r, err := s.lq.Search(ctx, 0, 1, link_querier.WithURL(u.String()))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if len(r.Links) > 0 {
		if err := s.queue.Publish(ctx, mq.ScrapeLink{URL: u}); err != nil {
			s.l.Error("failed to publish scrape link message",
				zap.Error(err),
				zap.String("url", u.String()))
		}

		return r.Links[0], nil
	}

	return s.ScrapeAndStore(ctx, u)
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

		err = s.QueueForItem(ctx, *parsed, item)
		if err != nil {
			continue
		}
	}
}

// QueueForItem queues a scrape request for a URL that is linked to an item.
// When the scrape job is done, the scraped link will be related to the item.
func (s *Fetcher) QueueForItem(ctx context.Context, u url.URL, item datagraph.Item) error {
	if err := s.queue.Publish(ctx, mq.ScrapeLink{
		URL:  u,
		Item: datagraph.NewRef(item),
	}); err != nil {
		s.l.Error("failed to publish scrape link message",
			zap.Error(err),
			zap.String("url", u.String()))
	}

	return nil
}

func (s *Fetcher) ScrapeAndStore(ctx context.Context, u url.URL) (*link_ref.LinkRef, error) {
	wc, err := s.sc.Scrape(ctx, u)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	opts := []link_writer.Option{}

	if wc.Favicon != "" {
		a, err := s.CopyAsset(ctx, wc.Favicon)
		if err != nil {
			s.l.Warn("failed to scrape web content favicon image", zap.Error(err), zap.String("url", u.String()))
		} else {
			opts = append(opts, link_writer.WithFaviconImage(a.ID))
		}
	}

	if wc.Image != "" {
		a, err := s.CopyAsset(ctx, wc.Image)
		if err != nil {
			s.l.Warn("failed to scrape web content primary image", zap.Error(err), zap.String("url", u.String()))
		} else {
			opts = append(opts, link_writer.WithPrimaryImage(a.ID))
		}
	}

	ln, err := s.lr.Store(ctx, u.String(), wc.Title, wc.Description, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return ln, nil
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
	name := slug.Make(url)

	a, err := s.uploader.Upload(ctx, resp.Body, resp.ContentLength, asset.NewFilename(name), asset_upload.Options{})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return a, nil
}
