package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	t.Run("missing config loads an empty config", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		store := NewFileStoreAt(filepath.Join(t.TempDir(), "storyden", "config.yaml"))

		cfg, err := store.Load()
		r.NoError(err)
		r.NotNil(cfg)

		a.Empty(cfg.CurrentContext)
		a.Empty(cfg.Contexts)
	})

	t.Run("save and load preserves contexts", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		store := NewFileStoreAt(filepath.Join(t.TempDir(), "storyden", "config.yaml"))
		cfg := New()
		cfg.UpsertContext("forum-example-com", Context{APIURL: "https://forum.example.com/api"})
		cfg.SetCurrentContext("forum-example-com")

		r.NoError(store.Save(cfg))

		loaded, err := store.Load()
		r.NoError(err)

		a.Equal("forum-example-com", loaded.CurrentContext)
		a.Equal("https://forum.example.com/api", loaded.Contexts["forum-example-com"].APIURL)
	})

	t.Run("file auth storage persists credentials in YAML", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		configPath := filepath.Join(t.TempDir(), "storyden", "config.yaml")
		store := NewFileStoreAt(configPath)
		cfg := New()
		cfg.UpsertContext("local", Context{
			APIURL:   "https://forum.example.com",
			AuthType: AuthStorageFile,
			Auth: &Auth{
				Method:       AuthMethodAccessKey,
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
			},
		})

		r.NoError(store.Save(cfg))

		data, err := os.ReadFile(configPath)
		r.NoError(err)
		a.Contains(string(data), "auth_type: file")
		a.Contains(string(data), "method: access_key")
		a.Contains(string(data), "access_token: access-token")
	})

	t.Run("credential store auth storage omits credentials from YAML and hydrates them on load", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		credentials := newMemoryCredentialStore()
		configPath := filepath.Join(t.TempDir(), "storyden", "config.yaml")
		store := NewStoreAtWithCredentialStore(configPath, credentials)
		cfg := New()
		cfg.UpsertContext("local", Context{
			APIURL:   "https://forum.example.com",
			AuthType: AuthStorageCredentialStore,
			Auth: &Auth{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresAt:    time.Now().Add(time.Hour),
			},
		})

		r.NoError(store.Save(cfg))

		data, err := os.ReadFile(configPath)
		r.NoError(err)
		a.Contains(string(data), "auth_type: credential_store")
		a.NotContains(string(data), "access_token")
		a.NotContains(string(data), "refresh_token")

		loaded, err := store.Load()
		r.NoError(err)
		auth := loaded.Contexts["local"].Auth
		r.NotNil(auth)
		a.Equal("access-token", auth.AccessToken)
		a.Equal("refresh-token", auth.RefreshToken)
	})

	t.Run("credential store auth storage deletes credentials when auth is nil", func(t *testing.T) {
		r := require.New(t)

		credentials := newMemoryCredentialStore()
		configPath := filepath.Join(t.TempDir(), "storyden", "config.yaml")
		store := NewStoreAtWithCredentialStore(configPath, credentials)
		r.NoError(credentials.SetAuth("local", Auth{AccessToken: "stale-token"}))

		cfg := New()
		cfg.UpsertContext("local", Context{
			APIURL:   "https://forum.example.com",
			AuthType: AuthStorageCredentialStore,
			Auth:     nil,
		})

		r.NoError(store.Save(cfg))

		_, ok, err := credentials.GetAuth("local")
		r.NoError(err)
		r.False(ok)
	})

	t.Run("credential store miss clears stale YAML auth", func(t *testing.T) {
		r := require.New(t)

		store := NewStoreAtWithCredentialStore(filepath.Join(t.TempDir(), "storyden", "config.yaml"), newMemoryCredentialStore())
		cfg := New()
		cfg.UpsertContext("local", Context{
			APIURL:   "https://forum.example.com",
			AuthType: AuthStorageCredentialStore,
			Auth:     &Auth{AccessToken: "stale-token"},
		})

		r.NoError(store.loadCredentials(cfg))
		r.Nil(cfg.Contexts["local"].Auth)
	})
}

type memoryCredentialStore struct {
	auth map[string]Auth
}

func newMemoryCredentialStore() *memoryCredentialStore {
	return &memoryCredentialStore{auth: map[string]Auth{}}
}

func (m *memoryCredentialStore) SetAuth(contextName string, auth Auth) error {
	m.auth[contextName] = auth
	return nil
}

func (m *memoryCredentialStore) GetAuth(contextName string) (Auth, bool, error) {
	auth, ok := m.auth[contextName]
	return auth, ok, nil
}

func (m *memoryCredentialStore) DeleteAuth(contextName string) error {
	delete(m.auth, contextName)
	return nil
}

func (m *memoryCredentialStore) Available() bool { return true }
