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

	qprofile pubsub.Topic[mq.IndexProfile],

	ic *indexerConsumer,
) {
	lc.Append(fx.StartHook(func(_ context.Context) error {
		profileChan, err := qprofile.Subscribe(ctx)
		if err != nil {
			panic(err)
		}

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
