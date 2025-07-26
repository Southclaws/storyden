package scrape_job

import (
	"context"
	"log/slog"

	"go.uber.org/fx"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

func runScrapeConsumer(
	ctx context.Context,
	lc fx.Lifecycle,
	logger *slog.Logger,

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
				if err := ic.scrapeLink(ctx, msg.Payload.URL, opt.NewPtr(msg.Payload.Item)); err != nil {
					logger.Error("failed to scrape link", slog.String("error", err.Error()))
				}

				msg.Ack()
			}
		}()

		return nil
	}))
}
