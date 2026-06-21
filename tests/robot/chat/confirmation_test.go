package chat_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/app/transports/sse"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
	"github.com/Southclaws/storyden/tests/robot"
)

func TestRobotDeleteRequiresConfirmation(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{
			LanguageModelProvider: "mock",
		},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelAck),
		sse.Build(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				victim := tests.AssertRequest(cl.RobotCreateWithResponse(root,
					openapi.RobotCreateJSONRequestBody{
						Name:        "delete-victim-" + xid.New().String(),
						Description: "robot that should require confirmation before deletion",
						Playbook:    "you are a delete victim",
					},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, victim.JSON200)
				victimID := string(victim.JSON200.Id)

				scriptName := "robot-chat-delete-" + xid.New().String() + ".yaml"
				scriptPath := filepath.Join("..", "scripts", scriptName)
				writeScript(t, scriptPath, `steps:
  - match:
      contains: "delete the robot"
    respond:
      tool_calls:
        - id: call_delete_1
          name: robot_delete
          args:
            id: "`+victimID+`"
  - match:
      tool_result: robot_delete
    respond:
      text: "Delete flow finished."
      finish: "stop"
`)
				defer os.Remove(scriptPath)

				deleter := tests.AssertRequest(cl.RobotCreateWithResponse(root,
					openapi.RobotCreateJSONRequestBody{
						Name:        "delete-actor-" + xid.New().String(),
						Description: "robot that asks to delete another robot",
						Playbook:    "you are a delete actor",
						Model:       robotModelPtr("mock/../scripts/" + scriptName),
						Tools:       robotToolsPtr("robot_delete"),
					},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, deleter.JSON200)

				sessionID := xid.New().String()
				first := doChat(t, root, ts, adminSession, sessionID, string(deleter.JSON200.Id), "delete the robot")
				inputs := collectToolInputs(first)
				require.Len(t, inputs, 1)
				assert.Equal(t, "robot_delete", string(inputs[0].ToolName))
				assert.Empty(t, collectToolOutputs(first), "delete must pause before server-side execution")

				sessionResponse := tests.AssertRequest(cl.RobotSessionGetWithResponse(root,
					openapi.RobotSessionIDParam(sessionID),
					&openapi.RobotSessionGetParams{},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, sessionResponse.JSON200)
				historicalInputs := collectSessionToolInputs(sessionResponse.JSON200.MessageList.Messages)
				require.Len(t, historicalInputs, 1)
				assert.Equal(t, "robot_delete", string(historicalInputs[0].ToolName))
				assert.Equal(t, string(inputs[0].ToolCallId), historicalInputs[0].ToolCallId)

				assert.Equal(t, http.StatusConflict, doChatStatus(t, root, ts, adminSession, sessionID, "perfect"))
				tests.AssertRequest(cl.RobotGetWithResponse(root, openapi.RobotIDParam(victimID), adminSession))(t, http.StatusOK)

				input := inputs[0].Input.(map[string]any)
				second := doChatToolOutput(t, root, ts, adminSession, sessionID, string(deleter.JSON200.Id), "robot_delete", string(inputs[0].ToolCallId), input, map[string]any{
					"_storyden_confirmation": map[string]any{
						"approved": true,
						"id":       victimID,
					},
				})

				assert.Equal(t, "Delete flow finished.", strings.Join(collectTextDeltas(second), ""))
				tests.AssertRequest(cl.RobotGetWithResponse(root, openapi.RobotIDParam(victimID), adminSession))(t, http.StatusNotFound)
			}))
		}),
	)
}

func TestRobotDeleteConfirmationAfterPriorToolCallWithoutRobotID(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{
			LanguageModelProvider: "mock",
		},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelAck),
		sse.Build(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				victim := tests.AssertRequest(cl.RobotCreateWithResponse(root,
					openapi.RobotCreateJSONRequestBody{
						Name:        "delete-victim-" + xid.New().String(),
						Description: "robot that should survive until approval",
						Playbook:    "you are a delete victim",
					},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, victim.JSON200)
				victimID := string(victim.JSON200.Id)

				scriptName := "robot-chat-delete-after-create-" + xid.New().String() + ".yaml"
				scriptPath := filepath.Join("..", "scripts", scriptName)
				writeScript(t, scriptPath, `steps:
  - match:
      contains: "create a robot"
    respond:
      tool_calls:
        - id: call_create_1
          name: robot_create
          args:
            name: "Temporary Test Robot"
            description: "temporary robot created before deletion"
            playbook: "you are temporary"
  - match:
      tool_result: robot_create
    respond:
      text: "Created Robot."
      finish: "stop"
  - match:
      contains: "delete it"
    respond:
      tool_calls:
        - id: call_delete_1
          name: robot_delete
          args:
            id: "`+victimID+`"
  - match:
      tool_result: robot_delete
    respond:
      text: "Delete flow finished."
      finish: "stop"
`)
				defer os.Remove(scriptPath)

				actor := tests.AssertRequest(cl.RobotCreateWithResponse(root,
					openapi.RobotCreateJSONRequestBody{
						Name:        "delete-after-create-actor-" + xid.New().String(),
						Description: "robot that creates before deleting",
						Playbook:    "you create first, then delete later",
						Model:       robotModelPtr("mock/../scripts/" + scriptName),
						Tools:       robotToolsPtr("robot_create", "robot_delete"),
					},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, actor.JSON200)
				actorID := string(actor.JSON200.Id)

				sessionID := xid.New().String()
				create := doChat(t, root, ts, adminSession, sessionID, actorID, "create a robot")
				assert.Equal(t, "Created Robot.", strings.Join(collectTextDeltas(create), ""))
				require.NotEmpty(t, collectToolOutputs(create))

				deleteTurn := doChat(t, root, ts, adminSession, sessionID, actorID, "delete it")
				inputs := collectToolInputs(deleteTurn)
				require.Len(t, inputs, 1)
				assert.Equal(t, "robot_delete", string(inputs[0].ToolName))
				assert.Empty(t, collectToolOutputs(deleteTurn), "delete must pause before server-side execution")

				tests.AssertRequest(cl.RobotGetWithResponse(root, openapi.RobotIDParam(victimID), adminSession))(t, http.StatusOK)

				input := inputs[0].Input.(map[string]any)
				confirmed := doChatToolOutput(t, root, ts, adminSession, sessionID, "", "robot_delete", string(inputs[0].ToolCallId), input, map[string]any{
					"_storyden_confirmation": map[string]any{
						"approved": true,
						"id":       victimID,
					},
				})

				assert.Empty(t, collectErrorParts(confirmed))
				assert.Equal(t, "Delete flow finished.", strings.Join(collectTextDeltas(confirmed), ""))
				tests.AssertRequest(cl.RobotGetWithResponse(root, openapi.RobotIDParam(victimID), adminSession))(t, http.StatusNotFound)
			}))
		}),
	)
}

func TestRobotDeleteMultipleConfirmationsInSameTurn(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{
			LanguageModelProvider: "mock",
		},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelAck),
		sse.Build(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				firstVictim := tests.AssertRequest(cl.RobotCreateWithResponse(root,
					openapi.RobotCreateJSONRequestBody{
						Name:        "multi-delete-victim-one-" + xid.New().String(),
						Description: "first robot that should require confirmation before deletion",
						Playbook:    "you are the first delete victim",
					},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, firstVictim.JSON200)
				firstVictimID := string(firstVictim.JSON200.Id)

				secondVictim := tests.AssertRequest(cl.RobotCreateWithResponse(root,
					openapi.RobotCreateJSONRequestBody{
						Name:        "multi-delete-victim-two-" + xid.New().String(),
						Description: "second robot that should require confirmation before deletion",
						Playbook:    "you are the second delete victim",
					},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, secondVictim.JSON200)
				secondVictimID := string(secondVictim.JSON200.Id)

				scriptName := "robot-chat-multi-delete-" + xid.New().String() + ".yaml"
				scriptPath := filepath.Join("..", "scripts", scriptName)
				writeScript(t, scriptPath, `steps:
  - match:
      contains: "delete both robots"
    respond:
      tool_calls:
        - id: call_delete_1
          name: robot_delete
          args:
            id: "`+firstVictimID+`"
        - id: call_delete_2
          name: robot_delete
          args:
            id: "`+secondVictimID+`"
  - match:
      tool_result: robot_delete
    respond:
      text: "Both delete flows finished."
      finish: "stop"
`)
				defer os.Remove(scriptPath)

				deleter := tests.AssertRequest(cl.RobotCreateWithResponse(root,
					openapi.RobotCreateJSONRequestBody{
						Name:        "multi-delete-actor-" + xid.New().String(),
						Description: "robot that asks to delete two robots",
						Playbook:    "you are a multi-delete actor",
						Model:       robotModelPtr("mock/../scripts/" + scriptName),
						Tools:       robotToolsPtr("robot_delete"),
					},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, deleter.JSON200)

				sessionID := xid.New().String()
				first := doChat(t, root, ts, adminSession, sessionID, string(deleter.JSON200.Id), "delete both robots")
				inputs := collectToolInputs(first)
				require.Len(t, inputs, 2)
				assert.Empty(t, collectToolOutputs(first), "deletes must pause before server-side execution")

				tests.AssertRequest(cl.RobotGetWithResponse(root, openapi.RobotIDParam(firstVictimID), adminSession))(t, http.StatusOK)
				tests.AssertRequest(cl.RobotGetWithResponse(root, openapi.RobotIDParam(secondVictimID), adminSession))(t, http.StatusOK)

				partial := doChatToolOutputsStatus(t, root, ts, adminSession, sessionID, string(deleter.JSON200.Id), []map[string]any{{
					"type":       "tool-robot_delete",
					"state":      "output-available",
					"toolCallId": string(inputs[0].ToolCallId),
					"toolName":   "robot_delete",
					"input":      inputs[0].Input.(map[string]any),
					"output": map[string]any{
						"_storyden_confirmation": map[string]any{
							"approved": true,
							"id":       firstVictimID,
						},
					},
				}})
				assert.Equal(t, http.StatusConflict, partial)
				tests.AssertRequest(cl.RobotGetWithResponse(root, openapi.RobotIDParam(firstVictimID), adminSession))(t, http.StatusOK)
				tests.AssertRequest(cl.RobotGetWithResponse(root, openapi.RobotIDParam(secondVictimID), adminSession))(t, http.StatusOK)

				confirmed := doChatToolOutputs(t, root, ts, adminSession, sessionID, string(deleter.JSON200.Id), []map[string]any{
					{
						"type":       "tool-robot_delete",
						"state":      "output-available",
						"toolCallId": string(inputs[0].ToolCallId),
						"toolName":   "robot_delete",
						"input":      inputs[0].Input.(map[string]any),
						"output": map[string]any{
							"_storyden_confirmation": map[string]any{
								"approved": true,
								"id":       firstVictimID,
							},
						},
					},
					{
						"type":       "tool-robot_delete",
						"state":      "output-available",
						"toolCallId": string(inputs[1].ToolCallId),
						"toolName":   "robot_delete",
						"input":      inputs[1].Input.(map[string]any),
						"output": map[string]any{
							"_storyden_confirmation": map[string]any{
								"approved": true,
								"id":       secondVictimID,
							},
						},
					},
				})

				assert.Empty(t, collectErrorParts(confirmed))
				assert.Equal(t, "Both delete flows finished.", strings.Join(collectTextDeltas(confirmed), ""))
				tests.AssertRequest(cl.RobotGetWithResponse(root, openapi.RobotIDParam(firstVictimID), adminSession))(t, http.StatusNotFound)
				tests.AssertRequest(cl.RobotGetWithResponse(root, openapi.RobotIDParam(secondVictimID), adminSession))(t, http.StatusNotFound)
			}))
		}),
	)
}

func TestRobotSessionCurrentRobotFollowsRobotSwitch(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{
			LanguageModelProvider: "mock",
		},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelAck),
		sse.Build(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				targetScriptName := "robot-chat-current-target-" + xid.New().String() + ".yaml"
				targetScriptPath := filepath.Join("..", "scripts", targetScriptName)
				writeScript(t, targetScriptPath, `steps:
  - match:
      contains: "target followup"
    respond:
      text: "Target robot handled followup."
      finish: "stop"
`)
				defer os.Remove(targetScriptPath)

				target := tests.AssertRequest(cl.RobotCreateWithResponse(root,
					openapi.RobotCreateJSONRequestBody{
						Name:        "current-target-" + xid.New().String(),
						Description: "target robot for current session state",
						Playbook:    "you are the target robot",
						Model:       robotModelPtr("mock/../scripts/" + targetScriptName),
					},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, target.JSON200)
				targetID := string(target.JSON200.Id)

				actorScriptName := "robot-chat-current-actor-" + xid.New().String() + ".yaml"
				actorScriptPath := filepath.Join("..", "scripts", actorScriptName)
				writeScript(t, actorScriptPath, `steps:
  - match:
      contains: "switch to target"
    respond:
      tool_calls:
        - id: call_switch_1
          name: robot_switch
          args:
            robot_id: "`+targetID+`"
  - match:
      contains: "ROBOT SWITCH"
    respond:
      text: "Switched to target."
      finish: "stop"
  - match:
      any: true
    respond:
      text: "Actor robot stayed active."
      finish: "stop"
`)
				defer os.Remove(actorScriptPath)

				actor := tests.AssertRequest(cl.RobotCreateWithResponse(root,
					openapi.RobotCreateJSONRequestBody{
						Name:        "current-actor-" + xid.New().String(),
						Description: "actor robot that switches to target",
						Playbook:    "you switch to the target robot",
						Model:       robotModelPtr("mock/../scripts/" + actorScriptName),
						Tools:       robotToolsPtr("robot_switch"),
					},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, actor.JSON200)

				sessionID := xid.New().String()
				first := doChat(t, root, ts, adminSession, sessionID, string(actor.JSON200.Id), "switch to target")
				inputs := collectToolInputs(first)
				require.Len(t, inputs, 1)
				assert.Equal(t, "robot_switch", string(inputs[0].ToolName))

				switched := doChatToolOutput(t, root, ts, adminSession, sessionID, string(actor.JSON200.Id), "robot_switch", string(inputs[0].ToolCallId), inputs[0].Input.(map[string]any), map[string]any{
					"success":  true,
					"robot_id": targetID,
				})
				assert.Empty(t, strings.Join(collectTextDeltas(switched), ""))

				currentSession := tests.AssertRequest(cl.RobotSessionGetWithResponse(root,
					openapi.RobotSessionIDParam(sessionID),
					&openapi.RobotSessionGetParams{},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, currentSession.JSON200)
				require.NotNil(t, currentSession.JSON200.ActiveRobotId)
				assert.Equal(t, targetID, *currentSession.JSON200.ActiveRobotId)

				followup := doChat(t, root, ts, adminSession, sessionID, targetID, "target followup")
				assert.Equal(t, "Target robot handled followup.", strings.Join(collectTextDeltas(followup), ""))
			}))
		}),
	)
}

func robotModelPtr(model string) *openapi.RobotModelRef {
	v := openapi.RobotModelRef(model)
	return &v
}
