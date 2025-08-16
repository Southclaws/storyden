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
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub/event"
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

	bus *event.Bus

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
	bus *event.Bus,
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
		bus:           bus,
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

	lc.Append(fx.StartHook(func(hctx context.Context) error {
		_, err := event.Subscribe(hctx, bus, "node_semdex.published", func(ctx context.Context, evt *mq.EventNodePublished) error {
			return bus.SendCommand(ctx, &mq.CommandNodeIndex{ID: evt.ID})
		})
		if err != nil {
			return err
		}

		_, err = event.Subscribe(hctx, bus, "node_semdex.unpublished", func(ctx context.Context, evt *mq.EventNodeUnpublished) error {
			return bus.SendCommand(ctx, &mq.CommandNodeDeindex{ID: evt.ID})
		})
		if err != nil {
			return err
		}

		_, err = event.Subscribe(hctx, bus, "node_semdex.deleted", func(ctx context.Context, evt *mq.EventNodeDeleted) error {
			return bus.SendCommand(ctx, &mq.CommandNodeDeindex{ID: evt.ID})
		})
		if err != nil {
			return err
		}

		return nil
	}))

	lc.Append(fx.StartHook(func(hctx context.Context) error {
		_, err := event.SubscribeCommand(hctx, bus, "node_semdex.index", func(ctx context.Context, cmd *mq.CommandNodeIndex) error {
			return re.index(ctx, cmd.ID)
		})
		if err != nil {
			return err
		}

		_, err = event.SubscribeCommand(hctx, bus, "node_semdex.deindex", func(ctx context.Context, cmd *mq.CommandNodeDeindex) error {
			return re.deindex(ctx, cmd.ID)
		})
		if err != nil {
			return err
		}

		return nil
	}))
}
