package pluginbuilder

import (
	"archive/zip"
	"bytes"
	"context"
	"testing"

	"github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider"
	workspacelocal "github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider/local"
	"github.com/stretchr/testify/require"
)

func TestImportPluginArchiveRestoresSourceAndManifestYAML(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	workspace := newImportTestWorkspace(t)
	archive := makeInstalledPluginArchive(t, map[string]string{
		"manifest.json": `{"id":"existing-plugin","author":"storyden","name":"Existing Plugin","command":"go","args":["run","."],"description":"Imported for editing","version":"1.0.0"}`,
		"go.mod":        "module storyden.local/plugins/existing-plugin\n",
		"main.go":       "package main\n",
	})

	imported, err := importPluginArchive(ctx, workspace, archive)
	require.NoError(t, err)
	require.Equal(t, "existing-plugin", imported.Manifest.Metadata.ID)
	require.Equal(t, 3, imported.Files)

	manifest, err := workspace.ReadFile(ctx, manifestYAMLFilename, -1)
	require.NoError(t, err)
	require.Contains(t, string(manifest.Content), "id: existing-plugin")
	require.Contains(t, string(manifest.Content), "name: Existing Plugin")

	source, err := workspace.ReadFile(ctx, "main.go", -1)
	require.NoError(t, err)
	require.Equal(t, "package main\n", string(source.Content))

	_, err = workspace.ReadFile(ctx, "manifest.json", -1)
	require.Error(t, err)
}

func TestImportPluginArchiveRejectsUnsafePathBeforeWriting(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	workspace := newImportTestWorkspace(t)
	archive := makeInstalledPluginArchive(t, map[string]string{
		"manifest.json": `{"id":"existing-plugin","author":"storyden","name":"Existing Plugin","command":"go","description":"Imported for editing","version":"1.0.0"}`,
		"../evil.go":    "package main\n",
	})

	_, err := importPluginArchive(ctx, workspace, archive)
	require.Error(t, err)
	require.ErrorContains(t, err, "escapes the workspace")

	files, err := workspace.List(ctx, workspaceprovider.ListOptions{MaxFiles: 10})
	require.NoError(t, err)
	require.Empty(t, files)
}

func TestRequireEmptyWorkspaceRejectsExistingFiles(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	workspace := newImportTestWorkspace(t)
	_, err := workspace.WriteFile(ctx, "main.go", []byte("package main\n"))
	require.NoError(t, err)

	err = requireEmptyWorkspace(ctx, workspace)
	require.Error(t, err)
	require.ErrorContains(t, err, "start a new chat")
}

func newImportTestWorkspace(t *testing.T) *workspacelocal.Workspace {
	t.Helper()

	workspace, err := workspacelocal.NewWorkspace(t.TempDir())
	require.NoError(t, err)
	return workspace
}

func makeInstalledPluginArchive(t *testing.T, files map[string]string) []byte {
	t.Helper()

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, content := range files {
		w, err := zw.Create(name)
		require.NoError(t, err)
		_, err = w.Write([]byte(content))
		require.NoError(t, err)
	}
	require.NoError(t, zw.Close())

	return buf.Bytes()
}
