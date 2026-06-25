package chat_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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

func TestRobotCreateWithoutToolsPersistsEmptyToolList(t *testing.T) {
	t.Parallel()
	integration.Test(t,
		&config.Config{
			LanguageModelProvider: "mock",
		},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelAck),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				created := tests.AssertRequest(cl.RobotCreateWithResponse(root,
					openapi.RobotCreateJSONRequestBody{
						Name:        "tools-test-robot-" + xid.New().String(),
						Description: "testing tool persistence",
						Playbook:    "you are a test robot",
					},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, created.JSON200)
				robotID := created.JSON200.Id

				fetched := tests.AssertRequest(cl.RobotGetWithResponse(root,
					openapi.RobotIDParam(robotID),
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, fetched.JSON200)

				assert.Empty(t, fetched.JSON200.Tools)
			}))
		}),
	)
}

func TestRobotCreateToolPersistsTools(t *testing.T) {
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

				scriptName := "robot-chat-create-tools-" + xid.New().String() + ".yaml"
				scriptPath := filepath.Join("..", "scripts", scriptName)
				writeScript(t, scriptPath, `steps:
  - match:
      contains: "create robot with tools"
    respond:
      tool_calls:
        - id: call_create_tools_1
          name: robot_create
          args:
            name: "Tool Persisted Robot"
            description: "robot that should retain selected tools"
            playbook: "you use the tools you were configured with"
            tools:
              - library_request_page
              - get_library_page
  - match:
      tool_result: robot_create
    respond:
      text: "Created robot with tools."
      finish: "stop"
`)
				defer os.Remove(scriptPath)

				actor := tests.AssertRequest(cl.RobotCreateWithResponse(root,
					openapi.RobotCreateJSONRequestBody{
						Name:        "create-tools-actor-" + xid.New().String(),
						Description: "robot that creates another robot with tools",
						Playbook:    "you create robots",
						Model:       robotModelPtr("mock/../scripts/" + scriptName),
						Tools:       robotToolsPtr("robot_create"),
					},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, actor.JSON200)

				stream := doChat(t, root, ts, adminSession, xid.New().String(), string(actor.JSON200.Id), "create robot with tools")
				outputs := collectToolOutputs(stream)
				require.NotEmpty(t, outputs)

				output, ok := outputs[0].Output.(map[string]any)
				require.True(t, ok)
				createdID, ok := output["id"].(string)
				require.True(t, ok)

				fetched := tests.AssertRequest(cl.RobotGetWithResponse(root,
					openapi.RobotIDParam(createdID),
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, fetched.JSON200)

				assert.ElementsMatch(t, []string{"library_request_page", "get_library_page"}, fetched.JSON200.Tools)
			}))
		}),
	)
}

func TestRobotWithNoToolsRejectsToolCalls(t *testing.T) {
	t.Parallel()
	integration.Test(t,
		&config.Config{
			LanguageModelProvider: "mock",
		},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelLibraryTool),
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

				created := tests.AssertRequest(cl.RobotCreateWithResponse(root,
					openapi.RobotCreateJSONRequestBody{
						Name:        "restricted-robot-" + xid.New().String(),
						Description: "robot that should only be able to search",
						Playbook:    "you are a restricted test robot",
					},
					adminSession,
				))(t, http.StatusOK)
				robotID := string(created.JSON200.Id)

				vis := openapi.VisibilityPublished
				tests.AssertRequest(cl.NodeCreateWithResponse(root,
					openapi.NodeCreateJSONRequestBody{
						Name:       "Restricted Robot Test Page",
						Visibility: &vis,
					},
					adminSession,
				))(t, http.StatusOK)

				t.Run("restricted_tool_should_not_be_callable", func(t *testing.T) {
					a := assert.New(t)

					stream := doChat(t, root, ts, adminSession, xid.New().String(), robotID, "list pages")
					toolOutputs := collectToolOutputs(stream)

					// The mock LLM always attempts the call, but the ADK rejects it
					// because library_page_list is not in this robot's toolset.
					a.NotEmpty(toolOutputs, "expected tool-output-available events in stream")
					hasRejection := false
					for _, out := range toolOutputs {
						if output, ok := out.Output.(map[string]any); ok {
							if _, ok := output["error"].(string); ok {
								hasRejection = true
								break
							}
						}
					}
					a.True(hasRejection, "expected library_page_list to be rejected (not in robot toolset)")
				})
			}))
		}),
	)
}

func TestDefaultToolsContainsThrowAnError(t *testing.T) {
	t.Parallel()
	integration.Test(t,
		&config.Config{
			LanguageModelProvider: "mock",
		},
		e2e.Setup(),
		robot.WithRobotSettings("mock/../scripts/robot-chat-tool-error.yaml"),
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

				t.Run("throw_an_error_is_available_on_default_agent", func(t *testing.T) {
					a := assert.New(t)

					stream := doChat(t, root, ts, adminSession, xid.New().String(), "", "trigger error")
					toolNames := collectToolCalls(stream)

					a.Contains(toolNames, "throw_an_error")
				})
			}))
		}),
	)
}
