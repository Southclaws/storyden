package settings_test

import (
	"testing"

	"github.com/Southclaws/opt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/settings"
)

func TestSettingsMerge(t *testing.T) {
	t.Parallel()

	t.Run("top_level_partial_update", func(t *testing.T) {
		r := require.New(t)
		a := assert.New(t)

		old := settings.Settings{
			Title:       opt.New("Old Title"),
			Description: opt.New("untouched description"),
		}

		updated := settings.Settings{
			Title: opt.New("New Title"),
		}

		err := old.Merge(updated)
		r.NoError(err)

		a.Equal("New Title", old.Title.OrZero())
		a.Equal("untouched description", old.Description.OrZero())
	})

	t.Run("service_patch_preserves_robot_providers_when_default_model_changes", func(t *testing.T) {
		r := require.New(t)
		a := assert.New(t)

		old := settings.Settings{
			Services: opt.New(settings.ServiceSettings{
				Robots: opt.New(settings.RobotServiceSettings{
					Enabled: opt.New(true),
					Providers: opt.New(map[string]settings.RobotProviderSettings{
						"anthropic": {
							Enabled: opt.New(true),
							APIKey:  opt.New("sk-ant-existing"),
						},
					}),
				}),
			}),
		}

		updated := settings.Settings{
			Services: opt.New(settings.ServiceSettings{
				Robots: opt.New(settings.RobotServiceSettings{
					DefaultModel: opt.New("anthropic/claude-haiku-4-5-20251001"),
				}),
			}),
		}

		err := old.Merge(updated)
		r.NoError(err)

		robots, ok := old.Services.OrZero().Robots.Get()
		r.True(ok)
		providers, ok := robots.Providers.Get()
		r.True(ok)

		a.True(robots.Enabled.Or(false))
		a.Equal("anthropic/claude-haiku-4-5-20251001", robots.DefaultModel.OrZero())
		a.Equal("sk-ant-existing", providers["anthropic"].APIKey.OrZero())
		a.True(providers["anthropic"].Enabled.Or(false))
	})

	t.Run("service_patch_preserves_sibling_service_blocks", func(t *testing.T) {
		r := require.New(t)
		a := assert.New(t)

		old := settings.Settings{
			Services: opt.New(settings.ServiceSettings{
				RateLimit: opt.New(settings.RateLimitServiceSettings{
					RateLimit: opt.New(100),
				}),
				Robots: opt.New(settings.RobotServiceSettings{
					DefaultModel: opt.New("mock/model"),
				}),
			}),
		}

		updated := settings.Settings{
			Services: opt.New(settings.ServiceSettings{
				ClientIP: opt.New(settings.ClientIPServiceSettings{
					ClientIPHeader: opt.New("X-Real-IP"),
				}),
			}),
		}

		err := old.Merge(updated)
		r.NoError(err)

		services := old.Services.OrZero()
		a.Equal(100, services.RateLimit.OrZero().RateLimit.OrZero())
		a.Equal("mock/model", services.Robots.OrZero().DefaultModel.OrZero())
		a.Equal("X-Real-IP", services.ClientIP.OrZero().ClientIPHeader.OrZero())
	})

	t.Run("robot_provider_patch_preserves_existing_api_key", func(t *testing.T) {
		r := require.New(t)
		a := assert.New(t)

		old := settings.Settings{
			Services: opt.New(settings.ServiceSettings{
				Robots: opt.New(settings.RobotServiceSettings{
					Providers: opt.New(map[string]settings.RobotProviderSettings{
						"anthropic": {
							Enabled: opt.New(false),
							APIKey:  opt.New("sk-ant-existing"),
						},
					}),
				}),
			}),
		}

		updated := settings.Settings{
			Services: opt.New(settings.ServiceSettings{
				Robots: opt.New(settings.RobotServiceSettings{
					Providers: opt.New(map[string]settings.RobotProviderSettings{
						"anthropic": {
							Enabled: opt.New(true),
						},
					}),
				}),
			}),
		}

		err := old.Merge(updated)
		r.NoError(err)

		provider := old.Services.OrZero().Robots.OrZero().Providers.OrZero()["anthropic"]
		a.True(provider.Enabled.Or(false))
		a.Equal("sk-ant-existing", provider.APIKey.OrZero())
	})

	t.Run("empty_map_patch_clears_existing_map", func(t *testing.T) {
		r := require.New(t)
		a := assert.New(t)

		old := settings.Settings{
			Services: opt.New(settings.ServiceSettings{
				RateLimit: opt.New(settings.RateLimitServiceSettings{
					CostOverrides: opt.New(map[string]int{
						"ThreadCreate": 10,
						"ReplyCreate":  3,
					}),
				}),
			}),
		}

		updated := settings.Settings{
			Services: opt.New(settings.ServiceSettings{
				RateLimit: opt.New(settings.RateLimitServiceSettings{
					CostOverrides: opt.New(map[string]int{}),
				}),
			}),
		}

		err := old.Merge(updated)
		r.NoError(err)

		overrides, ok := old.Services.OrZero().RateLimit.OrZero().CostOverrides.Get()
		r.True(ok)
		a.Empty(overrides)
	})
}
