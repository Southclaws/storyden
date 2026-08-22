package trail

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/trail/trail_action_registry"
	"github.com/Southclaws/storyden/app/services/trail/trail_manager"
	"github.com/Southclaws/storyden/app/services/trail/trail_runtime"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			trail_action_registry.NewRobotActionAdapter,
			trail_action_registry.New,
			trail_manager.New,
			trail_runtime.New,
		),
		fx.Invoke(func(*trail_runtime.Runtime) {}),
	)
}
