package openai

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

func TestConvertToOpenAIInputReplaysEmptyToolArgsAsObject(t *testing.T) {
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

	input := convertToOpenAIInput(req)
	require.Len(t, input, 1)

	toolCall := input[0].OfFunctionCall
	require.NotNil(t, toolCall)
	assert.Equal(t, "{}", toolCall.Arguments)
}

func TestConvertToOpenAIInputPassesThroughToolResult(t *testing.T) {
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

	input := convertToOpenAIInput(req)
	require.Len(t, input, 1)
	require.NotNil(t, input[0].OfFunctionCallOutput)
	assert.Equal(t, "call_123", input[0].OfFunctionCallOutput.CallID)
	require.True(t, input[0].OfFunctionCallOutput.Output.OfString.Valid())
	assert.JSONEq(t, `{"id":"d8818ueot5pfij6bvm90","name":"Documentation Hub","slug":"documentation-hub"}`, input[0].OfFunctionCallOutput.Output.OfString.Value)
}
