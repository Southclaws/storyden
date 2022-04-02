package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.uber.org/dig"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/api/src/config"
	"github.com/Southclaws/storyden/api/src/infra"
	"github.com/Southclaws/storyden/api/src/resources"
)

// Start starts the application and blocks until fatal error
// The server will shut down if the root context is cancelled
// nolint:errcheck
func Start(ctx context.Context) {
	app := fx.New(
		fx.NopLogger,

		config.Build(),
		infra.Build(),
		resources.Build(),
		// services.Build(),
		// interfaces.Build(),
	)

	err := app.Start(ctx)
	if err != nil {
		// Get the underlying error, without all the fx details.
		underlying := dig.RootCause(err)

		fmt.Println(underlying)

		os.Exit(1)
	}

	// Wait for context cancellation from the caller (interrupt signals etc.)
	<-ctx.Done()

	// Graceful shutdown time is 30 seconds. This context is passed to fx's stop
	// API which is then used to run all the OnStop hooks with a 30 sec timeout.
	ctx, cf := context.WithTimeout(context.Background(), time.Second*30)
	defer cf()

	if err := app.Stop(ctx); err != nil {
		zap.L().Error("fatal error occurred", zap.Error(err))
	}
}
