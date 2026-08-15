package sse

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReconcilePendingClientToolsAllowsNormalRequestWithoutPending(t *testing.T) {
	decision := reconcilePendingClientTools([]chatMessage{{
		Role:  "user",
		Parts: []chatPart{{Type: "text", Text: "continue"}},
	}}, pendingClientTools{})

	assert.Empty(t, decision.Pending.IDs)
	assert.Empty(t, decision.Provided)
	_, blocked := decision.BlockingInteraction.Get()
	assert.False(t, blocked)
}

func TestReconcilePendingClientToolsAcceptsResolvedApproval(t *testing.T) {
	decision := reconcilePendingClientTools([]chatMessage{{
		Role: "assistant",
		Parts: []chatPart{{
			Type:       "tool-robot_delete",
			State:      "approval-responded",
			ToolCallId: "call_confirm",
			ToolName:   "robot_delete",
			Approval:   &chatApproval{ID: "call_confirm", Approved: true},
		}},
	}}, pendingClientTools{IDs: []string{"call_confirm"}})

	_, ok := decision.Provided["call_confirm"]
	assert.True(t, ok)
	_, blocked := decision.BlockingInteraction.Get()
	assert.False(t, blocked)
}

func TestReconcilePendingClientToolsBlocksUnresolvedTool(t *testing.T) {
	decision := reconcilePendingClientTools([]chatMessage{{
		Role:  "user",
		Parts: []chatPart{{Type: "text", Text: "continue"}},
	}}, pendingClientTools{IDs: []string{"call_confirm"}})

	block, ok := decision.BlockingInteraction.Get()
	require.True(t, ok)
	assert.Equal(t, "pending_tool_interaction", block.Code)
	assert.Equal(t, http.StatusConflict, block.Status)
}

func TestReconcilePendingClientToolsBlocksPartialPendingResults(t *testing.T) {
	decision := reconcilePendingClientTools([]chatMessage{{
		Role: "assistant",
		Parts: []chatPart{{
			Type:       "tool-first",
			State:      "output-available",
			ToolCallId: "call_first",
			ToolName:   "first",
		}},
	}}, pendingClientTools{IDs: []string{"call_first", "call_second"}})

	block, ok := decision.BlockingInteraction.Get()
	require.True(t, ok)
	assert.Equal(t, "all pending tool interactions from the assistant turn must be resolved together", block.Message)
}

func TestClearPendingClientToolsPreservesOtherState(t *testing.T) {
	state := clearPendingClientTools(map[string]any{
		pendingClientToolsStateKey: []string{"call_1"},
		"robot_workspace":          map[string]any{"provider": "local"},
	})

	assert.NotContains(t, state, pendingClientToolsStateKey)
	assert.Equal(t, map[string]any{"provider": "local"}, state["robot_workspace"])
}

func TestWriteChatErrorReturnsJSONEnvelope(t *testing.T) {
	recorder := httptest.NewRecorder()
	writeChatError(recorder, chatError{
		Code:    "pending_tool_interaction",
		Message: "pending tool interaction must be resolved before continuing",
		Status:  http.StatusConflict,
	})

	require.Equal(t, http.StatusConflict, recorder.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
	assert.Equal(t, "pending_tool_interaction", body["error"])
}
