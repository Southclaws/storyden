package swagger

import (
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/api/src/config"
)

func Build() fx.Option {
	return fx.Options(
		fx.Invoke(func(
			cfg config.Config,
			r chi.Router,
		) {
			if cfg.Production {
				return
			}

			rtr := chi.NewRouter()
			r.Mount("/swagger", rtr)

			rtr.Get("/*", httpSwagger.Handler(
				httpSwagger.URL(cfg.ListenAddr+"/swagger/doc.json"),
			))
		}),
	)
}
