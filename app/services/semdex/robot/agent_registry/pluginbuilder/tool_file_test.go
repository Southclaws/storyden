package pluginbuilder

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	localworkspace "github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider/local"
)

func TestReadFileReturnsLineRangeRevisionAndTruncation(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	lines := make([]string, 0, 12)
	for i := 1; i <= 12; i++ {
		lines = append(lines, fmt.Sprintf("line %02d\n", i))
	}
	_, err = workspace.WriteFile(ctx, "notes.txt", []byte(strings.Join(lines, "")))
	require.NoError(t, err)

	agent := &Agent{workspace: workspace}
	result, err := agent.ReadFile(ctx, ReadFileInput{
		Path:      "notes.txt",
		StartLine: 4,
		MaxLines:  3,
	})
	require.NoError(t, err)
	require.Equal(t, "notes.txt", result.Path)
	require.NotEmpty(t, result.Revision)
	require.Equal(t, 4, result.StartLine)
	require.Equal(t, 6, result.EndLine)
	require.Equal(t, 12, result.TotalLines)
	require.True(t, result.Truncated)
	require.Equal(t, "line 04\nline 05\nline 06\n", result.Content)
}

func TestReadFileAroundSymbolReturnsGoFunctionWithContext(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	_, err = workspace.WriteFile(ctx, "main.go", []byte(`package main

func before() {}

func handleThreadCreated() {
	println("thread")
}

func after() {}
`))
	require.NoError(t, err)

	agent := &Agent{workspace: workspace}
	result, err := agent.ReadFile(ctx, ReadFileInput{
		Path:         "main.go",
		Symbol:       "handleThreadCreated",
		ContextLines: 1,
	})
	require.NoError(t, err)
	require.Equal(t, 4, result.StartLine)
	require.Equal(t, 8, result.EndLine)
	require.Contains(t, result.Content, "func handleThreadCreated()")
	require.Contains(t, result.Content, `println("thread")`)
}

func TestFileOutlineReportsGoRanges(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	_, err = workspace.WriteFile(ctx, "main.go", []byte(`package main

import (
	"context"
	"fmt"
)

type Plugin struct {
	name string
}

func New() *Plugin {
	return &Plugin{}
}

func (p *Plugin) Handle(ctx context.Context) error {
	fmt.Println(ctx)
	return nil
}
`))
	require.NoError(t, err)

	agent := &Agent{workspace: workspace}
	result, err := agent.FileOutline(ctx, FileOutlineInput{Path: "main.go"})
	require.NoError(t, err)
	require.Equal(t, "main.go", result.Path)
	require.Equal(t, "go", result.Language)
	require.Equal(t, "main", result.Package)
	require.Equal(t, []LineRange{{StartLine: 3, EndLine: 6}}, result.Imports)
	require.Contains(t, result.Symbols, OutlineSymbol{Kind: "type", Name: "Plugin", StartLine: 8, EndLine: 10})
	require.Contains(t, result.Symbols, OutlineSymbol{Kind: "func", Name: "New", StartLine: 12, EndLine: 14})
	require.Contains(t, result.Symbols, OutlineSymbol{Kind: "method", Name: "Handle", Receiver: "Plugin", StartLine: 16, EndLine: 19})
}

func TestSearchReturnsContextualSnippetsAndRespectsPath(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	_, err = workspace.WriteFile(ctx, "main.go", []byte(`package main

func Message() string {
	return "needle"
}
`))
	require.NoError(t, err)
	_, err = workspace.WriteFile(ctx, "README.md", []byte("needle in docs\n"))
	require.NoError(t, err)

	agent := &Agent{workspace: workspace}
	result, err := agent.Search(ctx, SearchInput{
		Query:        "needle",
		Path:         "main.go",
		ContextLines: 1,
	})
	require.NoError(t, err)
	require.Len(t, result.Matches, 1)
	match := result.Matches[0]
	require.Equal(t, "main.go", match.Path)
	require.NotEmpty(t, match.Revision)
	require.Equal(t, 4, match.Line)
	require.Equal(t, 3, match.StartLine)
	require.Equal(t, 5, match.EndLine)
	require.Equal(t, "func Message() string {\n\treturn \"needle\"\n}\n", match.Content)
}

func TestSearchSkipsBinaryFilesAndCapsResults(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	_, err = workspace.WriteFile(ctx, "one.txt", []byte("needle one\nneedle two\n"))
	require.NoError(t, err)
	_, err = workspace.WriteFile(ctx, "two.txt", []byte("needle three\n"))
	require.NoError(t, err)
	_, err = workspace.WriteFile(ctx, "data.bin", []byte{'n', 'e', 'e', 'd', 'l', 'e', 0})
	require.NoError(t, err)

	agent := &Agent{workspace: workspace}
	result, err := agent.Search(ctx, SearchInput{
		Query:      "needle",
		MaxResults: 2,
	})
	require.NoError(t, err)
	require.Len(t, result.Matches, 2)
	for _, match := range result.Matches {
		require.NotEqual(t, "data.bin", match.Path)
	}
}
