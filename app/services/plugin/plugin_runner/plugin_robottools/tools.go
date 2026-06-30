package plugin_robottools

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/Southclaws/opt"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/rs/xid"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner"
	robot_tools "github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	"github.com/Southclaws/storyden/lib/mcp"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

var unsafeToolNameChars = regexp.MustCompile(`[^a-zA-Z0-9_]`)

func FullyQualifiedName(installationID plugin.InstallationID, providerID, toolID string) string {
	return "plugin__" + safeSegment(installationID.String()) + "__" + safeSegment(providerID) + "__" + safeSegment(toolID)
}

func safeSegment(value string) string {
	value = unsafeToolNameChars.ReplaceAllString(value, "_")
	value = strings.Trim(value, "_")
	if value == "" {
		return "tool"
	}
	return value
}

func NewToolsForProvider(
	installationID plugin.InstallationID,
	provider rpc.RobotToolProviderCapabilityConfig,
	sess plugin_runner.Session,
) ([]*robot_tools.Tool, error) {
	out := make([]*robot_tools.Tool, 0, len(provider.Tools))
	for _, declaration := range provider.Tools {
		t, err := newTool(installationID, provider.ID, declaration, sess)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

func newTool(
	installationID plugin.InstallationID,
	providerID string,
	declaration rpc.RobotToolProviderToolConfig,
	sess plugin_runner.Session,
) (*robot_tools.Tool, error) {
	inputSchema, err := schemaFromManifest(declaration.InputSchema)
	if err != nil {
		return nil, fmt.Errorf("invalid input schema for plugin robot tool %q: %w", declaration.ID, err)
	}

	var outputSchema *jsonschema.Schema
	if raw, ok := declaration.OutputSchema.Get(); ok {
		outputSchema, err = schemaFromManifest(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid output schema for plugin robot tool %q: %w", declaration.ID, err)
		}
	}

	toolName := FullyQualifiedName(installationID, providerID, declaration.ID)
	def := &mcp.ToolDefinition{
		Name:                 toolName,
		Title:                declaration.Title.Or(declaration.Name),
		Description:          declaration.Description,
		InputSchema:          inputSchema,
		OutputSchema:         outputSchema,
		RequiresConfirmation: declaration.RequiresConfirmation.Or(false),
		Annotations:          annotationsFromManifest(declaration.Annotations),
	}

	return &robot_tools.Tool{
		Definition: def,
		Source:     "plugin",
		Builder: func(ctx context.Context) (tool.Tool, error) {
			run := robot_tools.RunContextFromContext(ctx)
			return functiontool.New(
				functiontool.Config{
					Name:                def.Name,
					Description:         def.Description,
					InputSchema:         def.InputSchema,
					RequireConfirmation: def.RequiresConfirmation && !robot_tools.ConfirmationDisabled(ctx),
				},
				func(ctx agent.Context, args map[string]interface{}) (robot_tools.ToolResult[map[string]interface{}], error) {
					return execute(ctx, sess, providerID, declaration.ID, run, args), nil
				},
			)
		},
	}, nil
}

func execute(
	ctx context.Context,
	sess plugin_runner.Session,
	providerID string,
	toolID string,
	run robot_tools.RunContext,
	args map[string]interface{},
) robot_tools.ToolResult[map[string]interface{}] {
	id := xid.New()
	callID := id.String()
	if toolCtx, ok := ctx.(interface{ FunctionCallID() string }); ok {
		if value := toolCtx.FunctionCallID(); value != "" {
			callID = value
		}
	}

	params := rpc.RPCRequestRobotToolCallParams{
		ProviderID: providerID,
		ToolID:     toolID,
		CallID:     callID,
		SessionID:  run.SessionID,
		AccountID:  run.AccountID,
		Arguments:  args,
	}
	if robotID, ok := run.RobotID.Get(); ok {
		params.RobotID = opt.New(robotID.String())
	}

	resp, err := sess.Send(ctx, id, &rpc.RPCRequestRobotToolCall{
		ID:      id,
		Jsonrpc: "2.0",
		Method:  "robot_tool_call",
		Params:  params,
	})
	if err != nil {
		return robot_tools.NewError[map[string]interface{}](err)
	}

	body, ok := resp.HostToPluginResponseUnionUnion.(*rpc.RPCResponseRobotToolCall)
	if !ok {
		return robot_tools.NewErrorMsg[map[string]interface{}](fmt.Sprintf("unexpected robot tool response: %T", resp.HostToPluginResponseUnionUnion))
	}
	if msg, ok := body.Error.Get(); ok && msg != "" {
		return robot_tools.NewErrorMsg[map[string]interface{}](msg)
	}

	return robot_tools.NewSuccess(body.Output)
}

func schemaFromManifest(raw rpc.RobotToolJSONSchema) (*jsonschema.Schema, error) {
	b, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}

	var schema jsonschema.Schema
	if err := json.Unmarshal(b, &schema); err != nil {
		return nil, err
	}
	if schema.Type != "object" {
		return nil, fmt.Errorf("schema type must be object")
	}

	return &schema, nil
}

func annotationsFromManifest(in opt.Optional[rpc.RobotToolAnnotations]) mcp.ToolAnnotations {
	annotations, ok := in.Get()
	if !ok {
		return mcp.ToolAnnotations{}
	}
	return mcp.ToolAnnotations{
		ReadOnlyHint:    annotations.ReadOnlyHint.Or(false),
		DestructiveHint: annotations.DestructiveHint.Or(false),
		IdempotentHint:  annotations.IdempotentHint.Or(false),
		OpenWorldHint:   annotations.OpenWorldHint.Or(false),
	}
}
