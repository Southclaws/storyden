package test

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/services/authentication"
)

type controller struct{}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(func() *controller { return &controller{} }),
		fx.Invoke(func(r chi.Router, c *controller, auth *authentication.State) {
			rtr := chi.NewRouter()
			r.Mount("/test", rtr)

			rtr.Get("/error", func(w http.ResponseWriter, r *http.Request) {
				web.StatusInternalServerError(w, web.WithSuggestion(
					errors.New("failed to exist"),
					"A problem occurred during the request and the process had to be cancelled. Your card was not charged.",
					"Try logging out and back in again, if this doesn't work, please contact our support team."))
			})
		}),
	)
}
