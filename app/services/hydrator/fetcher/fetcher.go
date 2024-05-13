// package fetcher simply grabs URLs, scrapes them and stores them and assets.
package fetcher

import (
	"context"
	"net/http"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/gosimple/slug"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/datagraph/cluster"
	"github.com/Southclaws/storyden/app/resources/datagraph/link"
	"github.com/Southclaws/storyden/app/services/asset_manager"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/app/services/url"
)

var errEmptyLink = fault.New("empty link")

type Service interface {
	Fetch(ctx context.Context, url string) (*datagraph.Link, error)
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l      *zap.Logger
	as     asset_manager.Service
	lr     link.Repository
	sc     url.Scraper
	semdex semdex.Indexer
}

func New(
	l *zap.Logger,
	as asset_manager.Service,
	cr cluster.Repository,
	lr link.Repository,
	sc url.Scraper,
	semdex semdex.Indexer,
) Service {
	return &service{
		l:      l.With(zap.String("service", "hydrator")),
		as:     as,
		lr:     lr,
		sc:     sc,
		semdex: semdex,
	}
}

func (s *service) Fetch(ctx context.Context, url string) (*datagraph.Link, error) {
	if url == "" {
		return nil, fault.Wrap(errEmptyLink, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	r, err := s.lr.Search(ctx, 0, 1, link.WithURL(url))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if len(r.Links) > 0 {
		// revalidate stale data asynchronously
		go s.scrapeAndStore(ctx, url)
		return r.Links[0], nil
	}

	return s.scrapeAndStore(ctx, url)
}

func (s *service) scrapeAndStore(ctx context.Context, url string) (*datagraph.Link, error) {
	wc, err := s.sc.Scrape(ctx, url)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	opts := []link.Option{}

	var a *asset.Asset
	if wc.Image != "" {
		a, err = s.copy(ctx, wc.Image)
		if err != nil {
			s.l.Warn("failed to scrape web content image", zap.Error(err), zap.String("url", url))
		} else {
			opts = append(opts, link.WithAssets(a.ID))
		}
	}

	ln, err := s.lr.Store(ctx, url, wc.Title, wc.Description, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if a != nil {
		ln.Assets = append(ln.Assets, a)
	}

	if err := s.semdex.Index(ctx, ln); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to index thread"))
	}

	return ln, nil
}

func (s *service) copy(ctx context.Context, url string) (*asset.Asset, error) {
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

	a, err := s.as.Upload(ctx, resp.Body, resp.ContentLength, asset.NewFilename(name), url)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return a, nil
}
