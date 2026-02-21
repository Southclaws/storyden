package local

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/resources/plugin/plugin_reader"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_auth"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/plugin_logger"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/supervised_runtime"
)

type commandType int

const (
	cmdStart commandType = iota
	cmdStop
	cmdProcessStarted
	cmdProcessExited

	gracefulStopTimeout = 3 * time.Second
)

type command struct {
	typ       commandType
	cmd       *exec.Cmd
	logWriter *plugin_logger.RotatingWriter
	exitError error
	details   map[string]any
	respch    chan error
}

type localRuntime struct {
	id             plugin.InstallationID
	logger         *slog.Logger
	baseURL        url.URL
	parentCtx      context.Context
	pluginReader   *plugin_reader.Reader
	dataPath       string
	pluginLogger   *plugin_logger.Writer
	bin            []byte
	manifest       *plugin.Validated
	commandCh      chan command
	eventCh        chan supervised_runtime.Event
	maxBackoff     time.Duration
	maxRestarts    int
	crashCutoff    time.Duration
	runtimeBackoff time.Duration

	procCtx    context.Context
	procCancel context.CancelFunc
}

func newRuntime(
	id plugin.InstallationID,
	baseURL url.URL,
	parentCtx context.Context,
	bin []byte,
	manifest *plugin.Validated,
	parentLogger *slog.Logger,
	pluginLogger *plugin_logger.Writer,
	pluginReader *plugin_reader.Reader,
	dataPath string,
	maxRestartAttempts int,
	maxBackoffDuration time.Duration,
	runtimeCrashThreshold time.Duration,
	runtimeCrashBackoff time.Duration,
) *localRuntime {
	logger := parentLogger.With(slog.String("plugin_id", id.String()))

	r := &localRuntime{
		id:             id,
		baseURL:        baseURL,
		parentCtx:      parentCtx,
		logger:         logger,
		bin:            bin,
		manifest:       manifest,
		pluginLogger:   pluginLogger,
		pluginReader:   pluginReader,
		dataPath:       dataPath,
		maxRestarts:    maxRestartAttempts,
		maxBackoff:     maxBackoffDuration,
		crashCutoff:    runtimeCrashThreshold,
		runtimeBackoff: runtimeCrashBackoff,
		commandCh:      make(chan command, 8),
		eventCh:        make(chan supervised_runtime.Event, 32),
	}

	go r.supervisor()

	return r
}

func (r *localRuntime) Events() <-chan supervised_runtime.Event {
	return r.eventCh
}

func (r *localRuntime) Start(ctx context.Context) error {
	respch := make(chan error, 1)

	select {
	case r.commandCh <- command{typ: cmdStart, respch: respch}:
	case <-ctx.Done():
		return ctx.Err()
	}

	select {
	case err := <-respch:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *localRuntime) Stop(ctx context.Context) error {
	respch := make(chan error, 1)

	select {
	case r.commandCh <- command{typ: cmdStop, respch: respch}:
	case <-ctx.Done():
		return ctx.Err()
	}

	select {
	case err := <-respch:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *localRuntime) supervisor() {
	var (
		cmd           *exec.Cmd
		running       bool
		stopping      bool
		stopWaiter    chan error
		processStarts time.Time
		restartCount  int
	)

	clearProcessCtx := func() {
		if r.procCancel != nil {
			r.procCancel()
			r.procCancel = nil
		}
		r.procCtx = nil
	}
	defer clearProcessCtx()

	for c := range r.commandCh {
		switch c.typ {
		case cmdStart:
			if running {
				c.respch <- nil
				continue
			}

			clearProcessCtx()

			running = true
			stopping = false
			restartCount = 0
			processStarts = time.Now()
			r.procCtx, r.procCancel = context.WithCancel(r.parentCtx)

			go r.runProcess(r.procCtx)

			c.respch <- nil

		case cmdStop:
			if !running {
				clearProcessCtx()
				r.reportState(plugin.ReportedStateInactive, "", nil)
				if c.respch != nil {
					c.respch <- nil
				}
				continue
			}

			r.reportState(plugin.ReportedStateStopping, "", nil)

			r.stopProcess(cmd)

			stopping = true
			stopWaiter = c.respch

		case cmdProcessStarted:
			if !running || stopping {
				continue
			}
			cmd = c.cmd

		case cmdProcessExited:
			if !running {
				continue
			}

			if c.logWriter != nil {
				if err := c.logWriter.Rotator.Rotate(); err != nil {
					r.logger.Error("failed to rotate log after process exit", slog.Any("error", err))
				}
			}

			if stopping {
				running = false
				stopping = false
				cmd = nil
				clearProcessCtx()
				r.reportState(plugin.ReportedStateInactive, "", nil)
				if stopWaiter != nil {
					stopWaiter <- nil
					stopWaiter = nil
				}
				continue
			}

			// Process exited because its context was canceled outside an explicit
			// stop transition (for example host/application shutdown). Treat this
			// as a normal non-crash shutdown and do not enter restart flow.
			if c.exitError == nil {
				running = false
				cmd = nil
				clearProcessCtx()
				r.reportState(plugin.ReportedStateInactive, "", nil)
				continue
			}

			processUptime := time.Since(processStarts)
			isStartupCrash := processUptime < r.crashCutoff

			if restartCount >= r.maxRestarts {
				r.logger.Error("plugin exceeded max restart attempts",
					slog.Int("restart_count", restartCount),
					slog.Any("error", c.exitError))
				r.reportState(plugin.ReportedStateError, c.exitError.Error(), c.details)
				running = false
				cmd = nil
				clearProcessCtx()
				continue
			}

			backoffSeconds := 1 << restartCount
			maxBackoffSeconds := int(r.maxBackoff.Seconds())
			if backoffSeconds > maxBackoffSeconds {
				backoffSeconds = maxBackoffSeconds
			}
			backoffDuration := time.Duration(backoffSeconds) * time.Second

			crashType := "startup"
			if !isStartupCrash {
				crashType = "runtime"
			}

			r.logger.Warn("plugin crashed, restarting with backoff",
				slog.String("crash_type", crashType),
				slog.Int("restart_count", restartCount),
				slog.Duration("backoff", backoffDuration),
				slog.Duration("uptime", processUptime),
				slog.Any("error", c.exitError))

			r.reportState(plugin.ReportedStateRestarting, c.exitError.Error(), c.details)

			restartCount++
			processStarts = time.Now()

			select {
			case <-time.After(backoffDuration):
			case <-r.procCtx.Done():
				r.logger.Info("plugin stop requested during backoff")
				running = false
				stopping = false
				cmd = nil
				clearProcessCtx()
				r.reportState(plugin.ReportedStateInactive, "", nil)
				if stopWaiter != nil {
					stopWaiter <- nil
					stopWaiter = nil
				}
				continue
			}

			r.logger.Info("attempting to restart plugin", slog.Int("attempt", restartCount))
			go r.runProcess(r.procCtx)
		}
	}
}

func (r *localRuntime) stopProcess(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}

	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		if errors.Is(err, os.ErrProcessDone) {
			return
		}

		r.logger.Warn("failed to send interrupt signal, falling back to kill", slog.Any("error", err))
		if killErr := cmd.Process.Kill(); killErr != nil && !errors.Is(killErr, os.ErrProcessDone) {
			r.logger.Error("failed to kill process after interrupt failure", slog.Any("error", killErr))
		}
		return
	}

	r.logger.Info("sent interrupt signal to plugin process")

	proc := cmd.Process
	go func() {
		timer := time.NewTimer(gracefulStopTimeout)
		defer timer.Stop()
		<-timer.C

		if err := proc.Kill(); err != nil {
			if !errors.Is(err, os.ErrProcessDone) {
				r.logger.Warn("failed to kill process after graceful shutdown timeout", slog.Any("error", err))
			}
			return
		}

		r.logger.Warn("plugin process did not exit after interrupt, killed",
			slog.Duration("timeout", gracefulStopTimeout))
	}()
}

func (r *localRuntime) runProcess(ctx context.Context) {
	cmd, logWriter, err := r.startProcess(ctx)
	if err != nil {
		r.logger.Error("failed to start process", slog.Any("error", err))
		r.commandCh <- command{
			typ:       cmdProcessExited,
			exitError: err,
			details:   runtimeErrorDetails("start", err),
		}
		return
	}

	r.commandCh <- command{typ: cmdProcessStarted, cmd: cmd}

	waitErr := cmd.Wait()

	if ctx.Err() != nil {
		r.logger.Info("plugin stopped")
		r.commandCh <- command{typ: cmdProcessExited, logWriter: logWriter}
		return
	}

	if waitErr == nil {
		waitErr = fault.New("process exited cleanly")
	}

	r.commandCh <- command{
		typ:       cmdProcessExited,
		logWriter: logWriter,
		exitError: waitErr,
		details:   runtimeErrorDetails("exit", waitErr),
	}
}

func (r *localRuntime) startProcess(ctx context.Context) (*exec.Cmd, *plugin_logger.RotatingWriter, error) {
	workdir, err := extractSDXArchive(r.bin, r.dataPath, r.id)
	if err != nil {
		return nil, nil, fault.Wrap(err, fmsg.With("failed to extract plugin"))
	}

	authSecret, err := r.pluginReader.GetAuthSecret(ctx, r.id)
	if err != nil {
		return nil, nil, fault.Wrap(err, fmsg.With("failed to get auth secret"))
	}

	rpcURL, err := plugin_auth.BuildConnectionURL(r.baseURL, r.id, authSecret)
	if err != nil {
		return nil, nil, fault.Wrap(err, fmsg.With("failed to build connection URL"))
	}

	if r.manifest.Metadata.Command == "" {
		return nil, nil, fault.New("command cannot be empty")
	}

	cmd := exec.CommandContext(ctx, r.manifest.Metadata.Command, r.manifest.Metadata.Args...)
	cmd.Dir = workdir
	env := os.Environ()
	env = append(env, "STORYDEN_RPC_URL="+rpcURL.String())
	cmd.Env = env

	logWriter, err := r.pluginLogger.NewWriter(r.dataPath, r.id)
	if err != nil {
		return nil, nil, fault.Wrap(err, fmsg.With("failed to create plugin log writer"))
	}

	cmd.Stdout = logWriter
	cmd.Stderr = logWriter

	if err := cmd.Start(); err != nil {
		return nil, nil, fault.Wrap(err, fmsg.With("failed to start process"))
	}

	return cmd, logWriter, nil
}

func (r *localRuntime) reportState(state plugin.ReportedState, message string, details map[string]any) {
	r.eventCh <- supervised_runtime.Event{
		State:   state,
		Message: message,
		Details: details,
	}
}

func runtimeErrorDetails(stage string, err error) map[string]any {
	details := map[string]any{
		"runtime_provider": "local",
		"stage":            stage,
		"error":            err.Error(),
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		details["exit_code"] = exitErr.ExitCode()
	}

	return details
}

func extractSDXArchive(bin []byte, dataPath string, pluginID plugin.InstallationID) (string, error) {
	zr, err := zip.NewReader(bytes.NewReader(bin), int64(len(bin)))
	if err != nil {
		return "", fault.Wrap(err, fmsg.With("failed to open archive"))
	}

	workdir := filepath.Join(dataPath, pluginID.String())
	if err := os.MkdirAll(workdir, 0o755); err != nil {
		return "", fault.Wrap(err, fmsg.With("failed to create plugin directory"))
	}

	for _, file := range zr.File {
		targetPath, err := joinWithin(workdir, file.Name)
		if err != nil {
			return "", err
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, 0o755); err != nil {
				return "", fault.Wrap(err, fmsg.With("failed to create dir"))
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
			return "", fault.Wrap(err, fmsg.With("failed to create parent dir"))
		}

		rc, err := file.Open()
		if err != nil {
			return "", fault.Wrap(err, fmsg.With("failed to open file in archive"))
		}

		mode := file.FileInfo().Mode()

		out, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
		if err != nil {
			rc.Close()
			return "", fault.Wrap(err, fmsg.With("failed to create extracted file"))
		}

		n, err := io.Copy(out, rc)
		if err != nil {
			out.Close()
			rc.Close()
			return "", fault.Wrap(err, fmsg.With(fmt.Sprintf("failed to copy archive contents for file %s (wrote %d bytes, expected %d)", file.Name, n, file.UncompressedSize64)))
		}

		out.Close()
		rc.Close()
	}

	return workdir, nil
}

func joinWithin(base, name string) (string, error) {
	clean := filepath.Join(base, name)
	if !strings.HasPrefix(clean, base+string(os.PathSeparator)) && clean != base {
		return "", fault.New("archive entry escapes extraction directory")
	}
	return clean, nil
}
