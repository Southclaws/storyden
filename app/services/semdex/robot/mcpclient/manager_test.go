package mcpclient

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/infrastructure/httpsafe"
)

func TestFetchServerCardRefusesInternalHosts(t *testing.T) {
	t.Parallel()

	client := httpsafe.NewClient(httpsafe.Config{DialTimeout: time.Second})

	for _, raw := range []string{
		"http://127.0.0.1/",
		"http://localhost/",
		"http://169.254.169.254/",
	} {
		u, err := url.Parse(raw)
		require.NoError(t, err)

		_, _, err = fetchServerCard(context.Background(), client, u)
		assert.Error(t, err, raw)
	}
}

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
