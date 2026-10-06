package list

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/spf13/cobra"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateVisibilities(t *testing.T) {
	r := require.New(t)

	r.NoError(validateVisibilities(nil))
	r.NoError(validateVisibilities([]string{}))
	r.NoError(validateVisibilities([]string{"draft", "review", "published", "unlisted"}))
	r.ErrorContains(validateVisibilities([]string{"private"}), "invalid --visibility: private")
	r.ErrorContains(validateVisibilities([]string{"draft", "garbage"}), "invalid --visibility: garbage")
}

func TestUnsupportedQueriesFailBeforeAuthentication(t *testing.T) {
	for _, params := range []cligen.NodeListParams{{Page: 1, Search: "release"}, {Page: 1, SearchSet: true}, {Page: 2}} {
		err := New(nil)(context.Background(), &cobra.Command{}, cligen.IO{}, params)
		require.ErrorContains(t, err, "does not support")
	}
}

func TestAuthorAliasUsesTheSameClientSideFilter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Empty(t, r.URL.Query().Get("author"))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"nodes":[{"slug":"root","owner":{"handle":"alice"},"children":[{"slug":"child","owner":{"handle":"bob"}}]}]}`)
	}))
	defer server.Close()
	store := config.NewFileStoreAt(filepath.Join(t.TempDir(), "config.yaml"))
	cfg := config.New()
	cfg.UpsertContext("test", config.Context{APIURL: server.URL, AuthType: config.AuthStorageFile, Auth: &config.Auth{Method: config.AuthMethodAccessKey, AccessToken: "test"}})
	cfg.SetCurrentContext("test")
	require.NoError(t, store.Save(cfg))
	var alias, owner bytes.Buffer
	for _, tc := range []struct {
		author, owner string
		out           *bytes.Buffer
	}{{"alice", "", &alias}, {"", "alice", &owner}} {
		params := cligen.NodeListParams{Page: 1, Columns: "default", Output: "json", PageFormat: "tree", Author: tc.author, OwnerHandle: tc.owner}
		require.NoError(t, New(store)(context.Background(), &cobra.Command{}, cligen.IO{Out: tc.out}, params))
	}
	require.JSONEq(t, owner.String(), alias.String())
	require.Contains(t, owner.String(), `"child"`)
}
