package plugin_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"unicode"

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
					"permissions": []string{"MANAGE_CATEGORIES"},
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
			accessAuth := accessKeyAuth(access.AccessKey)

			manualAccount := tests.AssertRequest(
				cl.AccountGetWithResponse(root, accessAuth),
			)(t, http.StatusOK)
			r.Contains(string(manualAccount.JSON200.Handle), botHandle)
			expectedRoleName := fmt.Sprintf(
				"External Access Bot (Bot %s)",
				pluginRoleShortIDForTest(string(manualAccount.JSON200.Handle)),
			)
			accessRole, found := findRoleByName(manualAccount.JSON200.Roles, expectedRoleName)
			r.True(found, "expected managed plugin role %q to exist", expectedRoleName)
			r.Contains(accessRole.Permissions, openapi.Permission("MANAGE_CATEGORIES"))

			tests.AssertRequest(
				cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
					Name:        "Plugin Access Category " + xid.New().String(),
					Description: "plugin bot category create check",
					Colour:      "hsl(120, 60%, 45%)",
				}, accessAuth),
			)(t, http.StatusOK)

			tests.AssertRequest(
				cl.RoleCreateWithResponse(root, openapi.RoleCreateJSONRequestBody{
					Name:        "plugin-access-no-manage-roles-" + xid.New().String(),
					Colour:      "hsl(10, 60%, 45%)",
					Permissions: openapi.PermissionList{openapi.MANAGECATEGORIES},
				}, accessAuth),
			)(t, http.StatusForbidden)

			builtClient, err := pl.BuildAPIClient(root)
			r.NoError(err)
			rawClient, ok := builtClient.ClientInterface.(*openapi.Client)
			r.True(ok, "expected *openapi.Client")
			rawClient.Server = strings.TrimRight(ts.URL, "/") + "/api/"

			builtAccount := tests.AssertRequest(
				builtClient.AccountGetWithResponse(root),
			)(t, http.StatusOK)
			r.Contains(string(builtAccount.JSON200.Handle), botHandle)
		}))
	}))
}

func TestExternalPluginAccessSubmitLibraryNodeInEmailMode(t *testing.T) {
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
			emailMode := openapi.Email
			tests.AssertRequest(
				cl.AdminSettingsUpdateWithResponse(root, openapi.AdminSettingsUpdateJSONRequestBody{
					AuthenticationMode: &emailMode,
				}, adminSession),
			)(t, http.StatusOK)

			manifest := openapi.PluginManifest(map[string]any{
				"id":              "test-access-submit-" + xid.New().String(),
				"name":            "External Access Submit",
				"author":          "test-author",
				"description":     "External plugin submit access test",
				"version":         "1.0.0",
				"command":         "./plugin",
				"events_consumed": []string{},
				"access": map[string]any{
					"handle":      "submit-" + xid.New().String(),
					"name":        "Submit Access Bot",
					"permissions": []string{"SUBMIT_LIBRARY_NODE"},
				},
			})
			_, installationID, token := addExternalPluginWithManifest(t, root, cl, adminSession, r, manifest)

			rpcURL := externalRPCURL(t, ts.URL, token)
			pl, stop, done := runExternalSDKAccessPlugin(root, t, rpcURL)
			defer waitForPluginStop(done)
			defer stop()

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateActive)

			access, err := pl.GetAccess(root)
			r.NoError(err)
			accessAuth := accessKeyAuth(access.AccessKey)

			review := openapi.Review
			nodeCreate := tests.AssertRequest(
				cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
					Name:       "plugin-submit-review-" + xid.New().String(),
					Visibility: &review,
				}, accessAuth),
			)(t, http.StatusOK)
			r.Equal(openapi.Review, nodeCreate.JSON200.Visibility)
		}))
	}))
}

func TestExternalPluginAccessPermissionsSyncOnManifestUpdate(t *testing.T) {
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
			botHandle := "sync-" + xid.New().String()

			manifestV1 := openapi.PluginManifest(map[string]any{
				"id":              "test-access-sync-" + xid.New().String(),
				"name":            "External Access Sync",
				"author":          "test-author",
				"description":     "External plugin access sync test",
				"version":         "1.0.0",
				"command":         "./plugin",
				"events_consumed": []string{},
				"access": map[string]any{
					"handle":      botHandle,
					"name":        "Sync Access Bot",
					"permissions": []string{"MANAGE_CATEGORIES"},
				},
			})

			pluginID, installationID, token := addExternalPluginWithManifest(t, root, cl, adminSession, r, manifestV1)
			rpcURL := externalRPCURL(t, ts.URL, token)

			pl1, stop1, done1 := runExternalSDKAccessPlugin(root, t, rpcURL)
			defer waitForPluginStop(done1)
			defer stop1()

			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateActive)

			accessV1, err := pl1.GetAccess(root)
			r.NoError(err)
			authV1 := accessKeyAuth(accessV1.AccessKey)

			tests.AssertRequest(
				cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
					Name:        "Plugin Access v1 Category " + xid.New().String(),
					Description: "manifest v1 manage categories check",
					Colour:      "hsl(145, 58%, 45%)",
				}, authV1),
			)(t, http.StatusOK)
			tests.AssertRequest(
				cl.RoleCreateWithResponse(root, openapi.RoleCreateJSONRequestBody{
					Name:        "plugin-sync-v1-" + xid.New().String(),
					Colour:      "hsl(20, 55%, 45%)",
					Permissions: openapi.PermissionList{openapi.MANAGECATEGORIES},
				}, authV1),
			)(t, http.StatusForbidden)

			manifestV2 := openapi.PluginManifest(map[string]any{
				"id":              manifestV1["id"],
				"name":            manifestV1["name"],
				"author":          manifestV1["author"],
				"description":     manifestV1["description"],
				"version":         "1.0.1",
				"command":         manifestV1["command"],
				"events_consumed": []string{},
				"access": map[string]any{
					"handle":      botHandle,
					"name":        "Sync Access Bot",
					"permissions": []string{"MANAGE_CATEGORIES", "MANAGE_ROLES"},
				},
			})
			tests.AssertRequest(
				cl.PluginUpdateManifestWithResponse(root, pluginID, openapi.PluginUpdateManifestJSONRequestBody(manifestV2), adminSession),
			)(t, http.StatusOK)

			select {
			case <-done1:
			case <-time.After(5 * time.Second):
				t.Fatal("expected plugin connection to close after manifest update")
			}
			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateConnecting)

			pl2, stop2, done2 := runExternalSDKAccessPlugin(root, t, rpcURL)
			defer waitForPluginStop(done2)
			defer stop2()
			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateActive)

			accessV2, err := pl2.GetAccess(root)
			r.NoError(err)
			authV2 := accessKeyAuth(accessV2.AccessKey)
			tests.AssertRequest(
				cl.RoleCreateWithResponse(root, openapi.RoleCreateJSONRequestBody{
					Name:        "plugin-sync-v2-" + xid.New().String(),
					Colour:      "hsl(200, 65%, 45%)",
					Permissions: openapi.PermissionList{openapi.MANAGECATEGORIES},
				}, authV2),
			)(t, http.StatusOK)

			manifestV3 := openapi.PluginManifest(map[string]any{
				"id":              manifestV1["id"],
				"name":            manifestV1["name"],
				"author":          manifestV1["author"],
				"description":     manifestV1["description"],
				"version":         "1.0.2",
				"command":         manifestV1["command"],
				"events_consumed": []string{},
				"access": map[string]any{
					"handle":      botHandle,
					"name":        "Sync Access Bot",
					"permissions": []string{"MANAGE_CATEGORIES"},
				},
			})
			tests.AssertRequest(
				cl.PluginUpdateManifestWithResponse(root, pluginID, openapi.PluginUpdateManifestJSONRequestBody(manifestV3), adminSession),
			)(t, http.StatusOK)

			select {
			case <-done2:
			case <-time.After(5 * time.Second):
				t.Fatal("expected plugin connection to close after manifest update")
			}
			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateConnecting)

			pl3, stop3, done3 := runExternalSDKAccessPlugin(root, t, rpcURL)
			defer waitForPluginStop(done3)
			defer stop3()
			requireSessionState(t, root, runner, installationID, resource_plugin.ReportedStateActive)

			accessV3, err := pl3.GetAccess(root)
			r.NoError(err)
			authV3 := accessKeyAuth(accessV3.AccessKey)

			tests.AssertRequest(
				cl.RoleCreateWithResponse(root, openapi.RoleCreateJSONRequestBody{
					Name:        "plugin-sync-v3-" + xid.New().String(),
					Colour:      "hsl(260, 55%, 45%)",
					Permissions: openapi.PermissionList{openapi.MANAGECATEGORIES},
				}, authV3),
			)(t, http.StatusForbidden)

			accountV3 := tests.AssertRequest(
				cl.AccountGetWithResponse(root, authV3),
			)(t, http.StatusOK)
			expectedRoleName := fmt.Sprintf(
				"Sync Access Bot (Bot %s)",
				pluginRoleShortIDForTest(string(accountV3.JSON200.Handle)),
			)
			accessRole, found := findRoleByName(accountV3.JSON200.Roles, expectedRoleName)
			r.True(found, "expected managed plugin role %q to exist", expectedRoleName)
			r.Contains(accessRole.Permissions, openapi.Permission("MANAGE_CATEGORIES"))
			r.NotContains(accessRole.Permissions, openapi.Permission("MANAGE_ROLES"))
		}))
	}))
}

func runExternalSDKAccessPlugin(ctx context.Context, t *testing.T, rpcURL string) (*sdk.Plugin, func(), <-chan error) {
	t.Helper()

	t.Setenv("STORYDEN_RPC_URL", rpcURL)
	pl, err := sdk.New(ctx)
	require.NoError(t, err)

	runCtx, cancel := context.WithCancel(ctx)
	done := make(chan error, 1)
	go func() {
		done <- pl.Run(runCtx)
	}()

	stop := func() {
		_ = pl.Shutdown()
		cancel()
	}

	return pl, stop, done
}

func accessKeyAuth(accessKey string) openapi.RequestEditorFn {
	return func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+accessKey)
		return nil
	}
}

func findRoleByName(roles openapi.AccountRoleList, name string) (openapi.AccountRole, bool) {
	for _, role := range roles {
		if role.Name == name {
			return role, true
		}
	}
	return openapi.AccountRole{}, false
}

func pluginRoleShortIDForTest(handle string) string {
	handle = strings.ToLower(strings.TrimSpace(handle))

	clean := strings.Builder{}
	for _, r := range handle {
		if unicode.IsLower(r) || unicode.IsDigit(r) {
			clean.WriteRune(r)
		}
	}

	shortID := clean.String()
	if len(shortID) >= 4 {
		return shortID[len(shortID)-4:]
	}
	if shortID == "" {
		return "0000"
	}
	for len(shortID) < 4 {
		shortID += "0"
	}
	return shortID
}
