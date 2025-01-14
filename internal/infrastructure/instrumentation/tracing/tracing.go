package tracing

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

type Tracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}

type Factory interface {
	Build(lc fx.Lifecycle, serviceName string) Tracer
}
