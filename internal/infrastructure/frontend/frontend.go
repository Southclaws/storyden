// Package frontend provides a simple service that executes the Node.js
// application. This is only used by all-in-one installations.
package frontend

import (
	"context"
	"log/slog"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/config"
)

type Frontend interface {
	Run(ctx context.Context, path string)
	Ready() <-chan struct{}
}

func Build() fx.Option {
	return fx.Provide(func(lc fx.Lifecycle, logger *slog.Logger, cfg config.Config) Frontend {
		if cfg.RunFrontend == "" {
			return nil
		}

		fe := &NextjsProcess{logger: logger}

		lc.Append(fx.StartHook(func(ctx context.Context) {
			go fe.Run(ctx, cfg.RunFrontend)
		}))

		return fe
	})
}
