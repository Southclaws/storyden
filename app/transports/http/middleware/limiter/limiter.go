package limiter

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/reqinfo"
	"github.com/Southclaws/storyden/app/transports/http/openapi/operation"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/infrastructure/rate"
)

const (
	RateLimitLimit      = "X-RateLimit-Limit"
	RateLimitRemaining  = "X-RateLimit-Remaining"
	RateLimitReset      = "X-RateLimit-Reset"
	RetryAfter          = "Retry-After"
	MaxRequestSizeBytes = 10 * 1024 * 1024
)

type Middleware struct {
	rl               atomic.Pointer[rate.Limiter]
	factory          *rate.LimiterFactory
	currentLimit     atomic.Int32
	currentPeriod    atomic.Int64
	currentBucket    atomic.Int64
	currentGuestCost atomic.Int32
	configLimit      int
	configPeriod     time.Duration
	configBucket     time.Duration
	configGuestCost  int
	kf               KeyFunc
	sizeLimit        int64
	settingsRepo     *settings.SettingsRepository
	logger           *slog.Logger
}

func New(
	ctx context.Context,
	lc fx.Lifecycle,
	cfg config.Config,
	f *rate.LimiterFactory,
	settingsRepo *settings.SettingsRepository,
	bus *pubsub.Bus,
	logger *slog.Logger,
) *Middleware {
	m := &Middleware{
		factory:         f,
		configLimit:     cfg.RateLimit,
		configPeriod:    cfg.RateLimitPeriod,
		configBucket:    cfg.RateLimitBucket,
		configGuestCost: cfg.RateLimitGuestCost,
		kf:              fromIP("CF-Connecting-IP", "X-Real-IP", "True-Client-IP"),
		sizeLimit:       MaxRequestSizeBytes,
		settingsRepo:    settingsRepo,
		logger:          logger,
	}

	lc.Append(fx.StartHook(func(hctx context.Context) error {
		m.reconfigureLimiter(hctx)

		_, err := pubsub.Subscribe(ctx, bus, "limiter.settings_updated", func(ctx context.Context, evt *message.EventSettingsUpdated) error {
			m.reconfigureLimiter(ctx)
			return nil
		})
		return err
	}))

	return m
}

func (m *Middleware) getConfiguration(ctx context.Context) (limit int, period time.Duration, bucket time.Duration, guestCost int) {
	limit = m.configLimit
	period = m.configPeriod
	bucket = m.configBucket
	guestCost = m.configGuestCost

	appSettings, err := m.settingsRepo.Get(ctx)
	if err != nil {
		m.logger.Warn("failed to fetch settings for rate limiter reconfiguration",
			slog.String("error", err.Error()),
		)
		return
	}

	svc, ok := appSettings.Services.Get()
	if !ok {
		return
	}

	rl, ok := svc.RateLimit.Get()
	if !ok {
		return
	}

	if v, ok := rl.RateLimit.Get(); ok && v > 0 {
		limit = v
	}
	if v, ok := rl.RateLimitPeriod.Get(); ok && v > 0 {
		period = v
	}
	if v, ok := rl.RateLimitBucket.Get(); ok && v > 0 {
		bucket = v
	}
	if v, ok := rl.RateLimitGuestCost.Get(); ok && v > 0 {
		guestCost = v
	}

	return limit, period, bucket, guestCost
}

func (m *Middleware) reconfigureLimiter(ctx context.Context) {
	currentLimit := int(m.currentLimit.Load())
	currentPeriod := time.Duration(m.currentPeriod.Load())
	currentBucket := time.Duration(m.currentBucket.Load())
	currentGuestCost := int(m.currentGuestCost.Load())

	limit, period, bucket, guestCost := m.getConfiguration(ctx)

	if m.rl.Load() != nil && limit == currentLimit && period == currentPeriod && bucket == currentBucket && guestCost == currentGuestCost {
		return
	}

	m.logger.Debug("reconfiguring rate limiter",
		slog.Int("limit", limit),
		slog.Duration("period", period),
		slog.Duration("bucket", bucket),
		slog.Int("guest_cost", guestCost),
	)

	m.currentLimit.Store(int32(limit))
	m.currentPeriod.Store(int64(period))
	m.currentBucket.Store(int64(bucket))
	m.currentGuestCost.Store(int32(guestCost))

	newLimiter := m.factory.NewLimiter(limit, period, bucket)
	m.rl.Store(&newLimiter)
}

func (m *Middleware) WithRateLimit() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			key, err := m.kf(r)
			if err != nil {
				m.logger.Error("failed to extract rate limit key from request",
					slog.String("error", err.Error()),
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
				)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			isAuthenticated := session.GetOptAccountID(ctx).Ok()

			cost := m.getCost(ctx)

			if !isAuthenticated {
				guestCost := int(m.currentGuestCost.Load())
				cost = cost * guestCost
			}

			limiter := m.rl.Load()
			if limiter == nil {
				m.logger.Error("rate limiter not initialized",
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
				)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				panic("limiter: rl pointer is null")
			}

			status, allowed, err := (*limiter).Increment(ctx, key, cost)
			if err != nil {
				m.logger.Error("failed to increment rate limit counter",
					slog.String("error", err.Error()),
					slog.String("key", key),
					slog.Int("cost", cost),
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
				)
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

func (m *Middleware) getCost(ctx context.Context) int {
	operationID, err := operation.NewOperationID(reqinfo.GetOperationID(ctx))
	if err != nil {
		return 1
	}

	overrides := m.getLookup(ctx)

	if cost, ok := overrides[operationID.String()]; ok {
		return cost
	}

	if cost, ok := operation.RateLimitCostOverrides[operationID]; ok {
		return cost
	}

	return 1
}

func (m *Middleware) getLookup(ctx context.Context) map[string]int {
	appSettings, err := m.settingsRepo.Get(ctx)
	if err != nil {
		return nil
	}
	svc, ok := appSettings.Services.Get()
	if !ok {
		return nil
	}

	rl, ok := svc.RateLimit.Get()
	if !ok {
		return nil
	}

	overrides, ok := rl.CostOverrides.Get()
	if !ok {
		return nil
	}

	return overrides
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
