package origin

import (
	"context"
)

type contextKey struct{}

func setOriginContext(ctx context.Context, origin string) context.Context {
	return context.WithValue(ctx, contextKey{}, origin)
}

func GetOrigin(ctx context.Context) string {
	if origin, ok := ctx.Value(contextKey{}).(string); ok {
		return origin
	}

	return ""
}
