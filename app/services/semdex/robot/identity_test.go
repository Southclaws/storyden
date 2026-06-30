package robot

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

func TestRobotIdentityInstructionIncludesCurrentRobot(t *testing.T) {
	text := robotIdentityInstruction(robotIdentityContext{
		Current: robotIdentity{
			Name:         "Content Moderator",
			Description:  "Reviews posts and replies.",
			Capabilities: []string{"thread_reply", "robot_switch"},
		},
		Participants: []robotParticipant{
			{Name: "Content Moderator", Active: true},
		},
	})

	assert.Contains(t, text, "## Current Robot")
	assert.Contains(t, text, "Name: Content Moderator")
	assert.Contains(t, text, "Description: Reviews posts and replies.")
	assert.Contains(t, text, "- robot_switch")
	assert.Contains(t, text, "This identity is authoritative for the current turn.")
	assert.NotContains(t, text, "## Session Robot Context")
}

func TestRobotIdentityInstructionIncludesParticipantsForMultiRobotSessions(t *testing.T) {
	text := robotIdentityInstruction(robotIdentityContext{
		Current: robotIdentity{Name: "Content Moderator"},
		Participants: []robotParticipant{
			{Name: "Content Moderator", Active: true},
			{Name: "Storyden Robot Builder"},
			{Name: "Library Curator"},
		},
	})

	assert.Contains(t, text, "## Session Robot Context")
	assert.Contains(t, text, "- Content Moderator - active for this turn")
	assert.Contains(t, text, "- Storyden Robot Builder - previously active")
	assert.Contains(t, text, "- Library Curator - previously active")
	assert.Contains(t, text, "Robot-switch markers in the conversation history indicate where the active Robot changed.")
}

func TestRobotIdentityInstructionIncludesUnavailableConfiguredTools(t *testing.T) {
	text := robotIdentityInstruction(robotIdentityContext{
		Current: robotIdentity{
			Name:             "Content Moderator",
			Capabilities:     []string{"thread_reply"},
			UnavailableTools: []string{"mcp:search:gone"},
		},
		Participants: []robotParticipant{
			{Name: "Content Moderator", Active: true},
		},
	})

	assert.Contains(t, text, "Configured tools currently unavailable:")
	assert.Contains(t, text, "- mcp:search:gone")
	assert.Contains(t, text, "this Robot's toolset has changed")
}

func TestProjectRobotSwitchesRewritesCompletedSwitchToolProtocol(t *testing.T) {
	req := &model.LLMRequest{
		Contents: []*genai.Content{
			genai.NewContentFromText("Please switch.", genai.RoleUser),
			genai.NewContentFromParts([]*genai.Part{
				{
					FunctionCall: &genai.FunctionCall{
						ID:   "call_switch_1",
						Name: "robot_switch",
						Args: map[string]any{"robot_id": "robot_123"},
					},
				},
			}, genai.RoleModel),
			genai.NewContentFromParts([]*genai.Part{
				{
					FunctionResponse: &genai.FunctionResponse{
						ID:   "call_switch_1",
						Name: "robot_switch",
						Response: map[string]any{
							"success":  true,
							"robot_id": "robot_123",
						},
					},
				},
			}, genai.RoleUser),
			genai.NewContentFromText("Done.", genai.RoleModel),
		},
	}

	count := projectRobotSwitches(context.Background(), req, func(context.Context, string) string {
		return "Content Moderator"
	})

	require.Equal(t, 1, count)
	require.Len(t, req.Contents, 4)
	require.Len(t, req.Contents[1].Parts, 1)
	assert.Equal(t, "Robot switch requested.", req.Contents[1].Parts[0].Text)
	assert.Equal(t, genai.RoleModel, req.Contents[1].Role)
	require.Len(t, req.Contents[2].Parts, 1)
	assert.Equal(t, genai.RoleUser, req.Contents[2].Role)
	assert.Nil(t, req.Contents[2].Parts[0].FunctionResponse)
	assert.Contains(t, req.Contents[2].Parts[0].Text, "ROBOT SWITCH")
	assert.Contains(t, req.Contents[2].Parts[0].Text, "Current Robot after this switch: Content Moderator")
	assertNoRobotSwitchToolParts(t, req)
}

func TestProjectRobotSwitchesLeavesPendingSwitchProtocolUntouched(t *testing.T) {
	req := &model.LLMRequest{
		Contents: []*genai.Content{
			genai.NewContentFromParts([]*genai.Part{
				{
					FunctionCall: &genai.FunctionCall{
						ID:   "call_switch_1",
						Name: "robot_switch",
						Args: map[string]any{"robot_id": "robot_123"},
					},
				},
			}, genai.RoleModel),
		},
	}

	count := projectRobotSwitches(context.Background(), req, nil)

	require.Equal(t, 0, count)
	require.NotNil(t, req.Contents[0].Parts[0].FunctionCall)
	assert.Equal(t, "robot_switch", req.Contents[0].Parts[0].FunctionCall.Name)
}

func TestProjectRobotSwitchesRewritesADKContextTextForOtherAgents(t *testing.T) {
	req := &model.LLMRequest{
		Contents: []*genai.Content{
			genai.NewContentFromText("For context:\n[Switcher] called tool `robot_switch` with parameters: {\"robot_id\":\"robot_123\"}", genai.RoleUser),
			genai.NewContentFromText("What next?", genai.RoleUser),
		},
	}

	count := projectRobotSwitches(context.Background(), req, func(context.Context, string) string {
		return "Content Moderator"
	})

	require.Equal(t, 1, count)
	require.Len(t, req.Contents, 2)
	assert.Contains(t, req.Contents[0].Parts[0].Text, "ROBOT SWITCH")
	assert.Contains(t, req.Contents[0].Parts[0].Text, "Current Robot after this switch: Content Moderator")
	assert.NotContains(t, req.Contents[0].Parts[0].Text, "called tool `robot_switch`")
	assert.Equal(t, "What next?", req.Contents[1].Parts[0].Text)
}

func TestProjectRobotSwitchesLeavesFailedSwitchProtocolUntouched(t *testing.T) {
	req := &model.LLMRequest{
		Contents: []*genai.Content{
			genai.NewContentFromParts([]*genai.Part{
				{
					FunctionCall: &genai.FunctionCall{
						ID:   "call_switch_1",
						Name: "robot_switch",
						Args: map[string]any{"robot_id": "robot_123"},
					},
				},
			}, genai.RoleModel),
			genai.NewContentFromParts([]*genai.Part{
				{
					FunctionResponse: &genai.FunctionResponse{
						ID:   "call_switch_1",
						Name: "robot_switch",
						Response: map[string]any{
							"success":  false,
							"robot_id": "robot_123",
						},
					},
				},
			}, genai.RoleUser),
		},
	}

	count := projectRobotSwitches(context.Background(), req, nil)

	require.Equal(t, 0, count)
	require.NotNil(t, req.Contents[0].Parts[0].FunctionCall)
	require.NotNil(t, req.Contents[1].Parts[0].FunctionResponse)
}

func assertNoRobotSwitchToolParts(t *testing.T, req *model.LLMRequest) {
	t.Helper()

	for _, content := range req.Contents {
		for _, part := range content.Parts {
			if part.FunctionCall != nil {
				assert.NotEqual(t, "robot_switch", part.FunctionCall.Name)
			}
			if part.FunctionResponse != nil {
				assert.NotEqual(t, "robot_switch", part.FunctionResponse.Name)
			}
			assert.False(t, strings.Contains(part.Text, "tool_response"))
		}
	}
}
