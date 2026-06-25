package tools

import (
	"context"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
)

type RunContext struct {
	RobotID   opt.Optional[xid.ID]
	AccountID string
	SessionID string
}

type runContextKey struct{}

func ContextWithRunContext(ctx context.Context, run RunContext) context.Context {
	return context.WithValue(ctx, runContextKey{}, run)
}

func RunContextFromContext(ctx context.Context) RunContext {
	v, _ := ctx.Value(runContextKey{}).(RunContext)
	return v
}
