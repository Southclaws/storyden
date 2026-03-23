package bindings

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildNetworkHeadersSampleWhitelistedHeaders(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("GET", "/api/admin", nil)
	r.RemoteAddr = "192.0.2.9:1234"
	r.Header.Set("X-Forwarded-For", "203.0.113.9")
	r.Header.Set("User-Agent", "storyden-test")

	sample := buildNetworkHeadersSample(r)
	require.NotNil(t, sample)
	require.NotNil(t, sample.Headers)
	require.NotNil(t, sample.RawClientAddress)

	headers := *sample.Headers
	assert.Equal(t, "203.0.113.9", headers["x-forwarded-for"])
	assert.Equal(t, "192.0.2.9:1234", *sample.RawClientAddress)
	assert.NotContains(t, headers, "user-agent")
}

func TestBuildNetworkHeadersSampleSSRForwardedHeaders(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("GET", "/api/admin", nil)
	r.RemoteAddr = "192.0.2.10:2345"
	r.Header.Set("X-Storyden-SSR", "1")
	r.Header.Set("X-Forwarded-For", "198.51.100.44")
	r.Header.Set("X-Custom-Client-IP", "198.51.100.55") // should be ignored

	sample := buildNetworkHeadersSample(r)
	require.NotNil(t, sample)
	require.NotNil(t, sample.Headers)
	require.NotNil(t, sample.HeadersSsr)
	require.NotNil(t, sample.RawClientAddress)

	direct := *sample.Headers
	headers := *sample.HeadersSsr
	assert.Equal(t, "198.51.100.44", direct["x-forwarded-for"])
	assert.Equal(t, "198.51.100.44", headers["x-forwarded-for"])
	assert.NotContains(t, direct, "x-custom-client-ip")
	assert.NotContains(t, headers, "x-custom-client-ip")
	assert.Equal(t, "192.0.2.10:2345", *sample.RawClientAddress)
}

func TestBuildNetworkHeadersSampleSSRForwardedHeadersDenylist(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("GET", "/api/admin", nil)
	r.Header.Set("X-Storyden-SSR", "1")
	r.Header.Set("X-Forwarded-For", "198.51.100.44")
	r.Header.Set("Cookie", "session=bad")
	r.Header.Set("Authorization", "Bearer bad")
	r.Header.Set("X-Nonstandard-IP", "198.51.100.20")
	r.Header.Set("X-Storyden-SSR-Client-IP", "203.0.113.2")

	sample := buildNetworkHeadersSample(r)
	require.NotNil(t, sample)
	require.NotNil(t, sample.Headers)
	require.NotNil(t, sample.HeadersSsr)

	direct := *sample.Headers
	headers := *sample.HeadersSsr
	assert.Equal(t, "198.51.100.44", direct["x-forwarded-for"])
	assert.Equal(t, "198.51.100.44", headers["x-forwarded-for"])
	assert.NotContains(t, direct, "cookie")
	assert.NotContains(t, direct, "authorization")
	assert.NotContains(t, direct, "x-storyden-ssr")
	assert.NotContains(t, direct, "x-storyden-ssr-client-ip")
	assert.NotContains(t, direct, "x-nonstandard-ip")
	assert.NotContains(t, headers, "cookie")
	assert.NotContains(t, headers, "authorization")
	assert.NotContains(t, headers, "x-storyden-ssr")
	assert.NotContains(t, headers, "x-storyden-ssr-client-ip")
	assert.NotContains(t, headers, "x-nonstandard-ip")
}

func TestBuildNetworkHeadersSampleWithoutSSRMarkerHasNoSSRHeaders(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("GET", "/api/admin", nil)
	r.Header.Set("X-Forwarded-For", "198.51.100.44")

	sample := buildNetworkHeadersSample(r)
	require.NotNil(t, sample)
	require.NotNil(t, sample.Headers)
	headers := *sample.Headers
	assert.Equal(t, "198.51.100.44", headers["x-forwarded-for"])
	assert.Nil(t, sample.HeadersSsr)
	require.NotNil(t, sample.RawClientAddress)
}

func TestBuildNetworkHeadersSampleRawClientAddressOnly(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest("GET", "/api/admin", nil)
	r.RemoteAddr = "198.51.100.99:4040"

	sample := buildNetworkHeadersSample(r)
	require.NotNil(t, sample)
	assert.Nil(t, sample.Headers)
	assert.Nil(t, sample.HeadersSsr)
	require.NotNil(t, sample.RawClientAddress)
	assert.Equal(t, "198.51.100.99:4040", *sample.RawClientAddress)
}
