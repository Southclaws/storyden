package ratelimit

import (
	"time"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// TODO: Make this configurable.
const (
	UnauthenticatedRPS = 50
	AuthenticatedRPS   = 100
)

type Limiter struct {
	Guest  *stdlib.Middleware
	Member *stdlib.Middleware
}

func New() *Limiter {
	store := memory.NewStore()

	// Rate limit for unauthenticated requests.
	guestRate := limiter.Rate{
		Period: time.Second,
		Limit:  UnauthenticatedRPS,
	}
	guest := stdlib.NewMiddleware(limiter.New(store, guestRate))

	// Rate limit for authenticated requests.
	memberRate := limiter.Rate{
		Period: time.Second,
		Limit:  AuthenticatedRPS,
	}
	member := stdlib.NewMiddleware(limiter.New(store, memberRate))

	return &Limiter{
		Guest:  guest,
		Member: member,
	}
}
