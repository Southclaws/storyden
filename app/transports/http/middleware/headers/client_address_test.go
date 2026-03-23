package headers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/reqinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestMiddleware(cfg clientIPConfiguration, trustedSSRSourceCIDRs ...string) *Middleware {
	mw := &Middleware{
		trustedSSRSourceRanges: parseTrustedProxyCIDRs(trustedSSRSourceCIDRs),
	}
	mw.clientIPConfig.Store(cfg)
	return mw
}

func TestClientAddressRemoteAddrMode(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "203.0.113.5:1234"
	req.Header.Set("X-Forwarded-For", "198.51.100.1")

	key := newTestMiddleware(
		clientIPConfiguration{Mode: settings.ClientIPModeRemoteAddr},
	).clientAddress(req)

	assert.Equal(t, "203.0.113.5", key)
}

func TestClientAddressSingleHeaderMode(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "203.0.113.5:1234"
	req.Header.Set("CF-Connecting-IP", "198.51.100.1")

	key := newTestMiddleware(
		clientIPConfiguration{
			Mode:   settings.ClientIPModeSingleHeader,
			Header: "CF-Connecting-IP",
		},
	).clientAddress(req)

	assert.Equal(t, "198.51.100.1", key)
}

func TestClientAddressSingleHeaderModeTrimsConfiguredHeaderName(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "203.0.113.5:1234"
	req.Header.Set("X-Real-IP", "198.51.100.77")

	key := newTestMiddleware(
		clientIPConfiguration{
			Mode:   settings.ClientIPModeSingleHeader,
			Header: "  X-Real-IP  ",
		},
	).clientAddress(req)

	assert.Equal(t, "198.51.100.77", key)
}

func TestClientAddressTrustedXFFMode(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "172.16.0.4:1234"
	req.Header.Add("X-Forwarded-For", "198.51.100.9, 203.0.113.7")

	cfg := clientIPConfiguration{
		Mode:               settings.ClientIPModeXFFTrustedProxies,
		trustedProxyRanges: parseTrustedProxyCIDRs([]string{"172.16.0.0/12", "203.0.113.0/24"}),
	}

	key := newTestMiddleware(cfg).clientAddress(req)
	assert.Equal(t, "198.51.100.9", key)
}

func TestClientAddressTrustedXFFFallbackWhenRemoteUntrusted(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.0.2.20:1234"
	req.Header.Add("X-Forwarded-For", "198.51.100.9, 203.0.113.7")

	cfg := clientIPConfiguration{
		Mode:               settings.ClientIPModeXFFTrustedProxies,
		trustedProxyRanges: parseTrustedProxyCIDRs([]string{"172.16.0.0/12", "203.0.113.0/24"}),
	}

	key := newTestMiddleware(cfg).clientAddress(req)
	assert.Equal(t, "192.0.2.20", key)
}

func TestClientAddressTrustedXFFModeAllowsSSRLoopbackRemoteIPv6(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "[::1]:1234"
	req.Header.Set(ssrRequestHeader, "1")
	req.Header.Add("X-Forwarded-For", "198.51.100.9, 203.0.113.7")

	cfg := clientIPConfiguration{
		Mode:               settings.ClientIPModeXFFTrustedProxies,
		trustedProxyRanges: parseTrustedProxyCIDRs([]string{"203.0.113.0/24"}),
	}

	key := newTestMiddleware(cfg, "::1/128").clientAddress(req)
	assert.Equal(t, "198.51.100.9", key)
}

func TestClientAddressTrustedXFFModeAllowsSSRLoopbackRemoteIPv4(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	req.Header.Set(ssrRequestHeader, "1")
	req.Header.Add("X-Forwarded-For", "198.51.100.9, 203.0.113.7")

	cfg := clientIPConfiguration{
		Mode:               settings.ClientIPModeXFFTrustedProxies,
		trustedProxyRanges: parseTrustedProxyCIDRs([]string{"203.0.113.0/24"}),
	}

	key := newTestMiddleware(cfg, "127.0.0.1/32").clientAddress(req)
	assert.Equal(t, "198.51.100.9", key)
}

func TestClientAddressTrustedXFFModeSSRHeaderIgnoredWhenSSRSourceUntrusted(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "[::1]:1234"
	req.Header.Set(ssrRequestHeader, "1")
	req.Header.Add("X-Forwarded-For", "198.51.100.9, 203.0.113.7")

	cfg := clientIPConfiguration{
		Mode:               settings.ClientIPModeXFFTrustedProxies,
		trustedProxyRanges: parseTrustedProxyCIDRs([]string{"203.0.113.0/24"}),
	}

	key := newTestMiddleware(cfg).clientAddress(req)
	assert.Equal(t, "::1", key)
}

func TestClientAddressTrustedXFFModeTrustedSSRSourceNotAppliedWithoutSSRHeader(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	req.Header.Add("X-Forwarded-For", "198.51.100.9, 203.0.113.7")

	cfg := clientIPConfiguration{
		Mode:               settings.ClientIPModeXFFTrustedProxies,
		trustedProxyRanges: parseTrustedProxyCIDRs([]string{"203.0.113.0/24"}),
	}

	key := newTestMiddleware(cfg, "127.0.0.1/32").clientAddress(req)
	assert.Equal(t, "127.0.0.1", key)
}

func TestClientAddressTrustedXFFModeUsesConfiguredExternalSSRSource(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.12.0.9:1234"
	req.Header.Set(ssrRequestHeader, "1")
	req.Header.Add("X-Forwarded-For", "198.51.100.9, 203.0.113.7")

	cfg := clientIPConfiguration{
		Mode:               settings.ClientIPModeXFFTrustedProxies,
		trustedProxyRanges: parseTrustedProxyCIDRs([]string{"203.0.113.0/24"}),
	}

	key := newTestMiddleware(cfg, "10.12.0.0/16").clientAddress(req)
	assert.Equal(t, "198.51.100.9", key)
}

func TestClientAddressSSRHeaderDoesNotChangeSingleHeaderModeBehaviour(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.0.2.2:1234"
	req.Header.Set(ssrRequestHeader, "1")
	req.Header.Set("CF-Connecting-IP", "198.51.100.55")

	key := newTestMiddleware(clientIPConfiguration{
		Mode:   settings.ClientIPModeSingleHeader,
		Header: "CF-Connecting-IP",
	}).clientAddress(req)
	assert.Equal(t, "198.51.100.55", key)
}

func TestClientAddressSSRHeaderDoesNotChangeRemoteAddrModeBehaviour(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.0.2.2:1234"
	req.Header.Set(ssrRequestHeader, "1")
	req.Header.Set("X-Forwarded-For", "203.0.113.99, 198.51.100.22")
	req.Header.Set("X-Real-IP", "198.51.100.22")

	key := newTestMiddleware(
		clientIPConfiguration{Mode: settings.ClientIPModeRemoteAddr},
	).clientAddress(req)
	assert.Equal(t, "192.0.2.2", key)
}

func TestClientAddressSSRHeaderDoesNotChangeXFFTrustedProxyModeBehaviour(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "172.16.0.4:1234"
	req.Header.Set(ssrRequestHeader, "1")
	req.Header.Add("X-Forwarded-For", "198.51.100.9, 203.0.113.7")

	cfg := clientIPConfiguration{
		Mode:               settings.ClientIPModeXFFTrustedProxies,
		trustedProxyRanges: parseTrustedProxyCIDRs([]string{"172.16.0.0/12", "203.0.113.0/24"}),
	}

	key := newTestMiddleware(cfg).clientAddress(req)
	assert.Equal(t, "198.51.100.9", key)
}

func TestWithHeaderContextStoresClientAddress(t *testing.T) {
	t.Parallel()

	mw := newTestMiddleware(clientIPConfiguration{
		Mode:   settings.ClientIPModeSingleHeader,
		Header: "CF-Connecting-IP",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
	req.RemoteAddr = "172.16.0.1:9999"
	req.Header.Set("CF-Connecting-IP", "198.51.100.44")

	rr := httptest.NewRecorder()
	called := false

	handler := mw.WithHeaderContext()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		assert.Equal(t, "198.51.100.44", reqinfo.GetClientAddress(r.Context()))
		assert.Equal(t, "", reqinfo.GetSSRClientAddress(r.Context()))
		w.WriteHeader(http.StatusNoContent)
	}))

	handler.ServeHTTP(rr, req)

	require.True(t, called)
	require.Equal(t, http.StatusNoContent, rr.Code)
}

func TestWithHeaderContextStoresSSRClientAddress(t *testing.T) {
	t.Parallel()

	mw := newTestMiddleware(clientIPConfiguration{
		Mode:   settings.ClientIPModeSingleHeader,
		Header: "CF-Connecting-IP",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
	req.RemoteAddr = "192.0.2.2:9999"
	req.Header.Set(ssrRequestHeader, "1")
	req.Header.Set("CF-Connecting-IP", "198.51.100.77")

	rr := httptest.NewRecorder()
	called := false

	handler := mw.WithHeaderContext()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		assert.Equal(t, "198.51.100.77", reqinfo.GetSSRClientAddress(r.Context()))
		w.WriteHeader(http.StatusNoContent)
	}))

	handler.ServeHTTP(rr, req)

	require.True(t, called)
	require.Equal(t, http.StatusNoContent, rr.Code)
}

func TestWithHeaderContextStoresSSRResolvedClientAddressWithoutSSRIPHeader(t *testing.T) {
	t.Parallel()

	mw := newTestMiddleware(clientIPConfiguration{
		Mode:   settings.ClientIPModeSingleHeader,
		Header: "CF-Connecting-IP",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
	req.RemoteAddr = "192.0.2.2:9999"
	req.Header.Set(ssrRequestHeader, "1")
	req.Header.Set("CF-Connecting-IP", "198.51.100.91")

	rr := httptest.NewRecorder()
	called := false

	handler := mw.WithHeaderContext()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		assert.Equal(t, "198.51.100.91", reqinfo.GetClientAddress(r.Context()))
		assert.Equal(t, "198.51.100.91", reqinfo.GetSSRClientAddress(r.Context()))
		w.WriteHeader(http.StatusNoContent)
	}))

	handler.ServeHTTP(rr, req)

	require.True(t, called)
	require.Equal(t, http.StatusNoContent, rr.Code)
}

func TestParseTrustedProxyCIDRsAcceptsHostPrefixAndPlainIP(t *testing.T) {
	t.Parallel()

	ranges := parseTrustedProxyCIDRs([]string{
		"172.16.38.226/24",
		"203.0.113.7",
		"2001:db8::7/64",
	})

	require.Len(t, ranges, 3)
	assert.Equal(t, "172.16.38.0/24", ranges[0].String())
	assert.Equal(t, "203.0.113.7/32", ranges[1].String())
	assert.Equal(t, "2001:db8::/64", ranges[2].String())
}
