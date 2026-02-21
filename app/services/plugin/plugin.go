package plugin

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/plugin/plugin_manager"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/plugin_host"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/plugin_logger"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/rpc_handler"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(plugin_manager.New),
		fx.Provide(plugin_host.New),
		fx.Provide(rpc_handler.NewFactory),
		plugin_logger.Build(),
	)
}
