package plugin_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestPluginsDisabledWhenProviderNone(t *testing.T) {
	integration.Test(t, &config.Config{
		PluginRuntimeProvider: "none",
	}, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		accountWrite *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)

			adminSession, _ := createAdminSession(root, cl, sh, accountWrite, r)

			infoResp := tests.AssertRequest(cl.GetInfoWithResponse(root))(t, http.StatusOK)
			r.NotContains(infoResp.JSON200.Capabilities, openapi.InstanceCapability("plugins"))

			tests.AssertRequest(
				cl.PluginListWithResponse(root, adminSession),
			)(t, http.StatusForbidden)

			body := openapi.PluginInitialProps{}
			err := body.FromPluginInitialExternal(openapi.PluginInitialExternal{
				Mode:     openapi.External,
				Manifest: buildTestManifest("Disabled Plugin", nil),
			})
			r.NoError(err)

			tests.AssertRequest(
				cl.PluginAddWithResponse(root, body, adminSession),
			)(t, http.StatusForbidden)
		}))
	}))
}
