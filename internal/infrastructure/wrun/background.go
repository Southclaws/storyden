package wrun

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/puzpuzpuz/xsync/v4"
	"github.com/rs/xid"
	"github.com/tetratelabs/wazero"

	"github.com/Southclaws/storyden/lib/plugin"
)

// TODO: Make this configurable in some way. Perhaps at the instance level or
// perhaps at the plugin level - if a plugin declares it implements particularly
// lengthy operations (such as an AI-driven plugin or something similar.)
var defaultRPCTimeout = time.Second * 30

type PluginSession struct {
	logger          *slog.Logger
	id              string
	bin             []byte
	runtime         wazero.Runtime
	inchan          chan []byte
	outchan         chan []byte
	pendingCommands *xsync.Map[xid.ID, pendingCommand]
}

type pendingCommand struct {
	request RPCRequest
	sent    time.Time
	respch  chan RPCResponse
}

type RPCRequest struct {
	ID     string `json:"id"`
	Method string `json:"method"`
	Params any    `json:"params"`
}

type RPCResponse struct {
	ID     string `json:"id"`
	Result any    `json:"result"`
}

func (r *wazeroRunner) NewSession(ctx context.Context, bin []byte) *PluginSession {
	id := xid.New().String()

	s := &PluginSession{
		logger:          r.logger,
		id:              id,
		bin:             bin,
		runtime:         r.runtime,
		inchan:          make(chan []byte, 1),
		outchan:         make(chan []byte, 1),
		pendingCommands: xsync.NewMap[xid.ID, pendingCommand](),
	}

	return s
}

func (s PluginSession) Send(ctx context.Context, method string, params any) (any, error) {
	id := xid.New()
	request := RPCRequest{
		ID:     id.String(),
		Method: method,
		Params: params,
	}

	pending := pendingCommand{
		request: request,
		sent:    time.Now(),
		respch:  make(chan RPCResponse),
	}

	s.pendingCommands.Store(id, pending)

	b, err := json.Marshal(request)
	if err != nil {
		panic(fault.Wrap(err, fmsg.With("failed to encode command")))
	}

	if _, ok := ctx.Deadline(); !ok {
		// If the context doesn't have a deadline, set a default timeout.
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, defaultRPCTimeout)
		defer cancel()
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	case s.inchan <- b:
		s.logger.Debug("send rpc",
			slog.String("id", id.String()),
			slog.String("method", method),
			slog.Any("params", params),
		)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	case resp := <-pending.respch:
		return resp.Result, nil
	}
}

func (s *PluginSession) Stop() {
	// Close channels to signal shutdown
	if s.inchan != nil {
		close(s.inchan)
	}
}

func (s *PluginSession) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	ir, iw := io.Pipe()
	or, ow := io.Pipe()

	errCh := make(chan error)

	mc := wazero.NewModuleConfig().
		WithStdin(ir).
		WithStdout(ow).
		WithStderr(os.Stderr) // TODO: Route error to persisted plugin logs

	handleError := func(err error) {
		fmt.Println("ERROR:", err)
	}

	go func() {
		defer cancel()

		mod, err := s.runtime.InstantiateWithConfig(ctx, s.bin, mc)
		if err != nil {
			errCh <- fault.Wrap(err, fmsg.With("failed to instantiate"))
			return
		}

		if err := mod.Close(ctx); err != nil {
			errCh <- fault.Wrap(err, fmsg.With("failed to close plugin"))
			return
		}

		errCh <- fault.New("plugin exited unexpectedly")
	}()

	manifestScan := bufio.NewScanner(or)
	manifestScan.Split(bufio.ScanLines)
	if !manifestScan.Scan() {
		return fault.New("failed to read plugin manifest")
	}

	line := manifestScan.Text()

	var manifest plugin.Manifest
	if err := json.Unmarshal([]byte(line), &manifest); err != nil {
		return fault.Wrap(err, fmsg.With("failed to parse plugin manifest"))
	}

	fmt.Println("LOADED", manifest)
	s.logger.Debug("loaded plugin manifest", slog.String("name", manifest.Name.String()), slog.String("version", manifest.Version.String()))

	go func() {
		scan := bufio.NewScanner(or)
		scan.Split(bufio.ScanLines)
		for scan.Scan() {
			s.logger.Debug("recv bytes", slog.String("raw", scan.Text()))
			s.outchan <- scan.Bytes()
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				s.logger.Debug("writer: context cancelled",
					slog.String("id", s.id))
				return

			case command := <-s.inchan:
				_, err := fmt.Fprintf(iw, "%s\n", command)
				if err != nil {
					handleError(fault.Wrap(err, fmsg.With("failed to write command to module")))
				}
				s.logger.Debug("send bytes", slog.String("raw", string(command)))
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			s.logger.Debug("reader-proc: context cancelled",
				slog.String("id", s.id))
			return ctx.Err()

		case err := <-errCh:
			s.logger.Error("plugin failed", slog.Any("error", err))
			return err

		case output := <-s.outchan:
			// TODO: Check if JSON first, if not, treat as regular log output.
			if output[0] != '{' {
				fmt.Println("LOG:", string(output))
				continue
			}

			var response RPCResponse
			if err := json.Unmarshal(output, &response); err != nil {
				handleError(fault.Wrap(err, fmsg.With("failed to decode response")))
				continue
			}

			ident, err := xid.FromString(response.ID)
			if err != nil {
				handleError(fault.Wrap(err, fmsg.With("failed to parse response ID")))
				continue
			}

			pending, ok := s.pendingCommands.LoadAndDelete(ident)
			if !ok {
				handleError(fault.New("received response for unknown command"))
				continue
			}

			s.logger.Debug("recv rpc",
				slog.String("id", response.ID),
				slog.String("method", pending.request.Method),
				slog.Any("params", pending.request.Params),
			)

			select {
			case pending.respch <- response:
				s.logger.Debug("send rpc response", slog.String("id", response.ID))

			case <-ctx.Done():
				s.logger.Debug("context cancelled while waiting for response",
					slog.String("id", s.id),
					slog.String("response_id", response.ID),
				)
			}
		}
	}
}
