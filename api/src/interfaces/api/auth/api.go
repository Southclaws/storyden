package auth

import (
	"github.com/go-chi/chi"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/api/src/services/authentication"
)

type controller struct{}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(func() *controller { return &controller{} }),
		fx.Invoke(func(
			r chi.Router,
			c *controller,
			auth *authentication.State,
		) {
			rtr := chi.NewRouter()
			r.Mount("/auth", rtr)

			// TODO: Auth endpoints
		}),
	)
}
