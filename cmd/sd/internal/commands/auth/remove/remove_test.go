package remove

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/cmd/sd/internal/config"
)

func TestRemoveCommandRemovesContextByName(t *testing.T) {
	r := require.New(t)

	credentials := newMemoryCredentialStore()
	store := config.NewStoreAtWithCredentialStore(filepath.Join(t.TempDir(), "config.yaml"), credentials)
	cfg := config.New()
	cfg.UpsertContext("localhost-8000", config.Context{
		APIURL:   "http://localhost:8000",
		AuthType: config.AuthStorageCredentialStore,
		Auth:     &config.Auth{AccessToken: "local-token"},
	})
	cfg.UpsertContext("makeroom-club", config.Context{
		APIURL:   "https://makeroom.club",
		AuthType: config.AuthStorageFile,
		Auth:     &config.Auth{AccessToken: "makeroom-token"},
	})
	cfg.SetCurrentContext("localhost-8000")
	r.NoError(store.Save(cfg))

	var out bytes.Buffer
	command := (*cobra.Command)(New(store))
	command.SetOut(&out)
	command.SetArgs([]string{"localhost-8000"})

	r.NoError(command.Execute())

	loaded, err := store.Load()
	r.NoError(err)
	r.NotContains(loaded.Contexts, "localhost-8000")
	r.Contains(loaded.Contexts, "makeroom-club")
	r.Equal("makeroom-club", loaded.CurrentContext)
	r.Contains(out.String(), "Removed context:")
	_, ok, err := credentials.GetAuth("localhost-8000")
	r.NoError(err)
	r.False(ok)
}

func TestRemoveCommandRejectsUnknownContext(t *testing.T) {
	r := require.New(t)

	store := config.NewFileStoreAt(filepath.Join(t.TempDir(), "config.yaml"))
	cfg := config.New()
	cfg.UpsertContext("localhost-8000", config.Context{APIURL: "http://localhost:8000"})
	r.NoError(store.Save(cfg))

	command := (*cobra.Command)(New(store))
	command.SetArgs([]string{"missing"})

	err := command.Execute()

	r.ErrorContains(err, "unknown context")
}

func TestRemoveContextClearsCurrentWhenLastContextRemoved(t *testing.T) {
	r := require.New(t)

	store := config.NewFileStoreAt(filepath.Join(t.TempDir(), "config.yaml"))
	cfg := config.New()
	cfg.UpsertContext("localhost-8000", config.Context{APIURL: "http://localhost:8000"})
	cfg.SetCurrentContext("localhost-8000")
	r.NoError(store.Save(cfg))

	r.NoError(removeContext(store, cfg, "localhost-8000"))

	loaded, err := store.Load()
	r.NoError(err)
	r.Empty(loaded.Contexts)
	r.Empty(loaded.CurrentContext)
}

type memoryCredentialStore struct {
	auth map[string]config.Auth
}

func newMemoryCredentialStore() *memoryCredentialStore {
	return &memoryCredentialStore{auth: map[string]config.Auth{}}
}

func (m *memoryCredentialStore) SetAuth(contextName string, auth config.Auth) error {
	m.auth[contextName] = auth
	return nil
}

func (m *memoryCredentialStore) GetAuth(contextName string) (config.Auth, bool, error) {
	auth, ok := m.auth[contextName]
	return auth, ok, nil
}

func (m *memoryCredentialStore) DeleteAuth(contextName string) error {
	delete(m.auth, contextName)
	return nil
}

func (m *memoryCredentialStore) Available() bool {
	return true
}
