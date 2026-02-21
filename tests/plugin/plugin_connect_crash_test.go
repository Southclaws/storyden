package plugin_test

import (
	"context"
	"log/slog"
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

func TestPluginCrashOnConnect(t *testing.T) {
	integration.Test(t, &config.Config{
		LogLevel:                    slog.LevelInfo,
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
			adminSession := sh.WithSession(e2e.WithAccountID(root, adminID))

			accountWrite.Update(root, adminID, account_writer.SetAdmin(true))

			archivePath := packageTestPluginArchive(t, "test_data/crash_connect", "crash_connect")
			pluginFile, err := os.Open(archivePath)
			r.NoError(err)
			defer pluginFile.Close()

			addResp, err := cl.PluginAddWithBodyWithResponse(root, "application/zip", pluginFile, adminSession)
			r.NoError(err)
			r.Equal(http.StatusOK, addResp.StatusCode())
			r.NotNil(addResp.JSON200)

			pluginID := string(addResp.JSON200.Id)
			pluginIDxid, _ := xid.FromString(pluginID)
			installationID := resource_plugin.InstallationID(pluginIDxid)

			activeResp, err := cl.PluginSetActiveStateWithResponse(root, pluginID, openapi.PluginSetActiveStateJSONRequestBody{
				Active: "active",
			}, adminSession)
			r.NoError(err)
			r.Equal(http.StatusInternalServerError, activeResp.StatusCode())
			a.Equal("exit status 42", activeResp.JSONDefault.Error)

			sess, err := runner.GetSession(root, installationID)
			r.NoError(err)
			r.NotNil(sess)

			currentState := sess.GetReportedState()
			a.Equal(resource_plugin.ReportedStateError, currentState, "plugin should be in error state due to startup crashes")

			errorMsg := sess.GetErrorMessage()
			a.Contains(errorMsg, "exit status", "error message should describe the crash")
		}))
	}))
}
