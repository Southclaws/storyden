package web

import (
	"net/http"
)

// RouteUse is router.User but for individual routes.
func RouteUse(mws ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(f http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var handler = f
			for _, mw := range mws {
				handler = mw(handler)
			}
			handler.ServeHTTP(w, r)
		})
	}
}
