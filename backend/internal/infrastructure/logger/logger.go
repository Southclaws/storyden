package logger

import (
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Southclaws/storyden/backend/internal/config"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(func(cfg config.Config) (*zap.Logger, error) {
			var zapconfig zap.Config
			if cfg.Production {
				zapconfig = zap.NewProductionConfig()
				zapconfig.InitialFields = map[string]interface{}{"v": config.Version}
			} else {
				zapconfig = zap.NewDevelopmentConfig()
			}

			zapconfig.Level.SetLevel(cfg.LogLevel)
			zapconfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

			logger, err := zapconfig.Build()
			if err != nil {
				return nil, err
			}

			return logger, nil
		}),
		fx.Invoke(func(c config.Config, l *zap.Logger) {
			// Use our logger for globals too, even though it's passed to
			// dependents most of the time using DI, the global logger is used
			// in a couple of places during startup/shutdown.
			zap.ReplaceGlobals(l)
			if !c.Production {
				l.Info("logger configured in development mode")
			}
		}),
	)
}

// WithLogger is simple Zap logger HTTP middleware
func WithLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zap.L().Info(
			"request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Any("query", r.URL.Query()),
			zap.Any("headers", r.Header),
			zap.Int64("body", r.ContentLength),
		)
		next.ServeHTTP(w, r)
	})
}
