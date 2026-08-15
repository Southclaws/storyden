package openai

import (
	"encoding/json"
	"testing"

	"github.com/openai/openai-go/v3/responses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/adk/v2/model"
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
	assert.Equal(t, "call_123", toolCall.CallID)
	assert.Equal(t, "library_request_page", toolCall.Name)
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

func TestConvertOpenAIResponseToGenaiContentPreservesReasoningItems(t *testing.T) {
	var output []responses.ResponseOutputItemUnion
	require.NoError(t, json.Unmarshal([]byte(`[
		{"type":"reasoning","id":"rs_123","summary":[],"encrypted_content":"encrypted-reasoning"},
		{"type":"function_call","call_id":"call_123","name":"library_request_page","arguments":"{\"page\":1}"},
		{"type":"message","id":"msg_123","role":"assistant","status":"completed","content":[{"type":"output_text","text":"Looking it up."}]}
	]`), &output))

	content, err := convertOpenAIResponseToGenaiContent(responses.Response{Output: output})
	require.NoError(t, err)
	require.Len(t, content.Parts, 3)
	assert.JSONEq(t, `{"type":"reasoning","id":"rs_123","summary":[],"encrypted_content":"encrypted-reasoning"}`, reasoningMetadata(t, content.Parts[0]))
	require.NotNil(t, content.Parts[1].FunctionCall)
	assert.Equal(t, map[string]any{"page": float64(1)}, content.Parts[1].FunctionCall.Args)
	assert.Equal(t, "Looking it up.", content.Parts[2].Text)

	input := convertToOpenAIInput(&model.LLMRequest{Contents: []*genai.Content{content}})
	require.Len(t, input, 3)
	require.NotNil(t, input[0].OfReasoning)
	assert.Equal(t, "rs_123", input[0].OfReasoning.ID)
	require.True(t, input[0].OfReasoning.EncryptedContent.Valid())
	assert.Equal(t, "encrypted-reasoning", input[0].OfReasoning.EncryptedContent.Value)
	require.NotNil(t, input[1].OfFunctionCall)
	assert.Equal(t, "call_123", input[1].OfFunctionCall.CallID)
	require.NotNil(t, input[2].OfMessage)
}

func TestConvertOpenAIResponseToGenaiContentRejectsMalformedFunctionArguments(t *testing.T) {
	var output []responses.ResponseOutputItemUnion
	require.NoError(t, json.Unmarshal([]byte(`[
		{"type":"function_call","call_id":"call_123","name":"library_request_page","arguments":"{not valid"}
	]`), &output))

	_, err := convertOpenAIResponseToGenaiContent(responses.Response{Output: output})
	require.Error(t, err)
	assert.ErrorContains(t, err, `failed to parse arguments for function call "library_request_page"`)
}

func TestConvertToOpenAIToolsUsesObjectSchemaWhenDeclarationHasNone(t *testing.T) {
	req := &model.LLMRequest{Config: &genai.GenerateContentConfig{Tools: []*genai.Tool{{
		FunctionDeclarations: []*genai.FunctionDeclaration{{Name: "library_request_page"}},
	}}}}

	tools := convertToOpenAITools(req)
	require.Len(t, tools, 1)
	require.NotNil(t, tools[0].OfFunction)
	assert.Equal(t, map[string]any{"type": "object", "properties": map[string]any{}}, tools[0].OfFunction.Parameters)
}

func TestConvertToOpenAIToolsNormalizesADKAgentSchema(t *testing.T) {
	req := &model.LLMRequest{Config: &genai.GenerateContentConfig{Tools: []*genai.Tool{{
		FunctionDeclarations: []*genai.FunctionDeclaration{{
			Name: "robot_d9n7eodo2dtgv230rpng",
			Parameters: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"request": {Type: genai.TypeString},
				},
				Required: []string{"request"},
			},
		}},
	}}}}

	tools := convertToOpenAITools(req)
	require.Len(t, tools, 1)
	require.NotNil(t, tools[0].OfFunction)
	assert.Equal(t, map[string]any{
		"type": "object",
		"properties": map[string]any{
			"request": map[string]any{"type": "string"},
		},
		"required": []any{"request"},
	}, tools[0].OfFunction.Parameters)
}

func TestConvertOpenAIResponseToGenaiFinishReason(t *testing.T) {
	assert.Equal(t, genai.FinishReasonStop, convertOpenAIResponseToGenaiFinishReason(responses.Response{Status: responses.ResponseStatusCompleted}))
	assert.Equal(t, genai.FinishReasonMaxTokens, convertOpenAIResponseToGenaiFinishReason(responses.Response{
		Status: responses.ResponseStatusIncomplete, IncompleteDetails: responses.ResponseIncompleteDetails{Reason: "max_output_tokens"},
	}))
	assert.Equal(t, genai.FinishReasonSafety, convertOpenAIResponseToGenaiFinishReason(responses.Response{
		Status: responses.ResponseStatusIncomplete, IncompleteDetails: responses.ResponseIncompleteDetails{Reason: "content_filter"},
	}))
	assert.Equal(t, genai.FinishReasonUnspecified, convertOpenAIResponseToGenaiFinishReason(responses.Response{Status: responses.ResponseStatusFailed}))
}

func reasoningMetadata(t *testing.T, part *genai.Part) string {
	t.Helper()
	require.NotNil(t, part)
	raw, ok := part.PartMetadata[openAIReasoningItemKey].(string)
	require.True(t, ok)
	return raw
}
