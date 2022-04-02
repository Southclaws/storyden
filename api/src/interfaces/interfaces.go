package interfaces

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/api/src/interfaces/api"
)

func Build() fx.Option {
	return fx.Options(
		api.Build(),
	)
}
