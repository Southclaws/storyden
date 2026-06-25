package robot

import (
	"context"

	"github.com/Southclaws/opt"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/settings"
)

func WithRobotSettings(model string) fx.Option {
	return fx.Invoke(func(lc fx.Lifecycle, root context.Context, repo *settings.SettingsRepository) {
		lc.Append(fx.StartHook(func() error {
			_, err := repo.Set(root, settings.Settings{
				Services: opt.New(settings.ServiceSettings{
					Robots: opt.New(settings.RobotServiceSettings{
						Enabled:      opt.New(true),
						DefaultModel: opt.New(model),
						Providers: opt.New(map[string]settings.RobotProviderSettings{
							"mock": {
								Enabled: opt.New(true),
							},
						}),
					}),
				}),
			})
			return err
		}))
	})
}
