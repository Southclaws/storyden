package spanner

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/kv"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/tracing"
)

// Build is to be used in services during initialisation.
func (i *service) Build() Instrumentation {
	pc, _, _, _ := runtime.Caller(1)
	path := runtime.FuncForPC(pc).Name()
	caller := filepath.Base(path)
	pkg := strings.Split(caller, ".")[0]
	tr := i.tf.Build(i.lc, pkg)

	return &impl{
		logger: i.lg,
		tracer: tr,
	}
}

type impl struct {
	logger *slog.Logger
	tracer tracing.Tracer
}

func (i *impl) Instrument(ctx context.Context, a ...kv.Attr) (context.Context, Span) {
	pkg, fn, loc := getLoc(runtime.Caller(1))

	a = append(a,
		kv.String("package", pkg),
		kv.String("location", loc),
	)

	return i.InstrumentNamed(ctx, fn, a...)
}

func (i *impl) InstrumentNamed(ctx context.Context, name string, a ...kv.Attr) (context.Context, Span) {
	logger := i.logger.With(kv.Attrs(a).ToSlog()...)

	// Create a child span with the KV data as attributes.
	ctx, span := i.tracer.Start(ctx, name, trace.WithAttributes(kv.Attrs(a).ToAttributes()...))

	// Create a new context with the KV data as fctx metadata.
	ctx = fctx.WithMeta(ctx, kv.Attrs(a).ToFault()...)

	// NOTE: We store the context into the tracking span so it can be mutated
	// and new child contexts can be returned with the new attributes.
	return ctx, &trackingSpan{span, logger, ctx, name}
}

type trackingSpan struct {
	span   trace.Span
	logger *slog.Logger
	// NOTE: We store ctx because we need to mutate it in Annotate.
	//nolint:containedctx
	ctx    context.Context
	caller string
}

func (t *trackingSpan) End() {
	t.span.End()
}

func (t *trackingSpan) Annotate(a ...kv.Attr) context.Context {
	//
	// NON-OBVIOUS MUTATIONS AHEAD
	//
	// We need to mutate the context and the logger to add the attributes.
	// The reason for this is the ctx and logger are passed in to the Instrument
	// function (which is always called once at the start of a procedure) and we
	// don't want to force users to explicitly pass in the ctx and logger for
	// every single annotation throughout the procedure. This means that when
	// a span is annotated, we can still return a new context and child logger.
	//
	t.ctx = fctx.WithMeta(t.ctx, kv.Attrs(a).ToFault()...)

	// Add the attributes to the span. This is a slightly more obvious mutation.
	t.span.SetAttributes(kv.Attrs(a).ToAttributes()...)

	// Mutate the stored logger with the same attributes.
	t.logger = t.logger.With(kv.Attrs(a).ToSlog()...)

	return t.ctx
}

func (t *trackingSpan) Event(name string, a ...kv.Attr) {
	t.span.AddEvent(name, trace.WithStackTrace(true), trace.WithAttributes(kv.Attrs(a).ToAttributes()...))
}

func (t *trackingSpan) Logger() *slog.Logger {
	return t.logger
}

func (t *trackingSpan) Wrap(err error, msg string, a ...kv.Attr) error {
	t.span.SetStatus(codes.Error, msg)
	t.span.RecordError(err, trace.WithAttributes(kv.Attrs(a).ToAttributes()...))

	ctx := fctx.WithMeta(t.ctx, kv.Attrs(a).ToFault()...)

	return fault.Wrap(err, fctx.With(ctx), fmsg.With(msg))
}

func getLoc(pc uintptr, file string, line int, ok bool) (string, string, string) {
	if !ok {
		return "", "", ""
	}

	file = file[strings.Index(file, "captain/app")+len("captain/"):]
	loc := fmt.Sprintf("%s:%d", file, line)

	path := runtime.FuncForPC(pc).Name() // get the fully qualified object-path to the caller function
	caller := filepath.Base(path)        // get the caller in "package.Function" format
	sep := strings.Index(caller, ".")    // split `package.Function``
	if sep == -1 {
		sep = len(caller)
	}
	pkg := caller[:sep] // get package name

	return pkg, caller, loc
}
