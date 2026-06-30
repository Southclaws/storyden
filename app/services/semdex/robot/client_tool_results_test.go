package robot

import (
	"net/url"
	"testing"

	"github.com/Southclaws/storyden/lib/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

func TestNormalizeClientToolResultsWrapsLibraryPageSelection(t *testing.T) {
	webAddress := mustParseURL(t, "https://example.com")
	req := &model.LLMRequest{
		Contents: []*genai.Content{
			{
				Role: genai.RoleUser,
				Parts: []*genai.Part{
					{
						FunctionResponse: &genai.FunctionResponse{
							ID:   "call_123",
							Name: mcp.GetLibraryRequestPageTool().Name,
							Response: map[string]any{
								"id":   "d8818ueot5pfij6bvm90",
								"name": "Documentation Hub",
								"slug": "documentation-hub",
							},
						},
					},
				},
			},
		},
	}

	count := normalizeClientToolResults(req, *webAddress)
	assert.Equal(t, 1, count)

	response := req.Contents[0].Parts[0].FunctionResponse.Response
	assert.Equal(t, "completed", response["status"])
	assert.Equal(t, `The user selected the Library page "Documentation Hub". Use this URL for normal inline Markdown links or lists when needed: https://example.com/_/resolve/node/d8818ueot5pfij6bvm90-documentation-hub. If the page is the single primary result and its ID is available, put one SDR Markdown link alone in its own paragraph instead of adding a separate URL.`, response["message"])
	assert.Equal(t, "https://example.com/_/resolve/node/d8818ueot5pfij6bvm90-documentation-hub", response["browser_url"])
	assert.True(t, isElicitationHydrated(req.Contents[0].Parts[0]))

	selection, ok := response["selection"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "Documentation Hub", selection["name"])
	assert.Equal(t, "documentation-hub", selection["slug"])
	assert.Equal(t, "https://example.com/_/resolve/node/d8818ueot5pfij6bvm90-documentation-hub", selection["browser_url"])
}

func TestNormalizeClientToolResultsSkipsAlreadyNormalisedSelection(t *testing.T) {
	response := map[string]any{
		"status":    "completed",
		"message":   "The user selected a Library page.",
		"selection": map[string]any{"name": "Documentation Hub"},
	}

	req := &model.LLMRequest{
		Contents: []*genai.Content{
			{
				Role: genai.RoleUser,
				Parts: []*genai.Part{
					{
						FunctionResponse: &genai.FunctionResponse{
							ID:       "call_123",
							Name:     mcp.GetLibraryRequestPageTool().Name,
							Response: response,
						},
						PartMetadata: map[string]any{
							elicitationHydratedMetadataKey: true,
						},
					},
				},
			},
		},
	}

	count := normalizeClientToolResults(req, url.URL{})
	assert.Equal(t, 0, count)
	assert.Equal(t, response, req.Contents[0].Parts[0].FunctionResponse.Response)
}

func TestNormalizeClientToolResultsUsesMetadataNotResponseShape(t *testing.T) {
	req := &model.LLMRequest{
		Contents: []*genai.Content{
			{
				Role: genai.RoleUser,
				Parts: []*genai.Part{
					{
						FunctionResponse: &genai.FunctionResponse{
							ID:   "call_123",
							Name: mcp.GetLibraryRequestPageTool().Name,
							Response: map[string]any{
								"status":    "draft",
								"message":   "Raw selected page field",
								"selection": "Raw selected page field",
								"name":      "Documentation Hub",
							},
						},
					},
				},
			},
		},
	}

	count := normalizeClientToolResults(req, url.URL{})
	assert.Equal(t, 1, count)

	response := req.Contents[0].Parts[0].FunctionResponse.Response
	assert.Equal(t, "completed", response["status"])

	selection, ok := response["selection"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "draft", selection["status"])
	assert.Equal(t, "Raw selected page field", selection["message"])
}

func TestNormalizeClientToolResultsLeavesOtherToolsUnchanged(t *testing.T) {
	response := map[string]any{"ok": true}
	req := &model.LLMRequest{
		Contents: []*genai.Content{
			{
				Role: genai.RoleUser,
				Parts: []*genai.Part{
					{
						FunctionResponse: &genai.FunctionResponse{
							ID:       "call_123",
							Name:     "robot_switch",
							Response: response,
						},
					},
				},
			},
		},
	}

	count := normalizeClientToolResults(req, url.URL{})
	assert.Equal(t, 0, count)
	assert.Equal(t, response, req.Contents[0].Parts[0].FunctionResponse.Response)
}

func TestNormalizeClientToolResultsHandlesSelectionWithoutURLFields(t *testing.T) {
	req := &model.LLMRequest{
		Contents: []*genai.Content{
			{
				Role: genai.RoleUser,
				Parts: []*genai.Part{
					{
						FunctionResponse: &genai.FunctionResponse{
							ID:   "call_123",
							Name: mcp.GetLibraryRequestPageTool().Name,
							Response: map[string]any{
								"name": "Documentation Hub",
							},
						},
					},
				},
			},
		},
	}

	count := normalizeClientToolResults(req, *mustParseURL(t, "https://example.com"))
	assert.Equal(t, 1, count)

	response := req.Contents[0].Parts[0].FunctionResponse.Response
	assert.Equal(t, `The user selected the Library page "Documentation Hub".`, response["message"])
	_, hasTopLevelURL := response["browser_url"]
	assert.False(t, hasTopLevelURL)

	selection, ok := response["selection"].(map[string]any)
	require.True(t, ok)
	_, hasURL := selection["browser_url"]
	assert.False(t, hasURL)
}

func mustParseURL(t *testing.T, raw string) *url.URL {
	t.Helper()
	u, err := url.Parse(raw)
	require.NoError(t, err)
	return u
}
