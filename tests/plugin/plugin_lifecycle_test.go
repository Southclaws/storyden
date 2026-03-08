package plugin_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	resource_plugin "github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/app/transports/rpc"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/utils"
)

func TestPluginLifecycle(t *testing.T) {
	t.Parallel()

	integration.Test(t, &config.Config{
		PluginMaxRestartAttempts:    1,
		PluginMaxBackoffDuration:    100 * time.Millisecond,
		PluginRuntimeCrashThreshold: 1 * time.Second,
		PluginRuntimeCrashBackoff:   100 * time.Millisecond,
	}, e2e.Setup(), rpc.Build(), fx.Invoke(func(
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
			a := assert.New(t)

			type serverURLSetter interface {
				SetServerURL(string)
			}
			if setter, ok := runner.(serverURLSetter); ok {
				t.Logf("setting server URL to: %s", ts.URL)
				setter.SetServerURL(ts.URL)
			} else {
				t.Logf("runner does not implement SetServerURL")
			}

			adminHandle := "admin-" + xid.New().String()

			admin, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{
				Identifier: adminHandle,
				Token:      "password",
			})
			r.NoError(err)
			r.Equal(http.StatusOK, admin.StatusCode())
			adminID := account.AccountID(utils.Must(xid.FromString(admin.JSON200.Id)))
			adminSession := sh.WithSession(e2e.WithAccountID(root, adminID))

			accountWrite.Update(root, adminID, account_writer.SetAdmin(true))

			archivePath := packageTestPluginArchive(t, "test_data/example", "example-plugin")
			pluginFile, err := os.Open(archivePath)
			r.NoError(err)
			defer pluginFile.Close()

			addResp, err := cl.PluginAddWithBodyWithResponse(root, "application/zip", pluginFile, adminSession)
			r.NoError(err)
			r.Equal(http.StatusOK, addResp.StatusCode())
			r.NotNil(addResp.JSON200)

			pluginID := string(addResp.JSON200.Id)
			t.Logf("created plugin: %s", pluginID)

			getResp1, err := cl.PluginGetWithResponse(root, pluginID, adminSession)
			r.NoError(err)
			r.Equal(http.StatusOK, getResp1.StatusCode())
			r.NotNil(getResp1.JSON200)

			status1, err := getResp1.JSON200.Status.AsPluginStatusInactive()
			r.NoError(err)
			a.Equal(openapi.Inactive, status1.ActiveState)
			t.Logf("plugin status: %s", status1.ActiveState)

			pluginIDxid, _ := xid.FromString(pluginID)
			installationID := resource_plugin.InstallationID(pluginIDxid)

			sess1, err := runner.GetSession(root, installationID)
			r.Error(err)
			r.Nil(sess1)

			activeResp, err := cl.PluginSetActiveStateWithResponse(root, pluginID, openapi.PluginSetActiveStateJSONRequestBody{
				Active: "active",
			}, adminSession)
			r.NoError(err)
			r.Equal(http.StatusOK, activeResp.StatusCode())
			r.NotNil(activeResp.JSON200)

			status2, err := activeResp.JSON200.Status.AsPluginStatusActive()
			r.NoError(err)
			a.Equal(openapi.PluginStatusActiveActiveStateActive, status2.ActiveState)
			t.Logf("plugin activated: %s", status2.ActiveState)

			sess2, err := runner.GetSession(root, installationID)
			r.NoError(err)
			r.NotNil(sess2)
			t.Logf("session found: %s", sess2.ID())

			s, _ := runner.GetSessions(root)
			fmt.Println(s)

			state := sess2.GetReportedState()
			a.Equal(resource_plugin.ReportedStateActive, state)
			t.Logf("session state: %s", state)

			inactiveResp, err := cl.PluginSetActiveStateWithResponse(root, pluginID, openapi.PluginSetActiveStateJSONRequestBody{
				Active: "inactive",
			}, adminSession)
			r.NoError(err)
			r.Equal(http.StatusOK, inactiveResp.StatusCode())
			r.NotNil(inactiveResp.JSON200)

			status3, err := inactiveResp.JSON200.Status.AsPluginStatusInactive()
			r.NoError(err)
			a.Equal(openapi.Inactive, status3.ActiveState)
			t.Logf("plugin deactivated: %s", status3.ActiveState)

			state2 := sess2.GetReportedState()
			a.NotEqual(resource_plugin.ReportedStateActive, state2)
			t.Logf("session state after deactivation: %s", state2)

			deleteResp, err := cl.PluginDeleteWithResponse(root, pluginID, adminSession)
			r.NoError(err)
			r.Equal(http.StatusNoContent, deleteResp.StatusCode())
			t.Logf("plugin deleted")
		}))
	}))
}
