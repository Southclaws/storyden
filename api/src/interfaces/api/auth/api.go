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
			auth *authentication.CookieAuth,
		) {
			rtr := chi.NewRouter()
			r.Mount("/auth", rtr)

			// Return a list of auth methods
			// rtr.Get("/methods", c.methods)

			// Initiate an auth flow.
			rtr.Post("/initiate/{method}", c.start)

			// Finish auth flow.
			rtr.Post("/callback/{method}", c.callback)
		}),
	)
}
