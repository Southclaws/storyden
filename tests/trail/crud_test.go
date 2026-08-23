package trail_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/rs/xid"
	"github.com/samber/lo"
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

func TestTrailCRUD(t *testing.T) {
	tests.ParallelIfNotSharedPostgres(t)

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		accounts *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, admin := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			firstAction := trailRobotAction(t, xid.New().String(), "Post the weekly prompt")
			secondAction := trailRobotAction(t, xid.New().String(), "Check the replies")

			var created openapi.Trail

			t.Run("create", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				props := trailProps(t, "Weekly community prompt", openapi.TrailMutableStatusActive, firstAction, secondAction)
				created = createTrail(t, root, cl, adminSession, props)

				a.Equal(props.Name, created.Name)
				a.Equal(*props.Description, created.Description)
				a.Equal(openapi.TrailStatusActive, created.Status)
				a.Equal(admin.ID.String(), created.CreatedBy.Id)
				a.Equal(admin.Handle, created.CreatedBy.Handle)
				r.Len(created.Actions, 2)
				r.NotNil(created.NextOccurrenceAt)
				a.Equal("2090-01-10T09:00:00Z", created.NextOccurrenceAt.UTC().Format("2006-01-02T15:04:05Z"))

				first := requireRobotAction(t, created.Actions[0].Action)
				second := requireRobotAction(t, created.Actions[1].Action)
				a.Equal("Post the weekly prompt", first.Instruction)
				a.Equal("Check the replies", second.Instruction)
				a.NotEqual(created.Actions[0].Id, created.Actions[1].Id)

				trigger := requireScheduleTrigger(t, created.Trigger)
				a.Equal(openapi.TrailTriggerTypeSchedule, trigger.Type)
				a.Equal(openapi.Daily, trigger.Schedule.Rule.Frequency)
				a.Equal("UTC", trigger.Schedule.Timezone)
			})

			t.Run("create_event_trigger", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				props := trailProps(t, "Published thread responder", openapi.TrailMutableStatusPaused, firstAction)
				props.Trigger = trailEventTrigger(t, "EventThreadPublished")
				created := createTrail(t, root, cl, adminSession, props)

				a.Equal(openapi.TrailStatusPaused, created.Status)
				a.Nil(created.NextOccurrenceAt)
				trigger := requireEventTrigger(t, created.Trigger)
				a.Equal(openapi.TrailTriggerTypeEvent, trigger.Type)
				a.Equal("EventThreadPublished", trigger.Event)
				r.Len(created.Actions, 1)
			})

			t.Run("get", func(t *testing.T) {
				a := assert.New(t)

				response := tests.AssertRequest(cl.TrailGetWithResponse(root, created.Id, adminSession))(t, http.StatusOK)

				a.Equal(created.Id, response.JSON200.Id)
				a.Equal(created.Name, response.JSON200.Name)
				a.Equal(created.CreatedBy, response.JSON200.CreatedBy)
				a.Len(response.JSON200.Actions, 2)
			})

			t.Run("list", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				response := tests.AssertRequest(cl.TrailListWithResponse(root, adminSession))(t, http.StatusOK)
				found, ok := lo.Find(response.JSON200.Trails, func(item openapi.Trail) bool {
					return item.Id == created.Id
				})

				r.True(ok)
				a.Equal(created.Name, found.Name)
				a.Equal(created.Status, found.Status)
				a.Equal(created.NextOccurrenceAt, found.NextOccurrenceAt)
			})

			t.Run("update_replaces_mutable_definition", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				description := "Updated Trail description"
				weekly := trailSchedule("2090-01-10T11:30:00", "Europe/London", openapi.Weekly)
				weekdays := openapi.RecurrenceWeekdayList{openapi.RecurrenceWeekday("tuesday"), openapi.RecurrenceWeekday("thursday")}
				weekly.Rule.Interval = 2
				weekly.Rule.ByWeekday = &weekdays
				newAction := trailRobotAction(t, xid.New().String(), "Review the fortnight")

				response := tests.AssertRequest(cl.TrailUpdateWithResponse(root, created.Id, openapi.TrailMutableProps{
					Name:        "Fortnightly community review",
					Description: &description,
					Status:      openapi.TrailMutableStatusPaused,
					Trigger:     trailTrigger(t, weekly),
					Actions:     []openapi.TrailAction{newAction},
				}, adminSession))(t, http.StatusOK)

				updated := openapi.Trail(*response.JSON200)
				a.Equal("Fortnightly community review", updated.Name)
				a.Equal(description, updated.Description)
				a.Equal(openapi.TrailStatusPaused, updated.Status)
				a.Equal(created.CreatedBy, updated.CreatedBy)
				r.Len(updated.Actions, 1)
				a.NotEqual(created.Actions[0].Id, updated.Actions[0].Id)
				a.Equal("Review the fortnight", requireRobotAction(t, updated.Actions[0].Action).Instruction)

				stored := tests.AssertRequest(cl.TrailGetWithResponse(root, created.Id, adminSession))(t, http.StatusOK)
				a.Equal(updated.Name, stored.JSON200.Name)
				a.Equal(updated.Actions, stored.JSON200.Actions)

				created = updated
			})

			t.Run("resume", func(t *testing.T) {
				a := assert.New(t)

				props := openapi.TrailMutableProps{
					Name:        created.Name,
					Description: &created.Description,
					Status:      openapi.TrailMutableStatusActive,
					Trigger:     created.Trigger,
					Actions:     []openapi.TrailAction{created.Actions[0].Action},
				}

				response := tests.AssertRequest(cl.TrailUpdateWithResponse(root, created.Id, props, adminSession))(t, http.StatusOK)

				a.Equal(openapi.TrailStatusActive, response.JSON200.Status)
				a.NotNil(response.JSON200.NextOccurrenceAt)

				created = openapi.Trail(*response.JSON200)
			})

			t.Run("change_to_event_trigger_clears_next_occurrence", func(t *testing.T) {
				a := assert.New(t)

				response := tests.AssertRequest(cl.TrailUpdateWithResponse(root, created.Id, openapi.TrailMutableProps{
					Name:        created.Name,
					Description: &created.Description,
					Status:      openapi.TrailMutableStatusActive,
					Trigger:     trailEventTrigger(t, "EventThreadPublished"),
					Actions:     []openapi.TrailAction{created.Actions[0].Action},
				}, adminSession))(t, http.StatusOK)

				a.Equal(openapi.TrailStatusActive, response.JSON200.Status)
				a.Nil(response.JSON200.NextOccurrenceAt)
				a.Equal("EventThreadPublished", requireEventTrigger(t, response.JSON200.Trigger).Event)

				created = openapi.Trail(*response.JSON200)
			})

			t.Run("archive_retains_definition", func(t *testing.T) {
				a := assert.New(t)

				response := tests.AssertRequest(cl.TrailUpdateWithResponse(root, created.Id, openapi.TrailMutableProps{
					Name:        created.Name,
					Description: &created.Description,
					Status:      openapi.TrailMutableStatusArchived,
					Trigger:     created.Trigger,
					Actions:     []openapi.TrailAction{created.Actions[0].Action},
				}, adminSession))(t, http.StatusOK)

				a.Equal(openapi.TrailStatusArchived, response.JSON200.Status)

				get := tests.AssertRequest(cl.TrailGetWithResponse(root, created.Id, adminSession))(t, http.StatusOK)
				a.Equal(openapi.TrailStatusArchived, get.JSON200.Status)

				list := tests.AssertRequest(cl.TrailListWithResponse(root, adminSession))(t, http.StatusOK)
				_, found := lo.Find(list.JSON200.Trails, func(item openapi.Trail) bool {
					return item.Id == created.Id
				})
				a.True(found)

				tests.AssertRequest(cl.TrailRunNowWithResponse(root, created.Id, adminSession))(t, http.StatusNotFound)
				tests.AssertRequest(cl.TrailUpdateWithResponse(root, created.Id, openapi.TrailMutableProps{
					Name:    created.Name,
					Status:  openapi.TrailMutableStatusActive,
					Trigger: created.Trigger,
					Actions: []openapi.TrailAction{created.Actions[0].Action},
				}, adminSession))(t, http.StatusBadRequest)
			})

			t.Run("missing_trail", func(t *testing.T) {
				missing := xid.New().String()

				tests.AssertRequest(cl.TrailGetWithResponse(root, missing, adminSession))(t, http.StatusNotFound)
				tests.AssertRequest(cl.TrailRunListWithResponse(root, missing, adminSession))(t, http.StatusNotFound)
				tests.AssertRequest(cl.TrailRunNowWithResponse(root, missing, adminSession))(t, http.StatusNotFound)
			})
		}))
	}))
}

func TestTrailValidation(t *testing.T) {
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
			validAction := trailRobotAction(t, xid.New().String(), "Valid instruction")

			tests := map[string]openapi.TrailInitialProps{
				"empty_name": trailProps(t, " ", openapi.TrailMutableStatusPaused, validAction),
				"long_name":  trailProps(t, strings.Repeat("a", 121), openapi.TrailMutableStatusPaused, validAction),
				"no_actions": trailProps(t, "No actions", openapi.TrailMutableStatusPaused),
				"empty_robot": trailProps(
					t,
					"Empty Robot",
					openapi.TrailMutableStatusPaused,
					trailRobotAction(t, " ", "Valid instruction"),
				),
				"empty_instruction": trailProps(
					t,
					"Empty instruction",
					openapi.TrailMutableStatusPaused,
					trailRobotAction(t, xid.New().String(), " "),
				),
			}

			invalidSchedule := trailProps(t, "Invalid timezone", openapi.TrailMutableStatusPaused, validAction)
			invalidSchedule.Trigger = trailTrigger(t, trailSchedule("2090-01-10T09:00:00", "Not/AZone", openapi.Daily))
			tests["invalid_schedule"] = invalidSchedule

			invalidEvent := trailProps(t, "Invalid event", openapi.TrailMutableStatusPaused, validAction)
			invalidEvent.Trigger = trailEventTrigger(t, "EventNotReal")
			tests["invalid_event"] = invalidEvent

			for name, props := range tests {
				t.Run(name, func(t *testing.T) {
					response, err := cl.TrailCreateWithResponse(root, props, adminSession)
					r := require.New(t)

					r.NoError(err)
					r.Equal(http.StatusBadRequest, response.StatusCode())
				})
			}
		}))
	}))
}

func TestTrailPermissions(t *testing.T) {
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

			memberCtx, _ := e2e.WithAccount(root, accounts, seed.Account_004_Loki)
			memberSession := sh.WithSession(memberCtx)

			action := trailRobotAction(t, xid.New().String(), "Permission fixture")
			created := createPausedTrail(t, root, cl, adminSession, action)
			queued := tests.AssertRequest(cl.TrailRunNowWithResponse(root, created.Id, adminSession))(t, http.StatusOK)
			require.NotNil(t, queued.JSON200)
			require.Len(t, queued.JSON200.Actions, 1)

			tests := map[string]func(*testing.T, openapi.RequestEditorFn) int{
				"list": func(t *testing.T, session openapi.RequestEditorFn) int {
					response, err := cl.TrailListWithResponse(root, session)
					require.NoError(t, err)
					return response.StatusCode()
				},
				"create": func(t *testing.T, session openapi.RequestEditorFn) int {
					response, err := cl.TrailCreateWithResponse(root, trailProps(t, "Forbidden Trail", openapi.TrailMutableStatusPaused, action), session)
					require.NoError(t, err)
					return response.StatusCode()
				},
				"preview": func(t *testing.T, session openapi.RequestEditorFn) int {
					response, err := cl.TrailSchedulePreviewWithResponse(root, openapi.TrailSchedulePreviewProps{
						Schedule: trailSchedule("2090-01-10T09:00:00", "UTC", openapi.Daily),
					}, session)
					require.NoError(t, err)
					return response.StatusCode()
				},
				"get": func(t *testing.T, session openapi.RequestEditorFn) int {
					response, err := cl.TrailGetWithResponse(root, created.Id, session)
					require.NoError(t, err)
					return response.StatusCode()
				},
				"update": func(t *testing.T, session openapi.RequestEditorFn) int {
					response, err := cl.TrailUpdateWithResponse(root, created.Id, openapi.TrailMutableProps{
						Name: created.Name, Status: openapi.TrailMutableStatusPaused, Trigger: created.Trigger, Actions: []openapi.TrailAction{action},
					}, session)
					require.NoError(t, err)
					return response.StatusCode()
				},
				"run_list": func(t *testing.T, session openapi.RequestEditorFn) int {
					response, err := cl.TrailRunListWithResponse(root, created.Id, session)
					require.NoError(t, err)
					return response.StatusCode()
				},
				"run_now": func(t *testing.T, session openapi.RequestEditorFn) int {
					response, err := cl.TrailRunNowWithResponse(root, created.Id, session)
					require.NoError(t, err)
					return response.StatusCode()
				},
				"run_get": func(t *testing.T, session openapi.RequestEditorFn) int {
					response, err := cl.TrailRunGetWithResponse(root, created.Id, queued.JSON200.Id, session)
					require.NoError(t, err)
					return response.StatusCode()
				},
				"action_cancel": func(t *testing.T, session openapi.RequestEditorFn) int {
					response, err := cl.TrailActionRunCancelWithResponse(
						root,
						created.Id,
						queued.JSON200.Id,
						queued.JSON200.Actions[0].Id,
						session,
					)
					require.NoError(t, err)
					return response.StatusCode()
				},
			}

			for name, request := range tests {
				t.Run(name, func(t *testing.T) {
					a := assert.New(t)

					a.Equal(http.StatusForbidden, request(t, memberSession))
					a.Equal(http.StatusOK, request(t, adminSession))
				})
			}

			t.Run("unauthenticated", func(t *testing.T) {
				r := require.New(t)

				response, err := cl.TrailListWithResponse(root)
				r.NoError(err)
				r.Equal(http.StatusUnauthorized, response.StatusCode())
			})
		}))
	}))
}
