package robotprojection

import (
	"encoding/json"
	"testing"

	"google.golang.org/adk/v2/model"
	adksession "google.golang.org/adk/v2/session"
	"google.golang.org/adk/v2/tool/toolconfirmation"
	"google.golang.org/genai"
)

func TestFunctionResponseToUIPartDoesNotCopyOutputIntoInput(t *testing.T) {
	part, err := FunctionResponseToUIPart(&genai.FunctionResponse{
		ID:   "call_1",
		Name: "content_search",
		Response: map[string]any{
			"items":   []any{},
			"results": 0,
		},
	})
	if err != nil {
		t.Fatalf("FunctionResponseToUIPart() error = %v", err)
	}

	data, err := json.Marshal(part)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var got struct {
		Input  map[string]any `json:"input"`
		Output map[string]any `json:"output"`
	}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if len(got.Input) != 0 {
		t.Fatalf("input = %#v, want empty map", got.Input)
	}
	if got.Output["results"] != float64(0) {
		t.Fatalf("output.results = %#v, want 0", got.Output["results"])
	}
}

func TestADKEventToUIMessagePartsProjectsConfirmationAsApprovalRequested(t *testing.T) {
	event := adksession.Event{
		LLMResponse: model.LLMResponse{
			Content: &genai.Content{
				Parts: []*genai.Part{
					{
						FunctionCall: &genai.FunctionCall{
							ID:   "approval_1",
							Name: toolconfirmation.FunctionCallName,
							Args: map[string]any{
								"originalFunctionCall": &genai.FunctionCall{
									ID:   "call_1",
									Name: "discord_send",
									Args: map[string]any{"message": "hi"},
								},
								"toolConfirmation": toolconfirmation.ToolConfirmation{},
							},
						},
					},
				},
			},
		},
	}

	parts, err := ADKEventToUIMessageParts(event, nil, func(toolName string) map[string]any {
		return map[string]any{
			"storyden": map[string]any{
				"requires_confirmation": true,
			},
		}
	})
	if err != nil {
		t.Fatalf("ADKEventToUIMessageParts() error = %v", err)
	}
	if len(parts) != 1 {
		t.Fatalf("len(parts) = %d, want 1", len(parts))
	}

	data, err := json.Marshal(parts[0])
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var got struct {
		Type                 string         `json:"type"`
		ToolCallID           string         `json:"toolCallId"`
		ToolName             string         `json:"toolName"`
		State                string         `json:"state"`
		Approval             map[string]any `json:"approval"`
		CallProviderMetadata map[string]any `json:"callProviderMetadata"`
	}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if got.Type != "tool-discord_send" {
		t.Fatalf("type = %q, want tool-discord_send", got.Type)
	}
	if got.ToolCallID != "approval_1" {
		t.Fatalf("toolCallId = %q, want approval_1", got.ToolCallID)
	}
	if got.ToolName != "discord_send" {
		t.Fatalf("toolName = %q, want discord_send", got.ToolName)
	}
	if got.State != "approval-requested" {
		t.Fatalf("state = %q, want approval-requested", got.State)
	}
	if got.Approval["id"] != "approval_1" {
		t.Fatalf("approval.id = %#v, want approval_1", got.Approval["id"])
	}
	if got.CallProviderMetadata["storyden"] == nil {
		t.Fatalf("callProviderMetadata.storyden missing: %#v", got.CallProviderMetadata)
	}
}

func TestToolApprovalRequestStreamPart(t *testing.T) {
	part := ToolApprovalRequestStreamPart("call_1", "approval_1")

	data, err := json.Marshal(part)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var got struct {
		Type       string `json:"type"`
		ToolCallID string `json:"toolCallId"`
		ApprovalID string `json:"approvalId"`
	}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if got.Type != "tool-approval-request" {
		t.Fatalf("type = %q, want tool-approval-request", got.Type)
	}
	if got.ToolCallID != "call_1" {
		t.Fatalf("toolCallId = %q, want call_1", got.ToolCallID)
	}
	if got.ApprovalID != "approval_1" {
		t.Fatalf("approvalId = %q, want approval_1", got.ApprovalID)
	}
}
