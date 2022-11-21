package http

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
