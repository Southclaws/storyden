package search_indexer

import (
	"context"
	"log/slog"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply_querier"
	"github.com/Southclaws/storyden/app/resources/post/thread_querier"
	"github.com/Southclaws/storyden/app/resources/profile/profile_querier"
	"github.com/Southclaws/storyden/app/services/search/searcher"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Indexer struct {
	logger *slog.Logger
	db     *ent.Client

	searchIndexer searcher.Indexer
	semdexMutator semdex.Mutator

	nodeQuerier    *node_querier.Querier
	threadQuerier  *thread_querier.Querier
	replyQuerier   *reply_querier.Querier
	profileQuerier *profile_querier.Querier

	bus       *pubsub.Bus
	chunkSize int
}

func runIndexerOnBoot(ctx context.Context, lc fx.Lifecycle, i *Indexer) {
	if i == nil {
		return
	}

	lc.Append(fx.StartHook(func(hctx context.Context) error {
		go func() {
			time.Sleep(time.Second)
			err := i.ReindexAll(ctx)
			if err != nil {
				i.logger.Error("failed to run initial reindex job", slog.String("error", err.Error()))
			}
		}()

		return nil
	}))
}

func newIndexer(
	ctx context.Context,
	lc fx.Lifecycle,
	cfg config.Config,
	logger *slog.Logger,
	db *ent.Client,

	searchIndexer searcher.Indexer,
	semdexMutator semdex.Mutator,

	nodeQuerier *node_querier.Querier,
	threadQuerier *thread_querier.Querier,
	replyQuerier *reply_querier.Querier,
	profileQuerier *profile_querier.Querier,

	bus *pubsub.Bus,
) *Indexer {
	if cfg.SearchProvider == "" || cfg.SearchProvider == "database" {
		return nil
	}

	idx := &Indexer{
		logger: logger,
		db:     db,

		searchIndexer: searchIndexer,
		semdexMutator: semdexMutator,

		nodeQuerier:    nodeQuerier,
		threadQuerier:  threadQuerier,
		replyQuerier:   replyQuerier,
		profileQuerier: profileQuerier,

		bus:       bus,
		chunkSize: cfg.SearchIndexChunkSize,
	}

	lc.Append(fx.StartHook(func(hctx context.Context) error {
		_, err := pubsub.Subscribe(ctx, idx.bus, "search_indexer.thread_published", func(ctx context.Context, evt *message.EventThreadPublished) error {
			return idx.bus.SendCommand(ctx, &message.CommandThreadIndex{ID: evt.ID})
		})
		if err != nil {
			return err
		}

		_, err = pubsub.Subscribe(ctx, idx.bus, "search_indexer.thread_updated", func(ctx context.Context, evt *message.EventThreadUpdated) error {
			return idx.bus.SendCommand(ctx, &message.CommandThreadIndex{ID: evt.ID})
		})
		if err != nil {
			return err
		}

		_, err = pubsub.Subscribe(ctx, idx.bus, "search_indexer.thread_unpublished", func(ctx context.Context, evt *message.EventThreadUnpublished) error {
			return idx.bus.SendCommand(ctx, &message.CommandThreadDeindex{ID: evt.ID})
		})
		if err != nil {
			return err
		}

		_, err = pubsub.Subscribe(ctx, idx.bus, "search_indexer.thread_deleted", func(ctx context.Context, evt *message.EventThreadDeleted) error {
			return idx.bus.SendCommand(ctx, &message.CommandThreadDeindex{ID: evt.ID})
		})
		if err != nil {
			return err
		}

		_, err = pubsub.Subscribe(ctx, idx.bus, "search_indexer.node_published", func(ctx context.Context, evt *message.EventNodePublished) error {
			return idx.bus.SendCommand(ctx, &message.CommandNodeIndex{ID: evt.ID})
		})
		if err != nil {
			return err
		}

		_, err = pubsub.Subscribe(ctx, idx.bus, "search_indexer.node_updated", func(ctx context.Context, evt *message.EventNodeUpdated) error {
			return idx.bus.SendCommand(ctx, &message.CommandNodeIndex{ID: evt.ID})
		})
		if err != nil {
			return err
		}

		_, err = pubsub.Subscribe(ctx, idx.bus, "search_indexer.node_unpublished", func(ctx context.Context, evt *message.EventNodeUnpublished) error {
			return idx.bus.SendCommand(ctx, &message.CommandNodeDeindex{ID: evt.ID})
		})
		if err != nil {
			return err
		}

		_, err = pubsub.Subscribe(ctx, idx.bus, "search_indexer.node_deleted", func(ctx context.Context, evt *message.EventNodeDeleted) error {
			return idx.bus.SendCommand(ctx, &message.CommandNodeDeindex{ID: evt.ID})
		})
		if err != nil {
			return err
		}

		_, err = pubsub.SubscribeCommand(ctx, idx.bus, "search_indexer.index_thread", func(ctx context.Context, cmd *message.CommandThreadIndex) error {
			return idx.indexThread(ctx, cmd.ID)
		})
		if err != nil {
			return err
		}

		_, err = pubsub.SubscribeCommand(ctx, idx.bus, "search_indexer.deindex_thread", func(ctx context.Context, cmd *message.CommandThreadDeindex) error {
			return idx.deindexThread(ctx, cmd.ID)
		})
		if err != nil {
			return err
		}

		_, err = pubsub.SubscribeCommand(ctx, idx.bus, "search_indexer.index_node", func(ctx context.Context, cmd *message.CommandNodeIndex) error {
			return idx.indexNode(ctx, cmd.ID)
		})
		if err != nil {
			return err
		}

		_, err = pubsub.SubscribeCommand(ctx, idx.bus, "search_indexer.deindex_node", func(ctx context.Context, cmd *message.CommandNodeDeindex) error {
			return idx.deindexNode(ctx, cmd.ID)
		})
		if err != nil {
			return err
		}

		_, err = pubsub.Subscribe(ctx, idx.bus, "search_indexer.reply_created", func(ctx context.Context, evt *message.EventThreadReplyCreated) error {
			return idx.bus.SendCommand(ctx, &message.CommandReplyIndex{ID: evt.ReplyID})
		})
		if err != nil {
			return err
		}

		_, err = pubsub.Subscribe(ctx, idx.bus, "search_indexer.reply_updated", func(ctx context.Context, evt *message.EventThreadReplyUpdated) error {
			return idx.bus.SendCommand(ctx, &message.CommandReplyIndex{ID: evt.ReplyID})
		})
		if err != nil {
			return err
		}

		_, err = pubsub.Subscribe(ctx, idx.bus, "search_indexer.reply_deleted", func(ctx context.Context, evt *message.EventThreadReplyDeleted) error {
			return idx.bus.SendCommand(ctx, &message.CommandReplyDeindex{ID: evt.ReplyID})
		})
		if err != nil {
			return err
		}

		_, err = pubsub.Subscribe(ctx, idx.bus, "search_indexer.reply_published", func(ctx context.Context, evt *message.EventThreadReplyPublished) error {
			return idx.bus.SendCommand(ctx, &message.CommandReplyIndex{ID: evt.ReplyID})
		})
		if err != nil {
			return err
		}

		_, err = pubsub.Subscribe(ctx, idx.bus, "search_indexer.reply_unpublished", func(ctx context.Context, evt *message.EventThreadReplyUnpublished) error {
			return idx.bus.SendCommand(ctx, &message.CommandReplyDeindex{ID: evt.ReplyID})
		})
		if err != nil {
			return err
		}

		_, err = pubsub.SubscribeCommand(ctx, idx.bus, "search_indexer.index_reply", func(ctx context.Context, cmd *message.CommandReplyIndex) error {
			return idx.indexReply(ctx, cmd.ID)
		})
		if err != nil {
			return err
		}

		_, err = pubsub.SubscribeCommand(ctx, idx.bus, "search_indexer.deindex_reply", func(ctx context.Context, cmd *message.CommandReplyDeindex) error {
			return idx.deindexReply(ctx, cmd.ID)
		})
		if err != nil {
			return err
		}

		_, err = pubsub.Subscribe(ctx, idx.bus, "search_indexer.account_created", func(ctx context.Context, evt *message.EventAccountCreated) error {
			return idx.bus.SendCommand(ctx, &message.CommandProfileIndex{ID: evt.ID})
		})
		if err != nil {
			return err
		}

		_, err = pubsub.Subscribe(ctx, idx.bus, "search_indexer.account_updated", func(ctx context.Context, evt *message.EventAccountUpdated) error {
			return idx.bus.SendCommand(ctx, &message.CommandProfileIndex{ID: evt.ID})
		})
		if err != nil {
			return err
		}

		_, err = pubsub.SubscribeCommand(ctx, idx.bus, "search_indexer.index_profile", func(ctx context.Context, cmd *message.CommandProfileIndex) error {
			return idx.indexProfile(ctx, cmd.ID)
		})
		if err != nil {
			return err
		}

		return nil
	}))

	return idx
}

func (idx *Indexer) indexThread(ctx context.Context, id post.ID) error {
	thread, err := idx.threadQuerier.Get(ctx, id, pagination.NewPageParams(1, 1), opt.NewEmpty[account.AccountID]())
	if err != nil {
		idx.logger.Error("failed to get thread for indexing", slog.String("id", id.String()), slog.String("error", err.Error()))
		return err
	}

	if err := idx.searchIndexer.Index(ctx, thread); err != nil {
		idx.logger.Error("failed to index thread", slog.String("id", id.String()), slog.String("error", err.Error()))
		return err
	}

	if _, err := idx.semdexMutator.Index(ctx, thread); err != nil {
		idx.logger.Error("failed to semdex thread", slog.String("id", id.String()), slog.String("error", err.Error()))
		return fault.Wrap(err, fctx.With(ctx))
	}

	idx.logger.Debug("indexed thread", slog.String("id", id.String()))
	return nil
}

func (idx *Indexer) deindexThread(ctx context.Context, id post.ID) error {
	if err := idx.searchIndexer.Deindex(ctx, &datagraph.Ref{
		ID:   xid.ID(id),
		Kind: datagraph.KindThread,
	}); err != nil {
		idx.logger.Error("failed to deindex thread", slog.String("id", id.String()), slog.String("error", err.Error()))
		return err
	}

	if _, err := idx.semdexMutator.Delete(ctx, xid.ID(id)); err != nil {
		idx.logger.Error("failed to desemdex thread", slog.String("id", id.String()), slog.String("error", err.Error()))
		return fault.Wrap(err, fctx.With(ctx))
	}

	idx.logger.Debug("deindexed thread", slog.String("id", id.String()))
	return nil
}

func (idx *Indexer) indexNode(ctx context.Context, id library.NodeID) error {
	node, err := idx.nodeQuerier.Get(ctx, library.NewID(xid.ID(id)))
	if err != nil {
		idx.logger.Error("failed to get node for indexing", slog.String("id", id.String()), slog.String("error", err.Error()))
		return err
	}

	if err := idx.searchIndexer.Index(ctx, node); err != nil {
		idx.logger.Error("failed to index node", slog.String("id", id.String()), slog.String("error", err.Error()))
		return err
	}

	if _, err := idx.semdexMutator.Index(ctx, node); err != nil {
		idx.logger.Error("failed to semdex node", slog.String("id", id.String()), slog.String("error", err.Error()))
		return fault.Wrap(err, fctx.With(ctx))
	}

	idx.logger.Debug("indexed node", slog.String("id", id.String()))
	return nil
}

func (idx *Indexer) deindexNode(ctx context.Context, id library.NodeID) error {
	if err := idx.searchIndexer.Deindex(ctx, &datagraph.Ref{
		ID:   xid.ID(id),
		Kind: datagraph.KindNode,
	}); err != nil {
		idx.logger.Error("failed to deindex node", slog.String("id", id.String()), slog.String("error", err.Error()))
		return err
	}

	if _, err := idx.semdexMutator.Delete(ctx, xid.ID(id)); err != nil {
		idx.logger.Error("failed to desemdex node", slog.String("id", id.String()), slog.String("error", err.Error()))
		return fault.Wrap(err, fctx.With(ctx))
	}

	idx.logger.Debug("deindexed node", slog.String("id", id.String()))
	return nil
}

func (idx *Indexer) indexReply(ctx context.Context, id post.ID) error {
	p, err := idx.replyQuerier.Get(ctx, id)
	if err != nil {
		idx.logger.Error("failed to get reply for indexing", slog.String("id", id.String()), slog.String("error", err.Error()))
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := idx.searchIndexer.Index(ctx, p); err != nil {
		idx.logger.Error("failed to index reply", slog.String("id", id.String()), slog.String("error", err.Error()))
		return err
	}

	if _, err := idx.semdexMutator.Index(ctx, p); err != nil {
		idx.logger.Error("failed to semdex reply", slog.String("id", id.String()), slog.String("error", err.Error()))
		return fault.Wrap(err, fctx.With(ctx))
	}

	idx.logger.Debug("indexed reply", slog.String("id", id.String()))
	return nil
}

func (idx *Indexer) deindexReply(ctx context.Context, id post.ID) error {
	if err := idx.searchIndexer.Deindex(ctx, &datagraph.Ref{
		ID:   xid.ID(id),
		Kind: datagraph.KindReply,
	}); err != nil {
		idx.logger.Error("failed to deindex reply", slog.String("id", id.String()), slog.String("error", err.Error()))
		return err
	}

	if _, err := idx.semdexMutator.Delete(ctx, xid.ID(id)); err != nil {
		idx.logger.Error("failed to desemdex reply", slog.String("id", id.String()), slog.String("error", err.Error()))
		return fault.Wrap(err, fctx.With(ctx))
	}

	idx.logger.Debug("deindexed reply", slog.String("id", id.String()))
	return nil
}

func (idx *Indexer) indexProfile(ctx context.Context, id account.AccountID) error {
	p, err := idx.profileQuerier.GetByID(ctx, id)
	if err != nil {
		idx.logger.Error("failed to get profile for indexing", slog.String("id", id.String()), slog.String("error", err.Error()))
		return fault.Wrap(err, fctx.With(ctx))
	}

	if p.GetContent().IsEmpty() {
		return nil
	}

	if _, err := idx.semdexMutator.Index(ctx, p); err != nil {
		idx.logger.Error("failed to semdex profile", slog.String("id", id.String()), slog.String("error", err.Error()))
		return fault.Wrap(err, fctx.With(ctx))
	}

	idx.logger.Debug("indexed profile", slog.String("id", id.String()))
	return nil
}
