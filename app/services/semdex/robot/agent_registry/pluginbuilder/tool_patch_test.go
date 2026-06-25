package pluginbuilder

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	localworkspace "github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider/local"
)

func TestApplyPatchAppliesGopatchToGoFile(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	_, err = workspace.WriteFile(ctx, "go.mod", []byte("module testplugin\n\ngo 1.24\n"))
	require.NoError(t, err)

	_, err = workspace.WriteFile(ctx, "main.go", []byte(`package main

func Message() string {
	return "old"
}
`))
	require.NoError(t, err)

	agent := &Agent{workspace: workspace}
	result, err := agent.ApplyPatch(ctx, ApplyPatchInput{
		Path: "main.go",
		Patch: `@@
@@
-return "old"
+return "new"
`,
	})
	require.NoError(t, err)
	require.True(t, result.Changed)
	require.Equal(t, "main.go", result.Path)
	require.Contains(t, result.Message, "patch applied")
	require.NotNil(t, result.Vet)
	require.True(t, result.Vet.Success, result.Vet.Output+result.Vet.Error)

	read, err := workspace.ReadFile(ctx, "main.go", -1)
	require.NoError(t, err)
	require.Contains(t, string(read.Content), `return "new"`)
}

func TestApplyPatchReportsNoMatch(t *testing.T) {
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
	result, err := agent.ApplyPatch(ctx, ApplyPatchInput{
		Path: "main.go",
		Patch: `@@
@@
-return "missing"
+return "new"
`,
	})
	require.NoError(t, err)
	require.False(t, result.Changed)
	require.Contains(t, result.Message, "did not match")

	read, err := workspace.ReadFile(ctx, "main.go", -1)
	require.NoError(t, err)
	require.Contains(t, string(read.Content), `return "old"`)
}

func TestApplyPatchRejectsNonGoFiles(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	agent := &Agent{workspace: workspace}
	_, err = agent.ApplyPatch(ctx, ApplyPatchInput{
		Path:  "manifest.yaml",
		Patch: "@@\n@@\n-a\n+b\n",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "only supports Go source files")
}
