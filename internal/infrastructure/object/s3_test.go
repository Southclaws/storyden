package object

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestS3StorerExists(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name    string
		status  int
		code    string
		exists  bool
		wantErr bool
	}{
		{name: "existing", status: http.StatusOK, exists: true},
		{name: "missing_object", status: http.StatusNotFound},
		{name: "missing_bucket", status: http.StatusNotFound, code: "NoSuchBucket", wantErr: true},
		{name: "forbidden", status: http.StatusForbidden, wantErr: true},
		{name: "server_error", status: http.StatusInternalServerError, wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodHead, r.Method)
				assert.Equal(t, "/storyden/assets/icon-512x512.png", r.URL.Path)

				w.Header().Set("Content-Length", "1")
				w.Header().Set("Last-Modified", "Mon, 02 Jan 2006 15:04:05 GMT")
				w.Header().Set("x-minio-error-code", tc.code)
				w.WriteHeader(tc.status)
			}))
			defer server.Close()

			client, err := minio.New(strings.TrimPrefix(server.URL, "http://"), &minio.Options{
				Region:     "us-east-1",
				MaxRetries: 1,
			})
			require.NoError(t, err)

			store := &s3Storer{bucket: "storyden", minioClient: client}
			exists, err := store.Exists(t.Context(), "assets/icon-512x512.png")
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.exists, exists)
		})
	}
}
