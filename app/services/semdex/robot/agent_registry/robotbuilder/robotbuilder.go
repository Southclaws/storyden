package robotbuilder

import (
	_ "embed"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_registry"
	"github.com/Southclaws/storyden/app/services/semdex/robot/tools"
)

const (
	AppName     = "storyden"
	AgentName   = "storyden"
	DisplayName = "Storyden Robot Builder"
	Description = "Storyden's default agent that helps users get started and manage their community knowledge base."
)

//go:embed default.md
var instruction string

func Build() fx.Option {
	return fx.Invoke(Register)
}

func Register(registry *agent_registry.Registry) error {
	return registry.Register(Definition())
}

func Definition() agent_registry.Definition {
	return agent_registry.Definition{
		ID:           agent_registry.RobotBuilderID,
		Name:         DisplayName,
		Description:  Description,
		AppName:      AppName,
		AgentName:    AgentName,
		Instruction:  instruction,
		ToolNames:    tools.DefaultTools,
		Capabilities: tools.DefaultTools,
	}
}
