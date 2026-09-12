package bindings

import (
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadPendingToolIDsHandlesStoredStringSlice(t *testing.T) {
	state := map[string]any{
		"pending_client_tools": []string{"call_1", "call_2"},
	}
	assert.Equal(t, []string{"call_1", "call_2"}, readPendingToolIDs(state))
}

func TestReadPendingToolIDsHandlesJSONDecodedSlice(t *testing.T) {
	var state map[string]any
	require.NoError(t, json.Unmarshal([]byte(`{"pending_client_tools":["call_1","call_2"]}`), &state))
	assert.Equal(t, []string{"call_1", "call_2"}, readPendingToolIDs(state))
}

func TestGetProvidedPendingToolIDsHandlesToolOutput(t *testing.T) {
	provided := getProvidedPendingToolIDs([]chatMessage{{
		Role: "assistant",
		Parts: []chatPart{{
			Type:       "tool-render_card",
			State:      "output-available",
			ToolCallId: "call_render",
			ToolName:   "render_card",
		}},
	}}, []string{"call_render"})

	_, ok := provided["call_render"]
	assert.True(t, ok)
}

func TestGetLastMessageConvertsDynamicClientToolOutput(t *testing.T) {
	content, err := getLastMessage([]chatMessage{{
		Role: "assistant",
		Parts: []chatPart{{
			Type:       "dynamic-tool",
			State:      "output-available",
			ToolCallId: "call_add",
			ToolName:   "library_page_block_add",
			Output:     json.RawMessage(`{"_storyden_confirmation":{"approved":true},"message":"Added the assets block."}`),
		}},
	}}, []string{"call_add"}, slog.Default())
	require.NoError(t, err)
	require.Len(t, content.Parts, 1)
	require.NotNil(t, content.Parts[0].FunctionResponse)
	assert.Equal(t, "user", content.Role)
	assert.Equal(t, "call_add", content.Parts[0].FunctionResponse.ID)
	assert.Equal(t, "library_page_block_add", content.Parts[0].FunctionResponse.Name)
	assert.Equal(t, "Added the assets block.", content.Parts[0].FunctionResponse.Response["message"])
	assert.Equal(t, map[string]any{"approved": true}, content.Parts[0].FunctionResponse.Response["_storyden_confirmation"])
}

func TestGetLastMessageConvertsDynamicClientToolError(t *testing.T) {
	content, err := getLastMessage([]chatMessage{{
		Role: "assistant",
		Parts: []chatPart{{
			Type:       "dynamic-tool",
			State:      "output-error",
			ToolCallId: "call_add",
			ToolName:   "library_page_block_add",
			ErrorText:  "assets block already exists",
		}},
	}}, []string{"call_add"}, slog.Default())
	require.NoError(t, err)
	require.Len(t, content.Parts, 1)
	require.NotNil(t, content.Parts[0].FunctionResponse)
	assert.Equal(t, map[string]any{"error": "assets block already exists"}, content.Parts[0].FunctionResponse.Response)
}

func TestGetLastMessageRejectsUnsolicitedDynamicClientToolOutput(t *testing.T) {
	_, err := getLastMessage([]chatMessage{{
		Role: "assistant",
		Parts: []chatPart{{
			Type:       "dynamic-tool",
			State:      "output-available",
			ToolCallId: "call_add",
			ToolName:   "library_page_block_add",
			Output:     json.RawMessage(`{"message":"Added the assets block."}`),
		}},
	}}, nil, slog.Default())

	require.EqualError(t, err, "user message has no content")
}

func TestGetLastMessageRejectsUnsolicitedToolApproval(t *testing.T) {
	_, err := getLastMessage([]chatMessage{{
		Role: "assistant",
		Parts: []chatPart{{
			Type:       "tool-robot_delete",
			State:      "approval-responded",
			ToolCallId: "call_confirm",
			ToolName:   "robot_delete",
			Approval:   &chatApproval{ID: "call_confirm", Approved: true},
		}},
	}}, nil, slog.Default())

	require.EqualError(t, err, "user message has no content")
}

func TestGetLastMessageWrapsNonObjectDynamicClientToolOutput(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected any
	}{
		{name: "text", output: `"done"`, expected: "done"},
		{name: "scalar", output: `42`, expected: float64(42)},
		{name: "array", output: `["first","second"]`, expected: []any{"first", "second"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := getLastMessage([]chatMessage{{
				Role: "assistant",
				Parts: []chatPart{{
					Type:       "dynamic-tool",
					State:      "output-available",
					ToolCallId: "call_value",
					ToolName:   "browser_value_get",
					Output:     json.RawMessage(tt.output),
				}},
			}}, []string{"call_value"}, slog.Default())
			require.NoError(t, err)
			require.Len(t, content.Parts, 1)
			assert.Equal(t, "browser_value_get", content.Parts[0].FunctionResponse.Name)
			assert.Equal(t, map[string]any{"result": tt.expected}, content.Parts[0].FunctionResponse.Response)
		})
	}
}
