package tools

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/services/link/scrape"
	"github.com/Southclaws/storyden/lib/mcp"
)

func TestWebFetchShapesContentForAudience(t *testing.T) {
	content, err := datagraph.NewRichText(`<h1>Example</h1><p>Fetched body.</p>`)
	require.NoError(t, err)
	tools := &linkTools{scraper: staticScraper{content: &scrape.WebContent{
		Title: "Example", Description: "Example description", Content: content,
	}}}
	input := mcp.ToolWebFetchInput{Url: "https://example.com/article"}

	robotOutput, err := tools.ExecuteWebFetch(ContextWithToolAudience(context.Background(), ToolAudienceRobot), input)
	require.NoError(t, err)
	assert.Nil(t, robotOutput.Content)
	require.NotNil(t, robotOutput.NextAction)
	assert.Contains(t, *robotOutput.NextAction, "web_open")

	mcpOutput, err := tools.ExecuteWebFetch(ContextWithToolAudience(context.Background(), ToolAudienceMCP), input)
	require.NoError(t, err)
	require.NotNil(t, mcpOutput.Content)
	assert.Contains(t, *mcpOutput.Content, "Fetched body")
	assert.Nil(t, mcpOutput.NextAction)
}

type staticScraper struct {
	content *scrape.WebContent
}

func (s staticScraper) Scrape(context.Context, url.URL) (*scrape.WebContent, error) {
	return s.content, nil
}
