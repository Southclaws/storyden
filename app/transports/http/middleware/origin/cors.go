package origin

import (
	"net/http"

	"github.com/rs/cors"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/internal/config"
)

func WithCORS(cfg config.Config) func(next http.Handler) http.Handler {
	allowedOrigins := []string{
		"http://localhost:3000", // Local development
		"http://localhost:8001", // Swagger UI
		cfg.PublicWebAddress,    // Live public website
	}

	table := lo.SliceToMap(allowedOrigins, func(s string) (string, struct{}) { return s, struct{}{} })

	allowOriginFunc := func(origin string) bool {
		_, present := table[origin]
		return present
	}

	cors := cors.New(cors.Options{
		AllowOriginFunc: allowOriginFunc,
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"Content-Length",
			"X-CSRF-Token",
			"X-Correlation-ID",
			"X-Forwarded-Host",
		},
		ExposedHeaders: []string{
			"Link",
			"Content-Type",
			"Content-Length",
			"X-Ratelimit-Limit",
			"X-Ratelimit-Reset",
		},
		AllowCredentials: true,
		MaxAge:           300,
	})

	return cors.Handler
}
