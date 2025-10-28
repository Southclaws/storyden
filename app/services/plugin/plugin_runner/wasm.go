package plugin_runner

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	resource_plugin "github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/lib/plugin"
	"github.com/puzpuzpuz/xsync/v4"
	"github.com/rs/xid"
	"github.com/tetratelabs/wazero"

	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

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

		runType:       RunTypeBackground,
		reportedState: resource_plugin.ReportedStateInactive,

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

	s.stop()

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
	o, err := w.readPluginManifest(ctx, bin)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to validate plugin"))
	}

	m := plugin.Manifest{}
	if err := json.Unmarshal(o, &m); err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to decode manifest"))
	}

	return &m, nil
}

func (w *wazeroRunner) StartPlugin(ctx context.Context, id plugin.ID) error {
	sess, ok := w.sessions.Load(id)
	if !ok {
		return fault.New("plugin session not found")
	}

	go func() {
		if err := sess.start(); err != nil {
			w.logger.Error("plugin start failed",
				slog.String("id", id.String()),
				slog.Any("error", err),
			)
		}
	}()

	return nil
}

func (w *wazeroRunner) StopPlugin(ctx context.Context, id plugin.ID) error {
	sess, ok := w.sessions.Load(id)
	if !ok {
		return fault.New("plugin session not found")
	}

	sess.stop()

	return nil
}

func (w *wazeroRunner) readPluginManifest(ctx context.Context, bin []byte) ([]byte, error) {
	pr, pw := io.Pipe()
	manifestCh := make(chan []byte, 1)
	errCh := make(chan error, 1)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		scanner := bufio.NewScanner(pr)
		scanner.Split(bufio.ScanLines)

		if scanner.Scan() {
			manifestCh <- scanner.Bytes()
		} else if err := scanner.Err(); err != nil {
			errCh <- fault.Wrap(err, fmsg.With("failed to read manifest line"))
		} else {
			errCh <- fault.New("no output received from module: expected a manifest")
		}
	}()

	mc := wazero.NewModuleConfig().
		WithStdout(pw)

	go func() {
		mod, err := w.runtime.InstantiateWithConfig(ctx, bin, mc)
		if err != nil {
			errCh <- fault.Wrap(err, fmsg.With("failed to instantiate"))
			return
		}

		if err := mod.Close(ctx); err != nil {
			errCh <- fault.Wrap(err, fmsg.With("failed to close module"))
		}
	}()

	select {
	case manifest := <-manifestCh:
		cancel()
		pw.Close()
		return manifest, nil

	case err := <-errCh:
		cancel()
		pw.Close()
		return nil, err

	case <-ctx.Done():
		pw.Close()
		return nil, fault.Wrap(ctx.Err(), fmsg.With("context cancelled while reading manifest"))
	}
}
