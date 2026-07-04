package local

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	robotresource "github.com/Southclaws/storyden/app/resources/robot"
	workspacecap "github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider/workspace"
	"github.com/Southclaws/storyden/internal/config"
)

const (
	defaultCommandLimit = 256_000
	defaultCommandWait  = 30 * time.Second
)

type Provider struct {
	root           string
	commandTimeout time.Duration
	outputLimit    int
}

func New(cfg config.Config) *Provider {
	return &Provider{
		root:           cfg.RobotWorkspaceDataPath,
		commandTimeout: defaultCommandWait,
		outputLimit:    defaultCommandLimit,
	}
}

func (*Provider) Provider() robotresource.WorkspaceProvider {
	return robotresource.WorkspaceProviderLocal
}

func RootPath(mount robotresource.WorkspaceMount) (string, error) {
	rootPath, ok := mount.ProviderState["root_path"].(string)
	if !ok || rootPath == "" {
		return "", fmt.Errorf("local workspace mount has no root_path")
	}

	return rootPath, nil
}

func (p *Provider) Mount(ctx context.Context, instance *robotresource.WorkspaceInstance) (map[string]any, error) {
	rootPath := path(p.root, instance.ID)
	if existing, ok := instance.ProviderState["root_path"].(string); ok && existing != "" {
		rootPath = existing
	}

	if err := os.MkdirAll(rootPath, 0o700); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	state := cloneMap(instance.ProviderState)
	state["root_path"] = rootPath
	return state, nil
}

func (p *Provider) Open(ctx context.Context, mount robotresource.WorkspaceMount) (workspacecap.Workspace, error) {
	rootPath, err := RootPath(mount)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(rootPath, 0o700); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	abs, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}

	return &Workspace{
		root:           abs,
		commandTimeout: p.commandTimeout,
		outputLimit:    p.outputLimit,
	}, nil
}

func (p *Provider) Cleanup(ctx context.Context, instance *robotresource.WorkspaceInstance) error {
	rootPath := path(p.root, instance.ID)
	if existing, ok := instance.ProviderState["root_path"].(string); ok && existing != "" {
		rootPath = existing
	}
	if rootPath == "" {
		return nil
	}
	if err := os.RemoveAll(rootPath); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	return nil
}

func path(root string, id robotresource.WorkspaceInstanceID) string {
	return filepath.Join(root, id.String())
}

func cloneMap(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

type Workspace struct {
	root           string
	commandTimeout time.Duration
	outputLimit    int
}

func NewWorkspace(root string) (*Workspace, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(abs, 0o700); err != nil {
		return nil, err
	}
	return &Workspace{
		root:           abs,
		commandTimeout: defaultCommandWait,
		outputLimit:    defaultCommandLimit,
	}, nil
}

func (w *Workspace) List(ctx context.Context, opts workspacecap.ListOptions) ([]workspacecap.FileInfo, error) {
	_ = ctx
	maxFiles := opts.MaxFiles
	if maxFiles <= 0 || maxFiles > 500 {
		maxFiles = 200
	}

	files := []workspacecap.FileInfo{}
	if err := filepath.WalkDir(w.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == w.root {
			return nil
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", "node_modules", ".next":
				return filepath.SkipDir
			}
			return nil
		}
		if len(files) >= maxFiles {
			return filepath.SkipDir
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(w.root, path)
		if err != nil {
			return err
		}
		files = append(files, workspacecap.FileInfo{
			Path:    filepath.ToSlash(rel),
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

func (w *Workspace) ReadFile(ctx context.Context, path string, maxBytes int) (workspacecap.ReadFileResult, error) {
	_ = ctx
	abs, clean, err := w.resolve(path)
	if err != nil {
		return workspacecap.ReadFileResult{}, err
	}

	data, err := os.ReadFile(abs)
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

	return workspacecap.ReadFileResult{
		Path:      clean,
		Content:   data,
		Truncated: truncated,
	}, nil
}

func (w *Workspace) WriteFile(ctx context.Context, path string, data []byte) (workspacecap.WriteFileResult, error) {
	_ = ctx
	abs, clean, err := w.resolve(path)
	if err != nil {
		return workspacecap.WriteFileResult{}, err
	}
	if len(data) > 512_000 {
		return workspacecap.WriteFileResult{}, errors.New("content exceeds 512KB write limit")
	}
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		return workspacecap.WriteFileResult{}, err
	}
	if err := os.WriteFile(abs, data, 0o644); err != nil {
		return workspacecap.WriteFileResult{}, err
	}
	return workspacecap.WriteFileResult{Path: clean, Bytes: len(data)}, nil
}

func (w *Workspace) Search(ctx context.Context, query string, maxMatches int) (workspacecap.SearchResult, error) {
	_ = ctx
	if strings.TrimSpace(query) == "" {
		return workspacecap.SearchResult{}, errors.New("query is required")
	}
	if maxMatches <= 0 || maxMatches > 100 {
		maxMatches = 50
	}

	needle := strings.ToLower(query)
	matches := []workspacecap.SearchMatch{}
	if err := filepath.WalkDir(w.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == w.root {
			return nil
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", "node_modules", ".next":
				return filepath.SkipDir
			}
			return nil
		}
		if len(matches) >= maxMatches {
			return filepath.SkipDir
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if bytes.IndexByte(data, 0) >= 0 {
			return nil
		}
		rel, err := filepath.Rel(w.root, path)
		if err != nil {
			return err
		}
		lines := strings.Split(string(data), "\n")
		for i, line := range lines {
			if strings.Contains(strings.ToLower(line), needle) {
				matches = append(matches, workspacecap.SearchMatch{
					Path: filepath.ToSlash(rel),
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

	return workspacecap.SearchResult{Matches: matches}, nil
}

func (w *Workspace) Run(ctx context.Context, spec workspacecap.CommandSpec) (workspacecap.CommandResult, error) {
	timeout := spec.Timeout
	if timeout == 0 {
		timeout = w.commandTimeout
	}
	cmdCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, spec.Command, spec.Args...)
	cmd.Dir = w.root
	cmd.Env = append(os.Environ(), spec.Env...)
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

func (w *Workspace) resolve(relpath string) (string, string, error) {
	relpath = filepath.Clean(strings.TrimSpace(relpath))
	if relpath == "." || relpath == "" {
		return "", "", errors.New("path is required")
	}
	if filepath.IsAbs(relpath) || strings.HasPrefix(relpath, ".."+string(filepath.Separator)) || relpath == ".." {
		return "", "", fmt.Errorf("path %q must stay inside the workspace", relpath)
	}
	target := filepath.Join(w.root, relpath)
	rel, err := filepath.Rel(w.root, target)
	if err != nil {
		return "", "", err
	}
	if rel == "." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." {
		return "", "", fmt.Errorf("path %q escapes workspace", relpath)
	}
	return target, filepath.ToSlash(filepath.Clean(relpath)), nil
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
