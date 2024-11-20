package rate

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/swirl"
)

type Limiter interface {
	Increment(ctx context.Context, key string, cost int) (*swirl.Status, bool, error)
	Check(ctx context.Context, key string, cost int) error
}

type swirlLimiter struct {
	rl *swirl.Limiter
}

func wrap(rl *swirl.Limiter) Limiter {
	return &swirlLimiter{
		rl: rl,
	}
}

func (l *swirlLimiter) Increment(ctx context.Context, key string, incr int) (*swirl.Status, bool, error) {
	status, allowed, err := l.rl.Increment(ctx, key, incr)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	if !allowed {
		return nil, false, fault.Wrap(status, fctx.With(ctx), ftag.With(ftag.PermissionDenied))
	}

	return status, true, nil
}

func (l *swirlLimiter) Check(ctx context.Context, key string, cost int) error {
	status, allowed, err := l.rl.Increment(ctx, key, cost)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if !allowed {
		return fault.Wrap(status, fctx.With(ctx), ftag.With(ftag.PermissionDenied))
	}

	return nil
}
