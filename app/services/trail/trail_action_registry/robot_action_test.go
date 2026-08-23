package trail_action_registry

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/trail"
)

func TestMaterialiseTrailTriggerContext(t *testing.T) {
	t.Parallel()

	trigger, err := trail.NewEventTrigger([]string{"EventThreadPublished", "EventNodeUpdated"})
	require.NoError(t, err)

	context := materialiseTrailTriggerContext(trail.TriggerEvent{
		TrailID:     "trail-1",
		TrailRunID:  "run-1",
		Kind:        trail.RunKindEvent,
		Trigger:     trigger,
		Payload:     json.RawMessage(`{"thread_id":"thread-1"}`),
		ObservedAt:  time.Date(2026, time.August, 23, 10, 15, 0, 0, time.UTC),
		InitiatedBy: "account-1",
	})

	encoded, err := json.Marshal(context)
	require.NoError(t, err)
	require.JSONEq(t, `{
		"trail_id":"trail-1",
		"trail_run_id":"run-1",
		"kind":"event",
		"trigger":{"type":"event","events":["EventThreadPublished","EventNodeUpdated"]},
		"payload":{"thread_id":"thread-1"},
		"observed_at":"2026-08-23T10:15:00Z",
		"initiated_by":"account-1"
	}`, string(encoded))
}
