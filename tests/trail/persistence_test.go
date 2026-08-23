package trail_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/resources/trail"
	"github.com/Southclaws/storyden/app/services/trail/trail_runtime"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestTrailSchedulerLease(t *testing.T) {
	tests.ParallelIfNotSharedPostgres(t)

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		repository *trail.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			leaseStart := time.Date(2040, time.January, 1, 0, 0, 0, 0, time.UTC)

			t.Run("acquire", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				lease, acquired, err := repository.AcquireSchedulerLease(root, leaseStart, 15*time.Second)

				r.NoError(err)
				r.True(acquired)
				r.NotNil(lease)
				a.WithinDuration(leaseStart.Add(15*time.Second), lease.ExpiresAt, 0)

				second, acquired, err := repository.AcquireSchedulerLease(root, leaseStart.Add(time.Second), 15*time.Second)
				r.NoError(err)
				a.False(acquired)
				a.Nil(second)

				renewed, err := repository.RenewSchedulerLease(root, "not-the-holder", leaseStart.Add(30*time.Second))
				r.NoError(err)
				a.False(renewed)

				renewed, err = repository.RenewSchedulerLease(root, lease.Token, leaseStart.Add(30*time.Second))
				r.NoError(err)
				a.True(renewed)

				second, acquired, err = repository.AcquireSchedulerLease(root, leaseStart.Add(20*time.Second), 15*time.Second)
				r.NoError(err)
				a.False(acquired)
				a.Nil(second)

				takeover, acquired, err := repository.AcquireSchedulerLease(root, leaseStart.Add(31*time.Second), 15*time.Second)
				r.NoError(err)
				r.True(acquired)
				r.NotNil(takeover)
				a.NotEqual(lease.Token, takeover.Token)
			})
		}))
	}))
}

func TestTrailScheduledOccurrenceMaterialisation(t *testing.T) {
	tests.ParallelIfNotSharedPostgres(t)

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		accounts *account_writer.Writer,
		repository *trail.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)

			_, creator := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
			firstOccurrence := time.Date(2040, time.January, 10, 9, 0, 0, 0, time.UTC)
			secondOccurrence := time.Date(2040, time.January, 17, 9, 0, 0, 0, time.UTC)
			config := json.RawMessage(`{"type":"robot_run","robot_ref":"robot-one","instruction":"Post the weekly prompt"}`)

			definition, err := repository.Create(
				root,
				xid.ID(creator.ID),
				"Weekly prompt",
				"",
				trail.StatusActive,
				domainScheduleTrigger(t, json.RawMessage(`{"start":"2040-01-04T09:00:00","timezone":"Europe/London","rule":{"frequency":"weekly","interval":1,"by_weekday":["wednesday"]}}`)),
				[]trail.ActionSpec{{Kind: trail.ActionKindRobotRun, Config: config}},
				opt.New(firstOccurrence),
			)
			r.NoError(err)

			t.Run("materialises_once_and_advances", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				run, created, err := repository.MaterialiseScheduled(
					root,
					definition.ID,
					firstOccurrence,
					firstOccurrence,
					&secondOccurrence,
					false,
					false,
				)

				r.NoError(err)
				r.True(created)
				r.NotNil(run)
				r.Len(run.ActionRuns, 1)
				a.Equal(trail.RunKindScheduled, run.Kind)
				a.Equal(trail.RunStatusQueued, run.Status)
				r.NotNil(run.ScheduledFor)
				a.WithinDuration(firstOccurrence, *run.ScheduledFor, 0)
				a.JSONEq(string(config), string(run.ActionRuns[0].Config))

				event := run.Trigger
				r.NotNil(event)
				a.Equal(definition.ID.String(), event.TrailID)
				a.Equal(run.ID.String(), event.TrailRunID)
				a.Equal(trail.RunKindScheduled, event.Kind)
				r.NotNil(event.ScheduledFor)
				a.Equal(firstOccurrence, *event.ScheduledFor)
				a.Empty(event.InitiatedBy)

				stored, err := repository.Get(root, definition.ID)
				r.NoError(err)
				r.NotNil(stored.NextOccurrenceAt)
				r.NotNil(stored.LastOccurrenceAt)
				a.WithinDuration(secondOccurrence, *stored.NextOccurrenceAt, 0)
				a.WithinDuration(firstOccurrence, *stored.LastOccurrenceAt, 0)

				duplicate, duplicateCreated, err := repository.MaterialiseScheduled(
					root,
					definition.ID,
					firstOccurrence,
					firstOccurrence,
					&secondOccurrence,
					false,
					false,
				)
				r.NoError(err)
				a.False(duplicateCreated)
				a.Nil(duplicate)

				cancelled, err := repository.CancelQueuedActionRun(root, run.ActionRuns[0].ID)
				r.NoError(err)
				a.True(cancelled)
			})
		}))
	}))
}

func TestTrailEventOccurrenceMaterialisation(t *testing.T) {
	tests.ParallelIfNotSharedPostgres(t)

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		accounts *account_writer.Writer,
		repository *trail.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			_, creator := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
			actionConfig := json.RawMessage(`{"type":"robot_run","robot_ref":"robot-one","instruction":"Respond to the published thread"}`)
			definition, err := repository.Create(
				root,
				xid.ID(creator.ID),
				"Published thread responder",
				"",
				trail.StatusActive,
				domainEventTrigger(t, "EventThreadPublished", "EventThreadUpdated"),
				[]trail.ActionSpec{{Kind: trail.ActionKindRobotRun, Config: actionConfig}},
				opt.NewEmpty[time.Time](),
			)
			r.NoError(err)

			observedAt := time.Date(2040, time.January, 10, 9, 0, 0, 0, time.UTC)
			sourcePayload := json.RawMessage(`{"id":"thread-one"}`)
			runID, created, err := repository.MaterialiseEvent(
				root,
				definition.ID,
				"EventThreadUpdated",
				sourcePayload,
				observedAt,
			)
			r.NoError(err)
			r.True(created)
			run, err := repository.GetRun(root, runID)
			r.NoError(err)
			r.NotNil(run)
			r.Len(run.ActionRuns, 1)
			a.Equal(trail.RunKindScheduled, run.Kind)
			a.Nil(run.ScheduledFor)
			a.JSONEq(string(actionConfig), string(run.ActionRuns[0].Config))

			event := run.Trigger
			r.NotNil(event)
			a.Equal(trail.RunKindEvent, event.Kind)
			a.Equal("EventThreadUpdated", event.EventName)
			a.Equal(observedAt, event.ObservedAt)
			a.JSONEq(`{
				"event":"EventThreadUpdated",
				"id":"thread-one"
			}`, string(event.Payload))
			a.Equal([]string{"EventThreadPublished", "EventThreadUpdated"}, event.Trigger.Event().Events)

			stored, err := repository.Get(root, definition.ID)
			r.NoError(err)
			r.NotNil(stored.LastOccurrenceAt)
			a.Equal(observedAt, *stored.LastOccurrenceAt)
			a.Nil(stored.NextOccurrenceAt)

			notCreated, ok, err := repository.MaterialiseEvent(
				root,
				definition.ID,
				"EventNodePublished",
				sourcePayload,
				observedAt.Add(time.Minute),
			)
			r.NoError(err)
			a.False(ok)
			a.Equal(trail.RunID{}, notCreated)
		}))
	}))
}

func TestTrailActionFanoutAndRecovery(t *testing.T) {
	tests.ParallelIfNotSharedPostgres(t)

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		accounts *account_writer.Writer,
		repository *trail.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)

			_, creator := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
			creatorID := xid.ID(creator.ID)
			nextOccurrence := time.Date(2040, time.January, 10, 9, 0, 0, 0, time.UTC)
			firstConfig := json.RawMessage(`{"type":"robot_run","robot_ref":"robot-a","instruction":"First consumer"}`)
			secondConfig := json.RawMessage(`{"type":"robot_run","robot_ref":"robot-b","instruction":"Second consumer"}`)

			definition, err := repository.Create(
				root,
				creatorID,
				"Independent consumers",
				"",
				trail.StatusPaused,
				domainScheduleTrigger(t, json.RawMessage(`{"start":"2040-01-10T09:00:00","timezone":"UTC","rule":{"frequency":"daily","interval":1}}`)),
				[]trail.ActionSpec{
					{Kind: trail.ActionKindRobotRun, Config: firstConfig},
					{Kind: trail.ActionKindRobotRun, Config: secondConfig},
				},
				opt.New(nextOccurrence),
			)
			r.NoError(err)

			manual, err := repository.MaterialiseManual(root, definition.ID, creatorID)
			r.NoError(err)
			r.Len(manual.ActionRuns, 2)

			t.Run("manual_trigger_and_config_snapshot", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				a.Equal(trail.RunKindManual, manual.Kind)
				a.Nil(manual.ScheduledFor)

				event := manual.Trigger
				r.NotNil(event)
				a.Equal(manual.ID.String(), event.TrailRunID)
				a.Equal(creatorID.String(), event.InitiatedBy)
				a.Equal(trail.RunKindManual, event.Kind)

				stored, err := repository.Get(root, definition.ID)
				r.NoError(err)
				r.NotNil(stored.NextOccurrenceAt)
				a.WithinDuration(nextOccurrence, *stored.NextOccurrenceAt, 0)

				updatedConfig := json.RawMessage(`{"type":"robot_run","robot_ref":"robot-new","instruction":"Replacement consumer"}`)
				_, err = repository.Update(
					root,
					definition.ID,
					definition.Name,
					definition.Description,
					trail.StatusPaused,
					definition.Trigger,
					[]trail.ActionSpec{{Kind: trail.ActionKindRobotRun, Config: updatedConfig}},
					opt.New(nextOccurrence),
				)
				r.NoError(err)

				snapshot, err := repository.GetRun(root, manual.ID)
				r.NoError(err)
				r.Len(snapshot.ActionRuns, 2)
				a.JSONEq(string(firstConfig), string(snapshot.ActionRuns[0].Config))
				a.JSONEq(string(secondConfig), string(snapshot.ActionRuns[1].Config))
			})

			t.Run("lease_recovery_and_independent_results", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				claimAt := time.Date(2040, time.February, 1, 12, 0, 0, 0, time.UTC)
				claimed, ok, err := repository.ClaimActionRunByID(root, manual.ActionRuns[0].ID, claimAt, time.Minute)
				r.NoError(err)
				r.True(ok)
				r.NotNil(claimed.LeaseToken)
				a.Equal(manual.ActionRuns[0].ID, claimed.ID)
				firstToken := *claimed.LeaseToken

				other, ok, err := repository.ClaimActionRunByID(root, manual.ActionRuns[1].ID, claimAt.Add(30*time.Second), time.Minute)
				r.NoError(err)
				r.True(ok)
				a.Equal(manual.ActionRuns[1].ID, other.ID)

				renewed, err := repository.RenewActionLease(root, claimed.ID, firstToken, claimAt.Add(2*time.Minute))
				r.NoError(err)
				a.True(renewed)

				recovered, ok, err := repository.ClaimActionRunByID(root, other.ID, claimAt.Add(91*time.Second), time.Minute)
				r.NoError(err)
				r.True(ok)
				a.Equal(other.ID, recovered.ID)

				recovered, ok, err = repository.ClaimActionRunByID(root, claimed.ID, claimAt.Add(121*time.Second), time.Minute)
				r.NoError(err)
				r.True(ok)
				a.Equal(claimed.ID, recovered.ID)
				a.NotEqual(firstToken, *recovered.LeaseToken)

				r.NoError(repository.FinishActionRun(
					root,
					recovered.ID,
					*recovered.LeaseToken,
					trail.ActionRunStatusCompleted,
					json.RawMessage(`{"status":"completed","summary":"first complete"}`),
					"",
				))

				otherRecovered, ok, err := repository.ClaimActionRunByID(root, other.ID, claimAt.Add(152*time.Second), time.Minute)
				r.NoError(err)
				r.True(ok)
				a.Equal(other.ID, otherRecovered.ID)

				r.NoError(repository.FinishActionRun(
					root,
					otherRecovered.ID,
					*otherRecovered.LeaseToken,
					trail.ActionRunStatusFailed,
					nil,
					"semantic failure",
				))

				finished, err := repository.GetRun(root, manual.ID)
				r.NoError(err)
				a.Equal(trail.RunStatusAttentionRequired, finished.Status)
				r.NotNil(finished.FinishedAt)
				r.Len(finished.ActionRuns, 2)
				a.Equal(trail.ActionRunStatusCompleted, finished.ActionRuns[0].Status)
				a.Equal(trail.ActionRunStatusFailed, finished.ActionRuns[1].Status)
				r.NotNil(finished.ActionRuns[1].ErrorText)
				a.Equal("semantic failure", *finished.ActionRuns[1].ErrorText)
			})

			t.Run("attention_notification_is_deduplicated", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				failed := manual.ActionRuns[1]
				target := manual.ID.String()

				recorded, err := repository.RecordAttentionNotification(root, failed.ID, creatorID, target)
				r.NoError(err)
				a.True(recorded)

				recorded, err = repository.RecordAttentionNotification(root, failed.ID, creatorID, target)
				r.NoError(err)
				a.False(recorded)
			})

			t.Run("all_cancelled_aggregates_cancelled", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				cancelledRun, err := repository.MaterialiseManual(root, definition.ID, creatorID)
				r.NoError(err)
				r.Len(cancelledRun.ActionRuns, 1)

				cancelled, err := repository.CancelQueuedActionRun(root, cancelledRun.ActionRuns[0].ID)
				r.NoError(err)
				a.True(cancelled)

				stored, err := repository.GetRun(root, cancelledRun.ID)
				r.NoError(err)
				a.Equal(trail.RunStatusCancelled, stored.Status)
				r.NotNil(stored.FinishedAt)
			})

			t.Run("history_is_newest_first_and_run_numbers_are_stable", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				history, err := repository.ListRuns(root, definition.ID, 10)
				r.NoError(err)
				r.Len(history, 2)
				a.Equal(manual.ID, history[1].ID)

				firstNumber, err := repository.RunNumber(root, manual.ID)
				r.NoError(err)
				a.Equal(1, firstNumber)

				secondNumber, err := repository.RunNumber(root, history[0].ID)
				r.NoError(err)
				a.Equal(2, secondNumber)
			})
		}))
	}))
}

func TestTrailMissedOneTimeOccurrence(t *testing.T) {
	tests.ParallelIfNotSharedPostgres(t)

	if tests.IsSharedPostgresDatabase() {
		t.Skip("the Trail scheduler lease is intentionally shared between backend replicas")
	}

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		accounts *account_writer.Writer,
		repository *trail.Repository,
		_ *trail_runtime.Runtime,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			_, creator := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
			missedAt := time.Date(2020, time.March, 31, 0, 30, 0, 0, time.UTC)

			definition, err := repository.Create(
				root,
				xid.ID(creator.ID),
				"One time announcement",
				"",
				trail.StatusActive,
				domainScheduleTrigger(t, json.RawMessage(`{"start":"2020-03-31T01:30:00","timezone":"Europe/London","rule":{"frequency":"daily","interval":1,"count":1}}`)),
				[]trail.ActionSpec{{
					Kind:   trail.ActionKindRobotRun,
					Config: json.RawMessage(`{"type":"robot_run","robot_ref":"robot-one","instruction":"Post once"}`),
				}},
				opt.New(missedAt),
			)
			r.NoError(err)

			var run *trail.Run
			r.Eventually(func() bool {
				history, err := repository.ListRuns(root, definition.ID, 10)
				if err != nil || len(history) != 1 {
					return false
				}

				run = history[0]

				return run.Status == trail.RunStatusSkipped
			}, 5*time.Second, 100*time.Millisecond)
			r.NotNil(run)
			a.Equal(trail.RunStatusSkipped, run.Status)
			a.Empty(run.ActionRuns)
			r.NotNil(run.FinishedAt)

			finished, err := repository.Get(root, definition.ID)
			r.NoError(err)
			a.Equal(trail.StatusFinished, finished.Status)
			a.Nil(finished.NextOccurrenceAt)
			r.NotNil(finished.LastOccurrenceAt)
			a.WithinDuration(missedAt, *finished.LastOccurrenceAt, 0)

			history, err := repository.ListRuns(root, definition.ID, 10)
			r.NoError(err)
			r.Len(history, 1)
			a.Equal(run.ID, history[0].ID)
		}))
	}))
}
