package admin_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestMediaSettings(t *testing.T) {
	if tests.IsSharedPostgresDatabase() {
		t.Skip("skipping global media settings mutations on shared postgres database")
	}

	t.Parallel()

	cfg := &config.Config{
		ImageWorkers:      1,
		ImageMemoryBudget: 512 << 20,
	}
	integration.Test(t, cfg, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			ctx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			auth := sh.WithSession(ctx)

			response := tests.AssertRequest(cl.AdminSettingsGetWithResponse(ctx, auth))(t, http.StatusOK)
			media := response.JSON200.Services.Media
			require.Equal(t, cfg.ImageMemoryBudget, *media.ImageMemoryBudget)
			require.Greater(t, *media.ImageEstimatedMemory, *media.ImageMemoryBudget)

			patch := openapi.AdminSettingsUpdateJSONRequestBody{
				Services: &openapi.AdminSettingsServiceProps{
					Media: &openapi.MediaServiceSettings{
						ImageMaxSize:   new(int64(1048576)),
						ImageMaxPixels: new(int64(250000)),
					},
				},
			}
			updated := tests.AssertRequest(cl.AdminSettingsUpdateWithResponse(ctx, patch, auth))(t, http.StatusOK)
			media = updated.JSON200.Services.Media
			require.Less(t, *media.ImageEstimatedMemory, *media.ImageMemoryBudget)
			require.Equal(t, response.JSON200.Services.Moderation, updated.JSON200.Services.Moderation)

			patch.Services.Media = &openapi.MediaServiceSettings{ImageMaxPixels: new(int64(2000000))}
			updated = tests.AssertRequest(cl.AdminSettingsUpdateWithResponse(ctx, patch, auth))(t, http.StatusOK)
			require.EqualValues(t, 1048576, *updated.JSON200.Services.Media.ImageMaxSize)
			require.EqualValues(t, 2000000, *updated.JSON200.Services.Media.ImageMaxPixels)

			for _, invalid := range []int64{0, -1} {
				patch.Services.Media = &openapi.MediaServiceSettings{ImageMaxSize: &invalid}
				tests.AssertRequest(cl.AdminSettingsUpdateWithResponse(ctx, patch, auth))(t, http.StatusBadRequest)
				patch.Services.Media = &openapi.MediaServiceSettings{ImageMaxPixels: &invalid}
				tests.AssertRequest(cl.AdminSettingsUpdateWithResponse(ctx, patch, auth))(t, http.StatusBadRequest)
			}

			stored := tests.AssertRequest(cl.AdminSettingsGetWithResponse(ctx, auth))(t, http.StatusOK)
			require.EqualValues(t, 1048576, *stored.JSON200.Services.Media.ImageMaxSize)
			require.EqualValues(t, 2000000, *stored.JSON200.Services.Media.ImageMaxPixels)
		}))
	}))
}
