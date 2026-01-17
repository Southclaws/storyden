package origin

import (
	"net/http"

	"github.com/rs/cors"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/internal/config"
)

type Middleware struct {
	cfg config.Config
}

func New(cfg config.Config) *Middleware {
	return &Middleware{cfg: cfg}
}

func (m *Middleware) WithCORS() func(next http.Handler) http.Handler {
	allowedMethods := []string{
		"GET",
		"POST",
		"PUT",
		"PATCH",
		"DELETE",
		"OPTIONS",
	}

	allowedHeaders := []string{
		"Accept",
		"Authorization",
		"Content-Type",
		"Content-Length",
		"X-CSRF-Token",
		"X-Correlation-ID",
		"X-Forwarded-Host",
	}

	exposedHeaders := []string{
		"Link",
		"Content-Type",
		"Content-Length",
		"X-Ratelimit-Limit",
		"X-Ratelimit-Remaining",
		"X-Ratelimit-Reset",
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// TODO: Provide a way to set multiple allowed origins via config.
			// NOTE: Currently, we allow all origins but not via "*".
			allowedOrigins := []string{
				"http://localhost:3000",
				origin,
			}

			table := lo.SliceToMap(allowedOrigins, func(s string) (string, struct{}) { return s, struct{}{} })

			allowOriginFunc := func(origin string) bool {
				_, present := table[origin]
				return present
			}

			corsConfig := cors.New(cors.Options{
				AllowOriginFunc:  allowOriginFunc,
				AllowedMethods:   allowedMethods,
				AllowedHeaders:   allowedHeaders,
				ExposedHeaders:   exposedHeaders,
				AllowCredentials: true,
				MaxAge:           300,
			})

			ctx := setOriginContext(r.Context(), origin)

			corsConfig.Handler(next).ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
