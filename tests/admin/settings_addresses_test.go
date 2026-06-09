package admin_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestAdminSettingsIncludesAddresses(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			t.Run("get_settings_includes_web_and_api_addresses", func(t *testing.T) {
				r := require.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				resp, err := cl.AdminSettingsGetWithResponse(adminCtx, adminSession)
				tests.Ok(t, err, resp)

				settings := resp.JSON200
				r.NotNil(settings)

				// Verify web and api addresses are present
				r.NotEmpty(settings.WebAddress, "web_address should be present in admin settings response")
				r.NotEmpty(settings.ApiAddress, "api_address should be present in admin settings response")

				// Verify they are valid URIs (basic check)
				r.Contains(settings.WebAddress, "http", "web_address should be a valid URI")
				r.Contains(settings.ApiAddress, "http", "api_address should be a valid URI")
			})

			t.Run("update_settings_returns_web_and_api_addresses", func(t *testing.T) {
				r := require.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				// Update some setting
				title := "Test Instance"
				updateBody := openapi.AdminSettingsUpdateJSONRequestBody{
					Title: &title,
				}

				resp, err := cl.AdminSettingsUpdateWithResponse(adminCtx, updateBody, adminSession)
				tests.Ok(t, err, resp)

				settings := resp.JSON200
				r.NotNil(settings)

				// Verify web and api addresses are present in update response
				r.NotEmpty(settings.WebAddress, "web_address should be present in update response")
				r.NotEmpty(settings.ApiAddress, "api_address should be present in update response")

				// Verify they are valid URIs (basic check)
				r.Contains(settings.WebAddress, "http", "web_address should be a valid URI")
				r.Contains(settings.ApiAddress, "http", "api_address should be a valid URI")
			})
		}))
	}))
}
