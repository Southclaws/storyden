package plugin_logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/internal/config"
)

func TestWriterRotation_DropsOldEntriesAtMaxBackups(t *testing.T) {
	tmpDir := t.TempDir()
	pluginID := plugin.InstallationID(xid.New())

	w := newWriter(config.Config{PluginDataPath: tmpDir})
	writer, err := w.NewWriter(tmpDir, pluginID)
	require.NoError(t, err)

	for i := 1; i <= 5; i++ {
		_, err := writer.Write([]byte(fmt.Sprintf("line %d\n", i)))
		require.NoError(t, err)
		require.NoError(t, writer.Rotator.Rotate())
		time.Sleep(5 * time.Millisecond)
	}

	logDir := getPluginLogDirectory(tmpDir, pluginID)
	entries, err := os.ReadDir(logDir)
	require.NoError(t, err)

	backupCount := 0
	combined := strings.Builder{}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".log" {
			continue
		}
		if entry.Name() != "output.log" {
			backupCount++
		}

		content, err := os.ReadFile(filepath.Join(logDir, entry.Name()))
		require.NoError(t, err)
		combined.Write(content)
	}

	// Existing behaviour: MaxBackups is 3 so older history is discarded.
	assert.LessOrEqual(t, backupCount, 3)
	assert.NotContains(t, combined.String(), "line 1")
}
