package credentials

import (
	"bytes"
	"path/filepath"
	"testing"
	"time"

	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"

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
	headers := credentialCommand(store, "headers")
	require.ErrorContains(t, headers.Execute(), "explicit --context")
	token := credentialCommand(store, "token")
	require.ErrorContains(t, token.Execute(), "does not use OAuth")
	store.SelectedContext = "bot"
	status := credentialCommand(store, "status")
	var out bytes.Buffer
	status.SetOut(&out)
	status.SetArgs([]string{"status", "--output", "json"})
	require.NoError(t, status.Execute())
	require.Contains(t, out.String(), `"context": "bot"`)
	for _, secret := range []string{"private-secret", "bot-secret", "poll-secret", "human-secret"} {
		require.NotContains(t, out.String(), secret)
	}
	store.SelectedContext = "human"
	headers = credentialCommand(store, "headers")
	out.Reset()
	headers.SetOut(&out)
	require.NoError(t, headers.Execute())
	require.JSONEq(t, `{"Authorization":"Bearer human-secret"}`, out.String())
}

func credentialCommand(store *config.Store, name string) *cobra.Command {
	cmd := (*cobra.Command)(cligen.NewAuthCommand(nil, nil, nil, nil, nil, nil, nil, NewToken(store), NewHeaders(store), NewStatus(store)))
	cmd.SetArgs([]string{name})
	return cmd
}
