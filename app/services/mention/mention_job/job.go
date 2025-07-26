package mention_job

import (
	"context"
	"log/slog"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

func runMentionConsumer(
	ctx context.Context,
	lc fx.Lifecycle,
	logger *slog.Logger,

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
				if err := ic.mention(ctx, msg.Payload.By, msg.Payload.Source, msg.Payload.Item); err != nil {
					logger.Error("failed to record mention", slog.String("error", err.Error()))
				}

				msg.Ack()
			}
		}()

		return nil
	}))
}
