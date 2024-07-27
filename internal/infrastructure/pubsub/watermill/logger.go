package watermill

import (
	"github.com/Southclaws/dt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type logAdapter struct {
	z *zap.Logger
}

func (l *logAdapter) Error(msg string, err error, fields watermill.LogFields) {
	l.z.Error(msg, l.zapToWatermill(fields)...)
}

func (l *logAdapter) Info(msg string, fields watermill.LogFields) {
	l.z.Info(msg, l.zapToWatermill(fields)...)
}

func (l *logAdapter) Debug(msg string, fields watermill.LogFields) {
	l.z.Debug(msg, l.zapToWatermill(fields)...)
}

func (l *logAdapter) Trace(msg string, fields watermill.LogFields) {
	l.z.Debug(msg, l.zapToWatermill(fields)...)
}

func (l *logAdapter) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return &logAdapter{z: l.z.With(l.zapToWatermill(fields)...)}
}

func (l *logAdapter) zapToWatermill(fs watermill.LogFields) []zap.Field {
	entries := lo.Entries(fs)

	return dt.Map(entries, func(e lo.Entry[string, any]) zap.Field {
		return zap.Any(e.Key, e.Value)
	})
}
