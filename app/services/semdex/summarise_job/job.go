package summarise_job

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

func runSummariseConsumer(
	ctx context.Context,
	lc fx.Lifecycle,
	l *zap.Logger,

	qnode pubsub.Topic[mq.SummariseNode],
	// qpost pubsub.Topic[mq.SummarisePost], // TODO
	// qprofile pubsub.Topic[mq.SummariseProfile], // TODO

	ic *summariseConsumer,
) {
	lc.Append(fx.StartHook(func(_ context.Context) error {
		nodeChan, err := qnode.Subscribe(ctx)
		if err != nil {
			panic(err)
		}

		go func() {
			for msg := range nodeChan {
				if err := ic.summariseNode(ctx, msg.Payload.ID); err != nil {
					l.Error("failed to summarise node", zap.Error(err))
				}

				msg.Ack()
			}
		}()

		return nil
	}))
}
