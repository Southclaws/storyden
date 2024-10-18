package limiter

import (
	"context"
	"net/http"

	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/httpserver/ratelimit"
)

func WithRateLimiter(rl *ratelimit.Limiter) func(next http.Handler) http.Handler {
	limiterForSession := func(ctx context.Context) func(next http.Handler) http.Handler {
		s := session.GetOptAccountID(ctx)
		isAuthenticated := s.Ok()

		if isAuthenticated {
			return rl.Member.Handler
		}

		return rl.Guest.Handler
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mw := limiterForSession(r.Context())
			mw(next).ServeHTTP(w, r)
		})
	}
}

func WithRequestSizeLimiter(bytes int64) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, bytes)
			h.ServeHTTP(w, r)
		})
	}
}
