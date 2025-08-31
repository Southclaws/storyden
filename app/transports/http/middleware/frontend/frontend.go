package frontend

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
	"github.com/Southclaws/storyden/internal/config"
	frontendService "github.com/Southclaws/storyden/internal/infrastructure/frontend"
)

type Provider struct {
	handler  func(http.ResponseWriter, *http.Request)
	frontend frontendService.Frontend
	logger   *slog.Logger
}

func New(
	cfg config.Config,
	logger *slog.Logger,
	mux *http.ServeMux,
	cj *session_cookie.Jar,
	fe frontendService.Frontend,
) *Provider {
	if cfg.FrontendProxy.String() == "" {
		return &Provider{}
	}

	handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err, ok := recover().(error); ok && err != nil {
					if errors.Is(err, http.ErrAbortHandler) {
						return
					}

					logger.Error("frontend proxy panic",
						slog.String("url", r.URL.String()),
						slog.String("method", r.Method),
						slog.String("remote_addr", r.RemoteAddr),
						slog.Any("error", err),
					)
					return
				}
			}()

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
		w.WriteHeader(http.StatusBadGateway)
	}

	return &Provider{
		handler:  handler(proxy),
		frontend: fe,
		logger:   logger,
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
				if p.frontend == nil {
					p.handler(w, r)
					return
				}

				// Wait for frontend to be ready before proxying
				select {
				case <-p.frontend.Ready():
					p.handler(w, r)

				case <-time.After(30 * time.Second):
					p.logger.Error("timeout waiting for frontend to be ready",
						slog.String("url", r.URL.String()),
						slog.String("method", r.Method),
						slog.String("remote_addr", r.RemoteAddr),
					)
					w.WriteHeader(http.StatusServiceUnavailable)
					w.Write([]byte("Frontend service is unavailable"))
				}

			}
		})
	}
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(New),
	)
}
