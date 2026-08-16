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
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
	"github.com/Southclaws/storyden/tests/robot"
)

func TestDenbotDelegatesToSpecialistInSameSession(t *testing.T) {
	t.Parallel()

	integration.Test(t,
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
      finish: "stop"
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
				specialist := tests.AssertRequest(cl.RobotCreateWithResponse(root,
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
      tool_result: `+agentName+`
    respond:
      text: "Denbot synthesized the specialist result."
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
				stream := doChat(t, root, ts, adminSession, sessionID, "", "delegate this")
				assert.NotContains(t, collectToolCalls(stream), agentName)
				assert.Equal(t, "Denbot synthesized the specialist result.", strings.Join(collectTextDeltas(stream), ""))
				assert.NotContains(t, strings.Join(collectReasoningDeltas(stream), ""), "Specialist found the delegated result.")

				delegations := collectDelegations(stream)
				require.NotEmpty(t, delegations)
				delegation := delegations[len(delegations)-1]
				assert.Equal(t, "call_delegate_1", delegation.StreamID)
				assert.Equal(t, "call_delegate_1", delegation.CallID)
				assert.Equal(t, specialistID, string(delegation.Robot.Id))
				assert.Equal(t, specialistName, delegation.Robot.Name)
				assert.Equal(t, "Find the delegated result.", delegation.Request)
				assert.Equal(t, "completed", delegation.Status)
				var delegatedText strings.Builder
				for _, message := range delegation.Messages {
					for _, part := range message.Parts {
						textPart, err := part.AsTextUIPart()
						if err == nil {
							delegatedText.WriteString(textPart.Text)
						}
					}
				}
				assert.Contains(t, delegatedText.String(), "Specialist found the delegated result.")

				session := tests.AssertRequest(cl.RobotSessionGetWithResponse(root,
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
					assert.Equal(t, "content_search", input.ToolName)
					toolCallIDs[input.ToolCallId] = struct{}{}
				}
				assert.Equal(t, map[string]struct{}{"call_specialist_search": {}}, toolCallIDs)
			}))
		}),
	)
}
