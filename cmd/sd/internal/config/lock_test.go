package config

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConcurrentProcessUpdates(t *testing.T) {
	if path := os.Getenv("SD_LOCK_TEST_CONFIG"); path != "" {
		store := NewFileStoreAt(path)
		name := os.Getenv("SD_LOCK_TEST_NAME")
		err := store.Update(context.Background(), func(cfg *Config) error {
			cfg.UpsertContext(name, Context{APIURL: "https://example.com", AuthType: AuthStorageFile, Auth: &Auth{Method: AuthMethodAccessKey, AccessToken: name}})
			return nil
		})
		require.NoError(t, err)
		return
	}
	path := filepath.Join(t.TempDir(), "config.yaml")
	store := NewFileStoreAt(path)
	var wg sync.WaitGroup
	for i := range 6 {
		wg.Go(func() {
			cmd := exec.Command(os.Args[0], "-test.run=^TestConcurrentProcessUpdates$")
			cmd.Env = append(os.Environ(), "SD_LOCK_TEST_CONFIG="+path, fmt.Sprintf("SD_LOCK_TEST_NAME=identity-%d", i))
			out, err := cmd.CombinedOutput()
			require.NoError(t, err, string(out))
		})
	}
	wg.Wait()
	cfg, err := store.Load()
	require.NoError(t, err)
	require.Len(t, cfg.Contexts, 6)
	require.Empty(t, cfg.CurrentContext)
	for i := range 6 {
		name := fmt.Sprintf("identity-%d", i)
		require.Equal(t, name, cfg.Contexts[name].Auth.AccessToken)
	}
	stat, err := os.Stat(path)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0o600), stat.Mode().Perm())
}
