package pluginbuilder

import (
	"context"

	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider"
)

type ListFilesInput struct {
	MaxFiles int `json:"max_files" jsonschema:"Maximum number of files to return"`
}

type FileInfo struct {
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	ModTime string `json:"mod_time"`
}

type ListFilesResult struct {
	Files []FileInfo `json:"files"`
}

type ReadFileInput struct {
	Path     string `json:"path" jsonschema:"Workspace-relative file path"`
	MaxBytes int    `json:"max_bytes" jsonschema:"Maximum bytes to read"`
}

type ReadFileResult struct {
	Path      string `json:"path"`
	Content   string `json:"content"`
	Truncated bool   `json:"truncated"`
}

type WriteFileInput struct {
	Path    string `json:"path" jsonschema:"Workspace-relative file path"`
	Content string `json:"content" jsonschema:"Complete new file content"`
}

type WriteFileResult struct {
	Path  string `json:"path"`
	Bytes int    `json:"bytes"`
}

type SearchInput struct {
	Query      string `json:"query" jsonschema:"Case-insensitive substring to search for"`
	MaxMatches int    `json:"max_matches" jsonschema:"Maximum number of matches to return"`
}

type SearchMatch struct {
	Path string `json:"path"`
	Line int    `json:"line"`
	Text string `json:"text"`
}

type SearchResult struct {
	Matches []SearchMatch `json:"matches"`
}

func (a *Agent) addFileTools(add toolAdder) error {
	if err := add(functiontool.New(functiontool.Config{
		Name:        "plugin_file_list",
		Description: "List files in a managed plugin workspace.",
	}, func(ctx adktool.Context, args ListFilesInput) (ListFilesResult, error) {
		result, err := a.ListFiles(ctx, args)
		if err != nil {
			return ListFilesResult{}, err
		}
		return result, nil
	})); err != nil {
		return err
	}

	if err := add(functiontool.New(functiontool.Config{
		Name:        "plugin_file_read",
		Description: "Read a workspace-relative text file from a managed plugin workspace.",
	}, func(ctx adktool.Context, args ReadFileInput) (ReadFileResult, error) {
		result, err := a.ReadFile(ctx, args)
		if err != nil {
			return ReadFileResult{}, err
		}
		return result, nil
	})); err != nil {
		return err
	}

	if err := add(functiontool.New(functiontool.Config{
		Name:        "plugin_file_write",
		Description: "Write complete content to a workspace-relative file. Prefer plugin_apply_patch for existing files.",
	}, func(ctx adktool.Context, args WriteFileInput) (WriteFileResult, error) {
		result, err := a.WriteFile(ctx, args)
		if err != nil {
			return WriteFileResult{}, err
		}
		return result, nil
	})); err != nil {
		return err
	}

	return add(functiontool.New(functiontool.Config{
		Name:        "plugin_file_search",
		Description: "Search text in workspace files using a case-insensitive substring.",
	}, func(ctx adktool.Context, args SearchInput) (SearchResult, error) {
		result, err := a.Search(ctx, args)
		if err != nil {
			return SearchResult{}, err
		}
		return result, nil
	}))
}

func (a *Agent) ListFiles(ctx context.Context, in ListFilesInput) (ListFilesResult, error) {
	workspace, err := a.Workspace(ctx)
	if err != nil {
		return ListFilesResult{}, err
	}

	providerFiles, err := workspace.List(ctx, workspaceprovider.ListOptions{MaxFiles: in.MaxFiles})
	if err != nil {
		return ListFilesResult{}, err
	}

	files := make([]FileInfo, 0, len(providerFiles))
	for _, file := range providerFiles {
		files = append(files, FileInfo{
			Path:    file.Path,
			Size:    file.Size,
			ModTime: file.ModTime,
		})
	}

	return ListFilesResult{Files: files}, nil
}

func (a *Agent) ReadFile(ctx context.Context, in ReadFileInput) (ReadFileResult, error) {
	workspace, err := a.Workspace(ctx)
	if err != nil {
		return ReadFileResult{}, err
	}

	result, err := workspace.ReadFile(ctx, in.Path, in.MaxBytes)
	if err != nil {
		return ReadFileResult{}, err
	}
	return ReadFileResult{
		Path:      result.Path,
		Content:   string(result.Content),
		Truncated: result.Truncated,
	}, nil
}

func (a *Agent) WriteFile(ctx context.Context, in WriteFileInput) (WriteFileResult, error) {
	workspace, err := a.Workspace(ctx)
	if err != nil {
		return WriteFileResult{}, err
	}

	result, err := workspace.WriteFile(ctx, in.Path, []byte(in.Content))
	if err != nil {
		return WriteFileResult{}, err
	}
	return WriteFileResult{Path: result.Path, Bytes: result.Bytes}, nil
}

func (a *Agent) Search(ctx context.Context, in SearchInput) (SearchResult, error) {
	workspace, err := a.Workspace(ctx)
	if err != nil {
		return SearchResult{}, err
	}

	result, err := workspace.Search(ctx, in.Query, in.MaxMatches)
	if err != nil {
		return SearchResult{}, err
	}

	matches := make([]SearchMatch, 0, len(result.Matches))
	for _, match := range result.Matches {
		matches = append(matches, SearchMatch{
			Path: match.Path,
			Line: match.Line,
			Text: match.Text,
		})
	}

	return SearchResult{Matches: matches}, nil
}
