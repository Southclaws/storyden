package tools

import (
	"context"
	"log/slog"
	"net/url"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"

	"github.com/Southclaws/storyden/app/services/link/fetcher"
	"github.com/Southclaws/storyden/lib/mcp"
)

type linkTools struct {
	logger  *slog.Logger
	fetcher *fetcher.Fetcher
}

func newLinkTools(
	logger *slog.Logger,
	registry *Registry,
	fetcher *fetcher.Fetcher,
) *linkTools {
	t := &linkTools{
		logger:  logger,
		fetcher: fetcher,
	}

	registry.Register(t.newLinkCreateTool())

	return t
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
				func(ctx agent.Context, args mcp.ToolLinkCreateInput) (*mcp.ToolLinkCreateOutput, error) {
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
