package sprites

import (
	"strings"
	"testing"

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
