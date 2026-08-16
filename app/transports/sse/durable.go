package sse

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
	switch event.Kind {
	case robot.SessionEventTurnQueued:
		start := openapi.StreamPart{}
		_ = start.FromStartPart(openapi.StartPart{MessageId: event.TurnID.String()})
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
			break
		}
		if adkEvent.LLMResponse.Content == nil {
			break
		}
		delegatedEvent := p.delegations.AppendEvent(adkEvent)
		for partIndex, part := range adkEvent.LLMResponse.Content.Parts {
			if part == nil {
				continue
			}
			if !delegatedEvent {
				presentationParts := robotprojection.PresentationPartStreamPartsDeterministic(adkEvent, part, fmt.Sprintf("%s-%d-%d", event.TurnID.String(), event.Sequence, partIndex))
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

	case robot.SessionEventTurnFailed, robot.SessionEventTurnCancelled:
		message := event.ErrorText.Or("Robot turn failed")
		p.delegations.Fail(message)
		errorPart := openapi.StreamPart{}
		_ = errorPart.FromErrorPart(openapi.ErrorPart{ErrorText: message})
		_ = p.collector.Send(errorPart)
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

func newDurableReadHandler(
	logger *slog.Logger,
	coordinator *session_coordinator.Coordinator,
	sessions *robot_session.Repository,
	robotQuerier *robot_querier.Querier,
	toolRegistry *tools.Registry,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		accountID, err := session.GetAccountID(ctx)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		sessionXID, err := xid.FromString(r.PathValue("sessionID"))
		if err != nil {
			http.Error(w, "invalid Robot session ID", http.StatusBadRequest)
			return
		}
		turnXID, err := xid.FromString(r.PathValue("turnID"))
		if err != nil {
			http.Error(w, "invalid Robot turn ID", http.StatusBadRequest)
			return
		}
		sessionID := robot.SessionID(sessionXID)
		turnID := robot.TurnID(turnXID)
		allowed, err := sessions.CanReadTurn(ctx, sessionID, turnID, accountID)
		if err != nil {
			logger.Error("Robot turn stream authorisation", slog.String("error", err.Error()))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if !allowed {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		offset, err := parseStreamOffset(r.URL.Query().Get("offset"))
		if err != nil {
			http.Error(w, "invalid durable stream offset", http.StatusBadRequest)
			return
		}
		switch r.URL.Query().Get("live") {
		case "":
			serveDurableCatchUp(ctx, w, r, logger, sessions, robotQuerier, toolRegistry, sessionID, turnID, offset)
		case "sse":
			serveDurableLive(ctx, w, logger, coordinator, sessions, robotQuerier, toolRegistry, sessionID, turnID, offset)
		default:
			http.Error(w, "unsupported durable stream live mode", http.StatusBadRequest)
		}
	})
}

func serveDurableCatchUp(ctx context.Context, w http.ResponseWriter, r *http.Request, logger *slog.Logger, sessions *robot_session.Repository, robotQuerier *robot_querier.Querier, toolRegistry *tools.Registry, sessionID robot.SessionID, turnID robot.TurnID, offset uint64) {
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
	if r.Method == http.MethodHead {
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

func writeDurableControl(w http.ResponseWriter, control durableControlEvent) error {
	data, err := json.Marshal(control)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "event: control\ndata: %s\n\n", data)
	return err
}

func jsonBatch(parts []openapi.StreamPart) []byte {
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
