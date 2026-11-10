package fetcher

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/enetx/surf"
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

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "4096")
		w.WriteHeader(http.StatusOK)
		w.Write(make([]byte, 4096))
	}))
	defer srv.Close()

	_, err := fetchBounded(context.Background(), srv.Client(), srv.URL, 16)

	require.Error(t, err)
	assert.True(t, errors.Is(err, errAssetTooLarge))
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

func TestFetchBoundedRejectsDecompressedOversizeBody(t *testing.T) {
	t.Parallel()

	var compressed bytes.Buffer
	writer, err := gzip.NewWriterLevel(&compressed, gzip.BestCompression)
	require.NoError(t, err)
	_, err = writer.Write(bytes.Repeat([]byte("x"), 8192))
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	require.Less(t, compressed.Len(), 1024)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.Write(compressed.Bytes())
	}))
	defer server.Close()

	result := surf.NewClient().Builder().Impersonate().MacOS().Chrome().SecureTLS().Build()
	require.False(t, result.IsErr())
	client := result.Ok().Std()
	defer client.CloseIdleConnections()

	_, err = fetchBounded(context.Background(), client, server.URL, 1024)
	require.ErrorIs(t, err, errAssetTooLarge)
}
