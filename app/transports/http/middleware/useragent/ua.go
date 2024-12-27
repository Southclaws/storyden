package useragent

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mileusna/useragent"
)

type Middleware struct{}

func New() *Middleware {
	return &Middleware{}
}

type uaKey struct{}

// WithUserAgentContext stores in the request context the user agent info.
func (m *Middleware) WithUserAgentContext() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			ua := useragent.Parse(r.Header.Get("User-Agent"))

			newctx := context.WithValue(ctx, uaKey{}, ua)

			r = r.WithContext(newctx)

			next.ServeHTTP(w, r)
		})
	}
}

func GetDeviceName(ctx context.Context) string {
	v := ctx.Value(uaKey{})
	ua, ok := v.(useragent.UserAgent)
	if !ok {
		return "Unknown"
	}

	return fmt.Sprintf("%s (%s)", ua.Name, ua.OS)
}
