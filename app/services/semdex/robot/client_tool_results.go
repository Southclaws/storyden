package robot

import (
	"fmt"
	"log/slog"
	"net/url"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/lib/mcp"
	agentpkg "google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

const elicitationHydratedMetadataKey = "storyden_elicitation_hydrated"

func normalizeClientToolResultsBeforeModel(logger *slog.Logger, webAddress url.URL) llmagent.BeforeModelCallback {
	return func(ctx agentpkg.Context, req *model.LLMRequest) (*model.LLMResponse, error) {
		normalised := normalizeClientToolResults(req, webAddress)
		if normalised > 0 {
			logger.Info("normalised client tool results",
				slog.String("agent", ctx.AgentName()),
				slog.String("invocation", ctx.InvocationID()),
				slog.String("session", ctx.SessionID()),
				slog.Int("count", normalised),
			)
		}
		return nil, nil
	}
}

func normalizeClientToolResults(req *model.LLMRequest, webAddress url.URL) int {
	if req == nil {
		return 0
	}

	var count int
	for _, content := range req.Contents {
		if content == nil {
			continue
		}

		for _, part := range content.Parts {
			if part == nil || part.FunctionResponse == nil {
				continue
			}

			switch part.FunctionResponse.Name {
			case mcp.GetLibraryRequestPageTool().Name:
				if isElicitationHydrated(part) {
					continue
				}
				selection := hydrateLibraryPageSelection(part.FunctionResponse.Response, webAddress)
				response := map[string]any{
					"status":    "completed",
					"message":   libraryPageSelectionMessage(selection),
					"selection": selection,
				}
				if browserURL, _ := selection["browser_url"].(string); browserURL != "" {
					response["browser_url"] = browserURL
				}
				part.FunctionResponse.Response = response
				markElicitationHydrated(part)
				count++
			}
		}
	}

	return count
}

func hydrateLibraryPageSelection(selection map[string]any, webAddress url.URL) map[string]any {
	hydrated := make(map[string]any, len(selection)+1)
	for k, v := range selection {
		hydrated[k] = v
	}

	if browserURL := libraryPageSelectionURL(selection, webAddress); browserURL != "" {
		hydrated["browser_url"] = browserURL
	}

	return hydrated
}

func libraryPageSelectionURL(selection map[string]any, webAddress url.URL) string {
	id, _ := selection["id"].(string)
	slug, _ := selection["slug"].(string)
	if id == "" || slug == "" {
		return ""
	}

	mark := fmt.Sprintf("%s-%s", id, slug)
	return datagraph.CanonicalResolveURL(webAddress, datagraph.KindNode, mark).String()
}

func isElicitationHydrated(part *genai.Part) bool {
	if part == nil || part.PartMetadata == nil {
		return false
	}
	hydrated, _ := part.PartMetadata[elicitationHydratedMetadataKey].(bool)
	return hydrated
}

func markElicitationHydrated(part *genai.Part) {
	if part.PartMetadata == nil {
		part.PartMetadata = make(map[string]any)
	}
	part.PartMetadata[elicitationHydratedMetadataKey] = true
}

func libraryPageSelectionMessage(selection map[string]any) string {
	name, _ := selection["name"].(string)
	browserURL, _ := selection["browser_url"].(string)
	if name == "" {
		return "The user selected a Library page."
	}
	if browserURL != "" {
		return fmt.Sprintf("The user selected the Library page %q. Use this URL for normal inline Markdown links or lists when needed: %s. If the page is the single primary result and its ID is available, put one SDR Markdown link alone in its own paragraph instead of adding a separate URL.", name, browserURL)
	}
	return fmt.Sprintf("The user selected the Library page %q.", name)
}
