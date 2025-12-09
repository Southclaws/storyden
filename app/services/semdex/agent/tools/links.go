package tools

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/Southclaws/storyden/app/resources/link/link_ref"
	"github.com/Southclaws/storyden/app/services/link/fetcher"
	"github.com/Southclaws/storyden/app/services/link/scrape"
)

type linkTools struct {
	tools []server.ServerTool

	fetcher *fetcher.Fetcher
}

func newLinkTools(
	fetcher *fetcher.Fetcher,
) *linkTools {
	handler := &linkTools{
		fetcher: fetcher,
	}

	handler.tools = []server.ServerTool{
		{Tool: linkCreateTool, Handler: handler.linkCreate},
	}

	return handler
}

var linkCreateTool = mcp.NewTool("createLink",
	mcp.WithDescription("Create or update a link in the shared bookmarks list and return its OpenGraph metadata"),
	mcp.WithString("url", mcp.Required()),
)

func (t *linkTools) linkCreate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	urlStr, err := request.RequireString("url")
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	link, wc, err := t.fetcher.ScrapeAndStore(ctx, *u)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx),
			fmsg.WithDesc("failed to fetch link",
				"The URL could not be fetched. It may be invalid or the server may be unreachable.",
			), ftag.With(ftag.InvalidArgument))
	}

	obj := mapLinkRef(*link, wc)
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return mcp.NewToolResultText(string(b)), nil
}

func mapLinkRef(in link_ref.LinkRef, wc *scrape.WebContent) map[string]any {
	return map[string]any{
		"slug":                  in.Slug,
		"url":                   in.URL,
		"opengraph_title":       in.Title.Ptr(),
		"opengraph_description": in.Description.Ptr(),
		"plain_text":            wc.Content.Plaintext(),
	}
}
