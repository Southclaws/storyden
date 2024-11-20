package rate

import (
	"time"

	"github.com/Southclaws/swirl"

	"github.com/Southclaws/storyden/internal/infrastructure/cache"
)

type LimiterFactory struct {
	store swirl.Store
}

func NewFactory(
	store cache.Store,
) *LimiterFactory {
	return &LimiterFactory{
		store: store,
	}
}

func (f *LimiterFactory) NewLimiter(
	limit int,
	period time.Duration,
	expiry time.Duration,
) Limiter {
	return wrap(swirl.New(f.store, limit, period, expiry))
}
