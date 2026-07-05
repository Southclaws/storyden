package sprites

import (
	"io/fs"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestResolvePathRejectsEscapes(t *testing.T) {
	for _, path := range []string{"../secret.txt", "/tmp/secret.txt", ".."} {
		_, err := resolvePath(path)
		require.Error(t, err)
		require.Contains(t, err.Error(), "must stay inside the workspace")
	}
}

func TestResolvePathCleansRelativePaths(t *testing.T) {
	path, err := resolvePath("./dir/../main.go")
	require.NoError(t, err)
	require.Equal(t, "main.go", path)
}

func TestDefaultCommandEnvIncludesSpritesToolchainPath(t *testing.T) {
	var path string
	for _, value := range defaultCommandEnv() {
		if strings.HasPrefix(value, "PATH=") {
			path = strings.TrimPrefix(value, "PATH=")
			break
		}
	}

	require.NotEmpty(t, path)
	require.Contains(t, strings.Split(path, ":"), "/.sprite/bin")
	require.Contains(t, strings.Split(path, ":"), "/.sprite/languages/go/current/bin")
	require.Contains(t, strings.Split(path, ":"), "/.sprite/languages/node/nvm/versions/node/v22.20.0/bin")
}

func TestParseFindFileList(t *testing.T) {
	files, err := parseFindFileList("go.mod\t115\t644\t1783267200.0000000000\nmain.go\t3151\t644\t1783267201.2500000000\n")
	require.NoError(t, err)
	require.Len(t, files, 2)
	require.Equal(t, "go.mod", files[0].Path)
	require.Equal(t, int64(115), files[0].Size)
	require.Equal(t, fs.FileMode(0o644), files[0].Mode)
	require.Equal(t, time.Unix(1783267200, 0).UTC().Format(time.RFC3339), files[0].ModTime)
	require.Equal(t, "main.go", files[1].Path)
	require.Equal(t, int64(3151), files[1].Size)
}

func TestParseFindFileListIgnoresKeepMarker(t *testing.T) {
	files, err := parseFindFileList(".keep\t0\t644\t1783267199.0000000000\ngo.mod\t115\t644\t1783267200.0000000000\n")
	require.NoError(t, err)
	require.Len(t, files, 1)
	require.Equal(t, "go.mod", files[0].Path)
}

func TestShellQuote(t *testing.T) {
	require.Equal(t, `'simple'`, shellQuote("simple"))
	require.Equal(t, `'dir/'"'"'file'`, shellQuote("dir/'file"))
}
