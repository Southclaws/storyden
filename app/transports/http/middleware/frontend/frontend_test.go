package frontend

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/config"
)

func TestWithFrontendProxyPreservesIncomingXFF(t *testing.T) {
	t.Parallel()

	upstreamXFF := make(chan string, 1)
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upstreamXFF <- r.Header.Get("X-Forwarded-For")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer upstream.Close()

	frontendProxyURL, err := url.Parse(upstream.URL)
	require.NoError(t, err)

	p := New(
		config.Config{FrontendProxy: *frontendProxyURL},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		http.NewServeMux(),
		nil,
		nil,
	)

	handler := p.WithFrontendProxy()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "172.16.38.226:50450"
	req.Header.Set("X-Forwarded-For", "2a00:23ee:1930:908e:b04a:176b:8bd:efd2, 2a09:8280:1::3:6ef5")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)
	assert.Equal(
		t,
		"2a00:23ee:1930:908e:b04a:176b:8bd:efd2, 2a09:8280:1::3:6ef5",
		<-upstreamXFF,
	)
}

func TestWithFrontendProxyBypassesAPIRequests(t *testing.T) {
	t.Parallel()

	upstreamCalled := false
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upstreamCalled = true
		w.WriteHeader(http.StatusNoContent)
	}))
	defer upstream.Close()

	frontendProxyURL, err := url.Parse(upstream.URL)
	require.NoError(t, err)

	p := New(
		config.Config{FrontendProxy: *frontendProxyURL},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		http.NewServeMux(),
		nil,
		nil,
	)

	handler := p.WithFrontendProxy()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusAccepted, rr.Code)
	assert.False(t, upstreamCalled)
}
