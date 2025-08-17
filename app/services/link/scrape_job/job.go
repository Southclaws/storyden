package scrape_job

import (
	"context"
	"log/slog"
	"net/url"

	"go.uber.org/fx"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/post_writer"
	"github.com/Southclaws/storyden/app/services/link/fetcher"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

func Build() fx.Option {
	return fx.Invoke(runScrapeConsumer)
}

type scrapeConsumer struct {
	fetcher    *fetcher.Fetcher
	posts      *post_writer.PostWriter
	nodeWriter *node_writer.Writer
}

func runScrapeConsumer(
	ctx context.Context,
	lc fx.Lifecycle,
	logger *slog.Logger,
	bus *pubsub.Bus,
	fetcher *fetcher.Fetcher,
	posts *post_writer.PostWriter,
	nodeWriter *node_writer.Writer,
) {
	ic := scrapeConsumer{
		fetcher:    fetcher,
		posts:      posts,
		nodeWriter: nodeWriter,
	}

	lc.Append(fx.StartHook(func(hctx context.Context) error {
		_, err := pubsub.SubscribeCommand(hctx, bus, "scrape_job.scrape", func(ctx context.Context, cmd *message.CommandScrapeLink) error {
			return ic.scrapeLink(ctx, cmd.URL, opt.NewPtr(cmd.Item))
		})
		if err != nil {
			return err
		}

		return nil
	}))
}

func (s *scrapeConsumer) scrapeLink(ctx context.Context, u url.URL, item opt.Optional[datagraph.Ref]) error {
	ln, _, err := s.fetcher.ScrapeAndStore(ctx, u)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if i, ok := item.Get(); ok {
		switch i.Kind {
		case datagraph.KindPost:
			_, err := s.posts.Update(ctx, post.ID(i.ID), post_writer.WithContentLinks(xid.ID(ln.ID)))
			if err != nil {
				return fault.Wrap(err, fctx.With(ctx))
			}

		case datagraph.KindNode:
			qk := library.QueryKey{mark.NewQueryKeyID(i.ID)}
			_, err := s.nodeWriter.Update(ctx, qk, node_writer.WithContentLinks(xid.ID(ln.ID)))
			if err != nil {
				return fault.Wrap(err, fctx.With(ctx))
			}
		}
	}

	return nil
}
