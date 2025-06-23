package tools

import (
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/fx"
)

type All []server.ServerTool

func newTools(
	nodeTools *nodeTools,
	linkTools *linkTools,
	tagTools *tagTools,
) All {
	tools := []server.ServerTool{}

	tools = append(tools, nodeTools.tools...)
	tools = append(tools, linkTools.tools...)
	tools = append(tools, tagTools.tools...)

	return tools
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			newNodeTools,
			newLinkTools,
			newTagTools,
			newTools,
		),
	)
}
