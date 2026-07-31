package switcher

import (
	"bytes"
	"context"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
)

func TestSwitchCommand(t *testing.T) {
	t.Run("switches to named context", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		store := config.NewFileStoreAt(filepath.Join(t.TempDir(), "storyden", "config.yaml"))
		cfg := config.New()
		cfg.UpsertContext("local", config.Context{APIURL: "http://localhost:8000"})
		cfg.UpsertContext("prod", config.Context{APIURL: "https://forum.example.com"})
		cfg.SetCurrentContext("local")
		r.NoError(store.Save(cfg))

		var stdout bytes.Buffer
		handler := New(store)
		r.NoError(handler(context.Background(), &cobra.Command{}, cligen.IO{Out: &stdout, Err: &bytes.Buffer{}}, cligen.AuthSwitchParams{ContextName: "prod"}))

		loaded, err := store.Load()
		r.NoError(err)
		a.Equal("prod", loaded.CurrentContext)
		a.Contains(stdout.String(), "prod")
	})
}
