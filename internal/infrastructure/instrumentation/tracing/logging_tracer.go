package tracing

import (
	"context"
	"log/slog"
	"sync"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type slogExporter struct {
	logger *slog.Logger

	stoppedMu sync.RWMutex
	stopped   bool
}

func newLoggingTracer(logger *slog.Logger) trace.SpanExporter {
	return &slogExporter{
		logger: logger,
	}
}

func (e *slogExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	if err := ctx.Err(); err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("span context error"))
	}
	e.stoppedMu.RLock()
	stopped := e.stopped
	e.stoppedMu.RUnlock()
	if stopped {
		return nil
	}

	if len(spans) == 0 {
		return nil
	}

	stubs := tracetest.SpanStubsFromReadOnlySpans(spans)

	for i := range stubs {
		stub := &stubs[i]

		duration := stub.EndTime.Sub(stub.StartTime)

		args := []any{
			slog.String("trace_id", stub.SpanContext.TraceID().String()),
			slog.String("span_id", stub.SpanContext.SpanID().String()),
			slog.Duration("duration", duration),
		}

		for _, attr := range stub.Attributes {
			args = append(args, toSlog(attr))
		}

		for _, ev := range stub.Events {
			for _, attr := range ev.Attributes {
				args = append(args, toSlog(attr))
			}
		}

		fn := e.stubLevel(stub)
		fn(stub.Name, args...)
	}
	return nil
}

func (e *slogExporter) stubLevel(stub *tracetest.SpanStub) func(string, ...any) {
	switch stub.Status.Code {
	case codes.Error:
		return e.logger.Error
	default:
		return e.logger.Info
	}
}

func (e *slogExporter) Shutdown(ctx context.Context) error {
	e.stoppedMu.Lock()
	e.stopped = true
	e.stoppedMu.Unlock()

	return nil
}

func (e *slogExporter) MarshalLog() interface{} {
	return struct {
		Type           string
		WithTimestamps bool
	}{
		Type:           "log",
		WithTimestamps: true,
	}
}

func toSlog(attr attribute.KeyValue) any {
	switch attr.Value.Type() {
	case attribute.BOOL:
		return slog.Bool(string(attr.Key), attr.Value.AsBool())

	case attribute.INT64:
		return slog.Int64(string(attr.Key), attr.Value.AsInt64())

	case attribute.FLOAT64:
		return slog.Float64(string(attr.Key), attr.Value.AsFloat64())

	case attribute.STRING:
		return slog.String(string(attr.Key), attr.Value.AsString())

	default:
		return slog.Any(string(attr.Key), attr.Value.AsInterface())
	}
}
