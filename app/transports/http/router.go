package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/transports/http/middleware/limiter"
	"github.com/Southclaws/storyden/app/transports/http/middleware/origin"
	"github.com/Southclaws/storyden/app/transports/http/middleware/reqlog"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
	"github.com/Southclaws/storyden/app/transports/http/middleware/useragent"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/httpserver"
	"github.com/Southclaws/storyden/internal/infrastructure/httpserver/chaos"
)

const MaxRequestSizeBytes = 10 * 1024 * 1024

// Invoked by fx at runtime to mount the Echo router onto the http multiplexer.
// This is where all global middleware (not OpenAPI specific) is applied.
func MountOpenAPI(
	lc fx.Lifecycle,

	cfg config.Config,
	logger *zap.Logger,
	mux *http.ServeMux,
	router *echo.Echo,

	// Middleware providers
	cj *session_cookie.Jar,
	rl *limiter.Middleware,
) {
	lc.Append(fx.StartHook(func() {
		applied := httpserver.Apply(router,
			origin.WithCORS(cfg),
			reqlog.WithLogger(logger),
			useragent.UserAgentContext,
			cj.WithAuth,
			rl.WithRateLimit,
			limiter.WithRequestSizeLimiter(MaxRequestSizeBytes),
			chaos.WithChaos(cfg),
		)

		// Mounting the Echo router must happen after all Echo's middleware and
		// routes have been set up so it's done inside the start lifecycle hook.
		mux.Handle("/", applied)
	}))
}
