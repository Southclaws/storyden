package tools

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/robot/robot_querier"
	"github.com/Southclaws/storyden/app/resources/robot/robot_ref"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/app/resources/robot/robot_writer"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/mcp"
)

type robotTools struct {
	logger           *slog.Logger
	robotQuerier     *robot_querier.Querier
	robotWriter      *robot_writer.Writer
	robotSessionRepo *robot_session.Repository
	registry         *Registry
}

func newRobotTools(
	logger *slog.Logger,
	robotQuerier *robot_querier.Querier,
	robotWriter *robot_writer.Writer,
	robotSessionRepo *robot_session.Repository,
	registry *Registry,
) *robotTools {
	t := &robotTools{
		logger:           logger,
		robotQuerier:     robotQuerier,
		robotWriter:      robotWriter,
		robotSessionRepo: robotSessionRepo,
		registry:         registry,
	}

	registry.Register(t.newRobotSwitchTool())
	registry.Register(t.newRobotGetAllToolNamesTool())
	registry.Register(t.newRobotCreateTool())
	registry.Register(t.newRobotListTool())
	registry.Register(t.newRobotGetTool())
	registry.Register(t.newRobotUpdateTool())
	registry.Register(t.newRobotDeleteTool())

	return t
}

func (rt *robotTools) newRobotSwitchTool() *Tool {
	toolDef := mcp.GetRobotSwitchTool()

	return &Tool{
		Definition:   toolDef,
		IsClientTool: true,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			result, err := rt.robotQuerier.List(ctx, pagination.Parameters{})
			if err != nil {
				return nil, err
			}

			robotIDs := make([]any, len(result.Items))
			for i, r := range result.Items {
				robotIDs[i] = r.ID.String()
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
				rt.ExecuteRobotSwitch,
			)
		},
	}
}

func (rt *robotTools) ExecuteRobotSwitch(ctx tool.Context, args mcp.ToolRobotSwitchInput) (*mcp.ToolRobotSwitchOutput, error) {
	robotID, err := robot_ref.NewID(args.RobotId)
	if err != nil {
		return nil, fmt.Errorf("invalid robot ID: " + args.RobotId)
	}

	_, err = rt.robotQuerier.Get(ctx, robotID)
	if err != nil {
		return nil, fmt.Errorf("robot not found: " + args.RobotId)
	}

	return &(mcp.ToolRobotSwitchOutput{
		Success: true,
		RobotId: args.RobotId,
	}), nil
}

func (rt *robotTools) newRobotGetAllToolNamesTool() *Tool {
	toolDef := mcp.GetRobotGetAllToolNamesTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				rt.ExecuteGetAllToolNames,
			)
		},
	}
}

func (rt *robotTools) ExecuteGetAllToolNames(ctx tool.Context, args map[string]any) (*mcp.ToolRobotGetAllToolNamesOutput, error) {
	allTools, err := rt.registry.GetTools(ctx)
	if err != nil {
		return nil, err
	}

	result := mcp.ToolRobotGetAllToolNamesOutput{
		Tools: make([]mcp.ToolInfo, 0, len(allTools)),
	}
	for _, t := range allTools {
		result.Tools = append(result.Tools, mcp.ToolInfo{
			Name:        t.Definition.Name,
			Description: t.Definition.Description,
		})
	}

	return &result, nil
}

func (rt *robotTools) newRobotCreateTool() *Tool {
	toolDef := mcp.GetRobotCreateTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			inputSchema := toolDef.InputSchema
			mcp.InjectToolNamesEnum(inputSchema, "tools")

			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: inputSchema,
				},
				rt.ExecuteCreateRobot,
			)
		},
	}
}

func (rt *robotTools) ExecuteCreateRobot(ctx tool.Context, args mcp.ToolRobotCreateInput) (*mcp.ToolRobotCreateOutput, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, err
	}

	var validationErrors []string

	if len(args.Tools) > 0 {
		allToolNames := mcp.AllToolNames()
		toolNameSet := make(map[string]bool)
		for _, name := range allToolNames {
			toolNameSet[name] = true
		}

		var invalidTools []string
		for _, t := range args.Tools {
			if !toolNameSet[t] {
				invalidTools = append(invalidTools, t)
			}
		}
		if len(invalidTools) > 0 {
			validationErrors = append(validationErrors, "invalid tool names: "+strings.Join(invalidTools, ", "))
		}
	}

	if len(validationErrors) > 0 {
		return nil, fmt.Errorf(strings.Join(validationErrors, "; "))
	}

	opts := []robot_writer.Option{
		robot_writer.WithMeta(map[string]any{
			"tools": args.Tools,
		}),
	}

	robot, err := rt.robotWriter.Create(ctx, args.Name, args.Description, args.Playbook, accountID, opts...)
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
				rt.ExecuteListRobots,
			)
		},
	}
}

func (rt *robotTools) ExecuteListRobots(ctx tool.Context, args mcp.ToolRobotListInput) (*mcp.ToolRobotListOutput, error) {
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
				rt.ExecuteGetRobot,
			)
		},
	}
}

func (rt *robotTools) ExecuteGetRobot(ctx tool.Context, args mcp.ToolRobotGetInput) (*mcp.ToolRobotGetOutput, error) {
	robotID, err := robot_ref.NewID(args.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid robot ID: " + args.Id)
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
		Tools:       robot.Tools,
	}), nil
}

func (rt *robotTools) newRobotUpdateTool() *Tool {
	toolDef := mcp.GetRobotUpdateTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			inputSchema := toolDef.InputSchema
			mcp.InjectToolNamesEnum(inputSchema, "tools")

			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: inputSchema,
				},
				rt.ExecuteUpdateRobot,
			)
		},
	}
}

func (rt *robotTools) ExecuteUpdateRobot(ctx tool.Context, args mcp.ToolRobotUpdateInput) (*mcp.ToolRobotUpdateOutput, error) {
	robotID, err := robot_ref.NewID(args.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid robot ID: " + args.Id)
	}

	var validationErrors []string

	if len(args.Tools) > 0 {
		allToolNames := mcp.AllToolNames()
		toolNameSet := make(map[string]bool)
		for _, name := range allToolNames {
			toolNameSet[name] = true
		}

		var invalidTools []string
		for _, t := range args.Tools {
			if !toolNameSet[t] {
				invalidTools = append(invalidTools, t)
			}
		}
		if len(invalidTools) > 0 {
			validationErrors = append(validationErrors, "invalid tool names: "+strings.Join(invalidTools, ", "))
		}
	}

	if len(validationErrors) > 0 {
		return nil, fmt.Errorf(strings.Join(validationErrors, "; "))
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
	if len(args.Tools) > 0 {
		opts = append(opts, robot_writer.WithMeta(map[string]any{
			"tools": args.Tools,
		}))
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
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				rt.ExecuteDeleteRobot,
			)
		},
	}
}

func (rt *robotTools) ExecuteDeleteRobot(ctx tool.Context, args mcp.ToolRobotDeleteInput) (*mcp.ToolRobotDeleteOutput, error) {
	robotID, err := robot_ref.NewID(args.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid robot ID: " + args.Id)
	}

	err = rt.robotWriter.Delete(ctx, robotID)
	if err != nil {
		return nil, err
	}

	return &(mcp.ToolRobotDeleteOutput{Success: true, Id: args.Id}), nil
}
