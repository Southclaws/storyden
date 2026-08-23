package trail_tools

import (
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/recurrence"
	"github.com/Southclaws/storyden/app/resources/trail"
	robottools "github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	"github.com/Southclaws/storyden/lib/mcp"
)

func TestTrailToolsRegisterAsDirectManageTrailsTools(t *testing.T) {
	registry := robottools.NewRegistry(slog.Default())
	New(nil, registry)

	toolNames := []string{
		"trail_create",
		"trail_list",
		"trail_get",
		"trail_update",
		"trail_schedule_preview",
		"trail_run_list",
		"trail_run_get",
		"trail_run_create",
		"trail_action_run_cancel",
	}
	for _, name := range toolNames {
		registered, err := registry.GetTool(t.Context(), name)
		require.NoError(t, err)
		assert.Equal(t, []string{"system.trails"}, registered.Definition.Toolsets)
		assert.NotNil(t, registered.Handler)

		permission, ok := registered.Definition.RequiredPermission.Get()
		require.True(t, ok)
		assert.Equal(t, rbac.PermissionManageTrails, permission)
	}

	for _, name := range []string{"trail_run_create", "trail_action_run_cancel"} {
		registered, err := registry.GetTool(t.Context(), name)
		require.NoError(t, err)
		assert.True(t, registered.Definition.RequiresConfirmation)
	}
}

func TestDeserialiseTrailToolActionsDefaultsToCurrentRobot(t *testing.T) {
	ctx := robottools.ContextWithRunContext(t.Context(), robottools.RunContext{
		RobotRef: "denbot",
	})

	actions, err := deserialiseTrailToolActions(ctx, []mcp.TrailToolActionYaml{{
		Type:        "robot_run",
		Instruction: "Continue the scheduled task.",
	}})
	require.NoError(t, err)
	require.Len(t, actions, 1)

	var config trail.RobotRunConfig
	require.NoError(t, json.Unmarshal(actions[0].Config, &config))
	assert.Equal(t, trail.ActionKindRobotRun, actions[0].Kind)
	assert.Equal(t, "denbot", config.RobotRef)
	assert.Equal(t, "Continue the scheduled task.", config.Instruction)
}

func TestDeserialiseTrailToolActionsRequiresRobotReferenceOutsideRobotRun(t *testing.T) {
	_, err := deserialiseTrailToolActions(context.Background(), []mcp.TrailToolActionYaml{{
		Type:        "robot_run",
		Instruction: "Continue the scheduled task.",
	}})
	require.Error(t, err)
	assert.ErrorContains(t, err, "robot_ref is required outside a Robot conversation")
}

func TestTrailToolOneShotScheduleRoundTrip(t *testing.T) {
	count := 1
	input := mcp.TrailToolScheduleYaml{
		Start:    "2026-08-24T09:30:00",
		Timezone: "Europe/London",
		Rule: mcp.TrailToolRecurrenceRuleYaml{
			Frequency: mcp.TrailToolRecurrenceRuleYamlFrequencyDaily,
			Interval:  1,
			Count:     &count,
		},
	}

	schedule, err := deserialiseTrailToolSchedule(input)
	require.NoError(t, err)
	assert.Equal(t, 1, *schedule.Rule.Count)

	result := serialiseTrailToolSchedule(schedule)
	assert.Equal(t, input.Start, result.Start)
	assert.Equal(t, input.Timezone, result.Timezone)
	assert.Equal(t, count, *result.Rule.Count)

	after := time.Date(2026, time.August, 23, 0, 0, 0, 0, time.UTC)
	occurrences := recurrence.Preview(schedule, after, 5)
	require.Len(t, occurrences, 1)
	assert.Equal(t, time.Date(2026, time.August, 24, 8, 30, 0, 0, time.UTC), occurrences[0])
}

func TestSerialiseTrailToolRunIncludesTriggerAndIndependentActionResult(t *testing.T) {
	trailID := trail.ID(xid.New())
	runID := trail.RunID(xid.New())
	actionID := trail.ActionID(xid.New())
	actionRunID := trail.ActionRunID(xid.New())
	sessionID := xid.New().String()
	now := time.Date(2026, time.August, 23, 12, 0, 0, 0, time.UTC)

	trigger, err := trail.NewEventTrigger([]string{"EventThreadPublished"})
	require.NoError(t, err)
	actionConfig, err := json.Marshal(trail.RobotRunConfig{
		Type:        trail.ActionKindRobotRun,
		RobotRef:    "moderator",
		Instruction: "Inspect the published thread.",
	})
	require.NoError(t, err)
	target, err := json.Marshal(trail.RobotInvocation{
		Type:           trail.ActionKindRobotRun,
		RobotSessionID: sessionID,
	})
	require.NoError(t, err)
	output, err := json.Marshal(trail.RobotInvocationOutput{
		Status:  trail.RobotInvocationOutputStatusBlocked,
		Summary: "A policy decision is required.",
		Attention: &trail.RobotInvocationAttention{
			Reason:  "approval_required",
			Message: "Choose whether to hide this thread.",
		},
	})
	require.NoError(t, err)

	run := &trail.Run{
		ID:        runID,
		TrailID:   trailID,
		Kind:      trail.RunKindScheduled,
		Status:    trail.RunStatusAttentionRequired,
		CreatedAt: now,
		UpdatedAt: now,
		Trigger: &trail.TriggerEvent{
			TrailID:    trailID.String(),
			TrailRunID: runID.String(),
			Kind:       trail.RunKindEvent,
			EventName:  "EventThreadPublished",
			Trigger:    trigger,
			Payload:    json.RawMessage(`{"thread_id":"thread-123"}`),
			ObservedAt: now,
		},
		ActionRuns: []*trail.ActionRun{{
			ID:        actionRunID,
			RunID:     runID,
			ActionID:  actionID,
			Kind:      trail.ActionKindRobotRun,
			Config:    actionConfig,
			Status:    trail.ActionRunStatusBlocked,
			Target:    target,
			Output:    output,
			CreatedAt: now,
			UpdatedAt: now,
		}},
	}

	result, err := serialiseTrailToolRun(run)
	require.NoError(t, err)
	assert.Equal(t, mcp.TrailToolRunSummaryYamlKindEvent, result.Summary.Kind)
	assert.Equal(t, 1, result.Summary.ActionStatus.Blocked)
	require.NotNil(t, result.Trigger)
	assert.Equal(t, "EventThreadPublished", *result.Trigger.EventName)
	assert.JSONEq(t, `{"thread_id":"thread-123"}`, *result.Trigger.EventPayloadJson)
	require.Len(t, result.Action, 1)
	assert.Equal(t, sessionID, *result.Action[0].RobotSessionId)
	require.NotNil(t, result.Action[0].Output)
	assert.Equal(t, mcp.TrailToolRobotOutputYamlStatusBlocked, result.Action[0].Output.Status)
	require.NotNil(t, result.Action[0].Output.Attention)
	assert.Equal(t, "approval_required", result.Action[0].Output.Attention.Reason)
}
