package bindings

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"google.golang.org/adk/v2/tool/toolconfirmation"

	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_querier"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/semdex/robot/session_coordinator"
	"github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/app/transports/http/robotprojection"
)

const (
	headerStreamNextOffset = "Stream-Next-Offset"
	headerStreamUpToDate   = "Stream-Up-To-Date"
	headerStreamClosed     = "Stream-Closed"
	headerStreamCursor     = "Stream-Cursor"
	streamContentTypeJSON  = "application/json"
	streamKeepAlivePeriod  = 15 * time.Second
)

var errInvalidStreamOffset = errors.New("invalid durable stream offset")

type durableControlEvent struct {
	StreamNextOffset string `json:"streamNextOffset"`
	StreamCursor     string `json:"streamCursor"`
	UpToDate         bool   `json:"upToDate"`
	StreamClosed     bool   `json:"streamClosed,omitempty"`
}

type partEmitter interface {
	Send(openapi.StreamPart) error
	Done()
}

type partCollector struct {
	parts []openapi.StreamPart
}

func (c *partCollector) Send(part openapi.StreamPart) error {
	c.parts = append(c.parts, part)
	return nil
}

func (c *partCollector) Done() {}

func (c *partCollector) take() []openapi.StreamPart {
	parts := c.parts
	c.parts = nil
	return parts
}

type turnProjector struct {
	ctx                context.Context
	logger             *slog.Logger
	sessions           *robot_session.Repository
	sessionID          robot.SessionID
	collector          *partCollector
	delegations        *delegationStream
	projectionFinished bool
	hasDelegatedEvents bool
}

func newTurnProjector(
	ctx context.Context,
	logger *slog.Logger,
	sessions *robot_session.Repository,
	robotQuerier *robot_querier.Querier,
	toolRegistry *tools.Registry,
	sessionID robot.SessionID,
) *turnProjector {
	collector := &partCollector{}
	return &turnProjector{
		ctx:         ctx,
		logger:      logger,
		sessions:    sessions,
		sessionID:   sessionID,
		collector:   collector,
		delegations: newDelegationStream(ctx, robotQuerier, robotprojection.ToolMetadataFromRegistry(ctx, toolRegistry), collector, logger),
	}
}

func (p *turnProjector) project(event robot.SessionEvent, toolRegistry *tools.Registry) []openapi.StreamPart {
	turnID, hasTurn := event.TurnID.Get()
	if !hasTurn {
		return []openapi.StreamPart{}
	}
	switch event.Kind {
	case robot.SessionEventTurnQueued:
		start := openapi.StreamPart{}
		_ = start.FromStartPart(openapi.StartPart{MessageId: turnID.String()})
		_ = p.collector.Send(start)
		_ = p.collector.Send(robotprojection.DataStreamPart("data-session_id", p.sessionID.String()))

	case robot.SessionEventMessage:
		if p.projectionFinished {
			break
		}
		message, ok := event.Message.Get()
		if !ok {
			break
		}
		adkEvent := &message.Event
		if adkEvent.Author == "user" {
			if adkEvent.LLMResponse.Content != nil {
				for _, part := range adkEvent.LLMResponse.Content.Parts {
					if part != nil && part.FunctionResponse != nil {
						p.delegations.Complete(part.FunctionResponse)
					}
				}
			}
			break
		}
		if adkEvent.LLMResponse.Content == nil {
			break
		}
		delegatedEvent := p.delegations.AppendEvent(adkEvent)
		p.hasDelegatedEvents = p.hasDelegatedEvents || delegatedEvent
		for partIndex, part := range adkEvent.LLMResponse.Content.Parts {
			if part == nil {
				continue
			}
			if !delegatedEvent {
				presentationParts := robotprojection.PresentationPartStreamPartsDeterministic(adkEvent, part, fmt.Sprintf("%s-%d-%d", turnID.String(), event.Sequence, partIndex))
				for _, presentationPart := range presentationParts {
					_ = p.collector.Send(presentationPart)
				}
			}
			if part.FunctionCall != nil {
				if !delegatedEvent && p.delegations.Start(part.FunctionCall) {
					continue
				}
				if part.FunctionCall.Name == toolconfirmation.FunctionCallName {
					sendToolConfirmationCall(p.ctx, adkEvent, part, p.collector, toolRegistry, p.logger)
				} else if toolRequiresConfirmation(p.ctx, toolRegistry, part.FunctionCall.Name) {
					continue
				} else {
					sendToolCall(p.ctx, adkEvent, part, p.collector, toolRegistry, p.logger)
				}
			}
			if part.FunctionResponse != nil {
				if !delegatedEvent && p.delegations.Complete(part.FunctionResponse) {
					continue
				}
				if adkEvent.Author == "user" || isClientSidePending(part.FunctionResponse.Response) || toolRequiresConfirmation(p.ctx, toolRegistry, part.FunctionResponse.Name) {
					continue
				}
				sendToolResult(part, p.collector, p.logger)
			}
		}
		for _, part := range adkEvent.LLMResponse.Content.Parts {
			if part != nil && part.FunctionResponse != nil && isClientSidePending(part.FunctionResponse.Response) {
				p.finish()
				break
			}
		}
		if hasPendingConfirmation(adkEvent) {
			p.finish()
		}

	case robot.SessionEventTurnCompleted:
		if p.hasDelegatedEvents {
			p.delegations.CompleteRunning()
		}
		if !p.projectionFinished {
			if sess, _, err := p.sessions.Get(p.ctx, p.sessionID, robot.NewMessageCursorParams(opt.NewEmpty[robot.MessageID](), 1)); err == nil {
				data := openapi.StreamPart{}
				_ = data.FromDataPart(openapi.DataPart{Data: sess.Name})
				data.Type = "data-session_name"
				_ = p.collector.Send(data)
			}
			p.finish()
		}

	case robot.SessionEventTurnBlocked:
		p.finish()

	case robot.SessionEventTurnFailed:
		message := event.ErrorText.Or("Robot turn failed")
		p.delegations.Fail(message)
		errorPart := openapi.StreamPart{}
		_ = errorPart.FromErrorPart(openapi.ErrorPart{ErrorText: message})
		_ = p.collector.Send(errorPart)

	case robot.SessionEventTurnCancelled:
		p.finish()
	}
	return p.collector.take()
}

func (p *turnProjector) finish() {
	if p.projectionFinished {
		return
	}
	finish := openapi.StreamPart{}
	_ = finish.FromFinishMessagePart(openapi.FinishMessagePart{})
	_ = p.collector.Send(finish)
	p.projectionFinished = true
}

type projectedRead struct {
	parts      []openapi.StreamPart
	nextOffset uint64
	closed     bool
}

type sessionProjector struct {
	ctx          context.Context
	logger       *slog.Logger
	sessions     *robot_session.Repository
	robotQuerier *robot_querier.Querier
	toolRegistry *tools.Registry
	sessionID    robot.SessionID
	turns        map[robot.TurnID]*turnProjector
}

func newSessionProjector(
	ctx context.Context,
	logger *slog.Logger,
	sessions *robot_session.Repository,
	robotQuerier *robot_querier.Querier,
	toolRegistry *tools.Registry,
	sessionID robot.SessionID,
) *sessionProjector {
	return &sessionProjector{
		ctx:          ctx,
		logger:       logger,
		sessions:     sessions,
		robotQuerier: robotQuerier,
		toolRegistry: toolRegistry,
		sessionID:    sessionID,
		turns:        map[robot.TurnID]*turnProjector{},
	}
}

func (p *sessionProjector) project(event robot.SessionEvent) (openapi.RobotSessionStreamEvent, error) {
	turnID, hasTurn := event.TurnID.Get()
	var parts []openapi.StreamPart
	if hasTurn {
		projector := p.turns[turnID]
		if projector == nil {
			projector = newTurnProjector(p.ctx, p.logger, p.sessions, p.robotQuerier, p.toolRegistry, p.sessionID)
			p.turns[turnID] = projector
		}
		parts = projector.project(event, p.toolRegistry)
	}
	if parts == nil {
		parts = []openapi.StreamPart{}
	}
	streamEvent := openapi.RobotSessionStreamEvent{
		Sequence:  event.Sequence,
		EventKind: openapi.RobotSessionStreamEventKind(event.Kind),
		Parts:     parts,
	}
	if hasTurn {
		id := openapi.Identifier(turnID.String())
		streamEvent.TurnId = &id
	}
	if len(event.InputIDs) > 0 {
		ids := make([]openapi.Identifier, len(event.InputIDs))
		for i, inputID := range event.InputIDs {
			ids[i] = openapi.Identifier(inputID.String())
		}
		streamEvent.InputIds = &ids
	}
	if message, ok := event.Message.Get(); ok && event.Kind == robot.SessionEventInputQueued {
		serialised, err := serialiseRobotSessionMessage(message, map[string]bool{}, robotprojection.ToolMetadataFromRegistry(p.ctx, p.toolRegistry))
		if err != nil {
			return openapi.RobotSessionStreamEvent{}, err
		}
		streamEvent.Message = &serialised
	}

	switch event.Kind {
	case robot.SessionEventTurnCompleted, robot.SessionEventTurnBlocked, robot.SessionEventTurnFailed, robot.SessionEventTurnCancelled:
		if hasTurn {
			delete(p.turns, turnID)
		}
	}
	return streamEvent, nil
}

type projectedSessionRead struct {
	events     []openapi.RobotSessionStreamEvent
	nextOffset uint64
}

func initialiseSessionProjection(
	ctx context.Context,
	logger *slog.Logger,
	sessions *robot_session.Repository,
	robotQuerier *robot_querier.Querier,
	toolRegistry *tools.Registry,
	sessionID robot.SessionID,
	offset uint64,
) (*sessionProjector, *projectedSessionRead, error) {
	projector := newSessionProjector(ctx, logger, sessions, robotQuerier, toolRegistry, sessionID)
	result := &projectedSessionRead{nextOffset: offset}
	after := uint64(0)
	for {
		page, next, err := sessions.ReadSessionEvents(ctx, sessionID, after, robot.SessionEventPageSize)
		if err != nil {
			return nil, nil, err
		}
		for _, event := range page {
			projected, err := projector.project(event)
			if err != nil {
				return nil, nil, err
			}
			if event.Sequence > offset {
				result.events = append(result.events, projected)
				result.nextOffset = event.Sequence
			}
		}
		after = next
		if len(page) < robot.SessionEventPageSize {
			break
		}
	}
	if offset > after {
		return nil, nil, errInvalidStreamOffset
	}
	return projector, result, nil
}

func readProjectedSessionAfter(ctx context.Context, sessions *robot_session.Repository, projector *sessionProjector, sessionID robot.SessionID, offset uint64) (*projectedSessionRead, error) {
	result := &projectedSessionRead{nextOffset: offset}
	after := offset
	for {
		page, next, err := sessions.ReadSessionEvents(ctx, sessionID, after, robot.SessionEventPageSize)
		if err != nil {
			return nil, err
		}
		for _, event := range page {
			projected, err := projector.project(event)
			if err != nil {
				return nil, err
			}
			result.events = append(result.events, projected)
			result.nextOffset = event.Sequence
		}
		after = next
		if len(page) < robot.SessionEventPageSize {
			break
		}
	}
	return result, nil
}

func readProjectedTurn(
	ctx context.Context,
	logger *slog.Logger,
	sessions *robot_session.Repository,
	robotQuerier *robot_querier.Querier,
	toolRegistry *tools.Registry,
	sessionID robot.SessionID,
	turnID robot.TurnID,
	offset uint64,
) (*projectedRead, error) {
	var events []robot.SessionEvent
	after := uint64(0)
	closed := false
	for {
		page, next, pageClosed, err := sessions.ReadTurnEvents(ctx, sessionID, turnID, after, robot.SessionEventPageSize)
		if err != nil {
			return nil, err
		}
		events = append(events, page...)
		after = next
		closed = closed || pageClosed
		if len(page) < robot.SessionEventPageSize || pageClosed {
			break
		}
	}
	if offset > after {
		return nil, errInvalidStreamOffset
	}
	projector := newTurnProjector(ctx, logger, sessions, robotQuerier, toolRegistry, sessionID)
	result := &projectedRead{nextOffset: offset, closed: closed}
	for _, event := range events {
		parts := projector.project(event, toolRegistry)
		if event.Sequence <= offset {
			continue
		}
		result.parts = append(result.parts, parts...)
		result.nextOffset = event.Sequence
	}
	return result, nil
}

type robotSessionTurnGetResponse struct {
	robots  *Robots
	ctx     context.Context
	request openapi.RobotSessionTurnGetRequestObject
}

func (response robotSessionTurnGetResponse) VisitRobotSessionTurnGetResponse(w http.ResponseWriter) error {
	offset := ""
	if response.request.Params.Offset != nil {
		offset = string(*response.request.Params.Offset)
	}
	live := ""
	if response.request.Params.Live != nil {
		live = string(*response.request.Params.Live)
	}
	response.robots.serveSessionTurn(
		response.ctx,
		w,
		string(response.request.SessionId),
		string(response.request.TurnId),
		offset,
		live,
		false,
	)
	return nil
}

type robotSessionTurnHeadResponse struct {
	robots  *Robots
	ctx     context.Context
	request openapi.RobotSessionTurnHeadRequestObject
}

func (response robotSessionTurnHeadResponse) VisitRobotSessionTurnHeadResponse(w http.ResponseWriter) error {
	offset := ""
	if response.request.Params.Offset != nil {
		offset = string(*response.request.Params.Offset)
	}
	response.robots.serveSessionTurn(
		response.ctx,
		w,
		string(response.request.SessionId),
		string(response.request.TurnId),
		offset,
		"",
		true,
	)
	return nil
}

func (r *Robots) RobotSessionTurnGet(ctx context.Context, request openapi.RobotSessionTurnGetRequestObject) (openapi.RobotSessionTurnGetResponseObject, error) {
	return robotSessionTurnGetResponse{robots: r, ctx: ctx, request: request}, nil
}

func (r *Robots) RobotSessionTurnHead(ctx context.Context, request openapi.RobotSessionTurnHeadRequestObject) (openapi.RobotSessionTurnHeadResponseObject, error) {
	return robotSessionTurnHeadResponse{robots: r, ctx: ctx, request: request}, nil
}

func (r *Robots) RobotSessionTurnCancel(ctx context.Context, request openapi.RobotSessionTurnCancelRequestObject) (openapi.RobotSessionTurnCancelResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return openapi.RobotSessionTurnCancel401Response{}, nil
	}
	sessionID, err := robot.NewSessionID(string(request.SessionId))
	if err != nil {
		return openapi.RobotSessionTurnCancel400Response{}, nil
	}
	turnXID, err := xid.FromString(string(request.TurnId))
	if err != nil {
		return openapi.RobotSessionTurnCancel400Response{}, nil
	}
	turnID := robot.TurnID(turnXID)
	allowed, err := r.sessionRepo.CanReadTurn(ctx, sessionID, turnID, accountID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return openapi.RobotSessionTurnCancel404Response{}, nil
	}
	if err := r.coordinator.Cancel(ctx, sessionID, turnID); err != nil {
		return nil, err
	}
	return openapi.RobotSessionTurnCancel202Response{}, nil
}

func (r *Robots) serveSessionTurn(ctx context.Context, w http.ResponseWriter, sessionValue, turnValue, offsetValue, live string, head bool) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	sessionXID, err := xid.FromString(sessionValue)
	if err != nil {
		http.Error(w, "invalid Robot session ID", http.StatusBadRequest)
		return
	}
	turnXID, err := xid.FromString(turnValue)
	if err != nil {
		http.Error(w, "invalid Robot turn ID", http.StatusBadRequest)
		return
	}
	sessionID := robot.SessionID(sessionXID)
	turnID := robot.TurnID(turnXID)
	allowed, err := r.sessionRepo.CanReadTurn(ctx, sessionID, turnID, accountID)
	if err != nil {
		r.logger.Error("Robot turn stream authorisation", slog.String("error", err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !allowed {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	offset, err := parseStreamOffset(offsetValue)
	if err != nil {
		http.Error(w, "invalid durable stream offset", http.StatusBadRequest)
		return
	}
	if head {
		serveDurableCatchUp(ctx, w, r.logger, r.sessionRepo, r.robotQuerier, r.tools, sessionID, turnID, offset, true)
		return
	}
	switch live {
	case "":
		serveDurableCatchUp(ctx, w, r.logger, r.sessionRepo, r.robotQuerier, r.tools, sessionID, turnID, offset, false)
	case "sse":
		serveDurableLive(ctx, w, r.logger, r.coordinator, r.sessionRepo, r.robotQuerier, r.tools, sessionID, turnID, offset)
	default:
		http.Error(w, "unsupported durable stream live mode", http.StatusBadRequest)
	}
}

func serveDurableCatchUp(ctx context.Context, w http.ResponseWriter, logger *slog.Logger, sessions *robot_session.Repository, robotQuerier *robot_querier.Querier, toolRegistry *tools.Registry, sessionID robot.SessionID, turnID robot.TurnID, offset uint64, head bool) {
	result, err := readProjectedTurn(ctx, logger, sessions, robotQuerier, toolRegistry, sessionID, turnID, offset)
	if err != nil {
		writeDurableReadError(w, err, logger)
		return
	}
	header := w.Header()
	header.Set("Content-Type", streamContentTypeJSON)
	header.Set("Cache-Control", headerNoCache)
	header.Set(headerStreamNextOffset, formatStreamOffset(result.nextOffset))
	header.Set(headerStreamUpToDate, "true")
	if result.closed {
		header.Set(headerStreamClosed, "true")
	}
	if head {
		w.WriteHeader(http.StatusOK)
		return
	}
	_, _ = w.Write(jsonBatch(result.parts))
}

func serveDurableLive(ctx context.Context, w http.ResponseWriter, logger *slog.Logger, coordinator *session_coordinator.Coordinator, sessions *robot_session.Repository, robotQuerier *robot_querier.Querier, toolRegistry *tools.Registry, sessionID robot.SessionID, turnID robot.TurnID, offset uint64) {
	flusher, ok := GetFlusher(w)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	wake := make(chan struct{}, 1)
	cursor := xid.New().String()
	subscription, err := coordinator.SubscribeTurn(ctx, sessionID.String(), turnID, "robot-turn-reader-"+cursor, wake)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer subscription.Close()
	header := w.Header()
	header.Set("Content-Type", contentTypeEventStream)
	header.Set("Cache-Control", headerNoCache)
	header.Set("Connection", headerKeepAlive)
	header.Set("X-Accel-Buffering", "no")
	header.Set(headerStreamNextOffset, formatStreamOffset(offset))
	header.Set(headerStreamCursor, cursor)
	w.WriteHeader(http.StatusOK)
	flusher.Flush()
	keepAlive := time.NewTicker(streamKeepAlivePeriod)
	defer keepAlive.Stop()
	for {
		result, err := readProjectedTurn(ctx, logger, sessions, robotQuerier, toolRegistry, sessionID, turnID, offset)
		if err != nil {
			logger.Error("Robot turn stream read", slog.String("error", err.Error()))
			return
		}
		if len(result.parts) > 0 {
			if _, err := fmt.Fprintf(w, "event: data\ndata: %s\n\n", jsonBatch(result.parts)); err != nil {
				return
			}
		}
		control := durableControlEvent{StreamNextOffset: formatStreamOffset(result.nextOffset), StreamCursor: cursor, UpToDate: true, StreamClosed: result.closed}
		if err := writeDurableControl(w, control); err != nil {
			return
		}
		flusher.Flush()
		offset = result.nextOffset
		if result.closed {
			return
		}
		select {
		case <-ctx.Done():
			return
		case <-wake:
		case <-keepAlive.C:
			if _, err := w.Write([]byte(": keepalive\n\n")); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func serveSessionCatchUp(ctx context.Context, w http.ResponseWriter, logger *slog.Logger, sessions *robot_session.Repository, robotQuerier *robot_querier.Querier, toolRegistry *tools.Registry, sessionID robot.SessionID, offset uint64, head bool) {
	_, result, err := initialiseSessionProjection(ctx, logger, sessions, robotQuerier, toolRegistry, sessionID, offset)
	if err != nil {
		writeDurableReadError(w, err, logger)
		return
	}
	header := w.Header()
	header.Set("Content-Type", streamContentTypeJSON)
	header.Set("Cache-Control", headerNoCache)
	header.Set(headerStreamNextOffset, formatStreamOffset(result.nextOffset))
	header.Set(headerStreamUpToDate, "true")
	if head {
		w.WriteHeader(http.StatusOK)
		return
	}
	_, _ = w.Write(jsonBatch(result.events))
}

func serveSessionLive(ctx context.Context, w http.ResponseWriter, logger *slog.Logger, coordinator *session_coordinator.Coordinator, sessions *robot_session.Repository, robotQuerier *robot_querier.Querier, toolRegistry *tools.Registry, sessionID robot.SessionID, offset uint64) {
	flusher, ok := GetFlusher(w)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	wake := make(chan struct{}, 1)
	cursor := xid.New().String()
	subscription, err := coordinator.SubscribeSession(ctx, sessionID.String(), "robot-session-reader-"+cursor, wake)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer subscription.Close()

	projector, initial, err := initialiseSessionProjection(ctx, logger, sessions, robotQuerier, toolRegistry, sessionID, offset)
	if err != nil {
		if errors.Is(err, context.Canceled) || ctx.Err() != nil {
			return
		}
		writeDurableReadError(w, err, logger)
		return
	}
	header := w.Header()
	header.Set("Content-Type", contentTypeEventStream)
	header.Set("Cache-Control", headerNoCache)
	header.Set("Connection", headerKeepAlive)
	header.Set("X-Accel-Buffering", "no")
	header.Set(headerStreamNextOffset, formatStreamOffset(offset))
	header.Set(headerStreamCursor, cursor)
	w.WriteHeader(http.StatusOK)
	if len(initial.events) > 0 {
		if _, err := fmt.Fprintf(w, "event: data\ndata: %s\n\n", jsonBatch(initial.events)); err != nil {
			return
		}
	}
	offset = initial.nextOffset
	if err := writeDurableControl(w, durableControlEvent{StreamNextOffset: formatStreamOffset(offset), StreamCursor: cursor, UpToDate: true}); err != nil {
		return
	}
	flusher.Flush()

	keepAlive := time.NewTicker(streamKeepAlivePeriod)
	defer keepAlive.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-wake:
		case <-keepAlive.C:
			if _, err := w.Write([]byte(": keepalive\n\n")); err != nil {
				return
			}
			flusher.Flush()
			continue
		}

		result, err := readProjectedSessionAfter(ctx, sessions, projector, sessionID, offset)
		if err != nil {
			if errors.Is(err, context.Canceled) || ctx.Err() != nil {
				return
			}
			logger.Error("Robot session stream read", slog.String("error", err.Error()))
			return
		}
		if len(result.events) == 0 {
			continue
		}
		if _, err := fmt.Fprintf(w, "event: data\ndata: %s\n\n", jsonBatch(result.events)); err != nil {
			return
		}
		offset = result.nextOffset
		if err := writeDurableControl(w, durableControlEvent{StreamNextOffset: formatStreamOffset(offset), StreamCursor: cursor, UpToDate: true}); err != nil {
			return
		}
		flusher.Flush()
	}
}

func writeDurableControl(w http.ResponseWriter, control durableControlEvent) error {
	data, err := json.Marshal(control)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "event: control\ndata: %s\n\n", data)
	return err
}

func jsonBatch[T any](parts []T) []byte {
	var body bytes.Buffer
	body.WriteByte('[')
	for i, part := range parts {
		if i > 0 {
			body.WriteByte(',')
		}
		encoded, _ := json.Marshal(part)
		body.Write(encoded)
	}
	body.WriteByte(']')
	return body.Bytes()
}

func formatStreamOffset(offset uint64) string {
	return fmt.Sprintf("%016d_%016d", 0, offset)
}

func parseStreamOffset(value string) (uint64, error) {
	if value == "" || value == "-1" {
		return 0, nil
	}
	parts := strings.Split(value, "_")
	if len(parts) != 2 || len(parts[0]) != 16 || len(parts[1]) != 16 {
		return 0, errInvalidStreamOffset
	}
	generation, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil || generation != 0 {
		return 0, errInvalidStreamOffset
	}
	offset, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return 0, errInvalidStreamOffset
	}
	return offset, nil
}

func writeDurableReadError(w http.ResponseWriter, err error, logger *slog.Logger) {
	if errors.Is(err, errInvalidStreamOffset) {
		http.Error(w, "invalid durable stream offset", http.StatusBadRequest)
		return
	}
	logger.Error("Robot turn stream read", slog.String("error", err.Error()))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
