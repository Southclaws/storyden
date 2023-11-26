package hydrator

import (
	"context"

	"github.com/rs/xid"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/cluster"
	"github.com/Southclaws/storyden/app/resources/item"
	"github.com/Southclaws/storyden/app/resources/reply"
	"github.com/Southclaws/storyden/app/resources/thread"
	"github.com/Southclaws/storyden/app/services/hydrator/extractor"
	"github.com/Southclaws/storyden/app/services/hydrator/fetcher"
)

type Service interface {
	HydrateThread(ctx context.Context, body, url string) []thread.Option
	HydrateReply(ctx context.Context, body, url string) []reply.Option
	HydrateCluster(ctx context.Context, body, url string) []cluster.Option
	HydrateItem(ctx context.Context, body, url string) []item.Option
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

func (s *service) HydrateThread(ctx context.Context, body, url string) []thread.Option {
	short, links, assets := s.hydrate(ctx, body, url)

	return []thread.Option{
		thread.WithAssets(assets),
		thread.WithLinks(links...),
		thread.WithSummary(short),
	}
}

func (s *service) HydrateReply(ctx context.Context, body, url string) []reply.Option {
	short, links, assets := s.hydrate(ctx, body, url)

	return []reply.Option{
		reply.WithAssets(assets...),
		reply.WithShort(short),
		reply.WithLinks(links...),
	}
}

func (s *service) HydrateCluster(ctx context.Context, body, url string) []cluster.Option {
	_, links, assets := s.hydrate(ctx, body, url)

	return []cluster.Option{
		cluster.WithAssets(assets),
		cluster.WithLinks(links...),
	}
}

func (s *service) HydrateItem(ctx context.Context, body, url string) []item.Option {
	_, links, assets := s.hydrate(ctx, body, url)

	return []item.Option{
		item.WithAssets(assets),
		item.WithLinks(links...),
	}
}

// hydrate takes the body and primary URL of a piece of content and fetches all
// the links and produces a short summary of the post's body text.
func (s *service) hydrate(ctx context.Context, body, url string) (string, []xid.ID, []asset.AssetID) {
	structured := extractor.Destructure(body)

	urls := append([]string{url}, structured.Links...)

	links := []xid.ID{}
	assets := []asset.AssetID{}

	for _, l := range urls {
		// TODO: async

		ln, err := s.f.Fetch(ctx, l)
		if err != nil {
			continue
		}

		links = append(links, ln.ID)
		assets = append(assets, ln.AssetIDs()...)
	}

	return structured.Short, links, assets
}
