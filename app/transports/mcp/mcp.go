package mcp

import (
	"go.uber.org/fx"
)

func Build() fx.Option {
	return fx.Options(
		// Mount the MCP server into the HTTP mux.
		fx.Invoke(MountMCP),
	)
}
