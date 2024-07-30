package hydrator

import (
	"context"
	"net/url"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/content"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/services/link/fetcher"
)

type Service interface {
	HydrateThread(ctx context.Context, structured content.Rich, url opt.Optional[string]) []thread.Option
	HydrateReply(ctx context.Context, structured content.Rich, url opt.Optional[string]) []reply.Option
	HydrateNode(ctx context.Context, structured content.Rich, url opt.Optional[string]) []library.Option
}

type service struct {
	l  *zap.Logger
	tr thread.Repository
	nr library.Repository
	f  fetcher.Fetcher
}

func New(
	l *zap.Logger,
	tr thread.Repository,
	nr library.Repository,
	f fetcher.Fetcher,
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

func (s *service) HydrateNode(ctx context.Context, structured content.Rich, url opt.Optional[string]) []library.Option {
	links, assets := s.hydrate(ctx, structured, url)

	return []library.Option{
		library.WithAssets(assets),
		library.WithLinks(links...),
	}
}

func (s *service) hydrate(ctx context.Context, structured content.Rich, optionalURL opt.Optional[string]) ([]xid.ID, []asset.AssetID) {
	urls := []string{}

	if u, ok := optionalURL.Get(); ok {
		urls = append(urls, u)
	}

	urls = append(urls, structured.Links()...)

	links := []xid.ID{}
	assets := []asset.AssetID{}

	for _, l := range urls {
		parsed, err := url.Parse(l)
		if err != nil {
			// TODO: Handle this in queue consumer later
			continue
		}

		ln, err := s.f.Fetch(ctx, *parsed)
		if err != nil {
			continue
		}

		links = append(links, xid.ID(ln.ID))
		assets = append(assets, ln.AssetIDs()...)
	}

	return links, assets
}
