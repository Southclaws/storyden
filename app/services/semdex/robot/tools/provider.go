package tools

import (
	"context"
	"log/slog"
	"slices"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/puzpuzpuz/xsync/v4"
	"go.uber.org/fx"
)

var ErrToolNotFound = fault.New("tool not found")

var DefaultTools = []string{
	"content_search",
	"system_robot_tool_catalog",
	"robot_switch",
	"robot_create",
	"robot_list",
	"robot_get",
	"robot_update",
	"robot_delete",
	"throw_an_error",
}

type Registry struct {
	logger  *slog.Logger
	tools   *xsync.Map[string, *Tool]
	aliases *xsync.Map[string, string]
}

type CatalogueTool struct {
	ID           string
	CallableName string
	Name         string
	Description  string
	Source       string
	Available    bool
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(NewRegistry),
		fx.Invoke(
			NewRegistry,
			newRobotTools,
			newSearchTools,
			newLibraryTools,
			newTagTools,
			newLinkTools,
			newThreadTools,
		),
	)
}

func NewRegistry(
	logger *slog.Logger,
) *Registry {
	return &Registry{
		logger:  logger,
		tools:   xsync.NewMap[string, *Tool](),
		aliases: xsync.NewMap[string, string](),
	}
}

func (p *Registry) Register(tool *Tool) {
	p.tools.Store(tool.Definition.Name, tool)
}

func (p *Registry) RegisterAlias(alias string, target string) {
	p.aliases.Store(alias, target)
}

func (p *Registry) Unregister(name string) {
	p.tools.Delete(name)
}

func (p *Registry) UnregisterPrefix(prefix string) {
	p.tools.Range(func(key string, tool *Tool) bool {
		if strings.HasPrefix(key, prefix) {
			p.tools.Delete(key)
		}
		return true
	})
}

func (p *Registry) HasTool(name string) bool {
	_, ok := p.loadTool(name)
	return ok
}

func (p *Registry) GetTool(ctx context.Context, name string) (*Tool, error) {
	tool, ok := p.loadTool(name)
	if !ok {
		return nil, fault.Wrap(ErrToolNotFound, fctx.With(ctx))
	}
	return tool, nil
}

func (p *Registry) GetTools(ctx context.Context, toolNames ...string) (Tools, error) {
	tools, _ := p.GetToolsWithMissing(ctx, toolNames...)
	return tools, nil
}

func (p *Registry) GetToolsWithMissing(ctx context.Context, toolNames ...string) (Tools, []string) {
	var tools []*Tool
	var missing []string

	if len(toolNames) == 0 {
		p.tools.Range(func(key string, tool *Tool) bool {
			tools = append(tools, tool)
			return true
		})
		slices.SortFunc(tools, func(a, b *Tool) int {
			if a.Name() < b.Name() {
				return -1
			}
			if a.Name() > b.Name() {
				return 1
			}
			return 0
		})
		return tools, nil
	}

	seenRequested := map[string]struct{}{}
	seenResolved := map[string]struct{}{}
	for _, name := range toolNames {
		if _, ok := seenRequested[name]; ok {
			continue
		}
		seenRequested[name] = struct{}{}
		resolved := p.resolveName(name)
		tool, ok := p.tools.Load(resolved)
		if !ok {
			missing = append(missing, name)
			continue
		}
		if _, ok := seenResolved[resolved]; ok {
			continue
		}
		seenResolved[resolved] = struct{}{}
		tools = append(tools, tool)
	}

	return tools, missing
}

func (p *Registry) ListCatalogue(ctx context.Context) []CatalogueTool {
	var tools []CatalogueTool
	p.tools.Range(func(key string, tool *Tool) bool {
		source := "native"
		if tool.CallableName != "" && tool.CallableName != tool.Definition.Name {
			source = "mcp"
		}
		tools = append(tools, CatalogueTool{
			ID:           tool.Definition.Name,
			CallableName: tool.ADKName(),
			Name:         tool.Definition.Title,
			Description:  tool.Definition.Description,
			Source:       source,
			Available:    true,
		})
		return true
	})
	p.aliases.Range(func(alias string, target string) bool {
		tool, ok := p.tools.Load(target)
		if !ok {
			return true
		}
		tools = append(tools, CatalogueTool{
			ID:           alias,
			CallableName: tool.ADKName(),
			Name:         tool.Definition.Title,
			Description:  tool.Definition.Description,
			Source:       "native",
			Available:    true,
		})
		return true
	})
	slices.SortFunc(tools, func(a, b CatalogueTool) int {
		if a.ID < b.ID {
			return -1
		}
		if a.ID > b.ID {
			return 1
		}
		return 0
	})
	return tools
}

func (p *Registry) AllToolIDs(ctx context.Context) []string {
	var names []string
	p.tools.Range(func(key string, tool *Tool) bool {
		names = append(names, key)
		return true
	})
	p.aliases.Range(func(alias string, target string) bool {
		if _, ok := p.tools.Load(target); ok {
			names = append(names, alias)
		}
		return true
	})
	slices.Sort(names)
	return names
}

func (p *Registry) resolveName(name string) string {
	if target, ok := p.aliases.Load(name); ok {
		return target
	}
	return name
}

func (p *Registry) loadTool(name string) (*Tool, bool) {
	return p.tools.Load(p.resolveName(name))
}
