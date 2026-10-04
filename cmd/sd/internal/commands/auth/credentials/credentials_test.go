package credentials

import (
	"bytes"
	"path/filepath"
	"testing"
	"time"

	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestCredentialOutput(t *testing.T) {
	store := config.NewFileStoreAt(filepath.Join(t.TempDir(), "config.yaml"))
	cfg := config.New()
	cfg.CurrentContext = "human"
	cfg.UpsertContext("human", config.Context{APIURL: "https://example.com", AuthType: config.AuthStorageFile, Auth: &config.Auth{Method: config.AuthMethodAccessKey, AccessToken: "human-secret"}})
	cfg.UpsertContext("bot", config.Context{APIURL: "https://example.com", AuthType: config.AuthStorageFile, Auth: &config.Auth{Method: config.AuthMethodOAuthClient, PrivateKey: "private-secret", AccessToken: "bot-secret", ExpiresAt: time.Now().Add(time.Hour), Registration: &config.Registration{State: "pending", Handle: "test-bot", Code: "poll-secret"}}})
	require.NoError(t, store.Save(cfg))
	headers := (*cobra.Command)(NewHeaders(store))
	require.ErrorContains(t, headers.Execute(), "explicit --context")
	token := (*cobra.Command)(NewToken(store))
	require.ErrorContains(t, token.Execute(), "does not use OAuth")
	store.SelectedContext = "bot"
	status := (*cobra.Command)(NewStatus(store))
	var out bytes.Buffer
	status.SetOut(&out)
	status.SetArgs([]string{"--format", "json"})
	require.NoError(t, status.Execute())
	require.Contains(t, out.String(), `"context": "bot"`)
	for _, secret := range []string{"private-secret", "bot-secret", "poll-secret", "human-secret"} {
		require.NotContains(t, out.String(), secret)
	}
	store.SelectedContext = "human"
	headers = (*cobra.Command)(NewHeaders(store))
	out.Reset()
	headers.SetOut(&out)
	require.NoError(t, headers.Execute())
	require.JSONEq(t, `{"Authorization":"Bearer human-secret"}`, out.String())
}
