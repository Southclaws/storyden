package trail_tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"
	adkagent "google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"

	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/recurrence"
	"github.com/Southclaws/storyden/app/resources/trail"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	robottools "github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	"github.com/Southclaws/storyden/app/services/trail/trail_manager"
	"github.com/Southclaws/storyden/lib/mcp"
)

type trailTools struct {
	manager *trail_manager.Manager
}

func makeHandler[T, O any](execute func(context.Context, T) (*O, error)) robottools.Handler {
	return func(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
		var input T
		if err := json.Unmarshal(args, &input); err != nil {
			return nil, err
		}
		output, err := execute(ctx, input)
		if err != nil {
			return nil, err
		}
		return json.Marshal(output)
	}
}

func New(manager *trail_manager.Manager, registry *robottools.Registry) *trailTools {
	t := &trailTools{manager: manager}

	registry.Register(t.newCreateTool())
	registry.Register(t.newListTool())
	registry.Register(t.newGetTool())
	registry.Register(t.newUpdateTool())
	registry.Register(t.newSchedulePreviewTool())
	registry.Register(t.newRunListTool())
	registry.Register(t.newRunGetTool())
	registry.Register(t.newRunCreateTool())
	registry.Register(t.newActionRunCancelTool())

	return t
}

func (t *trailTools) newCreateTool() *robottools.Tool {
	toolDef := mcp.GetTrailCreateTool()

	return &robottools.Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				func(ctx adkagent.Context, args mcp.ToolTrailCreateInput) (*mcp.ToolTrailCreateOutput, error) {
					return t.ExecuteCreate(ctx, args)
				},
			)
		},
		Handler: makeHandler(t.ExecuteCreate),
	}
}

func (t *trailTools) ExecuteCreate(ctx context.Context, args mcp.ToolTrailCreateInput) (*mcp.ToolTrailCreateOutput, error) {
	if err := authoriseTrailTools(ctx); err != nil {
		return nil, err
	}

	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	status, err := trail.NewStatus(string(args.Status))
	if err != nil {
		return nil, invalidTrailToolInput(ctx, err)
	}

	trigger, err := deserialiseTrailToolTrigger(args.Trigger)
	if err != nil {
		return nil, invalidTrailToolInput(ctx, err)
	}

	actions, err := deserialiseTrailToolActions(ctx, args.Action)
	if err != nil {
		return nil, invalidTrailToolInput(ctx, err)
	}

	created, err := t.manager.Create(
		ctx,
		xid.ID(accountID),
		args.Name,
		optionalString(args.Description),
		status,
		trigger,
		actions,
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result, err := serialiseTrailToolItem(created)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &mcp.ToolTrailCreateOutput{
		Trail:      result,
		Message:    trailToolConfirmation("Created", created),
		NextAction: "Use trail_get to inspect this Trail or trail_update to change its schedule, actions, or lifecycle state.",
	}, nil
}

func (t *trailTools) newListTool() *robottools.Tool {
	toolDef := mcp.GetTrailListTool()

	return &robottools.Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				func(ctx adkagent.Context, _ map[string]any) (*mcp.ToolTrailListOutput, error) {
					return t.ExecuteList(ctx, map[string]any{})
				},
			)
		},
		Handler: makeHandler(t.ExecuteList),
	}
}

func (t *trailTools) ExecuteList(ctx context.Context, _ map[string]any) (*mcp.ToolTrailListOutput, error) {
	if err := authoriseTrailTools(ctx); err != nil {
		return nil, err
	}

	items, err := t.manager.List(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	trails, err := dt.MapErr(items, serialiseTrailToolItem)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &mcp.ToolTrailListOutput{Trails: trails, Total: len(trails)}, nil
}

func (t *trailTools) newGetTool() *robottools.Tool {
	toolDef := mcp.GetTrailGetTool()

	return &robottools.Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				func(ctx adkagent.Context, args mcp.ToolTrailGetInput) (*mcp.ToolTrailGetOutput, error) {
					return t.ExecuteGet(ctx, args)
				},
			)
		},
		Handler: makeHandler(t.ExecuteGet),
	}
}

func (t *trailTools) ExecuteGet(ctx context.Context, args mcp.ToolTrailGetInput) (*mcp.ToolTrailGetOutput, error) {
	if err := authoriseTrailTools(ctx); err != nil {
		return nil, err
	}

	id, err := parseTrailToolID(args.TrailId)
	if err != nil {
		return nil, invalidTrailToolInput(ctx, err)
	}

	item, err := t.manager.Get(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result, err := serialiseTrailToolItem(item)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &mcp.ToolTrailGetOutput{Trail: result}, nil
}

func (t *trailTools) newUpdateTool() *robottools.Tool {
	toolDef := mcp.GetTrailUpdateTool()

	return &robottools.Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				func(ctx adkagent.Context, args mcp.ToolTrailUpdateInput) (*mcp.ToolTrailUpdateOutput, error) {
					return t.ExecuteUpdate(ctx, args)
				},
			)
		},
		Handler: makeHandler(t.ExecuteUpdate),
	}
}

func (t *trailTools) ExecuteUpdate(ctx context.Context, args mcp.ToolTrailUpdateInput) (*mcp.ToolTrailUpdateOutput, error) {
	if err := authoriseTrailTools(ctx); err != nil {
		return nil, err
	}

	id, err := parseTrailToolID(args.TrailId)
	if err != nil {
		return nil, invalidTrailToolInput(ctx, err)
	}

	status, err := trail.NewStatus(string(args.Status))
	if err != nil {
		return nil, invalidTrailToolInput(ctx, err)
	}

	trigger, err := deserialiseTrailToolTrigger(args.Trigger)
	if err != nil {
		return nil, invalidTrailToolInput(ctx, err)
	}

	actions, err := deserialiseTrailToolActions(ctx, args.Action)
	if err != nil {
		return nil, invalidTrailToolInput(ctx, err)
	}

	updated, err := t.manager.Update(
		ctx,
		id,
		args.Name,
		optionalString(args.Description),
		status,
		trigger,
		actions,
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result, err := serialiseTrailToolItem(updated)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &mcp.ToolTrailUpdateOutput{
		Trail:      result,
		Message:    trailToolConfirmation("Updated", updated),
		NextAction: "Use trail_get to confirm the persisted Trail definition after any later change.",
	}, nil
}

func (t *trailTools) newSchedulePreviewTool() *robottools.Tool {
	toolDef := mcp.GetTrailSchedulePreviewTool()

	return &robottools.Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				func(ctx adkagent.Context, args mcp.ToolTrailSchedulePreviewInput) (*mcp.ToolTrailSchedulePreviewOutput, error) {
					return t.ExecuteSchedulePreview(ctx, args)
				},
			)
		},
		Handler: makeHandler(t.ExecuteSchedulePreview),
	}
}

func (t *trailTools) ExecuteSchedulePreview(ctx context.Context, args mcp.ToolTrailSchedulePreviewInput) (*mcp.ToolTrailSchedulePreviewOutput, error) {
	if err := authoriseTrailTools(ctx); err != nil {
		return nil, err
	}

	schedule, err := deserialiseTrailToolSchedule(args.Schedule)
	if err != nil {
		return nil, invalidTrailToolInput(ctx, err)
	}

	after := time.Now().UTC()
	if args.After != nil {
		after = args.After.UTC()
	}

	return &mcp.ToolTrailSchedulePreviewOutput{
		Occurrences: t.manager.Preview(schedule, after),
	}, nil
}

func (t *trailTools) newRunListTool() *robottools.Tool {
	toolDef := mcp.GetTrailRunListTool()

	return &robottools.Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{Name: toolDef.Name, Description: toolDef.Description, InputSchema: toolDef.InputSchema},
				func(ctx adkagent.Context, args mcp.ToolTrailRunListInput) (*mcp.ToolTrailRunListOutput, error) {
					return t.ExecuteRunList(ctx, args)
				},
			)
		},
		Handler: makeHandler(t.ExecuteRunList),
	}
}

func (t *trailTools) ExecuteRunList(ctx context.Context, args mcp.ToolTrailRunListInput) (*mcp.ToolTrailRunListOutput, error) {
	if err := authoriseTrailTools(ctx); err != nil {
		return nil, err
	}

	id, err := parseTrailToolID(args.TrailId)
	if err != nil {
		return nil, invalidTrailToolInput(ctx, err)
	}

	limit := 20
	if args.Limit != nil {
		limit = *args.Limit
	}

	runs, err := t.manager.ListRunsLimited(ctx, id, limit+1)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	hasMore := len(runs) > limit
	if hasMore {
		runs = runs[:limit]
	}

	result := dt.Map(runs, serialiseTrailToolRunSummary)
	return &mcp.ToolTrailRunListOutput{
		Runs:       result,
		Returned:   len(result),
		HasMore:    hasMore,
		NextAction: "Use trail_run_get with a run ID to inspect its trigger context and independent action results.",
	}, nil
}

func (t *trailTools) newRunGetTool() *robottools.Tool {
	toolDef := mcp.GetTrailRunGetTool()

	return &robottools.Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{Name: toolDef.Name, Description: toolDef.Description, InputSchema: toolDef.InputSchema},
				func(ctx adkagent.Context, args mcp.ToolTrailRunGetInput) (*mcp.ToolTrailRunGetOutput, error) {
					return t.ExecuteRunGet(ctx, args)
				},
			)
		},
		Handler: makeHandler(t.ExecuteRunGet),
	}
}

func (t *trailTools) ExecuteRunGet(ctx context.Context, args mcp.ToolTrailRunGetInput) (*mcp.ToolTrailRunGetOutput, error) {
	if err := authoriseTrailTools(ctx); err != nil {
		return nil, err
	}

	trailID, err := parseTrailToolID(args.TrailId)
	if err != nil {
		return nil, invalidTrailToolInput(ctx, err)
	}
	runID, err := parseTrailToolRunID(args.RunId)
	if err != nil {
		return nil, invalidTrailToolInput(ctx, err)
	}

	run, err := t.manager.GetRun(ctx, trailID, runID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result, err := serialiseTrailToolRun(run)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &mcp.ToolTrailRunGetOutput{
		Run:        result,
		Message:    trailToolRunMessage(run),
		NextAction: trailToolRunNextAction(run),
	}, nil
}

func (t *trailTools) newRunCreateTool() *robottools.Tool {
	toolDef := mcp.GetTrailRunCreateTool()

	return &robottools.Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:                toolDef.Name,
					Description:         toolDef.Description,
					InputSchema:         toolDef.InputSchema,
					RequireConfirmation: toolDef.RequiresConfirmation && !robottools.ConfirmationDisabled(ctx),
				},
				func(ctx adkagent.Context, args mcp.ToolTrailRunCreateInput) (*mcp.ToolTrailRunCreateOutput, error) {
					return t.ExecuteRunCreate(ctx, args)
				},
			)
		},
		Handler: makeHandler(t.ExecuteRunCreate),
	}
}

func (t *trailTools) ExecuteRunCreate(ctx context.Context, args mcp.ToolTrailRunCreateInput) (*mcp.ToolTrailRunCreateOutput, error) {
	if err := authoriseTrailTools(ctx); err != nil {
		return nil, err
	}

	id, err := parseTrailToolID(args.TrailId)
	if err != nil {
		return nil, invalidTrailToolInput(ctx, err)
	}
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	run, err := t.manager.RunNow(ctx, id, xid.ID(accountID))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	result, err := serialiseTrailToolRun(run)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &mcp.ToolTrailRunCreateOutput{
		Run:        result,
		Message:    fmt.Sprintf("Started manual run %s without changing the Trail's schedule.", run.ID),
		NextAction: "Use trail_run_get to inspect progress and action results.",
	}, nil
}

func (t *trailTools) newActionRunCancelTool() *robottools.Tool {
	toolDef := mcp.GetTrailActionRunCancelTool()

	return &robottools.Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:                toolDef.Name,
					Description:         toolDef.Description,
					InputSchema:         toolDef.InputSchema,
					RequireConfirmation: toolDef.RequiresConfirmation && !robottools.ConfirmationDisabled(ctx),
				},
				func(ctx adkagent.Context, args mcp.ToolTrailActionRunCancelInput) (*mcp.ToolTrailActionRunCancelOutput, error) {
					return t.ExecuteActionRunCancel(ctx, args)
				},
			)
		},
		Handler: makeHandler(t.ExecuteActionRunCancel),
	}
}

func (t *trailTools) ExecuteActionRunCancel(ctx context.Context, args mcp.ToolTrailActionRunCancelInput) (*mcp.ToolTrailActionRunCancelOutput, error) {
	if err := authoriseTrailTools(ctx); err != nil {
		return nil, err
	}

	trailID, err := parseTrailToolID(args.TrailId)
	if err != nil {
		return nil, invalidTrailToolInput(ctx, err)
	}
	runID, err := parseTrailToolRunID(args.RunId)
	if err != nil {
		return nil, invalidTrailToolInput(ctx, err)
	}
	actionRunID, err := parseTrailToolActionRunID(args.ActionRunId)
	if err != nil {
		return nil, invalidTrailToolInput(ctx, err)
	}

	if err := t.manager.Cancel(ctx, trailID, runID, actionRunID); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	run, err := t.manager.GetRun(ctx, trailID, runID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	result, err := serialiseTrailToolRun(run)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &mcp.ToolTrailActionRunCancelOutput{
		Run:        result,
		Message:    fmt.Sprintf("Cancelled action run %s.", actionRunID),
		NextAction: "Inspect the other independent action results before deciding whether to start a new manual run.",
	}, nil
}

func authoriseTrailTools(ctx context.Context) error {
	return session.Authorise(ctx, nil, rbac.PermissionManageTrails)
}

func invalidTrailToolInput(ctx context.Context, err error) error {
	return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
}

func parseTrailToolID(value string) (trail.ID, error) {
	id, err := xid.FromString(value)
	if err != nil {
		return trail.ID{}, fmt.Errorf("invalid Trail ID: %w", err)
	}
	return trail.ID(id), nil
}

func parseTrailToolRunID(value string) (trail.RunID, error) {
	id, err := xid.FromString(value)
	if err != nil {
		return trail.RunID{}, fmt.Errorf("invalid Trail run ID: %w", err)
	}
	return trail.RunID(id), nil
}

func parseTrailToolActionRunID(value string) (trail.ActionRunID, error) {
	id, err := xid.FromString(value)
	if err != nil {
		return trail.ActionRunID{}, fmt.Errorf("invalid Trail action run ID: %w", err)
	}
	return trail.ActionRunID(id), nil
}

func optionalString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func deserialiseTrailToolTrigger(input mcp.TrailToolTriggerYaml) (trail.Trigger, error) {
	switch input.Type {
	case mcp.TrailToolTriggerYamlTypeSchedule:
		if input.Schedule == nil || len(input.Event) > 0 {
			return trail.Trigger{}, trail.ErrInvalidScheduleTrigger
		}

		schedule, err := deserialiseTrailToolSchedule(*input.Schedule)
		if err != nil {
			return trail.Trigger{}, err
		}
		return trail.NewScheduleTrigger(*schedule)

	case mcp.TrailToolTriggerYamlTypeEvent:
		if input.Schedule != nil {
			return trail.Trigger{}, trail.ErrInvalidEventTrigger
		}
		return trail.NewEventTrigger(input.Event)

	default:
		return trail.Trigger{}, trail.ErrUnsupportedTrigger
	}
}

func deserialiseTrailToolSchedule(input mcp.TrailToolScheduleYaml) (*recurrence.Schedule, error) {
	frequency, err := recurrence.NewFrequency(string(input.Rule.Frequency))
	if err != nil {
		return nil, err
	}

	weekdays, err := dt.MapErr(input.Rule.Weekday, func(value mcp.TrailToolRecurrenceRuleYamlWeekdayElem) (recurrence.Weekday, error) {
		return recurrence.NewWeekday(string(value))
	})
	if err != nil {
		return nil, err
	}

	return recurrence.Compile(recurrence.Schedule{
		Start:    input.Start,
		Timezone: input.Timezone,
		Rule: recurrence.Rule{
			Frequency:  frequency,
			Interval:   input.Rule.Interval,
			ByWeekday:  weekdays,
			ByMonth:    input.Rule.Month,
			ByMonthDay: input.Rule.MonthDay,
			Count:      input.Rule.Count,
		},
	})
}

func deserialiseTrailToolActions(ctx context.Context, input []mcp.TrailToolActionYaml) ([]trail.ActionSpec, error) {
	return dt.MapErr(input, func(action mcp.TrailToolActionYaml) (trail.ActionSpec, error) {
		if action.Type != trail.ActionKindRobotRun.String() {
			return trail.ActionSpec{}, fmt.Errorf("unsupported Trail action type %q", action.Type)
		}

		robotRef := ""
		if action.RobotRef != nil {
			robotRef = strings.TrimSpace(*action.RobotRef)
		}
		if robotRef == "" {
			run := robottools.RunContextFromContext(ctx)
			robotRef = strings.TrimSpace(run.RobotRef)
			if robotRef == "" {
				current, ok := run.RobotID.Get()
				if !ok {
					return trail.ActionSpec{}, fault.New("robot_ref is required outside a Robot conversation")
				}
				robotRef = current.String()
			}
		}

		config, err := json.Marshal(trail.RobotRunConfig{
			Type:        trail.ActionKindRobotRun,
			RobotRef:    robotRef,
			Instruction: action.Instruction,
		})
		if err != nil {
			return trail.ActionSpec{}, err
		}

		return trail.ActionSpec{Kind: trail.ActionKindRobotRun, Config: config}, nil
	})
}

func serialiseTrailToolItem(input *trail.Trail) (mcp.TrailToolItemYaml, error) {
	trigger, err := serialiseTrailToolTrigger(input.Trigger)
	if err != nil {
		return mcp.TrailToolItemYaml{}, err
	}

	actions, err := dt.MapErr(input.Actions, serialiseTrailToolAction)
	if err != nil {
		return mcp.TrailToolItemYaml{}, err
	}

	return mcp.TrailToolItemYaml{
		Id:               input.ID.String(),
		Name:             input.Name,
		Description:      input.Description,
		Status:           mcp.TrailToolItemYamlStatus(input.Status.String()),
		Trigger:          trigger,
		Action:           actions,
		NextOccurrenceAt: input.NextOccurrenceAt,
		LastOccurrenceAt: input.LastOccurrenceAt,
		CreatedAt:        input.CreatedAt,
		UpdatedAt:        input.UpdatedAt,
	}, nil
}

func serialiseTrailToolTrigger(input trail.Trigger) (mcp.TrailToolTriggerYaml, error) {
	switch input.Type() {
	case trail.TriggerTypeSchedule:
		schedule := input.Schedule()
		if schedule == nil {
			return mcp.TrailToolTriggerYaml{}, trail.ErrInvalidScheduleTrigger
		}
		return mcp.TrailToolTriggerYaml{
			Type:     mcp.TrailToolTriggerYamlTypeSchedule,
			Schedule: serialiseTrailToolSchedule(schedule),
		}, nil

	case trail.TriggerTypeEvent:
		event := input.Event()
		if event == nil {
			return mcp.TrailToolTriggerYaml{}, trail.ErrInvalidEventTrigger
		}
		return mcp.TrailToolTriggerYaml{
			Type:  mcp.TrailToolTriggerYamlTypeEvent,
			Event: event.Events,
		}, nil

	default:
		return mcp.TrailToolTriggerYaml{}, trail.ErrUnsupportedTrigger
	}
}

func serialiseTrailToolSchedule(input *recurrence.Schedule) *mcp.TrailToolScheduleYaml {
	weekdays := dt.Map(input.Rule.ByWeekday, func(value recurrence.Weekday) mcp.TrailToolRecurrenceRuleYamlWeekdayElem {
		return mcp.TrailToolRecurrenceRuleYamlWeekdayElem(value.String())
	})

	return &mcp.TrailToolScheduleYaml{
		Start:    input.Start,
		Timezone: input.Timezone,
		Rule: mcp.TrailToolRecurrenceRuleYaml{
			Frequency: mcp.TrailToolRecurrenceRuleYamlFrequency(input.Rule.Frequency.String()),
			Interval:  input.Rule.Interval,
			Weekday:   weekdays,
			Month:     input.Rule.ByMonth,
			MonthDay:  input.Rule.ByMonthDay,
			Count:     input.Rule.Count,
		},
	}
}

func serialiseTrailToolAction(input *trail.Action) (mcp.TrailToolActionYaml, error) {
	if input.Kind != trail.ActionKindRobotRun {
		return mcp.TrailToolActionYaml{}, fmt.Errorf("unsupported Trail action type %s", input.Kind)
	}

	var config trail.RobotRunConfig
	if err := json.Unmarshal(input.Config, &config); err != nil {
		return mcp.TrailToolActionYaml{}, err
	}

	robotRef := config.RobotRef
	return mcp.TrailToolActionYaml{
		Type:        trail.ActionKindRobotRun.String(),
		RobotRef:    &robotRef,
		Instruction: config.Instruction,
	}, nil
}

func serialiseTrailToolRunSummary(input *trail.Run) mcp.TrailToolRunSummaryYaml {
	kind := input.Kind
	if input.Trigger != nil {
		kind = input.Trigger.Kind
	}

	counts := mcp.TrailToolActionRunCountsYaml{}
	for _, action := range input.ActionRuns {
		switch action.Status {
		case trail.ActionRunStatusQueued:
			counts.Queued++
		case trail.ActionRunStatusRunning:
			counts.Running++
		case trail.ActionRunStatusCompleted:
			counts.Completed++
		case trail.ActionRunStatusBlocked:
			counts.Blocked++
		case trail.ActionRunStatusFailed:
			counts.Failed++
		case trail.ActionRunStatusCancelled:
			counts.Cancelled++
		}
	}

	return mcp.TrailToolRunSummaryYaml{
		Id:           input.ID.String(),
		TrailId:      input.TrailID.String(),
		Kind:         mcp.TrailToolRunSummaryYamlKind(kind.String()),
		Status:       mcp.TrailToolRunSummaryYamlStatus(input.Status.String()),
		ActionStatus: counts,
		ScheduledFor: input.ScheduledFor,
		FinishedAt:   input.FinishedAt,
		CreatedAt:    input.CreatedAt,
		UpdatedAt:    input.UpdatedAt,
	}
}

func serialiseTrailToolRun(input *trail.Run) (mcp.TrailToolRunYaml, error) {
	actions, err := dt.MapErr(input.ActionRuns, serialiseTrailToolActionRun)
	if err != nil {
		return mcp.TrailToolRunYaml{}, err
	}

	var trigger *mcp.TrailToolTriggerSnapshotYaml
	if input.Trigger != nil {
		serialised, err := serialiseTrailToolTriggerSnapshot(input.Trigger)
		if err != nil {
			return mcp.TrailToolRunYaml{}, err
		}
		trigger = &serialised
	}

	return mcp.TrailToolRunYaml{
		Summary: serialiseTrailToolRunSummary(input),
		Trigger: trigger,
		Action:  actions,
	}, nil
}

func serialiseTrailToolTriggerSnapshot(input *trail.TriggerEvent) (mcp.TrailToolTriggerSnapshotYaml, error) {
	trigger, err := serialiseTrailToolTrigger(input.Trigger)
	if err != nil {
		return mcp.TrailToolTriggerSnapshotYaml{}, err
	}

	var eventName *string
	if input.EventName != "" {
		eventName = &input.EventName
	}

	var payload *string
	if len(input.Payload) > 0 {
		if !json.Valid(input.Payload) {
			return mcp.TrailToolTriggerSnapshotYaml{}, fmt.Errorf("Trail run has an invalid event payload")
		}
		encoded := string(input.Payload)
		payload = &encoded
	}

	var initiatedBy *string
	if input.InitiatedBy != "" {
		initiatedBy = &input.InitiatedBy
	}

	return mcp.TrailToolTriggerSnapshotYaml{
		Kind:             mcp.TrailToolTriggerSnapshotYamlKind(input.Kind.String()),
		Trigger:          trigger,
		EventName:        eventName,
		EventPayloadJson: payload,
		ScheduledFor:     input.ScheduledFor,
		ObservedAt:       input.ObservedAt,
		InitiatedBy:      initiatedBy,
	}, nil
}

func serialiseTrailToolActionRun(input *trail.ActionRun) (mcp.TrailToolActionRunYaml, error) {
	action := input.Action
	if action == nil {
		action = &trail.Action{ID: input.ActionID, Kind: input.Kind, Config: input.Config}
	}
	actionConfig, err := serialiseTrailToolAction(action)
	if err != nil {
		return mcp.TrailToolActionRunYaml{}, err
	}

	var robotSessionID *string
	if len(input.Target) > 0 {
		var target trail.RobotInvocation
		if err := json.Unmarshal(input.Target, &target); err != nil {
			return mcp.TrailToolActionRunYaml{}, fmt.Errorf("decode Trail action target: %w", err)
		}
		if target.RobotSessionID != "" {
			robotSessionID = &target.RobotSessionID
		}
	}

	var output *mcp.TrailToolRobotOutputYaml
	if len(input.Output) > 0 {
		var result trail.RobotInvocationOutput
		if err := json.Unmarshal(input.Output, &result); err != nil {
			return mcp.TrailToolActionRunYaml{}, fmt.Errorf("decode Trail action output: %w", err)
		}

		var attention *mcp.TrailToolRobotAttentionYaml
		if result.Attention != nil {
			attention = &mcp.TrailToolRobotAttentionYaml{
				Reason:  result.Attention.Reason,
				Message: result.Attention.Message,
			}
		}
		output = &mcp.TrailToolRobotOutputYaml{
			Status:    mcp.TrailToolRobotOutputYamlStatus(result.Status.String()),
			Summary:   result.Summary,
			Attention: attention,
		}
	}

	return mcp.TrailToolActionRunYaml{
		Id:             input.ID.String(),
		ActionId:       input.ActionID.String(),
		Action:         actionConfig,
		Status:         mcp.TrailToolActionRunYamlStatus(input.Status.String()),
		RobotSessionId: robotSessionID,
		Output:         output,
		Error:          input.ErrorText,
		StartedAt:      input.StartedAt,
		FinishedAt:     input.FinishedAt,
		CreatedAt:      input.CreatedAt,
		UpdatedAt:      input.UpdatedAt,
	}, nil
}

func trailToolRunMessage(input *trail.Run) string {
	return fmt.Sprintf("Trail run %s is %s.", input.ID, input.Status)
}

func trailToolRunNextAction(input *trail.Run) string {
	switch input.Status {
	case trail.RunStatusQueued, trail.RunStatusRunning:
		return "Check this run again later; cancel only an individual action run that should stop."
	case trail.RunStatusAttentionRequired:
		return "Inspect blocked and failed action outputs, then resolve the reported attention or start a new manual run when appropriate."
	case trail.RunStatusCompleted:
		return "Review the action summaries; no follow-up is required unless another occurrence should run now."
	case trail.RunStatusCancelled:
		return "Inspect the independent action results before deciding whether to start a new manual run."
	case trail.RunStatusSkipped:
		return "Inspect the trigger snapshot and Trail definition to understand why the occurrence was skipped."
	default:
		return "Inspect the action results before choosing a follow-up."
	}
}

func trailToolConfirmation(verb string, input *trail.Trail) string {
	if input.NextOccurrenceAt != nil {
		return fmt.Sprintf("%s Trail %q; next occurrence is %s.", verb, input.Name, input.NextOccurrenceAt.UTC().Format(time.RFC3339))
	}
	return fmt.Sprintf("%s Trail %q with status %s and no future occurrence.", verb, input.Name, input.Status.String())
}
