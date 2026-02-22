package plugin_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	resource_plugin "github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/resources/plugin/plugin_writer"
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

func TestExternalPluginConfigurationLifecycle(t *testing.T) {
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

			manifest := buildConfigurationManifest("External Config Lifecycle")
			pluginID, installationID, token := addExternalPluginWithManifest(t, root, cl, adminSession, r, manifest)

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateConnecting)

			var (
				lastConfigure map[string]any
				configureMu   sync.Mutex
			)
			stopPlugin, done := runExternalSDKConfigPlugin(root, t, externalRPCURL(t, ts.URL, token), func(config map[string]any) error {
				configureMu.Lock()
				defer configureMu.Unlock()
				lastConfigure = config
				return nil
			})
			defer stopPlugin()
			defer waitForPluginStop(done)

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateActive)

			schemaResp := tests.AssertRequest(
				cl.PluginGetConfigurationSchemaWithResponse(root, pluginID, adminSession),
			)(t, http.StatusOK)

			fields := schemaResp.JSON200.Fields
			r.NotNil(fields)
			r.Len(*fields, 3)
			fieldTypes := map[string]string{}
			for _, field := range *fields {
				id, typ := parseConfigurationField(field)
				fieldTypes[id] = typ
			}
			r.Equal("string", fieldTypes["name"])
			r.Equal("boolean", fieldTypes["enabled"])
			r.Equal("number", fieldTypes["threshold"])

			updated := tests.AssertRequest(
				cl.PluginUpdateConfigurationWithResponse(root, pluginID, openapi.PluginUpdateConfigurationJSONRequestBody{
					"name":      "configured-plugin",
					"enabled":   true,
					"threshold": 2.5,
					"extra":     "allowed",
				}, adminSession),
			)(t, http.StatusOK)

			updatedCfg := openapi.PluginConfiguration(*updated.JSON200)
			r.Equal("configured-plugin", updatedCfg["name"])
			r.Equal(true, updatedCfg["enabled"])
			r.Equal(2.5, updatedCfg["threshold"])

			valueResp := tests.AssertRequest(
				cl.PluginGetConfigurationWithResponse(root, pluginID, adminSession),
			)(t, http.StatusOK)

			valueCfg := openapi.PluginConfiguration(*valueResp.JSON200)
			r.Equal("configured-plugin", valueCfg["name"])
			r.Equal(true, valueCfg["enabled"])
			r.Equal(2.5, valueCfg["threshold"])
			r.Equal("allowed", valueCfg["extra"])

			configureMu.Lock()
			r.Equal("configured-plugin", lastConfigure["name"])
			r.Equal(true, lastConfigure["enabled"])
			r.Equal(2.5, lastConfigure["threshold"])
			configureMu.Unlock()
		}))
	}))
}

func TestExternalPluginConfigurationValidationAndConnectedChecks(t *testing.T) {
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
			manifest := buildConfigurationManifest("External Config Validation")

			pluginIDDisconnected, installationIDDisconnected, _ := addExternalPluginWithManifest(t, root, cl, adminSession, r, manifest)
			requireSessionState(t, root, runner, installationIDDisconnected, resource_plugin.ReportedStateConnecting)

			tests.AssertRequest(
				cl.PluginUpdateConfigurationWithResponse(root, pluginIDDisconnected, openapi.PluginUpdateConfigurationJSONRequestBody{
					"name":      "valid",
					"enabled":   true,
					"threshold": 1.5,
				}, adminSession),
			)(t, http.StatusBadRequest)

			pluginIDConnected, installationIDConnected, token := addExternalPluginWithManifest(t, root, cl, adminSession, r, manifest)
			requireSessionState(t, root, runner, installationIDConnected, resource_plugin.ReportedStateConnecting)

			stopPlugin, done := runExternalSDKConfigPlugin(root, t, externalRPCURL(t, ts.URL, token), func(config map[string]any) error {
				return nil
			})
			defer stopPlugin()
			defer waitForPluginStop(done)

			requireSessionState(t, root, runner, installationIDConnected, resource_plugin.ReportedStateActive)

			tests.AssertRequest(
				cl.PluginUpdateConfigurationWithResponse(root, pluginIDConnected, openapi.PluginUpdateConfigurationJSONRequestBody{
					"name":      "valid",
					"enabled":   "not-bool",
					"threshold": 1.5,
				}, adminSession),
			)(t, http.StatusBadRequest)

			valueResp := tests.AssertRequest(
				cl.PluginGetConfigurationWithResponse(root, pluginIDConnected, adminSession),
			)(t, http.StatusOK)
			r.Empty(openapi.PluginConfiguration(*valueResp.JSON200))
		}))
	}))
}

func TestExternalPluginConfigurationRejectedByPluginDoesNotPersist(t *testing.T) {
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
			manifest := buildConfigurationManifest("External Config Reject")
			pluginID, installationID, token := addExternalPluginWithManifest(t, root, cl, adminSession, r, manifest)

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateConnecting)

			stopPlugin, done := runExternalSDKConfigPlugin(root, t, externalRPCURL(t, ts.URL, token), func(config map[string]any) error {
				reject, _ := config["reject"].(bool)
				if reject {
					return errors.New("rejected by test plugin")
				}
				return nil
			})
			defer stopPlugin()
			defer waitForPluginStop(done)

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateActive)

			tests.AssertRequest(
				cl.PluginUpdateConfigurationWithResponse(root, pluginID, openapi.PluginUpdateConfigurationJSONRequestBody{
					"name":      "valid",
					"enabled":   true,
					"threshold": 1.5,
					"reject":    true,
				}, adminSession),
			)(t, http.StatusBadRequest)

			valueResp := tests.AssertRequest(
				cl.PluginGetConfigurationWithResponse(root, pluginID, adminSession),
			)(t, http.StatusOK)
			r.Empty(openapi.PluginConfiguration(*valueResp.JSON200))
		}))
	}))
}

func TestExternalPluginConfigureHandlerCanCallAccessGet(t *testing.T) {
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
			manifest := buildConfigurationManifest("External Config Access In Configure")
			manifest["access"] = map[string]any{
				"handle":      "cfg-access-" + xid.New().String(),
				"name":        "Config Access Bot",
				"permissions": []string{"account_get"},
			}

			pluginID, installationID, token := addExternalPluginWithManifest(t, root, cl, adminSession, r, manifest)
			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateConnecting)

			rpcURL := externalRPCURL(t, ts.URL, token)
			t.Setenv("STORYDEN_RPC_URL", rpcURL)

			pl, err := sdk.New(root)
			r.NoError(err)

			accessSeen := make(chan string, 1)
			pl.OnConfigure(func(ctx context.Context, _ map[string]any) error {
				access, err := pl.GetAccess(ctx)
				if err != nil {
					return err
				}
				accessSeen <- access.AccessKey
				return nil
			})

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

			tests.AssertRequest(
				cl.PluginUpdateConfigurationWithResponse(root, pluginID, openapi.PluginUpdateConfigurationJSONRequestBody{
					"name":      "configured-plugin",
					"enabled":   true,
					"threshold": 2.5,
				}, adminSession),
			)(t, http.StatusOK)

			select {
			case accessKey := <-accessSeen:
				r.True(strings.HasPrefix(accessKey, "sdbak_"))
			case <-time.After(5 * time.Second):
				t.Fatal("configure handler did not receive access key")
			}
		}))
	}))
}

func TestExternalPluginGetConfigRPCReadsPersistedConfiguration(t *testing.T) {
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
		pluginWrite *plugin_writer.Writer,
		runner plugin_runner.Host,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)

			adminSession, _ := createAdminSession(root, cl, sh, accountWrite, r)
			manifest := buildConfigurationManifest("External Config Get RPC")
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

			_, err = pluginWrite.UpdateConfig(root, installationID, map[string]any{
				"name":      "rpc-config",
				"enabled":   true,
				"threshold": 3.14,
			})
			r.NoError(err)

			result, err := pl.Send(root, rpc.RPCRequestGetConfig{
				Jsonrpc: "2.0",
				Method:  "get_config",
			})
			r.NoError(err)

			getConfig, ok := result.(*rpc.RPCResponseGetConfig)
			r.True(ok, "expected *rpc.RPCResponseGetConfig, got %T", result)
			r.Equal("rpc-config", getConfig.Config["name"])
			r.Equal(true, getConfig.Config["enabled"])
			r.Equal(3.14, getConfig.Config["threshold"])
		}))
	}))
}

func addExternalPluginWithManifest(
	t *testing.T,
	ctx context.Context,
	cl *openapi.ClientWithResponses,
	adminSession openapi.RequestEditorFn,
	r *require.Assertions,
	manifest openapi.PluginManifest,
) (string, resource_plugin.InstallationID, string) {
	body := openapi.PluginInitialProps{}
	err := body.FromPluginInitialExternal(openapi.PluginInitialExternal{
		Mode:     openapi.External,
		Manifest: manifest,
	})
	r.NoError(err)

	addResp := tests.AssertRequest(
		cl.PluginAddWithResponse(ctx, body, adminSession),
	)(t, http.StatusOK)
	plugin := openapi.Plugin(*addResp.JSON200)

	pluginXID, err := xid.FromString(string(plugin.Id))
	r.NoError(err)

	getResp := tests.AssertRequest(
		cl.PluginGetWithResponse(ctx, string(plugin.Id), adminSession),
	)(t, http.StatusOK)
	ext, err := getResp.JSON200.Connection.AsPluginExternalProps()
	r.NoError(err)

	return string(plugin.Id), resource_plugin.InstallationID(pluginXID), ext.Token
}

func buildConfigurationManifest(name string) openapi.PluginManifest {
	id := "config-" + strings.ReplaceAll(strings.ToLower(name), " ", "-")

	return openapi.PluginManifest(map[string]any{
		"id":          id,
		"name":        name,
		"author":      "test-author",
		"description": "Config test plugin: " + name,
		"version":     "1.0.0",
		"command":     "./plugin",
		"configuration_schema": map[string]any{
			"fields": []map[string]any{
				{"type": "string", "id": "name", "label": "Name"},
				{"type": "boolean", "id": "enabled", "label": "Enabled"},
				{"type": "number", "id": "threshold", "label": "Threshold"},
			},
		},
	})
}

func parseConfigurationField(field openapi.PluginConfigurationFieldUnion) (string, string) {
	var raw map[string]any
	_ = json.Unmarshal(mustJSON(field), &raw)
	id, _ := raw["id"].(string)
	typ, _ := raw["type"].(string)
	return id, typ
}

func mustJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}

func runExternalSDKConfigPlugin(
	ctx context.Context,
	t *testing.T,
	rpcURL string,
	onConfigure func(map[string]any) error,
) (func(), <-chan error) {
	t.Helper()

	t.Setenv("STORYDEN_RPC_URL", rpcURL)
	pl, err := sdk.New(ctx)
	require.NoError(t, err)

	if onConfigure != nil {
		pl.OnConfigure(func(_ context.Context, config map[string]any) error {
			return onConfigure(config)
		})
	}

	runCtx, cancel := context.WithCancel(ctx)
	done := make(chan error, 1)
	go func() {
		done <- pl.Run(runCtx)
	}()

	stop := func() {
		_ = pl.Shutdown()
		cancel()
	}

	return stop, done
}
