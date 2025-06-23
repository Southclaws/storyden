package tools

import (
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/fx"
)

type All []server.ServerTool

func newTools(
	nodeTools *nodeTools,
) All {
	tools := []server.ServerTool{}

	tools = append(tools, nodeTools.tools...)

	return tools
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			newNodeTools,
			newTools,
		),
	)
}
