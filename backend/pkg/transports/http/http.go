package http

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/backend/pkg/transports/http/bindings"
)

//go:generate go run -mod=mod github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config ./api/config.yaml ./api/openapi.yaml

func Build() fx.Option {
	return fx.Options(
		bindings.Build(),
		fx.Provide(newRouter),
		fx.Invoke(newServer),
	)
}
