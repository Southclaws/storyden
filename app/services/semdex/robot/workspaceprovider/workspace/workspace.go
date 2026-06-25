package workspace

import (
	"context"
	"io/fs"
	"time"
)

type Workspace interface {
	List(ctx context.Context, opts ListOptions) ([]FileInfo, error)
	ReadFile(ctx context.Context, path string, maxBytes int) (ReadFileResult, error)
	WriteFile(ctx context.Context, path string, data []byte) (WriteFileResult, error)
	Search(ctx context.Context, query string, maxMatches int) (SearchResult, error)
	Run(ctx context.Context, spec CommandSpec) (CommandResult, error)
}

type ListOptions struct {
	MaxFiles int
}

type FileInfo struct {
	Path    string
	Size    int64
	Mode    fs.FileMode
	ModTime string
}

type ReadFileResult struct {
	Path      string
	Content   []byte
	Truncated bool
}

type WriteFileResult struct {
	Path  string
	Bytes int
}

type SearchMatch struct {
	Path string
	Line int
	Text string
}

type SearchResult struct {
	Matches []SearchMatch
}

type CommandSpec struct {
	Command string
	Args    []string
	Stdin   string
	Timeout time.Duration
}

type CommandResult struct {
	Command    string
	Success    bool
	Output     string
	Error      string
	Truncated  bool
	DurationMS int64
}
