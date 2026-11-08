package anthropic

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Southclaws/fault/fmsg"
	anthropicapi "github.com/anthropics/anthropic-sdk-go/v2"
	anthropicoption "github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/robot/model_ref"
)

func TestModelCapabilitiesMapsAnthropicV2Errors(t *testing.T) {
	cases := []struct {
		name   string
		status int
		body   string
		issue  string
	}{
		{
			name:   "structured API error",
			status: http.StatusBadRequest,
			body:   `{"type":"error","error":{"type":"invalid_request_error","message":"bad message"}}`,
			issue:  "bad message",
		},
		{
			name:   "unstructured gateway error",
			status: http.StatusBadGateway,
			body:   "Bad Gateway",
			issue:  "Bad Gateway (status 502)",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				_, err := io.WriteString(w, tc.body)
				assert.NoError(t, err)
			}))
			t.Cleanup(server.Close)

			client := anthropicapi.NewClient(
				anthropicoption.WithAPIKey("test-key"),
				anthropicoption.WithBaseURL(server.URL+"/"),
				anthropicoption.WithMaxRetries(0),
			)
			provider := &Anthropic{client: &client}

			_, err := provider.ModelCapabilities(context.Background(), model_ref.ModelRef{
				Provider: Provider,
				Model:    model_ref.NewModel("claude-test"),
			})
			require.ErrorIs(t, err, errAnthropicAPI)
			assert.Equal(t, tc.issue, string(fmsg.GetIssue(err)))
		})
	}
}
