package queue

import (
	"context"
	"log/slog"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub/watermill"
)

func Build() fx.Option {
	return fx.Provide(func(
		ctx context.Context,
		cfg config.Config,
		l *slog.Logger,
	) (*QueueFactory, error) {
		sub, pub, err := watermill.NewWatermillQueue(cfg, l)
		if err != nil {
			return nil, err
		}

		return &QueueFactory{
			logger: l,
			pub:    pub,
			sub:    sub,
		}, nil
	})
}
