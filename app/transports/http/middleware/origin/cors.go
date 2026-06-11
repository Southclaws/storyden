package origin

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/rs/cors"

	"github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
	"github.com/Southclaws/storyden/internal/config"
)

type Middleware struct {
	cfg config.Config
}

func New(cfg config.Config) *Middleware {
	return &Middleware{cfg: cfg}
}

func (m *Middleware) WithCORS() func(next http.Handler) http.Handler {
	allowedMethods := []string{
		"GET",
		"POST",
		"PUT",
		"PATCH",
		"DELETE",
		"OPTIONS",
	}

	allowedHeaders := []string{
		"Accept",
		"Authorization",
		"Content-Type",
		"Content-Length",
		"X-CSRF-Token",
		"X-Correlation-ID",
		"X-Forwarded-Host",
	}

	exposedHeaders := []string{
		"Link",
		"Content-Type",
		"Content-Length",
		"X-Ratelimit-Limit",
		"X-Ratelimit-Remaining",
		"X-Ratelimit-Reset",
	}

	allowOriginFunc := m.makeAllowOriginFunc()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			corsConfig := cors.New(cors.Options{
				AllowOriginFunc:  allowOriginFunc,
				AllowedMethods:   allowedMethods,
				AllowedHeaders:   allowedHeaders,
				ExposedHeaders:   exposedHeaders,
				AllowCredentials: true,
				MaxAge:           300,
			})

			ctx := setOriginContext(r.Context(), origin)

			corsConfig.Handler(next).ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (m *Middleware) makeAllowOriginFunc() func(origin string) bool {
	webAddr := m.cfg.PublicWebAddress
	apiAddr := m.cfg.PublicAPIAddress

	return func(origin string) bool {
		if origin == "" {
			return false
		}

		originURL, err := url.Parse(origin)
		if err != nil {
			return false
		}

		originHost := originURL.Hostname()
		if originHost == "" {
			return false
		}

		if isOriginAllowed(originHost, webAddr, apiAddr) {
			return true
		}

		return false
	}
}

func isOriginAllowed(originHost string, webAddr, apiAddr url.URL) bool {
	webHost := webAddr.Hostname()
	apiHost := apiAddr.Hostname()

	if originHost == webHost || originHost == apiHost {
		return true
	}

	if webHost == apiHost {
		return isSubdomainOfRoot(originHost, webHost)
	}

	webDomain, err := session_cookie.DomainFromString(webHost)
	if err != nil {
		return false
	}

	apiDomain, err := session_cookie.DomainFromString(apiHost)
	if err != nil {
		return false
	}

	webRoot := webDomain.GetETLDp1()
	apiRoot := apiDomain.GetETLDp1()

	if !webRoot.IsEqual(apiRoot) {
		return false
	}

	originDomain, err := session_cookie.DomainFromString(originHost)
	if err != nil {
		return false
	}

	return originDomain.IsSubdomainOf(webRoot) || originDomain.IsEqual(webRoot)
}

func isSubdomainOfRoot(host, rootHost string) bool {
	if host == rootHost {
		return true
	}

	if !strings.HasSuffix(host, "."+rootHost) {
		return false
	}

	prefix := strings.TrimSuffix(host, "."+rootHost)
	return !strings.Contains(prefix, ":")
}
