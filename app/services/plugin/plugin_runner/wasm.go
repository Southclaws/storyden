package plugin_runner

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/storyden/lib/plugin"
	"github.com/puzpuzpuz/xsync/v4"
	"github.com/rs/xid"
	"github.com/tetratelabs/wazero"

	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type Runner interface {
	Load(ctx context.Context, bin []byte) (*PluginSession, error)
	Unload(ctx context.Context, id plugin.ID) error
	Validate(ctx context.Context, bin []byte) (*plugin.Manifest, error)
	GetSession(ctx context.Context, id plugin.ID) (*PluginSession, error)
	GetSessions(ctx context.Context) ([]*PluginSession, error)
}

func New(ctx context.Context, logger *slog.Logger) Runner {
	return newWazeroRunner(ctx, logger)
}

type wazeroRunner struct {
	logger   *slog.Logger
	runtime  wazero.Runtime
	sessions *xsync.Map[plugin.ID, *PluginSession]
}

func newWazeroRunner(ctx context.Context, logger *slog.Logger) Runner {
	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfigInterpreter())

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	return &wazeroRunner{
		logger:   logger,
		runtime:  r,
		sessions: xsync.NewMap[plugin.ID, *PluginSession](),
	}
}

func (w *wazeroRunner) Load(ctx context.Context, bin []byte) (*PluginSession, error) {
	m, err := w.Validate(ctx, bin)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	key := m.ID

	session := &PluginSession{
		id:     key,
		logger: w.logger,

		bin:      bin,
		manifest: m,

		runtime: w.runtime,
		runner:  w,

		inchan:          make(chan []byte),
		outchan:         make(chan []byte),
		pendingCommands: xsync.NewMap[xid.ID, pendingCommand](),
	}

	w.sessions.Store(key, session)

	return session, nil
}

func (w *wazeroRunner) Unload(ctx context.Context, id plugin.ID) error {
	s, ok := w.sessions.Load(id)
	if !ok {
		return fault.New("plugin session not found")
	}

	s.Stop()

	w.sessions.Delete(id)
	return nil
}

func (w *wazeroRunner) GetSession(ctx context.Context, id plugin.ID) (*PluginSession, error) {
	s, ok := w.sessions.Load(id)
	if !ok {
		return nil, fault.New("plugin session not found")
	}
	return s, nil
}

func (w *wazeroRunner) GetSessions(ctx context.Context) ([]*PluginSession, error) {
	sessions := []*PluginSession{}
	w.sessions.Range(func(key plugin.ID, s *PluginSession) bool {
		sessions = append(sessions, s)
		return true
	})
	return sessions, nil
}

func (w *wazeroRunner) Validate(ctx context.Context, bin []byte) (*plugin.Manifest, error) {
	o, err := w.runOnce(ctx, bin, nil)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to validate plugin"))
	}

	m := plugin.Manifest{}
	if err := json.Unmarshal(o, &m); err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to decode manifest"))
	}

	return &m, nil
}

func (w *wazeroRunner) runOnce(ctx context.Context, bin []byte, command any) ([]byte, error) {
	ir, iw := io.Pipe()
	buf := bytes.NewBuffer(nil)
	done := make(chan struct{})
	errCh := make(chan error, 2)

	mc := wazero.NewModuleConfig().
		WithStdin(ir).
		WithStdout(buf)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		defer close(done)

		mod, err := w.runtime.InstantiateWithConfig(ctx, bin, mc)
		if err != nil {
			errCh <- fault.Wrap(err, fmsg.With("failed to instantiate"))
			return
		}

		if err := mod.Close(ctx); err != nil {
			errCh <- fault.Wrap(err, fmsg.With("failed to close module"))
			return
		}

		done <- struct{}{}
	}()

	if command != nil {
		go func() {
			defer iw.Close()

			cb, err := json.Marshal(command)
			if err != nil {
				errCh <- fault.Wrap(err, fmsg.With("failed to encode command"))
				return
			}

			_, err = iw.Write(cb)
			if err != nil {
				errCh <- fault.Wrap(err, fmsg.With("failed to write command to module"))
				return
			}
		}()
	}

	select {
	case <-done:
	case err := <-errCh:
		return nil, err
	}

	// Shut down the plugin.
	// TODO: Background plugins may try to init, we should do some
	// kind of "start" RPC to allow booting up post manifest.
	cancel()

	outputs := [][]byte{}

	s := bufio.NewScanner(buf)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		outputs = append(outputs, s.Bytes())
	}

	o := outputs[len(outputs)-1]

	if len(outputs) == 0 {
		if command == nil {
			return nil, fault.New("no output received from module: expected a manifest")
		} else {
			return nil, fault.New("no output received from module: expected a command response")
		}
	}

	return o, nil
}
