package tools

import (
	"context"
	"log/slog"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/tag/tag_querier"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/lib/mcp"
)

type tagTools struct {
	logger     *slog.Logger
	tagQuerier *tag_querier.Querier
}

func newTagTools(
	logger *slog.Logger,
	registry *Registry,
	tagQuerier *tag_querier.Querier,
) *tagTools {
	t := &tagTools{
		logger:     logger,
		tagQuerier: tagQuerier,
	}

	registry.Register(t.newTagListTool())

	return t
}

func (tt *tagTools) newTagListTool() *Tool {
	toolDef := mcp.GetTagListTool()

	return &Tool{
		Definition: toolDef,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(
				functiontool.Config{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					InputSchema: toolDef.InputSchema,
				},
				func(ctx tool.Context, args mcp.ToolTagListInput) (*mcp.ToolTagListOutput, error) {
					return tt.ExecuteTagList(ctx, args)
				},
			)
		},
		Handler: makeHandler(tt.ExecuteTagList),
	}
}

const toolTagPageSize = 100

func (tt *tagTools) ExecuteTagList(ctx context.Context, args mcp.ToolTagListInput) (*mcp.ToolTagListOutput, error) {
	params := pagination.NewPageParams(1, toolTagPageSize)

	var result *pagination.Result[*tag_ref.Tag]
	var err error

	if args.Query != nil && *args.Query != "" {
		result, err = tt.tagQuerier.Search(ctx, *args.Query, params)
	} else {
		result, err = tt.tagQuerier.List(ctx, params)
	}
	if err != nil {
		return nil, err
	}

	tags := make([]mcp.TagItem, 0, len(result.Items))
	for _, tag := range result.Items {
		tags = append(tags, mcp.TagItem{
			Name:      tag.Name.String(),
			ItemCount: tag.ItemCount,
		})
	}

	output := mcp.ToolTagListOutput{
		Tags: tags,
	}

	return &output, nil
}
