package mention_job

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

func runMentionConsumer(
	ctx context.Context,
	lc fx.Lifecycle,
	l *zap.Logger,

	queue pubsub.Topic[mq.Mention],

	ic *mentionConsumer,
) {
	lc.Append(fx.StartHook(func(_ context.Context) error {
		channel, err := queue.Subscribe(ctx)
		if err != nil {
			panic(err)
		}

		go func() {
			for msg := range channel {
				ctx = session.GetSessionFromMessage(ctx, msg)

				if err := ic.mention(ctx, msg.Payload.Source, msg.Payload.Item); err != nil {
					l.Error("failed to record mention", zap.Error(err))
				}

				msg.Ack()
			}
		}()

		return nil
	}))
}
