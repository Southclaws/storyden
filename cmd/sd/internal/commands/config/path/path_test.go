package path

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/cmd/sd/internal/config"
)

func TestPathCommand(t *testing.T) {
	t.Run("prints config file path", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		configPath := filepath.Join(t.TempDir(), "storyden", "config.yaml")
		var stdout bytes.Buffer
		command := (*cobra.Command)(New(config.NewFileStoreAt(configPath)))
		command.SetOut(&stdout)

		r.NoError(command.Execute())

		a.Equal(configPath+"\n", stdout.String())
	})
}
