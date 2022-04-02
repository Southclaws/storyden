package posts

import (
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/api/src/config"
	"github.com/Southclaws/storyden/api/src/infra/web/ratelimiter"
	"github.com/Southclaws/storyden/api/src/resources/notification"
	"github.com/Southclaws/storyden/api/src/resources/post"
	"github.com/Southclaws/storyden/api/src/services/authentication"
)

type controller struct {
	as            *authentication.State
	repo          post.Repository
	publicAddress string

	// TODO: This should be event-driven so forum posts result in an event being
	// emitted and a worker reacts to the event to create subscriptions.
	notifications notification.Repository
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(func(
			as *authentication.State,
			repo post.Repository,
			notifications notification.Repository,
			cfg config.Config,
		) *controller {
			return &controller{as, repo, cfg.PublicWebAddress, notifications}
		}),
		fx.Invoke(func(
			r chi.Router,
			c *controller,
		) {
			rtr := chi.NewRouter()
			r.Mount("/posts", rtr)

			rtr.
				Get("/{slug}", c.get)

			rtr.
				With(authentication.MustBeAuthenticated, ratelimiter.WithRateLimit(3, time.Minute)).
				Post("/{id}", c.post)

			rtr.
				With(authentication.MustBeAuthenticated, ratelimiter.WithRateLimit(3, time.Minute)).
				Patch("/{id}", c.patch)

			rtr.
				With(authentication.MustBeAuthenticated, ratelimiter.WithRateLimit(3, time.Minute)).
				Delete("/{id}", c.delete)
		}),
	)
}
