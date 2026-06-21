package mcpclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCallableNameUsesReadableProviderPrefix(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "notion_mcp_search", CallableName("Notion MCP (Beta)", "notion-search"))
	assert.Equal(t, "notion_mcp_fetch", CallableName("Notion MCP (Beta)", "notion-fetch"))
	assert.Equal(t, "notion_mcp_create_pages", CallableName("Notion MCP (Beta)", "notion-create-pages"))
	assert.Equal(t, "notion_mcp_update_page", CallableName("Notion MCP (Beta)", "notion-update-page"))
}

func TestCallableNameKeepsProviderWhenRemoteNameDoesNotRepeatIt(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "github_search_repositories", CallableName("GitHub", "search-repositories"))
}

func TestCallableNameCleansDomainFallbackNames(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "notion_search", CallableName("mcp.notion.com", "notion-search"))
	assert.Equal(t, "notion_search", CallableName("mcp-notion-com", "notion-search"))
}
