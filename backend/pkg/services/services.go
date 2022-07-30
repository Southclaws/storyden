package services

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/backend/pkg/services/authentication"
)

func Build() fx.Option {
	return fx.Options(
		authentication.Build(),
	)
}
