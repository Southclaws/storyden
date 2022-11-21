package ogen

import (
	"context"
	"net/http"

	"github.com/Southclaws/fault"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/internal/openapi/ogen"
)

// Bindings is a DI parameter struct that is used to compose together all of the
// individual service bindings in this package. When the provider below depends
// on this type, it provides all these composed bindings to the DI system so the
// invoke call can mount them onto the router using the `StrictServerInterface`.
//
// The reason this is done this way is so we split code up based on OpenAPI
// REST collections instead of bundling everything into one huge struct with
// loads of dependencies. This is just how the oapi-codegen tool works, by
// generating one big interface which the bindings layer must satisfy.
type Bindings struct {
	fx.In

	Tmp
}

// bindingsProviders provides to the application the necessary implementations
// that compose the `Bindings` parameter struct which implements the OpenAPI
// server interface. When you add a new collection, add it to Bindings and here.
func bindingsProviders() fx.Option {
	return fx.Provide(
	//
	)
}

// bindings provides to the application the above struct which binds the service
// layer to the transport layer. This uses `Bindings` as an fx parameter struct.
//
// ## WHY AM I GETTING AN ERROR HERE?
//
// When you edit `openapi.yaml` and re-run the code generation task, this will
// most likely change the declaration of `StrictServerInterface` inside the
// generated package `openapi`.
//
// The error you will see is most likely something along the lines of:
//
//	*Bindings does not implement openapi.StrictServerInterface
//
// and the underlying problem is either missing methods or methods that have
// changed signature due to changes to the parameters or request or response.
//
// This API follows RESTful design so a collection in the API specification
// (such as `/v1/accounts`) will map to a file, struct and constructor here
// (such as `accounts.go`, `Accounts` and `NewAccounts`) and everything is glued
// together in this file.
func bindings(s Bindings) ogen.Handler {
	return &s
}

// mounts the OpenAPI routes and middleware onto the /api path. Everything that
// is outside of the `/api` path is considered part of the proxied frontend app.
// Note: routes are mounted with the `OnStart` hook so that middleware is first.
func mount(lc fx.Lifecycle, l *zap.Logger, mux *http.ServeMux, router *http.ServeMux, h ogen.Handler) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			s, err := ogen.NewServer(h)
			if err != nil {
				return fault.Wrap(err)
			}

			router.Handle("/api", s)

			l.Info("mounted OpenAPI to service bindings")

			// mount onto / because this router already only cares about /api
			mux.Handle("/", router)

			return nil
		},
	})
}

func Build() fx.Option {
	return fx.Options(
		// Provide the bindings struct which implements the generated OpenAPI
		// interface by composing together all of the service bindings into a
		// single struct.
		fx.Provide(bindings),

		// Mount the bound OpenAPI routes onto the router.
		fx.Invoke(mount),

		// Provide all service layer bindings to the DI system.
		bindingsProviders(),
	)
}
