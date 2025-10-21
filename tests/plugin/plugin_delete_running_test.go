package plugin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	resource_plugin "github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	rpc_transport "github.com/Southclaws/storyden/app/transports/rpc"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestDeleteRunningExternalPluginDisconnectsSession(t *testing.T) {
	integration.Test(t, &config.Config{
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
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)

			adminSession, _ := createAdminSession(root, cl, sh, accountWrite, r)
			pluginID, installationID, token := addExternalPlugin(
				t,
				root,
				cl,
				adminSession,
				r,
				"External Delete Running",
				[]string{},
			)

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateConnecting)

			rpcURL := externalRPCURL(t, ts.URL, token)
			stopPlugin, pluginDone := runExternalSDKPlugin(root, t, rpcURL, nil)
			defer stopPlugin()

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateActive)

			tests.AssertRequest(
				cl.PluginDeleteWithResponse(root, pluginID, adminSession),
			)(t, http.StatusNoContent)

			select {
			case <-pluginDone:
			case <-time.After(2 * time.Second):
				t.Fatal("external plugin should be disconnected when deleted")
			}

			tests.AssertRequest(
				cl.PluginGetWithResponse(root, pluginID, adminSession),
			)(t, http.StatusNotFound)

			require.Eventually(t, func() bool {
				_, err := runner.GetSession(root, installationID)
				return err != nil
			}, 2*time.Second, 50*time.Millisecond)
		}))
	}))
}

func TestDeleteInactiveSupervisedPluginWithoutLoadedSession(t *testing.T) {
	integration.Test(t, &config.Config{
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
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)

			adminSession, _ := createAdminSession(root, cl, sh, accountWrite, r)

			archivePath := packageTestPluginArchive(t, "test_data/example", "example-plugin")
			pluginFile, err := os.Open(archivePath)
			r.NoError(err)
			defer pluginFile.Close()

			addResp := tests.AssertRequest(
				cl.PluginAddWithBodyWithResponse(root, "application/zip", pluginFile, adminSession),
			)(t, http.StatusOK)

			pluginID := string(addResp.JSON200.Id)
			pluginXID, err := xid.FromString(pluginID)
			r.NoError(err)
			installationID := resource_plugin.InstallationID(pluginXID)

			_, err = runner.GetSession(root, installationID)
			r.Error(err)

			tests.AssertRequest(
				cl.PluginDeleteWithResponse(root, pluginID, adminSession),
			)(t, http.StatusNoContent)
		}))
	}))
}
