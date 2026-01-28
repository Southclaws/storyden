package plugin_runner

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/puzpuzpuz/xsync/v4"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/plugin"
	resource_plugin "github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/resources/plugin/plugin_reader"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_auth"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	lib_plugin "github.com/Southclaws/storyden/lib/plugin"
)

const ArchiveManifestFileName = "manifest.json"

type sdxRunner struct {
	logger       *slog.Logger
	sessions     *xsync.Map[resource_plugin.InstallationID, *sdxSession]
	pluginReader *plugin_reader.Reader
	bus          *pubsub.Bus
	serverURL    url.URL
	mu           sync.RWMutex
}

func New(logger *slog.Logger, pluginReader *plugin_reader.Reader, bus *pubsub.Bus) plugin_runner.Runner {
	defaultURL, _ := url.Parse("http://localhost:8000")
	return &sdxRunner{
		logger:       logger,
		sessions:     xsync.NewMap[resource_plugin.InstallationID, *sdxSession](),
		pluginReader: pluginReader,
		bus:          bus,
		serverURL:    *defaultURL,
	}
}

func (r *sdxRunner) SetServerURL(baseURL string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return
	}
	r.serverURL = *parsedURL
}

func (r *sdxRunner) Load(ctx context.Context, id plugin.InstallationID, bin []byte) (plugin_runner.Session, error) {
	manifest, err := r.Validate(ctx, bin)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	sess := &sdxSession{
		id:     id,
		logger: r.logger,
		runner: r,

		reportedState: resource_plugin.ReportedStateInactive,

		bin:      bin,
		manifest: manifest,

		bus: r.bus,

		pendingCommands: xsync.NewMap[xid.ID, pendingCommand](),
	}

	r.sessions.Store(id, sess)

	return sess, nil
}

func (r *sdxRunner) Unload(ctx context.Context, id plugin.InstallationID) error {
	sess, ok := r.sessions.Load(id)
	if !ok {
		return fault.New("plugin session not found")
	}

	if err := sess.Stop(ctx); err != nil {
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

func (r *sdxRunner) Validate(ctx context.Context, bin []byte) (*resource_plugin.Validated, error) {
	mb, err := readSDXManifest(bin)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	m, err := lib_plugin.ParseManifest(mb)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &resource_plugin.Validated{
		Binary:   bin,
		Metadata: *m,
	}, nil
}

func (r *sdxRunner) GetSession(ctx context.Context, id plugin.InstallationID) (plugin_runner.Session, error) {
	sess, ok := r.sessions.Load(id)
	if !ok {
		return nil, fault.New("plugin session not found")
	}
	return sess, nil
}

func (r *sdxRunner) GetSessions(ctx context.Context) ([]plugin_runner.Session, error) {
	out := []plugin_runner.Session{}
	r.sessions.Range(func(_ plugin.InstallationID, sess *sdxSession) bool {
		out = append(out, sess)
		return true
	})
	return out, nil
}

func (s *sdxSession) runProcess() error {
	s.stateMu.RLock()
	currentState := s.reportedState
	s.stateMu.RUnlock()

	if currentState == resource_plugin.ReportedStateActive {
		return fault.New("plugin is already running")
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.ctx = ctx
	s.cancel = cancel

	workdir, binaryPath, err := extractSDXArchive(s.bin, s.manifest.Metadata.Command)
	if err != nil {
		cancel()
		return fault.Wrap(err, fmsg.With("failed to extract plugin"))
	}

	authSecret, err := s.runner.pluginReader.GetAuthSecret(ctx, s.id)
	if err != nil {
		cancel()
		return fault.Wrap(err, fmsg.With("failed to get auth secret"))
	}

	s.runner.mu.RLock()
	baseURL := s.runner.serverURL
	s.runner.mu.RUnlock()

	baseURL.Path = "/rpc"

	rpcURL, err := plugin_auth.BuildConnectionURL(baseURL.String(), s.id, authSecret)
	if err != nil {
		cancel()
		return fault.Wrap(err, fmsg.With("failed to build connection URL"))
	}

	if rpcURL.Scheme == "http" {
		rpcURL.Scheme = "ws"
	} else if rpcURL.Scheme == "https" {
		rpcURL.Scheme = "wss"
	}

	cmd := exec.CommandContext(ctx, binaryPath)
	cmd.Dir = filepath.Dir(binaryPath)
	env := os.Environ()
	env = append(env, "STORYDEN_RPC_URL="+rpcURL.String())
	cmd.Env = env

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

	if err := cmd.Start(); err != nil {
		cancel()
		return fault.Wrap(err, fmsg.With("failed to start process"))
	}

	s.mu.Lock()
	s.cmd = cmd
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
		s.mu.Unlock()
		cancel()
	}()

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			if ctx.Err() != nil {
				return
			}
			s.logger.Info("plugin stdout", slog.String("id", s.id.String()), slog.String("line", scanner.Text()))
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			if ctx.Err() != nil {
				return
			}
			s.logger.Warn("plugin stderr", slog.String("id", s.id.String()), slog.String("line", scanner.Text()))
		}
	}()

	s.setState(resource_plugin.ReportedStateActive, "")

	waitErr := cmd.Wait()

	s.mu.Lock()
	stopping := s.stopping
	s.mu.Unlock()

	if stopping {
		s.logger.Info("plugin stopped", slog.String("id", s.id.String()))
		s.setState(resource_plugin.ReportedStateInactive, "")
		return nil
	}

	if waitErr != nil {
		s.logger.Error("plugin exited with error", slog.String("id", s.id.String()), slog.Any("error", waitErr))
		s.setState(resource_plugin.ReportedStateError, waitErr.Error())
		return fault.Wrap(waitErr, fmsg.With("process exited"))
	}

	s.logger.Warn("plugin exited unexpectedly", slog.String("id", s.id.String()))
	s.setState(resource_plugin.ReportedStateInactive, "")
	return nil
}

func readSDXManifest(bin []byte) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(bin), int64(len(bin)))
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to open archive"))
	}

	for _, file := range zr.File {
		if filepath.Base(file.Name) != ArchiveManifestFileName {
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

func extractSDXArchive(bin []byte, commandPath string) (string, string, error) {
	zr, err := zip.NewReader(bytes.NewReader(bin), int64(len(bin)))
	if err != nil {
		return "", "", fault.Wrap(err, fmsg.With("failed to open archive"))
	}

	workdir, err := os.MkdirTemp("", "storyden-sdx-")
	if err != nil {
		return "", "", fault.Wrap(err, fmsg.With("failed to create temp dir"))
	}

	commandPath = strings.TrimPrefix(commandPath, "./")

	var commandExecutablePath string

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
		if mode&0o111 == 0 && file.Name == commandPath {
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

		if file.Name == commandPath {
			commandExecutablePath = targetPath
		}
	}

	if commandExecutablePath == "" {
		return "", "", fault.Newf("command executable '%s' not found in archive", commandPath)
	}

	return workdir, commandExecutablePath, nil
}

func joinWithin(base, name string) (string, error) {
	clean := filepath.Join(base, name)
	if !strings.HasPrefix(clean, base+string(os.PathSeparator)) && clean != base {
		return "", fault.New("archive entry escapes extraction directory")
	}
	return clean, nil
}
