package threads

import (
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/api/src/infra/web/ratelimiter"
	"github.com/Southclaws/storyden/api/src/resources/notification"
	"github.com/Southclaws/storyden/api/src/resources/thread"
	"github.com/Southclaws/storyden/api/src/services/authentication"
)

type controller struct {
	as      *authentication.State
	threads thread.Repository

	// TODO: This should be event-driven so forum posts result in an event being
	// emitted and a worker reacts to the event to create subscriptions.
	notifications notification.Repository
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(func(
			as *authentication.State,
			threads thread.Repository,
			notifications notification.Repository,
		) *controller {
			return &controller{as, threads, notifications}
		}),
		fx.Invoke(func(
			r chi.Router,
			c *controller,
		) {
			rtr := chi.NewRouter()
			r.Mount("/threads", rtr)

			rtr.
				Get("/", c.list)

			rtr.
				With(authentication.MustBeAuthenticated, ratelimiter.WithRateLimit(3, time.Minute)).
				Post("/", c.post)

			rtr.
				With(authentication.MustBeAuthenticated, ratelimiter.WithRateLimit(3, time.Minute)).
				Patch("/{id}", c.patch)

			rtr.
				With(authentication.MustBeAuthenticated, ratelimiter.WithRateLimit(3, time.Minute)).
				Delete("/{id}", c.delete)
		}),
	)
}
