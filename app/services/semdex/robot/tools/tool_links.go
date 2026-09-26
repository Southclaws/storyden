package tools

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	adkagent "google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"

	"github.com/Southclaws/storyden/app/services/link/fetcher"
	"github.com/Southclaws/storyden/app/services/link/scrape"
	"github.com/Southclaws/storyden/app/services/semdex/robot/documents"
	"github.com/Southclaws/storyden/lib/mcp"
)

type linkTools struct {
	logger  *slog.Logger
	fetcher *fetcher.Fetcher
	scraper scrape.Scraper
}

func newLinkTools(
	logger *slog.Logger,
	registry *Registry,
	fetcher *fetcher.Fetcher,
	scraper scrape.Scraper,
) *linkTools {
	t := &linkTools{
		logger:  logger,
		fetcher: fetcher,
		scraper: scraper,
	}

	registry.Register(t.newLinkCreateTool())
	registry.Register(t.newWebFetchTool())
	registry.Register(t.newWebOpenTool())

	return t
}

func (lt *linkTools) newWebFetchTool() *Tool {
	toolDef := mcp.GetWebFetchTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{Name: toolDef.Name, Description: toolDef.Description, InputSchema: toolDef.InputSchema},
				func(ctx adkagent.Context, args mcp.ToolWebFetchInput) (*mcp.ToolWebFetchOutput, error) {
					return lt.ExecuteWebFetch(ctx, args)
				},
			)
		},
		Handler: makeHandler(lt.ExecuteWebFetch),
	}
}

func (lt *linkTools) ExecuteWebFetch(ctx context.Context, args mcp.ToolWebFetchInput) (*mcp.ToolWebFetchOutput, error) {
	u, err := url.Parse(args.Url)
	if err != nil {
		return nil, err
	}
	wc, err := lt.scraper.Scrape(ctx, *u)
	if err != nil {
		return nil, err
	}
	content, nextAction := contentForAudience(ctx, wc.Content.Plaintext(), webResearchNextAction(ctx, u.String()))
	return &mcp.ToolWebFetchOutput{
		Url:         u.String(),
		Title:       nonEmptyString(wc.Title),
		Description: nonEmptyString(wc.Description),
		FaviconUrl:  nonEmptyString(wc.Favicon),
		ImageUrl:    nonEmptyString(wc.Image),
		Content:     content,
		NextAction:  nextAction,
	}, nil
}

func webResearchNextAction(ctx context.Context, address string) string {
	availability := toolAvailability(ctx, &Tool{Definition: mcp.GetWebOpenTool()})
	if availability == nil {
		return fmt.Sprintf("If source content is needed, discover and activate web_open, then call it with url %q.", address)
	}

	if *availability == mcp.RobotToolAvailabilityYamlCallable {
		if toolIsCallable(ctx, mcp.GetDocumentGetTool()) || toolIsCallable(ctx, mcp.GetDocumentSearchTool()) {
			return fmt.Sprintf("If source content is needed, call web_open with url %q, then use its document_id with document_search or document_get.", address)
		}

		return fmt.Sprintf("If source content is needed, call web_open with url %q and read its returned projection. Discover and activate document_search or document_get before using its document_id for further navigation.", address)
	}

	if toolIsCallable(ctx, mcp.GetToolsetLoadTool()) {
		return "If source content is needed, activate system.web_research with toolset_load. On the next model step, call web_open and use the included document navigation tools."
	}

	if toolIsCallable(ctx, mcp.GetToolLoadTool()) {
		return "If source content is needed, activate web_open with tool_load. On the next model step, call web_open with the source URL and read its returned projection."
	}

	if toolIsCallable(ctx, mcp.GetToolGetTool()) {
		return "Web content reading is blocked. Inspect web_open with tool_get for its preconditions and update this Robot's configured tools or Toolsets."
	}

	return "Web content reading is blocked. Update this Robot's configured tools or Toolsets."
}

func toolIsCallable(ctx context.Context, definition *mcp.ToolDefinition) bool {
	availability := toolAvailability(ctx, &Tool{Definition: definition})
	return availability != nil && *availability == mcp.RobotToolAvailabilityYamlCallable
}

func (lt *linkTools) newWebOpenTool() *Tool {
	toolDef := mcp.GetWebOpenTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{Name: toolDef.Name, Description: toolDef.Description, InputSchema: toolDef.InputSchema},
				func(ctx adkagent.Context, args mcp.ToolWebOpenInput) (*mcp.RobotDocumentProjectionYaml, error) {
					u, err := url.Parse(args.Url)
					if err != nil {
						return nil, err
					}
					wc, err := lt.scraper.Scrape(ctx, *u)
					if err != nil {
						return nil, err
					}
					projection, err := documents.Open(ctx.State(), documents.SourceTypeWeb, u.String(), wc.Title, wc.Content)
					if err != nil {
						return nil, err
					}
					return mapDocumentProjection(projection), nil
				},
			)
		},
	}
}

func nonEmptyString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func (lt *linkTools) newLinkCreateTool() *Tool {
	toolDef := mcp.GetLinkCreateTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				func(ctx adkagent.Context, args mcp.ToolLinkCreateInput) (*mcp.ToolLinkCreateOutput, error) {
					return lt.ExecuteLinkCreate(ctx, args)
				},
			)
		},
		Handler: makeHandler(lt.ExecuteLinkCreate),
	}
}

func (lt *linkTools) ExecuteLinkCreate(ctx context.Context, args mcp.ToolLinkCreateInput) (*mcp.ToolLinkCreateOutput, error) {
	u, err := url.Parse(args.Url)
	if err != nil {
		return nil, err
	}

	link, wc, err := lt.fetcher.ScrapeAndStore(ctx, *u)
	if err != nil {
		return nil, err
	}

	output := mcp.ToolLinkCreateOutput{
		Slug:                 link.Slug,
		Url:                  link.URL,
		OpengraphTitle:       link.Title.Ptr(),
		OpengraphDescription: link.Description.Ptr(),
		PlainText:            func() *string { s := wc.Content.Plaintext(); return &s }(),
	}

	return &(output), nil
}
