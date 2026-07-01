package sprites

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	spritesdk "github.com/superfly/sprites-go"

	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/account/authentication/access_key"
	robotresource "github.com/Southclaws/storyden/app/resources/robot"
	workspacecap "github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider/workspace"
	"github.com/Southclaws/storyden/internal/config"
)

const (
	defaultCommandLimit = 256_000
	defaultCommandWait  = 30 * time.Second
	spritesHTTPTimeout  = 2 * time.Minute

	spriteNamePrefix = "storyden-robot-workspace-"
	spriteWorkdir    = "/workspace"
	installMarkPath  = ".storyden/sd-installed"
	authMarkPath     = ".storyden/sd-authenticated"
	authKeyIDState   = "sd_auth_key_id"
)

type Provider struct {
	client         *spritesdk.Client
	enabled        bool
	publicAPI      url.URL
	accessKeys     *access_key.Repository
	commandTimeout time.Duration
	outputLimit    int
}

func New(cfg config.Config, accessKeys *access_key.Repository) *Provider {
	apiKey := strings.TrimSpace(cfg.SpritesAPIKey)
	p := &Provider{
		enabled:        apiKey != "",
		publicAPI:      cfg.PublicAPIAddress,
		accessKeys:     accessKeys,
		commandTimeout: defaultCommandWait,
		outputLimit:    defaultCommandLimit,
	}
	if p.enabled {
		p.client = spritesdk.New(
			apiKey,
			spritesdk.WithHTTPClient(&http.Client{
				Timeout: spritesHTTPTimeout,
			}),
			// https://github.com/superfly/sprites-go/issues/22
			spritesdk.WithDisableControl(),
		)
	}
	return p
}

func (p *Provider) Enabled() bool {
	return p.enabled
}

func (*Provider) Provider() robotresource.WorkspaceProvider {
	return robotresource.WorkspaceProviderSprites
}

func (p *Provider) Mount(ctx context.Context, instance *robotresource.WorkspaceInstance) (map[string]any, error) {
	spriteName := spriteNameForInstance(instance.ID)
	if existing, ok := instance.ProviderState["sprite_name"].(string); ok && existing != "" {
		spriteName = existing
	}
	workdir := spriteWorkdir
	if existing, ok := instance.ProviderState["workdir"].(string); ok && existing != "" {
		workdir = existing
	}

	sprite, err := p.ensureSprite(ctx, spriteName)
	if err != nil {
		return nil, err
	}
	state := cloneMap(instance.ProviderState)

	if err := sprite.Filesystem().MkdirAll(workdir, 0o755); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create sprite workspace directory"))
	}
	workspace := &Workspace{
		sprite:         sprite,
		fsys:           sprite.FilesystemAt(workdir),
		workdir:        workdir,
		commandTimeout: p.commandTimeout,
		outputLimit:    p.outputLimit,
	}
	if err := p.provisionWorkspace(ctx, instance, workspace, state); err != nil {
		return nil, err
	}

	state["sprite_name"] = spriteName
	state["workdir"] = workdir
	return state, nil
}

func (p *Provider) Open(ctx context.Context, mount robotresource.WorkspaceMount) (workspacecap.Workspace, error) {
	spriteName, workdir, err := mountState(mount)
	if err != nil {
		return nil, err
	}
	sprite, err := p.ensureSprite(ctx, spriteName)
	if err != nil {
		return nil, err
	}

	return &Workspace{
		sprite:         sprite,
		fsys:           sprite.FilesystemAt(workdir),
		workdir:        workdir,
		commandTimeout: p.commandTimeout,
		outputLimit:    p.outputLimit,
	}, nil
}

func (p *Provider) Cleanup(ctx context.Context, instance *robotresource.WorkspaceInstance) error {
	if err := p.revokeWorkspaceAccessKey(ctx, instance); err != nil {
		return err
	}

	spriteName := spriteNameForInstance(instance.ID)
	if existing, ok := instance.ProviderState["sprite_name"].(string); ok && existing != "" {
		spriteName = existing
	}
	if err := p.client.DeleteSprite(ctx, spriteName); err != nil && !isSpriteNotFoundError(err) {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to delete sprite workspace"))
	}
	return nil
}

func (p *Provider) provisionWorkspace(ctx context.Context, instance *robotresource.WorkspaceInstance, workspace *Workspace, state map[string]any) error {
	if _, err := workspace.fsys.Stat(installMarkPath); errors.Is(err, fs.ErrNotExist) {
		if result, err := workspace.Run(ctx, workspacecap.CommandSpec{
			Command: "go",
			Args:    []string{"install", "github.com/Southclaws/storyden/cmd/sd@main"},
			Timeout: 5 * time.Minute,
		}); err != nil {
			return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to install sd CLI in sprite workspace"))
		} else if !result.Success {
			return fault.Newf("failed to install sd CLI in sprite workspace: %s%s", result.Output, result.Error)
		}

		if err := workspace.fsys.WriteFileContext(ctx, installMarkPath, []byte(time.Now().UTC().Format(time.RFC3339)+"\n"), 0o644); err != nil {
			return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to write sprite workspace install marker"))
		}
	} else if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to inspect sprite workspace install marker"))
	}

	if _, err := workspace.fsys.Stat(authMarkPath); errors.Is(err, fs.ErrNotExist) {
		endpoint := strings.TrimRight(p.publicAPI.String(), "/")
		if endpoint == "" {
			return fault.New("PUBLIC_API_ADDRESS is required to authenticate sd in sprite workspace")
		}
		if p.accessKeys == nil {
			return fault.New("access key repository is required to authenticate sd in sprite workspace")
		}

		createdKey, err := p.accessKeys.Create(
			ctx,
			instance.Creator.ID,
			access_key.AccessKeyKindBot,
			fmt.Sprintf("Robot workspace CLI %s", instance.ID.String()),
			opt.NewEmpty[time.Time](),
		)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to mint sprite workspace access key"))
		}

		if result, err := workspace.Run(ctx, workspacecap.CommandSpec{
			Command: "sd",
			Args:    []string{"auth", "login", endpoint, "--access-key-stdin", "--auth-storage", "file"},
			Stdin:   createdKey.String() + "\n",
			Timeout: time.Minute,
		}); err != nil {
			_, _ = p.accessKeys.Revoke(ctx, instance.Creator.ID, createdKey.AuthID)
			return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to authenticate sd in sprite workspace"))
		} else if !result.Success {
			_, _ = p.accessKeys.Revoke(ctx, instance.Creator.ID, createdKey.AuthID)
			return fault.Newf("failed to authenticate sd in sprite workspace: %s%s", result.Output, result.Error)
		}

		if err := workspace.fsys.WriteFileContext(ctx, authMarkPath, []byte(time.Now().UTC().Format(time.RFC3339)+"\n"), 0o644); err != nil {
			_, _ = p.accessKeys.Revoke(ctx, instance.Creator.ID, createdKey.AuthID)
			return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to write sprite workspace auth marker"))
		}

		state[authKeyIDState] = createdKey.AuthID.String()
	} else if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to inspect sprite workspace auth marker"))
	}

	return nil
}

func (p *Provider) revokeWorkspaceAccessKey(ctx context.Context, instance *robotresource.WorkspaceInstance) error {
	if p.accessKeys == nil {
		return nil
	}

	raw, ok := instance.ProviderState[authKeyIDState].(string)
	if !ok || raw == "" {
		return nil
	}

	id, err := xid.FromString(raw)
	if err != nil {
		return nil
	}

	if _, err := p.accessKeys.Revoke(ctx, instance.Creator.ID, authentication.ID(id)); err != nil && ftag.Get(err) != ftag.NotFound {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to revoke sprite workspace access key"))
	}

	return nil
}

func (p *Provider) ensureSprite(ctx context.Context, name string) (*spritesdk.Sprite, error) {
	sprite, err := p.client.GetSprite(ctx, name)
	if err == nil {
		return sprite, nil
	}
	if !isSpriteNotFoundError(err) {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to get sprite workspace"))
	}
	sprite, err = p.client.CreateSprite(ctx, name, nil)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create sprite workspace"))
	}
	return sprite, nil
}

type Workspace struct {
	sprite         *spritesdk.Sprite
	fsys           spritesdk.FS
	workdir        string
	commandTimeout time.Duration
	outputLimit    int
}

func (w *Workspace) List(ctx context.Context, opts workspacecap.ListOptions) ([]workspacecap.FileInfo, error) {
	_ = ctx
	maxFiles := opts.MaxFiles
	if maxFiles <= 0 || maxFiles > 500 {
		maxFiles = 200
	}

	files := []workspacecap.FileInfo{}
	if err := fs.WalkDir(w.fsys, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if name == "." {
			return nil
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", "node_modules", ".next":
				return fs.SkipDir
			}
			return nil
		}
		if len(files) >= maxFiles {
			return fs.SkipDir
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		files = append(files, workspacecap.FileInfo{
			Path:    cleanOutputPath(name),
			Size:    info.Size(),
			Mode:    info.Mode(),
			ModTime: info.ModTime().UTC().Format(time.RFC3339),
		})
		return nil
	}); err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files, nil
}

func (w *Workspace) ReadFile(ctx context.Context, relpath string, maxBytes int) (workspacecap.ReadFileResult, error) {
	clean, err := resolvePath(relpath)
	if err != nil {
		return workspacecap.ReadFileResult{}, err
	}

	data, err := w.fsys.ReadFile(clean)
	if err != nil {
		return workspacecap.ReadFileResult{}, err
	}

	limit := maxBytes
	if limit == 0 || limit > 128_000 {
		limit = 64_000
	}

	truncated := false
	if limit > 0 && len(data) > limit {
		data = data[:limit]
		truncated = true
	}

	_ = ctx
	return workspacecap.ReadFileResult{
		Path:      clean,
		Content:   data,
		Truncated: truncated,
	}, nil
}

func (w *Workspace) WriteFile(ctx context.Context, relpath string, data []byte) (workspacecap.WriteFileResult, error) {
	clean, err := resolvePath(relpath)
	if err != nil {
		return workspacecap.WriteFileResult{}, err
	}
	if len(data) > 512_000 {
		return workspacecap.WriteFileResult{}, errors.New("content exceeds 512KB write limit")
	}
	if err := w.fsys.WriteFileContext(ctx, clean, data, 0o644); err != nil {
		return workspacecap.WriteFileResult{}, err
	}
	return workspacecap.WriteFileResult{Path: clean, Bytes: len(data)}, nil
}

func (w *Workspace) Search(ctx context.Context, query string, maxMatches int) (workspacecap.SearchResult, error) {
	if strings.TrimSpace(query) == "" {
		return workspacecap.SearchResult{}, errors.New("query is required")
	}
	if maxMatches <= 0 || maxMatches > 100 {
		maxMatches = 50
	}

	needle := strings.ToLower(query)
	matches := []workspacecap.SearchMatch{}
	if err := fs.WalkDir(w.fsys, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if name == "." {
			return nil
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", "node_modules", ".next":
				return fs.SkipDir
			}
			return nil
		}
		if len(matches) >= maxMatches {
			return fs.SkipDir
		}
		data, err := w.fsys.ReadFile(name)
		if err != nil {
			return err
		}
		if bytes.IndexByte(data, 0) >= 0 {
			return nil
		}
		lines := strings.Split(string(data), "\n")
		for i, line := range lines {
			if strings.Contains(strings.ToLower(line), needle) {
				matches = append(matches, workspacecap.SearchMatch{
					Path: cleanOutputPath(name),
					Line: i + 1,
					Text: line,
				})
				if len(matches) >= maxMatches {
					break
				}
			}
		}
		return nil
	}); err != nil {
		return workspacecap.SearchResult{}, err
	}

	_ = ctx
	return workspacecap.SearchResult{Matches: matches}, nil
}

func (w *Workspace) Run(ctx context.Context, spec workspacecap.CommandSpec) (workspacecap.CommandResult, error) {
	timeout := spec.Timeout
	if timeout == 0 {
		timeout = w.commandTimeout
	}
	cmdCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := w.sprite.CommandContext(cmdCtx, spec.Command, spec.Args...)
	cmd.Dir = w.workdir
	cmd.Env = append([]string{}, defaultCommandEnv()...)
	if spec.Stdin != "" {
		cmd.Stdin = strings.NewReader(spec.Stdin)
	}

	var out bytes.Buffer
	limited := &limitWriter{w: &out, limit: w.outputLimit}
	cmd.Stdout = limited
	cmd.Stderr = limited

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)
	if cmdCtx.Err() != nil {
		err = cmdCtx.Err()
	}

	result := workspacecap.CommandResult{
		Command:    spec.Command + " " + strings.Join(spec.Args, " "),
		Success:    err == nil,
		Output:     out.String(),
		Truncated:  limited.truncated,
		DurationMS: duration.Milliseconds(),
	}
	if err != nil {
		result.Error = err.Error()
	}

	return result, nil
}

func defaultCommandEnv() []string {
	return []string{
		"PATH=/root/go/bin:/.sprite/languages/python/pyenv/shims:/.sprite/languages/python/pyenv/bin:/.sprite/languages/deno/bin:/.sprite/languages/bun/bin:/.sprite/languages/rust/cargo/bin:/.sprite/languages/ruby/rbenv/shims:/.sprite/languages/ruby/rbenv/bin:/.sprite/languages/go/current/bin:/.sprite/languages/node/nvm/versions/node/v22.20.0/bin:/home/sprite/.local/bin:/.sprite/bin:/go/bin:/usr/local/go/bin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		"GOBIN=/root/go/bin",
		"GOPATH=/root/go",
	}
}

func mountState(mount robotresource.WorkspaceMount) (string, string, error) {
	spriteName, ok := mount.ProviderState["sprite_name"].(string)
	if !ok || spriteName == "" {
		return "", "", fmt.Errorf("sprites workspace mount has no sprite_name")
	}
	workdir, ok := mount.ProviderState["workdir"].(string)
	if !ok || workdir == "" {
		workdir = spriteWorkdir
	}
	return spriteName, workdir, nil
}

func spriteNameForInstance(id robotresource.WorkspaceInstanceID) string {
	return spriteNamePrefix + id.String()
}

func resolvePath(relpath string) (string, error) {
	trimmed := strings.TrimSpace(relpath)
	if trimmed == "" {
		return "", errors.New("path is required")
	}
	clean := path.Clean(trimmed)
	if clean == "." || clean == "" {
		return "", errors.New("path is required")
	}
	if path.IsAbs(trimmed) || strings.HasPrefix(clean, "../") || clean == ".." {
		return "", fmt.Errorf("path %q must stay inside the workspace", relpath)
	}
	return clean, nil
}

func cleanOutputPath(name string) string {
	return strings.TrimPrefix(path.Clean(name), "./")
}

func cloneMap(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func isSpriteNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), "sprite not found")
}

type limitWriter struct {
	w         io.Writer
	limit     int
	written   int
	truncated bool
}

func (l *limitWriter) Write(p []byte) (int, error) {
	original := len(p)
	if l.written >= l.limit {
		l.truncated = true
		return original, nil
	}
	remain := l.limit - l.written
	if len(p) > remain {
		l.truncated = true
		p = p[:remain]
	}
	n, err := l.w.Write(p)
	l.written += n
	if err != nil {
		return n, err
	}
	return original, nil
}
