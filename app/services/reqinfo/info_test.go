package reqinfo

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseConditionalRequestTimeAcceptsAllHTTPDateFormats(t *testing.T) {
	t.Parallel()

	want := time.Date(1994, time.November, 6, 8, 49, 37, 0, time.UTC)

	// the three date formats an http recipient must accept per rfc 7231
	cases := map[string]string{
		"imf-fixdate": "Sun, 06 Nov 1994 08:49:37 GMT",
		"rfc850":      "Sunday, 06-Nov-94 08:49:37 GMT",
		"asctime":     "Sun Nov  6 08:49:37 1994",
	}

	for name, raw := range cases {
		got, err := parseConditionalRequestTime(raw)
		assert.NoError(t, err, name)
		assert.True(t, got.Equal(want), "%s: got %v", name, got)
	}
}

func TestGettersPanicWithoutRequestInfo(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	assert.Panics(t, func() { _ = GetOperationID(ctx) })
	assert.Panics(t, func() { _ = GetDeviceName(ctx) })
	assert.Panics(t, func() { _ = GetCacheQuery(ctx) })
	assert.Panics(t, func() { _ = GetClientAddress(ctx) })
	assert.Panics(t, func() { _ = GetSSRClientAddress(ctx) })
}

func TestWithRequestInfoProvidesAllGetters(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("If-None-Match", `"t-2026-01-02T03:04:05Z"`)
	req.Header.Set("If-Modified-Since", time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC).Format(time.RFC1123))

	ctx := WithRequestInfo(context.Background(), req, "ThreadGet", "203.0.113.9", "198.51.100.22")

	assert.Equal(t, "ThreadGet", GetOperationID(ctx))
	assert.Equal(t, "203.0.113.9", GetClientAddress(ctx))
	assert.Equal(t, "198.51.100.22", GetSSRClientAddress(ctx))
	assert.NotZero(t, GetCacheQuery(ctx))
	assert.NotEmpty(t, GetDeviceName(ctx))
}
