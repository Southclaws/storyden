package wrun

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/puzpuzpuz/xsync/v4"
	"github.com/tetratelabs/wazero"

	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type Runner interface {
	RunOnce(ctx context.Context, bin []byte, command any) ([]byte, error)
	NewSession(ctx context.Context, bin []byte) *PluginSession
}

func New(ctx context.Context, logger *slog.Logger) Runner {
	return newWazeroRunner(ctx, logger)
}

type wazeroRunner struct {
	logger   *slog.Logger
	runtime  wazero.Runtime
	sessions *xsync.Map[string, *PluginSession]
}

func newWazeroRunner(ctx context.Context, logger *slog.Logger) Runner {
	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfigInterpreter())

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	return &wazeroRunner{
		logger:   logger,
		runtime:  r,
		sessions: xsync.NewMap[string, *PluginSession](),
	}
}

func (w *wazeroRunner) RunOnce(ctx context.Context, bin []byte, command any) ([]byte, error) {
	ir, iw := io.Pipe()
	buf := bytes.NewBuffer(nil)
	done := make(chan struct{})
	errCh := make(chan error, 2)

	mc := wazero.NewModuleConfig().
		WithStdin(ir).
		WithStdout(buf)

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
