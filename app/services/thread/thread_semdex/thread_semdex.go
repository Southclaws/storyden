package thread_semdex

import (
	"context"
	"log/slog"
	"time"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

func Build() fx.Option {
	return fx.Options(
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
	bus           *pubsub.Bus
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
	bus *pubsub.Bus,
) {
	if cfg.SemdexProvider == "" {
		return
	}

	re := semdexer{
		logger:        logger,
		db:            db,
		threadQuerier: threadQuerier,
		threadWriter:  threadWriter,
		semdexMutator: semdexMutator,
		semdexQuerier: semdexQuerier,
		bus:           bus,
	}

	lc.Append(fx.StartHook(func(hctx context.Context) error {
		_, err := pubsub.Subscribe(hctx, bus, "thread_semdex.index_published", func(ctx context.Context, evt *mq.EventThreadPublished) error {
			return bus.SendCommand(ctx, &mq.CommandThreadIndex{ID: evt.ID})
		})
		if err != nil {
			return err
		}

		_, err = pubsub.Subscribe(hctx, bus, "thread_semdex.update_indexed", func(ctx context.Context, evt *mq.EventThreadUpdated) error {
			return bus.SendCommand(ctx, &mq.CommandThreadIndex{ID: evt.ID})
		})
		if err != nil {
			return err
		}

		_, err = pubsub.Subscribe(hctx, bus, "thread_semdex.remove_unpublished", func(ctx context.Context, evt *mq.EventThreadUnpublished) error {
			return bus.SendCommand(ctx, &mq.CommandThreadDeindex{ID: evt.ID})
		})
		if err != nil {
			return err
		}

		_, err = pubsub.Subscribe(hctx, bus, "thread_semdex.remove_deleted", func(ctx context.Context, evt *mq.EventThreadDeleted) error {
			return bus.SendCommand(ctx, &mq.CommandThreadDeindex{ID: evt.ID})
		})
		if err != nil {
			return err
		}

		return nil
	}))

	lc.Append(fx.StartHook(func(hctx context.Context) error {
		_, err := pubsub.SubscribeCommand(hctx, bus, "thread_semdex.index", func(ctx context.Context, cmd *mq.CommandThreadIndex) error {
			return re.indexThread(ctx, cmd.ID)
		})
		if err != nil {
			return err
		}

		_, err = pubsub.SubscribeCommand(hctx, bus, "thread_semdex.deindex", func(ctx context.Context, cmd *mq.CommandThreadDeindex) error {
			return re.deindexThread(ctx, cmd.ID)
		})
		if err != nil {
			return err
		}

		return nil
	}))

	lc.Append(fx.StartHook(func(hctx context.Context) error {
		go re.schedule(ctx, DefaultReindexSchedule, DefaultReindexThreshold, DefaultReindexChunk)
		return nil
	}))
}
