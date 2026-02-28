package headers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Southclaws/storyden/app/services/reqinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		remoteAddr string
		headers    map[string][]string
		want       string
	}{
		{
			name:       "forwarded preferred and rightmost",
			remoteAddr: "172.16.0.1:1234",
			headers: map[string][]string{
				"Forwarded":       {`for=198.51.100.9;proto=https, for=198.51.100.10`},
				"X-Forwarded-For": {"203.0.113.7"},
			},
			want: "198.51.100.10",
		},
		{
			name:       "forwarded quoted ipv6 with port",
			remoteAddr: "172.16.0.1:1234",
			headers: map[string][]string{
				"Forwarded": {`for="[2001:db8:cafe::17]:4711";proto=https`},
			},
			want: "2001:db8:cafe::17",
		},
		{
			name:       "xff rightmost entry",
			remoteAddr: "172.16.0.1:1234",
			headers: map[string][]string{
				"X-Forwarded-For": {"203.0.113.4, 203.0.113.5"},
			},
			want: "203.0.113.5",
		},
		{
			name:       "xff strips port",
			remoteAddr: "172.16.0.1:1234",
			headers: map[string][]string{
				"X-Forwarded-For": {"203.0.113.4:5432"},
			},
			want: "203.0.113.4",
		},
		{
			name:       "multiple xff header values uses last header",
			remoteAddr: "172.16.0.1:1234",
			headers: map[string][]string{
				"X-Forwarded-For": {"203.0.113.4, 203.0.113.5", "198.51.100.77"},
			},
			want: "198.51.100.77",
		},
		{
			name:       "unknown skipped",
			remoteAddr: "172.16.0.1:1234",
			headers: map[string][]string{
				"Forwarded":       {`for=unknown`},
				"X-Forwarded-For": {"unknown, 203.0.113.99"},
			},
			want: "203.0.113.99",
		},
		{
			name:       "remote addr fallback strips port",
			remoteAddr: "172.16.0.1:1234",
			want:       "172.16.0.1",
		},
		{
			name:       "remote addr raw fallback",
			remoteAddr: "127.0.0.1",
			want:       "127.0.0.1",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest(http.MethodGet, "/test", nil)
			assert.NoError(t, err)
			req.RemoteAddr = tt.remoteAddr

			for k, values := range tt.headers {
				for _, v := range values {
					req.Header.Add(k, v)
				}
			}

			assert.Equal(t, tt.want, clientAddress(req))
		})
	}
}

func TestWithHeaderContextStoresClientAddress(t *testing.T) {
	t.Parallel()

	mw := New()

	req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
	req.RemoteAddr = "172.16.0.1:9999"
	req.Header.Add("X-Forwarded-For", "203.0.113.10, 198.51.100.1")

	rr := httptest.NewRecorder()
	called := false

	handler := mw.WithHeaderContext()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		assert.Equal(t, "198.51.100.1", reqinfo.GetClientAddress(r.Context()))
		w.WriteHeader(http.StatusNoContent)
	}))

	handler.ServeHTTP(rr, req)

	require.True(t, called)
	require.Equal(t, http.StatusNoContent, rr.Code)
}
