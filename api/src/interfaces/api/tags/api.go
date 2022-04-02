package tags

import (
	"github.com/go-chi/chi"
	"github.com/Southclaws/storyden/api/src/resources/tag"
	"go.uber.org/fx"
)

type controller struct {
	tags tag.Repository
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(func(tags tag.Repository) *controller {
			return &controller{tags}
		}),
		fx.Invoke(func(
			r chi.Router,
			c *controller,

		) {
			rtr := chi.NewRouter()
			r.Mount("/tags", rtr)

			rtr.
				Get("/", c.get)
		}),
	)
}
