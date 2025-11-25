package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iter"
	"log/slog"
	"strings"

	storydentools "github.com/Southclaws/storyden/app/services/semdex/agent/tools"
	"github.com/Southclaws/storyden/internal/infrastructure/ai"

	agentpkg "google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/workflowagents/loopagent"
	"google.golang.org/adk/artifact"
	"google.golang.org/adk/memory"
	"google.golang.org/adk/model"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	adktool "google.golang.org/adk/tool"
	"google.golang.org/genai"
)

const (
	defaultAgentName        = "storyden"
	defaultAgentDescription = "Storyden agent capable of working with pages, links, tags, and threads."
	defaultToolsetName      = "storyden_mcp_toolset"
)

const defaultInstruction = `You are Storyden's built-in automation agent.

Storyden is a combined forum + wiki for communities. "Library pages" are structured knowledge-base articles and "threads" are real discussions.

Your job is to inspect or maintain the knowledge base using the exact tools you have been given.

Rules:
1. Only call tools that actually exist in the supplied list.
2. Prefer read-only tools (tree/list/get/search) to gather context before making mutations.
3. Reference real Storyden concepts (pages, tags, threads, categories) when you describe results or ask clarifying questions. If a request is impossible with the tools you have, clearly state that limitation.
4. When a change is requested, explain what you are about to do and then call the appropriate tool.
`

// Agent exposes the functionality we need at the transport layer so that
// callers (CLI, HTTP transports, etc.) can stream agent events.
type Agent interface {
	Run(ctx context.Context, userID, sessionID string, content *genai.Content, cfg agentpkg.RunConfig) iter.Seq2[*session.Event, error]
}

// New builds a Storyden agent backed by Google's ADK. The returned Agent uses
// the same MCP tools used by the /mcp transport which keeps internal and
// external tooling in sync.
func New(logger *slog.Logger, prompter ai.Prompter, tools storydentools.All) (Agent, error) {
	if len(tools) == 0 {
		return nil, errors.New("at least one MCP tool is required")
	}
	llm := newPrompterModel(prompter)
	if llm == nil {
		return nil, errors.New("agent requires a language model provider")
	}
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	toolset, err := newAdapterToolset(tools)
	if err != nil {
		return nil, fmt.Errorf("build toolset: %w", err)
	}
	toolCatalog := describeToolCatalog(tools)
	instruction := buildInstruction(toolCatalog)

	llmAgent, err := llmagent.New(llmagent.Config{
		Name:        defaultAgentName,
		Description: defaultAgentDescription,
		Instruction: instruction,
		Model:       llm,
		Toolsets: []adktool.Toolset{
			toolset,
		},
		BeforeModelCallbacks: []llmagent.BeforeModelCallback{logBeforeModel(logger)},
		AfterModelCallbacks:  []llmagent.AfterModelCallback{logAfterModel(logger)},
		BeforeToolCallbacks:  []llmagent.BeforeToolCallback{logBeforeTool(logger)},
		AfterToolCallbacks:   []llmagent.AfterToolCallback{logAfterTool(logger)},
	})
	if err != nil {
		return nil, fmt.Errorf("create llm agent: %w", err)
	}

	loop, err := loopagent.New(loopagent.Config{
		MaxIterations: 10,
		AgentConfig: agentpkg.Config{
			Name:        "storyden_loop",
			Description: "Storyden loop agent orchestrator",
			SubAgents:   []agentpkg.Agent{llmAgent},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("create loop agent: %w", err)
	}

	sessionService := session.InMemoryService()
	artifactService := artifact.InMemoryService()
	memoryService := memory.InMemoryService()

	runner, err := runner.New(runner.Config{
		AppName:         defaultAgentName,
		Agent:           loop,
		SessionService:  sessionService,
		ArtifactService: artifactService,
		MemoryService:   memoryService,
	})
	if err != nil {
		return nil, fmt.Errorf("create agent runner: %w", err)
	}

	return &storydenAgent{
		runner:   runner,
		sessions: sessionService,
		artifact: artifactService,
		memory:   memoryService,
		appName:  defaultAgentName,
	}, nil
}

type storydenAgent struct {
	runner   *runner.Runner
	sessions session.Service
	artifact artifact.Service
	memory   memory.Service
	appName  string
}

func (s *storydenAgent) Run(ctx context.Context, userID, sessionID string, content *genai.Content, cfg agentpkg.RunConfig) iter.Seq2[*session.Event, error] {
	if err := s.ensureSession(ctx, userID, sessionID); err != nil {
		return func(yield func(*session.Event, error) bool) {
			yield(nil, err)
		}
	}

	return s.runner.Run(ctx, userID, sessionID, content, cfg)
}

func (s *storydenAgent) ensureSession(ctx context.Context, userID, sessionID string) error {
	if userID == "" || sessionID == "" {
		return errors.New("user id and session id are required")
	}

	if _, err := s.sessions.Get(ctx, &session.GetRequest{
		AppName:   s.appName,
		UserID:    userID,
		SessionID: sessionID,
	}); err == nil {
		return nil
	}

	if _, err := s.sessions.Create(ctx, &session.CreateRequest{
		AppName:   s.appName,
		UserID:    userID,
		SessionID: sessionID,
	}); err != nil {
		if _, getErr := s.sessions.Get(ctx, &session.GetRequest{
			AppName:   s.appName,
			UserID:    userID,
			SessionID: sessionID,
		}); getErr != nil {
			return fmt.Errorf("create session: %w", err)
		}
	}

	return nil
}

const maxLogPreview = 400

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
			slog.String("args", marshalDebug(args)),
		)
		return args, nil
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
		return result, err
	}
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

func summariseContentParts(parts []*genai.Part) string {
	var b strings.Builder
	for _, part := range parts {
		if part == nil {
			continue
		}
		piece := strings.TrimSpace(part.Text)
		if piece == "" && part.FunctionCall != nil {
			piece = fmt.Sprintf("[tool-call %s]", part.FunctionCall.Name)
		}
		if piece == "" && part.FunctionResponse != nil {
			piece = fmt.Sprintf("[tool-response %s]", part.FunctionResponse.Name)
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

func truncateForLog(text string) string {
	if len(text) <= maxLogPreview {
		return text
	}
	return text[:maxLogPreview] + "…"
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

func buildInstruction(toolCatalog string) string {
	toolCatalog = strings.TrimSpace(toolCatalog)
	if toolCatalog == "" {
		return defaultInstruction
	}
	return defaultInstruction + "\n\nAvailable tools:\n" + toolCatalog
}

func describeToolCatalog(all storydentools.All) string {
	if len(all) == 0 {
		return ""
	}
	var b strings.Builder
	for _, serverTool := range all {
		if serverTool.Tool.Name == "" {
			continue
		}
		desc := strings.TrimSpace(serverTool.Tool.Description)
		if desc == "" {
			desc = "(no description provided)"
		}
		fmt.Fprintf(&b, "- %s: %s\n", serverTool.Tool.Name, desc)
	}
	return b.String()
}
