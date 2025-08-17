package mention_job

import (
	"context"
	"log/slog"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

func runMentionConsumer(
	ctx context.Context,
	lc fx.Lifecycle,
	logger *slog.Logger,
	bus *pubsub.Bus,
	ic *mentionConsumer,
) {
	lc.Append(fx.StartHook(func(hctx context.Context) error {
		_, err := pubsub.Subscribe(hctx, bus, "mention_job.notify_mentions", func(ctx context.Context, evt *message.EventMemberMentioned) error {
			if err := ic.mention(ctx, evt.By, evt.Source, evt.Item); err != nil {
				logger.Error("failed to record mention", slog.String("error", err.Error()))
				return err
			}
			return nil
		})

		return err
	}))
}
