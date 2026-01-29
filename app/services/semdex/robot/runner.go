package robot

import (
	"context"
	"fmt"
	"iter"
	"log/slog"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/samber/lo"
	"go.uber.org/fx"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/artifact"
	"google.golang.org/adk/memory"
	"google.golang.org/adk/model"
	"google.golang.org/adk/runner"
	adksession "google.golang.org/adk/session"
	adktool "google.golang.org/adk/tool"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/mcp"
)

func Build() fx.Option {
	return fx.Options(
		tools.Build(),
		fx.Provide(NewSessionStorage, New, NewSessionNamer),
	)
}

type Agent struct {
	logger       *slog.Logger
	db           *ent.Client
	llm          model.LLM
	sessions     adksession.Service
	sessionNamer *SessionNamer
	artifact     artifact.Service
	memory       memory.Service
	toolProvider *tools.Registry
	runConfig    agent.RunConfig

	beforeModelCallbacks []llmagent.BeforeModelCallback
	afterModelCallbacks  []llmagent.AfterModelCallback
	beforeToolCallbacks  []llmagent.BeforeToolCallback
	afterToolCallbacks   []llmagent.AfterToolCallback
}

func New(
	cfg config.Config,
	logger *slog.Logger,
	db *ent.Client,
	sessionService adksession.Service,
	sessionNamer *SessionNamer,
	toolProvider *tools.Registry,
) (*Agent, error) {
	llm := newOpenAIModel(cfg, logger)

	artifactService := artifact.InMemoryService()
	memoryService := memory.InMemoryService()

	return &Agent{
		logger:       logger,
		db:           db,
		llm:          llm,
		sessions:     sessionService,
		sessionNamer: sessionNamer,
		artifact:     artifactService,
		memory:       memoryService,
		toolProvider: toolProvider,
		runConfig: agent.RunConfig{
			// We don't use streaming, it's cliche as fuck.
			StreamingMode: agent.StreamingModeNone,
		},
		beforeModelCallbacks: []llmagent.BeforeModelCallback{logBeforeModel(logger)},
		afterModelCallbacks:  []llmagent.AfterModelCallback{logAfterModel(logger)},
		beforeToolCallbacks:  []llmagent.BeforeToolCallback{interceptClientSideTools(logger), logBeforeTool(logger)},
		afterToolCallbacks:   []llmagent.AfterToolCallback{logAfterTool(logger)},
	}, nil
}

func (s *Agent) Run(
	ctx context.Context,
	robotID opt.Optional[xid.ID],
	userID,
	sessionID string,
	content *genai.Content,
	chatContext *mcp.RobotChatContext,
) iter.Seq2[*adksession.Event, error] {
	// Load Robot from DB if robotID is provided, otherwise use default agent
	name := defaultAgentName
	description := defaultAgentDescription
	instruction := defaultInstruction
	toolNames := tools.DefaultTools

	if id, ok := robotID.Get(); ok {
		robot, err := s.db.Robot.Get(ctx, id)
		if err != nil {
			return func(yield func(*adksession.Event, error) bool) {
				yield(nil, fmt.Errorf("load robot %s: %w", id, err))
			}
		}

		name = robot.Name
		description = robot.Description
		instruction = robot.Playbook
		toolNames = robot.Tools

		s.logger.Debug("loaded custom robot",
			slog.String("robot_id", id.String()),
			slog.String("robot_name", name),
		)
	}

	toolList, err := s.toolProvider.GetTools(ctx, toolNames...)
	if err != nil {
		return func(yield func(*adksession.Event, error) bool) {
			yield(nil, fmt.Errorf("failed to construct tools: %w", err))
		}
	}

	// NOTE: Inside the actual callback, we lose our tool definition info, so we
	// build a quick table of our own tool definitions here, keyed by name, then
	// in the hook, look up our tool using the name to get permissions to check.
	tooltable := lo.KeyBy(toolList, func(t *tools.Tool) string { return t.Name() })
	checkToolRBAC := func(ctx adktool.Context, tool adktool.Tool, args map[string]any) (map[string]any, error) {
		t := tooltable[tool.Name()]

		if err := session.Authorise(ctx, nil, t.Definition.RequiredPermission); err != nil {
			return nil, err
		}

		return nil, nil
	}

	adktools, err := toolList.ToADKTools(ctx)
	if err != nil {
		return func(yield func(*adksession.Event, error) bool) {
			yield(nil, fmt.Errorf("convert to adk tools: %w", err))
		}
	}
	toolset := &tools.Toolset{ToolList: adktools}

	beforeToolCallbacks := append([]llmagent.BeforeToolCallback{checkToolRBAC}, s.beforeToolCallbacks...)

	llmAgent, err := llmagent.New(llmagent.Config{
		Name:                      name,
		Description:               description,
		GlobalInstructionProvider: s.globalInstructionProvider(ctx, chatContext),
		Instruction:               instruction,
		Model:                     s.llm,
		Toolsets:                  []adktool.Toolset{toolset},
		BeforeModelCallbacks:      s.beforeModelCallbacks,
		AfterModelCallbacks:       s.afterModelCallbacks,
		BeforeToolCallbacks:       beforeToolCallbacks,
		AfterToolCallbacks:        s.afterToolCallbacks,
	})
	if err != nil {
		return func(yield func(*adksession.Event, error) bool) {
			yield(nil, fmt.Errorf("create adk agent: %w", err))
		}
	}

	runner, err := runner.New(runner.Config{
		AppName:         defaultAgentName,
		Agent:           llmAgent,
		SessionService:  s.sessions,
		ArtifactService: s.artifact,
		MemoryService:   s.memory,
	})
	if err != nil {
		return func(yield func(*adksession.Event, error) bool) {
			yield(nil, fmt.Errorf("create adk runner: %w", err))
		}
	}

	if err := s.ensureSession(ctx, defaultAgentName, userID, sessionID); err != nil {
		return func(yield func(*adksession.Event, error) bool) {
			yield(nil, err)
		}
	}

	return runner.Run(ctx, userID, sessionID, content, s.runConfig)
}

// interceptClientSideTools is a BeforeToolCallback that enables client-side tool execution.
//
// CONTEXT: Google ADK Go doesn't have a built-in "confirmation" or "client-side tool" API
// (unlike Python/JS ADK). This callback is a workaround to prevent server-side execution
// of tools that should run on the client (e.g., browser-only tools like console_log).
//
// HOW IT WORKS:
// 1. Intercepts tools marked as client-side (see clientSideTools map below)
// 2. Returns a special marker result {"_client_side_pending": true} instead of executing
// 3. SSE handler (sse.go) detects this marker and:
//   - Emits the tool call to the frontend
//   - Ends the stream without sending the marker to the LLM
//
// 4. Frontend executes the tool and POSTs back the real result
// 5. Backend continues the agent with the real result from the client
//
// This prevents:
// - Dummy results like {pending: true} from reaching the LLM
// - The LLM generating responses before the real tool result is available
// - Wasting tokens on back-and-forth with placeholder data
//
// RELATED: Vercel AI SDK client-side tools, Google ADK tool confirmation (Python-only)
func interceptClientSideTools(logger *slog.Logger) llmagent.BeforeToolCallback {
	return func(ctx adktool.Context, tool adktool.Tool, args map[string]any) (map[string]any, error) {
		if tool.IsLongRunning() {
			logger.Info("intercepting client-side tool",
				slog.String("tool_name", tool.Name()),
				slog.String("call_id", ctx.FunctionCallID()))

			// Return a special marker that tells SSE handler to pause and wait for client
			// This result will NOT be sent to the LLM - see sse.go:isClientSidePending
			return map[string]any{
				"_client_side_pending": true,
				"tool_call_id":         ctx.FunctionCallID(),
			}, nil
		}

		// Not a client-side tool, continue to next callback
		return nil, nil
	}
}

func (s *Agent) ensureSession(ctx context.Context, robotName, userID, sessionID string) error {
	get, err := s.sessions.Get(ctx, &adksession.GetRequest{
		AppName:   robotName,
		UserID:    userID,
		SessionID: sessionID,
	})
	if err == nil {
		s.logger.Debug("session exists",
			slog.String("userID", userID),
			slog.String("sessionID", sessionID),
			slog.Any("session", get.Session),
		)

		s.updateSessionName(ctx, get.Session)
		return nil
	}
	if !ent.IsNotFound(err) {
		return err
	}

	create, err := s.sessions.Create(ctx, &adksession.CreateRequest{
		AppName:   robotName,
		UserID:    userID,
		SessionID: sessionID,
	})
	if err != nil {
		return err
	}
	s.logger.Debug("session created",
		slog.String("userID", userID),
		slog.String("sessionID", sessionID),
		slog.Any("session", create.Session),
	)

	return nil
}

func (s *Agent) updateSessionName(ctx context.Context, sess adksession.Session) {
	text := ""
	for e := range sess.Events().All() {
		for _, part := range e.Content.Parts {
			if part != nil && part.Text != "" {
				text += fmt.Sprintf("%s: %s\n", e.Author, part.Text)
			}
		}
	}

	id, _ := xid.FromString(sess.ID())

	s.sessionNamer.MaybeNameSession(ctx, robot.SessionID(id), text)
}
