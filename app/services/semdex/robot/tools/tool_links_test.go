package tools

import (
	"context"
	"encoding/json"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/tool"
	"net/url"
	"strings"
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

func TestWebOpenCanSearchAndReadSourceContent(t *testing.T) {
	content, err := datagraph.NewRichText(`<h2>Storage</h2><p>` + strings.Repeat("General storage details. ", 50) + `</p><h2>Limitations</h2><p>Indexes must fit in memory.</p>`)
	require.NoError(t, err)
	links := &linkTools{scraper: staticScraper{content: &scrape.WebContent{Title: "Technical project", Content: content}}}
	ctx := &toolTestContext{state: toolTestState{}}
	open, err := links.newWebOpenTool().Builder(context.Background())
	require.NoError(t, err)
	opened := runDocumentTool(t, open, ctx, map[string]any{"url": "https://example.com/project"})
	documentID := opened["document_id"].(string)
	assert.Equal(t, true, opened["truncated"])
	docs := &documentTools{}
	search, err := docs.newSearchTool().Builder(context.Background())
	require.NoError(t, err)
	matches := runDocumentTool(t, search, ctx, map[string]any{"document_id": documentID, "query": "Indexes must fit"})
	items := matches["matches"].([]any)
	require.NotEmpty(t, items)
	nodeID := items[0].(map[string]any)["node_id"].(string)
	get, err := docs.newGetTool().Builder(context.Background())
	require.NoError(t, err)
	read := runDocumentTool(t, get, ctx, map[string]any{"document_id": documentID, "node_id": nodeID})
	assert.Contains(t, read["projection"], "Indexes must fit in memory.")
	assert.Equal(t, documentID, read["document_id"])
}

func runDocumentTool(t *testing.T, selected tool.Tool, ctx agent.Context, args map[string]any) map[string]any {
	t.Helper()
	runnable := selected.(interface {
		Run(agent.Context, any) (map[string]any, error)
	})
	output, err := runnable.Run(ctx, args)
	require.NoError(t, err)
	encoded, err := json.Marshal(output)
	require.NoError(t, err)
	var result map[string]any
	require.NoError(t, json.Unmarshal(encoded, &result))
	return result
}

func TestWebResearchGuidanceUsesConversationAvailability(t *testing.T) {
	ctx := &toolTestContext{state: toolTestState{}}
	_, err := CaptureCallableTools()(ctx, &model.LLMRequest{Tools: map[string]any{"document_search": nil, "web_open": nil}})
	require.NoError(t, err)
	assert.Contains(t, webResearchNextAction(ctx, "https://example.com"), "document_id")

	_, err = CaptureCallableTools()(ctx, &model.LLMRequest{Tools: map[string]any{"web_open": nil}})
	require.NoError(t, err)
	assert.Contains(t, webResearchNextAction(ctx, "https://example.com"), "returned projection")
	assert.Contains(t, webResearchNextAction(ctx, "https://example.com"), "Discover and activate")

	_, err = CaptureCallableTools()(ctx, &model.LLMRequest{Tools: map[string]any{"tool_load": nil}})
	require.NoError(t, err)
	assert.Contains(t, webResearchNextAction(ctx, "https://example.com"), "activate web_open with tool_load")
	assert.NotContains(t, webResearchNextAction(ctx, "https://example.com"), "toolset_load")

	_, err = CaptureCallableTools()(ctx, &model.LLMRequest{Tools: map[string]any{"toolset_load": nil}})
	require.NoError(t, err)
	assert.Contains(t, webResearchNextAction(ctx, "https://example.com"), "activate system.web_research with toolset_load")

	_, err = CaptureCallableTools()(ctx, &model.LLMRequest{Tools: map[string]any{"tool_load": nil, "toolset_load": nil}})
	require.NoError(t, err)
	assert.Contains(t, webResearchNextAction(ctx, "https://example.com"), "activate system.web_research with toolset_load")

	_, err = CaptureCallableTools()(ctx, &model.LLMRequest{Tools: map[string]any{}})
	require.NoError(t, err)
	assert.Contains(t, webResearchNextAction(ctx, "https://example.com"), "configured tools")
	assert.NotContains(t, webResearchNextAction(ctx, "https://example.com"), "toolset_load")
}
