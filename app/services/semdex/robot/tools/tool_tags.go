package tools

import (
	"context"
	"log/slog"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/Southclaws/storyden/app/resources/tag/tag_querier"
	"github.com/Southclaws/storyden/mcp"
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
				tt.ExecuteTagList,
			)
		},
	}
}

func (tt *tagTools) ExecuteTagList(ctx tool.Context, args mcp.ToolTagListInput) (*mcp.ToolTagListOutput, error) {
	var tags []mcp.TagItem

	if args.Query != nil && *args.Query != "" {
		tagList, err := tt.tagQuerier.Search(ctx, *args.Query)
		if err != nil {
			return nil, err
		}

		tags = make([]mcp.TagItem, 0, len(tagList))
		for _, tag := range tagList {
			tags = append(tags, mcp.TagItem{
				Name:      tag.Name.String(),
				ItemCount: tag.ItemCount,
			})
		}
	} else {
		tagList, err := tt.tagQuerier.List(ctx)
		if err != nil {
			return nil, err
		}

		tags = make([]mcp.TagItem, 0, len(tagList))
		for _, tag := range tagList {
			tags = append(tags, mcp.TagItem{
				Name:      tag.Name.String(),
				ItemCount: tag.ItemCount,
			})
		}
	}

	output := mcp.ToolTagListOutput{
		Tags: tags,
	}

	return &output, nil
}
