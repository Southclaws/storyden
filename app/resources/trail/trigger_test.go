package trail

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/recurrence"
	"github.com/Southclaws/storyden/internal/ent"
	entrun "github.com/Southclaws/storyden/internal/ent/trailrun"
)

func TestEventTriggerValidation(t *testing.T) {
	t.Parallel()

	events := []string{"EventThreadPublished", "EventNodeUpdated"}
	trigger, err := NewEventTrigger(events)
	require.NoError(t, err)

	events[0] = "EventPostLiked"
	require.Equal(t, []string{"EventThreadPublished", "EventNodeUpdated"}, trigger.Event().Events)

	materialised := trigger.Event()
	materialised.Events[0] = "EventPostLiked"
	require.Equal(t, []string{"EventThreadPublished", "EventNodeUpdated"}, trigger.Event().Events)

	for _, input := range [][]string{
		nil,
		{},
		{"EventThreadPublished", "EventThreadPublished"},
		{"not-an-event"},
	} {
		_, err := NewEventTrigger(input)
		require.ErrorIs(t, err, ErrInvalidEventTrigger)
	}
}

func TestTriggerStorageRoundTrip(t *testing.T) {
	t.Parallel()

	t.Run("event", func(t *testing.T) {
		trigger, err := NewEventTrigger([]string{"EventThreadPublished", "EventNodeUpdated"})
		require.NoError(t, err)

		kind, encoded, err := encodeTrigger(trigger)
		require.NoError(t, err)
		require.Equal(t, TriggerTypeEvent, kind)
		require.JSONEq(t, `{"events":["EventThreadPublished","EventNodeUpdated"]}`, string(encoded))

		decoded, err := decodeTrigger("event", encoded)
		require.NoError(t, err)
		require.Equal(t, TriggerTypeEvent, decoded.Type())
		require.Equal(t, []string{"EventThreadPublished", "EventNodeUpdated"}, decoded.Event().Events)
	})

	t.Run("schedule", func(t *testing.T) {
		trigger, err := NewScheduleTrigger(recurrence.Schedule{
			Start:    "2090-01-10T09:00:00",
			Timezone: "UTC",
			Rule: recurrence.Rule{
				Frequency: recurrence.FrequencyDaily,
				Interval:  1,
			},
		})
		require.NoError(t, err)

		kind, encoded, err := encodeTrigger(trigger)
		require.NoError(t, err)
		require.Equal(t, TriggerTypeSchedule, kind)
		require.JSONEq(t, `{"start":"2090-01-10T09:00:00","timezone":"UTC","rule":{"frequency":"daily","interval":1}}`, string(encoded))

		decoded, err := decodeTrigger("schedule", encoded)
		require.NoError(t, err)
		require.Equal(t, TriggerTypeSchedule, decoded.Type())
		require.Equal(t, "2090-01-10T09:00:00", decoded.Schedule().Start)
		require.Equal(t, "UTC", decoded.Schedule().Timezone)
		require.Equal(t, recurrence.FrequencyDaily, decoded.Schedule().Rule.Frequency)
	})
}

func TestScheduleTriggerDefensiveCopy(t *testing.T) {
	t.Parallel()

	count := 3
	input := recurrence.Schedule{
		Start:    "2090-01-10T09:00:00",
		Timezone: "UTC",
		Rule: recurrence.Rule{
			Frequency: recurrence.FrequencyWeekly,
			Interval:  1,
			ByWeekday: []recurrence.Weekday{recurrence.WeekdayWednesday},
			Count:     &count,
		},
	}
	trigger, err := NewScheduleTrigger(input)
	require.NoError(t, err)

	input.Rule.ByWeekday[0] = recurrence.WeekdayFriday
	count = 10
	require.Equal(t, []recurrence.Weekday{recurrence.WeekdayWednesday}, trigger.Schedule().Rule.ByWeekday)
	require.Equal(t, 3, *trigger.Schedule().Rule.Count)

	materialised := trigger.Schedule()
	materialised.Start = "2100-01-01T00:00:00"
	materialised.Rule.ByWeekday[0] = recurrence.WeekdayMonday
	*materialised.Rule.Count = 20
	require.Equal(t, "2090-01-10T09:00:00", trigger.Schedule().Start)
	require.Equal(t, []recurrence.Weekday{recurrence.WeekdayWednesday}, trigger.Schedule().Rule.ByWeekday)
	require.Equal(t, 3, *trigger.Schedule().Rule.Count)
}

func TestTriggerEventStorageRoundTrip(t *testing.T) {
	t.Parallel()

	trigger, err := NewEventTrigger([]string{"EventThreadPublished", "EventNodeUpdated"})
	require.NoError(t, err)
	observedAt := time.Date(2026, time.August, 23, 10, 15, 0, 0, time.UTC)
	payload := json.RawMessage(`{"thread_id":"thread-1"}`)

	encoded, err := encodeTriggerEvent(TriggerEvent{
		TrailID:     "trail-1",
		TrailRunID:  "run-1",
		Kind:        RunKindEvent,
		EventName:   "EventThreadPublished",
		Trigger:     trigger,
		Payload:     payload,
		ObservedAt:  observedAt,
		InitiatedBy: "account-1",
	})
	require.NoError(t, err)
	require.JSONEq(t, `{
		"trail_id":"trail-1",
		"trail_run_id":"run-1",
		"kind":"event",
		"event_name":"EventThreadPublished",
		"trigger":{"type":"event","events":["EventThreadPublished","EventNodeUpdated"]},
		"payload":{"thread_id":"thread-1"},
		"observed_at":"2026-08-23T10:15:00Z",
		"initiated_by":"account-1"
	}`, string(encoded))

	decoded, err := decodeTriggerEvent(encoded)
	require.NoError(t, err)
	require.Equal(t, "trail-1", decoded.TrailID)
	require.Equal(t, "run-1", decoded.TrailRunID)
	require.Equal(t, RunKindEvent, decoded.Kind)
	require.Equal(t, "EventThreadPublished", decoded.EventName)
	require.Equal(t, []string{"EventThreadPublished", "EventNodeUpdated"}, decoded.Trigger.Event().Events)
	require.JSONEq(t, string(payload), string(decoded.Payload))
	require.Equal(t, observedAt, decoded.ObservedAt)
	require.Equal(t, "account-1", decoded.InitiatedBy)
}

func TestMaterialiseEventPayload(t *testing.T) {
	t.Parallel()

	payload, err := materialiseEventPayload(
		"EventThreadUpdated",
		json.RawMessage(`{"event":"","id":"thread-1","title":"Welcome"}`),
	)
	require.NoError(t, err)
	require.JSONEq(t, `{
		"event":"EventThreadUpdated",
		"id":"thread-1",
		"title":"Welcome"
	}`, string(payload))
}

func TestDecodeTriggerRejectsInvalidStorage(t *testing.T) {
	t.Parallel()

	_, err := decodeTrigger("event", json.RawMessage(`{"events":["EventThreadPublished"],"unexpected":true}`))
	require.Error(t, err)

	_, err = decodeTrigger("event", json.RawMessage(`{"events":["not-an-event"]}`))
	require.ErrorIs(t, err, ErrInvalidEventTrigger)
}

func TestMapRunPreservesHistoryWithInvalidTrigger(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.August, 23, 18, 0, 0, 0, time.UTC)
	row := &ent.TrailRun{
		ID:             xid.New(),
		TrailID:        xid.New(),
		Kind:           entrun.KindScheduled,
		TriggerPayload: json.RawMessage(`{"invalid":true}`),
		Status:         entrun.StatusCompleted,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	run, err := mapRun(row)
	require.NoError(t, err)
	require.Nil(t, run.Trigger)
	require.Equal(t, RunStatusCompleted, run.Status)
	require.Equal(t, now, run.CreatedAt)
}
