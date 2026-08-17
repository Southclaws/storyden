package tools

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	adkagent "google.golang.org/adk/v2/agent"
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
	"github.com/Southclaws/storyden/app/services/semdex/robot/searchscore"
	"github.com/Southclaws/storyden/lib/mcp"
)

type robotTools struct {
	logger           *slog.Logger
	robotQuerier     *robot_querier.Querier
	robotWriter      *robot_writer.Writer
	robotSessionRepo *robot_session.Repository
	registry         *Registry
	toolsets         ToolsetResolver
	modelFactory     *llm_provider.Factory
}

func newRobotTools(
	logger *slog.Logger,
	robotQuerier *robot_querier.Querier,
	robotWriter *robot_writer.Writer,
	robotSessionRepo *robot_session.Repository,
	registry *Registry,
	toolsetResolver ToolsetResolver,
	modelFactory *llm_provider.Factory,
) *robotTools {
	t := &robotTools{
		logger:           logger,
		robotQuerier:     robotQuerier,
		robotWriter:      robotWriter,
		robotSessionRepo: robotSessionRepo,
		registry:         registry,
		toolsets:         toolsetResolver,
		modelFactory:     modelFactory,
	}

	registry.Register(t.newRobotCreateTool())
	registry.Register(t.newRobotListTool())
	registry.Register(t.newRobotGetTool())
	registry.Register(t.newRobotUpdateTool())
	registry.Register(t.newRobotDeleteTool())
	registry.Register(t.newRobotSearchTool())
	registry.Register(newThrowAnErrorTool())

	return t
}

func (rt *robotTools) invalidToolNames(names []string) []string {
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
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				func(ctx adkagent.Context, args mcp.ToolRobotCreateInput) (*mcp.ToolRobotCreateOutput, error) {
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

	if err := rt.toolsets.Validate(ctx, args.Toolsets); err != nil {
		return nil, fmt.Errorf("invalid Toolsets: %w", err)
	}
	directTools, err := rt.toolsets.RemoveContainedTools(ctx, args.Tools, args.Toolsets)
	if err != nil {
		return nil, fmt.Errorf("compose Robot tools: %w", err)
	}
	if invalid := rt.invalidToolNames(directTools); len(invalid) > 0 {
		return nil, fmt.Errorf("invalid tool names: %s", strings.Join(invalid, ", "))
	}
	if err := rt.toolsets.ValidateDirectTools(ctx, directTools, args.Toolsets); err != nil {
		return nil, err
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
		robot_writer.WithTools(dt.Map(directTools, func(t string) robot_ref.ToolName { return robot_ref.ToolName(t) })),
		robot_writer.WithToolsets(dt.Map(args.Toolsets, func(t string) robot_ref.ToolsetRef { return robot_ref.ToolsetRef(t) })),
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
				func(ctx adkagent.Context, args mcp.ToolRobotListInput) (*mcp.ToolRobotListOutput, error) {
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
				func(ctx adkagent.Context, args mcp.ToolRobotGetInput) (*mcp.ToolRobotGetOutput, error) {
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
		Toolsets:    dt.Map(robot.Toolsets, func(t robot_ref.ToolsetRef) string { return string(t) }),
	}), nil
}

func (rt *robotTools) newRobotUpdateTool() *Tool {
	toolDef := mcp.GetRobotUpdateTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				func(ctx adkagent.Context, args mcp.ToolRobotUpdateInput) (*mcp.ToolRobotUpdateOutput, error) {
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
	if len(args.Tools) > 0 || len(args.Toolsets) > 0 {
		existing, err := rt.robotQuerier.Get(ctx, robotID)
		if err != nil {
			return nil, err
		}
		existingTools := dt.Map(existing.Tools, func(t robot_ref.ToolName) string { return string(t) })
		existingToolsets := dt.Map(existing.Toolsets, func(t robot_ref.ToolsetRef) string { return string(t) })
		finalTools := existingTools
		if len(args.Tools) > 0 {
			finalTools = args.Tools
		}
		finalToolsets := existingToolsets
		if len(args.Toolsets) > 0 {
			if err := rt.toolsets.Validate(ctx, args.Toolsets); err != nil {
				return nil, fmt.Errorf("invalid Toolsets: %w", err)
			}
			finalToolsets = args.Toolsets
		}

		addedToolsets := addedRobotConfigurationValues(existingToolsets, finalToolsets)
		if len(addedToolsets) > 0 {
			finalTools, err = rt.toolsets.RemoveContainedTools(ctx, finalTools, addedToolsets)
			if err != nil {
				return nil, fmt.Errorf("compose Robot tools: %w", err)
			}
		}
		if len(args.Tools) > 0 {
			if invalid := rt.invalidToolNames(addedRobotConfigurationValues(existingTools, finalTools)); len(invalid) > 0 {
				return nil, fmt.Errorf("invalid tool names: %s", strings.Join(invalid, ", "))
			}
		}
		if err := rt.toolsets.ValidateDirectTools(ctx, addedRobotConfigurationValues(existingTools, finalTools), finalToolsets); err != nil {
			return nil, err
		}

		if len(args.Tools) > 0 || len(addedToolsets) > 0 {
			opts = append(opts, robot_writer.WithTools(dt.Map(finalTools, func(t string) robot_ref.ToolName { return robot_ref.ToolName(t) })))
		}
		if len(args.Toolsets) > 0 {
			opts = append(opts, robot_writer.WithToolsets(dt.Map(finalToolsets, func(t string) robot_ref.ToolsetRef { return robot_ref.ToolsetRef(t) })))
		}
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

func addedRobotConfigurationValues(existing, requested []string) []string {
	seen := make(map[string]struct{}, len(existing))
	for _, value := range existing {
		seen[value] = struct{}{}
	}
	added := make([]string, 0, len(requested))
	for _, value := range requested {
		if _, ok := seen[value]; !ok {
			added = append(added, value)
		}
	}
	return added
}

func (rt *robotTools) newRobotSearchTool() *Tool {
	toolDef := mcp.GetRobotSearchTool()
	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(functiontool.Config{
				Name:        toolDef.Name,
				Description: toolDef.Description,
				InputSchema: toolDef.InputSchema,
			}, func(ctx adkagent.Context, args mcp.ToolRobotSearchInput) (*mcp.ToolRobotSearchOutput, error) {
				return rt.ExecuteSearchRobots(ctx, args)
			})
		},
		Handler: makeHandler(rt.ExecuteSearchRobots),
	}
}

func (rt *robotTools) ExecuteSearchRobots(ctx context.Context, args mcp.ToolRobotSearchInput) (*mcp.ToolRobotSearchOutput, error) {
	result, err := rt.robotQuerier.List(ctx, pagination.NewPageParams(1, 100))
	if err != nil {
		return nil, err
	}
	query := searchscore.New(args.Query)
	type match struct {
		robot mcp.RobotSearchResult
		score int
	}
	matches := make([]match, 0, len(result.Items))
	for _, candidate := range result.Items {
		score := query.Score(
			searchscore.Field{Text: candidate.Name, SubstringMatch: 40, TermMatch: 8},
			searchscore.Field{Text: candidate.Description, SubstringMatch: 20, TermMatch: 4},
			searchscore.Field{Text: candidate.Playbook, SubstringMatch: 5},
		)
		if query.Empty() || score > 0 {
			matches = append(matches, match{robot: mcp.RobotSearchResult{
				Id:          candidate.ID.String(),
				DelegateTo:  "robot_" + candidate.ID.String(),
				Name:        candidate.Name,
				Description: candidate.Description,
			}, score: score})
		}
	}
	sort.SliceStable(matches, func(i, j int) bool { return matches[i].score > matches[j].score })
	limit := 8
	if args.MaxResults != nil {
		limit = *args.MaxResults
	}
	if len(matches) > limit {
		matches = matches[:limit]
	}
	robots := make([]mcp.RobotSearchResult, 0, len(matches))
	for _, match := range matches {
		robots = append(robots, match.robot)
	}
	return &mcp.ToolRobotSearchOutput{
		Robots:     robots,
		NextAction: "Delegate a bounded task to the best matching Robot when its specialised playbook or Toolsets add value; otherwise continue directly.",
	}, nil
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
				func(ctx adkagent.Context, args mcp.ToolRobotDeleteInput) (*mcp.ToolRobotDeleteOutput, error) {
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
		Title:       "Throw an Error",
		Description: "Always returns an error. Used for testing error handling in the agent pipeline.",
		InputSchema: &jsonschema.Schema{Type: "object"},
		Annotations: mcp.ToolAnnotations{
			ReadOnlyHint:   true,
			IdempotentHint: true,
		},
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
				func(ctx adkagent.Context, _ input) (map[string]any, error) {
					return nil, fmt.Errorf("intentional tool error for testing")
				},
			)
		},
	}
}
