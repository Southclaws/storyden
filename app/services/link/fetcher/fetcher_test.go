package fetcher

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchBoundedAcceptsSmallBody(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("tiny"))
	}))
	defer srv.Close()

	body, err := fetchBounded(context.Background(), srv.Client(), srv.URL, 16)

	require.NoError(t, err)
	assert.Equal(t, "tiny", string(body))
}

func TestFetchBoundedRejectsDeclaredOversizeBody(t *testing.T) {
	t.Parallel()

	served := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "4096")
		w.WriteHeader(http.StatusOK)
		served, _ = w.Write(make([]byte, 4096))
	}))
	defer srv.Close()

	_, err := fetchBounded(context.Background(), srv.Client(), srv.URL, 16)

	require.Error(t, err)
	assert.True(t, errors.Is(err, errAssetTooLarge))
	assert.NotZero(t, served)
}

func TestFetchBoundedRejectsUndeclaredOversizeBody(t *testing.T) {
	t.Parallel()

	// no Content-Length, the response streams until the reader gives up
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Transfer-Encoding", "chunked")
		w.WriteHeader(http.StatusOK)
		for range 64 {
			if _, err := w.Write(make([]byte, 1024)); err != nil {
				return
			}
			w.(http.Flusher).Flush()
		}
	}))
	defer srv.Close()

	_, err := fetchBounded(context.Background(), srv.Client(), srv.URL, 16)

	require.Error(t, err)
	assert.True(t, errors.Is(err, errAssetTooLarge))
}

func TestFetchBoundedRejectsNonOK(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := fetchBounded(context.Background(), srv.Client(), srv.URL, 16)

	assert.Error(t, err)
}
