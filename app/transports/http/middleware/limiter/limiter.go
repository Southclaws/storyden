package limiter

import (
	"context"
	"net/http"
	"time"

	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

type Limiter struct {
	Guest  *stdlib.Middleware
	Member *stdlib.Middleware
}

func WithRateLimiter(cfg config.Config) func(next http.Handler) http.Handler {
	store := memory.NewStore()

	// Rate limit for unauthenticated requests.
	guestRate := limiter.Rate{
		Period: time.Minute,
		Limit:  int64(cfg.UnauthenticatedRPM),
	}
	guest := stdlib.NewMiddleware(limiter.New(store, guestRate))

	// Rate limit for authenticated requests.
	memberRate := limiter.Rate{
		Period: time.Minute,
		Limit:  int64(cfg.AuthenticatedRPM),
	}
	member := stdlib.NewMiddleware(limiter.New(store, memberRate))

	limiterForSession := func(ctx context.Context) func(next http.Handler) http.Handler {
		s := session.GetOptAccountID(ctx)
		isAuthenticated := s.Ok()

		if isAuthenticated {
			return member.Handler
		}

		return guest.Handler
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
