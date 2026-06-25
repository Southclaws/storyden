package local

import (
	"context"
	"strings"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"

	robotresource "github.com/Southclaws/storyden/app/resources/robot"
	workspacecap "github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider/workspace"
	"github.com/Southclaws/storyden/internal/config"
)

func TestProviderOpenReturnsWorkspace(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	provider := New(config.Config{RobotWorkspaceDataPath: root})

	instance := &robotresource.WorkspaceInstance{
		ID:            robotresource.WorkspaceInstanceID(xid.New()),
		WorkspaceID:   robotresource.WorkspaceID(xid.New()),
		Provider:      robotresource.WorkspaceProviderLocal,
		ProviderState: map[string]any{},
		Metadata:      map[string]any{},
	}

	state, err := provider.Mount(ctx, instance)
	require.NoError(t, err)

	workspace, err := provider.Open(ctx, robotresource.WorkspaceMount{
		WorkspaceID:         instance.WorkspaceID,
		WorkspaceInstanceID: instance.ID,
		Provider:            robotresource.WorkspaceProviderLocal,
		ProviderState:       state,
	})
	require.NoError(t, err)

	written, err := workspace.WriteFile(ctx, "main.go", []byte("package main\n"))
	require.NoError(t, err)
	require.Equal(t, "main.go", written.Path)

	read, err := workspace.ReadFile(ctx, "main.go", -1)
	require.NoError(t, err)
	require.Equal(t, "package main\n", string(read.Content))
}

func TestWorkspaceRejectsEscapingPaths(t *testing.T) {
	ctx := context.Background()
	workspace, err := NewWorkspace(t.TempDir())
	require.NoError(t, err)

	_, err = workspace.ReadFile(ctx, "../secret.txt", -1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "must stay inside the workspace")

	_, err = workspace.WriteFile(ctx, "/tmp/secret.txt", []byte("nope"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "must stay inside the workspace")
}

func TestWorkspaceRunUsesWorkspaceRoot(t *testing.T) {
	ctx := context.Background()
	workspace, err := NewWorkspace(t.TempDir())
	require.NoError(t, err)

	result, err := workspace.Run(ctx, workspacecap.CommandSpec{Command: "go", Args: []string{"env", "GOOS"}})
	require.NoError(t, err)
	require.True(t, result.Success, result.Output)
	require.NotEmpty(t, strings.TrimSpace(result.Output))
}
