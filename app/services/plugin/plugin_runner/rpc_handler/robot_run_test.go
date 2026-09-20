package rpc_handler

import (
	"testing"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/adk/v2/model"
	adksession "google.golang.org/adk/v2/session"
	"google.golang.org/genai"

	robotresource "github.com/Southclaws/storyden/app/resources/robot"
	robotservice "github.com/Southclaws/storyden/app/services/semdex/robot"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func TestValidateRobotRunRequestBuildsRolePreservingContents(t *testing.T) {
	firstImage := xid.New()
	secondImage := xid.New()
	contents, err := validateRobotRunRequest(rpc.RPCRequestRobotRunParams{
		Mode:    rpc.RobotRunModeConversation,
		RobotID: "robot-one",
		Messages: []rpc.RobotRunMessage{
			{
				Role:    rpc.RobotRunMessageRoleUser,
				Content: "first",
				Author:  opt.New("Alice"),
			},
			{
				Role:    rpc.RobotRunMessageRoleAssistant,
				Content: "second",
			},
			{
				Role:    rpc.RobotRunMessageRoleUser,
				Content: "third",
				Author:  opt.New("Bob"),
				Media: []rpc.RobotRunMedia{
					{Type: rpc.RobotRunMediaTypeImage, AssetID: firstImage},
					{Type: rpc.RobotRunMediaTypeImage, AssetID: secondImage},
				},
			},
		},
	})

	require.NoError(t, err)
	require.Len(t, contents, 3)
	assert.Equal(t, []string{genai.RoleUser, genai.RoleModel, genai.RoleUser}, []string{contents[0].Role, contents[1].Role, contents[2].Role})
	assert.Equal(t, []string{"first", "second", "third"}, []string{contents[0].Parts[0].Text, contents[1].Parts[0].Text, contents[2].Parts[0].Text})
	assert.Equal(t, "Alice", contents[0].Parts[0].PartMetadata[robotservice.MessageSpeakerMetadataKey])
	assert.Equal(t, "Bob", contents[2].Parts[0].PartMetadata[robotservice.MessageSpeakerMetadataKey])
	require.Len(t, contents[2].Parts, 3)
	assert.Equal(t, firstImage.String(), contents[2].Parts[1].PartMetadata[robotresource.ImageAssetIDMetadataKey])
	assert.Equal(t, secondImage.String(), contents[2].Parts[2].PartMetadata[robotresource.ImageAssetIDMetadataKey])
}

func TestValidateRobotRunRequestRejectsMediaOnAssistantMessages(t *testing.T) {
	_, err := validateRobotRunRequest(rpc.RPCRequestRobotRunParams{
		Mode:    rpc.RobotRunModeConversation,
		RobotID: "robot-one",
		Messages: []rpc.RobotRunMessage{
			{
				Role:    rpc.RobotRunMessageRoleAssistant,
				Content: "prior response",
				Media: []rpc.RobotRunMedia{{
					Type:    rpc.RobotRunMediaTypeImage,
					AssetID: xid.New(),
				}},
			},
			{Role: rpc.RobotRunMessageRoleUser, Content: "continue"},
		},
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "media is only supported for user messages")
}

func TestValidateRobotRunRequestEnforcesModeContracts(t *testing.T) {
	validUser := rpc.RobotRunMessage{
		Role:    rpc.RobotRunMessageRoleUser,
		Content: "run",
	}
	tests := map[string]rpc.RPCRequestRobotRunParams{
		"conversation final assistant": {
			Mode:    rpc.RobotRunModeConversation,
			RobotID: "robot-one",
			Messages: []rpc.RobotRunMessage{
				{
					Role:    rpc.RobotRunMessageRoleAssistant,
					Content: "reply",
				},
			},
		},
		"automation history": {
			Mode:     rpc.RobotRunModeAutomation,
			RobotID:  "robot-one",
			Messages: []rpc.RobotRunMessage{validUser, validUser},
		},
		"automation continuation": {
			Mode:      rpc.RobotRunModeAutomation,
			RobotID:   "robot-one",
			Messages:  []rpc.RobotRunMessage{validUser},
			SessionID: opt.New(xid.New()),
		},
		"blank author": {
			Mode:    rpc.RobotRunModeConversation,
			RobotID: "robot-one",
			Messages: []rpc.RobotRunMessage{
				{
					Role:    rpc.RobotRunMessageRoleUser,
					Content: "run",
					Author:  opt.New("  "),
				},
			},
		},
	}

	for name, params := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := validateRobotRunRequest(params)
			assert.Error(t, err)
		})
	}
}

func TestAppendRobotRunEventsOmitsReasoningAndKeepsPartOrder(t *testing.T) {
	result := rpc.RPCResponseRobotRun{
		Mode: rpc.RobotRunModeConversation,
	}
	appendRobotRunEvents(&result, &adksession.Event{
		LLMResponse: model.LLMResponse{
			Content: &genai.Content{
				Role: genai.RoleModel,
				Parts: []*genai.Part{
					{
						Text:    "private",
						Thought: true,
					},
					{
						Text: "visible",
					},
					{
						FunctionCall: &genai.FunctionCall{
							ID:   "call-1",
							Name: "search",
							Args: map[string]any{
								"query": "robots",
							},
						},
					},
					{
						FunctionResponse: &genai.FunctionResponse{
							ID:   "call-1",
							Name: "search",
							Response: map[string]any{
								"count": 2,
							},
						},
					},
				},
			},
		},
	})

	assert.Equal(t, "visible", result.FinalText.OrZero())
	require.Len(t, result.Events, 3)
	text, ok := result.Events[0].RobotRunEventUnion.(*rpc.RobotRunTextEvent)
	require.True(t, ok)
	assert.Equal(t, "visible", text.Text)
	call, ok := result.Events[1].RobotRunEventUnion.(*rpc.RobotRunToolCallEvent)
	require.True(t, ok)
	assert.Equal(t, "call-1", call.CallID.OrZero())
	assert.Equal(t, map[string]any{"query": "robots"}, call.Arguments)
	toolResult, ok := result.Events[2].RobotRunEventUnion.(*rpc.RobotRunToolResultEvent)
	require.True(t, ok)
	assert.Equal(t, map[string]any{"count": 2}, toolResult.Result)
}
