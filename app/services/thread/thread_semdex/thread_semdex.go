package thread_semdex

import (
	"context"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/datagraph/semdex"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub/queue"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			queue.New[mq.IndexThread],
			queue.New[mq.DeleteThread],
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
	logger        *zap.Logger
	db            *ent.Client
	threadQuerier thread.Repository
	threadWriter  thread.Repository
	indexQueue    pubsub.Topic[mq.IndexThread]
	deleteQueue   pubsub.Topic[mq.DeleteThread]
	indexer       semdex.Indexer
	deleter       semdex.Deleter
	retriever     semdex.Retriever
}

func newSemdexer(
	ctx context.Context,
	lc fx.Lifecycle,
	cfg config.Config,
	l *zap.Logger,

	db *ent.Client,
	threadQuerier thread.Repository,
	threadWriter thread.Repository,
	indexQueue pubsub.Topic[mq.IndexThread],
	deleteQueue pubsub.Topic[mq.DeleteThread],
	indexer semdex.Indexer,
	deleter semdex.Deleter,
	retriever semdex.Retriever,
) {
	if !cfg.SemdexEnabled {
		return
	}

	re := semdexer{
		logger:        l,
		db:            db,
		threadQuerier: threadQuerier,
		threadWriter:  threadQuerier,
		indexQueue:    indexQueue,
		deleteQueue:   deleteQueue,
		indexer:       indexer,
		deleter:       deleter,
		retriever:     retriever,
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
				if err := re.indexThread(ctx, msg.Payload.ID); err != nil {
					l.Error("failed to index post", zap.Error(err))
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
				if err := re.deindexThread(ctx, msg.Payload.ID); err != nil {
					l.Error("failed to index post", zap.Error(err))
				}

				msg.Ack()
			}
		}()

		return nil
	}))
}
