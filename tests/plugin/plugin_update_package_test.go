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

func TestSupervisedPluginUpdatePackageRunningRestartsProcess(t *testing.T) {
	t.Parallel()

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

			setRunnerServerURLForTests(runner, ts.URL)

			adminSession, _ := createAdminSession(root, cl, sh, accountWrite, r)
			pluginID, installationID := addSupervisedPluginArchive(
				t,
				root,
				cl,
				adminSession,
				"test_data/supervised_config",
				"supervised_config",
			)

			tests.AssertRequest(
				cl.PluginSetActiveStateWithResponse(root, pluginID, openapi.PluginSetActiveStateJSONRequestBody{
					Active: openapi.PluginActiveStateActive,
				}, adminSession),
			)(t, http.StatusOK)
			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateActive)

			sessBefore, err := runner.GetSession(root, installationID)
			r.NoError(err)
			startedBefore := sessBefore.GetStartedAt().OrZero()
			r.False(startedBefore.IsZero())

			time.Sleep(50 * time.Millisecond)

			updateResp := tests.AssertRequest(
				updateSupervisedPluginArchive(
					t,
					root,
					cl,
					adminSession,
					pluginID,
					"test_data/supervised_config_v2",
					"supervised_config_v2",
				),
			)(t, http.StatusOK)
			r.NotNil(updateResp.JSON200)
			r.NotNil(updateResp.JSON200.Version)
			r.Equal("1.0.1", *updateResp.JSON200.Version)

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateActive)

			sessAfter, err := runner.GetSession(root, installationID)
			r.NoError(err)
			startedAfter := sessAfter.GetStartedAt().OrZero()
			r.False(startedAfter.IsZero())
			r.True(startedAfter.After(startedBefore))
		}))
	}))
}

func TestSupervisedPluginUpdatePackageInactiveKeepsInactive(t *testing.T) {
	t.Parallel()

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

			setRunnerServerURLForTests(runner, ts.URL)

			adminSession, _ := createAdminSession(root, cl, sh, accountWrite, r)
			pluginID, installationID := addSupervisedPluginArchive(
				t,
				root,
				cl,
				adminSession,
				"test_data/supervised_config",
				"supervised_config",
			)

			tests.AssertRequest(
				cl.PluginSetActiveStateWithResponse(root, pluginID, openapi.PluginSetActiveStateJSONRequestBody{
					Active: openapi.PluginActiveStateActive,
				}, adminSession),
			)(t, http.StatusOK)
			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateActive)

			tests.AssertRequest(
				cl.PluginSetActiveStateWithResponse(root, pluginID, openapi.PluginSetActiveStateJSONRequestBody{
					Active: openapi.PluginActiveStateInactive,
				}, adminSession),
			)(t, http.StatusOK)
			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateInactive)

			updateResp := tests.AssertRequest(
				updateSupervisedPluginArchive(
					t,
					root,
					cl,
					adminSession,
					pluginID,
					"test_data/supervised_config_v2",
					"supervised_config_v2",
				),
			)(t, http.StatusOK)
			r.NotNil(updateResp.JSON200)
			r.NotNil(updateResp.JSON200.Version)
			r.Equal("1.0.1", *updateResp.JSON200.Version)

			status, err := updateResp.JSON200.Status.AsPluginStatusInactive()
			r.NoError(err)
			r.Equal(openapi.Inactive, status.ActiveState)
			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateInactive)

			tests.AssertRequest(
				cl.PluginSetActiveStateWithResponse(root, pluginID, openapi.PluginSetActiveStateJSONRequestBody{
					Active: openapi.PluginActiveStateActive,
				}, adminSession),
			)(t, http.StatusOK)
			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateActive)
		}))
	}))
}

func TestSupervisedPluginUpdatePackageRejectsManifestIDMismatch(t *testing.T) {
	t.Parallel()

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

			setRunnerServerURLForTests(runner, ts.URL)

			adminSession, _ := createAdminSession(root, cl, sh, accountWrite, r)
			pluginID, _ := addSupervisedPluginArchive(
				t,
				root,
				cl,
				adminSession,
				"test_data/supervised_config",
				"supervised_config",
			)

			tests.AssertRequest(
				updateSupervisedPluginArchive(
					t,
					root,
					cl,
					adminSession,
					pluginID,
					"test_data/event_listener",
					"event_listener",
				),
			)(t, http.StatusBadRequest)

			getResp := tests.AssertRequest(
				cl.PluginGetWithResponse(root, pluginID, adminSession),
			)(t, http.StatusOK)
			r.NotNil(getResp.JSON200.Version)
			r.Equal("1.0.0", *getResp.JSON200.Version)
		}))
	}))
}

func setRunnerServerURLForTests(runner plugin_runner.Host, serverURL string) {
	type serverURLSetter interface {
		SetServerURL(string)
	}

	if setter, ok := runner.(serverURLSetter); ok {
		setter.SetServerURL(serverURL)
	}
}

func addSupervisedPluginArchive(
	t *testing.T,
	ctx context.Context,
	cl *openapi.ClientWithResponses,
	adminSession openapi.RequestEditorFn,
	fixtureDir string,
	pluginName string,
) (string, resource_plugin.InstallationID) {
	t.Helper()

	archivePath := packageTestPluginArchive(t, fixtureDir, pluginName)
	pluginFile, err := os.Open(archivePath)
	require.NoError(t, err)
	defer pluginFile.Close()

	addResp := tests.AssertRequest(
		cl.PluginAddWithBodyWithResponse(ctx, "application/zip", pluginFile, adminSession),
	)(t, http.StatusOK)

	pluginID := string(addResp.JSON200.Id)
	pluginXID, err := xid.FromString(pluginID)
	require.NoError(t, err)

	return pluginID, resource_plugin.InstallationID(pluginXID)
}

func updateSupervisedPluginArchive(
	t *testing.T,
	ctx context.Context,
	cl *openapi.ClientWithResponses,
	adminSession openapi.RequestEditorFn,
	pluginID string,
	fixtureDir string,
	pluginName string,
) (*openapi.PluginUpdatePackageResponse, error) {
	t.Helper()

	archivePath := packageTestPluginArchive(t, fixtureDir, pluginName)
	pluginFile, err := os.Open(archivePath)
	require.NoError(t, err)
	defer pluginFile.Close()

	return cl.PluginUpdatePackageWithBodyWithResponse(
		ctx,
		pluginID,
		"application/zip",
		pluginFile,
		adminSession,
	)
}
