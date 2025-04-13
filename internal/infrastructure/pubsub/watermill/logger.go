package watermill

import (
	"log/slog"

	"github.com/Southclaws/dt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/samber/lo"
)

type logAdapter struct {
	z *slog.Logger
}

func (l *logAdapter) Error(msg string, err error, fields watermill.LogFields) {
	l.z.Error(msg, l.slogToWatermill(fields)...)
}

func (l *logAdapter) Info(msg string, fields watermill.LogFields) {
	l.z.Info(msg, l.slogToWatermill(fields)...)
}

func (l *logAdapter) Debug(msg string, fields watermill.LogFields) {
	l.z.Debug(msg, l.slogToWatermill(fields)...)
}

func (l *logAdapter) Trace(msg string, fields watermill.LogFields) {
	l.z.Debug(msg, l.slogToWatermill(fields)...)
}

func (l *logAdapter) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return &logAdapter{z: l.z.With(l.slogToWatermill(fields)...)}
}

func (l *logAdapter) slogToWatermill(fs watermill.LogFields) []any {
	entries := lo.Entries(fs)

	return dt.Map(entries, func(e lo.Entry[string, any]) any {
		return slog.Any(e.Key, e.Value)
	})
}
