package logger

import (
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Southclaws/storyden/internal/config"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(newLogger),
		fx.Invoke(replaceGlobals),
	)
}

func newLogger(cfg config.Config) (*zap.Logger, error) {
	var zapconfig zap.Config
	if cfg.Production {
		zapconfig = zap.NewProductionConfig()
		zapconfig.InitialFields = map[string]interface{}{"v": config.Version}
	} else {
		zapconfig = zap.NewDevelopmentConfig()
	}

	zapconfig.DisableStacktrace = true

	zapconfig.Level.SetLevel(cfg.LogLevel)
	zapconfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := zapconfig.Build()
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to build zap config"))
	}

	return logger, nil
}

func replaceGlobals(c config.Config, l *zap.Logger) {
	// Use our logger for globals too, even though it's passed to
	// dependents most of the time using DI, the global logger is used
	// in a couple of places during startup/shutdown.
	zap.ReplaceGlobals(l)

	if !c.Production {
		l.Debug("logger configured in development mode")
	}
}
