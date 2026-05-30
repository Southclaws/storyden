package switcher

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
		root := &cobra.Command{Use: "sd"}
		root.SetOut(&stdout)
		root.SetErr(&bytes.Buffer{})
		root.AddCommand((*cobra.Command)(New(store)))
		root.SetArgs([]string{"switch", "prod"})

		r.NoError(root.Execute())

		loaded, err := store.Load()
		r.NoError(err)
		a.Equal("prod", loaded.CurrentContext)
		a.Contains(stdout.String(), "prod")
	})
}
