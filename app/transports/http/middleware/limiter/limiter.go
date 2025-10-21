package limiter

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/rate"
)

const (
	RateLimitLimit      = "X-RateLimit-Limit"
	RateLimitRemaining  = "X-RateLimit-Remaining"
	RateLimitReset      = "X-RateLimit-Reset"
	RetryAfter          = "Retry-After"
	MaxRequestSizeBytes = 50 * 1024 * 1024
)

type Middleware struct {
	rl        rate.Limiter
	kf        KeyFunc
	sizeLimit int64
}

func New(
	cfg config.Config,

	f *rate.LimiterFactory,
) *Middleware {
	rl := f.NewLimiter(cfg.RateLimit, cfg.RateLimitPeriod, cfg.RateLimitExpire)

	return &Middleware{
		rl:        rl,
		kf:        fromIP("CF-Connecting-IP", "X-Real-IP", "True-Client-IP"),
		sizeLimit: MaxRequestSizeBytes, // TODO: cfg.MaxRequestSize
	}
}

func (m *Middleware) WithRateLimit() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			key, err := m.kf(r)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			// TODO: Generate costs per-operation from OpenAPI spec
			cost := 1

			status, allowed, err := m.rl.Increment(ctx, key, cost)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			limit := status.Limit
			remaining := status.Remaining
			resetTime := status.Reset.UTC().Format(time.RFC1123)

			w.Header().Set(RateLimitLimit, strconv.FormatUint(uint64(limit), 10))
			w.Header().Set(RateLimitRemaining, strconv.FormatUint(uint64(remaining), 10))
			w.Header().Set(RateLimitReset, resetTime)

			if !allowed {
				w.Header().Set(RetryAfter, resetTime)
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

type KeyFunc func(r *http.Request) (string, error)

func fromIP(headers ...string) KeyFunc {
	return func(r *http.Request) (string, error) {
		for _, h := range headers {
			if v := r.Header.Get(h); v != "" {
				return v, nil
			}
		}

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return "", err
		}
		return ip, nil
	}
}

func (m *Middleware) WithRequestSizeLimiter() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, m.sizeLimit)
			h.ServeHTTP(w, r)
		})
	}
}
