package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/transports/graphql/server"
)

//go:generate go run github.com/99designs/gqlgen@v0.17.20

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
	srv := handler.NewDefaultServer(es)
	srv.Use(extension.Introspection{})

	return srv
}

func mount(lc fx.Lifecycle, l *zap.Logger, mux *http.ServeMux, s *handler.Server) {
	p := playground.Handler("Storyden", "/graphql/query")

	mux.Handle("/graphql/ui", p)
	mux.Handle("/graphql/query", s)

	l.Info("mounted GraphQL to service bindings")
}
