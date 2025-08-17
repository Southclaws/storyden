package notify_job

import (
	"context"
	"log/slog"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

func runNotifyConsumer(
	ctx context.Context,
	lc fx.Lifecycle,
	logger *slog.Logger,
	bus *pubsub.Bus,
	ic *notifyConsumer,
) {
	lc.Append(fx.StartHook(func(hctx context.Context) error {
		_, err := pubsub.SubscribeCommand(hctx, bus, "notify_job.send_notification", func(ctx context.Context, cmd *message.CommandSendNotification) error {
			if err := ic.notify(ctx, cmd.TargetID, cmd.SourceID, cmd.Event, cmd.Item); err != nil {
				logger.Error("failed to notify", slog.String("error", err.Error()))
				return err
			}
			return nil
		})

		return err
	}))
}
