package analyse_job

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

func runAnalyseConsumer(
	ctx context.Context,
	lc fx.Lifecycle,
	l *zap.Logger,

	queue pubsub.Topic[mq.AnalyseAsset],
	consumer *analyseConsumer,
) {
	lc.Append(fx.StartHook(func(_ context.Context) error {
		nodeChan, err := queue.Subscribe(ctx)
		if err != nil {
			panic(err)
		}

		go func() {
			for msg := range nodeChan {
				nctx := session.GetSessionFromMessage(ctx, msg)

				if err := consumer.analyseAsset(nctx, msg.Payload.AssetID, msg.Payload.ContentFillRule); err != nil {
					l.Error("failed to index node", zap.Error(err))
				}

				msg.Ack()
			}
		}()

		return nil
	}))
}
