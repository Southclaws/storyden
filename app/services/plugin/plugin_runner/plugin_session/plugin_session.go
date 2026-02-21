package plugin_session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"
	"github.com/puzpuzpuz/xsync/v4"
	"github.com/rs/xid"
	"golang.org/x/sync/errgroup"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/duplex"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/rpc_handler"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/supervised_runtime"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type sdxSession struct {
	id     plugin.InstallationID
	logger *slog.Logger

	// session run metadata
	reportedState      plugin.ReportedState
	desiredActive      bool
	stateMu            sync.RWMutex
	started            time.Time
	errorMessage       string
	errorDetails       map[string]any
	lastSuccessfulConn time.Time

	// plugin data
	bin      []byte
	manifest *plugin.Validated

	// supervised runtime management (nil for external plugins)
	supervisedRuntime supervised_runtime.Runtime

	// websocket communication
	connMu          sync.Mutex
	connCtx         context.Context
	connCancel      context.CancelFunc
	duplex          duplex.Duplex
	sendChan        chan []byte
	disconnectOnce  sync.Once
	disconnectMu    sync.Mutex
	pendingCommands *xsync.Map[xid.ID, pendingCommand]

	// event subscriptions
	bus               *pubsub.Bus
	subscriptions     []*pubsub.Subscription
	subscriptionsMu   sync.Mutex
	inboundRpcHandler *rpc_handler.Handler
}

func New(
	id plugin.InstallationID,
	bin []byte,
	manifest *plugin.Validated,
	bus *pubsub.Bus,
	parentLogger *slog.Logger,
	inboundRpcHandler *rpc_handler.Handler,
	supervisedRuntime supervised_runtime.Runtime,
) *sdxSession {
	logger := parentLogger.With(slog.String("plugin_id", id.String()))

	sess := &sdxSession{
		id:     id,
		logger: logger,

		reportedState: plugin.ReportedStateInactive,

		bin:      bin,
		manifest: manifest,

		bus:               bus,
		inboundRpcHandler: inboundRpcHandler,

		pendingCommands:   xsync.NewMap[xid.ID, pendingCommand](),
		supervisedRuntime: supervisedRuntime,
	}

	if supervisedRuntime != nil {
		go sess.watchRuntimeEvents()
	}

	return sess
}

type pendingCommand struct {
	sent   time.Time
	respch chan rpc.HostToPluginResponse
}

func (s *sdxSession) ID() plugin.InstallationID {
	return s.id
}

func (s *sdxSession) Supervised() plugin_runner.Supervised {
	if s.supervisedRuntime == nil {
		return nil
	}
	return s.supervisedRuntime
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

func (s *sdxSession) GetErrorDetails() map[string]any {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()

	if len(s.errorDetails) == 0 {
		return nil
	}

	out := make(map[string]any, len(s.errorDetails))
	for k, v := range s.errorDetails {
		out[k] = v
	}
	return out
}

func (s *sdxSession) SetActiveState(ctx context.Context, state plugin.ActiveState) error {
	if s.supervisedRuntime == nil {
		// External plugin lifecycle is websocket-only.
		// NOTE: Not a fan of this right now. Needs another look.
		// 1. Lots of mutex bs going around
		// 2. Weird state, connected, desiredState, reportedState interactions
		switch state {
		case plugin.ActiveStateActive:
			s.stateMu.Lock()
			s.desiredActive = true
			s.stateMu.Unlock()

			s.connMu.Lock()
			isConnected := s.sendChan != nil && s.connCtx != nil
			s.connMu.Unlock()

			if isConnected {
				s.reportState(plugin.ReportedStateActive, "", nil)
			} else {
				s.reportState(plugin.ReportedStateConnecting, "", nil)
			}

		case plugin.ActiveStateInactive:
			s.stateMu.Lock()
			s.desiredActive = false
			s.stateMu.Unlock()

			s.connMu.Lock()
			duplex := s.duplex
			s.duplex = nil
			cancel := s.connCancel
			s.connMu.Unlock()

			if duplex != nil {
				_ = duplex.Close(nil)
			}

			if cancel != nil {
				cancel()
			}

			s.reportState(plugin.ReportedStateInactive, "", nil)

		default:
			return fault.Newf("unknown active state: %s", state)
		}

		return nil
	}

	// Supervised plugin - delegate to runtime provider implementation.
	switch state {
	case plugin.ActiveStateActive:
		return s.supervisedRuntime.Start(ctx)
	case plugin.ActiveStateInactive:
		return s.supervisedRuntime.Stop(ctx)
	default:
		return fault.Newf("unknown active state: %s", state)
	}
}

func (s *sdxSession) Send(ctx context.Context, id xid.ID, payload rpc.HostToPluginRequestUnion) (rpc.HostToPluginResponseUnion, error) {
	s.stateMu.RLock()
	state := s.reportedState
	s.stateMu.RUnlock()

	if state != plugin.ReportedStateActive {
		return rpc.HostToPluginResponseUnion{}, fault.Newf("plugin is not running (state: %s)", state)
	}

	s.connMu.Lock()
	ch := s.sendChan
	connCtx := s.connCtx
	s.connMu.Unlock()

	if ch == nil || connCtx == nil {
		return rpc.HostToPluginResponseUnion{}, fault.New("plugin not connected via websocket")
	}

	pending := pendingCommand{
		sent:   time.Now(),
		respch: make(chan rpc.HostToPluginResponse, 1),
	}

	s.pendingCommands.Store(id, pending)

	b, err := json.Marshal(payload)
	if err != nil {
		s.pendingCommands.Delete(id)
		return rpc.HostToPluginResponseUnion{}, fault.Wrap(err, fmsg.With("failed to encode command"))
	}

	select {
	case <-ctx.Done():
		s.pendingCommands.Delete(id)
		return rpc.HostToPluginResponseUnion{}, ctx.Err()
	case <-connCtx.Done():
		s.pendingCommands.Delete(id)
		return rpc.HostToPluginResponseUnion{}, fault.New("connection closed")
	case ch <- b:
		s.logger.Debug("send rpc",
			slog.String("id", id.String()),
			slog.Any("payload", payload),
		)
	}

	select {
	case <-ctx.Done():
		s.pendingCommands.Delete(id)
		return rpc.HostToPluginResponseUnion{}, ctx.Err()

	case <-connCtx.Done():
		s.pendingCommands.Delete(id)
		return rpc.HostToPluginResponseUnion{}, fault.New("connection closed while waiting for response")

	case resp, ok := <-pending.respch:
		if !ok {
			return rpc.HostToPluginResponseUnion{}, fault.New("connection closed")
		}
		if err, hasErr := resp.Error.Get(); hasErr {
			if msg, hasMsg := err.Message.Get(); hasMsg {
				return rpc.HostToPluginResponseUnion{}, fault.Newf("RPC error: %s", msg)
			}
			return rpc.HostToPluginResponseUnion{}, fault.New("RPC error with no message")
		}
		return resp.Result, nil
	}
}

func (s *sdxSession) reportState(state plugin.ReportedState, message string, details map[string]any) {
	s.stateMu.Lock()
	oldState := s.reportedState
	s.reportedState = state
	s.errorMessage = message
	if len(details) == 0 {
		s.errorDetails = nil
	} else {
		s.errorDetails = make(map[string]any, len(details))
		for k, v := range details {
			s.errorDetails[k] = v
		}
	}
	if state == plugin.ReportedStateActive {
		s.started = time.Now()
	}
	s.stateMu.Unlock()

	// Close subscriptions when transitioning OUT of Active state
	if oldState == plugin.ReportedStateActive && state != plugin.ReportedStateActive {
		s.logger.Info("closing subscriptions due to state transition",
			slog.String("old_state", oldState.String()),
			slog.String("new_state", state.String()))
		s.closeSubscriptions()
	}
}

func (s *sdxSession) handleRPCResponse(ctx context.Context, data []byte) error {
	var resp rpc.HostToPluginResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}

	pending, ok := s.pendingCommands.LoadAndDelete(resp.ID)
	if !ok {
		s.logger.WarnContext(ctx, "received response for unknown command", slog.String("id", resp.ID.String()))
		return nil
	}

	s.logger.DebugContext(ctx, "recv rpc response",
		slog.String("id", resp.ID.String()),
	)

	select {
	case pending.respch <- resp:
	case <-ctx.Done():
	}

	return nil
}

func (s *sdxSession) Connect(ctx context.Context, conn duplex.Duplex) error {
	s.logger.Info("plugin connected")

	s.connMu.Lock()
	if s.sendChan != nil {
		s.connMu.Unlock()
		return fault.New("plugin already connected")
	}
	s.disconnectOnce = sync.Once{}

	connCtx, connCancel := context.WithCancel(ctx)
	sendChan := make(chan []byte, 16)
	s.connCtx = connCtx
	s.connCancel = connCancel
	s.duplex = conn
	s.sendChan = sendChan
	s.connMu.Unlock()

	for _, eventName := range s.manifest.Metadata.EventsConsumed {
		topicName := "rpc." + string(eventName)
		if err := s.Subscribe(ctx, s.bus, string(eventName), topicName); err != nil {
			connCancel()
			return fault.Wrap(err, fmsg.With("failed to subscribe to event"))
		}
	}

	s.lastSuccessfulConn = time.Now()

	s.reportState(plugin.ReportedStateActive, "", nil)

	g, gctx := errgroup.WithContext(connCtx)

	g.Go(func() error {
		return s.sendLoop(gctx, conn, sendChan)
	})

	g.Go(func() error {
		return s.recvLoop(gctx, conn)
	})

	err := g.Wait()

	if isExpectedDisconnect(err) {
		s.disconnect(nil)
	} else {
		s.disconnect(duplex.NewError(duplex.ErrFailed, "internal error"))
	}

	// check for regular disconnect/eof/cancelled "errors" (not real errors.)
	if !isExpectedDisconnect(err) {
		s.logger.Info("plugin disconnected with error",
			slog.Any("error", err))
		return err
	}

	s.logger.Info("plugin disconnected")

	return nil
}

func (s *sdxSession) disconnect(closeCause error) {
	s.disconnectMu.Lock()
	defer s.disconnectMu.Unlock()

	s.disconnectOnce.Do(func() {
		s.connMu.Lock()
		duplex := s.duplex
		if s.connCancel != nil {
			s.connCancel()
		}
		s.sendChan = nil
		s.connCtx = nil
		s.connCancel = nil
		s.duplex = nil
		s.connMu.Unlock()

		if duplex != nil {
			if err := duplex.Close(closeCause); err != nil {
				s.logger.Debug("failed to close duplex", slog.Any("error", err))
			}
		}

		s.pendingCommands.Range(func(key xid.ID, value pendingCommand) bool {
			close(value.respch)
			s.pendingCommands.Delete(key)
			return true
		})

		s.inboundRpcHandler.OnDisconnect()

		s.closeSubscriptions()

		if s.supervisedRuntime == nil {
			s.stateMu.RLock()
			desiredActive := s.desiredActive
			s.stateMu.RUnlock()

			if desiredActive {
				s.reportState(plugin.ReportedStateConnecting, "", nil)
			} else {
				s.reportState(plugin.ReportedStateInactive, "", nil)
			}
		}
	})
}

func (s *sdxSession) watchRuntimeEvents() {
	for event := range s.supervisedRuntime.Events() {
		s.reportState(event.State, event.Message, event.Details)
	}
}

func (s *sdxSession) sendLoop(ctx context.Context, conn duplex.Duplex, sendChan chan []byte) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-sendChan:
			if err := conn.Send(ctx, msg); err != nil {
				return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to send to duplex"))
			}
		}
	}
}

func (s *sdxSession) recvLoop(ctx context.Context, conn duplex.Duplex) error {
	for {
		data, err := conn.Recv(ctx)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to receive from duplex"))
		}

		s.handleIncomingMessage(ctx, data)
	}
}

func (s *sdxSession) handleIncomingMessage(ctx context.Context, data []byte) {
	var wrapper rpc.PluginToHostRequest
	if err := json.Unmarshal(data, &wrapper); err == nil && wrapper.PluginToHostRequestUnion != nil {
		s.handleInboundRPC(ctx, wrapper.PluginToHostRequestUnion)
		return
	}

	if err := s.handleRPCResponse(ctx, data); err != nil {
		s.logger.Error(err.Error())
	}
}

func (s *sdxSession) handleInboundRPC(ctx context.Context, req rpc.PluginToHostRequestUnion) {
	s.logger.Debug("received RPC from plugin",
		slog.String("method", req.PluginToHostRequestType()),
	)

	response, err := s.inboundRpcHandler.Handle(ctx, req)
	if err != nil {
		s.logger.Error("RPC handler error", slog.Any("error", err))
		return
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		s.logger.Error("failed to marshal RPC response", slog.Any("error", err))
		return
	}

	s.connMu.Lock()
	ch := s.sendChan
	connCtx := s.connCtx
	s.connMu.Unlock()

	if ch != nil && connCtx != nil {
		select {
		case ch <- responseBytes:
			s.logger.Debug("sent RPC response",
				slog.String("method", req.PluginToHostRequestType()),
			)
		case <-connCtx.Done():
			s.logger.Warn("connection closed while sending RPC response")
		case <-ctx.Done():
			s.logger.Warn("context canceled while sending RPC response")
		}
	}
}

func (s *sdxSession) closeSubscriptions() {
	s.subscriptionsMu.Lock()
	defer s.subscriptionsMu.Unlock()

	for i, sub := range s.subscriptions {
		if err := sub.Close(); err != nil {
			s.logger.Error("failed to close subscription",
				slog.Int("index", i),
				slog.Any("error", err))
		}
	}

	s.subscriptions = nil
}

func (s *sdxSession) Subscribe(ctx context.Context, bus *pubsub.Bus, eventName, topicName string) error {
	suffix := strings.TrimSuffix(eventName, "Event")
	handlerName := fmt.Sprintf("plugin_%s_%s", s.id.String(), suffix)

	s.logger.Info("subscribing plugin to event",
		slog.String("topic", topicName),
		slog.String("handler", handlerName),
	)

	sub, err := pubsub.SubscribeNamed(ctx, bus, topicName, handlerName, func(ctx context.Context, event json.RawMessage) error {
		defer func() {
			if err := recover(); err != nil {
				s.logger.Error("panic in event handler",
					slog.String("topic", topicName),
					slog.Any("error", err),
				)
			}
		}()
		s.logger.Info("plugin received event",
			slog.String("topic", topicName),
		)

		rpcEvent, err := mapEventToRPC(ctx, topicName, event)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		id := xid.New()

		command := rpc.RPCRequestEvent{
			ID:      id,
			Jsonrpc: "2.0",
			Method:  "event",
			Params: rpc.EventPayload{
				EventPayloadUnion: rpcEvent,
			},
		}

		_, err = s.Send(ctx, id, command)
		if err != nil {
			s.logger.WarnContext(ctx, "failed to send event to plugin",
				slog.String("topic", topicName),
				slog.Any("error", err),
			)
			s.logger.Info("handler returning nil after send error")
		} else {
			s.logger.Info("successfully sent event to plugin",
				slog.String("topic", topicName),
			)
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

func isExpectedDisconnect(err error) bool {
	return err == nil ||
		errors.Is(err, context.Canceled) ||
		errors.Is(err, io.EOF) ||
		duplex.IsExpectedDisconnect(err)
}

func mapEventToRPC(ctx context.Context, name string, event json.RawMessage) (rpc.EventPayloadUnion, error) {
	var v map[string]string
	if err := json.Unmarshal(event, &v); err != nil {
		return nil, err
	}

	// NOTE: A bit of a hack here, we inject the topic name (a string like
	// "EventThreadPublished") into the event data so that the discriminated
	// union is satisfied correctly when re-encoding into an RPC event payload.
	// TODO: Remove the package name from the event topic, maybe generate the
	// topic name from the JSONSchema spec instead of reflecting the Go struct.
	v["event"] = strings.TrimPrefix(name, "rpc.")

	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var evt rpc.EventPayload
	if err := json.Unmarshal(b, &evt); err != nil {
		return nil, err
	}

	return evt, nil
}
