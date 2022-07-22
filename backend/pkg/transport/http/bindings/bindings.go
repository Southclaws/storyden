package bindings

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/backend/pkg/transport/http/openapi"
)

// Bindings is a DI parameter struct that is used to compose together all of the
// individual service bindings in this package. When the provider below depends
// on this type, it provides all these composed bindings to the DI system so the
// invoke call can mount them onto the router using the OpenAPI ServerInterface.
type Bindings struct {
	fx.In
	Authentication
}

func Build() fx.Option {
	return fx.Options(
		// Provide the bindings struct which implements the generated OpenAPI
		// interface by composing together all of the service bindings into a
		// single struct.
		fx.Provide(func(s Bindings) openapi.ServerInterface { return &s }),

		// Mount the bound OpenAPI routes onto the router.
		fx.Invoke(func(router chi.Router, si openapi.ServerInterface) { openapi.HandlerFromMux(si, router) }),

		// Provide all service layer bindings to the DI system so they can be
		// depended upon during the binding provider above.
		fx.Provide(
			NewAuthentication,
		),
	)
}
