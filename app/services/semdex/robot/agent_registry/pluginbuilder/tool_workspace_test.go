package pluginbuilder

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderMainGoIncludesLiveConfigurationSkeleton(t *testing.T) {
	source := renderMainGo()

	_, err := parser.ParseFile(token.NewFileSet(), "main.go", source, parser.ParseComments)
	require.NoError(t, err)

	require.Contains(t, source, "pl.OnConfigure(app.handleConfigure)")
	require.Contains(t, source, "go app.syncInitialConfig(ctx)")
	require.Contains(t, source, "Missing user-provided configuration should leave the plugin running")
	require.Contains(t, source, "thread published event ignored because plugin is not configured")
}
