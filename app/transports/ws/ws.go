// Package ws provides a transport layer using WebSockets.
package ws

import (
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/transports/http/middleware/limiter"
	"github.com/Southclaws/storyden/app/transports/http/middleware/origin"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session"
	"github.com/Southclaws/storyden/app/transports/ws/subscription"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/httpserver"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(subscription.New),
		fx.Invoke(MountWebSocketHandler),
	)
}

func MountWebSocketHandler(
	lc fx.Lifecycle,

	cfg config.Config,
	logger *zap.Logger,
	mux *http.ServeMux,
	cj *session.Jar,

	handler *subscription.Handler,
) {
	lc.Append(fx.StartHook(func() {
		handler := http.HandlerFunc(handler.Handle)

		applied := httpserver.Apply(handler,
			origin.WithCORS(cfg),
			cj.WithAuth,
			limiter.WithRateLimiter(cfg),
		)

		mux.Handle("/ws", applied)
	}))
}
