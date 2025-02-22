package tracing

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/getsentry/sentry-go"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/internal/config"
)

type factory struct {
	provider string
	opts     []trace.TracerProviderOption
}

func Build() fx.Option {
	return fx.Provide(newExporter, newTracerFactory)
}

func newTracerFactory(ctx context.Context,
	cfg config.Config,
	logger *zap.Logger,
	opts []trace.TracerProviderOption,
) (Factory, error) {
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		logger.Error("otel error", zap.String("error", err.Error()))
	}))

	return factory{
		provider: cfg.OTELProvider,
		opts:     opts,
	}, nil
}

func newExporter(ctx context.Context,
	cfg config.Config,
	logger *zap.Logger,
) ([]trace.TracerProviderOption, error) {
	switch cfg.OTELProvider {
	case "sentry":

		if cfg.SentryDSN == "" {
			if cfg.OTELEndpoint.String() != "" {
				return nil, fault.New("OTEL_EXPORTER_OTLP_ENDPOINT is set but sentry DSN is required instead when using the sentry provider")
			}
			return nil, fault.New("sentry DSN is required when using the sentry provider")
		}

		err := sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.SentryDSN,
			EnableTracing:    true,
			TracesSampleRate: 1.0,
		})
		if err != nil {
			return nil, err
		}

		spanProc := sentryotel.NewSentrySpanProcessor()

		// for some reason, sentry is a "span processor" not a "span exporter".
		return []trace.TracerProviderOption{
			trace.WithSpanProcessor(spanProc),
		}, nil

	case "otlp":
		endpoint := cfg.OTELEndpoint.String()
		if endpoint == "" {
			return nil, fault.New("OTEL_EXPORTER_OTLP_ENDPOINT is required when using the otlp provider")
		}

		opts := []otlptracehttp.Option{
			otlptracehttp.WithEndpointURL(endpoint),
		}

		if cfg.OTELEndpoint.Scheme != "https" {
			opts = append(opts, otlptracehttp.WithInsecure())
		}

		otlp, err := otlptracehttp.New(ctx, opts...)
		if err != nil {
			return nil, fault.Wrap(err, fmsg.With("failed to create OTLP exporter"))
		}

		return []trace.TracerProviderOption{
			trace.WithBatcher(otlp),
		}, nil

	case "logger":
		return []trace.TracerProviderOption{
			trace.WithSyncer(newLoggingTracer(logger)),
		}, nil

	default:
		return []trace.TracerProviderOption{}, nil
	}
}

// Build constructs a new tracer for use within a system component.
func (f factory) Build(lc fx.Lifecycle, serviceName string) Tracer {
	opts := append(f.opts, trace.WithResource(resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
	)))

	tp := trace.NewTracerProvider(opts...)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if err := tp.Shutdown(ctx); err != nil {
				return fault.Wrap(err)
			}

			return nil
		},
	})

	tracer := tp.Tracer("storyden")

	return tracer
}
