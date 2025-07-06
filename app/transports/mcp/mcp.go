package mcp

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/transports/mcp/tools"
)

func Build() fx.Option {
	return fx.Options(
		// Provide all the tools to the MCP server.
		tools.Build(),

		// Mount the MCP server into the HTTP mux.
		fx.Invoke(MountMCP),
	)
}
