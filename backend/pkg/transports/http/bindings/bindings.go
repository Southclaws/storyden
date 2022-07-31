package bindings

import (
	"context"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/backend/internal/web"
	"github.com/Southclaws/storyden/backend/pkg/transports/http/openapi"
)

// Bindings is a DI parameter struct that is used to compose together all of the
// individual service bindings in this package. When the provider below depends
// on this type, it provides all these composed bindings to the DI system so the
// invoke call can mount them onto the router using the OpenAPI ServerInterface.
type Bindings struct {
	fx.In
	Version
	Authentication
}

func mountBindings(lc fx.Lifecycle, l *zap.Logger, router *echo.Echo, si openapi.StrictServerInterface) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			openapi.RegisterHandlers(router, openapi.NewStrictHandler(si, nil))
			router.GET("/openapi.json", spec)
			return nil
		},
	})

	l.Info("mounted OpenAPI to service bindings")
}

func addMiddleware(l *zap.Logger, router *echo.Echo, a Authentication) {
	router.Use(echo.WrapMiddleware(web.WithLogger))
	router.Use(echo.WrapMiddleware(a.middleware))

	l.Info("added router middleware")
}

func Build() fx.Option {
	return fx.Options(
		// Provide the bindings struct which implements the generated OpenAPI
		// interface by composing together all of the service bindings into a
		// single struct.
		fx.Provide(func(s Bindings) openapi.StrictServerInterface { return &s }),

		// Add the middleware bindings.
		fx.Invoke(addMiddleware),

		// Mount the bound OpenAPI routes onto the router.
		fx.Invoke(mountBindings),

		// Provide all service layer bindings to the DI system so they can be
		// depended upon during the binding provider above.
		fx.Provide(
			NewVersion,
			NewAuthentication,
		),
	)
}
