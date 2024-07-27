// Package frontend provides a simple service that executes the Node.js
// application. This is only used by all-in-one installations.
package frontend

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/internal/config"
)

type Frontend interface {
	Run(ctx context.Context, path string)
}

func Build() fx.Option {
	return fx.Invoke(func(lc fx.Lifecycle, l *zap.Logger, cfg config.Config) {
		if cfg.RunFrontend == "" {
			return
		}

		var fe Frontend = &NextjsProcess{l}

		lc.Append(fx.StartHook(func(ctx context.Context) {
			go fe.Run(ctx, cfg.RunFrontend)
		}))
	})
}
