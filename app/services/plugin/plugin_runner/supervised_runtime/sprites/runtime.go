package sprites

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/superfly/sprites-go"

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
	spriteWorkdir       = "/plugin"
	spriteNamePrefix    = "storyden-plugin-"
)

type command struct {
	typ       commandType
	cmd       *sprites.Cmd
	logWriter *plugin_logger.RotatingWriter
	exitError error
	details   map[string]any
	respch    chan error
}

type archiveExtractStats struct {
	Files       int
	Directories int
	Bytes       int64
}

type runtime struct {
	id        plugin.InstallationID
	logger    *slog.Logger
	baseURL   url.URL
	parentCtx context.Context

	client       *sprites.Client
	pluginReader *plugin_reader.Reader
	pluginLogger *plugin_logger.Writer

	bin      []byte
	manifest *plugin.Validated
	dataPath string

	commandCh chan command
	eventCh   chan supervised_runtime.Event

	maxBackoff     time.Duration
	maxRestarts    int
	crashCutoff    time.Duration
	runtimeBackoff time.Duration

	procCtx    context.Context
	procCancel context.CancelFunc

	spriteName string
}

func newRuntime(
	id plugin.InstallationID,
	baseURL url.URL,
	parentCtx context.Context,
	bin []byte,
	manifest *plugin.Validated,
	parentLogger *slog.Logger,
	client *sprites.Client,
	pluginLogger *plugin_logger.Writer,
	pluginReader *plugin_reader.Reader,
	dataPath string,
	maxRestartAttempts int,
	maxBackoffDuration time.Duration,
	runtimeCrashThreshold time.Duration,
	runtimeCrashBackoff time.Duration,
) *runtime {
	logger := parentLogger.With(slog.String("plugin_id", id.String()))

	r := &runtime{
		id:             id,
		logger:         logger,
		baseURL:        baseURL,
		parentCtx:      parentCtx,
		client:         client,
		pluginReader:   pluginReader,
		pluginLogger:   pluginLogger,
		bin:            bin,
		manifest:       manifest,
		dataPath:       dataPath,
		commandCh:      make(chan command, 8),
		eventCh:        make(chan supervised_runtime.Event, 32),
		maxBackoff:     maxBackoffDuration,
		maxRestarts:    maxRestartAttempts,
		crashCutoff:    runtimeCrashThreshold,
		runtimeBackoff: runtimeCrashBackoff,
		spriteName:     spriteNamePrefix + id.String(),
	}

	go r.supervisor()

	return r
}

func (r *runtime) Events() <-chan supervised_runtime.Event {
	return r.eventCh
}

func (r *runtime) Start(ctx context.Context) error {
	respch := make(chan error, 1)
	r.logger.Info("runtime start requested")

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

func (r *runtime) Stop(ctx context.Context) error {
	respch := make(chan error, 1)
	r.logger.Info("runtime stop requested")

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

func (r *runtime) supervisor() {
	var (
		cmd           *sprites.Cmd
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
				r.logger.Debug("start requested while already running")
				c.respch <- nil
				continue
			}

			clearProcessCtx()

			running = true
			stopping = false
			restartCount = 0
			processStarts = time.Now()
			r.procCtx, r.procCancel = context.WithCancel(r.parentCtx)
			r.logger.Info("starting plugin process supervisor")

			go r.runProcess(r.procCtx)
			c.respch <- nil

		case cmdStop:
			if !running {
				r.logger.Debug("stop requested while not running")
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

			backoffDuration := r.runtimeBackoff
			if isStartupCrash {
				backoffSeconds := 1 << restartCount
				maxBackoffSeconds := int(r.maxBackoff.Seconds())
				if backoffSeconds > maxBackoffSeconds {
					backoffSeconds = maxBackoffSeconds
				}
				backoffDuration = time.Duration(backoffSeconds) * time.Second
			}

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

func (r *runtime) stopProcess(cmd *sprites.Cmd) {
	if cmd == nil {
		return
	}

	if err := cmd.Signal("INT"); err != nil && !isAlreadyStoppedError(err) {
		r.logger.Warn("failed to send interrupt signal, falling back to kill", slog.Any("error", err))
		if killErr := cmd.Signal("KILL"); killErr != nil && !isAlreadyStoppedError(killErr) {
			r.logger.Error("failed to kill process after interrupt failure", slog.Any("error", killErr))
		}
		return
	}

	r.logger.Info("sent interrupt signal to plugin process")

	go func() {
		timer := time.NewTimer(gracefulStopTimeout)
		defer timer.Stop()
		<-timer.C

		if err := cmd.Signal("KILL"); err != nil {
			if !isAlreadyStoppedError(err) {
				r.logger.Warn("failed to kill process after graceful shutdown timeout", slog.Any("error", err))
			}
			return
		}

		r.logger.Warn("plugin process did not exit after interrupt, killed",
			slog.Duration("timeout", gracefulStopTimeout))
	}()
}

func (r *runtime) runProcess(ctx context.Context) {
	r.logger.Info("starting plugin process")
	cmd, logWriter, err := r.startProcess(ctx)
	if err != nil {
		r.logger.Error("failed to start process", slog.Any("error", err))
		r.commandCh <- command{
			typ:       cmdProcessExited,
			exitError: err,
			details:   r.runtimeErrorDetails("start", err),
		}
		return
	}

	r.commandCh <- command{typ: cmdProcessStarted, cmd: cmd}

	waitErr := cmd.Wait()
	r.logger.Debug("plugin process wait returned", slog.Any("error", waitErr))

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
		details:   r.runtimeErrorDetails("exit", waitErr),
	}
}

func (r *runtime) startProcess(ctx context.Context) (*sprites.Cmd, *plugin_logger.RotatingWriter, error) {
	r.logger.Info("preparing sprite runtime")
	sprite, workdir, err := r.ensureSpriteSynced(ctx)
	if err != nil {
		return nil, nil, fault.Wrap(err, fmsg.With("failed to prepare sprite runtime"))
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

	cmd := sprite.CommandContext(ctx, r.manifest.Metadata.Command, r.manifest.Metadata.Args...)
	cmd.Dir = workdir
	cmd.Env = []string{
		"STORYDEN_RPC_URL=" + rpcURL.String(),
	}
	r.logger.Info("starting sprite command",
		slog.String("command", r.manifest.Metadata.Command),
		slog.Any("args", r.manifest.Metadata.Args),
		slog.String("workdir", workdir))

	logWriter, err := r.pluginLogger.NewWriter(r.dataPath, r.id)
	if err != nil {
		return nil, nil, fault.Wrap(err, fmsg.With("failed to create plugin log writer"))
	}

	cmd.Stdout = logWriter
	cmd.Stderr = logWriter

	if err := cmd.Start(); err != nil {
		return nil, nil, fault.Wrap(err, fmsg.With("failed to start process"))
	}
	r.logger.Info("sprite command started")

	return cmd, logWriter, nil
}

func (r *runtime) ensureSpriteSynced(ctx context.Context) (*sprites.Sprite, string, error) {
	r.logger.Info("ensuring sprite runtime is available", slog.String("sprite", r.spriteName))
	sprite, err := r.client.GetSprite(ctx, r.spriteName)
	if err != nil {
		if !isSpriteNotFoundError(err) {
			return nil, "", err
		}

		r.logger.Info("sprite not found, creating", slog.String("sprite", r.spriteName))
		sprite, err = r.client.CreateSprite(ctx, r.spriteName, nil)
		if err != nil {
			return nil, "", fault.Wrap(err, fmsg.With("failed to create sprite"))
		}

		r.logger.Info("created sprite runtime", slog.String("sprite", r.spriteName))
	}

	fsys := sprite.Filesystem()

	r.logger.Debug("clearing sprite workdir", slog.String("workdir", spriteWorkdir))
	if err := fsys.RemoveAllContext(ctx, spriteWorkdir); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, "", fault.Wrap(err, fmsg.With("failed to clear sprite working directory"))
	}

	r.logger.Info("uploading plugin archive to sprite filesystem")
	stats, err := extractSDXArchiveToSprite(ctx, fsys, r.bin, spriteWorkdir)
	if err != nil {
		return nil, "", fault.Wrap(err, fmsg.With("failed to write plugin archive to sprite filesystem"))
	}
	r.logger.Info("sprite filesystem sync complete",
		slog.Int("files", stats.Files),
		slog.Int("directories", stats.Directories),
		slog.Int64("bytes", stats.Bytes),
		slog.String("workdir", spriteWorkdir))

	return sprite, spriteWorkdir, nil
}

func (r *runtime) reportState(state plugin.ReportedState, message string, details map[string]any) {
	r.eventCh <- supervised_runtime.Event{
		State:   state,
		Message: message,
		Details: details,
	}
}

func (r *runtime) runtimeErrorDetails(stage string, err error) map[string]any {
	details := map[string]any{
		"runtime_provider": "sprites",
		"sprite_name":      r.spriteName,
		"stage":            stage,
		"error":            err.Error(),
	}

	var exitErr *sprites.ExitError
	if errors.As(err, &exitErr) {
		details["exit_code"] = exitErr.ExitCode()
	}

	return details
}

func isSpriteNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), "sprite not found")
}

func isAlreadyStoppedError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "after process finished") ||
		strings.Contains(msg, "connection closed")
}

func extractSDXArchiveToSprite(ctx context.Context, fsys sprites.FS, bin []byte, workdir string) (archiveExtractStats, error) {
	stats := archiveExtractStats{}

	zr, err := zip.NewReader(bytes.NewReader(bin), int64(len(bin)))
	if err != nil {
		return stats, fault.Wrap(err, fmsg.With("failed to open archive"))
	}

	if err := fsys.MkdirAll(workdir, 0o755); err != nil {
		return stats, fault.Wrap(err, fmsg.With("failed to create workdir"))
	}

	for _, file := range zr.File {
		targetPath, err := joinWithinSprite(workdir, file.Name)
		if err != nil {
			return stats, err
		}

		if file.FileInfo().IsDir() {
			if err := fsys.MkdirAll(targetPath, 0o755); err != nil {
				return stats, fault.Wrap(err, fmsg.With("failed to create dir"))
			}
			stats.Directories++
			continue
		}

		rc, err := file.Open()
		if err != nil {
			return stats, fault.Wrap(err, fmsg.With("failed to open file in archive"))
		}

		data, err := io.ReadAll(rc)
		closeErr := rc.Close()
		if err != nil {
			return stats, fault.Wrap(err, fmsg.With("failed to read file from archive"))
		}
		if closeErr != nil {
			return stats, fault.Wrap(closeErr, fmsg.With("failed to close archive file"))
		}

		mode := file.FileInfo().Mode() & 0o777
		if mode == 0 {
			mode = 0o644
		}

		if err := fsys.WriteFileContext(ctx, targetPath, data, mode); err != nil {
			return stats, fault.Wrap(err, fmsg.With("failed to write file "+file.Name))
		}

		stats.Files++
		stats.Bytes += int64(len(data))
	}

	return stats, nil
}

func joinWithinSprite(base, name string) (string, error) {
	cleanBase := path.Clean(base)
	clean := path.Clean(path.Join(cleanBase, name))
	if !strings.HasPrefix(clean, cleanBase+"/") && clean != cleanBase {
		return "", fault.New("archive entry escapes extraction directory")
	}
	return clean, nil
}
