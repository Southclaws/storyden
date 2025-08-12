package index_job

import (
	"context"
	"log/slog"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

func runIndexConsumer(
	ctx context.Context,
	lc fx.Lifecycle,
	logger *slog.Logger,

	qnode pubsub.Topic[mq.IndexNode],
	qreply pubsub.Topic[mq.IndexReply],
	qprofile pubsub.Topic[mq.IndexProfile],

	ic *indexerConsumer,
) {
	lc.Append(fx.StartHook(func(_ context.Context) error {
		replyChan, err := qreply.Subscribe(ctx)
		if err != nil {
			panic(err)
		}

		profileChan, err := qprofile.Subscribe(ctx)
		if err != nil {
			panic(err)
		}

		go func() {
			for msg := range replyChan {
				if err := ic.indexReply(ctx, msg.Payload.ID); err != nil {
					logger.Error("failed to index reply",
						slog.String("error", err.Error()),
						slog.String("post_id", msg.Payload.ID.String()),
					)
					msg.Nack()
					continue
				}

				msg.Ack()
			}
		}()

		go func() {
			for msg := range profileChan {
				if err := ic.indexProfile(ctx, msg.Payload.ID); err != nil {
					logger.Error("failed to index profile",
						slog.String("error", err.Error()),
						slog.String("account_id", msg.Payload.ID.String()),
					)
					msg.Nack()
					continue
				}

				msg.Ack()
			}
		}()

		return nil
	}))
}
