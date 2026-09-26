package tools

import (
	"context"
	"encoding/json"
	"slices"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/session"

	"github.com/Southclaws/storyden/app/services/semdex/robot/workspacestate"
	"github.com/Southclaws/storyden/lib/mcp"
)

const callableToolsStateKey = "temp:callable_tools"
const ToolBlockersStateKey = "temp:tool_blockers"

func CaptureCallableTools() llmagent.BeforeModelCallback {
	return func(ctx agent.Context, request *model.LLMRequest) (*model.LLMResponse, error) {
		names := make([]string, 0, len(request.Tools))
		for name := range request.Tools {
			names = append(names, name)
		}

		slices.Sort(names)
		return nil, ctx.State().Set(callableToolsStateKey, names)
	}
}

func toolAvailability(ctx context.Context, selected *Tool) *mcp.RobotToolAvailabilityYaml {
	runtime, ok := ctx.(interface{ ReadonlyState() session.ReadonlyState })
	if !ok {
		return nil
	}

	state := runtime.ReadonlyState()
	raw, err := state.Get(callableToolsStateKey)
	if err != nil {
		return nil
	}

	data, err := json.Marshal(raw)
	if err != nil {
		return nil
	}

	var callable []string
	if json.Unmarshal(data, &callable) != nil {
		return nil
	}

	result := mcp.RobotToolAvailabilityYamlBlocked
	if slices.Contains(callable, selected.ADKName()) {
		result = mcp.RobotToolAvailabilityYamlCallable
		return &result
	}

	if selected.Definition.RequiresWorkspace && !workspacestate.Available(state) {
		return &result
	}

	if toolBlocker(state, selected.Definition.Name) != "" {
		return &result
	}

	loader := "tool_load"
	if selected.Definition.ToolsetOnly {
		loader = "toolset_load"
	}

	if slices.Contains(callable, loader) {
		result = mcp.RobotToolAvailabilityYamlLoadRequired
	}

	return &result
}

func toolBlocker(state session.ReadonlyState, id string) string {
	raw, err := state.Get(ToolBlockersStateKey)
	if err != nil {
		return ""
	}

	data, err := json.Marshal(raw)
	if err != nil {
		return ""
	}

	var blockers map[string]string
	if json.Unmarshal(data, &blockers) != nil {
		return ""
	}

	return blockers[id]
}
