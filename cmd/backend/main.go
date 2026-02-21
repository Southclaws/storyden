package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/dig"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/Southclaws/storyden/app/resources"
	"github.com/Southclaws/storyden/app/services"
	transport "github.com/Southclaws/storyden/app/transports"
	"github.com/Southclaws/storyden/internal/boot_time"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/fxlogger"
	"github.com/Southclaws/storyden/internal/infrastructure"
)

// Start starts the application and blocks until fatal error
// The server will shut down if the root context is cancelled
// nolint:errcheck
func Start(ctx context.Context) {
	// TEMPORARY: Do not release this in v1.25.9!!!
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	logger = logger.With("component", "bootstrap")

	app := fx.New(
		fx.WithLogger(func() fxevent.Logger {
			return fxlogger.New(logger)
		}),

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
		slog.Error("fatal error occurred", slog.String("error", err.Error()))
	}
}

func main() {
	boot_time.StartedAt = time.Now()

	godotenv.Load()

	ctx, cf := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cf()

	Start(ctx)
}
