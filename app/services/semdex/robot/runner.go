package robot

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"log/slog"
	"strings"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"go.uber.org/fx"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/artifact"
	"google.golang.org/adk/v2/memory"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/runner"
	adksession "google.golang.org/adk/v2/session"
	adktool "google.golang.org/adk/v2/tool"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/llm_provider"
	"github.com/Southclaws/storyden/app/resources/robot/model_ref"
	"github.com/Southclaws/storyden/app/resources/robot/robot_ref"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_history"
	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_registry"
	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_registry/denbot"
	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_registry/pluginbuilder"
	"github.com/Southclaws/storyden/app/services/semdex/robot/llmprovider"
	"github.com/Southclaws/storyden/app/services/semdex/robot/mcpclient"
	"github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	"github.com/Southclaws/storyden/app/services/semdex/robot/toolsets"
	"github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
)

var ErrInternalInvocationUnavailable = errors.New("internal Robot invocation is unavailable")

func Build() fx.Option {
	return fx.Options(
		tools.Build(),
		toolsets.Build(),
		llmprovider.Build(),
		workspaceprovider.Build(),
		denbot.Build(),
		pluginbuilder.Build(),
		fx.Provide(agent_registry.New, NewSessionStorage, New, NewSessionNamer, NewWorkspaceManager, mcpclient.NewHTTPClient, mcpclient.New),
	)
}

type Agent struct {
	logger       *slog.Logger
	db           *ent.Client
	modelFactory *llm_provider.Factory
	sessions     adksession.Service
	sessionRepo  *robot_session.Repository
	sessionNamer *SessionNamer
	workspaces   *WorkspaceManager
	artifact     artifact.Service
	memory       memory.Service
	toolProvider *tools.Registry
	toolsets     *toolsets.Registry
	agents       *agent_registry.Registry
	runConfig    agent.RunConfig

	beforeModelCallbacks []llmagent.BeforeModelCallback
	afterModelCallbacks  []llmagent.AfterModelCallback
	beforeToolCallbacks  []llmagent.BeforeToolCallback
	afterToolCallbacks   []llmagent.AfterToolCallback
}

type resolvedAgentSpec struct {
	RobotRef            string
	DatabaseRobotID     opt.Optional[xid.ID]
	ModelRef            opt.Optional[model_ref.ModelRef]
	AppName             string
	AgentName           string
	DisplayName         string
	Description         string
	Instruction         string
	InstructionProvider func(agent.ReadonlyContext) (string, error)
	ToolNames           []string
	ToolsetRefs         []string
	Capabilities        []string
	DefaultWorkspaceID  opt.Optional[xid.ID]
}

type RunMode = agent_registry.RunMode

const (
	ModeInteractive RunMode = agent_registry.ModeInteractive
	ModeUnattended  RunMode = agent_registry.ModeUnattended
)

type RunSource = agent_registry.RunSource

const (
	SourceInteractiveChat  RunSource = agent_registry.SourceInteractiveChat
	SourcePluginRPC        RunSource = agent_registry.SourcePluginRPC
	SourceScheduled        RunSource = agent_registry.SourceScheduled
	SourceDelegation       RunSource = agent_registry.SourceDelegation
	SourceDelegationResult RunSource = agent_registry.SourceDelegationResult
)

type (
	RunOptions    = agent_registry.RunOptions
	DelegationRun = agent_registry.DelegationRun
)

func defaultRunOptions() RunOptions {
	return RunOptions{
		Mode:   ModeInteractive,
		Source: SourceInteractiveChat,
	}
}

func resolveRunOptions(opts []RunOptions) RunOptions {
	if len(opts) == 0 {
		return defaultRunOptions()
	}
	out := opts[0]
	if out.Mode == "" {
		out.Mode = ModeInteractive
	}
	if out.Source == "" {
		out.Source = SourceInteractiveChat
	}
	return out
}

func New(
	cfg config.Config,
	logger *slog.Logger,
	db *ent.Client,
	factory *llm_provider.Factory,
	sessionService adksession.Service,
	sessionRepo *robot_session.Repository,
	sessionNamer *SessionNamer,
	workspaceManager *WorkspaceManager,
	toolProvider *tools.Registry,
	toolsetRegistry *toolsets.Registry,
	agentRegistry *agent_registry.Registry,
) (*Agent, error) {
	artifactService := artifact.InMemoryService()
	memoryService := memory.InMemoryService()

	return &Agent{
		logger:       logger,
		db:           db,
		modelFactory: factory,
		sessions:     sessionService,
		sessionRepo:  sessionRepo,
		sessionNamer: sessionNamer,
		workspaces:   workspaceManager,
		artifact:     artifactService,
		memory:       memoryService,
		toolProvider: toolProvider,
		toolsets:     toolsetRegistry,
		agents:       agentRegistry,
		runConfig: agent.RunConfig{
			// We don't use streaming, it's cliche as fuck.
			StreamingMode: agent.StreamingModeNone,
		},
		beforeModelCallbacks: []llmagent.BeforeModelCallback{agent_history.RepairInterruptedToolCallsBeforeModel(logger), normalizeClientToolResultsBeforeModel(logger, cfg.PublicWebAddress), logBeforeModel(logger)},
		afterModelCallbacks:  []llmagent.AfterModelCallback{logAfterModel(logger)},
		beforeToolCallbacks:  []llmagent.BeforeToolCallback{logBeforeTool(logger)},
		afterToolCallbacks:   []llmagent.AfterToolCallback{logAfterTool(logger)},
	}, nil
}

func (s *Agent) Run(
	ctx context.Context,
	robotRef string,
	userID,
	sessionID string,
	content *genai.Content,
	invocationContext InvocationContext,
	options ...RunOptions,
) iter.Seq2[*adksession.Event, error] {
	runOptions := resolveRunOptions(options)

	robotRef = strings.TrimSpace(robotRef)
	if robotRef == "" {
		return errorSeq(fmt.Errorf("robot ID is required"))
	}

	spec, err := s.resolveAgentSpec(ctx, robotRef)
	if err != nil {
		return errorSeq(err)
	}

	return s.runResolvedAgent(ctx, spec, userID, sessionID, content, invocationContext, runOptions)
}

// PrepareSession creates a shared Robot session when needed and records the
// current account's view before the coordinator attempts to lease it.
func (s *Agent) PrepareSession(ctx context.Context, robotRef, userID, sessionID string) error {
	spec, err := s.resolveAgentSpec(ctx, strings.TrimSpace(robotRef))
	if err != nil {
		return err
	}
	if err := s.ensureSession(ctx, spec.AppName, spec.RobotRef, userID, sessionID); err != nil {
		return err
	}
	sessionXID, err := xid.FromString(sessionID)
	if err != nil {
		return err
	}
	accountXID, err := xid.FromString(userID)
	if err != nil {
		return err
	}
	return s.sessionRepo.EnsureView(ctx, robot.SessionID(sessionXID), account.AccountID(accountXID))
}

func errorSeq(err error) iter.Seq2[*adksession.Event, error] {
	return func(yield func(*adksession.Event, error) bool) {
		yield(nil, err)
	}
}

func (s *Agent) resolveAgentSpec(ctx context.Context, robotRef string) (*resolvedAgentSpec, error) {
	if id, err := xid.FromString(robotRef); err == nil {
		entRobot, err := s.db.Robot.Get(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("load robot %s: %w", id, err)
		}

		modelRef, err := model_ref.ParseID(entRobot.Model)
		if err != nil {
			return nil, fmt.Errorf("parse robot model ref: %w", err)
		}

		spec := &resolvedAgentSpec{
			RobotRef:           robotRef,
			DatabaseRobotID:    opt.New(id),
			ModelRef:           opt.New(modelRef),
			AppName:            entRobot.Name,
			AgentName:          robot_ref.AgentName(robot_ref.ID(id)),
			DisplayName:        entRobot.Name,
			Description:        entRobot.Description,
			Instruction:        entRobot.Playbook,
			ToolNames:          entRobot.Tools,
			ToolsetRefs:        entRobot.Toolsets,
			Capabilities:       append(append([]string(nil), entRobot.Toolsets...), entRobot.Tools...),
			DefaultWorkspaceID: opt.NewEmpty[xid.ID](),
		}
		if entRobot.WorkspaceID != nil {
			spec.DefaultWorkspaceID = opt.New(*entRobot.WorkspaceID)
		}

		s.logger.Debug(
			"loaded custom robot",
			slog.String("robot_id", id.String()),
			slog.String("robot_name", entRobot.Name),
		)

		return spec, nil
	}

	def, ok := s.agents.Get(robotRef)
	if !ok {
		return nil, fmt.Errorf("unknown robot %q", robotRef)
	}

	capabilities := append([]string(nil), def.Capabilities...)
	if len(capabilities) == 0 {
		capabilities = append([]string(nil), def.ToolsetRefs...)
	}

	spec := &resolvedAgentSpec{
		RobotRef:            robotRef,
		DatabaseRobotID:     opt.NewEmpty[xid.ID](),
		ModelRef:            opt.NewEmpty[model_ref.ModelRef](),
		AppName:             def.AppName,
		AgentName:           def.AgentName,
		DisplayName:         def.Name,
		Description:         def.Description,
		Instruction:         def.Instruction,
		InstructionProvider: def.InstructionProvider,
		ToolsetRefs:         def.ToolsetRefs,
		Capabilities:        capabilities,
		DefaultWorkspaceID:  opt.NewEmpty[xid.ID](),
	}

	return spec, nil
}

func (s *Agent) runResolvedAgent(
	ctx context.Context,
	spec *resolvedAgentSpec,
	userID string,
	sessionID string,
	content *genai.Content,
	invocationContext InvocationContext,
	runOptions RunOptions,
) iter.Seq2[*adksession.Event, error] {
	var llm model.LLM
	if modelRef, ok := spec.ModelRef.Get(); ok {
		var err error
		llm, err = s.modelFactory.GetADKModelLLM(ctx, modelRef)
		if err != nil {
			return errorSeq(fmt.Errorf("initialise robots model: %w", err))
		}
	} else {
		var err error
		llm, err = s.DefaultLLM(ctx)
		if err != nil {
			return errorSeq(err)
		}
	}

	checkToolRBAC := func(ctx agent.Context, tool adktool.Tool, args map[string]any) (map[string]any, error) {
		t, ok := s.toolProvider.FindByADKName(tool.Name())
		if !ok {
			return nil, nil
		}

		if p, ok := t.Definition.RequiredPermission.Get(); ok {
			if err := session.Authorise(ctx, nil, p); err != nil {
				return nil, err
			}
		}

		return nil, nil
	}

	toolCtx := tools.ContextWithRunContext(ctx, tools.RunContext{
		RobotID:   spec.DatabaseRobotID,
		AccountID: userID,
		SessionID: sessionID,
	})
	if runOptions.Mode == ModeUnattended {
		toolCtx = tools.ContextWithConfirmationDisabled(toolCtx)
	}

	var missingTools []string
	resolvedToolsets := make([]adktool.Toolset, 0, len(spec.ToolsetRefs)+1)
	if spec.RobotRef == denbot.ID {
		discovery, err := s.toolsets.Discovery(toolCtx, spec.ToolsetRefs...)
		if err != nil {
			return errorSeq(fmt.Errorf("build Toolset discovery: %w", err))
		}
		resolvedToolsets = append(resolvedToolsets, discovery)
	} else {
		var err error
		resolvedToolsets, missingTools, err = s.toolsets.Build(toolCtx, spec.ToolsetRefs)
		if err != nil {
			return errorSeq(fmt.Errorf("build Robot Toolsets: %w", err))
		}
	}
	if len(spec.ToolNames) > 0 {
		effectiveToolNames, err := s.toolsets.RemoveContainedTools(toolCtx, spec.ToolNames, spec.ToolsetRefs)
		if err != nil {
			return errorSeq(fmt.Errorf("compose direct Robot tools: %w", err))
		}
		directTools, missingDirectTools := s.toolProvider.GetToolsWithMissing(toolCtx, effectiveToolNames...)
		missingTools = append(missingTools, missingDirectTools...)
		if len(directTools) > 0 {
			resolvedToolsets = append(resolvedToolsets, &tools.Toolset{Registered: directTools, BuildContext: toolCtx})
		}
	}
	if runOptions.Mode == ModeUnattended {
		finishTool, err := newUnattendedFinishTool()
		if err != nil {
			return errorSeq(fmt.Errorf("construct unattended finish tool: %w", err))
		}
		resolvedToolsets = append(resolvedToolsets, &tools.Toolset{ToolList: []adktool.Tool{finishTool}})
	}

	beforeToolCallbacks := append([]llmagent.BeforeToolCallback{checkToolRBAC, interceptClientSideTools(s.logger, s.toolProvider, runOptions)}, s.beforeToolCallbacks...)
	currentIdentity := robotIdentity{
		ID:               spec.DatabaseRobotID,
		Name:             spec.DisplayName,
		Description:      spec.Description,
		Capabilities:     spec.Capabilities,
		UnavailableTools: missingTools,
	}
	identityContext := s.buildRobotIdentityContext(sessionID, currentIdentity, runOptions.Delegation != nil)

	if runOptions.Delegation == nil {
		backgroundTools, err := newBackgroundToolset()
		if err != nil {
			return errorSeq(fmt.Errorf("build background tools: %w", err))
		}
		resolvedToolsets = append(resolvedToolsets, backgroundTools)
	}

	if spec.RobotRef == denbot.ID {
		delegationTools, err := s.buildDelegationToolset(ctx)
		if err != nil {
			return errorSeq(fmt.Errorf("build Robot delegation tools: %w", err))
		}
		resolvedToolsets = append(resolvedToolsets, delegationTools)
	}

	llmAgent, err := llmagent.New(llmagent.Config{
		Name:                      spec.AgentName,
		Description:               spec.Description,
		GlobalInstructionProvider: s.globalInstructionProvider(invocationContext, runOptions),
		InstructionProvider:       robotInstructionProvider(spec, identityContext),
		Model:                     llm,
		Mode:                      llmagent.ModeChat,
		Toolsets:                  resolvedToolsets,
		BeforeModelCallbacks:      s.beforeModelCallbacks,
		AfterModelCallbacks:       s.afterModelCallbacks,
		BeforeToolCallbacks:       beforeToolCallbacks,
		AfterToolCallbacks:        s.afterToolCallbacks,
	})
	if err != nil {
		return errorSeq(fmt.Errorf("create adk agent: %w", err))
	}

	return s.RunADKAgent(ctx, agent_registry.ADKRunRequest{
		AppName:            spec.AppName,
		Agent:              llmAgent,
		RobotRef:           spec.RobotRef,
		UserID:             userID,
		SessionID:          sessionID,
		Content:            content,
		DefaultWorkspaceID: spec.DefaultWorkspaceID,
		Options:            runOptions,
	})
}

func (s *Agent) DefaultLLM(ctx context.Context) (model.LLM, error) {
	defaultModel, err := s.modelFactory.DefaultModel(ctx)
	if err != nil {
		return nil, fmt.Errorf("resolve default robots model: %w", err)
	}

	llm, err := s.modelFactory.GetADKModelLLM(ctx, defaultModel)
	if err != nil {
		return nil, fmt.Errorf("initialise robots model: %w", err)
	}

	return llm, nil
}

func (s *Agent) RunADKAgent(
	ctx context.Context,
	req agent_registry.ADKRunRequest,
) iter.Seq2[*adksession.Event, error] {
	ctx = withDelegationAttribution(ctx, req.Options.Delegation)
	if req.Options.Source == SourceScheduled || req.Options.Source == SourceDelegationResult {
		ctx = withHiddenRuntimeInput(ctx)
	}
	r, err := runner.New(runner.Config{
		AppName:         req.AppName,
		Agent:           req.Agent,
		SessionService:  s.sessions,
		ArtifactService: s.artifact,
		MemoryService:   s.memory,
	})
	if err != nil {
		return errorSeq(fmt.Errorf("create adk runner: %w", err))
	}

	rootRobotRef := req.RobotRef
	if req.Options.Delegation != nil && req.Options.Delegation.RootRobotRef != "" {
		rootRobotRef = req.Options.Delegation.RootRobotRef
	}
	if err := s.ensureSession(ctx, req.AppName, rootRobotRef, req.UserID, req.SessionID); err != nil {
		return errorSeq(err)
	}
	if err := s.mountWorkspaceForRun(ctx, req.UserID, req.SessionID, req.DefaultWorkspaceID, req.Options); err != nil {
		return errorSeq(err)
	}
	if req.Options.Mode == ModeUnattended {
		if err := s.markSessionUnattended(ctx, req.SessionID, req.Options); err != nil {
			return errorSeq(err)
		}
	}

	return r.Run(ctx, req.UserID, req.SessionID, req.Content, s.runConfig)
}

func (s *Agent) mountWorkspaceForRun(
	ctx context.Context,
	userID string,
	sessionID string,
	defaultWorkspaceID opt.Optional[xid.ID],
	options RunOptions,
) error {
	sessionXID, err := xid.FromString(sessionID)
	if err != nil {
		return err
	}
	accountXID, err := xid.FromString(userID)
	if err != nil {
		return err
	}

	if spec, ok := options.Workspace.Get(); ok {
		reuse, err := s.canReuseWorkspaceMount(ctx, robot.SessionID(sessionXID), spec)
		if err != nil {
			return err
		}
		if reuse {
			return nil
		}
		_, err = s.workspaces.Mount(ctx, robot.SessionID(sessionXID), account.AccountID(accountXID), spec)
		return err
	}

	if _, ok := defaultWorkspaceID.Get(); !ok {
		return nil
	}

	sess, _, err := s.sessionRepo.Get(ctx, robot.SessionID(sessionXID), robot.NewMessageCursorParams(opt.NewEmpty[robot.MessageID](), 1))
	if err != nil {
		return err
	}
	if _, ok := WorkspaceMountFromState(sess.State).Get(); ok {
		return nil
	}

	workspaceID, _ := defaultWorkspaceID.Get()
	_, err = s.workspaces.Mount(ctx, robot.SessionID(sessionXID), account.AccountID(accountXID), WorkspaceMountSpec{
		WorkspaceID: opt.New(robot.WorkspaceID(workspaceID)),
		Metadata:    map[string]any{},
	})
	return err
}

func (s *Agent) canReuseWorkspaceMount(ctx context.Context, sessionID robot.SessionID, spec WorkspaceMountSpec) (bool, error) {
	sess, _, err := s.sessionRepo.Get(ctx, sessionID, robot.NewMessageCursorParams(opt.NewEmpty[robot.MessageID](), 1))
	if err != nil {
		return false, err
	}

	mount, ok := WorkspaceMountFromState(sess.State).Get()
	if !ok {
		return false, nil
	}

	if workspaceID, ok := spec.WorkspaceID.Get(); ok {
		return mount.WorkspaceID == robot.WorkspaceID(workspaceID), nil
	}
	if instanceID, ok := spec.WorkspaceInstanceID.Get(); ok {
		return mount.WorkspaceInstanceID == robot.WorkspaceInstanceID(instanceID), nil
	}

	return false, nil
}

// interceptClientSideTools is a BeforeToolCallback that enables client-side tool execution
// and human confirmation gates.
//
// CONTEXT: Google ADK Go doesn't have a built-in "confirmation" or "client-side tool" API
// (unlike Python/JS ADK). This callback is a workaround to prevent server-side execution
// of tools that should run on the client (e.g., browser-only tools like console_log).
//
// HOW IT WORKS:
// 1. Intercepts registered Storyden tools marked for client handling.
// 2. Returns a special marker result {"_client_side_pending": true} instead of executing
// 3. The session stream projection detects this marker and:
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
func interceptClientSideTools(logger *slog.Logger, registry *tools.Registry, options RunOptions) llmagent.BeforeToolCallback {
	return func(ctx agent.Context, tool adktool.Tool, args map[string]any) (map[string]any, error) {
		registered, isRegistered := registry.FindByADKName(tool.Name())
		if isRegistered && registered.IsClientTool {
			if options.Mode == ModeUnattended {
				logger.Info("blocking client-side tool in unattended run",
					slog.String("tool_name", tool.Name()),
					slog.String("call_id", ctx.FunctionCallID()))

				return map[string]any{
					"status": "blocked",
					"attention": map[string]any{
						"reason":  "missing_input",
						"message": fmt.Sprintf("The tool %q requires live user input and cannot be used during an unattended invocation.", tool.Name()),
					},
				}, nil
			}

			logger.Info("intercepting long-running tool",
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

func (s *Agent) ensureSession(ctx context.Context, robotName, rootRobotRef, userID, sessionID string) error {
	return s.ensureSessionAttempt(ctx, robotName, rootRobotRef, userID, sessionID, true)
}

func (s *Agent) ensureSessionAttempt(ctx context.Context, robotName, rootRobotRef, userID, sessionID string, retryCreateRace bool) error {
	get, err := s.sessions.Get(ctx, &adksession.GetRequest{
		AppName:   robotName,
		UserID:    userID,
		SessionID: sessionID,
	})
	if err == nil {
		storedRootRobotRef, stateErr := get.Session.State().Get(sessionRootRobotRefStateKey)
		if errors.Is(stateErr, adksession.ErrStateKeyNotExist) {
			storedRootRobotRef = denbot.ID
			stateErr = nil
		}
		if stateErr != nil {
			return stateErr
		}
		storedRootRobotRefString, ok := storedRootRobotRef.(string)
		if !ok {
			return fmt.Errorf("session %s has invalid root Robot state", sessionID)
		}
		if strings.TrimSpace(storedRootRobotRefString) != strings.TrimSpace(rootRobotRef) {
			return fmt.Errorf("session %s already has root Robot %q", sessionID, storedRootRobotRefString)
		}

		s.logger.Debug(
			"session exists",
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
		State:     sessionInitialState(rootRobotRef),
	})
	if err != nil {
		if retryCreateRace && ent.IsConstraintError(err) {
			// Another command may have created the shared session between Get and
			// Create. Re-read it and validate the root Robot through the normal path.
			return s.ensureSessionAttempt(ctx, robotName, rootRobotRef, userID, sessionID, false)
		}
		return err
	}
	s.logger.Debug(
		"session created",
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
