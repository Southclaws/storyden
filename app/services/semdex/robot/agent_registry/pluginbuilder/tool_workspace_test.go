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
	require.NotContains(t, source, "type pluginConfig struct")
	require.NotContains(t, source, "values map[string]any")
	require.NotContains(t, source, "Add manifest configuration_schema fields")
	require.Contains(t, source, "hasRuntimeConfig(raw)")
	require.Contains(t, source, "thread published event ignored because plugin is not configured")
}

func TestRenderGoModLetsGoTidyResolveStorydenVersion(t *testing.T) {
	source := renderGoMod("example-plugin")

	require.Contains(t, source, "module storyden.local/plugins/example-plugin")
	require.Contains(t, source, "go 1.26.4")
	require.NotContains(t, source, "github.com/Southclaws/storyden latest")
	require.NotContains(t, source, "github.com/Southclaws/storyden v0.0.0")
}
