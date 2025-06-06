package node_semdex

import (
	"context"
	"log/slog"
	"time"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/tag/tag_writer"
	"github.com/Southclaws/storyden/app/services/library/node_mutate"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/app/services/tag/autotagger"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub/queue"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			queue.New[mq.IndexNode],
			queue.New[mq.DeleteNode],
		),
		fx.Invoke(newSemdexer),
	)
}

// NOTE: If a reindex takes longer than the schedule time, there will be issues
// with duplicate messages since there's no checksum mechanism built currently.
// TODO: Make these parameters configurable by the SD instance administrator.
var (
	DefaultReindexSchedule  = time.Hour * 21 // how frequently do we reindex
	DefaultReindexThreshold = time.Hour * 24 // ignore indexed_at after this
	DefaultReindexChunk     = 100            // size of query per reindex
)

type semdexer struct {
	logger *slog.Logger
	db     *ent.Client

	nodeQuerier *node_querier.Querier
	nodeWriter  *node_writer.Writer
	nodeUpdater *node_mutate.Manager

	indexQueue  pubsub.Topic[mq.IndexNode]
	deleteQueue pubsub.Topic[mq.DeleteNode]

	semdexMutator semdex.Mutator
	semdexQuerier semdex.Querier

	tagger    *autotagger.Tagger
	tagWriter *tag_writer.Writer
}

func newSemdexer(
	ctx context.Context,
	lc fx.Lifecycle,
	cfg config.Config,
	logger *slog.Logger,

	db *ent.Client,
	nodeQuerier *node_querier.Querier,
	nodeWriter *node_writer.Writer,
	nodeUpdater *node_mutate.Manager,
	indexQueue pubsub.Topic[mq.IndexNode],
	deleteQueue pubsub.Topic[mq.DeleteNode],
	semdexMutator semdex.Mutator,
	semdexQuerier semdex.Querier,

	tagger *autotagger.Tagger,
	tagWriter *tag_writer.Writer,
) {
	if cfg.SemdexProvider == "" {
		return
	}

	re := semdexer{
		logger:        logger,
		db:            db,
		nodeQuerier:   nodeQuerier,
		nodeWriter:    nodeWriter,
		nodeUpdater:   nodeUpdater,
		indexQueue:    indexQueue,
		deleteQueue:   deleteQueue,
		semdexMutator: semdexMutator,
		semdexQuerier: semdexQuerier,

		tagger:    tagger,
		tagWriter: tagWriter,
	}

	lc.Append(fx.StartHook(func(hctx context.Context) error {
		err := re.reindex(hctx, DefaultReindexThreshold, DefaultReindexChunk)
		if err != nil {
			return err
		}

		go re.schedule(ctx, DefaultReindexSchedule, DefaultReindexThreshold, DefaultReindexChunk)

		return nil
	}))

	lc.Append(fx.StartHook(func(_ context.Context) error {
		sub, err := indexQueue.Subscribe(ctx)
		if err != nil {
			return err
		}

		go func() {
			for msg := range sub {
				if err := re.index(ctx, msg.Payload.ID); err != nil {
					logger.Error("failed to index node", slog.String("error", err.Error()))
				}

				msg.Ack()
			}
		}()

		return nil
	}))

	lc.Append(fx.StartHook(func(_ context.Context) error {
		sub, err := deleteQueue.Subscribe(ctx)
		if err != nil {
			return err
		}

		go func() {
			for msg := range sub {
				if err := re.deindex(ctx, msg.Payload.ID); err != nil {
					logger.Error("failed to index node", slog.String("error", err.Error()))
				}

				msg.Ack()
			}
		}()

		return nil
	}))
}
