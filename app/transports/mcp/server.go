package mcp

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/http/middleware/limiter"
	"github.com/Southclaws/storyden/app/transports/http/middleware/origin"
	"github.com/Southclaws/storyden/app/transports/http/middleware/reqlog"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
	"github.com/Southclaws/storyden/app/transports/mcp/tools"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/httpserver"
)

func MountMCP(
	lc fx.Lifecycle,
	ctx context.Context,
	logger *slog.Logger,
	cfg config.Config,

	settings *settings.SettingsRepository,
	allTools tools.All,

	mux *http.ServeMux,

	// NOTE: This is duplicated from the OpenAPI router because there's an issue
	// in the OpenAPI codegen that makes mounting it on a sub-router difficult.
	// Eventually, when that's fixed, middleware can be declared once at root.
	co *origin.Middleware,
	lo *reqlog.Middleware,
	cj *session_cookie.Jar,
	rl *limiter.Middleware,
) {
	if !cfg.MCPEnabled {
		return
	}

	lc.Append(fx.StartHook(func() error {
		set, err := settings.Get(ctx)
		if err != nil {
			return err
		}

		s := server.NewMCPServer(
			set.Title.Or("Storyden"),
			"rolling", // NOTE: Worth providing versioning yet?
			server.WithToolCapabilities(true),
			server.WithRecovery(),
			server.WithLogging(),
		)

		s.AddTools(allTools...)

		// MCP is mounted on the root `/mcp` path, not under `/api`.
		sse := server.NewSSEServer(s,
			server.WithSSEEndpoint("/mcp/sse"),
			server.WithMessageEndpoint("/mcp/message"),
		)

		applied := httpserver.Apply(sse,
			co.WithCORS(),
			lo.WithLogger(),
			cj.WithAuth(),
			rl.WithRequestSizeLimiter(),
			rl.WithRateLimit(),
			withStrictAuthMCP(),
		)

		mux.Handle("/mcp/", applied)

		return nil
	}))
}

// middleware for MCP-specific authentication checks. MCP is fully behind auth
// so any requests require either a session cookie or an access key.
func withStrictAuthMCP() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := session.GetAccountID(r.Context())
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			withFlusher, ok := GetFlusher(w)
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Failed to get flusher for MCP server"))
				return
			}

			next.ServeHTTP(withFlusher, r)
		})
	}
}

func GetFlusher(w http.ResponseWriter) (http.ResponseWriter, bool) {
	for {
		if _, ok := w.(http.Flusher); ok {
			return w, true
		}
		// Try to unwrap
		if unwrapper, ok := w.(interface{ Unwrap() http.ResponseWriter }); ok {
			w = unwrapper.Unwrap()
		} else {
			return nil, false
		}
	}
}
