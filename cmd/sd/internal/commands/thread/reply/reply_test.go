package reply

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

func TestReplyUsesContentAndReplyTarget(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/threads/thread-id/replies", r.URL.Path)
		require.Equal(t, "Bearer bot-credential", r.Header.Get("Authorization"))
		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Equal(t, "<p>Reply</p>", body["body"])
		require.Equal(t, "post-id", body["reply_to"])
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":"reply-id"}`)
	}))
	defer server.Close()
	store := config.NewFileStoreAt(filepath.Join(t.TempDir(), "config.yaml"))
	cfg := config.New()
	cfg.UpsertContext("bot", config.Context{APIURL: server.URL, AuthType: config.AuthStorageFile, Auth: &config.Auth{Method: config.AuthMethodAccessKey, AccessToken: "bot-credential"}})
	require.NoError(t, store.Save(cfg))
	store.SelectedContext = "bot"
	cmd := (*cobra.Command)(New(store))
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"thread-id", "--content", "<p>Reply</p>", "--reply-to", "post-id", "--format", "json"})
	require.NoError(t, cmd.Execute())
	var result map[string]any
	require.NoError(t, json.Unmarshal(out.Bytes(), &result))
	require.Equal(t, "reply-id", result["id"])
}
