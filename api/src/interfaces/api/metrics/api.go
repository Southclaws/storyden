package metrics

import (
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
)

func Build() fx.Option {
	return fx.Options(
		fx.Invoke(func(
			r chi.Router,
		) {
			rtr := chi.NewRouter()
			r.Mount("/metrics", rtr)

			rtr.Handle("/", promhttp.Handler())
		}),
	)
}
