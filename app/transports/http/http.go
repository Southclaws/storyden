// Package http provides a transport layer using HTTP. In this layer, most of
// the code is generated: low level HTTP handlers, request and response structs
// and object validation.
//
// This is wired up to the service layer using "Bindings" which are just glue
// code which call service APIs from the endpoint handlers and deal with errors.
package http

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/transports/http/bindings"
	"github.com/Southclaws/storyden/app/transports/http/middleware"
	"github.com/Southclaws/storyden/internal/infrastructure/httpserver"
)

func Build() fx.Option {
	return fx.Options(
		// Builds the *http.ServeMux and *http.Server dependencies and, once set
		// up, starts the server on the fx OnStart lifecycle event.
		httpserver.Build(),

		// Build all middleware dependencies.
		middleware.Build(),

		// Binds all the generated spec code for services to the *http.ServeMux.
		bindings.Build(),

		fx.Invoke(MountOpenAPI),
	)
}
