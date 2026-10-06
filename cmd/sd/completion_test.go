package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/carapace-sh/carapace"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/cmd/sd/internal/cli"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
)

func completionValues(t *testing.T, output string) []string {
	t.Helper()
	var suggestions []struct {
		Value string `json:"value"`
	}
	require.NoError(t, json.Unmarshal([]byte(output), &suggestions), output)
	values := make([]string, 0, len(suggestions))
	for _, item := range suggestions {
		values = append(values, strings.TrimSuffix(item.Value, " "))
	}
	return values
}

func completeFixture(t *testing.T, endpoint string, args ...string) string {
	t.Helper()
	out, err := executeFixture(t, endpoint, "", append([]string{"_carapace", "nushell", "sd"}, args...)...)
	require.NoError(t, err)
	return out
}

func TestCompletionStaticChoices(t *testing.T) {
	for _, tc := range []struct{ args, want []string }{
		{[]string{"page", "list", "--output", ""}, []string{"auto", "json", "jsonl", "plain"}},
		{[]string{"search", "hello", "--kind", "thread,n"}, []string{"thread,node"}},
		{[]string{"admin", "accounts", "list", "--admin", ""}, []string{"false", "true"}},
		{[]string{"robot", "sessions", "events", "session", "--live", ""}, []string{"sse"}},
		{[]string{"admin", "settings", "update", "--content-type", ""}, []string{"application/json"}},
		{[]string{"page", "visibility", "--from-stdin", ""}, []string{"draft", "published", "review", "unlisted"}},
		{[]string{"page", "properties", "schema", "children", "docs", "status:"}, []string{"status:boolean:", "status:number:", "status:text:", "status:timestamp:"}},
		{[]string{"page", "properties", "schema", "children", "docs", "status:text:"}, []string{"status:text:asc", "status:text:desc"}},
		{[]string{"auth", "login", "--auth-storage", ""}, []string{"auto", "file"}},
		{[]string{"api", "schema", "CategoryL"}, []string{"CategoryList"}},
		{[]string{"api", "request", "ThreadList", "--query", "visibility="}, []string{"visibility=draft", "visibility=published", "visibility=review", "visibility=unlisted"}},
		{[]string{"api", "request", "NodeGet", "--path", ""}, []string{"node_slug="}},
		{[]string{"auth", "switch", ""}, []string{"fixture"}},
		{[]string{"role", "list", "--context", ""}, []string{"fixture"}},
		{[]string{"robot", "chat", ""}, []string{}},
		{[]string{"admin", "settings", "update", "--data", ""}, []string{}},
	} {
		t.Run(strings.Join(tc.args, " "), func(t *testing.T) {
			// No static completion may depend on an instance or trigger a request.
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Errorf("unexpected request: %s %s", r.Method, r.URL)
				w.WriteHeader(500)
			}))
			defer server.Close()
			require.ElementsMatch(t, tc.want, completionValues(t, completeFixture(t, server.URL, tc.args...)))
		})
	}
}

func TestCompletionResourceLookups(t *testing.T) {
	for _, tc := range []struct {
		args                        []string
		path, query, response, want string
	}{
		{[]string{"page", "delete", ""}, "/api/nodes", "format=flat", `{"nodes":[{"id":"page-id","slug":"handbook","name":"Handbook"}]}`, "handbook"},
		{[]string{"page", "move", "handbook", "--parent", ""}, "/api/nodes", "format=flat", `{"nodes":[{"slug":"library","name":"Library"}]}`, "library"},
		{[]string{"collection", "pages", "add", "collection-id", ""}, "/api/nodes", "format=flat", `{"nodes":[{"id":"page-id","slug":"handbook","name":"Handbook"}]}`, "page-id"},
		{[]string{"page", "versions", "get", "notes / today", ""}, "/api/nodes/notes%20%2F%20today/versions", "", `{"versions":[{"id":"version-id","name":"Notes"}]}`, "version-id"},
		{[]string{"page", "assets", "remove", "handbook", ""}, "/api/nodes/handbook", "", `{"assets":[{"id":"asset-id","filename":"photo.png"}]}`, "asset-id"},
		{[]string{"thread", "get", ""}, "/api/threads", "", `{"threads":[{"id":"thread-id","slug":"topic","title":"Topic"}]}`, "thread-id"},
		{[]string{"thread", "create", "--category", ""}, "/api/categories", "", `{"categories":[{"id":"category-id","slug":"general","name":"General"}]}`, "category-id"},
		{[]string{"search", "hello", "--categories", ""}, "/api/categories", "", `{"categories":[{"id":"category-id","slug":"general","name":"General"}]}`, "category-id"},
		{[]string{"category", "delete", ""}, "/api/categories", "", `{"categories":[{"id":"category-id","slug":"general","name":"General"}]}`, "general"},
		{[]string{"profile", "follow", "al"}, "/api/profiles", "q=al", `{"profiles":[{"id":"account-id","handle":"alice","name":"Alice"}]}`, "alice"},
		{[]string{"admin", "accounts", "get", ""}, "/api/admin/accounts", "", `{"accounts":[{"id":"account-id","handle":"alice","email_addresses":[{"email_address":"PRIVATE"}]}]}`, "account-id"},
		{[]string{"role", "delete", ""}, "/api/roles", "", `{"roles":[{"id":"role-id","name":"Moderator"}]}`, "role-id"},
		{[]string{"tag", "get", ""}, "/api/tags", "", `{"tags":[{"name":"Go & tools"}]}`, `"Go & tools"`},
		{[]string{"plugin", "token", "rotate", ""}, "/api/plugins", "", `{"plugins":[{"id":"plugin-id","name":"Plugin","token":"PRIVATE"}]}`, "plugin-id"},
		{[]string{"robot", "chat", "--robot", ""}, "/api/robots", "", `{"robots":[{"id":"robot-id","name":"Assistant"}]}`, "robot-id"},
		{[]string{"robot", "chat", "--session", ""}, "/api/robots/sessions", "", `{"sessions":[{"id":"session-id","name":"Planning"}]}`, "session-id"},
		{[]string{"robot", "providers", "update", ""}, "/api/robots/providers", "", `{"providers":[{"provider":"openai","settings":{"key":"PRIVATE"}}]}`, "openai"},
		{[]string{"robot", "mcp", "delete", ""}, "/api/robots/mcp-servers", "", `{"servers":[{"id":"mcp-id","name":"Docs"}]}`, "mcp-id"},
		{[]string{"robot", "toolsets", "get", ""}, "/api/robots/toolsets", "", `{"toolsets":[{"id":"system.library","name":"Library"}]}`, "system.library"},
		{[]string{"trail", "runs", "get", "trail-id", ""}, "/api/trails/trail-id/runs", "", `{"runs":[{"id":"run-id","status":"running"}]}`, "run-id"},
		{[]string{"trail", "runs", "cancel-action", "trail-id", "run-id", ""}, "/api/trails/trail-id/runs/run-id", "", `{"actions":[{"id":"action-id","status":"running","output":"PRIVATE"}]}`, "action-id"},
		{[]string{"api", "request", "NodeVersionGet", "--path", "node_slug=handbook", "--path", "version_id="}, "/api/nodes/handbook/versions", "", `{"versions":[{"id":"version-id","name":"Handbook"}]}`, "version_id=version-id"},
	} {
		t.Run(strings.Join(tc.args, " "), func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				require.Equal(t, "GET", r.Method)
				require.Equal(t, tc.path, r.URL.EscapedPath())
				require.Equal(t, tc.query, r.URL.RawQuery)
				require.Equal(t, "Bearer fixture-key", r.Header.Get("Authorization"))
				fmt.Fprint(w, tc.response)
			}))
			defer server.Close()
			out := completeFixture(t, server.URL, tc.args...)
			require.Equal(t, 1, calls)
			require.Equal(t, []string{tc.want}, completionValues(t, out))
			require.NotContains(t, out, "PRIVATE")
		})
	}
}

func TestCompletionUsesSelectedContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "Bearer selected-key", r.Header.Get("Authorization"))
		fmt.Fprint(w, `{"roles":[{"id":"selected-role","name":"Moderator"}]}`)
	}))
	defer server.Close()
	store := config.NewFileStoreAt(filepath.Join(t.TempDir(), "config.yaml"))
	cfg := config.New()
	cfg.UpsertContext("default", config.Context{APIURL: "http://127.0.0.1:1"})
	cfg.UpsertContext("selected", config.Context{APIURL: server.URL, AuthType: config.AuthStorageFile, Auth: &config.Auth{Method: config.AuthMethodAccessKey, AccessToken: "selected-key"}})
	cfg.SetCurrentContext("default")
	require.NoError(t, store.Save(cfg))
	for _, args := range [][]string{{"--context", "selected", "role", "get", ""}, {"role", "get", "--context=selected", ""}} {
		var out, stderr bytes.Buffer
		var root *cobra.Command
		app := fx.New(fx.NopLogger, Build(), fx.Replace(store, cli.Streams{In: strings.NewReader(""), Out: &out, Err: &stderr}), fx.Populate(&root))
		require.NoError(t, app.Err())
		root.SetArgs(append([]string{"_carapace", "nushell", "sd"}, args...))
		require.NoError(t, root.ExecuteContext(context.Background()))
		require.Equal(t, []string{"selected-role"}, completionValues(t, out.String()))
		require.Empty(t, stderr.String())
		require.Empty(t, store.SelectedContext)
	}
	saved, err := store.Load()
	require.NoError(t, err)
	require.Equal(t, "default", saved.CurrentContext)
}

func TestCompletionFailuresAreQuietAndBounded(t *testing.T) {
	for _, kind := range []string{"unauthorized", "invalid-json", "timeout", "oversize", "redirect"} {
		t.Run(kind, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch kind {
				case "unauthorized":
					w.WriteHeader(403)
					fmt.Fprint(w, "PRIVATE")
				case "invalid-json":
					fmt.Fprint(w, "PRIVATE")
				case "timeout":
					<-r.Context().Done()
				case "oversize":
					fmt.Fprint(w, strings.Repeat("x", (2<<20)+1))
				case "redirect":
					http.Redirect(w, r, "/should-not-follow", http.StatusFound)
				}
			}))
			defer server.Close()
			start := time.Now()
			out := completeFixture(t, server.URL, "robot", "delete", "")
			require.Empty(t, completionValues(t, out))
			require.Less(t, time.Since(start), 4*time.Second)
		})
	}
}

func TestCompletionFilesAndShellSetup(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "payload with spaces.json")
	require.NoError(t, os.WriteFile(file, []byte("{}"), 0600))
	out := completeFixture(t, "http://127.0.0.1:1", "admin", "settings", "update", "--file", dir+string(os.PathSeparator))
	require.Contains(t, completionValues(t, out), `"`+file+`"`)
	out = completeFixture(t, "http://127.0.0.1:1", "robot", "chat", "--message-file", "-")
	require.Contains(t, completionValues(t, out), "-")
	out = completeFixture(t, "http://127.0.0.1:1", "page", "assets", "upload", "docs", "-")
	require.NotContains(t, completionValues(t, out), "-")
	for _, shell := range []string{"bash", "zsh", "fish", "nushell", "powershell"} {
		out, err := executeFixture(t, "http://127.0.0.1:1", "", "completion", shell)
		require.NoError(t, err)
		require.Contains(t, out, "_carapace")
	}
	carapace.Test(t)
}
