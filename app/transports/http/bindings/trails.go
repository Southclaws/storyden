package bindings

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/recurrence"
	"github.com/Southclaws/storyden/app/resources/trail"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/trail/trail_manager"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Trails struct {
	manager *trail_manager.Manager
}

func NewTrails(manager *trail_manager.Manager) Trails {
	return Trails{manager: manager}
}

func (h *Trails) TrailList(ctx context.Context, _ openapi.TrailListRequestObject) (openapi.TrailListResponseObject, error) {
	items, err := h.manager.List(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	trails, err := serialiseTrailList(items)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.TrailList200JSONResponse{
		TrailListOKJSONResponse: openapi.TrailListOKJSONResponse{
			Trails: trails,
		},
	}, nil
}

func (h *Trails) TrailCreate(ctx context.Context, request openapi.TrailCreateRequestObject) (openapi.TrailCreateResponseObject, error) {
	status, err := trail.NewStatus(string(request.Body.Status))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	trigger, err := deserialiseTrailTrigger(request.Body.Trigger)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	actions, err := deserialiseTrailActionList(request.Body.Actions)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	description := opt.NewPtr(request.Body.Description).Or("")

	trail, err := h.manager.Create(
		ctx,
		xid.ID(accountID),
		request.Body.Name,
		description,
		status,
		trigger,
		actions,
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result, err := serialiseTrail(trail)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.TrailCreate200JSONResponse{
		TrailCreateOKJSONResponse: openapi.TrailCreateOKJSONResponse(result),
	}, nil
}

func (h *Trails) TrailSchedulePreview(ctx context.Context, request openapi.TrailSchedulePreviewRequestObject) (openapi.TrailSchedulePreviewResponseObject, error) {
	schedule, err := deserialiseRecurrenceSchedule(request.Body.Schedule)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}
	schedule, err = recurrence.Compile(*schedule)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	after := time.Now().UTC()
	if request.Body.After != nil {
		after = request.Body.After.UTC()
	}

	occurrences := h.manager.Preview(schedule, after)

	return openapi.TrailSchedulePreview200JSONResponse{
		TrailSchedulePreviewOKJSONResponse: openapi.TrailSchedulePreviewOKJSONResponse{
			Occurrences: occurrences,
		},
	}, nil
}

func (h *Trails) TrailGet(ctx context.Context, request openapi.TrailGetRequestObject) (openapi.TrailGetResponseObject, error) {
	trail, err := h.manager.Get(ctx, trail.ID(deserialiseID(request.TrailId)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result, err := serialiseTrail(trail)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.TrailGet200JSONResponse{
		TrailGetOKJSONResponse: openapi.TrailGetOKJSONResponse(result),
	}, nil
}

func (h *Trails) TrailUpdate(ctx context.Context, request openapi.TrailUpdateRequestObject) (openapi.TrailUpdateResponseObject, error) {
	status, err := trail.NewStatus(string(request.Body.Status))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	trigger, err := deserialiseTrailTrigger(request.Body.Trigger)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	actions, err := deserialiseTrailActionList(request.Body.Actions)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	description := opt.NewPtr(request.Body.Description).Or("")

	trail, err := h.manager.Update(
		ctx,
		trail.ID(deserialiseID(request.TrailId)),
		request.Body.Name,
		description,
		status,
		trigger,
		actions,
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result, err := serialiseTrail(trail)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.TrailUpdate200JSONResponse{
		TrailUpdateOKJSONResponse: openapi.TrailUpdateOKJSONResponse(result),
	}, nil
}

func (h *Trails) TrailRunList(ctx context.Context, request openapi.TrailRunListRequestObject) (openapi.TrailRunListResponseObject, error) {
	runs, err := h.manager.ListRuns(ctx, trail.ID(deserialiseID(request.TrailId)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result, err := serialiseTrailRunList(runs)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.TrailRunList200JSONResponse{
		TrailRunListOKJSONResponse: openapi.TrailRunListOKJSONResponse{
			Runs: result,
		},
	}, nil
}

func (h *Trails) TrailRunNow(ctx context.Context, request openapi.TrailRunNowRequestObject) (openapi.TrailRunNowResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	run, err := h.manager.RunNow(ctx, trail.ID(deserialiseID(request.TrailId)), xid.ID(accountID))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result, err := serialiseTrailRun(run)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.TrailRunNow200JSONResponse{
		TrailRunNowOKJSONResponse: openapi.TrailRunNowOKJSONResponse(result),
	}, nil
}

func (h *Trails) TrailRunGet(ctx context.Context, request openapi.TrailRunGetRequestObject) (openapi.TrailRunGetResponseObject, error) {
	run, err := h.manager.GetRun(
		ctx,
		trail.ID(deserialiseID(request.TrailId)),
		trail.RunID(deserialiseID(request.RunId)),
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result, err := serialiseTrailRun(run)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.TrailRunGet200JSONResponse{
		TrailRunGetOKJSONResponse: openapi.TrailRunGetOKJSONResponse(result),
	}, nil
}

func (h *Trails) TrailActionRunCancel(ctx context.Context, request openapi.TrailActionRunCancelRequestObject) (openapi.TrailActionRunCancelResponseObject, error) {
	if err := h.manager.Cancel(
		ctx,
		trail.ID(deserialiseID(request.TrailId)),
		trail.RunID(deserialiseID(request.RunId)),
		trail.ActionRunID(deserialiseID(request.ActionRunId)),
	); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.TrailActionRunCancel200Response{}, nil
}

func deserialiseTrailTrigger(input openapi.TrailTrigger) (trail.Trigger, error) {
	value, err := input.ValueByDiscriminator()
	if err != nil {
		return trail.Trigger{}, err
	}

	switch trigger := value.(type) {
	case openapi.TrailTriggerSchedule:
		schedule, err := deserialiseRecurrenceSchedule(trigger.Schedule)
		if err != nil {
			return trail.Trigger{}, err
		}

		return trail.NewScheduleTrigger(*schedule)

	case openapi.TrailTriggerEvent:
		return trail.NewEventTrigger(trigger.Events)

	default:
		return trail.Trigger{}, fault.Newf("invalid Trail trigger type %T", value)
	}
}

func deserialiseRecurrenceSchedule(input openapi.RecurrenceSchedule) (*recurrence.Schedule, error) {
	frequency, err := recurrence.NewFrequency(string(input.Rule.Frequency))
	if err != nil {
		return nil, err
	}

	weekdays := []recurrence.Weekday{}
	if input.Rule.ByWeekday != nil {
		weekdays, err = dt.MapErr(*input.Rule.ByWeekday, func(value openapi.RecurrenceWeekday) (recurrence.Weekday, error) {
			return recurrence.NewWeekday(string(value))
		})
		if err != nil {
			return nil, err
		}
	}

	byMonth := []int{}
	if input.Rule.ByMonth != nil {
		byMonth = *input.Rule.ByMonth
	}

	byMonthDay := []int{}
	if input.Rule.ByMonthDay != nil {
		byMonthDay = *input.Rule.ByMonthDay
	}

	schedule := recurrence.Schedule{
		Start:    input.Start,
		Timezone: input.Timezone,
		Rule: recurrence.Rule{
			Frequency:  frequency,
			Interval:   input.Rule.Interval,
			ByWeekday:  weekdays,
			ByMonth:    byMonth,
			ByMonthDay: byMonthDay,
			Count:      input.Rule.Count,
		},
	}

	return &schedule, nil
}

func deserialiseTrailActionList(input []openapi.TrailAction) ([]trail.ActionSpec, error) {
	return dt.MapErr(input, deserialiseTrailAction)
}

func deserialiseTrailAction(input openapi.TrailAction) (trail.ActionSpec, error) {
	value, err := input.ValueByDiscriminator()
	if err != nil {
		return trail.ActionSpec{}, err
	}

	switch action := value.(type) {
	case openapi.TrailActionRobotInvocation:
		kind, err := trail.NewActionKind(string(action.Type))
		if err != nil {
			return trail.ActionSpec{}, err
		}

		config, err := json.Marshal(trail.RobotRunConfig{
			Type:        kind,
			RobotRef:    action.RobotRef,
			Instruction: action.Instruction,
		})
		if err != nil {
			return trail.ActionSpec{}, err
		}

		return trail.ActionSpec{Kind: kind, Config: config}, nil

	default:
		return trail.ActionSpec{}, fault.Newf("invalid Trail action type %T", value)
	}
}

func serialiseTrailList(input []*trail.Trail) (openapi.TrailList, error) {
	return dt.MapErr(input, serialiseTrail)
}

func serialiseTrail(input *trail.Trail) (openapi.Trail, error) {
	trigger, err := serialiseTrailTrigger(input.Trigger)
	if err != nil {
		return openapi.Trail{}, err
	}

	actions, err := serialiseTrailActionBindingList(input.Actions)
	if err != nil {
		return openapi.Trail{}, err
	}

	return openapi.Trail{
		Id:               input.ID.String(),
		CreatedAt:        input.CreatedAt,
		UpdatedAt:        input.UpdatedAt,
		CreatedBy:        serialiseProfileReferenceFromAccount(input.Creator),
		Name:             input.Name,
		Description:      input.Description,
		Status:           openapi.TrailStatus(input.Status.String()),
		Trigger:          trigger,
		Actions:          actions,
		NextOccurrenceAt: input.NextOccurrenceAt,
		LastOccurrenceAt: input.LastOccurrenceAt,
	}, nil
}

func serialiseTrailTrigger(input trail.Trigger) (openapi.TrailTrigger, error) {
	trigger := openapi.TrailTrigger{}
	switch input.Type() {
	case trail.TriggerTypeSchedule:
		schedule := input.Schedule()
		if schedule == nil {
			return openapi.TrailTrigger{}, trail.ErrInvalidScheduleTrigger
		}

		if err := trigger.FromTrailTriggerSchedule(openapi.TrailTriggerSchedule{
			Type:     openapi.TrailTriggerType(input.Type().String()),
			Schedule: serialiseRecurrenceSchedule(schedule),
		}); err != nil {
			return openapi.TrailTrigger{}, err
		}

	case trail.TriggerTypeEvent:
		event := input.Event()
		if event == nil {
			return openapi.TrailTrigger{}, trail.ErrInvalidEventTrigger
		}

		if err := trigger.FromTrailTriggerEvent(openapi.TrailTriggerEvent{
			Type:   openapi.TrailTriggerType(input.Type().String()),
			Events: event.Events,
		}); err != nil {
			return openapi.TrailTrigger{}, err
		}

	default:
		return openapi.TrailTrigger{}, trail.ErrUnsupportedTrigger
	}

	return trigger, nil
}

func serialiseRecurrenceSchedule(input *recurrence.Schedule) openapi.RecurrenceSchedule {
	weekdays := dt.Map(input.Rule.ByWeekday, func(value recurrence.Weekday) openapi.RecurrenceWeekday {
		return openapi.RecurrenceWeekday(value.String())
	})

	return openapi.RecurrenceSchedule{
		Start:    input.Start,
		Timezone: input.Timezone,
		Rule: openapi.RecurrenceRule{
			Frequency:  openapi.RecurrenceFrequency(input.Rule.Frequency.String()),
			Interval:   input.Rule.Interval,
			ByWeekday:  opt.NewSafe(openapi.RecurrenceWeekdayList(weekdays), len(weekdays) > 0).Ptr(),
			ByMonth:    opt.NewSafe(openapi.RecurrenceMonthList(input.Rule.ByMonth), len(input.Rule.ByMonth) > 0).Ptr(),
			ByMonthDay: opt.NewSafe(openapi.RecurrenceMonthDayList(input.Rule.ByMonthDay), len(input.Rule.ByMonthDay) > 0).Ptr(),
			Count:      input.Rule.Count,
		},
	}
}

func serialiseTrailActionBindingList(input []*trail.Action) (openapi.TrailActionBindingList, error) {
	return dt.MapErr(input, serialiseTrailActionBinding)
}

func serialiseTrailActionBinding(input *trail.Action) (openapi.TrailActionBinding, error) {
	action, err := serialiseTrailAction(input.Kind, input.Config)
	if err != nil {
		return openapi.TrailActionBinding{}, err
	}

	return openapi.TrailActionBinding{
		Id:        input.ID.String(),
		CreatedAt: input.CreatedAt,
		UpdatedAt: input.UpdatedAt,
		Action:    action,
	}, nil
}

func serialiseTrailAction(kind trail.ActionKind, config json.RawMessage) (openapi.TrailAction, error) {
	if kind != trail.ActionKindRobotRun {
		return openapi.TrailAction{}, fault.Newf("unsupported Trail action type %s", kind)
	}

	var robotRun trail.RobotRunConfig
	if err := json.Unmarshal(config, &robotRun); err != nil {
		return openapi.TrailAction{}, err
	}

	action := openapi.TrailAction{}
	if err := action.FromTrailActionRobotInvocation(openapi.TrailActionRobotInvocation{
		Type:        openapi.TrailActionType(kind.String()),
		RobotRef:    robotRun.RobotRef,
		Instruction: robotRun.Instruction,
	}); err != nil {
		return openapi.TrailAction{}, err
	}

	return action, nil
}

func serialiseTrailRunList(input []*trail.Run) (openapi.TrailRunList, error) {
	return dt.MapErr(input, serialiseTrailRun)
}

func serialiseTrailRun(input *trail.Run) (openapi.TrailRun, error) {
	var trigger *openapi.TrailTriggerSnapshot
	if input.Trigger != nil {
		value, err := serialiseTrailTriggerSnapshot(*input.Trigger)
		if err != nil {
			return openapi.TrailRun{}, err
		}
		trigger = &value
	}

	actions, err := serialiseTrailActionRunList(input.ActionRuns)
	if err != nil {
		return openapi.TrailRun{}, err
	}

	return openapi.TrailRun{
		Id:           input.ID.String(),
		CreatedAt:    input.CreatedAt,
		UpdatedAt:    input.UpdatedAt,
		TrailId:      input.TrailID.String(),
		Trigger:      trigger,
		Status:       openapi.TrailRunStatus(input.Status.String()),
		ScheduledFor: input.ScheduledFor,
		Actions:      actions,
		FinishedAt:   input.FinishedAt,
	}, nil
}

func serialiseTrailTriggerSnapshot(event trail.TriggerEvent) (openapi.TrailTriggerSnapshot, error) {
	trigger, err := serialiseTrailTrigger(event.Trigger)
	if err != nil {
		return openapi.TrailTriggerSnapshot{}, err
	}

	initiatedBy := opt.NewEmpty[openapi.Identifier]()
	if event.InitiatedBy != "" {
		initiatedBy = opt.New(openapi.Identifier(event.InitiatedBy))
	}

	var payload *openapi.ArbitraryData
	if len(event.Payload) > 0 {
		var value openapi.ArbitraryData
		if err := json.Unmarshal(event.Payload, &value); err != nil {
			return openapi.TrailTriggerSnapshot{}, err
		}
		payload = &value
	}

	return openapi.TrailTriggerSnapshot{
		TrailId:      event.TrailID,
		TrailRunId:   event.TrailRunID,
		Kind:         openapi.TrailRunKind(event.Kind.String()),
		Trigger:      trigger,
		ScheduledFor: event.ScheduledFor,
		ObservedAt:   event.ObservedAt,
		InitiatedBy:  initiatedBy.Ptr(),
		Payload:      payload,
	}, nil
}

func serialiseTrailActionRunList(input []*trail.ActionRun) (openapi.TrailActionRunList, error) {
	return dt.MapErr(input, serialiseTrailActionRun)
}

func serialiseTrailActionRun(input *trail.ActionRun) (openapi.TrailActionRun, error) {
	if input.Action == nil {
		return openapi.TrailActionRun{}, fault.New("Trail action run is missing its action binding")
	}

	action, err := serialiseTrailActionRunBinding(input)
	if err != nil {
		return openapi.TrailActionRun{}, err
	}

	output, err := serialiseTrailActionRunOutput(input.Output)
	if err != nil {
		return openapi.TrailActionRun{}, err
	}

	target, err := serialiseTrailActionRunTarget(input)
	if err != nil {
		return openapi.TrailActionRun{}, err
	}

	return openapi.TrailActionRun{
		Id:         input.ID.String(),
		CreatedAt:  input.CreatedAt,
		UpdatedAt:  input.UpdatedAt,
		Action:     action,
		Status:     openapi.TrailActionRunStatus(input.Status.String()),
		Output:     output,
		Target:     target,
		Error:      input.ErrorText,
		StartedAt:  input.StartedAt,
		FinishedAt: input.FinishedAt,
	}, nil
}

func serialiseTrailActionRunBinding(input *trail.ActionRun) (openapi.TrailActionBinding, error) {
	action, err := serialiseTrailAction(input.Kind, input.Config)
	if err != nil {
		return openapi.TrailActionBinding{}, err
	}

	return openapi.TrailActionBinding{
		Id:        input.Action.ID.String(),
		CreatedAt: input.Action.CreatedAt,
		UpdatedAt: input.Action.UpdatedAt,
		Action:    action,
	}, nil
}

func serialiseTrailActionRunOutput(input json.RawMessage) (*openapi.ArbitraryData, error) {
	if len(input) == 0 {
		return nil, nil
	}

	var output openapi.ArbitraryData
	if err := json.Unmarshal(input, &output); err != nil {
		return nil, err
	}

	return &output, nil
}

func serialiseTrailActionRunTarget(input *trail.ActionRun) (*openapi.TrailActionRunTarget, error) {
	if len(input.Target) == 0 {
		return nil, nil
	}

	switch input.Kind {
	case trail.ActionKindRobotRun:
		var invocation trail.RobotInvocation
		if err := json.Unmarshal(input.Target, &invocation); err != nil {
			return nil, err
		}

		output, err := serialiseRobotInvocationOutput(input.Output)
		if err != nil {
			return nil, err
		}

		target := openapi.TrailActionRunTarget{}
		if err := target.FromTrailActionRunTargetRobotInvocation(openapi.TrailActionRunTargetRobotInvocation{
			Type:           openapi.TrailActionType(input.Kind.String()),
			RobotSessionId: invocation.RobotSessionID,
			Output:         output,
		}); err != nil {
			return nil, err
		}

		return &target, nil

	default:
		return nil, fault.Newf("unsupported Trail action run target type %s", input.Kind)
	}
}

func serialiseRobotInvocationOutput(input json.RawMessage) (*openapi.RobotInvocationOutput, error) {
	if len(input) == 0 {
		return nil, nil
	}

	var output trail.RobotInvocationOutput
	if err := json.Unmarshal(input, &output); err != nil {
		return nil, err
	}

	var attention *openapi.RobotInvocationAttention
	if output.Attention != nil {
		attention = &openapi.RobotInvocationAttention{
			Reason:  output.Attention.Reason,
			Message: output.Attention.Message,
		}
	}

	return &openapi.RobotInvocationOutput{
		Status:    openapi.RobotInvocationOutputStatus(output.Status.String()),
		Summary:   output.Summary,
		Attention: attention,
	}, nil
}
