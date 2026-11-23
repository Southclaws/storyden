package tools

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/lib/mcp"
)

func Authorise(ctx context.Context, def *mcp.ToolDefinition) error {
	if acc, ok := session.GetOptAccount(ctx).Get(); ok {
		if err := acc.RejectSuspended(); err != nil {
			return err
		}
	}

	perm, ok := def.RequiredPermission.Get()
	if !ok {
		return nil
	}

	if !session.HasContext(ctx) {
		return fault.New("tool requires an authenticated session", ftag.With(ftag.Unauthenticated))
	}

	return session.Authorise(ctx, nil, perm)
}
