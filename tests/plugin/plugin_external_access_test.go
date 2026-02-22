package plugin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
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
	sdk "github.com/Southclaws/storyden/sdk/go/storyden"
	"github.com/Southclaws/storyden/tests"
)

func TestExternalPluginAccessKeyAndClientBuilder(t *testing.T) {
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

			botHandle := "plg-" + xid.New().String()
			manifest := openapi.PluginManifest(map[string]any{
				"id":              "test-access-" + xid.New().String(),
				"name":            "External Access Test",
				"author":          "test-author",
				"description":     "External plugin access test",
				"version":         "1.0.0",
				"command":         "./plugin",
				"events_consumed": []string{},
				"access": map[string]any{
					"handle":      botHandle,
					"name":        "External Access Bot",
					"permissions": []string{"account_get"},
				},
			})

			body := openapi.PluginInitialProps{}
			r.NoError(body.FromPluginInitialExternal(openapi.PluginInitialExternal{
				Mode:     openapi.External,
				Manifest: manifest,
			}))

			addResp := tests.AssertRequest(
				cl.PluginAddWithResponse(root, body, adminSession),
			)(t, http.StatusOK)
			plugin := openapi.Plugin(*addResp.JSON200)
			pluginXID, err := xid.FromString(string(plugin.Id))
			r.NoError(err)
			installationID := resource_plugin.InstallationID(pluginXID)

			getResp := tests.AssertRequest(
				cl.PluginGetWithResponse(root, string(plugin.Id), adminSession),
			)(t, http.StatusOK)
			ext, err := getResp.JSON200.Connection.AsPluginExternalProps()
			r.NoError(err)
			token := ext.Token

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

			access, err := pl.GetAccess(root)
			r.NoError(err)
			r.True(strings.HasPrefix(access.AccessKey, "sdbak_"))

			manualAccount := tests.AssertRequest(
				cl.AccountGetWithResponse(root, func(_ context.Context, req *http.Request) error {
					req.Header.Set("Authorization", "Bearer "+access.AccessKey)
					return nil
				}),
			)(t, http.StatusOK)
			r.Equal(botHandle, string(manualAccount.JSON200.Handle))

			builtClient, err := pl.BuildAPIClient(root)
			r.NoError(err)
			rawClient, ok := builtClient.ClientInterface.(*openapi.Client)
			r.True(ok, "expected *openapi.Client")
			rawClient.Server = strings.TrimRight(ts.URL, "/") + "/api/"

			builtAccount := tests.AssertRequest(
				builtClient.AccountGetWithResponse(root),
			)(t, http.StatusOK)
			r.Equal(botHandle, string(builtAccount.JSON200.Handle))
		}))
	}))
}
