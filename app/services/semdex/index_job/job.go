package index_job

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/internal/pubsub"
)

func runIndexConsumer(
	ctx context.Context,
	lc fx.Lifecycle,
	l *zap.Logger,

	qnode pubsub.Topic[mq.IndexNode],
	qpost pubsub.Topic[mq.IndexPost],
	qprofile pubsub.Topic[mq.IndexProfile],

	ic *indexerConsumer,
) {
	lc.Append(fx.StartHook(func(_ context.Context) error {
		nodeChan, err := qnode.Subscribe(ctx)
		if err != nil {
			panic(err)
		}

		postChan, err := qpost.Subscribe(ctx)
		if err != nil {
			panic(err)
		}

		profileChan, err := qprofile.Subscribe(ctx)
		if err != nil {
			panic(err)
		}

		go func() {
			for msg := range nodeChan {
				if err := ic.indexNode(ctx, msg.Payload.ID); err != nil {
					l.Error("failed to index node", zap.Error(err))
					msg.Nack()
					continue
				}

				msg.Ack()
			}
		}()

		go func() {
			for msg := range postChan {
				if err := ic.indexPost(ctx, msg.Payload.ID); err != nil {
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
