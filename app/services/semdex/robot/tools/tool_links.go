package tools

import (
	"context"
	"fmt"
	adkagent "google.golang.org/adk/v2/agent"
	"log/slog"
	"net/url"

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
	content, nextAction := contentForAudience(ctx,
		wc.Content.Plaintext(),
		fmt.Sprintf("If this page looks relevant to the current task, call web_open with url %q to inspect its document.", u.String()),
	)
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
