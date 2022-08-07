package http

import (
	"context"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/backend/internal/config"
)

func newServer(lc fx.Lifecycle, l *zap.Logger, cfg config.Config, router *echo.Echo) *http.Server {
	server := &http.Server{
		Handler: router,
		Addr:    cfg.ListenAddr,
	}

	wctx, cancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			l.Info("http server starting", zap.String("address", cfg.ListenAddr))
			server.BaseContext = func(ln net.Listener) context.Context { return wctx }
			go func() {
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

	l.Info("created http server")

	return server
}
