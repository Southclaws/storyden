package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	adkagent "google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"

	"github.com/Southclaws/storyden/lib/mcp"
)

type toolDiscoveryTools struct {
	registry *Registry
}

func newToolDiscoveryTools(registry *Registry) *toolDiscoveryTools {
	t := &toolDiscoveryTools{registry: registry}
	registry.Register(t.newSearchTool())
	registry.Register(t.newGetTool())
	return t
}

func (t *toolDiscoveryTools) newSearchTool() *Tool {
	def := mcp.GetToolSearchTool()
	return &Tool{
		Definition: def,
		Builder: func(context.Context) (tool.Tool, error) {
			return functiontool.New(functiontool.Config{Name: def.Name, Description: def.Description, InputSchema: def.InputSchema},
				func(ctx adkagent.Context, input mcp.ToolToolSearchInput) (*mcp.ToolToolSearchOutput, error) {
					return t.search(ctx, input), nil
				})
		},
		Handler: makeHandler(func(ctx context.Context, input mcp.ToolToolSearchInput) (*mcp.ToolToolSearchOutput, error) {
			return t.search(ctx, input), nil
		}),
	}
}

func (t *toolDiscoveryTools) search(_ context.Context, input mcp.ToolToolSearchInput) *mcp.ToolToolSearchOutput {
	maxResults := 0
	if input.MaxResults != nil {
		maxResults = *input.MaxResults
	}
	results := t.registry.SearchCatalogue(input.Query, maxResults)
	items := make([]mcp.RobotToolCatalogueItemYaml, 0, len(results))
	for _, result := range results {
		items = append(items, mcp.RobotToolCatalogueItemYaml{
			Id: result.ID, Name: result.Name, Description: result.Description,
		})
	}
	return &mcp.ToolToolSearchOutput{
		Tools:      items,
		NextAction: "Use tool_get with the selected tool ID to inspect its full schema. Then use tool_load for this conversation, or assign the ID when configuring a Robot or Toolset.",
	}
}

func (t *toolDiscoveryTools) newGetTool() *Tool {
	def := mcp.GetToolGetTool()
	return &Tool{
		Definition: def,
		Builder: func(context.Context) (tool.Tool, error) {
			return functiontool.New(functiontool.Config{Name: def.Name, Description: def.Description, InputSchema: def.InputSchema},
				func(ctx adkagent.Context, input mcp.ToolToolGetInput) (*mcp.ToolToolGetOutput, error) {
					return t.get(ctx, input)
				})
		},
		Handler: makeHandler(t.get),
	}
}

func (t *toolDiscoveryTools) get(ctx context.Context, input mcp.ToolToolGetInput) (*mcp.ToolToolGetOutput, error) {
	selected, err := t.registry.GetTool(ctx, input.Id)
	if err != nil {
		return nil, err
	}
	name := selected.Definition.Title
	if name == "" {
		name = selected.Definition.Name
	}
	if selected.Definition.InputSchema == nil {
		return nil, fmt.Errorf("tool %q does not expose an input schema", selected.Definition.Name)
	}
	inputSchema, err := schemaMap(selected.Definition.InputSchema)
	if err != nil {
		return nil, fmt.Errorf("encode input schema for tool %q: %w", selected.Definition.Name, err)
	}
	outputSchema, err := schemaMap(selected.Definition.OutputSchema)
	if err != nil {
		return nil, fmt.Errorf("encode output schema for tool %q: %w", selected.Definition.Name, err)
	}
	return &mcp.ToolToolGetOutput{
		Id:                   selected.Definition.Name,
		Name:                 name,
		Description:          selected.Definition.Description,
		InputSchema:          inputSchema,
		OutputSchema:         outputSchema,
		Toolsets:             append([]string{}, selected.Definition.Toolsets...),
		ToolsetOnly:          selected.Definition.ToolsetOnly,
		RequiresConfirmation: selected.Definition.RequiresConfirmation,
		RequiresWorkspace:    selected.Definition.RequiresWorkspace,
	}, nil
}

func schemaMap(schema *jsonschema.Schema) (map[string]interface{}, error) {
	if schema == nil {
		return nil, nil
	}
	encoded, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(encoded, &result); err != nil {
		return nil, err
	}
	return result, nil
}
