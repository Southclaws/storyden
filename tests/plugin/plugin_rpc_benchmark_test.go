package plugin_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func TestExternalPluginConfigureRPCRoughBenchmark(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping plugin RPC benchmark in short mode")
	}

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
			manifest := buildConfigurationManifest("External Configure RPC Benchmark")
			pluginID, installationID, token := addExternalPluginWithManifest(t, root, cl, adminSession, r, manifest)

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateConnecting)

			stopPlugin, done := runExternalSDKConfigPlugin(root, t, externalRPCURL(t, ts.URL, token), func(_ map[string]any) error {
				return nil
			})
			defer stopPlugin()
			defer waitForPluginStop(done)

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateActive)

			const iterations = 200

			for i := 0; i < 10; i++ {
				tests.AssertRequest(
					cl.PluginUpdateConfigurationWithResponse(root, pluginID, openapi.PluginUpdateConfigurationJSONRequestBody{
						"name":      "warmup",
						"enabled":   true,
						"threshold": float64(i),
					}, adminSession),
				)(t, http.StatusOK)
			}

			start := time.Now()
			for i := 0; i < iterations; i++ {
				tests.AssertRequest(
					cl.PluginUpdateConfigurationWithResponse(root, pluginID, openapi.PluginUpdateConfigurationJSONRequestBody{
						"name":      fmt.Sprintf("bench-%d", i),
						"enabled":   true,
						"threshold": float64(i),
					}, adminSession),
				)(t, http.StatusOK)
			}
			elapsed := time.Since(start)
			t.Logf("plugin configure RPC rough benchmark: total=%s iterations=%d avg=%s", elapsed, iterations, elapsed/time.Duration(iterations))
		}))
	}))
}
