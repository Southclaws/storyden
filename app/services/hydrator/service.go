package hydrator

import (
	"context"
	"net/http"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/cluster"
	"github.com/Southclaws/storyden/app/resources/item"
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
	as asset_svc.Service
	tr thread.Repository
	cr cluster.Repository
	ir item.Repository
	sc url.Scraper
}

func New(
	as asset_svc.Service,
	tr thread.Repository,
	cr cluster.Repository,
	ir item.Repository,
	sc url.Scraper,
) Service {
	return &service{
		as,
		tr,
		cr,
		ir,
		sc,
	}
}

func (s *service) HydrateThread(ctx context.Context, url string) ([]thread.Option, error) {
	a, wc, err := s.scrape(ctx, url)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return []thread.Option{
		thread.WithAssets([]asset.AssetID{a.ID}),
		thread.WithLink(url, wc.Title, wc.Description),
	}, nil
}

func (s *service) HydrateCluster(ctx context.Context, url string) ([]cluster.Option, error) {
	a, wc, err := s.scrape(ctx, url)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return []cluster.Option{
		cluster.WithAssets([]asset.AssetID{a.ID}),
		cluster.WithLink(url, wc.Title, wc.Description),
	}, nil
}

func (s *service) HydrateItem(ctx context.Context, url string) ([]item.Option, error) {
	a, wc, err := s.scrape(ctx, url)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return []item.Option{
		item.WithAssets([]asset.AssetID{a.ID}),
		item.WithLink(url, wc.Title, wc.Description),
	}, nil
}

func (s *service) scrape(ctx context.Context, url string) (*asset.Asset, *url.WebContent, error) {
	wc, err := s.sc.Scrape(ctx, url)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	if wc.Image == "" {
		return nil, nil, nil
	}

	a, err := s.copy(ctx, wc.Image)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	return a, wc, nil
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
