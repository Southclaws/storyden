package spanner

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/kv"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/tracing"
)

type Builder interface {
	Build() Instrumentation
}

type Instrumentation interface {
	Instrument(ctx context.Context, a ...kv.Attr) (context.Context, Span)
	InstrumentNamed(ctx context.Context, name string, a ...kv.Attr) (context.Context, Span)
}

type Span interface {
	End()
	Annotate(a ...kv.Attr) context.Context
	Event(name string, a ...kv.Attr)
	Wrap(err error, msg string, a ...kv.Attr) error
}

func TraceID(ctx context.Context) string {
	return trace.SpanFromContext(ctx).SpanContext().TraceID().String()
}

func SpanID(ctx context.Context) string {
	return trace.SpanFromContext(ctx).SpanContext().SpanID().String()
}

func New(lc fx.Lifecycle, logger, tf tracing.Factory) Builder {
	return &service{
		lc: lc,
		tf: tf,
	}
}

type service struct {
	lc fx.Lifecycle
	tf tracing.Factory
}
