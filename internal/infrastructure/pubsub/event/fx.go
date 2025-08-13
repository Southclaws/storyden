package event

import (
	"context"
	"log/slog"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub/watermill"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(func(
			lc fx.Lifecycle,
			ctx context.Context,
			cfg config.Config,
			l *slog.Logger,
		) (*Bus, error) {
			sub, pub, err := watermill.NewWatermillQueue(cfg, l)
			if err != nil {
				return nil, err
			}

			bus, err := New(lc, l, ctx, cfg, pub, sub)
			if err != nil {
				return nil, err
			}

			return bus, nil
		}),
	)
}
