package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/transports/http/middleware/chaos"
	"github.com/Southclaws/storyden/app/transports/http/middleware/frontend"
	"github.com/Southclaws/storyden/app/transports/http/middleware/headers"
	"github.com/Southclaws/storyden/app/transports/http/middleware/limiter"
	"github.com/Southclaws/storyden/app/transports/http/middleware/origin"
	"github.com/Southclaws/storyden/app/transports/http/middleware/reqlog"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/httpserver"
)

// Invoked by fx at runtime to mount the Echo router onto the http multiplexer.
// This is where all global middleware (not OpenAPI specific) is applied.
func MountOpenAPI(
	lc fx.Lifecycle,

	cfg config.Config,
	mux *http.ServeMux,
	router *echo.Echo,

	// Middleware providers
	co *origin.Middleware,
	lo *reqlog.Middleware,
	fe *frontend.Provider,
	ri *headers.Middleware,
	cj *session_cookie.Jar,
	rl *limiter.Middleware,
	cm *chaos.Middleware,
) {
	lc.Append(fx.StartHook(func() {
		applied := httpserver.Apply(router,
			co.WithCORS(),
			lo.WithLogger(),
			fe.WithFrontendProxy(),
			ri.WithHeaderContext(),
			cj.WithAuth(),
			rl.WithRequestSizeLimiter(),
			rl.WithRateLimit(),
			cm.WithChaos(),
		)

		// Health check endpoint does not need any middleware, mounted directly.
		mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// Mounting the Echo router must happen after all Echo's middleware and
		// routes have been set up so it's done inside the start lifecycle hook.
		mux.Handle("/", applied)
	}))
}
