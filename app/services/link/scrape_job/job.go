package scrape_job

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

func runScrapeConsumer(
	ctx context.Context,
	lc fx.Lifecycle,
	l *zap.Logger,

	queue pubsub.Topic[mq.ScrapeLink],

	ic *scrapeConsumer,
) {
	lc.Append(fx.StartHook(func(_ context.Context) error {
		channel, err := queue.Subscribe(ctx)
		if err != nil {
			panic(err)
		}

		go func() {
			for msg := range channel {
				if err := ic.scrapeLink(ctx, msg.Payload.URL, msg.Payload.Item); err != nil {
					l.Error("failed to scrape link", zap.Error(err))
				}

				msg.Ack()
			}
		}()

		return nil
	}))
}
