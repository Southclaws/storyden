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
	"github.com/Southclaws/storyden/app/transports/sse"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
	"github.com/Southclaws/storyden/tests/robot"
)

func TestDefaultSearchesLoadsAndUsesCustomToolset(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{LanguageModelProvider: "mock"},
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
			settingsRepo *settings.SettingsRepository,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				instruction := "Use this only when explicitly asked to exercise failure handling."
				toolsetName := "Failure diagnostics " + xid.New().String()
				toolset := tests.AssertRequest(cl.RobotToolsetCreateWithResponse(root,
					openapi.RobotToolsetCreateJSONRequestBody{
						Name:        toolsetName,
						Description: "Exercises a deliberately failing diagnostic tool",
						Instruction: &instruction,
						Tools:       openapi.RobotToolNameList{"throw_an_error"},
					},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, toolset.JSON200)
				toolsetID := toolset.JSON200.Id

				scriptName := "robot-chat-toolset-load-" + xid.New().String() + ".yaml"
				scriptPath := filepath.Join("..", "scripts", scriptName)
				writeScript(t, scriptPath, `steps:
  - match:
      tool_result: throw_an_error
    respond:
      text: "Loaded Toolset executed."
      finish: "stop"
  - match:
      tool_result: toolset_load
    respond:
      tool_calls:
        - id: call_loaded_tool_1
          name: throw_an_error
          args: {}
  - match:
      tool_result: toolset_get
    respond:
      tool_calls:
        - id: call_toolset_load_1
          name: toolset_load
          args:
            toolsets:
              - "`+toolsetID+`"
  - match:
      tool_result: toolset_search
    respond:
      tool_calls:
        - id: call_toolset_get_1
          name: toolset_get
          args:
            id: "`+toolsetID+`"
  - match:
      contains: "exercise failure diagnostics"
    respond:
      tool_calls:
        - id: call_toolset_search_1
          name: toolset_search
          args:
            query: "failure diagnostics"
`)
				defer os.Remove(scriptPath)
				require.NoError(t, robot.SetRobotSettings(root, settingsRepo, "mock/../scripts/"+scriptName))

				stream := doChat(t, root, ts, adminSession, xid.New().String(), "", "exercise failure diagnostics")
				calls := collectToolCalls(stream)
				assert.Equal(t, []string{"toolset_search", "toolset_get", "toolset_load", "throw_an_error"}, calls)
				assert.Equal(t, "Loaded Toolset executed.", strings.Join(collectTextDeltas(stream), ""))

				outputs := collectToolOutputs(stream)
				searchOutput := toolOutputByCallID(t, outputs, "call_toolset_search_1")
				searchResults, ok := searchOutput["toolsets"].([]any)
				require.True(t, ok)
				require.NotEmpty(t, searchResults)
				searchResult, ok := searchResults[0].(map[string]any)
				require.True(t, ok)
				assert.ElementsMatch(t, []string{"id", "name", "description"}, mapKeys(searchResult))
				assert.NotContains(t, searchResult, "source")
				assert.NotContains(t, searchResult, "tools")

				getOutput := toolOutputByCallID(t, outputs, "call_toolset_get_1")
				assert.Equal(t, toolsetID, getOutput["id"])
				assert.Equal(t, toolsetName, getOutput["name"])
				assert.Equal(t, "Exercises a deliberately failing diagnostic tool", getOutput["description"])
				assert.Equal(t, instruction, getOutput["instruction"])
				assert.Equal(t, []any{"throw_an_error"}, getOutput["tools"])
				assert.Equal(t, true, getOutput["editable"])
				assert.Equal(t, false, getOutput["requires_workspace"])
				assert.NotContains(t, getOutput, "source")

				var sawExpectedError bool
				for _, output := range outputs {
					payload, ok := output.Output.(map[string]any)
					if !ok {
						continue
					}
					message, _ := payload["error"].(string)
					if strings.Contains(message, "intentional tool error") {
						sawExpectedError = true
					}
				}
				assert.True(t, sawExpectedError)
			}))
		}),
	)
}

func TestDenbotDoesNotLoadWorkspaceToolsetWithoutWorkspace(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{LanguageModelProvider: "mock"},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelAck),
		sse.Build(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
			settingsRepo *settings.SettingsRepository,
			db *ent.Client,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				scriptName := "robot-chat-workspace-toolset-" + xid.New().String() + ".yaml"
				scriptPath := filepath.Join("..", "scripts", scriptName)
				writeScript(t, scriptPath, `steps:
  - match:
      tool_result: toolset_load
    respond:
      text: "Workspace requirement explained."
      finish: "stop"
  - match:
      tool_result: toolset_get
    respond:
      tool_calls:
        - id: call_workspace_toolset_load
          name: toolset_load
          args:
            toolsets:
              - system.plugin_studio
  - match:
      contains: "inspect plugin studio"
    respond:
      tool_calls:
        - id: call_workspace_toolset_get
          name: toolset_get
          args:
            id: system.plugin_studio
`)
				defer os.Remove(scriptPath)
				require.NoError(t, robot.SetRobotSettings(root, settingsRepo, "mock/../scripts/"+scriptName))

				sessionID := xid.New()
				stream := doChat(t, root, ts, adminSession, sessionID.String(), "", "inspect plugin studio")
				assert.Equal(t, []string{"toolset_get", "toolset_load"}, collectToolCalls(stream))
				assert.Equal(t, "Workspace requirement explained.", strings.Join(collectTextDeltas(stream), ""))

				outputs := collectToolOutputs(stream)
				getOutput := toolOutputByCallID(t, outputs, "call_workspace_toolset_get")
				assert.Equal(t, true, getOutput["requires_workspace"])

				loadOutput := toolOutputByCallID(t, outputs, "call_workspace_toolset_load")
				assert.Equal(t,
					`Toolset "Plugin studio" (system.plugin_studio) requires an active Robot workspace. Start or continue this conversation with a workspace, then try again.`,
					loadOutput["error"],
				)

				storedSession, err := db.RobotSession.Get(root, sessionID)
				require.NoError(t, err)
				assert.NotContains(t, storedSession.State, "active_toolsets")
			}))
		}),
	)
}

func toolOutputByCallID(t *testing.T, outputs []openapi.ToolOutputAvailablePart, callID string) map[string]any {
	t.Helper()
	for _, output := range outputs {
		if output.ToolCallId != callID {
			continue
		}
		payload, ok := output.Output.(map[string]any)
		require.True(t, ok)
		return payload
	}
	require.FailNow(t, "tool output not found", "call ID: %s", callID)
	return nil
}

func mapKeys(values map[string]any) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}

func TestDefaultSearchesInspectsLoadsAndUsesIndividualTool(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{LanguageModelProvider: "mock"},
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
			settingsRepo *settings.SettingsRepository,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				scriptName := "robot-chat-tool-load-" + xid.New().String() + ".yaml"
				scriptPath := filepath.Join("..", "scripts", scriptName)
				writeScript(t, scriptPath, `steps:
  - match:
      tool_result: throw_an_error
    respond:
      text: "Individual tool executed."
      finish: "stop"
  - match:
      tool_result: tool_load
    respond:
      tool_calls:
        - id: call_loaded_tool_1
          name: throw_an_error
          args: {}
  - match:
      tool_result: tool_get
    respond:
      tool_calls:
        - id: call_tool_load_1
          name: tool_load
          args:
            tools:
              - throw_an_error
  - match:
      tool_result: tool_search
    respond:
      tool_calls:
        - id: call_tool_get_1
          name: tool_get
          args:
            id: throw_an_error
  - match:
      contains: "exercise one error tool"
    respond:
      tool_calls:
        - id: call_tool_search_1
          name: tool_search
          args:
            query: "intentional error handling"
`)
				defer os.Remove(scriptPath)
				require.NoError(t, robot.SetRobotSettings(root, settingsRepo, "mock/../scripts/"+scriptName))

				stream := doChat(t, root, ts, adminSession, xid.New().String(), "", "exercise one error tool")
				assert.Equal(t, []string{"tool_search", "tool_get", "tool_load", "throw_an_error"}, collectToolCalls(stream))
				assert.Equal(t, "Individual tool executed.", strings.Join(collectTextDeltas(stream), ""))
			}))
		}),
	)
}
