package tracing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace/noop"

	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/servertiming"
)

func TestTracerRecordsServerTimingSpans(t *testing.T) {
	t.Parallel()

	collector := servertiming.NewCollector()
	ctx := servertiming.WithCollector(t.Context(), collector)

	tr := tracer{inner: noop.NewTracerProvider().Tracer("test")}

	_, span := tr.Start(ctx, "svc/work")
	time.Sleep(6 * time.Millisecond)
	span.End()

	header := collector.HeaderValue()
	assert.Contains(t, header, "svc;dur=")
	assert.Contains(t, header, "desc=\"work\"")
}
