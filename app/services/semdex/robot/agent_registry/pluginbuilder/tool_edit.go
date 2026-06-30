package pluginbuilder

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	adkagent "google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool/functiontool"
)

type EditFileInput struct {
	Path             string `json:"path" jsonschema:"Workspace-relative text file path to edit"`
	OldText          string `json:"old_text" jsonschema:"Exact current text to replace"`
	NewText          string `json:"new_text" jsonschema:"Replacement text"`
	ExpectedRevision string `json:"expected_revision,omitempty" jsonschema:"Optional file revision returned by plugin_file_read, plugin_file_search, or plugin_file_outline"`
	ExpectedLine     int    `json:"expected_line,omitempty" jsonschema:"Optional 1-based line near the intended replacement. Required when old_text appears multiple times."`
}

type EditFileResult struct {
	Path     string `json:"path"`
	Changed  bool   `json:"changed"`
	Bytes    int    `json:"bytes"`
	Revision string `json:"revision"`
	Message  string `json:"message,omitempty"`
}

func (a *Agent) addEditTools(add toolAdder) error {
	return add(functiontool.New(functiontool.Config{
		Name:        "plugin_file_edit",
		Description: editFileDescription,
	}, func(ctx adkagent.Context, args EditFileInput) (EditFileResult, error) {
		result, err := a.EditFile(ctx, args)
		if err != nil {
			return EditFileResult{}, err
		}
		return result, nil
	}))
}

const editFileDescription = `Replace exact text in one workspace-relative text file.

Use this for focused edits to existing text files, including Go files,
manifest.yaml, go.mod, README.md, and other supporting files. Use
plugin_file_write only for new files or complete rewrites.

Rules:
- Read or search the file first, then pass exact current text as old_text.
- Line numbers are advisory only; matching is based on old_text.
- Provide expected_revision from the latest read/search/outline result when available.
- If old_text appears more than once, provide expected_line near the intended occurrence.
- This tool only edits text files and rejects NUL bytes.
- This tool does not run validation. After code or manifest edits, run the validation tools in the normal workflow.

If expected_revision is stale, re-read the file before editing.`

func (a *Agent) EditFile(ctx context.Context, in EditFileInput) (EditFileResult, error) {
	if strings.TrimSpace(in.Path) == "" {
		return EditFileResult{}, errors.New("path is required")
	}
	if in.OldText == "" {
		return EditFileResult{}, errors.New("old_text is required")
	}
	if strings.Contains(in.OldText, "\x00") || strings.Contains(in.NewText, "\x00") {
		return EditFileResult{}, errors.New("edit text contains NUL byte")
	}

	snapshot, err := a.readTextSnapshot(ctx, in.Path)
	if err != nil {
		return EditFileResult{}, err
	}
	if in.ExpectedRevision != "" && in.ExpectedRevision != snapshot.Revision {
		return EditFileResult{}, fmt.Errorf("file %q changed since revision %s; re-read before editing", snapshot.Path, in.ExpectedRevision)
	}

	offset, err := selectReplacementOffset(snapshot.Content, in.OldText, in.ExpectedLine)
	if err != nil {
		if in.ExpectedLine > 0 && strings.Contains(err.Error(), "old_text was not found") {
			return EditFileResult{}, fmt.Errorf("%w; current content near expected_line %d:\n%s", err, in.ExpectedLine, editFailureContext(snapshot, in.ExpectedLine, 4))
		}
		return EditFileResult{}, err
	}

	next := snapshot.Content[:offset] + in.NewText + snapshot.Content[offset+len(in.OldText):]
	if next == snapshot.Content {
		return EditFileResult{
			Path:     snapshot.Path,
			Changed:  false,
			Bytes:    len(snapshot.Content),
			Revision: snapshot.Revision,
			Message:  "edit matched but produced no content changes",
		}, nil
	}

	workspace, err := a.Workspace(ctx)
	if err != nil {
		return EditFileResult{}, err
	}

	written, err := workspace.WriteFile(ctx, snapshot.Path, []byte(next))
	if err != nil {
		return EditFileResult{}, err
	}

	return EditFileResult{
		Path:     written.Path,
		Changed:  true,
		Bytes:    written.Bytes,
		Revision: contentRevision([]byte(next)),
		Message:  "edit applied",
	}, nil
}

func selectReplacementOffset(content, oldText string, expectedLine int) (int, error) {
	offsets := findAllOffsets(content, oldText)
	if len(offsets) == 0 {
		return 0, errors.New("old_text was not found in the current file")
	}
	if len(offsets) == 1 {
		return offsets[0], nil
	}
	if expectedLine <= 0 {
		return 0, fmt.Errorf("old_text appears %d times; provide expected_line to choose the intended occurrence", len(offsets))
	}

	bestOffset := offsets[0]
	bestDistance := math.MaxInt
	for _, offset := range offsets {
		line := lineForOffset(content, offset)
		distance := line - expectedLine
		if distance < 0 {
			distance = -distance
		}
		if distance < bestDistance {
			bestDistance = distance
			bestOffset = offset
		}
	}

	return bestOffset, nil
}

func findAllOffsets(content, oldText string) []int {
	offsets := []int{}
	searchStart := 0
	for {
		index := strings.Index(content[searchStart:], oldText)
		if index < 0 {
			return offsets
		}
		offset := searchStart + index
		offsets = append(offsets, offset)
		searchStart = offset + len(oldText)
	}
}

func lineForOffset(content string, offset int) int {
	if offset <= 0 {
		return 1
	}
	if offset > len(content) {
		offset = len(content)
	}
	return bytes.Count([]byte(content[:offset]), []byte("\n")) + 1
}

func editFailureContext(snapshot textFileSnapshot, expectedLine int, contextLines int) string {
	if snapshot.TotalLines == 0 {
		return ""
	}
	if expectedLine < 1 {
		expectedLine = 1
	}
	if expectedLine > snapshot.TotalLines {
		expectedLine = snapshot.TotalLines
	}
	start := expectedLine - contextLines
	if start < 1 {
		start = 1
	}
	end := expectedLine + contextLines
	if end > snapshot.TotalLines {
		end = snapshot.TotalLines
	}

	var b strings.Builder
	for line := start; line <= end; line++ {
		fmt.Fprintf(&b, "%d | %s", line, snapshot.Lines[line-1])
		if !strings.HasSuffix(snapshot.Lines[line-1], "\n") {
			b.WriteByte('\n')
		}
	}
	return strings.TrimRight(b.String(), "\n")
}
