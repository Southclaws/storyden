package pinecone

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	pineconesdk "github.com/pinecone-io/go-pinecone/v7/pinecone"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetOrCreateIndex(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name         string
		describeCode int
		createCode   int
		wantErr      bool
	}{
		{"existing", http.StatusOK, 0, false},
		{"missing", http.StatusNotFound, http.StatusCreated, false},
		{"forbidden", http.StatusForbidden, 0, true},
		{"creation rejected", http.StatusNotFound, http.StatusBadRequest, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			requests := make(chan string, 4)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests <- r.Method + " " + r.URL.Path
				assert.Equal(t, "test-key", r.Header.Get("Api-Key"))
				w.Header().Set("Content-Type", "application/json")
				status := tc.describeCode
				if r.Method == http.MethodPost {
					status = tc.createCode
					var body map[string]any
					if assert.NoError(t, json.NewDecoder(r.Body).Decode(&body)) {
						assert.Equal(t, "storyden-test", body["name"])
						assert.Equal(t, map[string]any{"fields": map[string]any{
							"_values": map[string]any{"type": "dense_vector", "dimension": float64(1536), "metric": "cosine"},
						}}, body["schema"])
						assert.Equal(t, map[string]any{"deployment_type": "managed", "cloud": "aws", "region": "us-east-1"}, body["deployment"])
					}
				}
				if status == 0 {
					status = http.StatusInternalServerError
				}
				w.WriteHeader(status)
				if status >= 400 {
					_, _ = io.WriteString(w, `{"message":"request rejected"}`)
					return
				}
				_, _ = io.WriteString(w, `{"name":"storyden-test","host":"127.0.0.1:1","status":{"ready":true,"state":"Ready"},"deployment":{"deployment_type":"managed","cloud":"aws","region":"us-east-1"},"schema":{"fields":{"_values":{"type":"dense_vector","dimension":1536,"metric":"cosine"}}}}`)
			}))
			defer server.Close()

			sdk, err := pineconesdk.NewClient(pineconesdk.NewClientParams{ApiKey: "test-key", Host: server.URL, RestClient: server.Client()})
			require.NoError(t, err)
			client := &Client{Client: sdk, size: 1536, cloud: "aws", region: "us-east-1"}
			index, err := client.GetOrCreateIndex(t.Context(), "storyden-test")
			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, index)
			} else {
				require.NoError(t, err)
				defer index.Close()
				require.Equal(t, "storyden", index.Namespace())
			}
			close(requests)
			var got []string
			for request := range requests {
				got = append(got, request)
			}
			want := []string{"GET /indexes/storyden-test"}
			if tc.describeCode == http.StatusNotFound {
				want = append(want, "POST /indexes")
			}
			require.Equal(t, want, got)
		})
	}
}
