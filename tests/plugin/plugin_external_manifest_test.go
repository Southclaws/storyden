package plugin_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	rpc_transport "github.com/Southclaws/storyden/app/transports/rpc"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/tests"
)

func TestExternalPluginManifestValidation(t *testing.T) {
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
	) {
		lc.Append(fx.StartHook(func() {
			adminHandle := "admin-" + xid.New().String()
			admin, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{
				Identifier: adminHandle,
				Token:      "password",
			})
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, admin.StatusCode())

			adminID := account.AccountID(utils.Must(xid.FromString(admin.JSON200.Id)))
			accountWrite.Update(root, adminID, account_writer.SetAdmin(true))
			adminSession := sh.WithSession(e2e.WithAccountID(root, adminID))

			addExternal := func(manifest map[string]any) (*openapi.PluginAddResponse, error) {
				require.NoError(t, err)
				body := openapi.PluginInitialProps{}
				require.NoError(t, body.FromPluginInitialExternal(openapi.PluginInitialExternal{
					Mode:     openapi.External,
					Manifest: manifest,
				}))
				return cl.PluginAddWithResponse(root, body, adminSession)
			}

			t.Run("valid manifest is accepted", func(t *testing.T) {
				tests.AssertRequest(addExternal(map[string]any{
					"id":              "valid-plugin",
					"name":            "Valid Plugin",
					"author":          "test-author",
					"description":     "A valid test plugin",
					"version":         "1.0.0",
					"command":         "./plugin",
					"events_consumed": []string{},
				}))(t, http.StatusOK)
			})

			t.Run("missing required name field is rejected", func(t *testing.T) {
				tests.AssertRequest(addExternal(map[string]any{
					"id":              "no-name-plugin",
					"author":          "test-author",
					"description":     "Missing name",
					"version":         "1.0.0",
					"command":         "./plugin",
					"events_consumed": []string{},
				}))(t, http.StatusBadRequest)
			})

			t.Run("invalid id pattern is rejected", func(t *testing.T) {
				tests.AssertRequest(addExternal(map[string]any{
					"id":              "invalid id with spaces",
					"name":            "Bad ID Plugin",
					"author":          "test-author",
					"description":     "Plugin with spaces in ID",
					"version":         "1.0.0",
					"command":         "./plugin",
					"events_consumed": []string{},
				}))(t, http.StatusBadRequest)
			})

			t.Run("invalid author pattern is rejected", func(t *testing.T) {
				tests.AssertRequest(addExternal(map[string]any{
					"id":              "bad-author-plugin",
					"name":            "Bad Author Plugin",
					"author":          "invalid author name",
					"description":     "Plugin with spaces in author",
					"version":         "1.0.0",
					"command":         "./plugin",
					"events_consumed": []string{},
				}))(t, http.StatusBadRequest)
			})

			t.Run("empty name is rejected", func(t *testing.T) {
				tests.AssertRequest(addExternal(map[string]any{
					"id":              "empty-name-plugin",
					"name":            "",
					"author":          "test-author",
					"description":     "Plugin with empty name",
					"version":         "1.0.0",
					"command":         "./plugin",
					"events_consumed": []string{},
				}))(t, http.StatusBadRequest)
			})
		}))
	}))
}
