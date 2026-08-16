package bindings

import (
	"encoding/json"
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
