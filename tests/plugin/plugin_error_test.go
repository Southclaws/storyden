package plugin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Southclaws/opt"
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
	"github.com/Southclaws/storyden/tests"
)

func TestPluginErrorStates(t *testing.T) {
	t.Parallel()

	integration.Test(t, &config.Config{
		// LogLevel: slog.LevelDebug,
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

			archivePath := packageTestPluginArchive(t, "test_data/crash_test", "crash_test")
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
			r.Equal(http.StatusOK, activeResp.StatusCode())

			sess, err := runner.GetSession(root, installationID)
			r.NoError(err)
			r.NotNil(sess)

			state := sess.GetReportedState()
			r.Equal(resource_plugin.ReportedStateActive, state, "plugin should be active")

			threadTitle := "Crash Test Thread " + xid.New().String()
			threadCreate := tests.AssertRequest(
				cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       opt.New("<p>test plugin crash handling</p>").Ptr(),
					Visibility: opt.New(openapi.Published).Ptr(),
					Title:      threadTitle,
				}, adminSession),
			)(t, http.StatusOK)

			_ = threadCreate.JSON200.Id

			foundError := false
			var errorMessage string
			var currentState resource_plugin.ReportedState
			for i := 0; i < 30; i++ {
				sess, err := runner.GetSession(root, installationID)
				r.NoError(err)
				r.NotNil(sess)

				currentState = sess.GetReportedState()
				if currentState == resource_plugin.ReportedStateRestarting {
					foundError = true
					errorMessage = sess.GetErrorMessage()
					break
				}

				time.Sleep(100 * time.Millisecond)
			}

			a.True(foundError, "plugin should report error state after crash")
			a.NotEmpty(errorMessage, "error message should be populated")

			getResp, err := cl.PluginGetWithResponse(root, pluginID, adminSession)
			r.NoError(err)
			r.Equal(http.StatusOK, getResp.StatusCode())
			r.NotNil(getResp.JSON200)

			statusRestarting, err := getResp.JSON200.Status.AsPluginStatusRestarting()
			r.NoError(err, "plugin status should be restarting")
			a.Equal(openapi.Restarting, statusRestarting.ActiveState, "API should report restarting state")
			a.NotEmpty(statusRestarting.Message, "status message should not be empty")
			a.Contains(statusRestarting.Message, "exit status", "error message should describe the crash")

			foundRecovery := false
			for i := 0; i < 50; i++ {
				sess, err := runner.GetSession(root, installationID)
				if err != nil {
					time.Sleep(100 * time.Millisecond)
					continue
				}

				currentState := sess.GetReportedState()
				if currentState == resource_plugin.ReportedStateActive {
					foundRecovery = true
					break
				}

				time.Sleep(100 * time.Millisecond)
			}

			a.True(foundRecovery, "plugin should automatically restart after crash")
		}))
	}))
}
