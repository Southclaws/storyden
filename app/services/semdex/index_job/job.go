package index_job

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

func runIndexConsumer(
	ctx context.Context,
	lc fx.Lifecycle,
	l *zap.Logger,

	qnode pubsub.Topic[mq.IndexNode],
	qthread pubsub.Topic[mq.IndexThread],
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
					l.Error("failed to index post", zap.Error(err))
					msg.Nack()
					continue
				}

				msg.Ack()
			}
		}()

		go func() {
			for msg := range profileChan {
				if err := ic.indexProfile(ctx, msg.Payload.ID); err != nil {
					l.Error("failed to index post", zap.Error(err))
					msg.Nack()
					continue
				}

				msg.Ack()
			}
		}()

		return nil
	}))
}
