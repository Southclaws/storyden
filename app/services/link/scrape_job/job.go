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

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/post_writer"
	"github.com/Southclaws/storyden/app/resources/post/reply_querier"
	"github.com/Southclaws/storyden/app/resources/post/thread_querier"
	"github.com/Southclaws/storyden/app/services/link/fetcher"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func Build() fx.Option {
	return fx.Invoke(runScrapeConsumer)
}

type scrapeConsumer struct {
	fetcher     *fetcher.Fetcher
	postWriter  *post_writer.PostWriter
	postQuery   *reply_querier.Querier
	nodeWriter  *node_writer.Writer
	threadQuery *thread_querier.Querier
	nodeQuery   *node_querier.Querier
	bus         *pubsub.Bus
}

func runScrapeConsumer(
	ctx context.Context,
	lc fx.Lifecycle,
	logger *slog.Logger,
	bus *pubsub.Bus,
	fetcher *fetcher.Fetcher,
	postWriter *post_writer.PostWriter,
	postQuery *reply_querier.Querier,
	nodeWriter *node_writer.Writer,
	threadQuery *thread_querier.Querier,
	nodeQuery *node_querier.Querier,
) {
	ic := scrapeConsumer{
		fetcher:     fetcher,
		postWriter:  postWriter,
		postQuery:   postQuery,
		nodeWriter:  nodeWriter,
		threadQuery: threadQuery,
		nodeQuery:   nodeQuery,
		bus:         bus,
	}

	lc.Append(fx.StartHook(func(hctx context.Context) error {
		// Subscribe to scrape commands
		_, err := pubsub.SubscribeCommand(ctx, bus, "scrape_job.scrape", func(ctx context.Context, cmd *message.CommandScrapeLink) error {
			return ic.scrapeLink(ctx, cmd.URL, opt.NewPtr(cmd.Item))
		})
		if err != nil {
			return err
		}

		// Subscribe to thread events for URL hydration
		_, err = pubsub.Subscribe(ctx, bus, "scrape_job.hydrate_thread_created", func(ctx context.Context, evt *rpc.EventThreadPublished) error {
			return ic.hydrateThreadURLs(ctx, evt.ID)
		})
		if err != nil {
			return err
		}

		_, err = pubsub.Subscribe(ctx, bus, "scrape_job.hydrate_thread_updated", func(ctx context.Context, evt *rpc.EventThreadUpdated) error {
			return ic.hydrateThreadURLs(ctx, evt.ID)
		})
		if err != nil {
			return err
		}

		// Subscribe to reply events for URL hydration
		_, err = pubsub.Subscribe(ctx, bus, "scrape_job.hydrate_reply_created", func(ctx context.Context, evt *rpc.EventThreadReplyCreated) error {
			return ic.hydratePostURLs(ctx, evt.ReplyID)
		})
		if err != nil {
			return err
		}

		_, err = pubsub.Subscribe(ctx, bus, "scrape_job.hydrate_reply_updated", func(ctx context.Context, evt *rpc.EventThreadReplyUpdated) error {
			return ic.hydratePostURLs(ctx, evt.ReplyID)
		})
		if err != nil {
			return err
		}

		// Subscribe to node events for URL hydration
		_, err = pubsub.Subscribe(ctx, bus, "scrape_job.hydrate_node_created", func(ctx context.Context, evt *rpc.EventNodeCreated) error {
			return ic.hydrateNodeURLs(ctx, evt.ID)
		})
		if err != nil {
			return err
		}

		_, err = pubsub.Subscribe(ctx, bus, "scrape_job.hydrate_node_updated", func(ctx context.Context, evt *rpc.EventNodeUpdated) error {
			return ic.hydrateNodeURLs(ctx, evt.ID)
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
			_, err := s.postWriter.Update(ctx, post.ID(i.ID), post_writer.WithContentLinks(xid.ID(ln.ID)))
			if err != nil {
				return fault.Wrap(err, fctx.With(ctx))
			}

		case datagraph.KindNode:
			qk := library.NewID(i.ID)
			_, err := s.nodeWriter.Update(ctx, qk, node_writer.WithContentLinks(xid.ID(ln.ID)))
			if err != nil {
				return fault.Wrap(err, fctx.With(ctx))
			}
		}
	}

	return nil
}

func (s *scrapeConsumer) hydrateThreadURLs(ctx context.Context, threadID post.ID) error {
	thread, err := s.threadQuery.Get(ctx, threadID, pagination.NewPageParams(1, 1), opt.NewEmpty[account.AccountID]())
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	s.fetcher.HydrateContentURLs(ctx, thread)
	return nil
}

func (s *scrapeConsumer) hydratePostURLs(ctx context.Context, postID post.ID) error {
	post, err := s.postQuery.Get(ctx, postID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	s.fetcher.HydrateContentURLs(ctx, post)
	return nil
}

func (s *scrapeConsumer) hydrateNodeURLs(ctx context.Context, nodeID library.NodeID) error {
	node, err := s.nodeQuery.Get(ctx, library.NewID(xid.ID(nodeID)))
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	s.fetcher.HydrateContentURLs(ctx, node)
	return nil
}
