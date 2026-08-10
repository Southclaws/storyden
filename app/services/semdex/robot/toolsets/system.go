package toolsets

import (
	"go.uber.org/fx"

	robottools "github.com/Southclaws/storyden/app/services/semdex/robot/tools"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(NewRegistry),
		fx.Provide(func(registry *Registry) robottools.ToolsetResolver { return registry }),
		fx.Invoke(registerSystemToolsets),
	)
}

func registerSystemToolsets(registry *Registry) error {
	for _, definition := range systemDefinitions() {
		if err := registry.Register(definition); err != nil {
			return err
		}
	}
	return nil
}
