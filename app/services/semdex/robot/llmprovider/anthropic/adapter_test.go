package anthropic

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	anthropicapi "github.com/anthropics/anthropic-sdk-go"
	anthropicoption "github.com/anthropics/anthropic-sdk-go/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/robot/llm_provider"
	"github.com/Southclaws/storyden/app/resources/robot/model_ref"
)

func TestModelCapabilitiesProbesAnthropicImageSupport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		supported := r.URL.Path == "/v1/models/claude-with-images"
		w.Header().Set("Content-Type", "application/json")
		_, err := fmt.Fprintf(w, `{
			"id":"model",
			"type":"model",
			"display_name":"Model",
			"created_at":"2025-01-01T00:00:00Z",
			"max_input_tokens":1000,
			"max_tokens":1000,
			"capabilities":{"image_input":{"supported":%t}}
		}`, supported)
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	client := anthropicapi.NewClient(
		anthropicoption.WithAPIKey("test-key"),
		anthropicoption.WithBaseURL(server.URL+"/"),
	)
	provider := &Anthropic{client: &client}

	supported, err := provider.ModelCapabilities(context.Background(), model_ref.ModelRef{Provider: Provider, Model: model_ref.NewModel("claude-with-images")})
	require.NoError(t, err)
	assert.Equal(t, llm_provider.CapabilitySupportSupported, supported.ImageInput)

	unsupported, err := provider.ModelCapabilities(context.Background(), model_ref.ModelRef{Provider: Provider, Model: model_ref.NewModel("claude-text-only")})
	require.NoError(t, err)
	assert.Equal(t, llm_provider.CapabilitySupportUnsupported, unsupported.ImageInput)
}
