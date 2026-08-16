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

func TestRobotCreateWithoutToolsetsPersistsEmptyToolsetList(t *testing.T) {
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

				assert.Empty(t, fetched.JSON200.Toolsets)
			}))
		}),
	)
}

func TestRobotCreateToolPersistsToolsets(t *testing.T) {
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
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
			settingsRepo *settings.SettingsRepository,
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
              - thread_get
              - throw_an_error
            toolsets:
              - system.discussions
  - match:
      tool_result: robot_create
    respond:
      text: "Created robot with tools."
      finish: "stop"
`)
				defer os.Remove(scriptPath)
				require.NoError(t, robot.SetRobotSettings(root, settingsRepo, "mock/../scripts/"+scriptName))

				actor := tests.AssertRequest(cl.RobotCreateWithResponse(root,
					openapi.RobotCreateJSONRequestBody{
						Name:        "create-tools-actor-" + xid.New().String(),
						Description: "robot that creates another robot with tools",
						Playbook:    "you create robots",
						Model:       robotModelPtr("mock/../scripts/" + scriptName),
						Toolsets:    robotToolsetsPtr("system.robot_studio"),
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

				assert.Equal(t, openapi.RobotToolNameList{"throw_an_error"}, fetched.JSON200.Tools)
				assert.Equal(t, openapi.RobotToolsetRefList{"system.discussions"}, fetched.JSON200.Toolsets)
			}))
		}),
	)
}

func TestCustomRobotCanRunAnIndividualTool(t *testing.T) {
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
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				scriptName := "robot-chat-individual-tool-" + xid.New().String() + ".yaml"
				scriptPath := filepath.Join("..", "scripts", scriptName)
				writeScript(t, scriptPath, `steps:
  - match:
      contains: "list configured tools"
    respond:
      tool_calls:
        - id: call_individual_tool_1
          name: throw_an_error
          args: {}
  - match:
      tool_result: throw_an_error
    respond:
      text: "Individual tool executed."
      finish: "stop"
`)
				defer os.Remove(scriptPath)

				created := tests.AssertRequest(cl.RobotCreateWithResponse(root,
					openapi.RobotCreateJSONRequestBody{
						Name:        "individual-tool-runtime-" + xid.New().String(),
						Description: "robot with one directly assigned tool",
						Playbook:    "use the directly assigned tool",
						Model:       robotModelPtr("mock/../scripts/" + scriptName),
						Tools:       robotToolsPtr("throw_an_error"),
					},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, created.JSON200)

				stream := doChat(t, root, ts, adminSession, xid.New().String(), string(created.JSON200.Id), "list configured tools")
				assert.Contains(t, collectToolCalls(stream), "throw_an_error")
				assert.Equal(t, "Individual tool executed.", strings.Join(collectTextDeltas(stream), ""))
			}))
		}),
	)
}

func TestRobotChatStartsAndContinuesWithRequestedCustomRobot(t *testing.T) {
	t.Parallel()
	integration.Test(t,
		&config.Config{
			LanguageModelProvider: "mock",
		},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelSimple),
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
						Name:        "direct-chat-robot-" + xid.New().String(),
						Description: "robot selected directly as the session root",
						Playbook:    "respond as the directly selected custom Robot",
						Model:       robotModelPtr("mock/../scripts/robot-chat-specialist.yaml"),
					},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, created.JSON200)
				robotID := string(created.JSON200.Id)
				sessionID := xid.New().String()

				first := doChat(t, root, ts, adminSession, sessionID, robotID, "hello specialist")
				assert.Equal(t, "Specialist found the delegated result.", strings.Join(collectTextDeltas(first), ""))

				second := doChat(t, root, ts, adminSession, sessionID, "", "continue")
				assert.Equal(t, "Specialist found the delegated result.", strings.Join(collectTextDeltas(second), ""))
				assert.Equal(t, http.StatusConflict, doChatWithRobotStatus(t, root, ts, adminSession, sessionID, "denbot", "switch roots"))

				session := tests.AssertRequest(cl.RobotSessionGetWithResponse(root,
					openapi.RobotSessionIDParam(sessionID),
					&openapi.RobotSessionGetParams{},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, session.JSON200)
				require.NotNil(t, session.JSON200.RootRobotId)
				assert.Equal(t, robotID, *session.JSON200.RootRobotId)

				var customResponseFound bool
				for _, message := range session.JSON200.MessageList.Messages {
					if message.Robot != nil && string(message.Robot.Id) == robotID {
						customResponseFound = true
						break
					}
				}
				assert.True(t, customResponseFound, "custom Robot response should retain its actor attribution")
			}))
		}),
	)
}

func TestDenbotIncludesToolsetDiscovery(t *testing.T) {
	t.Parallel()
	integration.Test(t,
		&config.Config{
			LanguageModelProvider: "mock",
		},
		e2e.Setup(),
		robot.WithRobotSettings("mock/../scripts/robot-chat-toolset-search.yaml"),
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

				t.Run("toolset_search_is_available_on_denbot", func(t *testing.T) {
					a := assert.New(t)

					stream := doChat(t, root, ts, adminSession, xid.New().String(), "", "trigger error")
					toolNames := collectToolCalls(stream)

					a.Contains(toolNames, "toolset_search")
					a.Equal("Toolsets searched.", strings.Join(collectTextDeltas(stream), ""))
				})
			}))
		}),
	)
}
