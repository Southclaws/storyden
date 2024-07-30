package scrape_job

import (
	"context"
	"net/url"

	"go.uber.org/zap"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/services/link/fetcher"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type scrapeConsumer struct {
	l       *zap.Logger
	fetcher *fetcher.Fetcher
	queue   pubsub.Topic[mq.ScrapeLink]
}

func newScrapeConsumer(
	l *zap.Logger,
	fetcher *fetcher.Fetcher,
	queue pubsub.Topic[mq.ScrapeLink],
) *scrapeConsumer {
	return &scrapeConsumer{
		l:       l,
		fetcher: fetcher,
		queue:   queue,
	}
}

func (s *scrapeConsumer) scrapeLink(ctx context.Context, u url.URL) error {
	_, err := s.fetcher.ScrapeAndStore(ctx, u)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
