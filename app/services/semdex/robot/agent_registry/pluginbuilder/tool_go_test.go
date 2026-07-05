package pluginbuilder

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	localworkspace "github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider/local"
)

func TestEnsureStorydenModuleRequirementRepairsUnresolvedZeroVersion(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	writeWorkspaceFile(t, ctx, workspace, "go.mod", `module example.com/plugin

go 1.26.4

require github.com/Southclaws/storyden v0.0.0
`)

	require.NoError(t, ensureStorydenModuleRequirement(ctx, workspace))

	read, err := workspace.ReadFile(ctx, "go.mod", -1)
	require.NoError(t, err)
	require.NotContains(t, string(read.Content), "require github.com/Southclaws/storyden")
	require.NotContains(t, string(read.Content), "github.com/Southclaws/storyden v0.0.0")
}

func TestEnsureStorydenModuleRequirementKeepsLocalReplace(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	writeWorkspaceFile(t, ctx, workspace, "go.mod", `module example.com/plugin

go 1.26.4

require github.com/Southclaws/storyden v0.0.0

replace github.com/Southclaws/storyden => /workspace/storyden
`)

	require.NoError(t, ensureStorydenModuleRequirement(ctx, workspace))

	read, err := workspace.ReadFile(ctx, "go.mod", -1)
	require.NoError(t, err)
	require.Contains(t, string(read.Content), "require github.com/Southclaws/storyden v0.0.0")
	require.Contains(t, string(read.Content), "replace github.com/Southclaws/storyden => /workspace/storyden")
}
