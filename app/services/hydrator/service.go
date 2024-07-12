package hydrator

import (
	"context"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/content"
	"github.com/Southclaws/storyden/app/resources/datagraph/node"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/services/hydrator/fetcher"
)

type Service interface {
	HydrateThread(ctx context.Context, structured content.Rich, url opt.Optional[string]) []thread.Option
	HydrateReply(ctx context.Context, structured content.Rich, url opt.Optional[string]) []reply.Option
	HydrateNode(ctx context.Context, structured content.Rich, url opt.Optional[string]) []node.Option
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l  *zap.Logger
	tr thread.Repository
	nr node.Repository
	f  fetcher.Service
}

func New(
	l *zap.Logger,
	tr thread.Repository,
	nr node.Repository,
	f fetcher.Service,
) Service {
	return &service{
		l:  l.With(zap.String("service", "hydrator")),
		tr: tr,
		nr: nr,
		f:  f,
	}
}

func (s *service) HydrateThread(ctx context.Context, structured content.Rich, url opt.Optional[string]) []thread.Option {
	links, assets := s.hydrate(ctx, structured, url)

	return []thread.Option{
		thread.WithAssets(assets),
		thread.WithLinks(links...),
	}
}

func (s *service) HydrateReply(ctx context.Context, structured content.Rich, url opt.Optional[string]) []reply.Option {
	links, assets := s.hydrate(ctx, structured, url)

	return []reply.Option{
		reply.WithAssets(assets...),
		reply.WithLinks(links...),
	}
}

func (s *service) HydrateNode(ctx context.Context, structured content.Rich, url opt.Optional[string]) []node.Option {
	links, assets := s.hydrate(ctx, structured, url)

	return []node.Option{
		node.WithAssets(assets),
		node.WithLinks(links...),
	}
}

// hydrate takes the body and primary URL of a piece of content and fetches all
// the links and produces a short summary of the post's body text.
func (s *service) hydrate(ctx context.Context, structured content.Rich, urls opt.Optional[string]) ([]xid.ID, []asset.AssetID) {
	urls = append(urls, structured.Links()...)

	links := []xid.ID{}
	assets := []asset.AssetID{}

	for _, l := range urls {
		// TODO: async

		ln, err := s.f.Fetch(ctx, l)
		if err != nil {
			continue
		}

		links = append(links, xid.ID(ln.ID))
		assets = append(assets, ln.AssetIDs()...)
	}

	return links, assets
}
