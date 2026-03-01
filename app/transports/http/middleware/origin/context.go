package origin

import (
	"context"
)

type contextKey struct{}

var originContextKey = contextKey{}

func setOriginContext(ctx context.Context, origin string) context.Context {
	return context.WithValue(ctx, originContextKey, origin)
}

func GetOrigin(ctx context.Context) string {
	if origin, ok := ctx.Value(originContextKey).(string); ok {
		return origin
	}

	return ""
}
