package session_coordinator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"log/slog"
	"strings"
	"time"

	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"go.uber.org/fx"
	adksession "google.golang.org/adk/v2/session"
	"google.golang.org/adk/v2/tool/toolconfirmation"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/account"
	robotresource "github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	robotservice "github.com/Southclaws/storyden/app/services/semdex/robot"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/mcp"
)

const (
	leaseDuration     = time.Minute
	heartbeatInterval = 15 * time.Second
)

type EventKind string

const (
	EventKindQueued    EventKind = "queued"
	EventKindOutput    EventKind = "output"
	EventKindBlocked   EventKind = "blocked"
	EventKindCompleted EventKind = "completed"
	EventKindFailed    EventKind = "failed"
	EventKindBusy      EventKind = "busy"
)

type CommandRobotTurnStart struct {
	TurnID      xid.ID                  `json:"turn_id"`
	SessionID   string                  `json:"session_id"`
	RobotRef    string                  `json:"robot_ref"`
	AccountID   string                  `json:"account_id"`
	Content     *genai.Content          `json:"content,omitempty"`
	ChatContext *mcp.RobotChatContext   `json:"chat_context,omitempty"`
	Options     robotservice.RunOptions `json:"options"`
}

type EventRobotTurn struct {
	TurnID          xid.ID            `json:"turn_id"`
	SessionID       string            `json:"session_id"`
	Sequence        uint64            `json:"sequence"`
	LeaseGeneration uint64            `json:"lease_generation,omitempty"`
	Kind            EventKind         `json:"kind"`
	Event           *adksession.Event `json:"event,omitempty"`
	ErrorText       string            `json:"error_text,omitempty"`
}

type Coordinator struct {
	logger   *slog.Logger
	bus      *pubsub.Bus
	agent    *robotservice.Agent
	sessions *robot_session.Repository
}

func Build() fx.Option {
	return fx.Provide(New)
}

func New(ctx context.Context, lc fx.Lifecycle, logger *slog.Logger, bus *pubsub.Bus, agent *robotservice.Agent, sessions *robot_session.Repository) *Coordinator {
	coordinator := &Coordinator{logger: logger, bus: bus, agent: agent, sessions: sessions}
	lc.Append(fx.StartHook(func(context.Context) error {
		_, err := pubsub.SubscribeCommand(ctx, coordinator.bus, "robot.session_coordinator.start", coordinator.handleStart)
		return err
	}))
	return coordinator
}

func (c *Coordinator) Run(
	ctx context.Context,
	robotRef string,
	userID,
	sessionID string,
	content *genai.Content,
	chatContext *mcp.RobotChatContext,
	options ...robotservice.RunOptions,
) iter.Seq2[*adksession.Event, error] {
	return func(yield func(*adksession.Event, error) bool) {
		turnID := xid.New()
		var runOptions robotservice.RunOptions
		if len(options) > 0 {
			runOptions = options[0]
		}
		command, err := c.prepareTurn(ctx, turnID, robotRef, userID, sessionID, content, chatContext, runOptions)
		if err != nil {
			yield(nil, err)
			return
		}
		topic := sessionTopic(sessionID)
		handlerName := "robot_turn_waiter_" + turnID.String() + "_" + strings.ReplaceAll(uuid.NewString(), "-", "")
		events := make(chan EventRobotTurn, 32)

		subscription, err := pubsub.SubscribeEphemeralNamed(ctx, c.bus, topic, handlerName, func(eventCtx context.Context, payload json.RawMessage) error {
			var event EventRobotTurn
			if err := json.Unmarshal(payload, &event); err != nil {
				return err
			}
			if event.TurnID != turnID {
				return nil
			}
			select {
			case events <- event:
				return nil
			case <-ctx.Done():
				return nil
			case <-eventCtx.Done():
				return nil
			}
		})
		if err != nil {
			yield(nil, err)
			return
		}
		defer subscription.Close()

		if err := c.bus.SendCommand(ctx, command); err != nil {
			if sessionXID, parseErr := xid.FromString(sessionID); parseErr == nil {
				_ = c.sessions.FinishTurn(context.WithoutCancel(ctx), robotresource.SessionID(sessionXID), robotresource.TurnID(turnID), robotresource.TurnStatusFailed, err.Error())
			}
			yield(nil, err)
			return
		}

		nextSequence := uint64(1)
		pending := map[uint64]EventRobotTurn{}
		for {
			select {
			case <-ctx.Done():
				yield(nil, ctx.Err())
				return
			case event := <-events:
				if event.Sequence < nextSequence {
					continue
				}
				pending[event.Sequence] = event
				for {
					ordered, ok := pending[nextSequence]
					if !ok {
						break
					}
					delete(pending, nextSequence)
					nextSequence++

					switch ordered.Kind {
					case EventKindOutput:
						if !yield(ordered.Event, nil) {
							return
						}
					case EventKindCompleted, EventKindBlocked:
						return
					case EventKindBusy:
						yield(nil, robot_session.ErrSessionBusy)
						return
					case EventKindFailed:
						yield(nil, errors.New(ordered.ErrorText))
						return
					default:
						yield(nil, fmt.Errorf("unknown Robot turn event kind %q", ordered.Kind))
						return
					}
				}
			}
		}
	}
}

func (c *Coordinator) Start(
	ctx context.Context,
	robotRef string,
	userID,
	sessionID string,
	content *genai.Content,
	chatContext *mcp.RobotChatContext,
	options ...robotservice.RunOptions,
) (robotresource.TurnID, error) {
	turnID := xid.New()
	var runOptions robotservice.RunOptions
	if len(options) > 0 {
		runOptions = options[0]
	}
	command, err := c.prepareTurn(ctx, turnID, robotRef, userID, sessionID, content, chatContext, runOptions)
	if err != nil {
		return robotresource.TurnID{}, err
	}
	if err := c.bus.SendCommand(ctx, command); err != nil {
		sessionXID, parseErr := xid.FromString(sessionID)
		if parseErr == nil {
			_ = c.sessions.FinishTurn(context.WithoutCancel(ctx), robotresource.SessionID(sessionXID), robotresource.TurnID(turnID), robotresource.TurnStatusFailed, err.Error())
		}
		return robotresource.TurnID{}, err
	}
	return robotresource.TurnID(turnID), nil
}

func (c *Coordinator) publishSessionWake(ctx context.Context, sessionID string, turnID xid.ID) {
	if err := c.bus.PublishNamed(ctx, sessionTopic(sessionID), EventRobotTurn{
		TurnID:    turnID,
		SessionID: sessionID,
		Kind:      EventKindQueued,
	}); err != nil {
		c.logger.Warn("failed to publish Robot session wake",
			slog.String("session_id", sessionID),
			slog.String("turn_id", turnID.String()),
			slog.String("error", err.Error()))
	}
}

func (c *Coordinator) prepareTurn(
	ctx context.Context,
	turnID xid.ID,
	robotRef,
	userID,
	sessionID string,
	content *genai.Content,
	chatContext *mcp.RobotChatContext,
	options robotservice.RunOptions,
) (*CommandRobotTurnStart, error) {
	if err := c.agent.PrepareSession(ctx, robotRef, userID, sessionID); err != nil {
		return nil, err
	}
	sessionXID, err := xid.FromString(sessionID)
	if err != nil {
		return nil, err
	}
	accountXID, err := xid.FromString(userID)
	if err != nil {
		return nil, err
	}
	robotSessionID := robotresource.SessionID(sessionXID)
	accountID := account.AccountID(accountXID)
	if err := c.sessions.EnsureView(ctx, robotSessionID, accountID); err != nil {
		return nil, err
	}
	command := &CommandRobotTurnStart{
		TurnID:      turnID,
		SessionID:   sessionID,
		RobotRef:    robotRef,
		AccountID:   userID,
		Content:     content,
		ChatContext: chatContext,
		Options:     options,
	}
	input, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}
	sourceKind := string(options.Source)
	if isContinuation(content) {
		sourceKind = "tool_result"
	}
	if sourceKind == "" {
		sourceKind = string(robotservice.SourceInteractiveChat)
	}
	continuationOf := opt.NewEmpty[robotresource.TurnID]()
	if isContinuation(content) {
		if blockedTurn, err := c.sessions.BlockedTurn(ctx, robotSessionID); err == nil {
			continuationOf = opt.New(blockedTurn)
		}
	}
	if err := c.sessions.EnqueueTurn(ctx, robot_session.EnqueueTurnParams{
		ID:             robotresource.TurnID(turnID),
		SessionID:      robotSessionID,
		InitiatorID:    opt.New(accountID),
		SourceKind:     sourceKind,
		RobotRef:       robotRef,
		InputData:      input,
		ContinuationOf: continuationOf,
	}); err != nil {
		return nil, err
	}
	c.publishSessionWake(ctx, sessionID, turnID)
	return command, nil
}

func (c *Coordinator) handleStart(ctx context.Context, command *CommandRobotTurnStart) error {
	topic := sessionTopic(command.SessionID)
	publish := func(event EventRobotTurn) bool {
		if err := c.bus.PublishNamed(ctx, topic, event); err != nil {
			c.logger.Error("failed to publish Robot turn event",
				slog.String("session_id", command.SessionID),
				slog.String("turn_id", command.TurnID.String()),
				slog.String("kind", string(event.Kind)),
				slog.String("error", err.Error()))
			return false
		}
		return true
	}
	sessionXID, sessionParseErr := xid.FromString(command.SessionID)
	sessionID := robotresource.SessionID(sessionXID)
	fail := func(err error, sequence uint64, leaseGeneration uint64) {
		if sessionParseErr == nil {
			_ = c.sessions.FinishTurn(context.WithoutCancel(ctx), sessionID, robotresource.TurnID(command.TurnID), robotresource.TurnStatusFailed, err.Error())
		}
		publish(EventRobotTurn{
			TurnID:          command.TurnID,
			SessionID:       command.SessionID,
			Sequence:        sequence,
			LeaseGeneration: leaseGeneration,
			Kind:            EventKindFailed,
			ErrorText:       err.Error(),
		})
	}

	if sessionParseErr != nil {
		fail(sessionParseErr, 1, 0)
		return nil
	}

	lease, err := c.sessions.AcquireQueuedTurnExecution(ctx, sessionID, command.TurnID, isContinuation(command.Content), leaseDuration)
	if err != nil {
		if errors.Is(err, robot_session.ErrTurnHandled) {
			return nil
		}
		kind := EventKindFailed
		if errors.Is(err, robot_session.ErrSessionBusy) {
			kind = EventKindBusy
		}
		_ = c.sessions.FinishTurn(context.WithoutCancel(ctx), sessionID, robotresource.TurnID(command.TurnID), robotresource.TurnStatusFailed, err.Error())
		publish(EventRobotTurn{
			TurnID:    command.TurnID,
			SessionID: command.SessionID,
			Sequence:  1,
			Kind:      kind,
			ErrorText: err.Error(),
		})
		return nil
	}

	runCtx, cancel := context.WithCancel(robot_session.WithExecutionLease(ctx, *lease))
	defer cancel()
	heartbeatDone := make(chan struct{})
	go c.heartbeat(runCtx, cancel, *lease, heartbeatDone)
	defer close(heartbeatDone)

	sequence := uint64(0)
	blocked := false
	for event, runErr := range c.agent.Run(
		runCtx,
		command.RobotRef,
		command.AccountID,
		command.SessionID,
		command.Content,
		command.ChatContext,
		command.Options,
	) {
		if runErr != nil {
			_ = c.sessions.ReleaseExecution(context.WithoutCancel(ctx), *lease)
			fail(runErr, sequence+1, lease.Generation)
			return nil
		}
		if event == nil {
			continue
		}

		pendingIDs := pendingClientToolIDs(event)
		if len(pendingIDs) > 0 {
			if err := c.sessions.StorePendingClientTools(runCtx, sessionID, pendingIDs); err != nil {
				_ = c.sessions.ReleaseExecution(context.WithoutCancel(ctx), *lease)
				fail(err, sequence+1, lease.Generation)
				return nil
			}
			if err := c.sessions.BlockExecution(context.WithoutCancel(ctx), *lease); err != nil {
				fail(err, sequence+1, lease.Generation)
				return nil
			}
			blocked = true
		}

		sequence++
		if !publish(EventRobotTurn{
			TurnID:          command.TurnID,
			SessionID:       command.SessionID,
			Sequence:        sequence,
			LeaseGeneration: lease.Generation,
			Kind:            EventKindOutput,
			Event:           event,
		}) {
			if !blocked {
				_ = c.sessions.ReleaseExecution(context.WithoutCancel(ctx), *lease)
			}
			_ = c.sessions.FinishTurn(context.WithoutCancel(ctx), sessionID, robotresource.TurnID(command.TurnID), robotresource.TurnStatusFailed, "failed to publish Robot turn output")
			return nil
		}
		if blocked {
			break
		}
	}

	terminal := EventKindCompleted
	turnStatus := robotresource.TurnStatusCompleted
	if blocked {
		terminal = EventKindBlocked
		turnStatus = robotresource.TurnStatusBlocked
	} else if err := c.sessions.ReleaseExecution(context.WithoutCancel(ctx), *lease); err != nil {
		fail(err, sequence+1, lease.Generation)
		return nil
	}
	if err := c.sessions.FinishTurn(context.WithoutCancel(ctx), sessionID, robotresource.TurnID(command.TurnID), turnStatus, ""); err != nil {
		fail(err, sequence+1, lease.Generation)
		return nil
	}

	publish(EventRobotTurn{
		TurnID:          command.TurnID,
		SessionID:       command.SessionID,
		Sequence:        sequence + 1,
		LeaseGeneration: lease.Generation,
		Kind:            terminal,
	})
	return nil
}

func (c *Coordinator) heartbeat(ctx context.Context, cancel context.CancelFunc, lease robot_session.ExecutionLease, done <-chan struct{}) {
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-done:
			return
		case <-ticker.C:
			ok, err := c.sessions.HeartbeatExecution(ctx, lease, leaseDuration)
			if err != nil || !ok {
				c.logger.Error("Robot session execution lease lost",
					slog.String("session_id", lease.SessionID.String()),
					slog.String("turn_id", lease.TurnID.String()),
					slog.Any("error", err))
				cancel()
				return
			}
		}
	}
}

func sessionTopic(sessionID string) string {
	return "robot.session." + sessionID + ".turns"
}

func (c *Coordinator) SubscribeTurn(ctx context.Context, sessionID string, turnID robotresource.TurnID, handlerName string, wake chan<- struct{}) (*pubsub.Subscription, error) {
	return pubsub.SubscribeEphemeralNamed(ctx, c.bus, sessionTopic(sessionID), handlerName, func(_ context.Context, payload json.RawMessage) error {
		var event EventRobotTurn
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}
		if event.TurnID != xid.ID(turnID) {
			return nil
		}
		select {
		case wake <- struct{}{}:
		default:
		}
		return nil
	})
}

func (c *Coordinator) SubscribeSession(ctx context.Context, sessionID string, handlerName string, wake chan<- struct{}) (*pubsub.Subscription, error) {
	return pubsub.SubscribeEphemeralNamed(ctx, c.bus, sessionTopic(sessionID), handlerName, func(_ context.Context, payload json.RawMessage) error {
		var event EventRobotTurn
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}
		select {
		case wake <- struct{}{}:
		default:
		}
		return nil
	})
}

func isContinuation(content *genai.Content) bool {
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

func pendingClientToolIDs(event *adksession.Event) []string {
	if event == nil || event.LLMResponse.Content == nil {
		return nil
	}
	seen := map[string]struct{}{}
	var ids []string
	add := func(id string) {
		if id == "" {
			return
		}
		if _, ok := seen[id]; ok {
			return
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	for _, part := range event.LLMResponse.Content.Parts {
		if part == nil {
			continue
		}
		if call := part.FunctionCall; call != nil && call.Name == toolconfirmation.FunctionCallName {
			add(call.ID)
		}
		if response := part.FunctionResponse; response != nil {
			if pending, ok := response.Response["_client_side_pending"].(bool); ok && pending {
				add(response.ID)
			}
		}
	}
	return ids
}
