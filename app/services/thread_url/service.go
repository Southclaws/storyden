package thread_url

import (
	"context"
	"net/http"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/thread"
	asset_svc "github.com/Southclaws/storyden/app/services/asset"
	"github.com/Southclaws/storyden/app/services/url"
)

type Service interface {
	Hydrate(ctx context.Context, thr *thread.Thread) (*thread.Thread, error)
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	as asset_svc.Service
	tr thread.Repository
	sc url.Scraper
}

func New(
	as asset_svc.Service,
	tr thread.Repository,
	sc url.Scraper,
) Service {
	return &service{
		as,
		tr,
		sc,
	}
}

func (s *service) Hydrate(ctx context.Context, thr *thread.Thread) (*thread.Thread, error) {
	u, ok := thr.URL.Get()
	if !ok {
		return thr, nil
	}

	wc, err := s.sc.Scrape(ctx, u)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if wc.Image == "" {
		return nil, nil
	}

	a, err := s.copy(ctx, wc.Image)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	thr, err = s.tr.Update(ctx, thr.ID, thread.WithAssets([]asset.AssetID{a.ID}))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return thr, nil
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
