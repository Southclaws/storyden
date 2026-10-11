package admin_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestDiscoverySettingsProjection(t *testing.T) {
	if tests.IsSharedPostgresDatabase() {
		t.Skip("skipping global discovery settings mutations on shared postgres database")
	}

	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		repo *settings.SettingsRepository,
		db *ent.Client,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			expected := openapi.DiscoveryServiceSettings{
				SearchIndexing:   new(true),
				TrainingCrawling: new(false),
				LlmsTxt:          new(false),
				AuthMd:           new(false),
				AiCatalog:        new(false),
				ApiCatalog:       new(false),
				McpServerCard:    new(false),
				RobotsTxt:        new(""),
			}

			t.Run("defaults_are_opt_in", func(t *testing.T) {
				admin := tests.AssertRequest(cl.AdminSettingsGetWithResponse(adminCtx, adminSession))(t, http.StatusOK)
				require.NotNil(t, admin.JSON200.Services)
				require.Equal(t, &expected, admin.JSON200.Services.Discovery)

				public := tests.AssertRequest(cl.GetInfoWithResponse(root))(t, http.StatusOK)
				require.Equal(t, &expected, public.JSON200.Discovery)
			})

			t.Run("save_and_read_back_public_projection", func(t *testing.T) {
				expected.SearchIndexing = new(false)
				expected.TrainingCrawling = new(true)
				expected.LlmsTxt = new(true)
				expected.AuthMd = new(true)
				expected.AiCatalog = new(true)
				expected.ApiCatalog = new(true)
				expected.McpServerCard = new(true)
				expected.RobotsTxt = new("User-agent: *\nDisallow: /handwritten/\n")

				patch := openapi.AdminSettingsUpdateJSONRequestBody{
					Services: &openapi.AdminSettingsServiceProps{
						Discovery: &expected,
						Moderation: &openapi.ModerationServiceSettings{
							ThreadBodyLengthMax: new(12345),
						},
					},
				}
				updated := tests.AssertRequest(cl.AdminSettingsUpdateWithResponse(adminCtx, patch, adminSession))(t, http.StatusOK)
				require.Equal(t, &expected, updated.JSON200.Services.Discovery)

				admin := tests.AssertRequest(cl.AdminSettingsGetWithResponse(adminCtx, adminSession))(t, http.StatusOK)
				require.Equal(t, &expected, admin.JSON200.Services.Discovery)

				public := tests.AssertRequest(cl.GetInfoWithResponse(root))(t, http.StatusOK)
				require.Equal(t, &expected, public.JSON200.Discovery)
				var body map[string]json.RawMessage
				require.NoError(t, json.Unmarshal(public.Body, &body))
				require.NotContains(t, body, "services")
				require.NotContains(t, string(public.Body), "thread_body_length_max")
			})

			t.Run("partial_update_preserves_other_fields_and_services", func(t *testing.T) {
				patch := openapi.AdminSettingsUpdateJSONRequestBody{
					Services: &openapi.AdminSettingsServiceProps{
						Discovery: &openapi.DiscoveryServiceSettings{
							LlmsTxt:   new(false),
							RobotsTxt: new(""),
						},
					},
				}
				expected.LlmsTxt = new(false)
				expected.RobotsTxt = new("")

				updated := tests.AssertRequest(cl.AdminSettingsUpdateWithResponse(adminCtx, patch, adminSession))(t, http.StatusOK)
				require.Equal(t, &expected, updated.JSON200.Services.Discovery)
				require.Equal(t, 12345, *updated.JSON200.Services.Moderation.ThreadBodyLengthMax)

				public := tests.AssertRequest(cl.GetInfoWithResponse(root))(t, http.StatusOK)
				require.Equal(t, &expected, public.JSON200.Discovery)
			})

			t.Run("existing_install_without_discovery_settings", func(t *testing.T) {
				stored, err := db.Setting.Get(root, settings.StorydenPrimarySettingsKey)
				require.NoError(t, err)
				var legacy settings.Settings
				require.NoError(t, json.Unmarshal([]byte(stored.Value), &legacy))
				services := legacy.Services.OrZero()
				services.Discovery = nil
				legacy.Services = opt.New(services)
				value, err := json.Marshal(legacy)
				require.NoError(t, err)
				require.NoError(t, db.Setting.UpdateOneID(settings.StorydenPrimarySettingsKey).SetValue(string(value)).Exec(root))

				// Set reads the persisted snapshot and refreshes the repository cache.
				_, err = repo.Set(root, settings.Settings{})
				require.NoError(t, err)

				admin := tests.AssertRequest(cl.AdminSettingsGetWithResponse(adminCtx, adminSession))(t, http.StatusOK)
				require.Nil(t, admin.JSON200.Services.Discovery)
				public := tests.AssertRequest(cl.GetInfoWithResponse(root))(t, http.StatusOK)
				require.Nil(t, public.JSON200.Discovery)
			})
		}))
	}))
}
