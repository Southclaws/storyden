package http

import (
	"net/http"

	"github.com/Southclaws/storyden/backend/internal/config"
	"github.com/Southclaws/storyden/backend/internal/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

func newRouter(l *zap.Logger, cfg config.Config) chi.Router {
	router := chi.NewRouter()

	origins := []string{
		"http://localhost:3000", // Local development
		cfg.PublicWebAddress,    // Live public website
	}

	l.Debug("preparing router", zap.Strings("origins", origins))

	router.Use(
		web.WithLogger,
		cors.Handler(cors.Options{
			AllowedOrigins:   origins,
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Content-Length", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link", "Content-Length", "X-Ratelimit-Limit", "X-Ratelimit-Reset"},
			AllowCredentials: true,
			MaxAge:           300,
		}),
	)

	router.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		web.Write(w, map[string]string{"version": config.Version}) //nolint:errcheck
	})

	router.HandleFunc(
		"/{rest:[a-zA-Z0-9=\\-\\/]+}",
		func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte("no module found for that route")); err != nil {
				zap.L().Warn("failed to write error", zap.Error(err))
			}
		})

	return router
}
