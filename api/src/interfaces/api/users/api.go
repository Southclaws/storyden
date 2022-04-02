package users

import (
	"github.com/go-chi/chi"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/api/src/resources/user"
	"github.com/Southclaws/storyden/api/src/services/authentication"
)

type controller struct {
	auth *authentication.State
	repo user.Repository
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(func(
			auth *authentication.State,
			repo user.Repository,
		) *controller {
			return &controller{auth, repo}
		}),
		fx.Invoke(func(
			r chi.Router,
			c *controller,
			auth *authentication.State,
		) {
			rtr := chi.NewRouter()
			r.Mount("/users", rtr)

			rtr.
				With(authentication.MustBeAuthenticated, auth.MustBeAdmin).
				Get("/", c.list)

			rtr.
				Get("/{id}", c.get)

			rtr.
				Get("/{id}/image", c.image)

			rtr.
				With(authentication.MustBeAuthenticated).
				Patch("/{id}", c.patch)

			rtr.
				With(authentication.MustBeAuthenticated).
				Get("/self", c.self)

			rtr.
				With(authentication.MustBeAuthenticated, auth.MustBeAdmin).
				Patch("/{id}/banstatus", c.banstatus)

			rtr.
				With(authentication.MustBeAuthenticated, auth.MustBeAdmin).
				Patch("/{id}/adminstatus", c.patchAdmin)

			rtr.Get("/dev", c.dev)
		}),
	)
}
