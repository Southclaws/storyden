package logger

import (
	"log/slog"
	"os"

	"github.com/golang-cz/devslog"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/config"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(newLogger),
		fx.Invoke(replaceGlobals),
	)
}

func newLogger(cfg config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}

	logger := slog.New(func() slog.Handler {
		switch cfg.LogFormat {
		case "json":
			return slog.NewJSONHandler(os.Stdout, opts)

		case "dev":
			return devslog.NewHandler(os.Stdout, &devslog.Options{HandlerOptions: opts})

		default:
			return slog.NewTextHandler(os.Stdout, opts)
		}
	}())

	return logger
}

func replaceGlobals(c config.Config, l *slog.Logger) {
	// Use our logger for globals too, even though it's passed to
	// dependents most of the time using DI, the global logger is used
	// in a couple of places during startup/shutdown.
	slog.SetDefault(l)
}
