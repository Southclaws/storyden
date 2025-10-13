package beacon_listener

import (
	"context"
	"log/slog"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/post_read_state"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type listener struct {
	logger              *slog.Logger
	postReadStateWriter *post_read_state.Writer
}

func newListener(logger *slog.Logger, postReadStateWriter *post_read_state.Writer) *listener {
	return &listener{
		logger:              logger,
		postReadStateWriter: postReadStateWriter,
	}
}

func (l *listener) handleBeacon(ctx context.Context, cmd *message.CommandSendBeacon) error {
	log := l.logger.With(
		slog.String("datagraph_kind", cmd.Item.Kind.String()),
		slog.String("datagraph_id", cmd.Item.ID.String()),
	)

	log.Debug("received beacon")

	switch cmd.Item.Kind {
	case datagraph.KindThread:
		if subject, ok := cmd.Subject.Get(); ok {
			log = log.With(slog.String("account_id", subject.String()))

			err := l.postReadStateWriter.UpsertReadState(ctx, subject, post.ID(cmd.Item.ID))
			if err != nil {
				log.Error("failed to update read state", slog.String("error", err.Error()))
				// NOTE: No error here, we don't care enough about read states
				// to re-queue the message in the DLQ and try again.
			}
		}
	}

	return nil
}

func runBeaconListener(
	ctx context.Context,
	lc fx.Lifecycle,
	logger *slog.Logger,
	bus *pubsub.Bus,
	l *listener,
) {
	lc.Append(fx.StartHook(func(hctx context.Context) error {
		_, err := pubsub.SubscribeCommand(hctx, bus, "beacon_listener.handler", func(ctx context.Context, cmd *message.CommandSendBeacon) error {
			return l.handleBeacon(ctx, cmd)
		})

		return err
	}))
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(newListener),
		fx.Invoke(runBeaconListener),
	)
}
