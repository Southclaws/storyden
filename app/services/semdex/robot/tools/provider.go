package tools

import (
	"context"
	"log/slog"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/puzpuzpuz/xsync/v4"
	"github.com/samber/lo"
	"go.uber.org/fx"
)

var ErrToolNotFound = fault.New("tool not found")

var DefaultTools = []string{
	"search",
	"system_all_tool_names",
	"robot_switch",
	"robot_create",
	"robot_list",
	"robot_get",
	"robot_update",
	"robot_delete",
	"throw_an_error",
}

type Registry struct {
	logger *slog.Logger
	tools  *xsync.Map[string, *Tool]
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(newRegistry),
		fx.Invoke(
			newRegistry,
			newRobotTools,
			newSearchTools,
			newLibraryTools,
			newTagTools,
			newLinkTools,
			newThreadTools,
		),
	)
}

func newRegistry(
	logger *slog.Logger,
) *Registry {
	return &Registry{
		logger: logger,
		tools:  xsync.NewMap[string, *Tool](),
	}
}

func (p *Registry) Register(tool *Tool) {
	p.tools.Store(tool.Definition.Name, tool)
}

func (p *Registry) GetTool(ctx context.Context, name string) (*Tool, error) {
	tool, ok := p.tools.Load(name)
	if !ok {
		return nil, fault.Wrap(ErrToolNotFound, fctx.With(ctx))
	}
	return tool, nil
}

func (p *Registry) GetTools(ctx context.Context, toolNames ...string) (Tools, error) {
	var tools []*Tool

	if len(toolNames) == 0 {
		p.tools.Range(func(key string, tool *Tool) bool {
			tools = append(tools, tool)
			return true
		})
		return tools, nil
	}

	ht := lo.KeyBy(toolNames, func(name string) string { return name })
	p.tools.Range(func(key string, tool *Tool) bool {
		if _, ok := ht[key]; ok {
			tools = append(tools, tool)
		}
		return true
	})

	return tools, nil
}
