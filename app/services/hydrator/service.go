package hydrator

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
	"github.com/Southclaws/storyden/app/resources/thread"
	asset_svc "github.com/Southclaws/storyden/app/services/asset"
	"github.com/Southclaws/storyden/app/services/url"
)

type Service interface {
	HydrateThread(ctx context.Context, url string) ([]thread.Option, error)
	HydrateCluster(ctx context.Context, url string) ([]cluster.Option, error)
	HydrateItem(ctx context.Context, url string) ([]item.Option, error)
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l  *zap.Logger
	as asset_svc.Service
	tr thread.Repository
	cr cluster.Repository
	ir item.Repository
	lr link.Repository
	sc url.Scraper
}

func New(
	l *zap.Logger,
	as asset_svc.Service,
	tr thread.Repository,
	cr cluster.Repository,
	ir item.Repository,
	lr link.Repository,
	sc url.Scraper,
) Service {
	return &service{
		l:  l.With(zap.String("service", "hydrator")),
		as: as,
		tr: tr,
		cr: cr,
		ir: ir,
		lr: lr,
		sc: sc,
	}
}

func (s *service) HydrateThread(ctx context.Context, url string) ([]thread.Option, error) {
	ln, err := s.scrape(ctx, url)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return []thread.Option{
		thread.WithAssets(ln.AssetIDs()),
		thread.WithLinks(ln.ID),
	}, nil
}

func (s *service) HydrateCluster(ctx context.Context, url string) ([]cluster.Option, error) {
	ln, err := s.scrape(ctx, url)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return []cluster.Option{
		cluster.WithAssets(ln.AssetIDs()),
		cluster.WithLinks(ln.ID),
	}, nil
}

func (s *service) HydrateItem(ctx context.Context, url string) ([]item.Option, error) {
	ln, err := s.scrape(ctx, url)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return []item.Option{
		item.WithAssets(ln.AssetIDs()),
		item.WithLinks(ln.ID),
	}, nil
}

func (s *service) scrape(ctx context.Context, url string) (*link.Link, error) {
	wc, err := s.sc.Scrape(ctx, url)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	opts := []link.Option{}

	if wc.Image != "" {
		a, err := s.copy(ctx, wc.Image)
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
