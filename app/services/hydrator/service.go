package hydrator

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/cluster"
	"github.com/Southclaws/storyden/app/resources/item"
	"github.com/Southclaws/storyden/app/resources/thread"
	"github.com/Southclaws/storyden/app/services/hydrator/fetcher"
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
	tr thread.Repository
	cr cluster.Repository
	ir item.Repository
	f  fetcher.Service
}

func New(
	l *zap.Logger,
	tr thread.Repository,
	cr cluster.Repository,
	ir item.Repository,
	f fetcher.Service,
) Service {
	return &service{
		l:  l.With(zap.String("service", "hydrator")),
		tr: tr,
		cr: cr,
		ir: ir,
		f:  f,
	}
}

func (s *service) HydrateThread(ctx context.Context, url string) ([]thread.Option, error) {
	ln, err := s.f.Fetch(ctx, url)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return []thread.Option{
		thread.WithAssets(ln.AssetIDs()),
		thread.WithLinks(ln.ID),
	}, nil
}

func (s *service) HydrateCluster(ctx context.Context, url string) ([]cluster.Option, error) {
	ln, err := s.f.Fetch(ctx, url)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return []cluster.Option{
		cluster.WithAssets(ln.AssetIDs()),
		cluster.WithLinks(ln.ID),
	}, nil
}

func (s *service) HydrateItem(ctx context.Context, url string) ([]item.Option, error) {
	ln, err := s.f.Fetch(ctx, url)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return []item.Option{
		item.WithAssets(ln.AssetIDs()),
		item.WithLinks(ln.ID),
	}, nil
}
