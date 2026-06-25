package openai

import (
	"encoding/json"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

func TestConvertToOpenAIMessagesReplaysEmptyToolArgsAsObject(t *testing.T) {
	req := &model.LLMRequest{
		Contents: []*genai.Content{
			{
				Role: genai.RoleModel,
				Parts: []*genai.Part{
					{
						FunctionCall: &genai.FunctionCall{
							ID:   "call_123",
							Name: "library_request_page",
						},
					},
				},
			},
		},
	}

	messages := convertToOpenAIMessages(req)
	require.Len(t, messages, 1)
	require.NotNil(t, messages[0].OfAssistant)
	require.Len(t, messages[0].OfAssistant.ToolCalls, 1)

	toolCall := messages[0].OfAssistant.ToolCalls[0].OfFunction
	require.NotNil(t, toolCall)
	assert.Equal(t, "{}", toolCall.Function.Arguments)
}

func TestConvertToOpenAIMessagesPassesThroughToolResult(t *testing.T) {
	req := &model.LLMRequest{
		Contents: []*genai.Content{
			{
				Role: genai.RoleUser,
				Parts: []*genai.Part{
					{
						FunctionResponse: &genai.FunctionResponse{
							ID:   "call_123",
							Name: "library_request_page",
							Response: map[string]any{
								"id":   "d8818ueot5pfij6bvm90",
								"name": "Documentation Hub",
								"slug": "documentation-hub",
							},
						},
					},
				},
			},
		},
	}

	messages := convertToOpenAIMessages(req)
	require.Len(t, messages, 1)
	require.NotNil(t, messages[0].OfTool)
	assert.Equal(t, "call_123", messages[0].OfTool.ToolCallID)

	content := toolMessageContent(t, messages[0].OfTool.Content)
	var payload map[string]any
	require.NoError(t, json.Unmarshal([]byte(content), &payload))

	assert.Equal(t, "Documentation Hub", payload["name"])
	assert.Equal(t, "documentation-hub", payload["slug"])
}

func toolMessageContent(t *testing.T, content openai.ChatCompletionToolMessageParamContentUnion) string {
	t.Helper()

	raw, err := json.Marshal(content)
	require.NoError(t, err)

	var value string
	require.NoError(t, json.Unmarshal(raw, &value))
	return value
}
