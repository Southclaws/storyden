package robot

import (
	"fmt"
	"strings"

	adkagent "google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model"
	adktool "google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"

	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_registry"
	"github.com/Southclaws/storyden/app/services/semdex/robot/tools"
)

func annotateClientToolCalls(clientTools *agent_registry.ClientToolContext, scope InvocationContext) llmagent.AfterModelCallback {
	return func(ctx adkagent.Context, response *model.LLMResponse, responseErr error) (*model.LLMResponse, error) {
		if responseErr != nil || clientTools == nil || response == nil || response.Content == nil {
			return nil, nil
		}

		for _, part := range response.Content.Parts {
			if part == nil || part.FunctionCall == nil {
				continue
			}
			definition, ok := clientTools.Find(part.FunctionCall.Name)
			if !ok {
				continue
			}
			if part.PartMetadata == nil {
				part.PartMetadata = map[string]any{}
			}
			part.PartMetadata[agent_registry.ClientToolMetadataKey] = map[string]any{
				"source":      agent_registry.ClientToolSourceWebMCP,
				"client_id":   clientTools.ClientID,
				"scope":       scope,
				"title":       definition.Title,
				"annotations": definition.Annotations,
			}
		}

		return nil, nil
	}
}

func NewClientToolContext(clientID string, definitions []agent_registry.ClientToolDefinition) (*agent_registry.ClientToolContext, error) {
	return agent_registry.NewClientToolContext(clientID, definitions)
}

func validateClientToolNames(registry *tools.Registry, clientTools *agent_registry.ClientToolContext) error {
	if clientTools == nil {
		return nil
	}
	for _, definition := range clientTools.Tools {
		reservedRuntimeName := definition.Name == checkBackLaterToolName || definition.Name == unattendedFinishToolName || definition.Name == "transfer_to_agent"
		reservedPrefix := strings.HasPrefix(definition.Name, "adk_") || strings.HasPrefix(definition.Name, "robot_")
		if reservedRuntimeName || reservedPrefix {
			return fmt.Errorf("client tool name %q is reserved", definition.Name)
		}
		if _, exists := registry.FindByADKName(definition.Name); exists {
			return fmt.Errorf("client tool name %q conflicts with a registered Storyden tool", definition.Name)
		}
	}
	return nil
}

func buildClientToolset(clientTools *agent_registry.ClientToolContext) (adktool.Toolset, error) {
	if clientTools == nil || len(clientTools.Tools) == 0 {
		return nil, nil
	}

	toolList := make([]adktool.Tool, 0, len(clientTools.Tools))
	for _, definition := range clientTools.Tools {
		inputSchema, err := definition.JSONSchema()
		if err != nil {
			return nil, err
		}
		name := definition.Name
		clientTool, err := functiontool.New(
			functiontool.Config{
				Name:          name,
				Description:   definition.Description,
				InputSchema:   inputSchema,
				IsLongRunning: true,
			},
			func(context adkagent.Context, input map[string]any) (map[string]any, error) {
				return nil, fmt.Errorf("client tool %q cannot execute on the server", name)
			},
		)
		if err != nil {
			return nil, fmt.Errorf("build client tool %q: %w", name, err)
		}
		toolList = append(toolList, clientTool)
	}

	return &tools.Toolset{ToolList: toolList}, nil
}

func buildClientToolsets(registry *tools.Registry, clientTools *agent_registry.ClientToolContext) ([]adktool.Toolset, error) {
	if err := validateClientToolNames(registry, clientTools); err != nil {
		return nil, err
	}
	toolset, err := buildClientToolset(clientTools)
	if err != nil {
		return nil, err
	}
	if toolset == nil {
		return nil, nil
	}
	return []adktool.Toolset{toolset}, nil
}
