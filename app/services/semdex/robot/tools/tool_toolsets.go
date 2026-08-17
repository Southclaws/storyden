package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/Southclaws/opt"
	adkagent "google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"

	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/robot/robot_toolset"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/lib/mcp"
)

type toolsetTools struct {
	repo     *robot_toolset.Repository
	registry *Registry
	resolver ToolsetResolver
}

func newToolsetTools(repo *robot_toolset.Repository, registry *Registry, resolver ToolsetResolver) *toolsetTools {
	t := &toolsetTools{repo: repo, registry: registry, resolver: resolver}
	registry.Register(t.newCreateTool())
	registry.Register(t.newSearchTool())
	registry.Register(t.newListTool())
	registry.Register(t.newGetTool())
	registry.Register(t.newUpdateTool())
	registry.Register(t.newDeleteTool())
	return t
}

func (t *toolsetTools) newSearchTool() *Tool {
	def := mcp.GetToolsetSearchTool()
	return &Tool{
		Definition: def,
		Builder: func(context.Context) (tool.Tool, error) {
			return functiontool.New(functiontool.Config{Name: def.Name, Description: def.Description, InputSchema: def.InputSchema},
				func(ctx adkagent.Context, input mcp.ToolToolsetSearchInput) (*mcp.ToolToolsetSearchOutput, error) {
					return t.search(ctx, input)
				})
		},
		Handler: makeHandler(t.search),
	}
}

func (t *toolsetTools) search(ctx context.Context, input mcp.ToolToolsetSearchInput) (*mcp.ToolToolsetSearchOutput, error) {
	maxResults := 0
	if input.MaxResults != nil {
		maxResults = *input.MaxResults
	}
	items, err := t.resolver.SearchToolsets(ctx, input.Query, maxResults)
	if err != nil {
		return nil, err
	}
	return &mcp.ToolToolsetSearchOutput{
		Toolsets:   mapToolsetCatalogue(items),
		NextAction: "Use toolset_get with a returned ID to inspect its tools and instruction before loading or assigning it.",
	}, nil
}

func (t *toolsetTools) newCreateTool() *Tool {
	def := mcp.GetToolsetCreateTool()
	return &Tool{
		Definition: def,
		Builder: func(context.Context) (tool.Tool, error) {
			return functiontool.New(functiontool.Config{Name: def.Name, Description: def.Description, InputSchema: def.InputSchema},
				func(ctx adkagent.Context, input mcp.ToolToolsetCreateInput) (*mcp.ToolToolsetCreateOutput, error) {
					return t.create(ctx, input)
				})
		},
		Handler: makeHandler(t.create),
	}
}

func (t *toolsetTools) create(ctx context.Context, input mcp.ToolToolsetCreateInput) (*mcp.ToolToolsetCreateOutput, error) {
	if invalid := t.invalidTools(input.Tools); len(invalid) > 0 {
		return nil, fmt.Errorf("unknown tools: %s", strings.Join(invalid, ", "))
	}
	if err := t.validateStandaloneTools(input.Tools); err != nil {
		return nil, err
	}
	authorID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, err
	}
	created, err := t.repo.Create(ctx, authorID, input.Name, input.Description, opt.NewPtr(input.Instruction).Or(""), input.Tools, map[string]any{})
	if err != nil {
		return nil, err
	}
	return &mcp.ToolToolsetCreateOutput{Id: created.ID.String(), Name: created.Name}, nil
}

func (t *toolsetTools) newListTool() *Tool {
	def := mcp.GetToolsetListTool()
	return &Tool{
		Definition: def,
		Builder: func(context.Context) (tool.Tool, error) {
			return functiontool.New(functiontool.Config{Name: def.Name, Description: def.Description, InputSchema: def.InputSchema},
				func(ctx adkagent.Context, _ struct{}) (*mcp.ToolToolsetListOutput, error) {
					return t.list(ctx)
				})
		},
	}
}

func (t *toolsetTools) list(ctx context.Context) (*mcp.ToolToolsetListOutput, error) {
	items, err := t.resolver.ListToolsets(ctx)
	if err != nil {
		return nil, err
	}
	return &mcp.ToolToolsetListOutput{Toolsets: mapToolsetCatalogue(items)}, nil
}

func (t *toolsetTools) newGetTool() *Tool {
	def := mcp.GetToolsetGetTool()
	return &Tool{
		Definition: def,
		Builder: func(context.Context) (tool.Tool, error) {
			return functiontool.New(functiontool.Config{Name: def.Name, Description: def.Description, InputSchema: def.InputSchema},
				func(ctx adkagent.Context, input mcp.ToolToolsetGetInput) (*mcp.ToolToolsetGetOutput, error) {
					return t.get(ctx, input)
				})
		},
		Handler: makeHandler(t.get),
	}
}

func (t *toolsetTools) get(ctx context.Context, input mcp.ToolToolsetGetInput) (*mcp.ToolToolsetGetOutput, error) {
	item, err := t.resolver.GetToolset(ctx, input.Id)
	if err != nil {
		return nil, err
	}
	editable := item.Editable
	if editable {
		id, err := robot_toolset.ParseID(item.ID)
		editable = err == nil && t.canMutate(ctx, id)
	}
	return &mcp.ToolToolsetGetOutput{
		Id: item.ID, Name: item.Name, Description: item.Description, Instruction: item.Instruction,
		Tools: item.Tools, Editable: editable, RequiresWorkspace: item.RequiresWorkspace,
	}, nil
}

func (t *toolsetTools) newUpdateTool() *Tool {
	def := mcp.GetToolsetUpdateTool()
	return &Tool{
		Definition: def,
		Builder: func(context.Context) (tool.Tool, error) {
			return functiontool.New(functiontool.Config{Name: def.Name, Description: def.Description, InputSchema: def.InputSchema},
				func(ctx adkagent.Context, input mcp.ToolToolsetUpdateInput) (*mcp.ToolToolsetUpdateOutput, error) {
					return t.update(ctx, input)
				})
		},
		Handler: makeHandler(t.update),
	}
}

func (t *toolsetTools) update(ctx context.Context, input mcp.ToolToolsetUpdateInput) (*mcp.ToolToolsetUpdateOutput, error) {
	id, err := robot_toolset.ParseID(input.Id)
	if err != nil {
		return nil, fmt.Errorf("only custom Toolsets can be edited: %w", err)
	}
	if len(input.Tools) > 0 {
		if invalid := t.invalidTools(input.Tools); len(invalid) > 0 {
			return nil, fmt.Errorf("unknown tools: %s", strings.Join(invalid, ", "))
		}
		if err := t.validateStandaloneTools(input.Tools); err != nil {
			return nil, err
		}
	}
	if err := t.authoriseMutation(ctx, id); err != nil {
		return nil, err
	}
	patch := robot_toolset.Update{Name: input.Name, Description: input.Description, Instruction: input.Instruction}
	if input.Tools != nil {
		patch.Tools = &input.Tools
	}
	updated, err := t.repo.Update(ctx, id, patch)
	if err != nil {
		return nil, err
	}
	return &mcp.ToolToolsetUpdateOutput{Id: updated.ID.String(), Name: updated.Name}, nil
}

func (t *toolsetTools) validateStandaloneTools(names []string) error {
	for _, name := range names {
		if err := t.registry.ValidateStandaloneTool(name); err != nil {
			return err
		}
	}
	return nil
}

func (t *toolsetTools) newDeleteTool() *Tool {
	def := mcp.GetToolsetDeleteTool()
	return &Tool{
		Definition: def,
		Builder: func(ctx context.Context) (tool.Tool, error) {
			return functiontool.New(functiontool.Config{
				Name: def.Name, Description: def.Description, InputSchema: def.InputSchema,
				RequireConfirmation: def.RequiresConfirmation && !confirmationDisabled(ctx),
			}, func(ctx adkagent.Context, input mcp.ToolToolsetDeleteInput) (*mcp.ToolToolsetDeleteOutput, error) {
				return t.delete(ctx, input)
			})
		},
		Handler: makeHandler(t.delete),
	}
}

func (t *toolsetTools) delete(ctx context.Context, input mcp.ToolToolsetDeleteInput) (*mcp.ToolToolsetDeleteOutput, error) {
	id, err := robot_toolset.ParseID(input.Id)
	if err != nil {
		return nil, fmt.Errorf("only custom Toolsets can be deleted: %w", err)
	}
	if err := t.authoriseMutation(ctx, id); err != nil {
		return nil, err
	}
	if err := t.repo.Delete(ctx, id); err != nil {
		return nil, err
	}
	return &mcp.ToolToolsetDeleteOutput{Id: input.Id, Deleted: true}, nil
}

func (t *toolsetTools) authoriseMutation(ctx context.Context, id robot_toolset.ID) error {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return err
	}
	return session.Authorise(ctx, func() error {
		isAuthor, err := t.repo.IsAuthor(ctx, id, accountID)
		if err != nil {
			return err
		}
		if !isAuthor {
			return fmt.Errorf("only the Toolset author can change it")
		}
		return nil
	}, rbac.PermissionManageRobots)
}

func (t *toolsetTools) canMutate(ctx context.Context, id robot_toolset.ID) bool {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return false
	}
	return session.Authorise(ctx, func() error {
		isAuthor, err := t.repo.IsAuthor(ctx, id, accountID)
		if err != nil {
			return err
		}
		if !isAuthor {
			return rbac.ErrPermissions
		}
		return nil
	}, rbac.PermissionManageRobots) == nil
}

func (t *toolsetTools) invalidTools(names []string) []string {
	var invalid []string
	for _, name := range names {
		if !t.registry.HasTool(name) {
			invalid = append(invalid, name)
		}
	}
	return invalid
}

func mapToolsetCatalogue(items []ToolsetInfo) []mcp.RobotToolsetCatalogueItemYaml {
	out := make([]mcp.RobotToolsetCatalogueItemYaml, 0, len(items))
	for _, item := range items {
		out = append(out, mcp.RobotToolsetCatalogueItemYaml{
			Id: item.ID, Name: item.Name, Description: item.Description,
		})
	}
	return out
}
