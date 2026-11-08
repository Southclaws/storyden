package reqinfo

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	ctx := WithRequestInfo(context.Background(), req, "ThreadGet", "203.0.113.9", "198.51.100.22")

	assert.Equal(t, "ThreadGet", GetOperationID(ctx))
	assert.Equal(t, "203.0.113.9", GetClientAddress(ctx))
	assert.Equal(t, "198.51.100.22", GetSSRClientAddress(ctx))
	assert.True(t, GetCacheQuery(ctx).MatchesETag(`"t-2026-01-02T03:04:05Z"`))
	assert.NotEmpty(t, GetDeviceName(ctx))
}

func TestDeviceNamesFromRequestHeaders(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		header string
		want   string
	}{
		{
			name:   "desktop Edge",
			header: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36 Edg/130.0.0.0",
			want:   "Edge (Windows)",
		},
		{
			name:   "iOS Safari",
			header: "Mozilla/5.0 (iPhone; CPU iPhone OS 18_0 like Mac OS X) AppleWebKit/605.1.15 Version/18.0 Mobile/15E148 Safari/604.1",
			want:   "Safari (iOS)",
		},
		{
			name:   "Android Chrome",
			header: "Mozilla/5.0 (Linux; Android 14; K) AppleWebKit/537.36 Chrome/130.0.0.0 Mobile Safari/537.36",
			want:   "Chrome (Android)",
		},
		{name: "missing header", want: "Unknown (Unknown)"},
		{name: "unknown browser", header: "CustomClient/1 (Linux)", want: "Unknown (Linux)"},
		{name: "HTML injection", header: "<script>alert(1)</script>", want: "Unknown (Unknown)"},
		{name: "line injection", header: "Chrome/130\r\nSet-Cookie: stolen=1", want: "Unknown (Unknown)"},
		{name: "format injection", header: "Chrome/130\u202e", want: "Unknown (Unknown)"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("User-Agent", tc.header)
			ctx := WithRequestInfo(context.Background(), req, "AccountAuthProviderAdd", "203.0.113.9", "")

			require.Equal(t, tc.want, GetDeviceName(ctx))
			require.Equal(t, tc.header, req.Header.Get("User-Agent"))
		})
	}
}

func TestRequestDeviceNamesAreIndependent(t *testing.T) {
	t.Parallel()

	first := httptest.NewRequest("GET", "/test", nil)
	first.Header.Set("User-Agent", "Firefox/130 (Linux)")
	firstContext := WithRequestInfo(context.Background(), first, "First", "", "")

	second := httptest.NewRequest("GET", "/test", nil)
	second.Header.Set("User-Agent", "Chrome/130 (Windows NT 10.0)")
	secondContext := WithRequestInfo(firstContext, second, "Second", "", "")

	first.Header.Set("User-Agent", "Changed/1")

	require.Equal(t, "Firefox (Linux)", GetDeviceName(firstContext))
	require.Equal(t, "Chrome (Windows)", GetDeviceName(secondContext))
}
