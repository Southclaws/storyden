package tools

import (
	"context"
	"log/slog"
	"net/url"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/Southclaws/storyden/app/services/link/fetcher"
	"github.com/Southclaws/storyden/mcp"
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
				lt.ExecuteLinkCreate,
			)
		},
	}
}

func (lt *linkTools) ExecuteLinkCreate(ctx tool.Context, args mcp.ToolLinkCreateInput) ToolResult[mcp.ToolLinkCreateOutput] {
	u, err := url.Parse(args.Url)
	if err != nil {
		return NewError[mcp.ToolLinkCreateOutput](err)
	}

	link, wc, err := lt.fetcher.ScrapeAndStore(ctx, *u)
	if err != nil {
		return NewError[mcp.ToolLinkCreateOutput](err)
	}

	output := mcp.ToolLinkCreateOutput{
		Slug:                  link.Slug,
		Url:                   link.URL,
		OpengraphTitle:        link.Title.Ptr(),
		OpengraphDescription:  link.Description.Ptr(),
		PlainText:             func() *string { s := wc.Content.Plaintext(); return &s }(),
	}

	return NewSuccess(output)
}
