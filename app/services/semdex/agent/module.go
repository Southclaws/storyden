package agent

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/semdex/agent/tools"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(New),
		tools.Build(),
	)
}
