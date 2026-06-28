package chat_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/app/transports/sse"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
	robot_tests "github.com/Southclaws/storyden/tests/robot"
)

func TestRobotLibraryRequestPagePausesAndResumes(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{},
		e2e.Setup(),
		robot_tests.WithRobotSettings(mockModelAck),
		sse.Build(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
			sessionRepo *robot_session.Repository,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				scriptName := "robot-chat-request-page-" + xid.New().String() + ".yaml"
				scriptPath := filepath.Join("..", "scripts", scriptName)
				writeScript(t, scriptPath, `steps:
  - match:
      contains: "choose a page"
    respond:
      tool_calls:
        - id: call_request_page_1
          name: library_request_page
          args: {}
  - match:
      tool_result: library_request_page
    respond:
      text: "Selected page received."
      finish: "stop"
`)
				defer os.Remove(scriptPath)

				actor := tests.AssertRequest(cl.RobotCreateWithResponse(root,
					openapi.RobotCreateJSONRequestBody{
						Name:        "request-page-actor-" + xid.New().String(),
						Description: "robot that asks the user to select a Library page",
						Playbook:    "you ask the user to select a page when required",
						Model:       robotModelPtr("mock/../scripts/" + scriptName),
						Tools:       robotToolsPtr("library_request_page"),
					},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, actor.JSON200)

				sessionID := xid.New().String()
				first := doChat(t, root, ts, adminSession, sessionID, string(actor.JSON200.Id), "choose a page")
				inputs := collectToolInputs(first)
				require.Len(t, inputs, 1)
				assert.Equal(t, "library_request_page", string(inputs[0].ToolName))
				assert.Empty(t, collectToolOutputs(first), "page selection must pause until the client supplies a result")

				second := doChatToolOutput(t, root, ts, adminSession, sessionID, "", "library_request_page", string(inputs[0].ToolCallId), map[string]any{}, map[string]any{
					"id":          xid.New().String(),
					"slug":        "selected-page",
					"name":        "Selected Page",
					"description": "Chosen from the client UI",
				})

				assert.Empty(t, collectErrorParts(second))
				assert.Empty(t, collectToolOutputs(second), "the backend must not echo the client-supplied tool result back to the UI")
				assert.Equal(t, "Selected page received.", strings.Join(collectTextDeltas(second), ""))

				parsedSessionID, err := xid.FromString(sessionID)
				require.NoError(t, err)

				sess, _, err := sessionRepo.Get(root, robot.SessionID(parsedSessionID), robot.NewMessageCursorParams(opt.NewEmpty[robot.MessageID](), 50))
				require.NoError(t, err)

				var matchingToolResponses int
				for _, message := range sess.Messages {
					require.NotNil(t, message.Event.LLMResponse.Content)
					for _, part := range message.Event.LLMResponse.Content.Parts {
						if part == nil || part.FunctionResponse == nil {
							continue
						}

						assert.False(t, isPendingClientToolResponse(part.FunctionResponse.Response), "client-side pending marker must not be persisted as a real tool result")

						if part.FunctionResponse.ID != string(inputs[0].ToolCallId) || part.FunctionResponse.Name != "library_request_page" {
							continue
						}

						matchingToolResponses++
						assert.Equal(t, "selected-page", part.FunctionResponse.Response["slug"])
						assert.NotContains(t, part.FunctionResponse.Response, "selection", "model-facing elicitation wrapper must not be persisted")
						assert.NotContains(t, part.PartMetadata, "storyden_elicitation_hydrated", "model-facing hydration marker must not be persisted")
						assert.Equal(t, "user", message.Event.Author)
						assert.Equal(t, "user", message.Event.LLMResponse.Content.Role)
					}
				}
				assert.Equal(t, 1, matchingToolResponses, "expected exactly one persisted tool result for the selected page")
			}))
		}),
	)
}

func isPendingClientToolResponse(response map[string]any) bool {
	pending, _ := response["_client_side_pending"].(bool)
	return pending
}
