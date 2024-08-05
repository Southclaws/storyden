package httpserver

import (
	"context"
	"net"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/internal/config"
)

func NewServer(lc fx.Lifecycle, l *zap.Logger, cfg config.Config, router *http.ServeMux) *http.Server {
	server := &http.Server{
		Handler: router,
		Addr:    cfg.ListenAddr,
	}

	wctx, cancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			server.BaseContext = func(ln net.Listener) context.Context { return wctx }
			go func() {
				// The HTTP server is the root node of dependency tree, because
				// it depends on everything else being initialised first. So, if
				// the app reaches this point, its considered a successful boot!

				l.Info("storyden http server starting",
					zap.String("address", cfg.ListenAddr),
					zap.String("cookie_domain", cfg.CookieDomain),
					zap.String("frontend_address", cfg.PublicWebAddress),
					zap.Bool("frontend_address", cfg.Production),
					zap.String("log_level", cfg.LogLevel.String()),
				)

				if err := server.ListenAndServe(); err != nil {
					l.Fatal("http server stopped unexpectedly", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(_ context.Context) error {
			cancel()
			return nil
		},
	})

	return server
}
