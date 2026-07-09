package headers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithAssetSecurityHeaders(t *testing.T) {
	t.Parallel()

	handler := (&Middleware{}).WithAssetSecurityHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("neutralises active content on asset responses", func(t *testing.T) {
		t.Parallel()

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/assets/evil-svg", nil))

		assert.Equal(t, "nosniff", rec.Header().Get("X-Content-Type-Options"))
		assert.Equal(t, "default-src 'none'; sandbox", rec.Header().Get("Content-Security-Policy"))
	})

	t.Run("leaves non-asset responses untouched", func(t *testing.T) {
		t.Parallel()

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/threads", nil))

		assert.Empty(t, rec.Header().Get("X-Content-Type-Options"))
		assert.Empty(t, rec.Header().Get("Content-Security-Policy"))
	})
}
