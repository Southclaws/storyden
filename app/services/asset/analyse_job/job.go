package analyse_job

import (
	"context"
	"log/slog"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

func runAnalyseConsumer(
	ctx context.Context,
	lc fx.Lifecycle,
	logger *slog.Logger,

	analyseQueue pubsub.Topic[mq.AnalyseAsset],
	downloadQueue pubsub.Topic[mq.DownloadAsset],
	consumer *analyseConsumer,
) {
	lc.Append(fx.StartHook(func(_ context.Context) error {
		analyseChan, err := analyseQueue.Subscribe(ctx)
		if err != nil {
			return err
		}

		go func() {
			for msg := range analyseChan {
				nctx := session.GetSessionFromMessage(ctx, msg)

				if err := consumer.analyseAsset(nctx, msg.Payload.AssetID, msg.Payload.ContentFillRule); err != nil {
					logger.Error("failed to analyse asset", slog.String("error", err.Error()))
				}

				msg.Ack()
			}
		}()

		downloadChan, err := downloadQueue.Subscribe(ctx)
		if err != nil {
			return err
		}

		go func() {
			for msg := range downloadChan {
				nctx := session.GetSessionFromMessage(ctx, msg)

				if err := consumer.downloadAsset(nctx, msg.Payload.URL, msg.Payload.ContentFillRule); err != nil {
					logger.Error("failed to download asset", slog.String("error", err.Error()))
				}

				msg.Ack()
			}
		}()

		return nil
	}))
}
