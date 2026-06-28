package pluginbuilder

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	localworkspace "github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider/local"
)

func TestEditFileReplacesExactText(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	_, err = workspace.WriteFile(ctx, "manifest.yaml", []byte(`id: welcome-plugin
name: Welcome Plugin
description: Sends a welcome message
command: go
`))
	require.NoError(t, err)

	agent := &Agent{workspace: workspace}
	read, err := agent.ReadFile(ctx, ReadFileInput{Path: "manifest.yaml"})
	require.NoError(t, err)

	result, err := agent.EditFile(ctx, EditFileInput{
		Path:             "manifest.yaml",
		ExpectedRevision: read.Revision,
		OldText:          "description: Sends a welcome message\n",
		NewText:          "description: Sends a welcome message\nversion: 0.2.0\n",
	})
	require.NoError(t, err)
	require.True(t, result.Changed)
	require.Equal(t, "manifest.yaml", result.Path)
	require.NotEqual(t, read.Revision, result.Revision)
	require.Contains(t, result.Message, "edit applied")

	updated, err := workspace.ReadFile(ctx, "manifest.yaml", -1)
	require.NoError(t, err)
	require.Contains(t, string(updated.Content), "version: 0.2.0\n")
}

func TestEditFileRejectsStaleRevision(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	_, err = workspace.WriteFile(ctx, "main.go", []byte(`package main

func Message() string {
	return "old"
}
`))
	require.NoError(t, err)

	agent := &Agent{workspace: workspace}
	read, err := agent.ReadFile(ctx, ReadFileInput{Path: "main.go"})
	require.NoError(t, err)

	_, err = workspace.WriteFile(ctx, "main.go", []byte(`package main

func Message() string {
	return "changed elsewhere"
}
`))
	require.NoError(t, err)

	_, err = agent.EditFile(ctx, EditFileInput{
		Path:             "main.go",
		ExpectedRevision: read.Revision,
		OldText:          `return "old"`,
		NewText:          `return "new"`,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "changed since revision")
	require.Contains(t, err.Error(), "re-read before editing")
}

func TestEditFileRejectsAmbiguousReplacementWithoutLineHint(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	_, err = workspace.WriteFile(ctx, "main.go", []byte(`package main

func First() string {
	return "same"
}

func Second() string {
	return "same"
}
`))
	require.NoError(t, err)

	agent := &Agent{workspace: workspace}
	_, err = agent.EditFile(ctx, EditFileInput{
		Path:    "main.go",
		OldText: `return "same"`,
		NewText: `return "updated"`,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "appears 2 times")
	require.Contains(t, err.Error(), "expected_line")
}

func TestEditFileMissingOldTextReturnsLineContext(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	_, err = workspace.WriteFile(ctx, "main.go", []byte(`package main

func Message() string {
	return "current"
}
`))
	require.NoError(t, err)

	agent := &Agent{workspace: workspace}
	_, err = agent.EditFile(ctx, EditFileInput{
		Path:         "main.go",
		OldText:      `return "old"`,
		NewText:      `return "new"`,
		ExpectedLine: 4,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "old_text was not found")
	require.Contains(t, err.Error(), "current content near expected_line 4")
	require.Contains(t, err.Error(), `4 | 	return "current"`)
}

func TestEditFileUsesExpectedLineForRepeatedText(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	_, err = workspace.WriteFile(ctx, "main.go", []byte(`package main

func First() string {
	return "same"
}

func Second() string {
	return "same"
}
`))
	require.NoError(t, err)

	agent := &Agent{workspace: workspace}
	result, err := agent.EditFile(ctx, EditFileInput{
		Path:         "main.go",
		OldText:      `return "same"`,
		NewText:      `return "updated"`,
		ExpectedLine: 8,
	})
	require.NoError(t, err)
	require.True(t, result.Changed)

	updated, err := workspace.ReadFile(ctx, "main.go", -1)
	require.NoError(t, err)
	require.Contains(t, string(updated.Content), `func First() string {
	return "same"
}`)
	require.Contains(t, string(updated.Content), `func Second() string {
	return "updated"
}`)
}

func TestEditFileRejectsBinaryFiles(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	_, err = workspace.WriteFile(ctx, "data.bin", []byte{'a', 0, 'b'})
	require.NoError(t, err)

	agent := &Agent{workspace: workspace}
	_, err = agent.EditFile(ctx, EditFileInput{
		Path:    "data.bin",
		OldText: "a",
		NewText: "b",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "NUL byte")
}
