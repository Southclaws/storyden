package notify_job

import (
	"context"
	"log/slog"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub/event"
)

func runNotifyConsumer(
	ctx context.Context,
	lc fx.Lifecycle,
	logger *slog.Logger,
	bus *event.Bus,
	ic *notifyConsumer,
) {
	lc.Append(fx.StartHook(func(hctx context.Context) error {
		_, err := event.SubscribeCommand(hctx, bus, "notify_job.send_notification", func(ctx context.Context, cmd *mq.CommandSendNotification) error {
			if err := ic.notify(ctx, cmd.TargetID, cmd.SourceID, cmd.Event, cmd.Item); err != nil {
				logger.Error("failed to notify", slog.String("error", err.Error()))
				return err
			}
			return nil
		})

		return err
	}))
}
