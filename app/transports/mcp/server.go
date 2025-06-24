package mcp

import (
	"net/http"

	"github.com/mark3labs/mcp-go/server"

	"github.com/Southclaws/storyden/app/transports/http/middleware/reqlog"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
	"github.com/Southclaws/storyden/app/transports/mcp/tools"
	"github.com/Southclaws/storyden/internal/config"
)

func New(cfg config.Config, allTools tools.All) *server.SSEServer {
	s := server.NewMCPServer(
		"Storyden", // TODO: Load title from Settings
		"rolling",  // NOTE: Worth providing versioning yet?
		server.WithToolCapabilities(true),
		server.WithRecovery(),
		server.WithLogging(),
	)

	s.AddTools(allTools...)

	return server.NewSSEServer(s, server.WithSSEEndpoint("/mcp/sse"), server.WithMessageEndpoint("/mcp/message"))
}

func MountMCP(
	mux *http.ServeMux,
	s *server.SSEServer,
	cj *session_cookie.Jar,
	reqlog *reqlog.Middleware, // TODO: Make Flushable
) {
	// TODO: Add middleware to this (until we fix oapi-codegen mount)
	// also, ensure every request requires auth, since the ctx checks are used
	// in handlers to perform RBAC checks already so no auth will just error.
	// TODO: Introduce config MCP_ENABLED default to disabled
	mux.Handle("/mcp/", cj.WithAuth()(s))
}
