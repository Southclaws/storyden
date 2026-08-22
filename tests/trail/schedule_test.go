package trail_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestTrailSchedulePreview(t *testing.T) {
	tests.ParallelIfNotSharedPostgres(t)

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		accounts *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			t.Run("weekly_interval", func(t *testing.T) {
				a := assert.New(t)

				after := time.Date(2026, time.August, 18, 9, 0, 0, 0, time.UTC)
				weekdays := openapi.RecurrenceWeekdayList{openapi.RecurrenceWeekday("tuesday")}
				schedule := trailSchedule("2026-08-18T10:00:00", "Europe/London", openapi.Weekly)
				schedule.Rule.Interval = 2
				schedule.Rule.ByWeekday = &weekdays

				response := tests.AssertRequest(cl.TrailSchedulePreviewWithResponse(root, openapi.TrailSchedulePreviewProps{
					After:    &after,
					Schedule: schedule,
				}, adminSession))(t, http.StatusOK)

				a.Equal([]time.Time{
					time.Date(2026, time.September, 1, 9, 0, 0, 0, time.UTC),
					time.Date(2026, time.September, 15, 9, 0, 0, 0, time.UTC),
					time.Date(2026, time.September, 29, 9, 0, 0, 0, time.UTC),
					time.Date(2026, time.October, 13, 9, 0, 0, 0, time.UTC),
					time.Date(2026, time.October, 27, 10, 0, 0, 0, time.UTC),
				}, response.JSON200.Occurrences)
			})

			t.Run("monthly_last_day", func(t *testing.T) {
				a := assert.New(t)

				after := time.Date(2026, time.January, 30, 13, 0, 0, 0, time.UTC)
				monthDays := openapi.RecurrenceMonthDayList{-1}
				schedule := trailSchedule("2026-01-31T12:15:00", "UTC", openapi.Monthly)
				schedule.Rule.ByMonthDay = &monthDays

				response := tests.AssertRequest(cl.TrailSchedulePreviewWithResponse(root, openapi.TrailSchedulePreviewProps{
					After:    &after,
					Schedule: schedule,
				}, adminSession))(t, http.StatusOK)

				a.Equal([]time.Time{
					time.Date(2026, time.January, 31, 12, 15, 0, 0, time.UTC),
					time.Date(2026, time.February, 28, 12, 15, 0, 0, time.UTC),
					time.Date(2026, time.March, 31, 12, 15, 0, 0, time.UTC),
					time.Date(2026, time.April, 30, 12, 15, 0, 0, time.UTC),
					time.Date(2026, time.May, 31, 12, 15, 0, 0, time.UTC),
				}, response.JSON200.Occurrences)
			})

			t.Run("nonexistent_local_time_is_skipped", func(t *testing.T) {
				a := assert.New(t)

				after := time.Date(2026, time.March, 22, 2, 0, 0, 0, time.UTC)
				weekdays := openapi.RecurrenceWeekdayList{openapi.RecurrenceWeekday("sunday")}
				schedule := trailSchedule("2026-03-22T01:30:00", "Europe/London", openapi.Weekly)
				schedule.Rule.ByWeekday = &weekdays

				response := tests.AssertRequest(cl.TrailSchedulePreviewWithResponse(root, openapi.TrailSchedulePreviewProps{
					After:    &after,
					Schedule: schedule,
				}, adminSession))(t, http.StatusOK)

				a.Equal(time.Date(2026, time.April, 5, 0, 30, 0, 0, time.UTC), response.JSON200.Occurrences[0])
				a.Equal(time.Date(2026, time.April, 12, 0, 30, 0, 0, time.UTC), response.JSON200.Occurrences[1])
			})

			t.Run("ambiguous_local_time_fires_once", func(t *testing.T) {
				a := assert.New(t)

				after := time.Date(2026, time.October, 24, 0, 0, 0, 0, time.UTC)
				weekdays := openapi.RecurrenceWeekdayList{openapi.RecurrenceWeekday("sunday")}
				schedule := trailSchedule("2026-10-18T01:30:00", "Europe/London", openapi.Weekly)
				schedule.Rule.ByWeekday = &weekdays

				response := tests.AssertRequest(cl.TrailSchedulePreviewWithResponse(root, openapi.TrailSchedulePreviewProps{
					After:    &after,
					Schedule: schedule,
				}, adminSession))(t, http.StatusOK)

				a.Equal(time.Date(2026, time.October, 25, 1, 30, 0, 0, time.UTC), response.JSON200.Occurrences[0])
				a.Equal(time.Date(2026, time.November, 1, 1, 30, 0, 0, time.UTC), response.JSON200.Occurrences[1])
			})

			t.Run("invalid_schedules", func(t *testing.T) {
				ts := map[string]openapi.RecurrenceSchedule{
					"timezone": trailSchedule("2090-01-10T09:00:00", "Not/AZone", openapi.Daily),
					"interval": {
						Start: "2090-01-10T09:00:00", Timezone: "UTC",
						Rule: openapi.RecurrenceRule{Frequency: openapi.Daily, Interval: 0},
					},
					"offset_start": trailSchedule("2090-01-10T09:00:00Z", "UTC", openapi.Daily),
				}

				for name, schedule := range ts {
					t.Run(name, func(t *testing.T) {
						tests.AssertRequest(cl.TrailSchedulePreviewWithResponse(root, openapi.TrailSchedulePreviewProps{
							Schedule: schedule,
						}, adminSession))(t, http.StatusBadRequest)
					})
				}
			})
		}))
	}))
}

func TestTrailScheduleLifecycle(t *testing.T) {
	tests.ParallelIfNotSharedPostgres(t)

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		accounts *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)

			adminCtx, admin := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)
			action := trailRobotAction(t, xid.New().String(), "Lifecycle fixture")

			paused := createPausedTrail(t, root, cl, adminSession, action)
			r.NotNil(paused.NextOccurrenceAt)
			originalNext := *paused.NextOccurrenceAt
			var manualRun openapi.TrailRun

			t.Run("run_now_does_not_advance_schedule", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				response := tests.AssertRequest(cl.TrailRunNowWithResponse(root, paused.Id, adminSession))(t, http.StatusOK)
				manualRun = openapi.TrailRun(*response.JSON200)

				a.Equal(openapi.TrailRunKindManual, response.JSON200.Trigger.Kind)
				a.Equal(paused.Id, response.JSON200.TrailId)
				a.Equal(response.JSON200.Id, response.JSON200.Trigger.TrailRunId)
				r.NotNil(response.JSON200.Trigger.InitiatedBy)
				a.Equal(admin.ID.String(), string(*response.JSON200.Trigger.InitiatedBy))
				r.Len(response.JSON200.Actions, 1)

				stored := tests.AssertRequest(cl.TrailGetWithResponse(root, paused.Id, adminSession))(t, http.StatusOK)
				r.NotNil(stored.JSON200.NextOccurrenceAt)
				a.Equal(originalNext, *stored.JSON200.NextOccurrenceAt)
			})

			t.Run("queued_action_can_be_cancelled", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				r.Len(manualRun.Actions, 1)
				tests.AssertRequest(cl.TrailActionRunCancelWithResponse(
					root,
					paused.Id,
					manualRun.Id,
					manualRun.Actions[0].Id,
					adminSession,
				))(t, http.StatusOK)

				stored := tests.AssertRequest(cl.TrailRunGetWithResponse(root, paused.Id, manualRun.Id, adminSession))(t, http.StatusOK)
				a.Equal(openapi.TrailRunStatusCancelled, stored.JSON200.Status)
				r.Len(stored.JSON200.Actions, 1)
				a.Equal(openapi.TrailActionRunStatusCancelled, stored.JSON200.Actions[0].Status)
				r.NotNil(stored.JSON200.FinishedAt)
			})

			t.Run("past_finite_schedule_finishes", func(t *testing.T) {
				a := assert.New(t)

				count := 1
				schedule := trailSchedule("2020-01-10T09:00:00", "UTC", openapi.Daily)
				schedule.Rule.Count = &count

				props := trailProps(t, "Past one-time Trail "+xid.New().String(), openapi.TrailMutableStatusActive, action)
				props.Trigger = trailTrigger(t, schedule)
				finished := createTrail(t, root, cl, adminSession, props)

				a.Equal(openapi.TrailStatusFinished, finished.Status)
				a.Nil(finished.NextOccurrenceAt)
			})
		}))
	}))
}
