package origin

import (
	"net/http"
	"net/url"
	"strings"

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

	table := lo.SliceToMap(allowedOrigins(m.cfg), func(s string) (string, struct{}) { return s, struct{}{} })

	corsConfig := cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			_, present := table[normaliseOrigin(origin)]
			return present
		},
		AllowedMethods:   allowedMethods,
		AllowedHeaders:   allowedHeaders,
		ExposedHeaders:   exposedHeaders,
		AllowCredentials: true,
		MaxAge:           300,
	})

	return func(next http.Handler) http.Handler {
		handler := corsConfig.Handler(next)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := setOriginContext(r.Context(), r.Header.Get("Origin"))
			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func allowedOrigins(cfg config.Config) []string {
	origins := []string{
		originOf(cfg.PublicWebAddress.Scheme, cfg.PublicWebAddress.Host),
		originOf(cfg.PublicAPIAddress.Scheme, cfg.PublicAPIAddress.Host),
	}

	for _, raw := range cfg.CORSAllowedOrigins {
		origins = append(origins, normaliseOrigin(raw))
	}

	return lo.Filter(lo.Uniq(origins), func(s string, _ int) bool { return s != "" })
}

func normaliseOrigin(raw string) string {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return ""
	}

	return originOf(u.Scheme, u.Host)
}

func originOf(scheme, host string) string {
	if scheme == "" || host == "" {
		return ""
	}

	return strings.ToLower(scheme) + "://" + strings.ToLower(host)
}
