package mcp

import (
	"context"
	"net/http"

	"github.com/mark3labs/mcp-go/server"

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

func New(ctx context.Context, cfg config.Config, settings *settings.SettingsRepository, allTools tools.All) (*server.SSEServer, error) {
	set, err := settings.Get(ctx)
	if err != nil {
		return nil, err
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

	return sse, nil
}

func MountMCP(
	cfg config.Config,
	mux *http.ServeMux,
	s *server.SSEServer,

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

	applied := httpserver.Apply(s,
		co.WithCORS(),
		lo.WithLogger(),
		cj.WithAuth(),
		rl.WithRequestSizeLimiter(),
		rl.WithRateLimit(),
		withStrictAuthMCP(),
	)

	mux.Handle("/mcp/", applied)
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

			next.ServeHTTP(w, r)
		})
	}
}
