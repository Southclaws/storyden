package http

import (
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/transports/http/middleware/origin"
	"github.com/Southclaws/storyden/app/transports/http/middleware/reqlog"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session"
	"github.com/Southclaws/storyden/app/transports/http/middleware/useragent"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/httpserver"
)

// Invoked by fx at runtime to mount the Echo router onto the http multiplexer.
// This is where all global middleware (not OpenAPI specific) is applied.
func MountOpenAPI(
	lc fx.Lifecycle,

	cfg config.Config,
	logger *zap.Logger,
	mux *http.ServeMux,
	router http.Handler,

	// Middleware providers
	cj *session.Jar,
) {
	lc.Append(fx.StartHook(func() {
		applied := httpserver.Apply(router,
			reqlog.WithLogger(logger),
			origin.WithCORS(cfg),
			useragent.UserAgentContext,
			cj.WithAuth,
		)

		// Mounting the Echo router must happen after all Echo's middleware and
		// routes have been set up so it's done inside the start lifecycle hook.
		mux.Handle("/", applied)
	}))
}
