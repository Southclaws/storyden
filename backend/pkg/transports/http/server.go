package http

import (
	"context"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/backend/internal/config"
)

func newServer(lc fx.Lifecycle, l *zap.Logger, cfg config.Config, router chi.Router) *http.Server {
	server := &http.Server{
		Handler: router,
		Addr:    cfg.ListenAddr,
	}

	wctx, cancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			l.Info("http server starting", zap.String("address", cfg.ListenAddr))
			server.BaseContext = func(ln net.Listener) context.Context { return wctx }
			go func() {
				if err := server.ListenAndServe(); err != nil {
					l.Fatal("http server stopped unexpectedly", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			cancel()
			return nil
		},
	})

	l.Info("created http server")

	return server
}
