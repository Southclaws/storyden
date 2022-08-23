// Package http provides a transport layer using HTTP. In this layer, most of
// the code is generated: low level HTTP handlers, request and response structs
// and object validation.
//
// This is wired up to the service layer using "Bindings" which are just glue
// code which call service APIs from the endpoint handlers and deal with errors.
//
// Within this package, all you will find is fx providers for the actual HTTP
// server, the router and the bindings. The bindings package is where most of
// the logic is and it depends on the router to mount the bindings.
//
package http

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/pkg/transports/http/bindings"
)

//go:generate go run -mod=mod github.com/deepmap/oapi-codegen/cmd/oapi-codegen@master --config ./api/config.yaml ./api/openapi.yaml

func Build() fx.Option {
	return fx.Options(
		bindings.Build(),
		fx.Provide(newRouter),
		fx.Invoke(newServer),
	)
}
