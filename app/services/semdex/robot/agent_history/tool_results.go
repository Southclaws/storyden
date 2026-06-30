package agent_history

import (
	"log/slog"

	agentpkg "google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

const interruptedToolResultStatus = "interrupted"

type pendingToolCall struct {
	id   string
	name string
}

func RepairInterruptedToolCallsBeforeModel(logger *slog.Logger) llmagent.BeforeModelCallback {
	if logger == nil {
		logger = slog.Default()
	}

	return func(ctx agentpkg.Context, req *model.LLMRequest) (*model.LLMResponse, error) {
		repaired := RepairInterruptedToolCalls(req)
		if repaired > 0 {
			logger.Warn("repaired interrupted tool call history",
				slog.String("agent", ctx.AgentName()),
				slog.String("invocation", ctx.InvocationID()),
				slog.String("session", ctx.SessionID()),
				slog.Int("count", repaired),
			)
		}
		return nil, nil
	}
}

func RepairInterruptedToolCalls(req *model.LLMRequest) int {
	if req == nil {
		return 0
	}

	var repaired int
	var pending []pendingToolCall
	out := make([]*genai.Content, 0, len(req.Contents))

	for _, content := range req.Contents {
		if content == nil {
			continue
		}

		if len(pending) > 0 {
			if hasFunctionResponses(content) {
				missing := missingToolResponses(pending, content)
				if len(missing) > 0 {
					content.Parts = append(content.Parts, syntheticToolResponseParts(missing)...)
					repaired += len(missing)
				}
				pending = nil
				out = append(out, content)
				continue
			}

			out = append(out, syntheticToolResponseContent(pending))
			repaired += len(pending)
			pending = nil
		}

		out = append(out, content)

		calls := extractToolCalls(content)
		if len(calls) > 0 {
			pending = calls
		}
	}

	if len(pending) > 0 {
		out = append(out, syntheticToolResponseContent(pending))
		repaired += len(pending)
	}

	if repaired > 0 {
		req.Contents = out
	}

	return repaired
}

func extractToolCalls(content *genai.Content) []pendingToolCall {
	if content == nil || content.Role != genai.RoleModel {
		return nil
	}

	calls := make([]pendingToolCall, 0, len(content.Parts))
	for _, part := range content.Parts {
		if part == nil || part.FunctionCall == nil || part.FunctionCall.ID == "" {
			continue
		}
		calls = append(calls, pendingToolCall{
			id:   part.FunctionCall.ID,
			name: part.FunctionCall.Name,
		})
	}
	return calls
}

func hasFunctionResponses(content *genai.Content) bool {
	if content == nil {
		return false
	}
	for _, part := range content.Parts {
		if part != nil && part.FunctionResponse != nil {
			return true
		}
	}
	return false
}

func missingToolResponses(calls []pendingToolCall, content *genai.Content) []pendingToolCall {
	seen := make(map[string]struct{}, len(content.Parts))
	for _, part := range content.Parts {
		if part == nil || part.FunctionResponse == nil {
			continue
		}
		seen[part.FunctionResponse.ID] = struct{}{}
	}

	missing := make([]pendingToolCall, 0, len(calls))
	for _, call := range calls {
		if _, ok := seen[call.id]; ok {
			continue
		}
		missing = append(missing, call)
	}
	return missing
}

func syntheticToolResponseContent(calls []pendingToolCall) *genai.Content {
	return &genai.Content{
		Role:  genai.RoleUser,
		Parts: syntheticToolResponseParts(calls),
	}
}

func syntheticToolResponseParts(calls []pendingToolCall) []*genai.Part {
	parts := make([]*genai.Part, 0, len(calls))
	for _, call := range calls {
		parts = append(parts, &genai.Part{
			FunctionResponse: &genai.FunctionResponse{
				ID:   call.id,
				Name: call.name,
				Response: map[string]any{
					"status":      interruptedToolResultStatus,
					"interrupted": true,
					"error":       "The tool call was interrupted before a result was recorded. Treat the operation as incomplete and inspect current state before retrying.",
				},
			},
		})
	}
	return parts
}
