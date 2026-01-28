package plugin

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/plugin/plugin_manager"
	sdx_runner "github.com/Southclaws/storyden/app/services/plugin/plugin_runner/sdx_runner"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(sdx_runner.New),
		fx.Provide(plugin_manager.New),
	)
}
