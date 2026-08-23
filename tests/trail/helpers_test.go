package trail_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/tests"
)

func trailSchedule(start, timezone string, frequency openapi.RecurrenceFrequency) openapi.RecurrenceSchedule {
	return openapi.RecurrenceSchedule{
		Start:    start,
		Timezone: timezone,
		Rule: openapi.RecurrenceRule{
			Frequency: frequency,
			Interval:  1,
		},
	}
}

func trailTrigger(t *testing.T, schedule openapi.RecurrenceSchedule) openapi.TrailTrigger {
	t.Helper()

	trigger := openapi.TrailTrigger{}
	require.NoError(t, trigger.FromTrailTriggerSchedule(openapi.TrailTriggerSchedule{
		Type:     openapi.TrailTriggerTypeSchedule,
		Schedule: schedule,
	}))

	return trigger
}

func trailEventTrigger(t *testing.T, event string) openapi.TrailTrigger {
	t.Helper()

	trigger := openapi.TrailTrigger{}
	require.NoError(t, trigger.FromTrailTriggerEvent(openapi.TrailTriggerEvent{
		Type:  openapi.TrailTriggerTypeEvent,
		Event: event,
	}))

	return trigger
}

func trailRobotAction(t *testing.T, robotRef, instruction string) openapi.TrailAction {
	t.Helper()

	action := openapi.TrailAction{}
	require.NoError(t, action.FromTrailActionRobotInvocation(openapi.TrailActionRobotInvocation{
		Type:        openapi.RobotRun,
		RobotRef:    robotRef,
		Instruction: instruction,
	}))

	return action
}

func trailProps(t *testing.T, name string, status openapi.TrailMutableStatus, actions ...openapi.TrailAction) openapi.TrailInitialProps {
	t.Helper()

	description := "Trail integration fixture"

	return openapi.TrailInitialProps{
		Name:        name,
		Description: &description,
		Status:      status,
		Trigger:     trailTrigger(t, trailSchedule("2090-01-10T09:00:00", "UTC", openapi.Daily)),
		Actions:     actions,
	}
}

func createTrail(
	t *testing.T,
	ctx context.Context,
	cl *openapi.ClientWithResponses,
	session openapi.RequestEditorFn,
	props openapi.TrailInitialProps,
) openapi.Trail {
	t.Helper()

	response := tests.AssertRequest(cl.TrailCreateWithResponse(ctx, props, session))(t, http.StatusOK)
	require.NotNil(t, response.JSON200)

	return openapi.Trail(*response.JSON200)
}

func createPausedTrail(
	t *testing.T,
	ctx context.Context,
	cl *openapi.ClientWithResponses,
	session openapi.RequestEditorFn,
	actions ...openapi.TrailAction,
) openapi.Trail {
	t.Helper()

	return createTrail(t, ctx, cl, session, trailProps(
		t,
		"Trail "+xid.New().String(),
		openapi.TrailMutableStatusPaused,
		actions...,
	))
}

func requireRobotAction(t *testing.T, action openapi.TrailAction) openapi.TrailActionRobotInvocation {
	t.Helper()

	result, err := action.AsTrailActionRobotInvocation()
	require.NoError(t, err)

	return result
}

func requireScheduleTrigger(t *testing.T, trigger openapi.TrailTrigger) openapi.TrailTriggerSchedule {
	t.Helper()

	result, err := trigger.AsTrailTriggerSchedule()
	require.NoError(t, err)

	return result
}

func requireEventTrigger(t *testing.T, trigger openapi.TrailTrigger) openapi.TrailTriggerEvent {
	t.Helper()

	result, err := trigger.AsTrailTriggerEvent()
	require.NoError(t, err)

	return result
}
