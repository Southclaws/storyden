package chat_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
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

// TestRobotToolCallError verifies that when a tool returns an error the ADK
// wraps it as {"error": "..."} in the FunctionResponse, the SSE handler emits
// a tool-output-available part carrying that error, and the agent then
// continues to produce a text response based on the tool result.
func TestRobotToolCallError(t *testing.T) {
	t.Parallel()
	integration.Test(t,
		&config.Config{
			LanguageModelProvider: "mock",
		},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelToolError),
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

				rob := tests.AssertRequest(cl.RobotCreateWithResponse(root, openapi.RobotCreateJSONRequestBody{
					Name:        "error-robot-" + xid.New().String(),
					Description: "robot for error tests",
					Playbook:    "you are a test robot",
					Tools:       robotToolsPtr("throw_an_error"),
				}, adminSession))(t, http.StatusOK)
				robotID := string(rob.JSON200.Id)

				t.Run("tool_error_wrapped_as_result", func(t *testing.T) {
					a := assert.New(t)

					stream := doChat(t, root, ts, adminSession, xid.New().String(), robotID, "trigger error")

					toolNames := collectToolCalls(stream)
					toolOutputs := collectToolOutputs(stream)
					textDeltas := collectTextDeltas(stream)

					a.Contains(toolNames, "throw_an_error")

					// The ADK wraps tool errors as {"error": "..."} in the
					// FunctionResponse.Response, which becomes the output field.
					errorSeen := false
					for _, out := range toolOutputs {
						if output, ok := out.Output.(map[string]any); ok {
							if errVal, ok := output["error"].(string); ok && strings.Contains(errVal, "intentional tool error") {
								errorSeen = true
								break
							}
						}
					}
					a.True(errorSeen, "expected tool-output-available to carry the error")

					// The mock's second step fires on the tool result and responds
					// with a text message, confirming the agent loop continued.
					a.Equal("The tool returned an error.", strings.Join(textDeltas, ""))
				})
			}))
		}),
	)
}

// TestRobotLLMError verifies that when the mock LLM itself returns an error
// (simulating a provider failure) the SSE handler emits an "error" stream part
// and closes the stream without a finish-message.
func TestRobotLLMError(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{
			LanguageModelProvider: "mock",
		},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelLLMError),
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

				t.Run("llm_error_emits_error_part", func(t *testing.T) {
					a := assert.New(t)

					stream := doChat(t, root, ts, adminSession, xid.New().String(), "", "trigger llm error")

					errorParts := collectErrorParts(stream)
					finishParts := collectPartsOfType(stream, "finish-message")

					a.NotEmpty(errorParts, "expected at least one error part in the stream")
					a.Empty(finishParts, "stream should end on error without a finish-message")
				})
			}))
		}),
	)
}
