package pluginbuilder

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	adkagent "google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool/functiontool"

	"github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider"
)

const (
	defaultReadLines      = 160
	maxReadLines          = 400
	defaultAroundContext  = 40
	defaultSymbolContext  = 20
	defaultSearchContext  = 2
	maxSearchContext      = 20
	defaultSearchResults  = 50
	maxSearchResults      = 100
	maxTextToolTargetSize = 512_000
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
	Path         string `json:"path" jsonschema:"Workspace-relative text file path"`
	StartLine    int    `json:"start_line,omitempty" jsonschema:"First line to read, 1-based. Defaults to 1."`
	MaxLines     int    `json:"max_lines,omitempty" jsonschema:"Maximum lines to return. Defaults to 160 and is capped at 400."`
	AroundLine   int    `json:"around_line,omitempty" jsonschema:"Read around this 1-based line number instead of start_line."`
	Symbol       string `json:"symbol,omitempty" jsonschema:"Go symbol name to read with surrounding context."`
	ContextLines int    `json:"context_lines,omitempty" jsonschema:"Context lines around around_line or symbol."`
}

type ReadFileResult struct {
	Path       string `json:"path"`
	Revision   string `json:"revision"`
	StartLine  int    `json:"start_line"`
	EndLine    int    `json:"end_line"`
	TotalLines int    `json:"total_lines"`
	Content    string `json:"content"`
	Truncated  bool   `json:"truncated"`
}

type WriteFileInput struct {
	Path    string `json:"path" jsonschema:"Workspace-relative file path"`
	Content string `json:"content" jsonschema:"Complete new file content"`
}

type WriteFileResult struct {
	Path     string `json:"path"`
	Bytes    int    `json:"bytes"`
	Revision string `json:"revision"`
}

type SearchInput struct {
	Query        string `json:"query" jsonschema:"Case-insensitive substring to search for"`
	Path         string `json:"path,omitempty" jsonschema:"Optional workspace-relative file path to search"`
	MaxResults   int    `json:"max_results,omitempty" jsonschema:"Maximum number of snippet results to return. Defaults to 50 and is capped at 100."`
	ContextLines int    `json:"context_lines,omitempty" jsonschema:"Snippet context lines around each match. Defaults to 2 and is capped at 20."`
}

type SearchMatch struct {
	Path      string `json:"path"`
	Revision  string `json:"revision"`
	Line      int    `json:"line"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
	Content   string `json:"content"`
}

type SearchResult struct {
	Matches []SearchMatch `json:"matches"`
}

type FileOutlineInput struct {
	Path string `json:"path" jsonschema:"Workspace-relative Go source file path"`
}

type FileOutlineResult struct {
	Path       string          `json:"path"`
	Revision   string          `json:"revision"`
	Language   string          `json:"language"`
	TotalLines int             `json:"total_lines"`
	Package    string          `json:"package,omitempty"`
	Imports    []LineRange     `json:"imports,omitempty"`
	Symbols    []OutlineSymbol `json:"symbols,omitempty"`
}

type LineRange struct {
	StartLine int `json:"start_line"`
	EndLine   int `json:"end_line"`
}

type OutlineSymbol struct {
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Receiver  string `json:"receiver,omitempty"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
}

type textFileSnapshot struct {
	Path       string
	Content    string
	Revision   string
	Lines      []string
	TotalLines int
}

func (a *Agent) addFileTools(add toolAdder) error {
	if err := add(functiontool.New(functiontool.Config{
		Name:        "plugin_file_list",
		Description: "List files in a managed plugin workspace.",
	}, func(ctx adkagent.Context, args ListFilesInput) (ListFilesResult, error) {
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
		Description: "Read a line range from a workspace-relative text file. Use start_line/max_lines, around_line/context_lines, or symbol/context_lines for Go files.",
	}, func(ctx adkagent.Context, args ReadFileInput) (ReadFileResult, error) {
		result, err := a.ReadFile(ctx, args)
		if err != nil {
			return ReadFileResult{}, err
		}
		return result, nil
	})); err != nil {
		return err
	}

	if err := add(functiontool.New(functiontool.Config{
		Name:        "plugin_file_outline",
		Description: "Return a compact outline of a Go source file with import, type, function, and method line ranges.",
	}, func(ctx adkagent.Context, args FileOutlineInput) (FileOutlineResult, error) {
		result, err := a.FileOutline(ctx, args)
		if err != nil {
			return FileOutlineResult{}, err
		}
		return result, nil
	})); err != nil {
		return err
	}

	if err := add(functiontool.New(functiontool.Config{
		Name:        "plugin_file_write",
		Description: "Write complete content to a workspace-relative file. Prefer plugin_file_edit for focused changes to existing files.",
	}, func(ctx adkagent.Context, args WriteFileInput) (WriteFileResult, error) {
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
		Description: "Search text in workspace files and return contextual snippets with line numbers and file revisions.",
	}, func(ctx adkagent.Context, args SearchInput) (SearchResult, error) {
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
	snapshot, err := a.readTextSnapshot(ctx, in.Path)
	if err != nil {
		return ReadFileResult{}, err
	}

	start, maxLines, err := readWindow(snapshot, in)
	if err != nil {
		return ReadFileResult{}, err
	}

	return snapshot.readRange(start, maxLines), nil
}

func (a *Agent) FileOutline(ctx context.Context, in FileOutlineInput) (FileOutlineResult, error) {
	snapshot, err := a.readTextSnapshot(ctx, in.Path)
	if err != nil {
		return FileOutlineResult{}, err
	}
	if !strings.HasSuffix(snapshot.Path, ".go") {
		return FileOutlineResult{}, fmt.Errorf("plugin_file_outline currently supports Go source files only")
	}

	outline, err := parseGoOutline(snapshot.Path, snapshot.Content)
	if err != nil {
		return FileOutlineResult{}, err
	}

	return FileOutlineResult{
		Path:       snapshot.Path,
		Revision:   snapshot.Revision,
		Language:   "go",
		TotalLines: snapshot.TotalLines,
		Package:    outline.Package,
		Imports:    outline.Imports,
		Symbols:    outline.Symbols,
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
	return WriteFileResult{
		Path:     result.Path,
		Bytes:    result.Bytes,
		Revision: contentRevision([]byte(in.Content)),
	}, nil
}

func (a *Agent) Search(ctx context.Context, in SearchInput) (SearchResult, error) {
	if strings.TrimSpace(in.Query) == "" {
		return SearchResult{}, errors.New("query is required")
	}

	if strings.TrimSpace(in.Path) != "" {
		match, err := a.searchFile(ctx, in.Path, in)
		if err != nil {
			return SearchResult{}, err
		}
		return SearchResult{Matches: match}, nil
	}

	workspace, err := a.Workspace(ctx)
	if err != nil {
		return SearchResult{}, err
	}

	files, err := workspace.List(ctx, workspaceprovider.ListOptions{MaxFiles: 1000})
	if err != nil {
		return SearchResult{}, err
	}

	matches := []SearchMatch{}
	limit := normaliseSearchLimit(in.MaxResults)
	for _, file := range files {
		fileMatches, err := a.searchFile(ctx, file.Path, in)
		if err != nil {
			continue
		}
		for _, match := range fileMatches {
			matches = append(matches, match)
			if len(matches) >= limit {
				return SearchResult{Matches: matches}, nil
			}
		}
	}

	return SearchResult{Matches: matches}, nil
}

func (a *Agent) searchFile(ctx context.Context, path string, in SearchInput) ([]SearchMatch, error) {
	snapshot, err := a.readTextSnapshot(ctx, path)
	if err != nil {
		return nil, err
	}

	needle := strings.ToLower(in.Query)
	contextLines := normaliseSearchContext(in.ContextLines)
	limit := normaliseSearchLimit(in.MaxResults)

	matches := []SearchMatch{}
	for i, line := range snapshot.Lines {
		if !strings.Contains(strings.ToLower(line), needle) {
			continue
		}
		lineNumber := i + 1
		start := lineNumber - contextLines
		if start < 1 {
			start = 1
		}
		end := lineNumber + contextLines
		if end > snapshot.TotalLines {
			end = snapshot.TotalLines
		}
		matches = append(matches, SearchMatch{
			Path:      snapshot.Path,
			Revision:  snapshot.Revision,
			Line:      lineNumber,
			StartLine: start,
			EndLine:   end,
			Content:   snapshot.contentRange(start, end),
		})
		if len(matches) >= limit {
			break
		}
	}

	return matches, nil
}

func (a *Agent) readTextSnapshot(ctx context.Context, path string) (textFileSnapshot, error) {
	workspace, err := a.Workspace(ctx)
	if err != nil {
		return textFileSnapshot{}, err
	}

	result, err := workspace.ReadFile(ctx, path, -1)
	if err != nil {
		return textFileSnapshot{}, err
	}
	if len(result.Content) > maxTextToolTargetSize {
		return textFileSnapshot{}, fmt.Errorf("file %q exceeds %d byte text tool limit", result.Path, maxTextToolTargetSize)
	}
	if bytes.Contains(result.Content, []byte{0}) {
		return textFileSnapshot{}, fmt.Errorf("file %q contains NUL byte and cannot be read as text", result.Path)
	}

	lines := splitLinesPreserve(string(result.Content))
	return textFileSnapshot{
		Path:       result.Path,
		Content:    string(result.Content),
		Revision:   contentRevision(result.Content),
		Lines:      lines,
		TotalLines: len(lines),
	}, nil
}

func readWindow(snapshot textFileSnapshot, in ReadFileInput) (int, int, error) {
	maxLines := normaliseReadLimit(in.MaxLines)
	contextLines := in.ContextLines
	if contextLines < 0 {
		return 0, 0, errors.New("context_lines cannot be negative")
	}

	if strings.TrimSpace(in.Symbol) != "" {
		if !strings.HasSuffix(snapshot.Path, ".go") {
			return 0, 0, errors.New("symbol reads are currently supported for Go source files only")
		}
		outline, err := parseGoOutline(snapshot.Path, snapshot.Content)
		if err != nil {
			return 0, 0, err
		}
		symbol, ok := outline.findSymbol(in.Symbol)
		if !ok {
			return 0, 0, fmt.Errorf("symbol %q not found in %s", in.Symbol, snapshot.Path)
		}
		if contextLines == 0 {
			contextLines = defaultSymbolContext
		}
		symbolLines := symbol.EndLine - symbol.StartLine + 1
		if symbolLines >= maxReadLines {
			return symbol.StartLine, maxReadLines, nil
		}
		availableContext := maxReadLines - symbolLines
		if contextLines*2 > availableContext {
			contextLines = availableContext / 2
		}
		start := symbol.StartLine - contextLines
		if start < 1 {
			start = 1
		}
		end := symbol.EndLine + contextLines
		if end > snapshot.TotalLines {
			end = snapshot.TotalLines
		}
		maxLines = end - start + 1
		return start, maxLines, nil
	}

	if in.AroundLine > 0 {
		if contextLines == 0 {
			contextLines = defaultAroundContext
		}
		maxLines = normaliseReadLimit(in.MaxLines)
		if in.MaxLines <= 0 {
			maxLines = contextLines*2 + 1
			if maxLines > maxReadLines {
				maxLines = maxReadLines
			}
		}
		start := in.AroundLine - maxLines/2
		if start < 1 {
			start = 1
		}
		return start, maxLines, nil
	}

	start := in.StartLine
	if start <= 0 {
		start = 1
	}
	return start, maxLines, nil
}

func normaliseReadLimit(maxLines int) int {
	if maxLines <= 0 {
		return defaultReadLines
	}
	if maxLines > maxReadLines {
		return maxReadLines
	}
	return maxLines
}

func normaliseSearchLimit(maxResults int) int {
	if maxResults <= 0 {
		return defaultSearchResults
	}
	if maxResults > maxSearchResults {
		return maxSearchResults
	}
	return maxResults
}

func normaliseSearchContext(contextLines int) int {
	if contextLines < 0 {
		return 0
	}
	if contextLines == 0 {
		return defaultSearchContext
	}
	if contextLines > maxSearchContext {
		return maxSearchContext
	}
	return contextLines
}

func (s textFileSnapshot) readRange(start, maxLines int) ReadFileResult {
	if start < 1 {
		start = 1
	}
	if maxLines <= 0 {
		maxLines = defaultReadLines
	}
	if s.TotalLines == 0 {
		return ReadFileResult{
			Path:       s.Path,
			Revision:   s.Revision,
			StartLine:  1,
			EndLine:    0,
			TotalLines: 0,
			Content:    "",
			Truncated:  false,
		}
	}
	if start > s.TotalLines {
		start = s.TotalLines
	}

	end := start + maxLines - 1
	truncated := end < s.TotalLines
	if end > s.TotalLines {
		end = s.TotalLines
	}

	return ReadFileResult{
		Path:       s.Path,
		Revision:   s.Revision,
		StartLine:  start,
		EndLine:    end,
		TotalLines: s.TotalLines,
		Content:    s.contentRange(start, end),
		Truncated:  truncated,
	}
}

func (s textFileSnapshot) contentRange(start, end int) string {
	if start < 1 {
		start = 1
	}
	if end > s.TotalLines {
		end = s.TotalLines
	}
	if s.TotalLines == 0 || start > end {
		return ""
	}
	return strings.Join(s.Lines[start-1:end], "")
}

type goFileOutline struct {
	Package string
	Imports []LineRange
	Symbols []OutlineSymbol
}

func (o goFileOutline) findSymbol(name string) (OutlineSymbol, bool) {
	name = strings.TrimSpace(name)
	for _, symbol := range o.Symbols {
		if symbol.Name == name || symbol.Receiver+"."+symbol.Name == name {
			return symbol, true
		}
	}
	return OutlineSymbol{}, false
}

func parseGoOutline(path, source string) (goFileOutline, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, source, parser.ParseComments)
	if err != nil {
		return goFileOutline{}, err
	}

	out := goFileOutline{Package: file.Name.Name}

	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			symbol := OutlineSymbol{
				Kind:      "func",
				Name:      d.Name.Name,
				StartLine: fset.Position(d.Pos()).Line,
				EndLine:   fset.Position(d.End()).Line,
			}
			if d.Recv != nil && len(d.Recv.List) > 0 {
				symbol.Kind = "method"
				symbol.Receiver = receiverName(d.Recv.List[0].Type)
			}
			out.Symbols = append(out.Symbols, symbol)
		case *ast.GenDecl:
			switch d.Tok {
			case token.IMPORT:
				out.Imports = append(out.Imports, LineRange{
					StartLine: fset.Position(d.Pos()).Line,
					EndLine:   fset.Position(d.End()).Line,
				})
			case token.TYPE:
				for _, spec := range d.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					out.Symbols = append(out.Symbols, OutlineSymbol{
						Kind:      "type",
						Name:      ts.Name.Name,
						StartLine: fset.Position(ts.Pos()).Line,
						EndLine:   fset.Position(ts.End()).Line,
					})
				}
			}
		}
	}

	return out, nil
}

func receiverName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return receiverName(t.X)
	case *ast.IndexExpr:
		return receiverName(t.X)
	case *ast.IndexListExpr:
		return receiverName(t.X)
	default:
		return strings.TrimSpace(fmt.Sprint(expr))
	}
}

func contentRevision(content []byte) string {
	sum := sha256.Sum256(content)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func splitLinesPreserve(text string) []string {
	if text == "" {
		return nil
	}
	lines := strings.SplitAfter(text, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}
