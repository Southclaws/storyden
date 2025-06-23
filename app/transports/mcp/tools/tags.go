package tools

import (
	"context"
	"encoding/json"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/Southclaws/storyden/app/resources/tag/tag_querier"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
)

type tagTools struct {
	tools []server.ServerTool

	tagQuerier *tag_querier.Querier
}

func newTagTools(
	tagQuerier *tag_querier.Querier,
) *tagTools {
	handler := &tagTools{
		tagQuerier: tagQuerier,
	}

	handler.tools = []server.ServerTool{
		{Tool: tagListTool, Handler: handler.tagList},
	}

	return handler
}

var tagListTool = mcp.NewTool("listTags",
	mcp.WithDescription("Get a list of all tags on the site or search for tags by name using the optional 'query' argument."),
	mcp.WithString("query"),
)

func (t *tagTools) tagList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query := request.GetString("query", "")

	var list tag_ref.Tags
	var err error

	if query == "" {
		list, err = t.tagQuerier.List(ctx)
	} else {
		list, err = t.tagQuerier.Search(ctx, query)
	}

	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	obj := mapTagRefs(list)
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return mcp.NewToolResultText(string(b)), nil
}

func mapTagRef(in *tag_ref.Tag) map[string]any {
	return map[string]any{
		"name":       in.Name.String(),
		"item_count": in.ItemCount,
	}
}

func mapTagRefs(tags tag_ref.Tags) []map[string]any {
	return dt.Map(tags, mapTagRef)
}
