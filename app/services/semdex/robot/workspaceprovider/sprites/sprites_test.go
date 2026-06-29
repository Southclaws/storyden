package sprites

import (
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
