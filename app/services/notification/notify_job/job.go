package notify_job

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

func runNotifyConsumer(
	ctx context.Context,
	lc fx.Lifecycle,
	l *zap.Logger,

	queue pubsub.Topic[mq.Notification],

	ic *notifyConsumer,
) {
	lc.Append(fx.StartHook(func(_ context.Context) error {
		channel, err := queue.Subscribe(ctx)
		if err != nil {
			panic(err)
		}

		go func() {
			for msg := range channel {
				ctx = session.GetSessionFromMessage(ctx, msg)

				if err := ic.notify(ctx, msg.Payload.TargetID, msg.Payload.Event, msg.Payload.Item); err != nil {
					l.Error("failed to notify", zap.Error(err))
				}

				msg.Ack()
			}
		}()

		return nil
	}))
}
