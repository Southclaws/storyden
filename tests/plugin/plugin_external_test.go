package plugin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Southclaws/opt"
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
	"github.com/Southclaws/storyden/lib/plugin/rpc"
	sdk "github.com/Southclaws/storyden/sdk/go/storyden"
	"github.com/Southclaws/storyden/tests"
)

func TestExternalPluginEventSubscription(t *testing.T) {
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
			pluginID, installationID, token := addExternalPlugin(t, root, cl, adminSession, r, "External Event Listener", []string{"EventThreadPublished"})

			r.True(strings.HasPrefix(token, "sdprt_"))

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateConnecting)

			eventCh := make(chan string, 1)
			rpcURL := externalRPCURL(t, ts.URL, token)
			stopPlugin, pluginDone := runExternalSDKPlugin(root, t, rpcURL, func(eventID string) {
				select {
				case eventCh <- eventID:
				default:
				}
			})
			defer stopPlugin()
			defer waitForPluginStop(pluginDone)

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateActive)

			thread := tests.AssertRequest(
				cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       opt.New("<p>external plugin event test</p>").Ptr(),
					Visibility: opt.New(openapi.Published).Ptr(),
					Title:      "External Plugin Event " + xid.New().String(),
				}, adminSession),
			)(t, http.StatusOK)

			select {
			case got := <-eventCh:
				r.Equal(thread.JSON200.Id, got)
			case <-time.After(5 * time.Second):
				r.FailNow("timed out waiting for event from external plugin")
			}

			getResp := tests.AssertRequest(
				cl.PluginGetWithResponse(root, pluginID, adminSession),
			)(t, http.StatusOK)

			r.Equal("External Event Listener", getResp.JSON200.Name)
			r.NotNil(getResp.JSON200.Version)

			ext, err := getResp.JSON200.Connection.AsPluginExternalProps()
			r.NoError(err)
			r.Equal(token, ext.Token)
		}))
	}))
}

func TestExternalPluginDisconnectReconnect(t *testing.T) {
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
			_, installationID, token := addExternalPlugin(t, root, cl, adminSession, r, "External Reconnect", []string{})

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateConnecting)

			rpcURL := externalRPCURL(t, ts.URL, token)

			stopPlugin1, pluginDone1 := runExternalSDKPlugin(root, t, rpcURL, nil)
			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateActive)

			stopPlugin1()
			waitForPluginStop(pluginDone1)
			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateConnecting)

			stopPlugin2, pluginDone2 := runExternalSDKPlugin(root, t, rpcURL, nil)
			defer stopPlugin2()
			defer waitForPluginStop(pluginDone2)

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateActive)
		}))
	}))
}

func TestExternalPluginSetActiveStateRejected(t *testing.T) {
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
			pluginID, installationID, _ := addExternalPlugin(t, root, cl, adminSession, r, "External Active State", []string{})

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateConnecting)

			tests.AssertRequest(
				cl.PluginSetActiveStateWithResponse(root, pluginID, openapi.PluginSetActiveStateJSONRequestBody{
					Active: openapi.PluginActiveStateActive,
				}, adminSession),
			)(t, http.StatusBadRequest)

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateConnecting)
		}))
	}))
}

func TestExternalPluginGetLogsRejected(t *testing.T) {
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
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)

			adminSession, _ := createAdminSession(root, cl, sh, accountWrite, r)
			pluginID, _, _ := addExternalPlugin(t, root, cl, adminSession, r, "External Logs", []string{})

			tests.AssertRequest(
				cl.PluginGetLogsWithResponse(root, pluginID, adminSession),
			)(t, http.StatusBadRequest)
		}))
	}))
}

func TestExternalPluginManifestUpdateReconnectsWithNewSubscriptions(t *testing.T) {
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
				"External Manifest Update",
				[]string{"EventThreadPublished"},
			)

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateConnecting)

			eventCh1 := make(chan string, 1)
			rpcURL := externalRPCURL(t, ts.URL, token)
			stopPlugin1, pluginDone1 := runExternalSDKPlugin(root, t, rpcURL, func(eventID string) {
				select {
				case eventCh1 <- eventID:
				default:
				}
			})
			defer stopPlugin1()

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateActive)

			thread1 := tests.AssertRequest(
				cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       opt.New("<p>manifest update before change</p>").Ptr(),
					Visibility: opt.New(openapi.Published).Ptr(),
					Title:      "Manifest Update Before " + xid.New().String(),
				}, adminSession),
			)(t, http.StatusOK)

			select {
			case got := <-eventCh1:
				r.Equal(thread1.JSON200.Id, got)
			case <-time.After(5 * time.Second):
				r.FailNow("timed out waiting for event from external plugin before manifest update")
			}

			updatedManifest := buildTestManifest("External Manifest Update", []string{})
			updateResp := tests.AssertRequest(
				cl.PluginUpdateManifestWithResponse(
					root,
					pluginID,
					openapi.PluginUpdateManifestJSONRequestBody(updatedManifest),
					adminSession,
				),
			)(t, http.StatusOK)
			r.NotNil(updateResp.JSON200)

			updatedConn, err := updateResp.JSON200.Connection.AsPluginExternalProps()
			r.NoError(err)
			r.Equal(token, updatedConn.Token)

			select {
			case <-pluginDone1:
			case <-time.After(3 * time.Second):
				t.Fatal("expected connected client to be disconnected when manifest is updated")
			}

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateConnecting)

			eventCh2 := make(chan string, 1)
			stopPlugin2, pluginDone2 := runExternalSDKPlugin(root, t, rpcURL, func(eventID string) {
				select {
				case eventCh2 <- eventID:
				default:
				}
			})
			defer stopPlugin2()
			defer waitForPluginStop(pluginDone2)

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateActive)

			tests.AssertRequest(
				cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       opt.New("<p>manifest update after change</p>").Ptr(),
					Visibility: opt.New(openapi.Published).Ptr(),
					Title:      "Manifest Update After " + xid.New().String(),
				}, adminSession),
			)(t, http.StatusOK)

			select {
			case got := <-eventCh2:
				t.Fatalf("did not expect thread published event after manifest update, got: %s", got)
			case <-time.After(1500 * time.Millisecond):
			}
		}))
	}))
}

func TestPluginUpdateManifestRejectsSupervisedPlugin(t *testing.T) {
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
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)

			adminSession, _ := createAdminSession(root, cl, sh, accountWrite, r)

			archivePath := packageTestPluginArchive(t, "test_data/event_listener", "event_listener")
			pluginFile, err := os.Open(archivePath)
			r.NoError(err)
			defer pluginFile.Close()

			addResp := tests.AssertRequest(
				cl.PluginAddWithBodyWithResponse(root, "application/zip", pluginFile, adminSession),
			)(t, http.StatusOK)
			pluginID := string(addResp.JSON200.Id)

			tests.AssertRequest(
				cl.PluginUpdateManifestWithResponse(
					root,
					pluginID,
					openapi.PluginUpdateManifestJSONRequestBody(
						buildTestManifest("Should Fail", []string{}),
					),
					adminSession,
				),
			)(t, http.StatusBadRequest)
		}))
	}))
}

func TestExternalPluginCycleTokenInvalidatesConnectedSession(t *testing.T) {
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
				"External Rotate Token",
				[]string{},
			)
			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateConnecting)

			oldRPCURL := externalRPCURL(t, ts.URL, token)
			stopOld, oldDone := runExternalSDKPlugin(root, t, oldRPCURL, nil)
			defer stopOld()

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateActive)

			cycleResp := tests.AssertRequest(
				cl.PluginCycleTokenWithResponse(root, pluginID, adminSession),
			)(t, http.StatusOK)
			newToken := cycleResp.JSON200.Token
			r.NotEqual(token, newToken)

			select {
			case <-oldDone:
			case <-time.After(3 * time.Second):
				t.Fatal("expected connected client to be disconnected when token is cycled")
			}

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateConnecting)

			t.Setenv("STORYDEN_RPC_URL", oldRPCURL)
			oldClient, err := sdk.New(root)
			r.NoError(err)
			reconnectCtx, cancel := context.WithTimeout(root, 2*time.Second)
			defer cancel()
			err = oldClient.Run(reconnectCtx)
			r.Error(err)

			newRPCURL := externalRPCURL(t, ts.URL, newToken)
			stopNew, newDone := runExternalSDKPlugin(root, t, newRPCURL, nil)
			defer stopNew()
			defer waitForPluginStop(newDone)

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateActive)
		}))
	}))
}

func createAdminSession(
	ctx context.Context,
	cl *openapi.ClientWithResponses,
	sh *e2e.SessionHelper,
	accountWrite *account_writer.Writer,
	r *require.Assertions,
) (openapi.RequestEditorFn, account.AccountID) {
	adminHandle := "admin-" + xid.New().String()

	admin, err := cl.AuthPasswordSignupWithResponse(ctx, nil, openapi.AuthPair{
		Identifier: adminHandle,
		Token:      "password",
	})
	r.NoError(err)
	r.Equal(http.StatusOK, admin.StatusCode())

	adminID := account.AccountID(utils.Must(xid.FromString(admin.JSON200.Id)))
	accountWrite.Update(ctx, adminID, account_writer.SetAdmin(true))

	return sh.WithSession(e2e.WithAccountID(ctx, adminID)), adminID
}

func addExternalPlugin(
	t *testing.T,
	ctx context.Context,
	cl *openapi.ClientWithResponses,
	adminSession openapi.RequestEditorFn,
	r *require.Assertions,
	name string,
	events []string,
) (string, resource_plugin.InstallationID, string) {
	manifest := buildTestManifest(name, events)

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
	r.Equal(name, getResp.JSON200.Name)
	r.NotNil(getResp.JSON200.Version)

	ext, err := getResp.JSON200.Connection.AsPluginExternalProps()
	r.NoError(err)
	token := ext.Token
	r.NotEmpty(token)

	return string(plugin.Id), resource_plugin.InstallationID(pluginXID), token
}

func buildTestManifest(name string, events []string) openapi.PluginManifest {
	id := "test-" + strings.ReplaceAll(strings.ToLower(name), " ", "-")

	return openapi.PluginManifest(map[string]any{
		"id":              id,
		"name":            name,
		"author":          "test-author",
		"description":     "Test plugin: " + name,
		"version":         "1.0.0",
		"command":         "./plugin",
		"events_consumed": events,
	})
}

func externalRPCURL(t *testing.T, serverURL string, token string) string {
	t.Helper()

	u, err := url.Parse(serverURL)
	require.NoError(t, err)

	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	default:
		t.Fatalf("unexpected test server scheme: %s", u.Scheme)
	}

	u.Path = "/rpc"
	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()

	return u.String()
}

func runExternalSDKPlugin(
	ctx context.Context,
	t *testing.T,
	rpcURL string,
	onThreadPublished func(string),
) (func(), <-chan error) {
	t.Helper()

	t.Setenv("STORYDEN_RPC_URL", rpcURL)
	pl, err := sdk.New(ctx)
	require.NoError(t, err)

	if onThreadPublished != nil {
		pl.OnThreadPublished(func(_ context.Context, event *rpc.EventThreadPublished) error {
			onThreadPublished(event.ID.String())
			return nil
		})
	}

	runCtx, cancel := context.WithCancel(ctx)
	done := make(chan error, 1)
	go func() {
		done <- pl.Run(runCtx)
	}()

	return cancel, done
}

func waitForPluginStop(done <-chan error) {
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
	}
}

func requireSessionState(
	t *testing.T,
	ctx context.Context,
	runner plugin_runner.Host,
	installationID resource_plugin.InstallationID,
	want resource_plugin.ReportedState,
) {
	t.Helper()

	var last resource_plugin.ReportedState
	require.Eventually(t, func() bool {
		sess, err := runner.GetSession(ctx, installationID)
		if err != nil || sess == nil {
			return false
		}

		last = sess.GetReportedState()
		return last == want
	}, 8*time.Second, 50*time.Millisecond, "expected state %s, last state %s", want, last)
}
