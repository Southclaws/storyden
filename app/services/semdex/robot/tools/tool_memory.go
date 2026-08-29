package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/Southclaws/opt"
	adkagent "google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"

	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_memory"
	"github.com/Southclaws/storyden/lib/mcp"
)

type memoryTools struct {
	repo *robot_memory.Repository
}

func newMemoryTools(registry *Registry, repo *robot_memory.Repository) *memoryTools {
	t := &memoryTools{repo: repo}
	registry.Register(t.newListTool())
	registry.Register(t.newOpenTool())
	registry.Register(t.newSearchTool())
	registry.Register(t.newCreateTool())
	registry.Register(t.newUpdateTool())
	registry.Register(t.newMoveTool())
	return t
}

func memoryRobotRef(ctx context.Context) (string, error) {
	robotRef := RunContextFromContext(ctx).RobotRef
	if robotRef == "" {
		return "", fmt.Errorf("memory tools require an active Robot context")
	}
	return robotRef, nil
}

func (t *memoryTools) newListTool() *Tool {
	definition := mcp.GetMemoryListTool()
	return &Tool{
		Definition: definition,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			robotRef, err := memoryRobotRef(ctx)
			if err != nil {
				return nil, err
			}
			return functiontool.New(functiontool.Config{Name: definition.Name, Description: definition.Description, InputSchema: definition.InputSchema}, func(ctx adkagent.Context, input mcp.ToolMemoryListInput) (*mcp.ToolMemoryListOutput, error) {
				return t.executeList(ctx, robotRef, input)
			})
		},
		Handler: makeHandler(func(ctx context.Context, input mcp.ToolMemoryListInput) (*mcp.ToolMemoryListOutput, error) {
			robotRef, err := memoryRobotRef(ctx)
			if err != nil {
				return nil, err
			}
			return t.executeList(ctx, robotRef, input)
		}),
	}
}

func (t *memoryTools) executeList(ctx context.Context, robotRef string, input mcp.ToolMemoryListInput) (*mcp.ToolMemoryListOutput, error) {
	parentID, err := optionalMemoryID(input.ParentId)
	if err != nil {
		return nil, err
	}
	items, hasMore, err := t.repo.List(ctx, robotRef, parentID, robot_memory.ListLimit)
	if err != nil {
		return nil, err
	}
	return &mcp.ToolMemoryListOutput{Memories: mapMemoryItems(items), Returned: len(items), HasMore: hasMore, NextAction: listNextAction(hasMore)}, nil
}

func (t *memoryTools) newOpenTool() *Tool {
	definition := mcp.GetMemoryOpenTool()
	return &Tool{
		Definition: definition,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			robotRef, err := memoryRobotRef(ctx)
			if err != nil {
				return nil, err
			}
			return functiontool.New(functiontool.Config{Name: definition.Name, Description: definition.Description, InputSchema: definition.InputSchema}, func(ctx adkagent.Context, input mcp.ToolMemoryOpenInput) (*mcp.ToolMemoryOpenOutput, error) {
				return t.executeOpen(ctx, robotRef, input)
			})
		},
		Handler: makeHandler(func(ctx context.Context, input mcp.ToolMemoryOpenInput) (*mcp.ToolMemoryOpenOutput, error) {
			robotRef, err := memoryRobotRef(ctx)
			if err != nil {
				return nil, err
			}
			return t.executeOpen(ctx, robotRef, input)
		}),
	}
}

func (t *memoryTools) executeOpen(ctx context.Context, robotRef string, input mcp.ToolMemoryOpenInput) (*mcp.ToolMemoryOpenOutput, error) {
	id, err := robot.NewMemoryID(input.Id)
	if err != nil {
		return nil, err
	}
	detail, err := t.repo.Open(ctx, robotRef, id)
	if err != nil {
		return nil, err
	}
	return &mcp.ToolMemoryOpenOutput{
		Memory: mapMemoryRecord(detail.Memory), Path: mapMemoryItems(detail.Path), Children: mapMemoryItems(detail.Children),
		NextAction: "Use memory_update to correct this knowledge graph node, memory_list to navigate its children, or a focused memory_search to find related facts.",
	}, nil
}

func (t *memoryTools) newSearchTool() *Tool {
	definition := mcp.GetMemorySearchTool()
	return &Tool{
		Definition: definition,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			robotRef, err := memoryRobotRef(ctx)
			if err != nil {
				return nil, err
			}
			return functiontool.New(functiontool.Config{Name: definition.Name, Description: definition.Description, InputSchema: definition.InputSchema}, func(ctx adkagent.Context, input mcp.ToolMemorySearchInput) (*mcp.ToolMemorySearchOutput, error) {
				return t.executeSearch(ctx, robotRef, input)
			})
		},
		Handler: makeHandler(func(ctx context.Context, input mcp.ToolMemorySearchInput) (*mcp.ToolMemorySearchOutput, error) {
			robotRef, err := memoryRobotRef(ctx)
			if err != nil {
				return nil, err
			}
			return t.executeSearch(ctx, robotRef, input)
		}),
	}
}

func (t *memoryTools) executeSearch(ctx context.Context, robotRef string, input mcp.ToolMemorySearchInput) (*mcp.ToolMemorySearchOutput, error) {
	parentID, err := optionalMemoryID(input.ParentId)
	if err != nil {
		return nil, err
	}
	query := ""
	if input.Query != nil {
		query = *input.Query
	}
	items, hasMore, err := t.repo.Search(ctx, robotRef, robot_memory.SearchFilter{
		Query:     query,
		Subject:   opt.NewPtr(input.Subject),
		Predicate: opt.NewPtr(input.Predicate),
		Object:    opt.NewPtr(input.Object),
	}, parentID, robot_memory.SearchLimit)
	if err != nil {
		return nil, err
	}
	results := make([]mcp.RobotMemorySearchResultYaml, len(items))
	for i, item := range items {
		subject, predicate, object := memoryFactFields(item.Memory)
		results[i] = mcp.RobotMemorySearchResultYaml{
			MemoryId: item.Memory.ID.String(), Excerpt: item.Excerpt,
			Subject: subject, Predicate: predicate, Object: object,
		}
	}
	nextAction := "No matching memory was found. Continue the current task; save a useful lasting fact if this conversation revealed one."
	if len(items) > 0 {
		nextAction = "Use the matching evidence when relevant and continue the current task."
	}
	if hasMore {
		nextAction = "Results were truncated. Narrow the search only if the current task needs a more specific memory."
	}
	return &mcp.ToolMemorySearchOutput{Results: results, Returned: len(results), HasMore: hasMore, NextAction: nextAction}, nil
}

func (t *memoryTools) newCreateTool() *Tool {
	definition := mcp.GetMemoryCreateTool()
	return &Tool{
		Definition: definition,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			robotRef, err := memoryRobotRef(ctx)
			if err != nil {
				return nil, err
			}
			return functiontool.New(functiontool.Config{Name: definition.Name, Description: definition.Description, InputSchema: definition.InputSchema}, func(ctx adkagent.Context, input mcp.ToolMemoryCreateInput) (*mcp.ToolMemoryCreateOutput, error) {
				return t.executeCreate(ctx, robotRef, input)
			})
		},
		Handler: makeHandler(func(ctx context.Context, input mcp.ToolMemoryCreateInput) (*mcp.ToolMemoryCreateOutput, error) {
			robotRef, err := memoryRobotRef(ctx)
			if err != nil {
				return nil, err
			}
			return t.executeCreate(ctx, robotRef, input)
		}),
	}
}

func (t *memoryTools) executeCreate(ctx context.Context, robotRef string, input mcp.ToolMemoryCreateInput) (*mcp.ToolMemoryCreateOutput, error) {
	parentID, err := optionalMemoryID(input.ParentId)
	if err != nil {
		return nil, err
	}
	options, err := memoryFactOptions(input.Subject, input.Predicate, input.Object, false)
	if err != nil {
		return nil, err
	}
	created, err := t.repo.Create(ctx, robotRef, parentID, input.Content, options...)
	if err != nil {
		return nil, err
	}
	subject, predicate, object := memoryFactFields(created)
	return &mcp.ToolMemoryCreateOutput{
		MemoryId: created.ID.String(), Subject: subject, Predicate: predicate, Object: object,
		Message: "Memory saved.", NextAction: "Continue the current task; organize memories only when necessary.",
	}, nil
}

func (t *memoryTools) newUpdateTool() *Tool {
	definition := mcp.GetMemoryUpdateTool()
	return &Tool{
		Definition: definition,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			robotRef, err := memoryRobotRef(ctx)
			if err != nil {
				return nil, err
			}
			return functiontool.New(functiontool.Config{Name: definition.Name, Description: definition.Description, InputSchema: definition.InputSchema}, func(ctx adkagent.Context, input mcp.ToolMemoryUpdateInput) (*mcp.ToolMemoryUpdateOutput, error) {
				return t.executeUpdate(ctx, robotRef, input)
			})
		},
		Handler: makeHandler(func(ctx context.Context, input mcp.ToolMemoryUpdateInput) (*mcp.ToolMemoryUpdateOutput, error) {
			robotRef, err := memoryRobotRef(ctx)
			if err != nil {
				return nil, err
			}
			return t.executeUpdate(ctx, robotRef, input)
		}),
	}
}

func (t *memoryTools) executeUpdate(ctx context.Context, robotRef string, input mcp.ToolMemoryUpdateInput) (*mcp.ToolMemoryUpdateOutput, error) {
	id, err := robot.NewMemoryID(input.Id)
	if err != nil {
		return nil, err
	}
	options := make([]robot_memory.Option, 0, 5)
	if input.Content != nil {
		options = append(options, robot_memory.WithContent(*input.Content))
	}
	factOptions, err := memoryFactOptions(input.Subject, input.Predicate, input.Object, true)
	if err != nil {
		return nil, err
	}
	options = append(options, factOptions...)
	if input.State != nil {
		parsed, err := robot.NewMemoryState(string(*input.State))
		if err != nil {
			return nil, err
		}
		options = append(options, robot_memory.WithState(parsed))
	}
	updated, err := t.repo.Update(ctx, robotRef, id, options...)
	if err != nil {
		return nil, err
	}
	return &mcp.ToolMemoryUpdateOutput{
		Memory: mapMemoryRecord(updated), Message: fmt.Sprintf("Updated shared memory %s.", updated.ID),
		NextAction: "Use memory_open to verify the complete node or a focused memory_search to confirm its knowledge graph fields.",
	}, nil
}

func (t *memoryTools) newMoveTool() *Tool {
	definition := mcp.GetMemoryMoveTool()
	return &Tool{
		Definition: definition,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			robotRef, err := memoryRobotRef(ctx)
			if err != nil {
				return nil, err
			}
			return functiontool.New(functiontool.Config{Name: definition.Name, Description: definition.Description, InputSchema: definition.InputSchema}, func(ctx adkagent.Context, input mcp.ToolMemoryMoveInput) (*mcp.ToolMemoryMoveOutput, error) {
				return t.executeMove(ctx, robotRef, input)
			})
		},
		Handler: makeHandler(func(ctx context.Context, input mcp.ToolMemoryMoveInput) (*mcp.ToolMemoryMoveOutput, error) {
			robotRef, err := memoryRobotRef(ctx)
			if err != nil {
				return nil, err
			}
			return t.executeMove(ctx, robotRef, input)
		}),
	}
}

func (t *memoryTools) executeMove(ctx context.Context, robotRef string, input mcp.ToolMemoryMoveInput) (*mcp.ToolMemoryMoveOutput, error) {
	id, err := robot.NewMemoryID(input.Id)
	if err != nil {
		return nil, err
	}
	parentID, err := optionalMemoryID(input.ParentId)
	if err != nil {
		return nil, err
	}
	moved, err := t.repo.Move(ctx, robotRef, id, parentID)
	if err != nil {
		return nil, err
	}
	return &mcp.ToolMemoryMoveOutput{
		Memory: mapMemoryRecord(moved), Message: fmt.Sprintf("Moved shared memory %s.", moved.ID),
		NextAction: "Use memory_list on the new parent, or list the root when the memory was moved to the top level.",
	}, nil
}

func optionalMemoryID(value *string) (opt.Optional[robot.MemoryID], error) {
	if value == nil {
		return opt.NewEmpty[robot.MemoryID](), nil
	}
	id, err := robot.NewMemoryID(*value)
	if err != nil {
		return opt.NewEmpty[robot.MemoryID](), err
	}
	return opt.New(id), nil
}

func mapMemoryRecord(item *robot.Memory) mcp.RobotMemoryRecordYaml {
	subject, predicate, object := memoryFactFields(item)
	return mcp.RobotMemoryRecordYaml{
		Id: item.ID.String(), RobotRef: item.RobotRef, ParentId: opt.Map(item.ParentID, func(id robot.MemoryID) string { return id.String() }).Ptr(), Content: item.Content,
		Subject: subject, Predicate: predicate, Object: object, State: mcp.RobotMemoryRecordYamlState(item.State.String()),
		CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt, LastAccessedAt: item.LastAccessedAt, AccessCount: int(item.AccessCount),
	}
}

func mapMemoryItems(items []*robot_memory.Item) []mcp.RobotMemorySummaryYaml {
	result := make([]mcp.RobotMemorySummaryYaml, len(items))
	for i, item := range items {
		memory := item.Memory
		subject, predicate, object := memoryFactFields(memory)
		result[i] = mcp.RobotMemorySummaryYaml{
			Id: memory.ID.String(), Excerpt: item.Excerpt,
			Subject: subject, Predicate: predicate, Object: object,
			State: mcp.RobotMemorySummaryYamlState(memory.State.String()), Children: item.Children,
		}
	}
	return result
}

func memoryFactOptions(subject, predicate, object *string, allowClear bool) ([]robot_memory.Option, error) {
	provided := 0
	for _, value := range []*string{subject, predicate, object} {
		if value != nil {
			provided++
		}
	}
	if provided == 0 {
		return nil, nil
	}
	if provided != 3 {
		return nil, fmt.Errorf("subject, predicate, and object must be provided together")
	}
	values := []string{strings.TrimSpace(*subject), strings.TrimSpace(*predicate), strings.TrimSpace(*object)}
	empty := 0
	for _, value := range values {
		if value == "" {
			empty++
		}
	}
	if empty == 3 && allowClear {
		return []robot_memory.Option{robot_memory.WithoutFact()}, nil
	}
	if empty > 0 {
		return nil, fmt.Errorf("subject, predicate, and object must all be non-empty")
	}
	return []robot_memory.Option{robot_memory.WithFact(values[0], values[1], values[2])}, nil
}

func memoryFactFields(memory *robot.Memory) (subject, predicate, object *string) {
	fact, ok := memory.Fact.Get()
	if !ok {
		return nil, nil, nil
	}
	return &fact.Subject, &fact.Predicate, &fact.Object
}

func listNextAction(hasMore bool) string {
	if hasMore {
		return "This knowledge graph branch was truncated. Use a focused memory_search or open a known child to narrow navigation."
	}
	return "Open a memory to read its complete content, or list one of its children to navigate deeper."
}
