package plugin_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	resource_plugin "github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	rpc_transport "github.com/Southclaws/storyden/app/transports/rpc"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/tests"
)

func TestSupervisedPluginConfigurationLifecycle(t *testing.T) {
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

			type serverURLSetter interface {
				SetServerURL(string)
			}
			if setter, ok := runner.(serverURLSetter); ok {
				setter.SetServerURL(ts.URL)
			}

			adminHandle := "admin-" + xid.New().String()
			admin, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{
				Identifier: adminHandle,
				Token:      "password",
			})
			r.NoError(err)
			r.Equal(http.StatusOK, admin.StatusCode())
			adminID := account.AccountID(utils.Must(xid.FromString(admin.JSON200.Id)))
			accountWrite.Update(root, adminID, account_writer.SetAdmin(true))
			adminSession := sh.WithSession(e2e.WithAccountID(root, adminID))

			archivePath := packageTestPluginArchive(t, "test_data/supervised_config", "supervised_config")
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

			tests.AssertRequest(
				cl.PluginSetActiveStateWithResponse(root, pluginID, openapi.PluginSetActiveStateJSONRequestBody{
					Active: openapi.PluginActiveStateActive,
				}, adminSession),
			)(t, http.StatusOK)

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateActive)

			schemaResp := tests.AssertRequest(
				cl.PluginGetConfigurationSchemaWithResponse(root, pluginID, adminSession),
			)(t, http.StatusOK)
			r.NotNil(schemaResp.JSON200.Fields)
			r.Len(*schemaResp.JSON200.Fields, 3)

			updateResp := tests.AssertRequest(
				cl.PluginUpdateConfigurationWithResponse(root, pluginID, openapi.PluginUpdateConfigurationJSONRequestBody{
					"name":      "supervised-configured",
					"enabled":   true,
					"threshold": 7.5,
				}, adminSession),
			)(t, http.StatusOK)
			updateCfg := openapi.PluginConfiguration(*updateResp.JSON200)
			r.Equal("supervised-configured", updateCfg["name"])
			r.Equal(true, updateCfg["enabled"])
			r.Equal(7.5, updateCfg["threshold"])

			getResp := tests.AssertRequest(
				cl.PluginGetConfigurationWithResponse(root, pluginID, adminSession),
			)(t, http.StatusOK)
			getCfg := openapi.PluginConfiguration(*getResp.JSON200)
			r.Equal("supervised-configured", getCfg["name"])
			r.Equal(true, getCfg["enabled"])
			r.Equal(7.5, getCfg["threshold"])

			configuredFile := filepath.Join("data", "plugins", pluginID, "configured.json")
			require.Eventually(t, func() bool {
				_, err := os.Stat(configuredFile)
				return err == nil
			}, 5*time.Second, 100*time.Millisecond)

			b, err := os.ReadFile(configuredFile)
			r.NoError(err)

			var applied map[string]any
			err = json.Unmarshal(b, &applied)
			r.NoError(err)
			r.Equal("supervised-configured", applied["name"])
			r.Equal(true, applied["enabled"])
			r.Equal(7.5, applied["threshold"])
		}))
	}))
}
