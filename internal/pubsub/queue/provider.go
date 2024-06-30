package queue

import (
	"context"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/pubsub/watermill"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Build() fx.Option {
	return fx.Provide(func(
		ctx context.Context,
		cfg config.Config,
		l *zap.Logger,
	) (*QueueFactory, error) {
		sub, pub, err := watermill.NewWatermillQueue(cfg, l)
		if err != nil {
			return nil, err
		}

		return &QueueFactory{
			log: l,
			pub: pub,
			sub: sub,
		}, nil
	})
}
