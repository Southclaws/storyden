package rate

import (
	"sync"
	"time"

	"github.com/Southclaws/swirl"

	"github.com/Southclaws/storyden/internal/infrastructure/cache"
)

type LimiterFactory struct {
	store swirl.Store
	mu    sync.RWMutex
	cache map[string]Limiter // cache of operation-specific limiters
}

func NewFactory(
	store cache.Store,
) *LimiterFactory {
	return &LimiterFactory{
		store: store,
		cache: make(map[string]Limiter),
	}
}

func (f *LimiterFactory) NewLimiter(
	limit int,
	period time.Duration,
	expiry time.Duration,
) Limiter {
	return wrap(swirl.New(f.store, limit, period, expiry))
}

// GetOrCreateLimiter returns a cached limiter for the given key, or creates a new one
func (f *LimiterFactory) GetOrCreateLimiter(
	key string,
	limit int,
	period time.Duration,
	expiry time.Duration,
) Limiter {
	f.mu.RLock()
	limiter, ok := f.cache[key]
	f.mu.RUnlock()
	
	if ok {
		return limiter
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	// Double-check in case another goroutine created it
	if limiter, ok := f.cache[key]; ok {
		return limiter
	}

	limiter = wrap(swirl.New(f.store, limit, period, expiry))
	f.cache[key] = limiter
	return limiter
}
