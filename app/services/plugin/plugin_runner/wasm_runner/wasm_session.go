package plugin_runner

// import (
// 	"bufio"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log/slog"
// 	"os"
// 	"sync"
// 	"time"

// 	"github.com/Southclaws/fault"
// 	"github.com/Southclaws/fault/fmsg"
// 	"github.com/Southclaws/opt"
// 	resource_plugin "github.com/Southclaws/storyden/app/resources/plugin"
// 	lib_plugin "github.com/Southclaws/storyden/lib/plugin"
// 	"github.com/puzpuzpuz/xsync/v4"
// 	"github.com/rs/xid"
// 	"github.com/tetratelabs/wazero"
// 	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
// )

// // TODO: Make this configurable in some way. Perhaps at the instance level or
// // perhaps at the plugin level - if a plugin declares it implements particularly
// // lengthy operations (such as an AI-driven plugin or something similar.)
// var defaultRPCTimeout = time.Second * 30

// type RunType int

// const (
// 	RunTypeBackground RunType = iota
// 	RunTypeOnDemand
// )

// type PluginSession struct {
// 	id     resource_plugin.InstallationID
// 	logger *slog.Logger

// 	// session run metadata
// 	runType       RunType                       // one-shot or background
// 	reportedState resource_plugin.ReportedState // current runtime state
// 	stateMu       sync.RWMutex
// 	started       time.Time
// 	errorMessage  string

// 	// plugin data
// 	bin      []byte
// 	manifest *resource_plugin.Validated

// 	// runtime
// 	runtime wazero.Runtime
// 	runner  *wazeroRunner

// 	// lifecycle management
// 	ctx    context.Context
// 	cancel context.CancelFunc

// 	// in/out comms
// 	inchan          chan []byte
// 	outchan         chan []byte
// 	pendingCommands *xsync.Map[xid.ID, pendingCommand]
// }

// type pendingCommand struct {
// 	request RPCRequest
// 	sent    time.Time
// 	respch  chan RPCResponse
// }

// type RPCRequest struct {
// 	ID     string `json:"id"`
// 	Method string `json:"method"`
// 	Params any    `json:"params"`
// }

// type RPCResponse struct {
// 	ID     string `json:"id"`
// 	Result any    `json:"result"`
// }

// func (s *PluginSession) ID() resource_plugin.InstallationID {
// 	return s.id
// }

// func (s *PluginSession) GetReportedState() resource_plugin.ReportedState {
// 	s.stateMu.RLock()
// 	defer s.stateMu.RUnlock()
// 	return s.reportedState
// }

// func (s *PluginSession) GetStartedAt() opt.Optional[time.Time] {
// 	s.stateMu.RLock()
// 	defer s.stateMu.RUnlock()

// 	if s.reportedState != resource_plugin.ReportedStateActive {
// 		return opt.NewEmpty[time.Time]()
// 	}
// 	return opt.New(s.started)
// }

// func (s *PluginSession) GetErrorMessage() string {
// 	s.stateMu.RLock()
// 	defer s.stateMu.RUnlock()
// 	return s.errorMessage
// }

// func (s *PluginSession) setState(state resource_plugin.ReportedState, errMsg string) {
// 	s.stateMu.Lock()
// 	defer s.stateMu.Unlock()
// 	s.reportedState = state
// 	s.errorMessage = errMsg

// 	if state == resource_plugin.ReportedStateActive && s.started.IsZero() {
// 		s.started = time.Now()
// 	}
// }

// func (s *PluginSession) Send(ctx context.Context, method string, params any) (any, error) {
// 	s.stateMu.RLock()
// 	state := s.reportedState
// 	s.stateMu.RUnlock()

// 	if state != resource_plugin.ReportedStateActive {
// 		return nil, fault.Newf("plugin is not running (state: %s)", state)
// 	}

// 	if s.runType == RunTypeOnDemand {
// 		return nil, fault.New("todo: on-demand run call")
// 	}

// 	id := xid.New()
// 	request := RPCRequest{
// 		ID:     id.String(),
// 		Method: method,
// 		Params: params,
// 	}

// 	pending := pendingCommand{
// 		request: request,
// 		sent:    time.Now(),
// 		respch:  make(chan RPCResponse),
// 	}

// 	s.pendingCommands.Store(id, pending)

// 	b, err := json.Marshal(request)
// 	if err != nil {
// 		panic(fault.Wrap(err, fmsg.With("failed to encode command")))
// 	}

// 	if _, ok := ctx.Deadline(); !ok {
// 		// If the context doesn't have a deadline, set a default timeout.
// 		var cancel context.CancelFunc
// 		ctx, cancel = context.WithTimeout(ctx, defaultRPCTimeout)
// 		defer cancel()
// 	}

// 	select {
// 	case <-ctx.Done():
// 		return nil, ctx.Err()

// 	case s.inchan <- b:
// 		s.logger.Debug("send rpc",
// 			slog.String("id", id.String()),
// 			slog.String("method", method),
// 			slog.Any("params", params),
// 		)
// 	}

// 	select {
// 	case <-ctx.Done():
// 		return nil, ctx.Err()

// 	case resp := <-pending.respch:
// 		return resp.Result, nil
// 	}
// }

// func (s *PluginSession) stop() {
// 	s.setState(resource_plugin.ReportedStateInactive, "")

// 	if s.cancel != nil {
// 		s.cancel()
// 	}

// 	if s.inchan != nil {
// 		close(s.inchan)
// 	}
// 	if s.outchan != nil {
// 		close(s.outchan)
// 	}
// }

// func (s *PluginSession) start() error {
// 	s.stateMu.RLock()
// 	currentState := s.reportedState
// 	s.stateMu.RUnlock()

// 	if currentState == resource_plugin.ReportedStateActive {
// 		return fault.New("plugin is already running")
// 	}

// 	if s.runType != RunTypeBackground {
// 		return fault.New("only background plugins can be started")
// 	}

// 	s.ctx, s.cancel = context.WithCancel(context.Background())

// 	ir, iw := io.Pipe()
// 	or, ow := io.Pipe()

// 	errCh := make(chan error, 1)

// 	mc := wazero.NewModuleConfig().
// 		WithStdin(ir).
// 		WithStdout(ow).
// 		WithStderr(os.Stderr) // TODO: Route error to persisted plugin logs

// 	handleError := func(err error) {
// 		s.logger.Error("plugin error", slog.Any("error", err))
// 		s.setState(resource_plugin.ReportedStateError, err.Error())
// 	}

// 	go func() {
// 		defer s.cancel()

// 		closer, err := wasi_snapshot_preview1.Instantiate(s.ctx, s.runtime)
// 		if err != nil {
// 			errCh <- fault.Wrap(err, fmsg.With("failed to instantiate wasi"))
// 			return
// 		}
// 		defer closer.Close(s.ctx)

// 		mod, err := s.runtime.InstantiateWithConfig(s.ctx, s.bin, mc)
// 		if err != nil {
// 			errCh <- fault.Wrap(err, fmsg.With("failed to instantiate"))
// 			return
// 		}

// 		if err := mod.Close(s.ctx); err != nil {
// 			errCh <- fault.Wrap(err, fmsg.With("failed to close plugin"))
// 			return
// 		}

// 		errCh <- fault.New("plugin exited unexpectedly")
// 	}()

// 	// Scan the first line of output immediately for manifest.
// 	// TODO: Probably not needed any more.
// 	manifestScan := bufio.NewScanner(or)
// 	manifestScan.Split(bufio.ScanLines)
// 	if !manifestScan.Scan() {
// 		return fault.New("failed to read plugin manifest")
// 	}

// 	line := manifestScan.Bytes()
// 	m := lib_plugin.Manifest{}
// 	if err := json.Unmarshal(line, &m); err != nil {
// 		return fault.Wrap(err, fmsg.With("failed to decode plugin manifest"))
// 	}
// 	// TODO: Validate.
// 	// s.manifest = &m

// 	s.logger.Debug("plugin started",
// 		slog.String("id", s.id.String()),
// 		slog.String("manifest", string(line)),
// 	)

// 	go func() {
// 		scan := bufio.NewScanner(or)
// 		scan.Split(bufio.ScanLines)
// 		for scan.Scan() {
// 			s.logger.Debug("recv bytes", slog.String("raw", scan.Text()))
// 			s.outchan <- scan.Bytes()
// 		}
// 	}()

// 	go func() {
// 		for {
// 			select {
// 			case <-s.ctx.Done():
// 				s.logger.Debug("writer: context cancelled",
// 					slog.String("id", s.id.String()))
// 				return

// 			case command := <-s.inchan:
// 				_, err := fmt.Fprintf(iw, "%s\n", command)
// 				if err != nil {
// 					handleError(fault.Wrap(err, fmsg.With("failed to write command to module")))
// 				}
// 				s.logger.Debug("send bytes", slog.String("raw", string(command)))
// 			}
// 		}
// 	}()

// 	s.setState(resource_plugin.ReportedStateActive, "")

// 	for {
// 		select {
// 		case <-s.ctx.Done():
// 			s.logger.Debug("reader-proc: context cancelled",
// 				slog.String("id", s.id.String()))
// 			s.setState(resource_plugin.ReportedStateInactive, "")
// 			return s.ctx.Err()

// 		case err := <-errCh:
// 			s.logger.Error("plugin failed", slog.Any("error", err))
// 			s.setState(resource_plugin.ReportedStateError, err.Error())
// 			return err

// 		case output := <-s.outchan:
// 			// TODO: Check if JSON first, if not, treat as regular log output.
// 			if output[0] != '{' {
// 				fmt.Println("LOG:", string(output))
// 				continue
// 			}

// 			var response RPCResponse
// 			if err := json.Unmarshal(output, &response); err != nil {
// 				handleError(fault.Wrap(err, fmsg.With("failed to decode response")))
// 				continue
// 			}

// 			ident, err := xid.FromString(response.ID)
// 			if err != nil {
// 				handleError(fault.Wrap(err, fmsg.With("failed to parse response ID")))
// 				continue
// 			}

// 			pending, ok := s.pendingCommands.LoadAndDelete(ident)
// 			if !ok {
// 				handleError(fault.New("received response for unknown command"))
// 				continue
// 			}

// 			s.logger.Debug("recv rpc",
// 				slog.String("id", response.ID),
// 				slog.String("method", pending.request.Method),
// 				slog.Any("params", pending.request.Params),
// 			)

// 			select {
// 			case pending.respch <- response:
// 				s.logger.Debug("send rpc response", slog.String("id", response.ID))

// 			case <-s.ctx.Done():
// 				s.logger.Debug("context cancelled while waiting for response",
// 					slog.String("id", s.id.String()),
// 					slog.String("response_id", response.ID),
// 				)
// 			}
// 		}
// 	}
// }
