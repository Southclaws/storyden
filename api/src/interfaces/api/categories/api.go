package categories

import (
	"github.com/go-chi/chi"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/api/src/resources/category"
	"github.com/Southclaws/storyden/api/src/services/authentication"
)

type controller struct {
	as   *authentication.State
	repo category.Repository
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(func(as *authentication.State, repo category.Repository) *controller { return &controller{as, repo} }),
		fx.Invoke(func(
			r chi.Router,
			c *controller,

		) {
			rtr := chi.NewRouter()
			r.Mount("/categories", rtr)

			rtr.
				Get("/", c.get)

			rtr.
				With(authentication.MustBeAuthenticated, c.as.MustBeAdmin).
				Patch("/", c.patchAll)

			rtr.
				With(authentication.MustBeAuthenticated, c.as.MustBeAdmin).
				Delete("/{id}", c.delete)

			rtr.
				With(authentication.MustBeAuthenticated, c.as.MustBeAdmin).
				Patch("/{id}", c.patch)
		}),
	)
}
