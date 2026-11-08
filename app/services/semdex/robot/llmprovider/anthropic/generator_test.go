package anthropic

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	anthropicapi "github.com/anthropics/anthropic-sdk-go/v2"
	anthropicoption "github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

func TestGenerateContentUsesAnthropicV2(t *testing.T) {
	cases := []struct {
		name        string
		stream      bool
		contentType string
		response    string
	}{
		{
			name:        "synchronous",
			contentType: "application/json",
			response: `{
				"id":"msg_test","type":"message","role":"assistant","model":"claude-test",
				"content":[{"type":"text","text":"Done"}],"stop_reason":"end_turn",
				"usage":{"input_tokens":5,"output_tokens":1}
			}`,
		},
		{
			name:        "streaming",
			stream:      true,
			contentType: "text/event-stream",
			response: `event: message_start
data: {"type":"message_start","message":{"id":"msg_test","type":"message","role":"assistant","model":"claude-test","content":[],"stop_reason":null,"usage":{"input_tokens":5,"output_tokens":0}}}

event: content_block_start
data: {"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Done"}}

event: content_block_stop
data: {"type":"content_block_stop","index":0}

event: message_delta
data: {"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":1}}

event: message_stop
data: {"type":"message_stop"}

`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/v1/messages", r.URL.Path)
				assert.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))
				assert.Empty(t, r.Header.Get("X-Api-Key"))
				assert.Equal(t, "2026-10-21", r.Header.Get("Anthropic-Version"))

				var body map[string]json.RawMessage
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Errorf("decode request: %v", err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				assert.JSONEq(t, `"claude-test"`, string(body["model"]))
				assert.JSONEq(t, `8096`, string(body["max_tokens"]))
				assert.JSONEq(t, `[{"type":"text","text":"Be helpful."}]`, string(body["system"]))
				assert.JSONEq(t, `[
					{"role":"user","content":[{"type":"text","text":"Find it."}]},
					{"role":"assistant","content":[{"type":"tool_use","id":"call_1","name":"lookup","input":{"query":"storyden"}}]},
					{"role":"user","content":[{"type":"tool_result","tool_use_id":"call_1","content":[{"type":"text","text":"{\"found\":true}"}],"is_error":false}]}
				]`, string(body["messages"]))
				assert.NotContains(t, body, "temperature")
				assert.NotContains(t, body, "top_p")
				assert.NotContains(t, body, "top_k")

				if tc.stream {
					assert.JSONEq(t, `true`, string(body["stream"]))
				}

				w.Header().Set("Content-Type", tc.contentType)
				_, err := io.WriteString(w, tc.response)
				assert.NoError(t, err)
			}))
			t.Cleanup(server.Close)

			client := anthropicapi.NewClient(
				anthropicoption.WithAPIKey("test-key"),
				anthropicoption.WithBaseURL(server.URL+"/"),
				anthropicoption.WithMaxRetries(0),
			)
			provider := &Anthropic{client: &client, modelName: "claude-test"}
			request := &model.LLMRequest{
				Config: &genai.GenerateContentConfig{
					SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: "Be helpful."}}},
				},
				Contents: []*genai.Content{
					{Role: genai.RoleUser, Parts: []*genai.Part{{Text: "Find it."}}},
					{Role: genai.RoleModel, Parts: []*genai.Part{{FunctionCall: &genai.FunctionCall{
						ID: "call_1", Name: "lookup", Args: map[string]any{"query": "storyden"},
					}}}},
					{Role: genai.RoleUser, Parts: []*genai.Part{{FunctionResponse: &genai.FunctionResponse{
						ID: "call_1", Name: "lookup", Response: map[string]any{"found": true},
					}}}},
				},
			}

			var responses []*model.LLMResponse
			for response, err := range provider.GenerateContent(context.Background(), request, tc.stream) {
				require.NoError(t, err)
				responses = append(responses, response)
			}

			expectedResponses := 1
			if tc.stream {
				expectedResponses = 2
			}

			require.Len(t, responses, expectedResponses)
			if tc.stream {
				assert.True(t, responses[0].Partial)
				assert.Equal(t, "Done", responses[0].Content.Parts[0].Text)
			}

			final := responses[len(responses)-1]
			assert.True(t, final.TurnComplete)
			assert.Equal(t, genai.FinishReasonStop, final.FinishReason)
			assert.Equal(t, genai.RoleModel, final.Content.Role)
			require.Len(t, final.Content.Parts, 1)
			assert.Equal(t, "Done", final.Content.Parts[0].Text)
		})
	}
}
