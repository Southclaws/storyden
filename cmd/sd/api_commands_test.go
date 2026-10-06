package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/rs/xid"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/cmd/sd/internal/cli"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
)

func executeFixture(t *testing.T, endpoint, stdin string, args ...string) (string, error) {
	t.Helper()
	store := config.NewFileStoreAt(filepath.Join(t.TempDir(), "config.yaml"))
	cfg := config.New()
	cfg.UpsertContext("fixture", config.Context{APIURL: endpoint, AuthType: config.AuthStorageFile, Auth: &config.Auth{Method: config.AuthMethodAccessKey, AccessToken: "fixture-key"}})
	cfg.SetCurrentContext("fixture")
	require.NoError(t, store.Save(cfg))
	var out, stderr bytes.Buffer
	var root *cobra.Command
	app := fx.New(fx.NopLogger, Build(), fx.Replace(store, cli.Streams{In: strings.NewReader(stdin), Out: &out, Err: &stderr}), fx.Populate(&root))
	require.NoError(t, app.Err())
	root.SetArgs(args)
	err := root.ExecuteContext(context.Background())
	return out.String(), err
}

func TestAPICommandsSendExactRequests(t *testing.T) {
	const resourceID = "d3v6n8cr5gq000000001"
	for _, tc := range []struct {
		name                                       string
		args                                       []string
		stdin, method, path, query, body, response string
	}{
		{name: "admin settings", args: []string{"admin", "settings", "update", "--file", "-"}, stdin: `{"title":"Our community","description":"","metadata":{"enabled":false,"counter":9007199254740993}}`, method: "PATCH", path: "/api/admin", body: `{"title":"Our community","description":"","metadata":{"enabled":false,"counter":9007199254740993}}`, response: `{"title":"Our community","nested":{"counter":9007199254740993}}`},
		{name: "category create", args: []string{"category", "create", "--data", `{"name":"General","description":"Welcome","colour":"#aabbcc"}`}, method: "POST", path: "/api/categories", body: `{"name":"General","description":"Welcome","colour":"#aabbcc"}`, response: `{"id":"category","name":"General"}`},
		{name: "thread reply pagination", args: []string{"thread", "page", resourceID, "--page", "2"}, method: "GET", path: "/api/threads/" + resourceID, query: "page=2", response: `{"replies":{"current_page":2,"next_page":3,"replies":[]}}`},
		{name: "trail definition", args: []string{"trail", "create", "--file", "-"}, stdin: `{"name":"Review posts","status":"paused","trigger":{"type":"event","events":["EventThreadPublished"]},"actions":[{"type":"robot_run","robot_ref":"denbot","instruction":"Review this thread"}]}`, method: "POST", path: "/api/trails", body: `{"name":"Review posts","status":"paused","trigger":{"type":"event","events":["EventThreadPublished"]},"actions":[{"type":"robot_run","robot_ref":"denbot","instruction":"Review this thread"}]}`, response: `{"id":"trail","status":"paused"}`},
		{name: "manual trail run", args: []string{"trail", "runs", "start", resourceID}, method: "POST", path: "/api/trails/" + resourceID + "/runs", response: `{"id":"run","actions":[]}`},
		{name: "cancel turn", args: []string{"robot", "sessions", "turns", "cancel", resourceID, resourceID}, method: "PUT", path: "/api/robots/sessions/" + resourceID + "/turns/" + resourceID + "/cancellation", response: `{"status":"cancelling"}`},
		{name: "delete with no body", args: []string{"role", "delete", resourceID}, method: "DELETE", path: "/api/roles/" + resourceID, response: ""},
		{name: "query escape", args: []string{"api", "request", "TagGet", "--path", "tag_name=hello /?&"}, method: "GET", path: "/api/tags/hello%20%2F%3F&", response: `{"name":"hello /?&"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				require.Equal(t, "Bearer fixture-key", r.Header.Get("Authorization"))
				require.Equal(t, tc.method, r.Method)
				require.Equal(t, tc.path, r.URL.EscapedPath())
				require.Equal(t, tc.query, r.URL.RawQuery)
				data, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.Equal(t, tc.body, string(data))
				w.Header().Set("Content-Type", "application/json")
				if tc.response == "" {
					w.WriteHeader(204)
				} else {
					fmt.Fprint(w, tc.response)
				}
			}))
			defer server.Close()
			out, err := executeFixture(t, server.URL, tc.stdin, tc.args...)
			require.NoError(t, err)
			require.Equal(t, 1, calls)
			if tc.response == "" {
				require.Equal(t, "null\n", out)
			} else {
				require.JSONEq(t, tc.response, out)
				if strings.Contains(tc.response, "9007199254740993") {
					require.Contains(t, out, "9007199254740993")
				}
			}
		})
	}
}

func TestAPIValidationBeforeRequests(t *testing.T) {
	for _, tc := range []struct {
		args []string
		want string
	}{
		{[]string{"robot", "create", "--data", `{"name":"Missing playbook"}`}, "invalid body"},
		{[]string{"trail", "update", "d3v6n8cr5gq000000001", "--data", `{"status":"paused"}`}, "invalid body"},
		{[]string{"category", "create", "--data", "{"}, "invalid JSON"},
		{[]string{"admin", "settings", "update"}, "requires a body"},
		{[]string{"admin", "settings", "update", "--data", "{}", "--file", "-"}, "only one"},
		{[]string{"thread", "page", "d3v6n8cr5gq000000001", "--page", "banana"}, "invalid request"},
		{[]string{"api", "request", "CategoryList", "--query", "typo=true"}, "unknown query"},
		{[]string{"api", "request", "CategoryList", "--timeout", "-1s"}, "non-negative"},
		{[]string{"api", "request", "TrailGet", "--path", "trail_id=.."}, "invalid path"},
		{[]string{"robot", "chat", "--session", "not-an-xid", "hello"}, "valid session XID"},
	} {
		t.Run(strings.Join(tc.args, " "), func(t *testing.T) {
			out, err := executeFixture(t, "http://127.0.0.1:1", "", tc.args...)
			require.ErrorContains(t, err, tc.want)
			require.Empty(t, out)
		})
	}
}

func TestRobotChatSubmitsUserMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)
		require.Equal(t, "/api/robots/sessions", r.URL.Path)
		var body struct {
			ID        string `json:"id"`
			SessionID string `json:"sessionId"`
			RobotID   string `json:"robotId"`
			Messages  []struct {
				ID    string                        `json:"id"`
				Role  string                        `json:"role"`
				Parts []struct{ Type, Text string } `json:"parts"`
			} `json:"messages"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		_, err := xid.FromString(body.ID)
		require.NoError(t, err)
		require.Equal(t, body.ID, body.SessionID)
		require.Equal(t, "custom-robot", body.RobotID)
		require.Len(t, body.Messages, 1)
		require.Equal(t, "user", body.Messages[0].Role)
		require.Len(t, body.Messages[0].Parts, 1)
		require.Equal(t, "text", body.Messages[0].Parts[0].Type)
		require.Equal(t, "Review the recent discussions", body.Messages[0].Parts[0].Text)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(202)
		fmt.Fprint(w, `{"session_id":"accepted","status":"queued"}`)
	}))
	defer server.Close()
	out, err := executeFixture(t, server.URL, "Review the recent discussions", "robot", "chat", "--robot", "custom-robot", "--message-file", "-")
	require.NoError(t, err)
	require.JSONEq(t, `{"session_id":"accepted","status":"queued"}`, out)
}

func TestUploadsDownloadsAndStreams(t *testing.T) {
	const id = "d3v6n8cr5gq000000001"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method + " " + r.URL.Path {
		case "POST /api/assets":
			require.Equal(t, int64(4), r.ContentLength)
			require.Equal(t, "application/octet-stream", r.Header.Get("Content-Type"))
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			require.Equal(t, []byte{0, 1, 2, 255}, body)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"id":"asset"}`)
		case "GET /api/assets/file.bin":
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write([]byte{0, 1, 2, 255})
		case "GET /api/robots/sessions/" + id + "/stream":
			require.Equal(t, "-1", r.URL.Query().Get("offset"))
			if r.URL.Query().Get("live") == "sse" {
				w.Header().Set("Content-Type", "text/event-stream")
				fmt.Fprint(w, "data: {\"type\":\"text\"}\n\n")
			} else {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `[{"type":"text","text":"Hello"}]`)
			}
		case "HEAD /api/robots/sessions/" + id + "/stream":
			w.Header().Set("Stream-Next-Offset", "42")
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL)
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	out, err := executeFixture(t, server.URL, string([]byte{0, 1, 2, 255}), "asset", "upload", "--file", "-")
	require.NoError(t, err)
	require.JSONEq(t, `{"id":"asset"}`, out)
	target := filepath.Join(t.TempDir(), "file.bin")
	out, err = executeFixture(t, server.URL, "", "asset", "download", "file.bin", "--output-file", target)
	require.NoError(t, err)
	require.Empty(t, out)
	data, err := os.ReadFile(target)
	require.NoError(t, err)
	require.Equal(t, []byte{0, 1, 2, 255}, data)
	out, err = executeFixture(t, server.URL, "", "robot", "sessions", "events", id, "--offset", "-1")
	require.NoError(t, err)
	require.JSONEq(t, `[{"type":"text","text":"Hello"}]`, out)
	out, err = executeFixture(t, server.URL, "", "robot", "sessions", "events", id, "--offset", "-1", "--live", "sse", "--output", "raw")
	require.NoError(t, err)
	require.Equal(t, "data: {\"type\":\"text\"}\n\n", out)
	out, err = executeFixture(t, server.URL, "", "robot", "sessions", "head", id)
	require.NoError(t, err)
	require.Contains(t, out, `"Stream-Next-Offset"`)
	require.Contains(t, out, `"42"`)
}

func TestAPISchemaIsOfflineAndSelfContained(t *testing.T) {
	out, err := executeFixture(t, "http://127.0.0.1:1", "", "api", "schema", "TrailCreate")
	require.NoError(t, err)
	schema, err := openapi3.NewLoader().LoadFromData([]byte(out))
	require.NoError(t, err)
	require.NotNil(t, schema.Paths.Value("/trails").Post.RequestBody.Value.Content.Get("application/json").Schema.Value.Properties["actions"].Value.Items.Value)
	require.NotContains(t, out, "AdminSettingsMutableProps")
}

func TestHTTPFailuresDoNotProduceSuccessOutputOrFollowRedirects(t *testing.T) {
	for _, status := range []int{302, 400, 401, 403, 429, 500} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				w.Header().Set("Location", "/api/redirected")
				w.WriteHeader(status)
				fmt.Fprint(w, `{"error":"rejected"}`)
			}))
			defer server.Close()
			out, err := executeFixture(t, server.URL, "", "category", "list")
			require.Error(t, err)
			require.Empty(t, out)
			require.Equal(t, 1, calls)
		})
	}
}

func TestAPIQueryValuesPreserveFalseAndRepeatedFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/admin/accounts", r.URL.Path)
		require.Equal(t, []string{"false"}, r.URL.Query()["admin"])
		require.Equal(t, []string{"false"}, r.URL.Query()["suspended"])
		require.Equal(t, []string{"alice", "bob"}, r.URL.Query()["invited_by"])
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"accounts":[]}`)
	}))
	defer server.Close()
	out, err := executeFixture(t, server.URL, "", "admin", "accounts", "list", "--admin", "false", "--suspended", "false", "--invited-by", "alice", "--invited-by", "bob")
	require.NoError(t, err)
	require.JSONEq(t, `{"accounts":[]}`, out)
}

func TestAPIRequestDeadlineCancelsRemoteRead(t *testing.T) {
	cancelled := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
		close(cancelled)
	}))
	defer server.Close()
	out, err := executeFixture(t, server.URL, "", "category", "list", "--timeout", "100ms")
	require.ErrorContains(t, err, "context deadline exceeded")
	require.Empty(t, out)
	select {
	case <-cancelled:
	case <-time.After(time.Second):
		t.Fatal("request deadline did not cancel the server connection")
	}
}
