package tools

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/lib/mcp"
)

func TestRegistryReportsMissingTools(t *testing.T) {
	registry := NewRegistry(slog.New(slog.NewTextHandler(io.Discard, nil)))
	registry.Register(testTool("content_search", "content_search"))

	found, missing := registry.GetToolsWithMissing(context.Background(), "content_search", "mcp:search:gone", "content_search")

	require.Len(t, found, 1)
	assert.Equal(t, "content_search", found[0].Name())
	assert.Equal(t, []string{"mcp:search:gone"}, missing)
}

func TestRegistryResolvesAliases(t *testing.T) {
	registry := NewRegistry(slog.New(slog.NewTextHandler(io.Discard, nil)))
	registry.Register(testTool("library_search_pages", "library_search_pages"))
	registry.RegisterAlias("library_search", "library_search_pages")

	found, missing := registry.GetToolsWithMissing(context.Background(), "library_search")

	require.Empty(t, missing)
	require.Len(t, found, 1)
	assert.Equal(t, "library_search_pages", found[0].Name())
	assert.True(t, registry.HasTool("library_search"))
	assert.Contains(t, registry.AllToolIDs(context.Background()), "library_search")
}

func TestRegistryCatalogueUsesCallableNameForDynamicTools(t *testing.T) {
	registry := NewRegistry(slog.New(slog.NewTextHandler(io.Discard, nil)))
	registry.Register(testTool("mcp:search:remote-tool", "search_remote_tool"))

	catalogue := registry.ListCatalogue(context.Background())

	require.Len(t, catalogue, 1)
	assert.Equal(t, "mcp:search:remote-tool", catalogue[0].ID)
	assert.Equal(t, "search_remote_tool", catalogue[0].CallableName)
	assert.Equal(t, "mcp", catalogue[0].Source)
	assert.True(t, catalogue[0].Available)
}

func TestRegistryCatalogueReportsRequiresConfirmation(t *testing.T) {
	registry := NewRegistry(slog.New(slog.NewTextHandler(io.Discard, nil)))
	tool := testTool("discord_send", "discord_send")
	tool.Definition.RequiresConfirmation = true
	registry.Register(tool)

	catalogue := registry.ListCatalogue(context.Background())

	require.Len(t, catalogue, 1)
	assert.True(t, catalogue[0].RequiresConfirmation)
}

func TestRegistrySearchCatalogueReturnsCanonicalToolSummaries(t *testing.T) {
	registry := NewRegistry(slog.New(slog.NewTextHandler(io.Discard, nil)))
	tool := testTool("thread_get", "thread_get")
	tool.Definition.Title = "Get thread"
	tool.Definition.Description = "Read one community discussion by ID"
	require.NoError(t, registry.Register(tool))
	registry.RegisterAlias("legacy_thread_read", "thread_get")

	results := registry.SearchCatalogue("read discussion", 5)

	require.Len(t, results, 1)
	assert.Equal(t, "thread_get", results[0].ID)
	assert.Equal(t, "Get thread", results[0].Name)
	assert.Equal(t, "Read one community discussion by ID", results[0].Description)
	assert.NotEqual(t, "legacy_thread_read", results[0].ID)
}

func TestRegistryOmitsToolsetOnlyToolsFromIndividualDiscovery(t *testing.T) {
	registry := NewRegistry(slog.New(slog.NewTextHandler(io.Discard, nil)))
	tool := testTool("document_get", "document_get")
	tool.Definition.ToolsetOnly = true
	tool.Definition.Toolsets = []string{"system.documents"}
	require.NoError(t, registry.Register(tool))

	assert.Empty(t, registry.SearchCatalogue("document", 5))
	err := registry.ValidateStandaloneTool("document_get")
	require.Error(t, err)
	assert.EqualError(t, err, `tool "document_get" is Toolset-only and can only be provided by its declared Toolset; use Toolset system.documents instead`)
}

func TestToolGetReturnsEmptyToolsetArray(t *testing.T) {
	registry := NewRegistry(slog.New(slog.NewTextHandler(io.Discard, nil)))
	require.NoError(t, registry.Register(&Tool{Definition: mcp.GetWebOpenTool()}))
	discovery := newToolDiscoveryTools(registry)

	result, err := discovery.get(context.Background(), mcp.ToolToolGetInput{Id: "web_open"})
	require.NoError(t, err)
	require.NotNil(t, result.Toolsets)
	assert.Empty(t, result.Toolsets)
}

func TestRegistryRejectsToolsWithoutDisplayTitles(t *testing.T) {
	registry := NewRegistry(slog.New(slog.NewTextHandler(io.Discard, nil)))

	err := registry.Register(&Tool{Definition: &mcp.ToolDefinition{
		Name:        "plugin__example__search",
		Description: "Search an external source.",
	}})

	assert.EqualError(t, err, `register Robot tool "plugin__example__search": title is required`)
	assert.False(t, registry.HasTool("plugin__example__search"))
}

func testTool(id, callableName string) *Tool {
	return &Tool{
		Definition: &mcp.ToolDefinition{
			Name:        id,
			Title:       id,
			Description: "test tool",
		},
		CallableName: callableName,
	}
}
