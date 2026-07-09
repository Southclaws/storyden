package headers

import (
	"net/http"
	"strings"
)

const assetPathPrefix = "/api/assets/"

// user-uploaded files are served inline from the same origin, so neutralise active content to prevent stored xss
const assetContentSecurityPolicy = "default-src 'none'; sandbox"

func (m *Middleware) WithAssetSecurityHeaders() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, assetPathPrefix) {
				h := w.Header()
				h.Set("X-Content-Type-Options", "nosniff")
				h.Set("Content-Security-Policy", assetContentSecurityPolicy)
			}

			next.ServeHTTP(w, r)
		})
	}
}
