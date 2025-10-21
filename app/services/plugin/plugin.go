package plugin

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/plugin/plugin_manager"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			plugin_manager.New,
			plugin_runner.New,
		),
	)
}
