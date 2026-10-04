package create

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestCreateUsesSelectedIdentityAndMarkdown(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/threads", r.URL.Path)
		require.Equal(t, "Bearer bot-credential", r.Header.Get("Authorization"))
		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Equal(t, "Hello", body["title"])
		require.Equal(t, "draft", body["visibility"])
		require.Contains(t, body["body"], "<h1>Heading</h1>")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":"thread-id","title":"Hello"}`)
	}))
	defer server.Close()
	store := config.NewFileStoreAt(filepath.Join(t.TempDir(), "config.yaml"))
	cfg := config.New()
	cfg.CurrentContext = "human"
	for _, name := range []string{"human", "bot"} {
		cfg.UpsertContext(name, config.Context{APIURL: server.URL, AuthType: config.AuthStorageFile, Auth: &config.Auth{Method: config.AuthMethodAccessKey, AccessToken: name + "-credential"}})
	}
	require.NoError(t, store.Save(cfg))
	store.SelectedContext = "bot"
	cmd := (*cobra.Command)(New(store))
	var out bytes.Buffer
	cmd.SetIn(bytes.NewBufferString("# Heading"))
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--title", "Hello", "--content-file", "-", "--markdown", "--visibility", "draft", "--format", "json"})
	require.NoError(t, cmd.Execute())
	var result map[string]any
	require.NoError(t, json.Unmarshal(out.Bytes(), &result))
	require.Equal(t, "thread-id", result["id"])
	saved, err := store.Load()
	require.NoError(t, err)
	require.Equal(t, "human", saved.CurrentContext)
}
