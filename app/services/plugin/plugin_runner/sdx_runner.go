package plugin_runner

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	resource_plugin "github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/lib/plugin"
	"github.com/puzpuzpuz/xsync/v4"
	"github.com/rs/xid"
)

type sdxRunner struct {
	logger   *slog.Logger
	sessions *xsync.Map[plugin.ID, *sdxSession]
}

type sdxSession struct {
	*PluginSession

	mu       sync.Mutex
	workdir  string
	cmd      *exec.Cmd
	stdin    io.WriteCloser
	stopping bool
}

func newSDXRunner(logger *slog.Logger) Runner {
	return &sdxRunner{
		logger:   logger,
		sessions: xsync.NewMap[plugin.ID, *sdxSession](),
	}
}

func (r *sdxRunner) Load(ctx context.Context, bin []byte) (*PluginSession, error) {
	manifest, err := r.Validate(ctx, bin)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	sess := &PluginSession{
		id:     manifest.ID,
		logger: r.logger,

		runType:       RunTypeBackground,
		reportedState: resource_plugin.ReportedStateInactive,

		bin:      bin,
		manifest: manifest,

		inchan:          make(chan []byte),
		outchan:         make(chan []byte),
		pendingCommands: xsync.NewMap[xid.ID, pendingCommand](),
	}

	sdxSess := &sdxSession{PluginSession: sess}

	r.sessions.Store(manifest.ID, sdxSess)

	return sess, nil
}

func (r *sdxRunner) Unload(ctx context.Context, id plugin.ID) error {
	sess, ok := r.sessions.Load(id)
	if !ok {
		return fault.New("plugin session not found")
	}

	if err := r.StopPlugin(ctx, id); err != nil {
		r.logger.Warn("failed to stop plugin during unload", slog.String("id", id.String()), slog.Any("error", err))
	}

	r.sessions.Delete(id)

	sess.mu.Lock()
	if sess.workdir != "" {
		_ = os.RemoveAll(sess.workdir)
		sess.workdir = ""
	}
	sess.mu.Unlock()

	return nil
}

func (r *sdxRunner) Validate(ctx context.Context, bin []byte) (*plugin.Manifest, error) {
	manifestBytes, err := readSDXManifest(bin)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to read manifest"))
	}

	m := plugin.Manifest{}
	if err := json.Unmarshal(manifestBytes, &m); err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to decode manifest"))
	}

	if err := ensureExecutablePresent(bin); err != nil {
		return nil, fault.Wrap(err, fmsg.With("plugin binary missing"))
	}

	return &m, nil
}

func (r *sdxRunner) GetSession(ctx context.Context, id plugin.ID) (*PluginSession, error) {
	sess, ok := r.sessions.Load(id)
	if !ok {
		return nil, fault.New("plugin session not found")
	}
	return sess.PluginSession, nil
}

func (r *sdxRunner) GetSessions(ctx context.Context) ([]*PluginSession, error) {
	out := []*PluginSession{}
	r.sessions.Range(func(_ plugin.ID, sess *sdxSession) bool {
		out = append(out, sess.PluginSession)
		return true
	})
	return out, nil
}

func (r *sdxRunner) StartPlugin(ctx context.Context, id plugin.ID) error {
	sess, ok := r.sessions.Load(id)
	if !ok {
		return fault.New("plugin session not found")
	}

	go func() {
		if err := sess.runProcess(); err != nil {
			r.logger.Error("plugin start failed", slog.String("id", id.String()), slog.Any("error", err))
		}
	}()

	return nil
}

func (r *sdxRunner) StopPlugin(ctx context.Context, id plugin.ID) error {
	sess, ok := r.sessions.Load(id)
	if !ok {
		return fault.New("plugin session not found")
	}

	sess.mu.Lock()
	if sess.cmd == nil {
		sess.mu.Unlock()
		sess.stop()
		return nil
	}

	sess.stopping = true
	cmd := sess.cmd
	stdin := sess.stdin
	sess.mu.Unlock()

	if stdin != nil {
		_ = stdin.Close()
	}

	if cmd.ProcessState == nil || !cmd.ProcessState.Exited() {
		if err := cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			return fault.Wrap(err, fmsg.With("failed to stop plugin process"))
		}
	}

	sess.stop()

	if sess.ctx != nil {
		select {
		case <-sess.ctx.Done():
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func (s *sdxSession) runProcess() error {
	s.stateMu.RLock()
	currentState := s.reportedState
	s.stateMu.RUnlock()

	if currentState == resource_plugin.ReportedStateActive {
		return fault.New("plugin is already running")
	}

	if s.runType != RunTypeBackground {
		return fault.New("only background plugins can be started")
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.ctx = ctx
	s.cancel = cancel

	workdir, binaryPath, err := extractSDXArchive(s.bin)
	if err != nil {
		cancel()
		return fault.Wrap(err, fmsg.With("failed to extract plugin"))
	}

	cmd := exec.CommandContext(ctx, binaryPath)
	cmd.Dir = filepath.Dir(binaryPath)
	cmd.Env = os.Environ()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return fault.Wrap(err, fmsg.With("failed to get stdout"))
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		cancel()
		return fault.Wrap(err, fmsg.With("failed to get stderr"))
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		cancel()
		return fault.Wrap(err, fmsg.With("failed to get stdin"))
	}

	errCh := make(chan error, 1)
	sendErr := func(err error) {
		select {
		case errCh <- err:
		default:
		}
	}

	if err := cmd.Start(); err != nil {
		cancel()
		return fault.Wrap(err, fmsg.With("failed to start process"))
	}

	s.mu.Lock()
	s.cmd = cmd
	s.stdin = stdin
	s.workdir = workdir
	s.stopping = false
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		if s.workdir != "" {
			_ = os.RemoveAll(s.workdir)
			s.workdir = ""
		}
		s.cmd = nil
		s.stdin = nil
		s.mu.Unlock()
		cancel()
	}()

	handleError := func(err error) {
		s.logger.Error("plugin error", slog.Any("error", err))
		s.setState(resource_plugin.ReportedStateError, err.Error())
	}

	go func() {
		defer stdin.Close()
		for {
			select {
			case <-ctx.Done():
				return
			case payload, ok := <-s.inchan:
				if !ok {
					return
				}
				if len(payload) == 0 {
					continue
				}
				if _, err := fmt.Fprintf(stdin, "%s\n", payload); err != nil {
					if ctx.Err() == nil {
						sendErr(fault.Wrap(err, fmsg.With("failed to write command")))
					}
					return
				}
				s.logger.Debug("send bytes", slog.String("raw", string(payload)))
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			if ctx.Err() != nil {
				return
			}
			line := append([]byte(nil), scanner.Bytes()...)
			select {
			case <-ctx.Done():
				return
			case s.outchan <- line:
				s.logger.Debug("recv bytes", slog.String("raw", string(line)))
			}
		}
		if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) {
			if ctx.Err() == nil {
				sendErr(fault.Wrap(err, fmsg.With("failed to read stdout")))
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			if ctx.Err() != nil {
				return
			}
			s.logger.Warn("plugin stderr", slog.String("line", scanner.Text()))
		}
		if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) {
			s.logger.Error("stderr read error", slog.Any("error", err))
		}
	}()

	go func() {
		waitErr := cmd.Wait()
		s.mu.Lock()
		stopping := s.stopping
		s.mu.Unlock()
		if stopping {
			sendErr(nil)
			return
		}
		if waitErr != nil {
			sendErr(fault.Wrap(waitErr, fmsg.With("process exited")))
			return
		}
		sendErr(fault.New("plugin exited unexpectedly"))
	}()

	s.setState(resource_plugin.ReportedStateActive, "")

	for {
		select {
		case <-ctx.Done():
			s.logger.Debug("context cancelled", slog.String("id", s.id.String()))
			s.setState(resource_plugin.ReportedStateInactive, "")
			return ctx.Err()
		case err := <-errCh:
			if err == nil {
				s.logger.Info("plugin stopped", slog.String("id", s.id.String()))
				s.setState(resource_plugin.ReportedStateInactive, "")
				return nil
			}
			handleError(err)
			return err
		case output, ok := <-s.outchan:
			if !ok {
				return nil
			}
			if len(output) == 0 {
				continue
			}
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

			s.logger.Debug("recv rpc", slog.String("id", response.ID), slog.String("method", pending.request.Method), slog.Any("params", pending.request.Params))

			select {
			case pending.respch <- response:
				s.logger.Debug("send rpc response", slog.String("id", response.ID))
			case <-ctx.Done():
				s.logger.Debug("context cancelled while waiting for response", slog.String("id", s.id.String()), slog.String("response_id", response.ID))
			}
		}
	}
}

func readSDXManifest(bin []byte) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(bin), int64(len(bin)))
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to open archive"))
	}

	for _, file := range zr.File {
		if filepath.Base(file.Name) != "manifest.json" {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			return nil, fault.Wrap(err, fmsg.With("failed to open manifest"))
		}
		defer rc.Close()

		data, err := io.ReadAll(rc)
		if err != nil {
			return nil, fault.Wrap(err, fmsg.With("failed to read manifest"))
		}

		return data, nil
	}

	return nil, fault.New("manifest.json not found in plugin archive")
}

func ensureExecutablePresent(bin []byte) error {
	zr, err := zip.NewReader(bytes.NewReader(bin), int64(len(bin)))
	if err != nil {
		return fault.Wrap(err, fmsg.With("failed to open archive"))
	}

	for _, file := range zr.File {
		base := filepath.Base(file.Name)
		if base == "main" || base == "main.exe" {
			return nil
		}
	}

	return fault.New("plugin archive missing main executable")
}

func extractSDXArchive(bin []byte) (string, string, error) {
	zr, err := zip.NewReader(bytes.NewReader(bin), int64(len(bin)))
	if err != nil {
		return "", "", fault.Wrap(err, fmsg.With("failed to open archive"))
	}

	workdir, err := os.MkdirTemp("", "storyden-sdx-")
	if err != nil {
		return "", "", fault.Wrap(err, fmsg.With("failed to create temp dir"))
	}

	var binaryPath string

	for _, file := range zr.File {
		targetPath, err := joinWithin(workdir, file.Name)
		if err != nil {
			return "", "", err
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, 0o755); err != nil {
				return "", "", fault.Wrap(err, fmsg.With("failed to create dir"))
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
			return "", "", fault.Wrap(err, fmsg.With("failed to create parent dir"))
		}

		rc, err := file.Open()
		if err != nil {
			return "", "", fault.Wrap(err, fmsg.With("failed to open file in archive"))
		}

		mode := file.FileInfo().Mode()
		if mode&0o111 == 0 && (filepath.Base(file.Name) == "main" || filepath.Base(file.Name) == "main.exe") {
			mode |= 0o755
		}

		out, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
		if err != nil {
			rc.Close()
			return "", "", fault.Wrap(err, fmsg.With("failed to create extracted file"))
		}

		if _, err := io.Copy(out, rc); err != nil {
			out.Close()
			rc.Close()
			return "", "", fault.Wrap(err, fmsg.With("failed to copy archive contents"))
		}

		out.Close()
		rc.Close()

		base := filepath.Base(file.Name)
		if base == "main" || base == "main.exe" {
			binaryPath = targetPath
		}
	}

	if binaryPath == "" {
		return "", "", fault.New("extracted archive missing main executable")
	}

	return workdir, binaryPath, nil
}

func joinWithin(base, name string) (string, error) {
	clean := filepath.Join(base, name)
	if !strings.HasPrefix(clean, base+string(os.PathSeparator)) && clean != base {
		return "", fault.New("archive entry escapes extraction directory")
	}
	return clean, nil
}
