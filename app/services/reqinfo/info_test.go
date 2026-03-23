package reqinfo

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
