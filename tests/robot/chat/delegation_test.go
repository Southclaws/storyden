package chat_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	ent_robot_session_message "github.com/Southclaws/storyden/internal/ent/robotsessionmessage"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
	"github.com/Southclaws/storyden/tests/robot"
)

func TestDenbotDelegatesToSpecialistInSameSession(t *testing.T) {
	t.Parallel()

	integration.Test(
		t,
		&config.Config{LanguageModelProvider: "mock"},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelAck),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
			settingsRepo *settings.SettingsRepository,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				specialistScriptName := "robot-chat-specialist-tools-" + xid.New().String() + ".yaml"
				specialistScriptPath := filepath.Join("..", "scripts", specialistScriptName)
				writeScript(t, specialistScriptPath, `steps:
  - match:
      tool_result: content_search
    respond:
      text: "Specialist found the delegated result."
      tool_calls:
        - id: call_specialist_finish
          name: robot_run_finish
          args:
            status: completed
            summary: "Specialist found the delegated result."
  - match:
      any: true
    respond:
      tool_calls:
        - id: call_specialist_search
          name: content_search
          args:
            query: "delegation fixture with no results"
`)
				defer os.Remove(specialistScriptPath)

				specialistName := "research-specialist-" + xid.New().String()
				specialist := tests.AssertRequest(cl.RobotCreateWithResponse(
					root,
					openapi.RobotCreateJSONRequestBody{
						Name:        specialistName,
						Description: "Finds a bounded piece of evidence for Denbot",
						Playbook:    "Complete only the delegated research request and return evidence.",
						Model:       robotModelPtr("mock/../scripts/" + specialistScriptName),
						Tools:       robotToolsPtr("content_search"),
					},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, specialist.JSON200)

				specialistID := string(specialist.JSON200.Id)
				agentName := "robot_" + specialistID
				scriptName := "robot-chat-delegation-" + xid.New().String() + ".yaml"
				scriptPath := filepath.Join("..", "scripts", scriptName)
				writeScript(t, scriptPath, `steps:
  - match:
      contains: "asynchronous specialist result"
    respond:
      text: "Denbot synthesized the specialist result."
      finish: "stop"
  - match:
      tool_result: `+agentName+`
      tool_result_status: pending
    respond:
      text: "The specialist is working asynchronously."
      finish: "stop"
  - match:
      contains: "delegate this"
    respond:
      tool_calls:
        - id: call_delegate_1
          name: `+agentName+`
          args:
            request: "Find the delegated result."
`)
				defer os.Remove(scriptPath)
				require.NoError(t, robot.SetRobotSettings(root, settingsRepo, "mock/../scripts/"+scriptName))

				sessionID := xid.New().String()
				accepted := enqueueChatMessage(t, root, ts, adminSession, sessionID, "", "delegate this")
				readCtx, cancel := context.WithTimeout(root, 20*time.Second)
				defer cancel()
				completedTurns := 0
				read, err := robot.ReadDurableJSONUntil(
					readCtx,
					ts.URL+accepted.location,
					adminSession,
					"-1",
					func(event openapi.RobotSessionStreamEvent) bool {
						if event.EventKind == openapi.TurnCompleted {
							completedTurns++
						}
						return completedTurns == 3
					},
				)
				require.NoError(t, err)
				stream := &fullResponse{}
				for _, event := range read.Items {
					stream.parts = append(stream.parts, event.Parts...)
				}
				assert.NotContains(t, collectToolCalls(stream), agentName)
				assert.Contains(t, strings.Join(collectTextDeltas(stream), ""), "The specialist is working asynchronously.")
				assert.Contains(t, strings.Join(collectTextDeltas(stream), ""), "Denbot synthesized the specialist result.")
				assert.NotContains(t, strings.Join(collectReasoningDeltas(stream), ""), "Specialist found the delegated result.")

				delegations := collectDelegations(stream)
				require.NotEmpty(t, delegations)
				delegation := delegations[0]
				assert.Equal(t, "call_delegate_1", delegation.StreamID)
				assert.Equal(t, "call_delegate_1", delegation.CallID)
				assert.Equal(t, specialistID, string(delegation.Robot.Id))
				assert.Equal(t, specialistName, delegation.Robot.Name)
				assert.Equal(t, "Find the delegated result.", delegation.Request)
				assert.Condition(t, func() bool {
					for _, observed := range delegations {
						if observed.Status == "completed" {
							return true
						}
					}
					return false
				}, "delegation should complete in a later turn")
				var delegatedText strings.Builder
				for _, observed := range delegations {
					for _, message := range observed.Messages {
						for _, part := range message.Parts {
							textPart, err := part.AsTextUIPart()
							if err == nil {
								delegatedText.WriteString(textPart.Text)
							}
						}
					}
				}
				assert.Contains(t, delegatedText.String(), "Specialist found the delegated result.")

				session := tests.AssertRequest(cl.RobotSessionGetWithResponse(
					root,
					openapi.RobotSessionIDParam(sessionID),
					&openapi.RobotSessionGetParams{},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, session.JSON200)

				var delegatedMessages []openapi.RobotSessionMessage
				for i := range session.JSON200.MessageList.Messages {
					message := &session.JSON200.MessageList.Messages[i]
					if message.Robot != nil && string(message.Robot.Id) == specialistID {
						delegatedMessages = append(delegatedMessages, *message)
					}
				}
				require.NotEmpty(t, delegatedMessages, "specialist output should be attributed inside the root session")

				var persistedDelegatedText strings.Builder
				for _, message := range delegatedMessages {
					require.NotNil(t, message.Branch)
					assert.Contains(t, *message.Branch, agentName)
					require.NotNil(t, message.IsolationScope)
					assert.Equal(t, "call_delegate_1", *message.IsolationScope)
					for _, part := range message.Parts {
						textPart, err := part.AsTextUIPart()
						if err == nil {
							persistedDelegatedText.WriteString(textPart.Text)
						}
					}
				}
				assert.Contains(t, persistedDelegatedText.String(), "Specialist found the delegated result.")

				toolInputs := collectSessionToolInputs(delegatedMessages)
				toolCallIDs := make(map[string]struct{})
				for _, input := range toolInputs {
					if input.ToolName == "robot_run_finish" {
						continue
					}
					assert.Equal(t, "content_search", input.ToolName)
					toolCallIDs[input.ToolCallId] = struct{}{}
				}
				assert.Equal(t, map[string]struct{}{"call_specialist_search": {}}, toolCallIDs)
			}))
		}),
	)
}

func TestFailedDelegationDoesNotCaptureLaterUserMessages(t *testing.T) {
	t.Parallel()

	integration.Test(
		t,
		&config.Config{LanguageModelProvider: "mock"},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelAck),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
			settingsRepo *settings.SettingsRepository,
			db *ent.Client,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				specialistScriptName := "robot-chat-failing-specialist-" + xid.New().String() + ".yaml"
				specialistScriptPath := filepath.Join("..", "scripts", specialistScriptName)
				writeScript(t, specialistScriptPath, `steps:
  - match:
      any: true
    respond:
      delay_ms: 500
      error: "invalid_request_error: openai api error"
`)
				defer os.Remove(specialistScriptPath)

				specialist := tests.AssertRequest(cl.RobotCreateWithResponse(
					root,
					openapi.RobotCreateJSONRequestBody{
						Name:        "failing-specialist-" + xid.New().String(),
						Description: "Fails a delegated request after accepting it",
						Playbook:    "Attempt only the delegated request.",
						Model:       robotModelPtr("mock/../scripts/" + specialistScriptName),
					},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, specialist.JSON200)

				agentName := "robot_" + string(specialist.JSON200.Id)
				rootScriptName := "robot-chat-failed-delegation-root-" + xid.New().String() + ".yaml"
				rootScriptPath := filepath.Join("..", "scripts", rootScriptName)
				writeScript(t, rootScriptPath, `steps:
  - match:
      contains: "asynchronous specialist result"
    respond:
      text: "The failed specialist result was handled."
      finish: "stop"
  - match:
      contains: "unrelated queued message"
    respond:
      text: "I can read the unrelated queued message."
      finish: "stop"
  - match:
      tool_result: `+agentName+`
      tool_result_status: pending
    respond:
      text: "The specialist is working asynchronously."
      finish: "stop"
  - match:
      contains: "delegate the failing task"
    respond:
      tool_calls:
        - id: call_failing_delegate
          name: `+agentName+`
          args:
            request: "Attempt the failing task."
`)
				defer os.Remove(rootScriptPath)
				require.NoError(t, robot.SetRobotSettings(root, settingsRepo, "mock/../scripts/"+rootScriptName))

				sessionID := xid.New().String()
				delegation := enqueueChatMessage(t, root, ts, adminSession, sessionID, "", "delegate the failing task")
				firstTurn := readDurableChatParts(t, root, ts, adminSession, delegation)
				assert.Contains(t, strings.Join(collectTextDeltas(firstTurn), ""), "The specialist is working asynchronously.")

				followUp := enqueueChatMessage(t, root, ts, adminSession, sessionID, "", "unrelated queued message")
				followUpTurn := readDurableChatParts(t, root, ts, adminSession, followUp)
				assert.Contains(t, strings.Join(collectTextDeltas(followUpTurn), ""), "I can read the unrelated queued message.")

				inputID, err := xid.FromString(string(followUp.reference.MessageId))
				require.NoError(t, err)
				input, err := db.RobotSessionInput.Get(root, inputID)
				require.NoError(t, err)
				require.NotNil(t, input.TurnID)
				sessionXID, err := xid.FromString(sessionID)
				require.NoError(t, err)
				messages, err := db.RobotSessionMessage.Query().Where(
					ent_robot_session_message.SessionIDEQ(sessionXID),
					ent_robot_session_message.TurnIDEQ(*input.TurnID),
					ent_robot_session_message.HiddenFromProjectionEQ(true),
				).All(root)
				require.NoError(t, err)

				var runtimeInput *ent.RobotSessionMessage
				for _, message := range messages {
					if message.EventData.ADK().Author == "user" {
						runtimeInput = message
						break
					}
				}
				require.NotNil(t, runtimeInput)
				assert.Empty(t, runtimeInput.IsolationScope)
			}))
		}),
	)
}
