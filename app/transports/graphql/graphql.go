package graphql

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/transports/graphql/server"
)

//go:generate go run github.com/99designs/gqlgen@latest

func Build() fx.Option {
	return fx.Options(
		fx.Provide(bindings),
		fx.Provide(newServer),
		fx.Invoke(mount),
	)
}

func bindings() graphql.ExecutableSchema {
	return server.NewExecutableSchema(server.Config{
		Resolvers: &Resolver{},
	})
}

func newServer(router *echo.Echo, es graphql.ExecutableSchema) *handler.Server {
	return handler.NewDefaultServer(es)
}

func mount(lc fx.Lifecycle, l *zap.Logger, router *echo.Echo, s *handler.Server) {
	// router.Use(echo.WrapMiddleware(func(h http.Handler) http.Handler {
	// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		if strings.HasPrefix(r.URL.Path, "/graphql") {
	// 			s.ServeHTTP(w, r)
	// 			return
	// 		}

	// 		h.ServeHTTP(w, r)
	// 	})
	// }))

	p := playground.Handler("Storyden", "/graphql/query")

	router.Any("/graphql/ui", func(c echo.Context) error {
		p.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	router.Any("/graphql/query", func(c echo.Context) error {
		s.ServeHTTP(c.Response(), c.Request())

		return nil
	})

	l.Info("mounted GraphQL to service bindings")
}
