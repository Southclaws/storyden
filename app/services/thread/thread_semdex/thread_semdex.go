package thread_semdex

import (
	"context"
	"log/slog"
	"time"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub/event"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
		// queue.New[mq.EventThreadCreated],
		// queue.New[mq.DeleteThread],
		),
		fx.Invoke(newSemdexer),
	)
}

// NOTE: If a reindex takes longer than the schedule time, there will be issues
// with duplicate messages since there's no checksum mechanism built currently.
// TODO: Make these parameters configurable by the SD instance administrator.
var (
	DefaultReindexSchedule  = time.Hour      // how frequently do we reindex
	DefaultReindexThreshold = time.Hour * 24 // ignore indexed_at after this
	DefaultReindexChunk     = 100            // size of query per reindex
)

type semdexer struct {
	logger        *slog.Logger
	db            *ent.Client
	threadQuerier thread.Repository
	threadWriter  thread.Repository
	semdexMutator semdex.Mutator
	semdexQuerier semdex.Querier
	bus           *event.Bus
}

func newSemdexer(
	ctx context.Context,
	lc fx.Lifecycle,
	cfg config.Config,
	logger *slog.Logger,

	db *ent.Client,
	threadQuerier thread.Repository,
	threadWriter thread.Repository,
	semdexMutator semdex.Mutator,
	semdexQuerier semdex.Querier,
	bus *event.Bus,
) {
	if cfg.SemdexProvider == "" {
		return
	}

	// re := semdexer{
	// 	logger:        logger,
	// 	db:            db,
	// 	threadQuerier: threadQuerier,
	// 	threadWriter:  threadQuerier,
	// 	// indexQueue:    indexQueue,
	// 	// deleteQueue:   deleteQueue,
	// 	semdexMutator: semdexMutator,
	// 	semdexQuerier: semdexQuerier,
	// }

	// lc.Append(fx.StartHook(func(_ context.Context) error {
	// 	sub, err := indexQueue.Subscribe(ctx)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	go func() {
	// 		for msg := range sub {
	// 			if err := re.indexThread(ctx, msg.Payload.ID); err != nil {
	// 				logger.Error("failed to index thread",
	// 					slog.String("error", err.Error()),
	// 					slog.String("post_id", msg.Payload.ID.String()),
	// 				)
	// 			}

	// 			msg.Ack()
	// 		}
	// 	}()

	// 	return nil
	// }))

	// lc.Append(fx.StartHook(func(_ context.Context) error {
	// 	sub, err := deleteQueue.Subscribe(ctx)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	go func() {
	// 		for msg := range sub {
	// 			if err := re.deindexThread(ctx, msg.Payload.ID); err != nil {
	// 				logger.Error("failed to deindex post", slog.String("error", err.Error()))
	// 			}

	// 			msg.Ack()
	// 		}
	// 	}()

	// 	return nil
	// }))

	// lc.Append(fx.StartHook(func(hctx context.Context) error {
	// 	// err := re.reindex(hctx, DefaultReindexThreshold, DefaultReindexChunk)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }

	// 	go re.schedule(ctx, DefaultReindexSchedule, DefaultReindexThreshold, DefaultReindexChunk)

	// 	return nil
	// }))
}
