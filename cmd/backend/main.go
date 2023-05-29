package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/dig"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources"
	"github.com/Southclaws/storyden/app/services"
	transport "github.com/Southclaws/storyden/app/transports"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure"
)

// Start starts the application and blocks until fatal error
// The server will shut down if the root context is cancelled
// nolint:errcheck
func Start(ctx context.Context) {
	app := fx.New(
		fx.NopLogger,

		fx.Provide(func() context.Context { return ctx }),

		config.Build(),
		infrastructure.Build(),
		resources.Build(),
		services.Build(),
		transport.Build(),
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

func main() {
	godotenv.Load()

	ctx, cf := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cf()

	Start(ctx)
}
