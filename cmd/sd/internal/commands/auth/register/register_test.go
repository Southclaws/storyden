package register

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"

	"github.com/Southclaws/storyden/cmd/sd/internal/cli"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestRegistrationCommands(t *testing.T) {
	var endpoint string
	polls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/.well-known/openid-configuration" {
			fmt.Fprintf(w, `{"issuer":%q,"token_endpoint":%q,"registration_endpoint":%q}`, endpoint, endpoint+"/api/token", endpoint+"/api/register")
			return
		}
		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		if body["registration_code"] != nil {
			require.Equal(t, "registration-secret", body["registration_code"])
			require.Equal(t, true, body["cancel_registration"])
			require.Empty(t, r.Header.Get("Authorization"))
			polls++
			w.WriteHeader(204)
			return
		}
		require.Equal(t, "Bearer initial-secret", r.Header.Get("Authorization"))
		require.NotContains(t, body, "scope")
		w.WriteHeader(202)
		fmt.Fprint(w, `{"registration_code":"registration-secret","verification_code":"ABCD-EFGH","verification_uri":"https://example.com/verify","expires_in":600,"interval":30}`)
	}))
	defer server.Close()
	endpoint = server.URL
	store := config.NewFileStoreAt(filepath.Join(t.TempDir(), "config.yaml"))
	run := func(args ...string) (map[string]any, error) {
		cmd := (*cobra.Command)(cligen.NewAuthCommand(nil, nil, nil, New(store), NewCheck(store), NewWait(store), NewCancel(store), nil, nil, nil))
		var stdout, stderr bytes.Buffer
		cmd.SetOut(&stdout)
		cmd.SetErr(&stderr)
		cmd.SetIn(bytes.NewBufferString("initial-secret\n"))
		cmd.SetArgs(append([]string{"register"}, args...))
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true
		err := cmd.Execute()
		var result map[string]any
		require.NoError(t, json.Unmarshal(stdout.Bytes(), &result), stdout.String())
		require.NotContains(t, stdout.String(), "registration-secret")
		require.NotContains(t, stdout.String(), "initial-secret")
		require.NotContains(t, stdout.String(), "PRIVATE KEY")
		return result, err
	}
	result, err := run(endpoint, "--name", "bot", "--handle", "test-agent", "--auth-storage", "file", "--registration-token-stdin", "--output", "json")
	require.NoError(t, err)
	require.Equal(t, "pending", result["state"])
	saved, err := os.ReadFile(store.Path())
	require.NoError(t, err)
	require.NotContains(t, string(saved), "initial-secret")
	store.SelectedContext = "bot"
	result, err = run("check", "--output", "json")
	require.NoError(t, err)
	require.Equal(t, "pending", result["state"])
	require.Zero(t, polls)
	result, err = run("wait", "--timeout", "10ms", "--output", "json")
	require.Error(t, err)
	require.Equal(t, 2, cli.ExitCode(err))
	require.Equal(t, "pending", result["state"])
	require.Zero(t, polls)
	result, err = run("cancel", "--output", "json")
	require.NoError(t, err)
	require.Equal(t, "cancelled", result["state"])
	require.Equal(t, 1, polls)
}
