package frontend

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"strings"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
	"github.com/Southclaws/storyden/internal/config"
)

type Provider struct {
	handler func(http.ResponseWriter, *http.Request)
}

func New(
	cfg config.Config,
	logger *slog.Logger,
	mux *http.ServeMux,
	cj *session_cookie.Jar,
) *Provider {
	if cfg.FrontendProxy.String() == "" {
		return &Provider{}
	}

	handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			p.ServeHTTP(w, r)
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(&cfg.FrontendProxy)

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Error("frontend proxy error",
			slog.String("url", r.URL.String()),
			slog.String("method", r.Method),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("error", err.Error()),
		)
	}

	return &Provider{
		handler: handler(proxy),
	}
}

func (p *Provider) WithFrontendProxy() func(next http.Handler) http.Handler {
	if p.handler == nil {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api") {
				next.ServeHTTP(w, r)
			} else {
				p.handler(w, r)
			}
		})
	}
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(New),
	)
}
