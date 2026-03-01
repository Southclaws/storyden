package noop

import (
	"context"
	"log/slog"

	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/kv"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/spanner"
)

func NewBuilder() spanner.Builder {
	return builder{}
}

type builder struct{}

func (builder) Build() spanner.Instrumentation {
	return instrumentation{}
}

type instrumentation struct{}

func (instrumentation) Instrument(ctx context.Context, _ ...kv.Attr) (context.Context, spanner.Span) {
	return ctx, span{ctx: ctx}
}

func (instrumentation) InstrumentNamed(ctx context.Context, _ string, _ ...kv.Attr) (context.Context, spanner.Span) {
	return ctx, span{ctx: ctx}
}

type span struct {
	ctx context.Context
}

func (span) End() {}

func (s span) Annotate(_ ...kv.Attr) context.Context {
	return s.ctx
}

func (span) Event(_ string, _ ...kv.Attr) {}

func (span) Logger() *slog.Logger {
	return slog.Default()
}

func (span) Wrap(err error, _ string, _ ...kv.Attr) error {
	return err
}
