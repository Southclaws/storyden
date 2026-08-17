package tools

import (
	"context"
	"fmt"

	adkagent "google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"

	"github.com/Southclaws/storyden/app/services/semdex/robot/documents"
	"github.com/Southclaws/storyden/lib/mcp"
)

type documentTools struct{}

func newDocumentTools(registry *Registry) *documentTools {
	t := &documentTools{}
	registry.Register(t.newGetTool())
	registry.Register(t.newSearchTool())
	registry.Register(t.newListTool())
	registry.Register(t.newCloseTool())
	return t
}

func (t *documentTools) newGetTool() *Tool {
	def := mcp.GetDocumentGetTool()
	return &Tool{
		Definition: def,
		Builder: func(context.Context) (tool.Tool, error) {
			return functiontool.New(functiontool.Config{
				Name: def.Name, Description: def.Description, InputSchema: def.InputSchema,
			}, func(ctx adkagent.Context, input mcp.ToolDocumentGetInput) (*mcp.RobotDocumentProjectionYaml, error) {
				documentID := ""
				if input.DocumentId != nil {
					documentID = *input.DocumentId
				}
				nodeID := ""
				if input.NodeId != nil {
					nodeID = *input.NodeId
				}
				page := 0
				if input.Page != nil {
					page = *input.Page
				}
				projection, err := documents.Get(ctx.State(), documentID, nodeID, page)
				if err != nil {
					return nil, err
				}
				return mapDocumentProjection(projection), nil
			})
		},
	}
}

func (t *documentTools) newSearchTool() *Tool {
	def := mcp.GetDocumentSearchTool()
	return &Tool{
		Definition: def,
		Builder: func(context.Context) (tool.Tool, error) {
			return functiontool.New(functiontool.Config{
				Name: def.Name, Description: def.Description, InputSchema: def.InputSchema,
			}, func(ctx adkagent.Context, input mcp.ToolDocumentSearchInput) (*mcp.ToolDocumentSearchOutput, error) {
				documentID := ""
				if input.DocumentId != nil {
					documentID = *input.DocumentId
				}
				maxResults := 10
				if input.MaxResults != nil {
					maxResults = *input.MaxResults
				}
				resolvedID, query, matches, truncated, err := documents.Search(ctx.ReadonlyState(), documentID, input.Query, input.NodeIds, maxResults)
				if err != nil {
					return nil, err
				}
				items := make([]mcp.DocumentSearchMatch, 0, len(matches))
				for _, match := range matches {
					items = append(items, mcp.DocumentSearchMatch{
						NodeId: match.NodeID, Kind: match.Kind, Path: match.Path, Preview: match.Preview,
					})
				}
				nextAction := "No matching locations were found; broaden the query or search a different open document."
				if len(items) > 0 {
					nextAction = fmt.Sprintf("Use the smallest sufficient match. If its preview answers the question, continue from this result; otherwise call document_get with document_id %q and that node_id.", resolvedID)
				}
				return &mcp.ToolDocumentSearchOutput{
					DocumentId: resolvedID,
					Query:      query,
					Matches:    items,
					Truncated:  truncated,
					NextAction: nextAction,
				}, nil
			})
		},
	}
}

func (t *documentTools) newListTool() *Tool {
	def := mcp.GetDocumentListTool()
	return &Tool{
		Definition: def,
		Builder: func(context.Context) (tool.Tool, error) {
			return functiontool.New(functiontool.Config{
				Name: def.Name, Description: def.Description, InputSchema: def.InputSchema,
			}, func(ctx adkagent.Context, _ map[string]any) (*mcp.ToolDocumentListOutput, error) {
				documentsList, err := documents.List(ctx.ReadonlyState())
				if err != nil {
					return nil, err
				}
				items := make([]mcp.OpenDocument, 0, len(documentsList))
				for _, item := range documentsList {
					itemStart, itemEnd, totalItems := documentItemWindow(item.ItemStart, item.ItemEnd, item.TotalItems)
					items = append(items, mcp.OpenDocument{
						DocumentId:   item.DocumentID,
						SourceType:   mcp.OpenDocumentSourceType(item.SourceType),
						SourceId:     item.SourceID,
						Title:        item.Title,
						Active:       item.Active,
						ActiveNodeId: item.ActiveNodeID,
						Page:         item.Page,
						TotalPages:   item.TotalPages,
						ItemStart:    itemStart,
						ItemEnd:      itemEnd,
						TotalItems:   totalItems,
					})
				}
				nextAction := "Open a Library page, thread, or web page to begin document exploration."
				if len(items) > 0 {
					nextAction = "Call document_get with a selected document_id, or omit it to inspect the active document."
				}
				return &mcp.ToolDocumentListOutput{Documents: items, NextAction: nextAction}, nil
			})
		},
	}
}

func (t *documentTools) newCloseTool() *Tool {
	def := mcp.GetDocumentCloseTool()
	return &Tool{
		Definition: def,
		Builder: func(context.Context) (tool.Tool, error) {
			return functiontool.New(functiontool.Config{
				Name: def.Name, Description: def.Description, InputSchema: def.InputSchema,
			}, func(ctx adkagent.Context, input mcp.ToolDocumentCloseInput) (*mcp.ToolDocumentCloseOutput, error) {
				documentID := ""
				if input.DocumentId != nil {
					documentID = *input.DocumentId
				}
				closed, active, remaining, err := documents.Close(ctx.State(), documentID)
				if err != nil {
					return nil, err
				}
				var activeID *string
				nextAction := "Open another source when more document exploration is needed."
				if active != "" {
					activeID = &active
					nextAction = fmt.Sprintf("Document %q is now active; call document_get to inspect it.", active)
				}
				return &mcp.ToolDocumentCloseOutput{
					DocumentId: closed, ActiveDocumentId: activeID, Remaining: remaining, NextAction: nextAction,
				}, nil
			})
		},
	}
}

func mapDocumentProjection(projection documents.Projection) *mcp.RobotDocumentProjectionYaml {
	itemStart, itemEnd, totalItems := documentItemWindow(projection.ItemStart, projection.ItemEnd, projection.TotalItems)
	return &mcp.RobotDocumentProjectionYaml{
		DocumentId:   projection.DocumentID,
		SourceType:   mcp.RobotDocumentProjectionYamlSourceType(projection.SourceType),
		SourceId:     projection.SourceID,
		Title:        projection.Title,
		NodeId:       projection.NodeID,
		Page:         projection.Page,
		TotalPages:   projection.TotalPages,
		PreviousPage: projection.Previous,
		NextPage:     projection.Next,
		ItemStart:    itemStart,
		ItemEnd:      itemEnd,
		TotalItems:   totalItems,
		Projection:   projection.Text,
		Truncated:    projection.Truncated,
		NextAction:   projection.NextAction,
	}
}

func documentItemWindow(start, end, total int) (*int, *int, *int) {
	if total <= 0 {
		return nil, nil, nil
	}
	return &start, &end, &total
}
