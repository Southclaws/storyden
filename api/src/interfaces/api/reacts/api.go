package reacts

import (
	"github.com/go-chi/chi"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/api/src/resources/react"
	"github.com/Southclaws/storyden/api/src/services/authentication"
)

type controller struct {
	reacts react.Repository
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(func(reacts react.Repository) *controller {
			return &controller{reacts}
		}),
		fx.Invoke(func(
			r chi.Router,
			c *controller,

		) {
			rtr := chi.NewRouter()
			r.Mount("/reacts", rtr)

			rtr.
				With(authentication.MustBeAuthenticated).
				Post("/", c.post)

			rtr.
				With(authentication.MustBeAuthenticated).
				Delete("/{react_id}", c.delete)
		}),
	)
}
