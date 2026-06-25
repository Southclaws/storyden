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
	decision := reconcilePendingClientTools([]chatMessage{
		{
			Role: "user",
			Parts: []chatPart{
				{Type: "text", Text: "continue"},
			},
		},
	}, pendingClientTools{})

	assert.Empty(t, decision.Pending.IDs)
	assert.Empty(t, decision.Provided)
	_, blocked := decision.BlockingInteraction.Get()
	assert.False(t, blocked)
	_, recovered := decision.StaleRobotSwitch.Get()
	assert.False(t, recovered)
}

func TestReconcilePendingClientToolsUsesOwnerRobotForResolvedPendingTool(t *testing.T) {
	decision := reconcilePendingClientTools([]chatMessage{
		{
			Role: "assistant",
			Parts: []chatPart{
				{
					Type:       "tool-robot_switch",
					State:      "output-available",
					ToolCallId: "call_switch",
					ToolName:   "robot_switch",
					Output:     json.RawMessage(`{"success":true,"robot_id":"plugin_builder"}`),
				},
			},
		},
	}, pendingClientTools{
		IDs:    []string{"call_switch"},
		Robots: map[string]string{"call_switch": "robot_builder"},
	})

	_, ok := decision.Provided["call_switch"]
	assert.True(t, ok)
	owner, ok := decision.OwnerRobotID.Get()
	require.True(t, ok)
	assert.Equal(t, "robot_builder", owner)
	_, blocked := decision.BlockingInteraction.Get()
	assert.False(t, blocked)
}

func TestReconcilePendingClientToolsRecoversStaleRobotSwitch(t *testing.T) {
	decision := reconcilePendingClientTools([]chatMessage{
		{
			Role: "assistant",
			Parts: []chatPart{
				{
					Type:       "tool-robot_switch",
					State:      "input-available",
					ToolCallId: "call_switch",
					Input:      json.RawMessage(`{"robot_id":"plugin_builder"}`),
				},
			},
		},
		{
			Role: "user",
			Parts: []chatPart{
				{Type: "text", Text: "continue after reload"},
			},
		},
	}, pendingClientTools{
		IDs:    []string{"call_switch"},
		Robots: map[string]string{"call_switch": "robot_builder"},
	})

	recovery, ok := decision.StaleRobotSwitch.Get()
	require.True(t, ok)
	assert.Equal(t, "call_switch", recovery.ToolCallID)
	assert.Equal(t, "plugin_builder", recovery.RobotID)
	assert.Empty(t, decision.Provided)
	_, blocked := decision.BlockingInteraction.Get()
	assert.False(t, blocked)
}

func TestReconcilePendingClientToolsBlocksUnresolvedNonRecoverableTool(t *testing.T) {
	decision := reconcilePendingClientTools([]chatMessage{
		{
			Role: "user",
			Parts: []chatPart{
				{Type: "text", Text: "continue"},
			},
		},
	}, pendingClientTools{
		IDs:    []string{"call_confirm"},
		Robots: map[string]string{"call_confirm": "robot_builder"},
	})

	block, ok := decision.BlockingInteraction.Get()
	require.True(t, ok)
	assert.Equal(t, "pending_tool_interaction", block.Code)
	assert.Equal(t, http.StatusConflict, block.Status)
	assert.Equal(t, "pending tool interaction must be resolved before continuing", block.Message)
}

func TestReconcilePendingClientToolsBlocksPartialPendingResults(t *testing.T) {
	decision := reconcilePendingClientTools([]chatMessage{
		{
			Role: "assistant",
			Parts: []chatPart{
				{
					Type:       "tool-first",
					State:      "output-available",
					ToolCallId: "call_first",
					ToolName:   "first",
					Output:     json.RawMessage(`{"ok":true}`),
				},
			},
		},
	}, pendingClientTools{
		IDs: []string{"call_first", "call_second"},
	})

	block, ok := decision.BlockingInteraction.Get()
	require.True(t, ok)
	assert.Equal(t, "pending_tool_interaction", block.Code)
	assert.Equal(t, http.StatusConflict, block.Status)
	assert.Equal(t, "all pending tool interactions from the assistant turn must be resolved together", block.Message)
}

func TestClearPendingClientToolsPreservesOtherState(t *testing.T) {
	state := clearPendingClientTools(map[string]any{
		pendingClientToolsStateKey:      []string{"call_1"},
		pendingClientToolRobotsStateKey: map[string]string{"call_1": "robot_builder"},
		"current_robot_id":              "plugin_builder",
		"robot_workspace":               map[string]any{"provider": "local"},
	})

	assert.NotContains(t, state, pendingClientToolsStateKey)
	assert.NotContains(t, state, pendingClientToolRobotsStateKey)
	assert.Equal(t, "plugin_builder", state["current_robot_id"])
	assert.Equal(t, map[string]any{"provider": "local"}, state["robot_workspace"])
}

func TestRobotSwitchToolResultContentBuildsClientResultEventContent(t *testing.T) {
	content := robotSwitchToolResultContent("call_switch", "plugin_builder")

	require.Equal(t, "user", content.Role)
	require.Len(t, content.Parts, 1)
	response := content.Parts[0].FunctionResponse
	require.NotNil(t, response)
	assert.Equal(t, "call_switch", response.ID)
	assert.Equal(t, "robot_switch", response.Name)
	assert.Equal(t, true, response.Response["success"])
	assert.Equal(t, "plugin_builder", response.Response["robot_id"])
}

func TestWriteChatErrorReturnsJSONEnvelope(t *testing.T) {
	rec := httptest.NewRecorder()

	writeChatError(rec, chatError{
		Code:    "pending_tool_interaction",
		Message: "pending tool interaction must be resolved before continuing",
		Status:  http.StatusConflict,
	})

	require.Equal(t, http.StatusConflict, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "pending_tool_interaction", body["error"])
	assert.Equal(t, "pending tool interaction must be resolved before continuing", body["message"])
	assert.Equal(t, float64(http.StatusConflict), body["status"])
}
