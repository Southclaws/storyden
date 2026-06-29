package plugin_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	resource_plugin "github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	rpc_transport "github.com/Southclaws/storyden/app/transports/rpc"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
	sdk "github.com/Southclaws/storyden/sdk/go/storyden"
	"github.com/Southclaws/storyden/tests"
)

func TestExternalPluginRobotRunCompletedStructuredOutput(t *testing.T) {
	withRobotRunPlugin(t, []string{"USE_ROBOTS"}, func(
		root context.Context,
		cl *openapi.ClientWithResponses,
		adminSession openapi.RequestEditorFn,
		pl *sdk.Plugin,
		r *require.Assertions,
	) {
		scriptName := writeRobotRunScript(t, `steps:
  - match:
      contains: "summarise"
    respond:
      tool_calls:
        - id: call_finish_1
          name: robot_run_finish
          args:
            status: "completed"
            summary: "Robot run complete."
      finish: "stop"
`)
		robotID := createRobotForRun(t, root, cl, adminSession, r, scriptName)

		result := runRobotRPC(t, root, pl, robotID, "summarise this")

		output, ok := result.Output.Get()
		r.True(ok, "expected structured output")
		r.Equal(rpc.RobotRunStatusCompleted, output.Status)
		r.Equal("Robot run complete.", output.Summary)
		_, ok = result.SessionID.Get()
		r.True(ok, "expected session_id")
		r.False(result.Error.Ok(), "unexpected error: %v", result.Error)
	})
}

func TestExternalPluginRobotRunBypassesConfirmationForAssignedTool(t *testing.T) {
	withRobotRunPlugin(t, []string{"USE_ROBOTS", "MANAGE_ROBOTS"}, func(
		root context.Context,
		cl *openapi.ClientWithResponses,
		adminSession openapi.RequestEditorFn,
		pl *sdk.Plugin,
		r *require.Assertions,
	) {
		victim := tests.AssertRequest(cl.RobotCreateWithResponse(root,
			openapi.RobotCreateJSONRequestBody{
				Name:        "robot-run-delete-victim-" + xid.New().String(),
				Description: "victim for unattended delete test",
				Playbook:    "you are a victim",
			},
			adminSession,
		))(t, http.StatusOK)
		r.NotNil(victim.JSON200)
		victimID := string(victim.JSON200.Id)

		scriptName := writeRobotRunScript(t, fmt.Sprintf(`steps:
  - match:
      contains: "delete"
    respond:
      tool_calls:
        - id: call_delete_1
          name: robot_delete
          args:
            id: "%s"
  - match:
      tool_result: robot_delete
    respond:
      tool_calls:
        - id: call_finish_1
          name: robot_run_finish
          args:
            status: "completed"
            summary: "Robot deleted without live confirmation."
      finish: "stop"
`, victimID))
		robotID := createRobotForRun(t, root, cl, adminSession, r, scriptName, "robot_delete")

		result := runRobotRPC(t, root, pl, robotID, "delete it")

		output, ok := result.Output.Get()
		r.True(ok, "expected structured output")
		r.Equal(rpc.RobotRunStatusCompleted, output.Status)
		r.Equal("Robot deleted without live confirmation.", output.Summary)
		tests.AssertRequest(cl.RobotGetWithResponse(root, openapi.RobotIDParam(victimID), adminSession))(t, http.StatusNotFound)
	})
}

func TestExternalPluginRobotRunBlocksElicitation(t *testing.T) {
	withRobotRunPlugin(t, []string{"USE_ROBOTS", "READ_PUBLISHED_LIBRARY"}, func(
		root context.Context,
		cl *openapi.ClientWithResponses,
		adminSession openapi.RequestEditorFn,
		pl *sdk.Plugin,
		r *require.Assertions,
		sessionRepo *robot_session.Repository,
	) {
		scriptName := writeRobotRunScript(t, `steps:
  - match:
      contains: "choose"
    respond:
      tool_calls:
        - id: call_request_page_1
          name: library_request_page
          args: {}
  - match:
      tool_result: library_request_page
    respond:
      tool_calls:
        - id: call_finish_1
          name: robot_run_finish
          args:
            status: "blocked"
            summary: "A Library page selection is required."
            attention:
              reason: "missing_input"
              message: "Select a Library page before running unattended."
      finish: "stop"
`)
		robotID := createRobotForRun(t, root, cl, adminSession, r, scriptName, "library_request_page")

		result := runRobotRPC(t, root, pl, robotID, "choose a page")

		output, ok := result.Output.Get()
		r.True(ok, "expected structured output")
		r.Equal(rpc.RobotRunStatusBlocked, output.Status)
		attention, ok := output.Attention.Get()
		r.True(ok, "expected attention")
		r.Equal(rpc.RobotRunAttentionReasonMissingInput, attention.Reason)

		sessionID, ok := result.SessionID.Get()
		r.True(ok, "expected session_id")
		sess, _, err := sessionRepo.Get(root, robot.SessionID(sessionID), robot.NewMessageCursorParams(opt.NewEmpty[robot.MessageID](), 20))
		r.NoError(err)
		r.NotContains(sess.State, "pending_client_tool_ids")
		r.NotContains(sess.State, "pending_client_tool_robots")
	})
}

func TestExternalPluginRobotRunMissingToolPermissionFails(t *testing.T) {
	withRobotRunPlugin(t, []string{"USE_ROBOTS"}, func(
		root context.Context,
		cl *openapi.ClientWithResponses,
		adminSession openapi.RequestEditorFn,
		pl *sdk.Plugin,
		r *require.Assertions,
	) {
		victim := tests.AssertRequest(cl.RobotCreateWithResponse(root,
			openapi.RobotCreateJSONRequestBody{
				Name:        "robot-run-rbac-victim-" + xid.New().String(),
				Description: "victim for unattended rbac test",
				Playbook:    "you are a victim",
			},
			adminSession,
		))(t, http.StatusOK)
		r.NotNil(victim.JSON200)
		victimID := string(victim.JSON200.Id)

		scriptName := writeRobotRunScript(t, fmt.Sprintf(`steps:
  - match:
      contains: "delete"
    respond:
      tool_calls:
        - id: call_delete_1
          name: robot_delete
          args:
            id: "%s"
`, victimID))
		robotID := createRobotForRun(t, root, cl, adminSession, r, scriptName, "robot_delete")

		result := runRobotRPC(t, root, pl, robotID, "delete it")

		r.True(result.Error.Ok(), "expected error")
		output, ok := result.Output.Get()
		r.True(ok, "expected failed output")
		r.Equal(rpc.RobotRunStatusFailed, output.Status)
		tests.AssertRequest(cl.RobotGetWithResponse(root, openapi.RobotIDParam(victimID), adminSession))(t, http.StatusOK)
	})
}

func TestExternalPluginRobotRunMalformedOutputFails(t *testing.T) {
	withRobotRunPlugin(t, []string{"USE_ROBOTS"}, func(
		root context.Context,
		cl *openapi.ClientWithResponses,
		adminSession openapi.RequestEditorFn,
		pl *sdk.Plugin,
		r *require.Assertions,
	) {
		scriptName := writeRobotRunScript(t, `steps:
  - match:
      contains: "bad"
    respond:
      text: "not json"
      finish: "stop"
`)
		robotID := createRobotForRun(t, root, cl, adminSession, r, scriptName)

		result := runRobotRPC(t, root, pl, robotID, "bad output please")

		r.True(result.Error.Ok(), "expected error")
		output, ok := result.Output.Get()
		r.True(ok, "expected failed output")
		r.Equal(rpc.RobotRunStatusFailed, output.Status)
	})
}

func withRobotRunPlugin(t *testing.T, permissions []string, run any) {
	t.Helper()

	integration.Test(t, &config.Config{
		LanguageModelProvider:       "mock",
		PluginMaxRestartAttempts:    1,
		PluginMaxBackoffDuration:    100 * time.Millisecond,
		PluginRuntimeCrashThreshold: 1 * time.Second,
		PluginRuntimeCrashBackoff:   100 * time.Millisecond,
	}, e2e.Setup(), rpc_transport.Build(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		accountWrite *account_writer.Writer,
		runner plugin_runner.Host,
		ts *httptest.Server,
		sessionRepo *robot_session.Repository,
		settingsRepo *settings.SettingsRepository,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			_, err := settingsRepo.Set(root, settings.Settings{
				Services: opt.New(settings.ServiceSettings{
					Robots: opt.New(settings.RobotServiceSettings{
						Enabled:      opt.New(true),
						DefaultModel: opt.New("mock/../robot/scripts/robot-chat-ack.yaml"),
						Providers: opt.New(map[string]settings.RobotProviderSettings{
							"mock": {
								Enabled: opt.New(true),
							},
						}),
					}),
				}),
			})
			r.NoError(err)

			adminCtx, _ := e2e.WithAccount(root, accountWrite, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			manifest := robotRunManifest("External Robot Run "+xid.New().String(), permissions)
			_, installationID, token := addExternalPluginWithManifest(t, root, cl, adminSession, r, manifest)
			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateConnecting)

			rpcURL := externalRPCURL(t, ts.URL, token)
			t.Setenv("STORYDEN_RPC_URL", rpcURL)

			pl, err := sdk.New(root)
			r.NoError(err)

			runCtx, cancel := context.WithCancel(root)
			done := make(chan error, 1)
			go func() {
				done <- pl.Run(runCtx)
			}()
			defer func() {
				_ = pl.Shutdown()
				cancel()
				waitForPluginStop(done)
			}()

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateActive)

			switch f := run.(type) {
			case func(context.Context, *openapi.ClientWithResponses, openapi.RequestEditorFn, *sdk.Plugin, *require.Assertions):
				f(root, cl, adminSession, pl, r)
			case func(context.Context, *openapi.ClientWithResponses, openapi.RequestEditorFn, *sdk.Plugin, *require.Assertions, *robot_session.Repository):
				f(root, cl, adminSession, pl, r, sessionRepo)
			default:
				r.Failf("unsupported robot run test function", "%T", run)
			}
		}))
	}))
}

func robotRunManifest(name string, permissions []string) openapi.PluginManifest {
	return openapi.PluginManifest(map[string]any{
		"id":              "robot-run-" + xid.New().String(),
		"name":            name,
		"author":          "test-author",
		"description":     "Robot run test plugin",
		"version":         "1.0.0",
		"command":         "./plugin",
		"events_consumed": []string{},
		"access": map[string]any{
			"handle":      "robot-run-" + xid.New().String(),
			"name":        "Robot Run Bot",
			"permissions": permissions,
		},
	})
}

func createRobotForRun(
	t *testing.T,
	ctx context.Context,
	cl *openapi.ClientWithResponses,
	adminSession openapi.RequestEditorFn,
	r *require.Assertions,
	scriptName string,
	tools ...string,
) xid.ID {
	var toolList *openapi.RobotToolNameList
	if len(tools) > 0 {
		v := openapi.RobotToolNameList(tools)
		toolList = &v
	}

	created := tests.AssertRequest(cl.RobotCreateWithResponse(ctx,
		openapi.RobotCreateJSONRequestBody{
			Name:        "robot-run-actor-" + xid.New().String(),
			Description: "robot_run test actor",
			Playbook:    "you are a robot_run test actor",
			Model:       ptr("mock/../robot/scripts/" + scriptName),
			Tools:       toolList,
		},
		adminSession,
	))(t, http.StatusOK)
	r.NotNil(created.JSON200)
	id, err := xid.FromString(string(created.JSON200.Id))
	r.NoError(err)
	return id
}

func runRobotRPC(t *testing.T, ctx context.Context, pl *sdk.Plugin, robotID xid.ID, message string) *rpc.RPCResponseRobotRun {
	t.Helper()

	result, err := pl.Send(ctx, rpc.RPCRequestRobotRun{
		Jsonrpc: "2.0",
		Method:  "robot_run",
		Params: rpc.RPCRequestRobotRunParams{
			Message: message,
			RobotID: robotID.String(),
		},
	})
	require.NoError(t, err)

	typed, ok := result.(*rpc.RPCResponseRobotRun)
	require.True(t, ok, "expected *rpc.RPCResponseRobotRun, got %T", result)
	return typed
}

func writeRobotRunScript(t *testing.T, content string) string {
	t.Helper()

	name := "robot-run-" + xid.New().String() + ".yaml"
	path := filepath.Join("..", "robot", "scripts", name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	t.Cleanup(func() {
		_ = os.Remove(path)
	})
	return name
}

func ptr[T any](v T) *T {
	return &v
}
