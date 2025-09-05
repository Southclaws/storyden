package plugin

import (
	"github.com/Southclaws/storyden/app/services/plugin/plugin_manager"
	"go.uber.org/fx"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			plugin_manager.New,
		),
	)
}
