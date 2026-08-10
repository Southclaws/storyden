package tracing

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/servertiming"
)

type Tracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}

type Factory interface {
	Build(lc fx.Lifecycle) Tracer
}

type tracer struct {
	inner trace.Tracer
}

func (t tracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	ctx, span := t.inner.Start(ctx, spanName, opts...)

	collector := servertiming.FromContext(ctx)
	if collector == nil {
		return ctx, span
	}

	return ctx, &timedSpan{
		Span:      span,
		name:      spanName,
		collector: collector,
		start:     time.Now(),
	}
}

type timedSpan struct {
	trace.Span
	name      string
	collector *servertiming.Collector
	start     time.Time
}

func (s *timedSpan) End(options ...trace.SpanEndOption) {
	s.collector.Observe(s.name, time.Since(s.start))
	s.Span.End(options...)
}
