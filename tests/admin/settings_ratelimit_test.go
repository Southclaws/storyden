package admin_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestRateLimitCostOverrides(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			t.Run("get_default_settings", func(t *testing.T) {
				r := require.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				resp, err := cl.AdminSettingsGetWithResponse(adminCtx, adminSession)
				tests.Ok(t, err, resp)

				settings := resp.JSON200
				r.NotNil(settings)

				// Settings exist but cost overrides should be nil or empty initially
				if settings.Services != nil && settings.Services.RateLimiting != nil && settings.Services.RateLimiting.CostOverrides != nil {
					r.Empty(*settings.Services.RateLimiting.CostOverrides)
				}
			})

			t.Run("set_rate_limit_overrides", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				// Set custom rate limit costs for specific operations
				costOverrides := map[string]int{
					"ThreadCreate": 5,
					"ThreadList":   2,
					"ReplyCreate":  3,
				}

				updateReq := openapi.AdminSettingsUpdateJSONRequestBody{
					Services: &openapi.AdminSettingsServiceProps{
						RateLimiting: &openapi.RateLimitServiceSettings{
							CostOverrides: &costOverrides,
						},
					},
				}

				resp, err := cl.AdminSettingsUpdateWithResponse(adminCtx, updateReq, adminSession)
				tests.Ok(t, err, resp)

				settings := resp.JSON200
				r.NotNil(settings)
				r.NotNil(settings.Services)
				r.NotNil(settings.Services.RateLimiting)
				r.NotNil(settings.Services.RateLimiting.CostOverrides)

				overrides := *settings.Services.RateLimiting.CostOverrides
				a.Equal(5, overrides["ThreadCreate"])
				a.Equal(2, overrides["ThreadList"])
				a.Equal(3, overrides["ReplyCreate"])
			})

			t.Run("get_updated_settings", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				resp, err := cl.AdminSettingsGetWithResponse(adminCtx, adminSession)
				tests.Ok(t, err, resp)

				settings := resp.JSON200
				r.NotNil(settings)
				r.NotNil(settings.Services)
				r.NotNil(settings.Services.RateLimiting)
				r.NotNil(settings.Services.RateLimiting.CostOverrides)

				overrides := *settings.Services.RateLimiting.CostOverrides
				a.Equal(5, overrides["ThreadCreate"])
				a.Equal(2, overrides["ThreadList"])
				a.Equal(3, overrides["ReplyCreate"])
			})

			t.Run("update_partial_overrides", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				// Update only one override, should merge with existing
				partialOverrides := map[string]int{
					"ThreadCreate": 10, // Update existing
					"PostUpdate":   7,  // Add new
				}

				updateReq := openapi.AdminSettingsUpdateJSONRequestBody{
					Services: &openapi.AdminSettingsServiceProps{
						RateLimiting: &openapi.RateLimitServiceSettings{
							CostOverrides: &partialOverrides,
						},
					},
				}

				resp, err := cl.AdminSettingsUpdateWithResponse(adminCtx, updateReq, adminSession)
				tests.Ok(t, err, resp)

				settings := resp.JSON200
				r.NotNil(settings)
				r.NotNil(settings.Services)
				r.NotNil(settings.Services.RateLimiting)
				r.NotNil(settings.Services.RateLimiting.CostOverrides)

				overrides := *settings.Services.RateLimiting.CostOverrides
				a.Equal(10, overrides["ThreadCreate"]) // Updated
				a.Equal(7, overrides["PostUpdate"])    // Added
			})

			t.Run("clear_overrides", func(t *testing.T) {
				r := require.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				emptyOverrides := map[string]int{}

				updateReq := openapi.AdminSettingsUpdateJSONRequestBody{
					Services: &openapi.AdminSettingsServiceProps{
						RateLimiting: &openapi.RateLimitServiceSettings{
							CostOverrides: &emptyOverrides,
						},
					},
				}

				resp, err := cl.AdminSettingsUpdateWithResponse(adminCtx, updateReq, adminSession)
				tests.Ok(t, err, resp)

				settings := resp.JSON200
				r.NotNil(settings)

				// Verify overrides are cleared
				if settings.Services != nil && settings.Services.RateLimiting != nil && settings.Services.RateLimiting.CostOverrides != nil {
					overrides := *settings.Services.RateLimiting.CostOverrides
					r.Empty(overrides)
				}
			})

			t.Run("set_guest_cost_multiplier", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				guestCost := 5

				updateReq := openapi.AdminSettingsUpdateJSONRequestBody{
					Services: &openapi.AdminSettingsServiceProps{
						RateLimiting: &openapi.RateLimitServiceSettings{
							RateLimitGuestCost: &guestCost,
						},
					},
				}

				resp, err := cl.AdminSettingsUpdateWithResponse(adminCtx, updateReq, adminSession)
				tests.Ok(t, err, resp)

				settings := resp.JSON200
				r.NotNil(settings)
				r.NotNil(settings.Services)
				r.NotNil(settings.Services.RateLimiting)
				r.NotNil(settings.Services.RateLimiting.RateLimitGuestCost)

				a.Equal(5, *settings.Services.RateLimiting.RateLimitGuestCost)
			})

			t.Run("get_guest_cost_setting", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				resp, err := cl.AdminSettingsGetWithResponse(adminCtx, adminSession)
				tests.Ok(t, err, resp)

				settings := resp.JSON200
				r.NotNil(settings)
				r.NotNil(settings.Services)
				r.NotNil(settings.Services.RateLimiting)
				r.NotNil(settings.Services.RateLimiting.RateLimitGuestCost)

				a.Equal(5, *settings.Services.RateLimiting.RateLimitGuestCost)
			})
		}))
	}))
}
