package headers

import (
	"net/http"

	"github.com/Southclaws/storyden/app/services/reqinfo"
)

type Middleware struct{}

func New() *Middleware {
	return &Middleware{}
}

// WithHeaderContext stores in the request context header info.
func (m *Middleware) WithHeaderContext() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			newctx := reqinfo.WithRequestInfo(ctx, r)

			r = r.WithContext(newctx)

			next.ServeHTTP(w, r)
		})
	}
}
