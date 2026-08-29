package denbot

import (
	_ "embed"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_registry"
	"github.com/Southclaws/storyden/app/services/semdex/robot/toolsets/system_content_search"
	"github.com/Southclaws/storyden/app/services/semdex/robot/toolsets/system_discussions"
	"github.com/Southclaws/storyden/app/services/semdex/robot/toolsets/system_documents"
	"github.com/Southclaws/storyden/app/services/semdex/robot/toolsets/system_library"
	"github.com/Southclaws/storyden/app/services/semdex/robot/toolsets/system_memory"
	"github.com/Southclaws/storyden/app/services/semdex/robot/toolsets/system_robot_studio"
)

const (
	ID          = "denbot"
	AppName     = "denbot"
	AgentName   = "denbot"
	DisplayName = "Denbot"
	Description = "Storyden's general-purpose coordinator for community work, reusable tools, and specialist Robot delegation."
)

//go:embed playbook.md
var instruction string

func Build() fx.Option {
	return fx.Invoke(Register)
}

func Register(registry *agent_registry.Registry) error {
	return registry.Register(Definition())
}

func Definition() agent_registry.Definition {
	return agent_registry.Definition{
		ID:          ID,
		Name:        DisplayName,
		Description: Description,
		AppName:     AppName,
		AgentName:   AgentName,
		Instruction: instruction,
		ToolsetRefs: []string{
			system_content_search.ID,
			system_documents.ID,
			system_discussions.ID,
			system_library.ID,
			system_memory.ID,
			system_robot_studio.ID,
		},
		Capabilities: []string{"tool_search", "toolset_search", "robot_search", "delegation"},
	}
}
