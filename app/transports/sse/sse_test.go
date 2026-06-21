package sse

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

func TestGetRobotSwitchTargetIDFromPendingToolOutput(t *testing.T) {
	output := json.RawMessage(`{"success":true,"robot_id":"plugin_builder"}`)
	target := getRobotSwitchTargetID([]chatMessage{
		{
			Role: "assistant",
			Parts: []chatPart{
				{
					Type:       "tool-robot_switch",
					State:      "output-available",
					ToolCallId: "call_switch",
					ToolName:   "robot_switch",
					Output:     output,
				},
			},
		},
	}, []string{"call_switch"})

	got, ok := target.Get()
	require.True(t, ok)
	assert.Equal(t, "plugin_builder", got)
}

func TestGetRobotSwitchTargetIDInfersToolNameFromPartType(t *testing.T) {
	output := json.RawMessage(`{"success":true,"robot_id":"plugin_builder"}`)
	target := getRobotSwitchTargetID([]chatMessage{
		{
			Role: "assistant",
			Parts: []chatPart{
				{
					Type:       "tool-robot_switch",
					State:      "output-available",
					ToolCallId: "call_switch",
					Output:     output,
				},
			},
		},
	}, []string{"call_switch"})

	got, ok := target.Get()
	require.True(t, ok)
	assert.Equal(t, "plugin_builder", got)
}

func TestGetProvidedPendingToolIDsHandlesPendingToolOutput(t *testing.T) {
	provided := getProvidedPendingToolIDs([]chatMessage{
		{
			Role: "assistant",
			Parts: []chatPart{
				{
					Type:       "tool-robot_switch",
					State:      "output-available",
					ToolCallId: "call_switch",
					ToolName:   "robot_switch",
				},
			},
		},
	}, []string{"call_switch"})

	_, ok := provided["call_switch"]
	assert.True(t, ok)
}

func TestGetPendingRobotSwitchInputTargetIDFindsStalePendingSwitch(t *testing.T) {
	input := json.RawMessage(`{"robot_id":"plugin_builder"}`)
	target := getPendingRobotSwitchInputTargetID([]chatMessage{
		{
			Role: "assistant",
			Parts: []chatPart{
				{
					Type:       "tool-robot_switch",
					State:      "input-available",
					ToolCallId: "call_switch",
					Input:      input,
				},
			},
		},
		{
			Role: "user",
			Parts: []chatPart{
				{
					Type: "text",
					Text: "continue after reload",
				},
			},
		},
	}, []string{"call_switch"})

	got, ok := target.Get()
	require.True(t, ok)
	assert.Equal(t, "plugin_builder", got)
}

func TestGetPendingRobotSwitchInputTargetIDDoesNotRecoverNonSwitchTool(t *testing.T) {
	input := json.RawMessage(`{"path":"main.go"}`)
	target := getPendingRobotSwitchInputTargetID([]chatMessage{
		{
			Role: "assistant",
			Parts: []chatPart{
				{
					Type:       "tool-plugin_file_write",
					State:      "input-available",
					ToolCallId: "call_write",
					Input:      input,
				},
			},
		},
	}, []string{"call_write"})

	_, ok := target.Get()
	assert.False(t, ok)
}
