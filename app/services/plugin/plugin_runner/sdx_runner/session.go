package plugin_runner

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"
	"github.com/puzpuzpuz/xsync/v4"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type sdxSession struct {
	id     plugin.InstallationID
	logger *slog.Logger
	runner *sdxRunner

	// session run metadata
	reportedState plugin.ReportedState // current runtime state
	stateMu       sync.RWMutex
	started       time.Time
	errorMessage  string

	// plugin data
	bin      []byte
	manifest *plugin.Validated

	// lifecycle management
	ctx    context.Context
	cancel context.CancelFunc

	// websocket communication
	sendChan        chan rpcMessage
	recvChan        chan rpcMessage
	pendingCommands *xsync.Map[xid.ID, pendingCommand]

	// event subscriptions
	bus             *pubsub.Bus
	subscriptions   []*pubsub.Subscription
	subscriptionsMu sync.Mutex

	mu       sync.Mutex
	workdir  string
	cmd      *exec.Cmd
	stopping bool
}

type rpcMessage struct {
	data []byte
	err  error
}

type pendingCommand struct {
	request rpc.RPCRequest
	sent    time.Time
	respch  chan rpc.RPCResponse
}

func (s *sdxSession) ID() plugin.InstallationID {
	return s.id
}

func (s *sdxSession) GetReportedState() plugin.ReportedState {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()
	return s.reportedState
}

func (s *sdxSession) GetStartedAt() opt.Optional[time.Time] {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()
	return opt.New(s.started)
}

func (s *sdxSession) GetErrorMessage() string {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()
	return s.errorMessage
}

func (s *sdxSession) Send(ctx context.Context, method string, params any) (any, error) {
	s.stateMu.RLock()
	state := s.reportedState
	s.stateMu.RUnlock()

	if state != plugin.ReportedStateActive {
		return nil, fault.Newf("plugin is not running (state: %s)", state)
	}

	if s.sendChan == nil {
		return nil, fault.New("plugin not connected via websocket")
	}

	xidVal := xid.New()

	var rpcParams rpc.RPCRequestParams
	switch v := params.(type) {
	case json.RawMessage:
		if err := json.Unmarshal(v, &rpcParams); err != nil {
			return nil, fault.Wrap(err, fmsg.With("failed to unmarshal json.RawMessage to params"))
		}
	case map[string]interface{}:
		rpcParams = rpc.RPCRequestParams(v)
	case rpc.RPCRequestParams:
		rpcParams = v
	default:
		return nil, fault.Newf("invalid params type: %T", params)
	}

	request := rpc.RPCRequest{
		Jsonrpc: rpc.RPCRequestJsonrpcA20,
		Method:  method,
		Params:  rpcParams,
		Id:      int(xidVal.Counter()),
	}

	pending := pendingCommand{
		request: request,
		sent:    time.Now(),
		respch:  make(chan rpc.RPCResponse, 1),
	}

	s.pendingCommands.Store(xidVal, pending)

	b, err := json.Marshal(request)
	if err != nil {
		s.pendingCommands.Delete(xidVal)
		return nil, fault.Wrap(err, fmsg.With("failed to encode command"))
	}

	select {
	case <-ctx.Done():
		s.pendingCommands.Delete(xidVal)
		return nil, ctx.Err()
	case s.sendChan <- rpcMessage{data: b}:
		s.logger.Debug("send rpc",
			"id", xidVal.String(),
			"method", method,
			"params", params,
		)
	}

	select {
	case <-ctx.Done():
		s.pendingCommands.Delete(xidVal)
		return nil, ctx.Err()
	case resp := <-pending.respch:
		if resp.Error != nil {
			return nil, fault.Newf("RPC error: %s", *resp.Error.Message)
		}
		return resp.Result, nil
	}
}

func (s *sdxSession) setState(state plugin.ReportedState, message string) {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	s.reportedState = state
	s.errorMessage = message
	if state == plugin.ReportedStateActive {
		s.started = time.Now()
	}
}

func (s *sdxSession) handleRPCResponse(ctx context.Context, data []byte) {
	var resp rpc.RPCResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		s.logger.ErrorContext(ctx, "failed to decode RPC response", "error", err, "raw", string(data))
		return
	}

	requestID := xid.ID{}
	for i := 0; i < 12; i++ {
		requestID[i] = byte(resp.Id >> (i * 8))
	}

	pending, ok := s.pendingCommands.LoadAndDelete(requestID)
	if !ok {
		s.logger.WarnContext(ctx, "received response for unknown command", "id", resp.Id)
		return
	}

	s.logger.DebugContext(ctx, "recv rpc response",
		"id", resp.Id,
		"method", pending.request.Method,
	)

	select {
	case pending.respch <- resp:
	case <-ctx.Done():
	}
}

func (s *sdxSession) Connect(ctx context.Context, duplex plugin_runner.Duplex) error {
	s.logger.Info("plugin connected", slog.String("plugin_id", s.id.String()))

	s.mu.Lock()
	if s.sendChan != nil {
		s.mu.Unlock()
		return fault.New("plugin already connected")
	}

	sendChan := make(chan rpcMessage, 16)
	recvChan := make(chan rpcMessage, 16)
	s.sendChan = sendChan
	s.recvChan = recvChan
	s.mu.Unlock()

	for _, eventName := range s.manifest.Metadata.EventsConsumed {
		topicName := "message." + eventName
		if err := s.Subscribe(ctx, s.bus, topicName); err != nil {
			return fault.Wrap(err, fmsg.With("failed to subscribe to event"))
		}
	}

	errChan := make(chan error, 3)

	go s.sendLoop(ctx, duplex, sendChan, errChan)
	go s.recvLoop(ctx, duplex, errChan)

	err := <-errChan

	s.mu.Lock()
	if s.sendChan != nil {
		close(s.sendChan)
		s.sendChan = nil
	}
	if s.recvChan != nil {
		close(s.recvChan)
		s.recvChan = nil
	}
	s.mu.Unlock()

	s.logger.Info("plugin disconnected", slog.String("plugin_id", s.id.String()), slog.Any("error", err))

	return err
}

func (s *sdxSession) sendLoop(ctx context.Context, duplex plugin_runner.Duplex, sendChan chan rpcMessage, errChan chan error) {
	defer func() { errChan <- nil }()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-sendChan:
			if !ok {
				return
			}

			if msg.err != nil {
				continue
			}

			if err := duplex.Send(ctx, msg.data); err != nil {
				s.logger.Error("failed to send to duplex", "error", err)
				errChan <- err
				return
			}
		}
	}
}

func (s *sdxSession) recvLoop(ctx context.Context, duplex plugin_runner.Duplex, errChan chan error) {
	defer func() { errChan <- nil }()

	for {
		data, err := duplex.Recv(ctx)
		if err != nil {
			s.logger.Error("failed to receive from duplex", "error", err)
			errChan <- err
			return
		}

		s.handleIncomingMessage(ctx, data)
	}
}

func (s *sdxSession) handleIncomingMessage(ctx context.Context, data []byte) {
	var rpcReq rpc.RPCRequest
	if err := json.Unmarshal(data, &rpcReq); err == nil && rpcReq.Method != "" {
		s.handleInboundRPC(ctx, rpcReq)
		return
	}

	s.handleRPCResponse(ctx, data)
}

func (s *sdxSession) handleInboundRPC(ctx context.Context, req rpc.RPCRequest) {
	s.logger.Debug("received RPC from plugin",
		"plugin_id", s.id,
		"method", req.Method,
		"id", req.Id,
	)
}

func (s *sdxSession) Subscribe(ctx context.Context, bus *pubsub.Bus, topicName string) error {
	handlerName := "plugin_" + s.id.String()

	s.logger.Info("subscribing plugin to event",
		"plugin_id", s.id,
		"topic", topicName,
		"handler", handlerName,
	)

	sub, err := pubsub.SubscribeNamed(ctx, bus, topicName, handlerName, func(ctx context.Context, event json.RawMessage) error {
		defer func() {
			err := recover()
			s.logger.Error("panic in event handler",
				slog.Any("plugin_id", s.id),
				slog.Any("topic", topicName),
				slog.Any("error", err),
			)
		}()
		s.logger.Info("plugin received event",
			"plugin_id", s.id,
			"topic", topicName,
		)

		s.logger.Info("sending event to plugin",
			"plugin_id", s.id,
			"topic", topicName,
		)

		_, err := s.Send(ctx, "event", event)
		if err != nil {
			s.logger.WarnContext(ctx, "failed to send event to plugin",
				"plugin_id", s.id,
				"topic", topicName,
				"error", err,
			)
			s.logger.Info("handler returning nil after send error")
		} else {
			s.logger.Info("successfully sent event to plugin",
				"plugin_id", s.id,
				"topic", topicName,
			)
			s.logger.Info("handler returning nil after successful send")
		}

		return nil
	})
	if err != nil {
		return fault.Wrap(err, fmsg.With("failed to subscribe to topic"))
	}

	s.subscriptionsMu.Lock()
	s.subscriptions = append(s.subscriptions, sub)
	s.subscriptionsMu.Unlock()

	return nil
}

func (s *sdxSession) Start(ctx context.Context) error {
	s.stateMu.RLock()
	currentState := s.reportedState
	s.stateMu.RUnlock()

	if currentState == plugin.ReportedStateActive {
		return fault.New("plugin is already running")
	}

	go func() {
		if err := s.runProcess(); err != nil {
			s.logger.Error("plugin start failed", slog.String("id", s.id.String()), slog.Any("error", err))
		}
	}()

	return nil
}

func (s *sdxSession) Stop(ctx context.Context) error {
	s.mu.Lock()
	if s.cmd == nil {
		s.mu.Unlock()
		s.stop()
		return nil
	}

	s.stopping = true
	cmd := s.cmd
	s.mu.Unlock()

	if cmd != nil {
		err := cmd.Process.Kill()
		if err != nil && !errors.Is(err, os.ErrProcessDone) {
			return fault.Wrap(err, fmsg.With("failed to stop plugin process"))
		}
	}

	s.stop()

	if s.ctx != nil {
		select {
		case <-s.ctx.Done():
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func (s *sdxSession) stop() {
	s.subscriptionsMu.Lock()
	for _, sub := range s.subscriptions {
		sub.Close()
	}
	s.subscriptions = nil
	s.subscriptionsMu.Unlock()

	if s.cancel != nil {
		s.cancel()
	}
}
