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
	l.logger.Debug("received beacon",
		slog.String("kind", cmd.Item.Kind.String()),
		slog.String("id", cmd.Item.ID.String()),
		slog.Any("subject", cmd.Subject.String()),
	)

	switch cmd.Item.Kind {
	case datagraph.KindThread:
		if subject, ok := cmd.Subject.Get(); ok {
			err := l.postReadStateWriter.UpsertReadState(ctx, subject, post.ID(cmd.Item.ID))
			if err != nil {
				l.logger.Error("failed to update read state",
					slog.String("error", err.Error()),
					slog.String("account_id", subject.String()),
					slog.String("thread_id", cmd.Item.ID.String()),
				)
				return err
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
			if err := l.handleBeacon(ctx, cmd); err != nil {
				logger.Error("failed to handle beacon", slog.String("error", err.Error()))
				return err
			}
			return nil
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
