package notifications

import (
	"github.com/go-chi/chi"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/api/src/resources/notification"
	"github.com/Southclaws/storyden/api/src/services/authentication"
)

type controller struct {
	repo notification.Repository
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(func(repo notification.Repository) *controller { return &controller{repo} }),
		fx.Invoke(func(r chi.Router, c *controller) {
			rtr := chi.NewRouter()
			r.Mount("/subscriptions/notifications", rtr)

			rtr.
				With(authentication.MustBeAuthenticated).
				Get("/", c.get)

			rtr.
				With(authentication.MustBeAuthenticated).
				Patch("/{id}", c.patch)

			rtr.
				With(authentication.MustBeAuthenticated).
				Delete("/{id}", c.delete)
		}),
	)
}
