package tools

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/robot/llm_provider"
	"github.com/Southclaws/storyden/app/resources/robot/model_ref"
	"github.com/Southclaws/storyden/app/resources/robot/robot_querier"
	"github.com/Southclaws/storyden/app/resources/robot/robot_ref"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/app/resources/robot/robot_writer"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_registry"
	"github.com/Southclaws/storyden/lib/mcp"
)

type robotTools struct {
	logger           *slog.Logger
	robotQuerier     *robot_querier.Querier
	robotWriter      *robot_writer.Writer
	robotSessionRepo *robot_session.Repository
	registry         *Registry
	agentRegistry    *agent_registry.Registry
	modelFactory     *llm_provider.Factory
}

func newRobotTools(
	logger *slog.Logger,
	robotQuerier *robot_querier.Querier,
	robotWriter *robot_writer.Writer,
	robotSessionRepo *robot_session.Repository,
	registry *Registry,
	agentRegistry *agent_registry.Registry,
	modelFactory *llm_provider.Factory,
) *robotTools {
	t := &robotTools{
		logger:           logger,
		robotQuerier:     robotQuerier,
		robotWriter:      robotWriter,
		robotSessionRepo: robotSessionRepo,
		registry:         registry,
		agentRegistry:    agentRegistry,
		modelFactory:     modelFactory,
	}

	registry.Register(t.newRobotSwitchTool())
	registry.Register(t.newSystemRobotToolCatalogTool())
	registry.Register(t.newRobotCreateTool())
	registry.Register(t.newRobotListTool())
	registry.Register(t.newRobotGetTool())
	registry.Register(t.newRobotUpdateTool())
	registry.Register(t.newRobotDeleteTool())
	registry.Register(newThrowAnErrorTool())

	return t
}

func (rt *robotTools) newRobotSwitchTool() *Tool {
	toolDef := mcp.GetRobotSwitchTool()

	return &Tool{
		Definition:   toolDef,
		IsClientTool: true,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			result, err := rt.robotQuerier.List(ctx, pagination.NewPageParams(1, 20))
			if err != nil {
				return nil, err
			}

			robotIDs := make([]any, 0, len(result.Items)+len(rt.agentRegistry.List(false)))
			for _, agent := range rt.agentRegistry.List(false) {
				robotIDs = append(robotIDs, agent.ID)
			}
			for _, r := range result.Items {
				robotIDs = append(robotIDs, r.ID.String())
			}

			inputSchema := toolDef.InputSchema
			if inputSchema.Properties != nil {
				if robotIDProp, ok := inputSchema.Properties["robot_id"]; ok {
					robotIDProp.Enum = robotIDs
				}
			}

			return functiontool.New(
				functiontool.Config{
					Name:          toolDef.Name,
					Description:   toolDef.Description,
					InputSchema:   inputSchema,
					IsLongRunning: true,
				},
				func(ctx agent.Context, args mcp.ToolRobotSwitchInput) (*mcp.ToolRobotSwitchOutput, error) {
					return rt.ExecuteRobotSwitch(ctx, args)
				},
			)
		},
	}
}

func (rt *robotTools) ExecuteRobotSwitch(ctx context.Context, args mcp.ToolRobotSwitchInput) (*mcp.ToolRobotSwitchOutput, error) {
	robotID, err := robot_ref.NewID(args.RobotId)
	if err != nil {
		if def, ok := rt.agentRegistry.Get(args.RobotId); ok && !def.Hidden {
			return &(mcp.ToolRobotSwitchOutput{
				Success: true,
				RobotId: args.RobotId,
			}), nil
		}

		return nil, fmt.Errorf("invalid robot ID: %s", args.RobotId)
	}

	_, err = rt.robotQuerier.Get(ctx, robotID)
	if err != nil {
		return nil, fmt.Errorf("robot not found: %s", args.RobotId)
	}

	return &(mcp.ToolRobotSwitchOutput{
		Success: true,
		RobotId: args.RobotId,
	}), nil
}

func (rt *robotTools) newSystemRobotToolCatalogTool() *Tool {
	toolDef := mcp.GetSystemRobotToolCatalogTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				func(ctx agent.Context, args map[string]any) (*mcp.ToolSystemRobotToolCatalogOutput, error) {
					return rt.ExecuteGetAllToolNames(ctx, args)
				},
			)
		},
		Handler: makeHandler(rt.ExecuteGetAllToolNames),
	}
}

func (rt *robotTools) ExecuteGetAllToolNames(ctx context.Context, args map[string]any) (*mcp.ToolSystemRobotToolCatalogOutput, error) {
	allTools, err := rt.registry.GetTools(ctx)
	if err != nil {
		return nil, err
	}

	result := mcp.ToolSystemRobotToolCatalogOutput{
		Tools: make([]mcp.ToolInfo, 0, len(allTools)),
	}
	for _, t := range allTools {
		result.Tools = append(result.Tools, mcp.ToolInfo{
			Name:                 t.Definition.Name,
			Description:          t.Definition.Description,
			RequiresConfirmation: t.Definition.RequiresConfirmation,
		})
	}

	return &result, nil
}

func (rt *robotTools) injectToolNamesEnum(ctx context.Context, schema *jsonschema.Schema, propertyName string) {
	if schema == nil || schema.Properties == nil {
		return
	}

	prop, ok := schema.Properties[propertyName]
	if !ok || prop.Items == nil {
		return
	}

	ids := rt.registry.AllToolIDs(ctx)
	prop.Items.Enum = dt.Map(ids, func(name string) any { return name })
}

func (rt *robotTools) validateToolNames(names []string) []string {
	var invalid []string
	for _, name := range names {
		if !rt.registry.HasTool(name) {
			invalid = append(invalid, name)
		}
	}
	return invalid
}

func (rt *robotTools) newRobotCreateTool() *Tool {
	toolDef := mcp.GetRobotCreateTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			inputSchema := toolDef.InputSchema
			rt.injectToolNamesEnum(ctx, inputSchema, "tools")

			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: inputSchema,
				},
				func(ctx agent.Context, args mcp.ToolRobotCreateInput) (*mcp.ToolRobotCreateOutput, error) {
					return rt.ExecuteCreateRobot(ctx, args)
				},
			)
		},
		Handler: makeHandler(rt.ExecuteCreateRobot),
	}
}

func (rt *robotTools) ExecuteCreateRobot(ctx context.Context, args mcp.ToolRobotCreateInput) (*mcp.ToolRobotCreateOutput, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, err
	}

	var validationErrors []string

	if len(args.Tools) > 0 {
		if invalidTools := rt.validateToolNames(args.Tools); len(invalidTools) > 0 {
			validationErrors = append(validationErrors, "invalid tool names: "+strings.Join(invalidTools, ", "))
		}
	}

	if len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation errors: %v", strings.Join(validationErrors, "; "))
	}

	model, err := rt.modelFactory.DefaultModel(ctx)
	if err != nil {
		return nil, err
	}
	if args.Model != nil {
		model, err = model_ref.ParseID(*args.Model)
		if err != nil {
			return nil, fmt.Errorf("invalid model ref: %w", err)
		}
	}
	if err := rt.modelFactory.EnsureModelAvailable(ctx, model); err != nil {
		return nil, err
	}

	opts := []robot_writer.Option{
		robot_writer.WithTools(dt.Map(args.Tools, func(t string) robot_ref.ToolName { return robot_ref.ToolName(t) })),
	}

	robot, err := rt.robotWriter.Create(ctx, args.Name, args.Description, args.Playbook, model, accountID, opts...)
	if err != nil {
		return nil, err
	}

	return &(mcp.ToolRobotCreateOutput{
		Id:   robot.ID.String(),
		Name: robot.Name,
	}), nil
}

func (rt *robotTools) newRobotListTool() *Tool {
	toolDef := mcp.GetRobotListTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				func(ctx agent.Context, args mcp.ToolRobotListInput) (*mcp.ToolRobotListOutput, error) {
					return rt.ExecuteListRobots(ctx, args)
				},
			)
		},
		Handler: makeHandler(rt.ExecuteListRobots),
	}
}

func (rt *robotTools) ExecuteListRobots(ctx context.Context, args mcp.ToolRobotListInput) (*mcp.ToolRobotListOutput, error) {
	limit := uint(20)
	if args.Limit != nil {
		limit = uint(*args.Limit)
	}

	params := pagination.NewPageParams(1, limit)

	result, err := rt.robotQuerier.List(ctx, params)
	if err != nil {
		return nil, err
	}

	items := make([]mcp.RobotItem, 0, len(result.Items))
	for _, r := range result.Items {
		desc := r.Description
		items = append(items, mcp.RobotItem{
			Id:          r.ID.String(),
			Name:        r.Name,
			Description: &desc,
		})
	}

	return &(mcp.ToolRobotListOutput{
		Robots: items,
		Total:  len(items),
	}), nil
}

func (rt *robotTools) newRobotGetTool() *Tool {
	toolDef := mcp.GetRobotGetTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				func(ctx agent.Context, args mcp.ToolRobotGetInput) (*mcp.ToolRobotGetOutput, error) {
					return rt.ExecuteGetRobot(ctx, args)
				},
			)
		},
		Handler: makeHandler(rt.ExecuteGetRobot),
	}
}

func (rt *robotTools) ExecuteGetRobot(ctx context.Context, args mcp.ToolRobotGetInput) (*mcp.ToolRobotGetOutput, error) {
	robotID, err := robot_ref.NewID(args.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid robot ID: %s", args.Id)
	}

	robot, err := rt.robotQuerier.Get(ctx, robotID)
	if err != nil {
		return nil, err
	}

	desc := robot.Description

	return &(mcp.ToolRobotGetOutput{
		Id:          robot.ID.String(),
		Name:        robot.Name,
		Description: &desc,
		Playbook:    robot.Playbook,
		Model:       robot.Model.String(),
		Tools:       dt.Map(robot.Tools, func(t robot_ref.ToolName) string { return string(t) }),
	}), nil
}

func (rt *robotTools) newRobotUpdateTool() *Tool {
	toolDef := mcp.GetRobotUpdateTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			inputSchema := toolDef.InputSchema
			rt.injectToolNamesEnum(ctx, inputSchema, "tools")

			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: inputSchema,
				},
				func(ctx agent.Context, args mcp.ToolRobotUpdateInput) (*mcp.ToolRobotUpdateOutput, error) {
					return rt.ExecuteUpdateRobot(ctx, args)
				},
			)
		},
		Handler: makeHandler(rt.ExecuteUpdateRobot),
	}
}

func (rt *robotTools) ExecuteUpdateRobot(ctx context.Context, args mcp.ToolRobotUpdateInput) (*mcp.ToolRobotUpdateOutput, error) {
	robotID, err := robot_ref.NewID(args.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid robot ID: %s", args.Id)
	}

	var validationErrors []string

	if len(args.Tools) > 0 {
		if invalidTools := rt.validateToolNames(args.Tools); len(invalidTools) > 0 {
			validationErrors = append(validationErrors, "invalid tool names: "+strings.Join(invalidTools, ", "))
		}
	}

	if len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation errors: %s", strings.Join(validationErrors, "; "))
	}

	opts := []robot_writer.Option{}
	if args.Name != nil {
		opts = append(opts, robot_writer.WithName(*args.Name))
	}
	if args.Description != nil {
		opts = append(opts, robot_writer.WithDescription(*args.Description))
	}
	if args.Playbook != nil {
		opts = append(opts, robot_writer.WithPlaybook(*args.Playbook))
	}
	if args.Model != nil {
		model, err := model_ref.ParseID(*args.Model)
		if err != nil {
			return nil, fmt.Errorf("invalid model ref: %w", err)
		}
		if err := rt.modelFactory.EnsureModelAvailable(ctx, model); err != nil {
			return nil, err
		}
		opts = append(opts, robot_writer.WithModel(model))
	}
	if len(args.Tools) > 0 {
		opts = append(opts, robot_writer.WithTools(dt.Map(args.Tools, func(t string) robot_ref.ToolName { return robot_ref.ToolName(t) })))
	}

	robot, err := rt.robotWriter.Update(ctx, robotID, opts...)
	if err != nil {
		return nil, err
	}

	return &(mcp.ToolRobotUpdateOutput{
		Id:   robot.ID.String(),
		Name: robot.Name,
	}), nil
}

func (rt *robotTools) newRobotDeleteTool() *Tool {
	toolDef := mcp.GetRobotDeleteTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:                toolDef.Name,
					Description:         toolDef.Description,
					InputSchema:         toolDef.InputSchema,
					RequireConfirmation: toolDef.RequiresConfirmation && !confirmationDisabled(ctx),
				},
				func(ctx agent.Context, args mcp.ToolRobotDeleteInput) (*mcp.ToolRobotDeleteOutput, error) {
					return rt.ExecuteDeleteRobot(ctx, args)
				},
			)
		},
		Handler: makeHandler(rt.ExecuteDeleteRobot),
	}
}

func (rt *robotTools) ExecuteDeleteRobot(ctx context.Context, args mcp.ToolRobotDeleteInput) (*mcp.ToolRobotDeleteOutput, error) {
	robotID, err := robot_ref.NewID(args.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid robot ID: %s", args.Id)
	}

	err = rt.robotWriter.Delete(ctx, robotID)
	if err != nil {
		return nil, err
	}

	return &(mcp.ToolRobotDeleteOutput{Success: true, Id: args.Id}), nil
}

func newThrowAnErrorTool() *Tool {
	def := &mcp.ToolDefinition{
		Name:        "throw_an_error",
		Description: "Always returns an error. Used for testing error handling in the agent pipeline.",
	}

	type input struct{}

	return &Tool{
		Definition: def,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        def.Name,
					Description: def.Description,
				},
				func(ctx agent.Context, _ input) (map[string]any, error) {
					return nil, fmt.Errorf("intentional tool error for testing")
				},
			)
		},
	}
}
