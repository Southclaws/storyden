package trail_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/trail/trail_runtime"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
	"github.com/Southclaws/storyden/tests"
	testrobot "github.com/Southclaws/storyden/tests/robot"
)

func TestTrailRobotRuntime(t *testing.T) {
	tests.ParallelIfNotSharedPostgres(t)

	integration.Test(
		t,
		&config.Config{LanguageModelProvider: "mock"},
		e2e.Setup(),
		testrobot.WithRobotSettings("mock/testdata/robot-completed.yaml"),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			accounts *account_writer.Writer,
			bus *pubsub.Bus,
			_ *trail_runtime.Runtime,
		) {
			lc.Append(fx.StartHook(func() {
				r := require.New(t)

				adminCtx, admin := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				model := openapi.RobotModelRef("mock/testdata/robot-completed.yaml")
				robot := tests.AssertRequest(cl.RobotCreateWithResponse(root, openapi.RobotCreateJSONRequestBody{
					Name:        "Trail runtime Robot " + xid.New().String(),
					Description: "Completes Trail runtime fixtures",
					Playbook:    "Complete the unattended instruction.",
					Model:       &model,
				}, adminSession))(t, http.StatusOK)
				r.NotNil(robot.JSON200)

				trailName := "Scheduled moderation check " + xid.New().String()
				definition := createTrail(t, root, cl, adminSession, trailProps(
					t,
					trailName,
					openapi.TrailMutableStatusPaused,
					trailRobotAction(t, string(robot.JSON200.Id), "Perform the scheduled check"),
				))

				t.Run("event_trigger_materialises_source_payload", func(t *testing.T) {
					r := require.New(t)
					a := assert.New(t)

					props := trailProps(
						t,
						"Published thread responder "+xid.New().String(),
						openapi.TrailMutableStatusActive,
						trailRobotAction(t, string(robot.JSON200.Id), "Respond to the published thread"),
					)
					props.Trigger = trailEventTrigger(t, "EventThreadPublished", "EventThreadUpdated")
					eventTrail := createTrail(t, root, cl, adminSession, props)
					a.Nil(eventTrail.NextOccurrenceAt)

					threadID := post.ID(xid.New())
					var runID openapi.Identifier
					r.Eventually(func() bool {
						bus.Publish(root, &rpc.EventThreadUpdated{ID: threadID})
						history, err := cl.TrailRunListWithResponse(root, eventTrail.Id, adminSession)
						if err != nil || history == nil || history.JSON200 == nil || len(history.JSON200.Runs) == 0 {
							return false
						}

						runID = history.JSON200.Runs[0].Id
						return true
					}, 5*time.Second, 100*time.Millisecond)

					run := waitForTrailRun(t, root, cl, adminSession, eventTrail.Id, runID)
					a.Equal(openapi.TrailRunKindEvent, run.Trigger.Kind)
					a.Nil(run.ScheduledFor)
					trigger := requireEventTrigger(t, run.Trigger.Trigger)
					a.Equal([]string{"EventThreadPublished", "EventThreadUpdated"}, trigger.Events)
					r.NotNil(run.Trigger.Payload)

					payload, err := json.Marshal(*run.Trigger.Payload)
					r.NoError(err)
					var source struct {
						ID string `json:"id"`
					}
					r.NoError(json.Unmarshal(payload, &source))
					a.Equal(threadID.String(), source.ID)
				})

				var first openapi.TrailRun

				t.Run("completed_invocation", func(t *testing.T) {
					r := require.New(t)
					a := assert.New(t)

					started := tests.AssertRequest(cl.TrailRunNowWithResponse(root, definition.Id, adminSession))(t, http.StatusOK)
					first = waitForTrailRun(t, root, cl, adminSession, definition.Id, started.JSON200.Id)

					a.Equal(openapi.TrailRunStatusCompleted, first.Status)
					r.Len(first.Actions, 1)
					a.Equal(openapi.TrailActionRunStatusCompleted, first.Actions[0].Status)
					a.Empty(first.Actions[0].Error)
					r.NotNil(first.Actions[0].Target)

					target, err := first.Actions[0].Target.AsTrailActionRunTargetRobotInvocation()
					r.NoError(err)
					a.Equal(openapi.RobotRun, target.Type)
					a.Equal(first.Actions[0].Id, target.RobotSessionId)
					r.NotNil(target.Output)
					a.Equal(openapi.RobotInvocationOutputStatus("completed"), target.Output.Status)
					a.Equal("Scheduled work complete.", target.Output.Summary)

					session := tests.AssertRequest(cl.RobotSessionGetWithResponse(
						root,
						target.RobotSessionId,
						&openapi.RobotSessionGetParams{},
						adminSession,
					))(t, http.StatusOK)
					r.NotNil(session.JSON200)
					a.Equal(trailName+" (Run 1)", session.JSON200.Name)
					a.Equal(admin.ID.String(), session.JSON200.CreatedBy.Id)

					_, hasUserInput := lo.Find(session.JSON200.MessageList.Messages, func(message openapi.RobotSessionMessage) bool {
						return message.Role == openapi.RobotSessionMessageRoleUser
					})
					a.True(hasUserInput)

					list := tests.AssertRequest(cl.RobotSessionsListWithResponse(
						root,
						&openapi.RobotSessionsListParams{},
						adminSession,
					))(t, http.StatusOK)
					_, listed := lo.Find(list.JSON200.Sessions, func(session openapi.RobotSessionRef) bool {
						return session.Id == target.RobotSessionId
					})
					a.False(listed)
				})

				t.Run("fresh_session_per_occurrence", func(t *testing.T) {
					r := require.New(t)
					a := assert.New(t)

					started := tests.AssertRequest(cl.TrailRunNowWithResponse(root, definition.Id, adminSession))(t, http.StatusOK)
					second := waitForTrailRun(t, root, cl, adminSession, definition.Id, started.JSON200.Id)

					r.Len(second.Actions, 1)
					r.NotNil(second.Actions[0].Target)
					firstTarget, err := first.Actions[0].Target.AsTrailActionRunTargetRobotInvocation()
					r.NoError(err)
					secondTarget, err := second.Actions[0].Target.AsTrailActionRunTargetRobotInvocation()
					r.NoError(err)
					a.NotEqual(firstTarget.RobotSessionId, secondTarget.RobotSessionId)

					session := tests.AssertRequest(cl.RobotSessionGetWithResponse(
						root,
						secondTarget.RobotSessionId,
						&openapi.RobotSessionGetParams{},
						adminSession,
					))(t, http.StatusOK)
					a.Equal(trailName+" (Run 2)", session.JSON200.Name)
				})

				t.Run("blocked_invocation_records_structured_attention", func(t *testing.T) {
					r := require.New(t)
					a := assert.New(t)

					blockedModel := openapi.RobotModelRef("mock/testdata/robot-blocked.yaml")
					blockedRobot := tests.AssertRequest(cl.RobotCreateWithResponse(root, openapi.RobotCreateJSONRequestBody{
						Name:        "Blocked Trail Robot " + xid.New().String(),
						Description: "Returns structured attention",
						Playbook:    "Report when the run needs attention.",
						Model:       &blockedModel,
					}, adminSession))(t, http.StatusOK)
					r.NotNil(blockedRobot.JSON200)

					blockedTrail := createPausedTrail(
						t,
						root,
						cl,
						adminSession,
						trailRobotAction(t, string(blockedRobot.JSON200.Id), "Check for a moderator decision"),
					)
					started := tests.AssertRequest(cl.TrailRunNowWithResponse(root, blockedTrail.Id, adminSession))(t, http.StatusOK)
					blocked := waitForTrailRun(t, root, cl, adminSession, blockedTrail.Id, started.JSON200.Id)

					a.Equal(openapi.TrailRunStatusAttentionRequired, blocked.Status)
					r.Len(blocked.Actions, 1)
					a.Equal(openapi.TrailActionRunStatusBlocked, blocked.Actions[0].Status)
					r.NotNil(blocked.Actions[0].Target)

					target, err := blocked.Actions[0].Target.AsTrailActionRunTargetRobotInvocation()
					r.NoError(err)
					r.NotNil(target.Output)
					a.Equal(openapi.RobotInvocationOutputStatus("blocked"), target.Output.Status)
					a.Equal("A moderator decision is required.", target.Output.Summary)
					r.NotNil(target.Output.Attention)
					a.Equal("missing_input", target.Output.Attention.Reason)
					a.Equal("Select a moderation outcome.", target.Output.Attention.Message)

					r.Eventually(func() bool {
						response, err := cl.NotificationListWithResponse(root, &openapi.NotificationListParams{}, adminSession)
						if err != nil || response == nil || response.JSON200 == nil {
							return false
						}

						_, ok := lo.Find(response.JSON200.Notifications, func(item openapi.Notification) bool {
							return item.Event == openapi.TrailRunAttention && item.Target != nil && *item.Target == blocked.Id
						})

						return ok
					}, 5*time.Second, 100*time.Millisecond)
				})

				t.Run("missing_structured_finish_fails", func(t *testing.T) {
					r := require.New(t)
					a := assert.New(t)

					malformedModel := openapi.RobotModelRef("mock/testdata/robot-malformed.yaml")
					malformedRobot := tests.AssertRequest(cl.RobotCreateWithResponse(root, openapi.RobotCreateJSONRequestBody{
						Name:        "Malformed Trail Robot " + xid.New().String(),
						Description: "Omits the structured finish",
						Playbook:    "Finish without the required tool.",
						Model:       &malformedModel,
					}, adminSession))(t, http.StatusOK)
					r.NotNil(malformedRobot.JSON200)

					malformedTrail := createPausedTrail(
						t,
						root,
						cl,
						adminSession,
						trailRobotAction(t, string(malformedRobot.JSON200.Id), "Return an invalid result"),
					)
					started := tests.AssertRequest(cl.TrailRunNowWithResponse(root, malformedTrail.Id, adminSession))(t, http.StatusOK)
					failed := waitForTrailRun(t, root, cl, adminSession, malformedTrail.Id, started.JSON200.Id)

					a.Equal(openapi.TrailRunStatusAttentionRequired, failed.Status)
					r.Len(failed.Actions, 1)
					a.Equal(openapi.TrailActionRunStatusFailed, failed.Actions[0].Status)
					r.NotNil(failed.Actions[0].Error)
					a.NotEmpty(*failed.Actions[0].Error)
					r.NotNil(failed.Actions[0].Target)
				})

				t.Run("history_and_scoped_run_lookup", func(t *testing.T) {
					r := require.New(t)
					a := assert.New(t)

					other := createPausedTrail(
						t,
						root,
						cl,
						adminSession,
						trailRobotAction(t, string(robot.JSON200.Id), "Other Trail"),
					)

					history := tests.AssertRequest(cl.TrailRunListWithResponse(root, definition.Id, adminSession))(t, http.StatusOK)
					r.Len(history.JSON200.Runs, 2)
					a.NotEqual(history.JSON200.Runs[0].Id, history.JSON200.Runs[1].Id)

					get := tests.AssertRequest(cl.TrailRunGetWithResponse(root, definition.Id, first.Id, adminSession))(t, http.StatusOK)
					a.Equal(first.Id, get.JSON200.Id)
					a.Equal(definition.Id, get.JSON200.TrailId)

					tests.AssertRequest(cl.TrailRunGetWithResponse(root, other.Id, first.Id, adminSession))(t, http.StatusNotFound)
				})
			}))
		}),
	)
}

func TestTrailRobotPermissionLoss(t *testing.T) {
	tests.ParallelIfNotSharedPostgres(t)

	integration.Test(
		t,
		&config.Config{LanguageModelProvider: "mock"},
		e2e.Setup(),
		testrobot.WithRobotSettings("mock/testdata/robot-completed.yaml"),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			accounts *account_writer.Writer,
			_ *trail_runtime.Runtime,
		) {
			lc.Append(fx.StartHook(func() {
				r := require.New(t)
				a := assert.New(t)

				adminCtx, _ := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				memberCtx, member := e2e.WithAccount(root, accounts, seed.Account_004_Loki)
				memberSession := sh.WithSession(memberCtx)

				role := tests.AssertRequest(cl.RoleCreateWithResponse(root, openapi.RoleCreateJSONRequestBody{
					Name: "Trail runner " + xid.New().String(),
					Permissions: openapi.PermissionList{
						openapi.MANAGETRAILS,
						openapi.USEROBOTS,
					},
				}, adminSession))(t, http.StatusOK)
				r.NotNil(role.JSON200)

				tests.AssertRequest(cl.AccountAddRoleWithResponse(root, member.Handle, role.JSON200.Id, adminSession))(t, http.StatusOK)

				model := openapi.RobotModelRef("mock/testdata/robot-completed.yaml")
				robot := tests.AssertRequest(cl.RobotCreateWithResponse(root, openapi.RobotCreateJSONRequestBody{
					Name:        "Permission loss Robot " + xid.New().String(),
					Description: "Must not run without current permission",
					Playbook:    "Complete the unattended instruction.",
					Model:       &model,
				}, adminSession))(t, http.StatusOK)
				r.NotNil(robot.JSON200)

				definition := createPausedTrail(
					t,
					root,
					cl,
					memberSession,
					trailRobotAction(t, string(robot.JSON200.Id), "This must not execute"),
				)

				tests.AssertRequest(cl.AccountRemoveRoleWithResponse(root, member.Handle, role.JSON200.Id, adminSession))(t, http.StatusOK)

				started := tests.AssertRequest(cl.TrailRunNowWithResponse(root, definition.Id, adminSession))(t, http.StatusOK)
				failed := waitForTrailRun(t, root, cl, adminSession, definition.Id, started.JSON200.Id)

				a.Equal(openapi.TrailRunStatusAttentionRequired, failed.Status)
				r.Len(failed.Actions, 1)
				a.Equal(openapi.TrailActionRunStatusFailed, failed.Actions[0].Status)
				r.NotNil(failed.Actions[0].Error)
				a.NotEmpty(*failed.Actions[0].Error)
				a.Nil(failed.Actions[0].Target)

				var notification openapi.Notification
				r.Eventually(func() bool {
					response, err := cl.NotificationListWithResponse(root, &openapi.NotificationListParams{}, memberSession)
					if err != nil || response == nil || response.JSON200 == nil {
						return false
					}

					found, ok := lo.Find(response.JSON200.Notifications, func(item openapi.Notification) bool {
						return item.Event == openapi.TrailRunAttention && item.Target != nil && *item.Target == failed.Id
					})
					if ok {
						notification = found
					}

					return ok
				}, 5*time.Second, 100*time.Millisecond)

				a.Equal(openapi.TrailRunAttention, notification.Event)
				r.NotNil(notification.Target)
				a.Equal(failed.Id, *notification.Target)
			}))
		}),
	)
}

func waitForTrailRun(
	t *testing.T,
	ctx context.Context,
	cl *openapi.ClientWithResponses,
	session openapi.RequestEditorFn,
	trailID openapi.Identifier,
	runID openapi.Identifier,
) openapi.TrailRun {
	t.Helper()

	var result openapi.TrailRun
	require.Eventually(t, func() bool {
		response, err := cl.TrailRunGetWithResponse(ctx, trailID, runID, session)
		if err != nil || response == nil || response.JSON200 == nil {
			return false
		}

		result = openapi.TrailRun(*response.JSON200)

		return result.Status == openapi.TrailRunStatusCompleted ||
			result.Status == openapi.TrailRunStatusAttentionRequired ||
			result.Status == openapi.TrailRunStatusCancelled
	}, 15*time.Second, 100*time.Millisecond)

	return result
}
