package chat_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/semdex/robot/mcpclient"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
	robottest "github.com/Southclaws/storyden/tests/robot"
)

func TestMCPServerToolsetCanBeDiscoveredAndAssigned(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{LanguageModelProvider: "mock"},
		e2e.Setup(),
		robottest.WithRobotSettings(mockModelAck),
		fx.Replace(mcpclient.HTTPClient{Client: http.DefaultClient}),
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
				endpoint := newToolsetMCPServer(t)
				slug := "calendar-" + xid.New().String()
				name := "Calendar MCP " + xid.New().String()
				enabled := true

				created := tests.AssertRequest(cl.RobotMCPServerCreateWithResponse(root,
					openapi.RobotMCPServerCreateJSONRequestBody{
						Name:        name,
						Slug:        &slug,
						Description: stringPointer("Schedule appointments and echo event details."),
						EndpointUrl: endpoint,
						Enabled:     &enabled,
					},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, created.JSON200)

				toolsetID := created.JSON200.ToolsetId
				callableName := mcpclient.CallableName(name, "echo")
				denbotScriptName := "robot-chat-mcp-toolset-" + xid.New().String() + ".yaml"
				denbotScriptPath := filepath.Join("..", "scripts", denbotScriptName)
				writeScript(t, denbotScriptPath, `steps:
  - match:
      tool_result: `+callableName+`
    respond:
      text: "Denbot used the MCP Toolset."
      finish: "stop"
  - match:
      tool_result: toolset_load
    respond:
      tool_calls:
        - id: call_mcp_echo
          name: `+callableName+`
          args:
            message: "meeting"
  - match:
      tool_result: toolset_get
    respond:
      tool_calls:
        - id: call_mcp_toolset_load
          name: toolset_load
          args:
            toolsets:
              - "`+toolsetID+`"
  - match:
      tool_result: toolset_search
    respond:
      tool_calls:
        - id: call_mcp_toolset_get
          name: toolset_get
          args:
            id: "`+toolsetID+`"
  - match:
      contains: "schedule a meeting"
    respond:
      tool_calls:
        - id: call_mcp_toolset_search
          name: toolset_search
          args:
            query: "schedule appointments"
`)
				defer os.Remove(denbotScriptPath)
				require.NoError(t, robottest.SetRobotSettings(root, settingsRepo, "mock/../scripts/"+denbotScriptName))

				denbotStream := doChat(t, root, ts, adminSession, xid.New().String(), "", "schedule a meeting")
				assert.Equal(t, []string{"toolset_search", "toolset_get", "toolset_load", callableName}, collectToolCalls(denbotStream))
				assert.Equal(t, "Denbot used the MCP Toolset.", strings.Join(collectTextDeltas(denbotStream), ""))

				customScriptName := "robot-chat-custom-mcp-toolset-" + xid.New().String() + ".yaml"
				customScriptPath := filepath.Join("..", "scripts", customScriptName)
				writeScript(t, customScriptPath, `steps:
  - match:
      tool_result: `+callableName+`
    respond:
      text: "Custom Robot used the MCP Toolset."
      finish: "stop"
  - match:
      contains: "use the assigned server"
    respond:
      tool_calls:
        - id: call_custom_mcp_echo
          name: `+callableName+`
          args:
            message: "assigned"
`)
				defer os.Remove(customScriptPath)

				customRobot := tests.AssertRequest(cl.RobotCreateWithResponse(root,
					openapi.RobotCreateJSONRequestBody{
						Name:        "Assigned MCP Robot " + xid.New().String(),
						Description: "Uses one connected MCP server.",
						Playbook:    "Use the assigned server when asked.",
						Model:       robotModelPtr("mock/../scripts/" + customScriptName),
						Toolsets:    robotToolsetsPtr(toolsetID),
					},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, customRobot.JSON200)

				customStream := doChat(t, root, ts, adminSession, xid.New().String(), string(customRobot.JSON200.Id), "use the assigned server")
				assert.Equal(t, []string{callableName}, collectToolCalls(customStream))
				assert.Equal(t, "Custom Robot used the MCP Toolset.", strings.Join(collectTextDeltas(customStream), ""))

				tests.AssertRequest(cl.RobotUpdateWithResponse(root,
					customRobot.JSON200.Id,
					openapi.RobotUpdateJSONRequestBody{Model: robotModelPtr(mockModelAck)},
					adminSession,
				))(t, http.StatusOK)
				disabled := false
				tests.AssertRequest(cl.RobotMCPServerUpdateWithResponse(root,
					created.JSON200.Id,
					openapi.RobotMCPServerUpdateJSONRequestBody{Enabled: &disabled},
					adminSession,
				))(t, http.StatusOK)

				degradedStream := doChat(t, root, ts, adminSession, xid.New().String(), string(customRobot.JSON200.Id), "continue safely")
				assert.Empty(t, collectErrorParts(degradedStream))
				assert.Equal(t, "ack", strings.Join(collectTextDeltas(degradedStream), ""))

				enabled = true
				tests.AssertRequest(cl.RobotMCPServerUpdateWithResponse(root,
					created.JSON200.Id,
					openapi.RobotMCPServerUpdateJSONRequestBody{Enabled: &enabled},
					adminSession,
				))(t, http.StatusOK)
				tests.AssertRequest(cl.RobotUpdateWithResponse(root,
					customRobot.JSON200.Id,
					openapi.RobotUpdateJSONRequestBody{Model: robotModelPtr("mock/../scripts/" + customScriptName)},
					adminSession,
				))(t, http.StatusOK)
				restoredStream := doChat(t, root, ts, adminSession, xid.New().String(), string(customRobot.JSON200.Id), "use the assigned server")
				assert.Equal(t, []string{callableName}, collectToolCalls(restoredStream))

				tests.AssertRequest(cl.RobotMCPServerDeleteWithResponse(root, created.JSON200.Id, adminSession))(t, http.StatusOK)
				tests.AssertRequest(cl.RobotUpdateWithResponse(root,
					customRobot.JSON200.Id,
					openapi.RobotUpdateJSONRequestBody{Model: robotModelPtr(mockModelAck)},
					adminSession,
				))(t, http.StatusOK)
				deletedStream := doChat(t, root, ts, adminSession, xid.New().String(), string(customRobot.JSON200.Id), "continue after deletion")
				assert.Empty(t, collectErrorParts(deletedStream))
				assert.Equal(t, "ack", strings.Join(collectTextDeltas(deletedStream), ""))

				persisted := tests.AssertRequest(cl.RobotGetWithResponse(root, customRobot.JSON200.Id, adminSession))(t, http.StatusOK)
				require.NotNil(t, persisted.JSON200)
				assert.Equal(t, openapi.RobotToolsetRefList{toolsetID}, persisted.JSON200.Toolsets)
			}))
		}),
	)
}

func newToolsetMCPServer(t *testing.T) string {
	t.Helper()

	server := sdkmcp.NewServer(&sdkmcp.Implementation{Name: "toolset-mcp", Version: "v1"}, nil)
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "echo",
		Title:       "Echo event details",
		Description: "Echo calendar event details.",
	}, func(ctx context.Context, req *sdkmcp.CallToolRequest, args map[string]any) (*sdkmcp.CallToolResult, map[string]any, error) {
		return nil, map[string]any{"message": args["message"]}, nil
	})

	handler := sdkmcp.NewStreamableHTTPHandler(func(r *http.Request) *sdkmcp.Server {
		return server
	}, nil)
	httpServer := httptest.NewServer(handler)
	t.Cleanup(httpServer.Close)
	return httpServer.URL
}

func stringPointer(value string) *string {
	return &value
}
