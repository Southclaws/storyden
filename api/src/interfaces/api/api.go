package api

import (
	"context"
	"net"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/api/src/config"
	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/interfaces/api/auth"
	"github.com/Southclaws/storyden/api/src/interfaces/api/swagger"
	"github.com/Southclaws/storyden/api/src/services/authentication"
)

func Build() fx.Option {
	return fx.Options(
		swagger.Build(),
		auth.Build(),
		// users.Build(),
		// categories.Build(),
		// metrics.Build(),
		// posts.Build(),
		// reacts.Build(),
		// subscriptions.Build(),
		// tags.Build(),
		// test.Build(),
		// threads.Build(),

		// Starts the HTTP server in a goroutine and fatals if it errors.
		fx.Invoke(func(l *zap.Logger, server *http.Server) {
			l.Debug("http server starting")
			go func() {
				if err := server.ListenAndServe(); err != nil {
					l.Fatal("HTTP server failed", zap.Error(err))
				}
			}()
		}),

		fx.Provide(func(as *authentication.CookieAuth, l *zap.Logger, cfg config.Config) chi.Router {
			router := chi.NewRouter()

			origins := []string{
				"http://localhost:3000", // Local development, `npm run dev`
				cfg.PublicWebAddress,    // Live public website
			}

			l.Debug("preparing router", zap.Strings("origins", origins))

			router.Use(
				web.WithLogger,
				cors.Handler(cors.Options{
					AllowedOrigins:   origins,
					AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
					AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Content-Length", "X-CSRF-Token"},
					ExposedHeaders:   []string{"Link", "Content-Length", "X-Ratelimit-Limit", "X-Ratelimit-Reset"},
					AllowCredentials: true,
					MaxAge:           300,
				}),
				as.WithAuthentication,
			)

			router.Get("/version", func(w http.ResponseWriter, r *http.Request) {
				web.Write(w, map[string]string{"version": config.Version}) //nolint:errcheck
			})

			router.HandleFunc(
				"/{rest:[a-zA-Z0-9=\\-\\/]+}",
				func(w http.ResponseWriter, r *http.Request) {
					if _, err := w.Write([]byte("no module found for that route")); err != nil {
						zap.L().Warn("failed to write error", zap.Error(err))
					}
				})

			return router
		}),

		fx.Provide(func(lc fx.Lifecycle, cfg config.Config, l *zap.Logger, router chi.Router) *http.Server {
			server := &http.Server{
				Handler: router,
				Addr:    cfg.ListenAddr,
			}

			lc.Append(fx.Hook{
				// Inject the global context into each request handler for
				// graceful shutdowns.
				// Note: The server isn't started here, instead, it's started
				// via the Invoke call above.
				OnStart: func(ctx context.Context) error {
					server.BaseContext = func(net.Listener) context.Context { return ctx }
					return nil
				},
				// Graceful shutdowns using the signal context.
				OnStop: func(ctx context.Context) error {
					return server.Shutdown(ctx)
				},
			})

			return server
		}),
	)
}
