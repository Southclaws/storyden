package agent_history

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

func TestRepairInterruptedToolCallsInsertsResultBeforeNextUserMessage(t *testing.T) {
	req := &model.LLMRequest{
		Contents: []*genai.Content{
			modelToolCall("call_1", "plugin_install"),
			userText("try again"),
		},
	}

	count := RepairInterruptedToolCalls(req)
	require.Equal(t, 1, count)
	require.Len(t, req.Contents, 3)

	response := onlyFunctionResponse(t, req.Contents[1])
	assert.Equal(t, "call_1", response.ID)
	assert.Equal(t, "plugin_install", response.Name)
	assert.Equal(t, interruptedToolResultStatus, response.Response["status"])
	assert.Equal(t, true, response.Response["interrupted"])
	assert.Equal(t, "try again", req.Contents[2].Parts[0].Text)
}

func TestRepairInterruptedToolCallsAppendsResultForTrailingToolCall(t *testing.T) {
	req := &model.LLMRequest{
		Contents: []*genai.Content{
			modelToolCall("call_1", "plugin_install"),
		},
	}

	count := RepairInterruptedToolCalls(req)
	require.Equal(t, 1, count)
	require.Len(t, req.Contents, 2)

	response := onlyFunctionResponse(t, req.Contents[1])
	assert.Equal(t, "call_1", response.ID)
}

func TestRepairInterruptedToolCallsLeavesCompletedCallAlone(t *testing.T) {
	req := &model.LLMRequest{
		Contents: []*genai.Content{
			modelToolCall("call_1", "plugin_install"),
			toolResponse("call_1", "plugin_install", map[string]any{"success": true}),
			userText("thanks"),
		},
	}

	count := RepairInterruptedToolCalls(req)
	require.Equal(t, 0, count)
	require.Len(t, req.Contents, 3)
	assert.Equal(t, "thanks", req.Contents[2].Parts[0].Text)
}

func TestRepairInterruptedToolCallsAddsMissingResultToPartialToolResultMessage(t *testing.T) {
	req := &model.LLMRequest{
		Contents: []*genai.Content{
			{
				Role: genai.RoleModel,
				Parts: []*genai.Part{
					{FunctionCall: &genai.FunctionCall{ID: "call_1", Name: "plugin_validate"}},
					{FunctionCall: &genai.FunctionCall{ID: "call_2", Name: "plugin_install"}},
				},
			},
			toolResponse("call_1", "plugin_validate", map[string]any{"success": true}),
		},
	}

	count := RepairInterruptedToolCalls(req)
	require.Equal(t, 1, count)
	require.Len(t, req.Contents, 2)
	require.Len(t, req.Contents[1].Parts, 2)

	assert.Equal(t, "call_1", req.Contents[1].Parts[0].FunctionResponse.ID)
	assert.Equal(t, "call_2", req.Contents[1].Parts[1].FunctionResponse.ID)
	assert.Equal(t, interruptedToolResultStatus, req.Contents[1].Parts[1].FunctionResponse.Response["status"])
}

func TestRepairInterruptedToolCallsHandlesMultipleInterruptedCalls(t *testing.T) {
	req := &model.LLMRequest{
		Contents: []*genai.Content{
			{
				Role: genai.RoleModel,
				Parts: []*genai.Part{
					{FunctionCall: &genai.FunctionCall{ID: "call_1", Name: "plugin_validate"}},
					{FunctionCall: &genai.FunctionCall{ID: "call_2", Name: "plugin_install"}},
				},
			},
			userText("resume"),
		},
	}

	count := RepairInterruptedToolCalls(req)
	require.Equal(t, 2, count)
	require.Len(t, req.Contents, 3)
	require.Len(t, req.Contents[1].Parts, 2)
	assert.Equal(t, "call_1", req.Contents[1].Parts[0].FunctionResponse.ID)
	assert.Equal(t, "call_2", req.Contents[1].Parts[1].FunctionResponse.ID)
}

func modelToolCall(id, name string) *genai.Content {
	return &genai.Content{
		Role: genai.RoleModel,
		Parts: []*genai.Part{
			{
				FunctionCall: &genai.FunctionCall{
					ID:   id,
					Name: name,
				},
			},
		},
	}
}

func toolResponse(id, name string, response map[string]any) *genai.Content {
	return &genai.Content{
		Role: genai.RoleUser,
		Parts: []*genai.Part{
			{
				FunctionResponse: &genai.FunctionResponse{
					ID:       id,
					Name:     name,
					Response: response,
				},
			},
		},
	}
}

func userText(text string) *genai.Content {
	return &genai.Content{
		Role:  genai.RoleUser,
		Parts: []*genai.Part{{Text: text}},
	}
}

func onlyFunctionResponse(t *testing.T, content *genai.Content) *genai.FunctionResponse {
	t.Helper()
	require.NotNil(t, content)
	require.Len(t, content.Parts, 1)
	require.NotNil(t, content.Parts[0].FunctionResponse)
	return content.Parts[0].FunctionResponse
}
