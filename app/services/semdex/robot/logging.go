package robot

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	agentpkg "google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	adktool "google.golang.org/adk/tool"
	"google.golang.org/genai"
)

const maxLogPreview = 400

func truncateForLog(text string) string {
	if len(text) <= maxLogPreview {
		return text
	}
	return text[:maxLogPreview] + "…"
}

func summariseContentParts(parts []*genai.Part) string {
	var b strings.Builder
	for _, part := range parts {
		if part == nil {
			continue
		}
		piece := strings.TrimSpace(part.Text)
		if piece == "" && part.FunctionCall != nil {
			piece = fmt.Sprintf("[tool-call %s: %v]", part.FunctionCall.Name, part.FunctionCall.Args)
		}
		if piece == "" && part.FunctionResponse != nil {
			piece = fmt.Sprintf("[tool-response %s: %v]", part.FunctionResponse.Name, part.FunctionResponse.Response)
		}
		if piece == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		b.WriteString(piece)
	}
	return b.String()
}

func summariseLLMRequest(req *model.LLMRequest) string {
	if req == nil {
		return ""
	}
	pieces := make([]string, 0, len(req.Contents))
	for _, content := range req.Contents {
		if content == nil {
			continue
		}
		text := summariseContentParts(content.Parts)
		if text == "" {
			continue
		}
		pieces = append(pieces, fmt.Sprintf("%s: %s", content.Role, truncateForLog(text)))
	}
	return strings.Join(pieces, " | ")
}

func summariseLLMResponse(resp *model.LLMResponse) string {
	if resp == nil || resp.Content == nil {
		return ""
	}
	return truncateForLog(summariseContentParts(resp.Content.Parts))
}

func marshalDebug(v any) string {
	if v == nil {
		return ""
	}
	b, err := json.Marshal(v)
	if err != nil {
		return "<unserializable>"
	}
	return truncateForLog(string(b))
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func logBeforeModel(logger *slog.Logger) llmagent.BeforeModelCallback {
	return func(ctx agentpkg.CallbackContext, req *model.LLMRequest) (*model.LLMResponse, error) {
		logger.Info("agent model request",
			slog.String("agent", ctx.AgentName()),
			slog.String("invocation", ctx.InvocationID()),
			slog.String("session", ctx.SessionID()),
			slog.String("prompt", summariseLLMRequest(req)),
		)
		return nil, nil
	}
}

func logAfterModel(logger *slog.Logger) llmagent.AfterModelCallback {
	return func(ctx agentpkg.CallbackContext, resp *model.LLMResponse, respErr error) (*model.LLMResponse, error) {
		if respErr != nil {
			logger.Error("agent model error",
				slog.String("agent", ctx.AgentName()),
				slog.String("invocation", ctx.InvocationID()),
				slog.String("session", ctx.SessionID()),
				slog.String("error", respErr.Error()),
			)
			return nil, respErr
		}

		logger.Info("agent model response",
			slog.String("agent", ctx.AgentName()),
			slog.String("invocation", ctx.InvocationID()),
			slog.String("session", ctx.SessionID()),
			slog.String("finish_reason", fmt.Sprint(resp.FinishReason)),
			slog.String("text", summariseLLMResponse(resp)),
		)
		return nil, nil
	}
}

func logBeforeTool(logger *slog.Logger) llmagent.BeforeToolCallback {
	return func(ctx adktool.Context, tl adktool.Tool, args map[string]any) (map[string]any, error) {
		logger.Info("agent tool start",
			slog.String("tool", tl.Name()),
			slog.String("call_id", ctx.FunctionCallID()),
			slog.String("agent", ctx.AgentName()),
			slog.String("session", ctx.SessionID()),
			slog.Any("args", args),
		)
		return nil, nil
	}
}

func logAfterTool(logger *slog.Logger) llmagent.AfterToolCallback {
	return func(ctx adktool.Context, tl adktool.Tool, args map[string]any, result map[string]any, err error) (map[string]any, error) {
		level := slog.LevelInfo
		msg := "agent tool complete"
		if err != nil {
			level = slog.LevelError
			msg = "agent tool error"
		}

		logger.Log(ctx, level, msg,
			slog.String("tool", tl.Name()),
			slog.String("call_id", ctx.FunctionCallID()),
			slog.String("agent", ctx.AgentName()),
			slog.String("session", ctx.SessionID()),
			slog.String("args", marshalDebug(args)),
			slog.String("result", marshalDebug(result)),
			slog.String("error", errString(err)),
		)
		return nil, err
	}
}

func logBeforeAgent(logger *slog.Logger) agentpkg.BeforeAgentCallback {
	return func(ctx agentpkg.CallbackContext) (*genai.Content, error) {
		logger.Info("agent iteration start",
			slog.String("agent", ctx.AgentName()),
			slog.String("invocation", ctx.InvocationID()),
			slog.String("session", ctx.SessionID()),
		)
		return nil, nil
	}
}

func logAfterAgent(logger *slog.Logger) agentpkg.AfterAgentCallback {
	return func(ctx agentpkg.CallbackContext) (*genai.Content, error) {
		if ctx.Err() != nil {
			logger.Error("agent iteration error",
				slog.String("agent", ctx.AgentName()),
				slog.String("invocation", ctx.InvocationID()),
				slog.String("session", ctx.SessionID()),
				slog.String("error", ctx.Err().Error()),
			)
			return nil, ctx.Err()
		}

		logger.Info("agent iteration complete",
			slog.String("agent", ctx.AgentName()),
			slog.String("invocation", ctx.InvocationID()),
			slog.String("session", ctx.SessionID()),
		)

		return nil, nil
	}
}
