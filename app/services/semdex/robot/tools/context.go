package tools

import (
	"context"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
)

type RunContext struct {
	RobotID   opt.Optional[xid.ID]
	RobotRef  string
	AccountID string
	SessionID string
}

// ToolAudience identifies the kind of caller receiving a tool's result.
type ToolAudience string

const (
	ToolAudienceUnknown ToolAudience = ""
	ToolAudienceRobot   ToolAudience = "robot"
	ToolAudienceMCP     ToolAudience = "mcp"
)

type toolAudienceContextKey struct{}

// ContextWithToolAudience records the declared tool result audience on ctx.
func ContextWithToolAudience(ctx context.Context, audience ToolAudience) context.Context {
	return context.WithValue(ctx, toolAudienceContextKey{}, audience)
}

// ToolAudienceFromContext returns ToolAudienceUnknown when no caller declared an audience.
func ToolAudienceFromContext(ctx context.Context) ToolAudience {
	v, _ := ctx.Value(toolAudienceContextKey{}).(ToolAudience)
	return v
}

type runContextKey struct{}

func ContextWithRunContext(ctx context.Context, run RunContext) context.Context {
	return context.WithValue(ctx, runContextKey{}, run)
}

func RunContextFromContext(ctx context.Context) RunContext {
	v, _ := ctx.Value(runContextKey{}).(RunContext)
	return v
}
