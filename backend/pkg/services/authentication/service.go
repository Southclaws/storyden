package authentication

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/backend/pkg/services/authentication/provider"
)

func Build() fx.Option {
	return fx.Options(
		provider.Build(),
	)
}
