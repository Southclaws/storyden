// package fetcher simply grabs URLs, scrapes them and stores them and assets.
package fetcher

import (
	"context"
	"net/http"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/cluster"
	"github.com/Southclaws/storyden/app/resources/item"
	"github.com/Southclaws/storyden/app/resources/link"
	asset_svc "github.com/Southclaws/storyden/app/services/asset"
	"github.com/Southclaws/storyden/app/services/url"
)

type Service interface {
	Fetch(ctx context.Context, url string) (*link.Link, error)
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l  *zap.Logger
	as asset_svc.Service
	lr link.Repository
	sc url.Scraper
}

func New(
	l *zap.Logger,
	as asset_svc.Service,
	cr cluster.Repository,
	ir item.Repository,
	lr link.Repository,
	sc url.Scraper,
) Service {
	return &service{
		l:  l.With(zap.String("service", "hydrator")),
		as: as,
		lr: lr,
		sc: sc,
	}
}

func (s *service) Fetch(ctx context.Context, url string) (*link.Link, error) {
	r, err := s.lr.Search(ctx, link.WithURL(url))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if len(r) > 0 {
		// TODO: revalidate stale link async
		return r[0], nil
	}

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
			opts = append(opts, link.WithAssets(string(a.ID)))
		}
	}

	ln, err := s.lr.Store(ctx, url, wc.Title, wc.Description, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if a != nil {
		ln.Assets = append(ln.Assets, a)
	}

	return ln, nil
}

func (s *service) copy(ctx context.Context, url string) (*asset.Asset, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	a, err := s.as.Upload(ctx, resp.Body, resp.ContentLength)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return a, nil
}
