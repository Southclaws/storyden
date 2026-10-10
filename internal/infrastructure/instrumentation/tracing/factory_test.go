package tracing

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/trace"
	collectortrace "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/protobuf/proto"

	"github.com/Southclaws/storyden/internal/config"
)

func TestSentryExporterSendsSpans(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	type request struct {
		path string
		auth string
		body []byte
		err  error
	}
	requests := make(chan request, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		requests <- request{r.URL.Path, r.Header.Get("X-Sentry-Auth"), body, err}
		w.Header().Set("Content-Type", "application/x-protobuf")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	previousClient := sentry.CurrentHub().Client()
	defer sentry.CurrentHub().BindClient(previousClient)

	opts, err := newExporter(ctx, config.Config{
		OTELProvider: "sentry",
		SentryDSN:    strings.Replace(server.URL, "://", "://test-public-key@", 1) + "/123",
	}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	require.NoError(t, err)

	provider := trace.NewTracerProvider(opts...)
	t.Cleanup(func() { require.NoError(t, provider.Shutdown(context.Background())) })
	_, span := provider.Tracer("storyden-test").Start(ctx, "dependency-upgrade-span")
	span.End()
	require.NoError(t, provider.ForceFlush(ctx))

	select {
	case received := <-requests:
		require.NoError(t, received.err)
		require.Equal(t, "/api/123/integration/otlp/v1/traces/", received.path)
		require.Contains(t, received.auth, "sentry_key=test-public-key")
		var payload collectortrace.ExportTraceServiceRequest
		require.NoError(t, proto.Unmarshal(received.body, &payload))
		require.Len(t, payload.ResourceSpans, 1)
		require.Len(t, payload.ResourceSpans[0].ScopeSpans, 1)
		spans := payload.ResourceSpans[0].ScopeSpans[0].Spans
		require.Len(t, spans, 1)
		require.Equal(t, "dependency-upgrade-span", spans[0].Name)
	case <-ctx.Done():
		t.Fatal("Sentry exporter did not send the completed span")
	}
}
