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

			var initialStartedAt time.Time
			r.Eventually(func() bool {
				s, err := runner.GetSession(root, installationID)
				if err != nil || s == nil {
					return false
				}
				initialStartedAt = s.GetStartedAt().OrZero()
				return !initialStartedAt.IsZero()
			}, 3*time.Second, 20*time.Millisecond, "plugin should report initial started_at timestamp")

			threadTitle := "Crash Test Thread " + xid.New().String()
			threadCreate := tests.AssertRequest(
				cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       opt.New("<p>test plugin crash handling</p>").Ptr(),
					Visibility: opt.New(openapi.Published).Ptr(),
					Title:      threadTitle,
				}, adminSession),
			)(t, http.StatusOK)

			_ = threadCreate.JSON200.Id

			foundRestarting := false
			foundRecovery := false
			var restartMessage string
			for i := 0; i < 300; i++ {
				sess, err := runner.GetSession(root, installationID)
				r.NoError(err)
				r.NotNil(sess)

				currentState := sess.GetReportedState()
				if currentState == resource_plugin.ReportedStateRestarting {
					foundRestarting = true
					restartMessage = sess.GetErrorMessage()
				}

				if currentState == resource_plugin.ReportedStateActive {
					startedAt := sess.GetStartedAt().OrZero()
					if startedAt.After(initialStartedAt) {
						foundRecovery = true
						if foundRestarting {
							break
						}
					}
				}

				time.Sleep(10 * time.Millisecond)
			}

			a.True(foundRecovery, "plugin should automatically restart after crash")
			if foundRestarting {
				a.NotEmpty(restartMessage, "restarting state should include an error message")
				a.Contains(restartMessage, "exit status", "error message should describe the crash")
			}

			getResp, err := cl.PluginGetWithResponse(root, pluginID, adminSession)
			r.NoError(err)
			r.Equal(http.StatusOK, getResp.StatusCode())
			r.NotNil(getResp.JSON200)

			statusKind, err := getResp.JSON200.Status.Discriminator()
			r.NoError(err)
			a.Contains([]string{
				resource_plugin.ReportedStateActive.String(),
				resource_plugin.ReportedStateRestarting.String(),
			}, statusKind, "API status should be active or restarting")

			if statusKind == resource_plugin.ReportedStateRestarting.String() {
				statusRestarting, err := getResp.JSON200.Status.AsPluginStatusRestarting()
				r.NoError(err, "plugin status should be restarting")
				a.Equal(openapi.Restarting, statusRestarting.ActiveState, "API should report restarting state")
				a.NotEmpty(statusRestarting.Message, "status message should not be empty")
				a.Contains(statusRestarting.Message, "exit status", "error message should describe the crash")
			}
		}))
	}))
}
