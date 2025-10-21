package rpc

import (
	"context"
	"log/slog"
	"net/http"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/plugin/plugin_reader"
	"github.com/Southclaws/storyden/app/resources/plugin/plugin_writer"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner"
	"github.com/Southclaws/storyden/app/transports/http/middleware/headers"
	"github.com/Southclaws/storyden/app/transports/http/middleware/limiter"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/httpserver"
)

func MountRPC(
	lc fx.Lifecycle,
	ctx context.Context,
	logger *slog.Logger,
	cfg config.Config,

	settings *settings.SettingsRepository,
	runner plugin_runner.Host,
	pluginReader *plugin_reader.Reader,
	pluginWriter *plugin_writer.Writer,

	mux *http.ServeMux,

	ri *headers.Middleware,
	rl *limiter.Middleware,
) {
	lc.Append(fx.StartHook(func() error {
		// set, err := settings.Get(ctx)
		// if err != nil {
		// 	return err
		// }

		handler := NewWebSocketHandler(logger, runner, pluginReader, pluginWriter)
		lo := newRequestLogger(logger)

		applied := httpserver.Apply(http.HandlerFunc(handler.HandleWebSocket),
			ri.WithHeaderContext(),
			lo.WithLogger(),
			rl.WithRateLimit(),
		)

		mux.Handle("/rpc", applied)

		return nil
	}))
}
