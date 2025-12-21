package limiter

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/rate"
	"github.com/Southclaws/storyden/app/resources/settings"
)

const (
	RateLimitLimit      = "X-RateLimit-Limit"
	RateLimitRemaining  = "X-RateLimit-Remaining"
	RateLimitReset      = "X-RateLimit-Reset"
	RetryAfter          = "Retry-After"
	MaxRequestSizeBytes = 10 * 1024 * 1024
)

type Middleware struct {
	rl          rate.Limiter
	factory     *rate.LimiterFactory
	settings    *settings.SettingsRepository
	cfg         config.Config
	kf          KeyFunc
	sizeLimit   int64
}

func New(
	cfg config.Config,
	f *rate.LimiterFactory,
	settingsRepo *settings.SettingsRepository,
) *Middleware {
	rl := f.NewLimiter(cfg.RateLimit, cfg.RateLimitPeriod, cfg.RateLimitExpire)

	return &Middleware{
		rl:        rl,
		factory:   f,
		settings:  settingsRepo,
		cfg:       cfg,
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

			// Try to get operation ID from Echo context
			var operationID string
			if echoCtx, ok := r.Context().Value("echo").(echo.Context); ok {
				// Get the route path which we'll use as a lookup key
				routePath := echoCtx.Path()
				operationID = m.getOperationIDFromPath(routePath, r.Method)
			}

			// Get the appropriate limiter and cost for this operation
			limiter, cost := m.getLimiterAndCost(ctx, operationID, key)

			status, allowed, err := limiter.Increment(ctx, key, cost)
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

// getLimiterAndCost returns the appropriate limiter and cost for an operation
func (m *Middleware) getLimiterAndCost(ctx context.Context, operationID string, key string) (rate.Limiter, int) {
	// Default cost
	cost := 1
	
	// If no operation ID, use global limiter
	if operationID == "" {
		return m.rl, cost
	}

	// Check for settings override first
	settingsData, err := m.settings.Get(ctx)
	if err == nil && settingsData.RateLimitOverrides.Ok() {
		overrides, _ := settingsData.RateLimitOverrides.Get()
		if override, ok := overrides[operationID]; ok {
			return m.createLimiterFromOverride(operationID, override), overrideCostOrDefault(override.Cost, cost)
		}
	}

	// Check for OpenAPI spec config
	if opConfig := GetOperationConfig(operationID); opConfig != nil {
		cost = opConfig.Cost
		
		// If limit is specified in spec, create a specific limiter
		if opConfig.Limit > 0 {
			period := opConfig.Period
			if period == 0 {
				period = m.cfg.RateLimitPeriod
			}
			
			limiterKey := operationID
			return m.factory.GetOrCreateLimiter(limiterKey, opConfig.Limit, period, m.cfg.RateLimitExpire), cost
		}
	}

	// Use global limiter with potentially custom cost
	return m.rl, cost
}

// createLimiterFromOverride creates a limiter from settings override
func (m *Middleware) createLimiterFromOverride(operationID string, override settings.RateLimitOverride) rate.Limiter {
	limit := override.Limit
	if limit == 0 {
		limit = m.cfg.RateLimit
	}

	period := m.cfg.RateLimitPeriod
	if override.Period != "" {
		if d, err := time.ParseDuration(override.Period); err == nil {
			period = d
		}
	}

	limiterKey := "override:" + operationID
	return m.factory.GetOrCreateLimiter(limiterKey, limit, period, m.cfg.RateLimitExpire)
}

func overrideCostOrDefault(overrideCost, defaultCost int) int {
	if overrideCost > 0 {
		return overrideCost
	}
	return defaultCost
}

// getOperationIDFromPath maps a route path to an operation ID using generated mapping
func (m *Middleware) getOperationIDFromPath(path, method string) string {
	return GetOperationIDFromRoute(method, path)
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
