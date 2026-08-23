package tools

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/puzpuzpuz/xsync/v4"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/semdex/robot/searchscore"
)

var (
	ErrToolNotFound          = fault.New("tool not found")
	ErrToolAlreadyRegistered = fault.New("tool already registered")
)

type ToolsetResolver interface {
	Validate(context.Context, []string) error
	RemoveContainedTools(context.Context, []string, []string) ([]string, error)
	ValidateDirectTools(context.Context, []string, []string) error
	SearchToolsets(context.Context, string, int) ([]ToolsetInfo, error)
	ListToolsets(context.Context) ([]ToolsetInfo, error)
	GetToolset(context.Context, string) (ToolsetInfo, error)
}

type ToolsetInfo struct {
	ID                string
	Name              string
	Description       string
	Instruction       string
	Tools             []string
	Editable          bool
	RequiresWorkspace bool
}

type Registry struct {
	logger  *slog.Logger
	tools   *xsync.Map[string, *Tool]
	aliases *xsync.Map[string, string]
}

type CatalogueTool struct {
	ID                   string
	CallableName         string
	Name                 string
	Description          string
	Source               string
	Available            bool
	RequiresConfirmation bool
	RequiresWorkspace    bool
	ToolsetOnly          bool
}

type ToolsetOnlyError struct {
	ToolName string
	Toolsets []string
}

func (e *ToolsetOnlyError) Error() string {
	if len(e.Toolsets) == 0 {
		return fmt.Sprintf("tool %q is Toolset-only and can only be provided by its declared Toolset", e.ToolName)
	}
	return fmt.Sprintf(
		"tool %q is Toolset-only and can only be provided by its declared Toolset; use Toolset %s instead",
		e.ToolName,
		strings.Join(e.Toolsets, " or "),
	)
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(NewRegistry),
		fx.Invoke(
			NewRegistry,
			newDocumentTools,
			newToolDiscoveryTools,
			newRobotTools,
			newToolsetTools,
			newSearchTools,
			newLibraryTools,
			newTagTools,
			newLinkTools,
			newThreadTools,
			newModerationTools,
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

func (p *Registry) Register(tool *Tool) error {
	if tool == nil || tool.Definition == nil {
		return fmt.Errorf("register Robot tool: definition is required")
	}
	if strings.TrimSpace(tool.Definition.Name) == "" {
		return fmt.Errorf("register Robot tool: ID is required")
	}
	if strings.TrimSpace(tool.Definition.Title) == "" {
		return fmt.Errorf("register Robot tool %q: title is required", tool.Definition.Name)
	}
	if _, ok := p.tools.LoadOrStore(tool.Definition.Name, tool); ok {
		return fault.Wrap(ErrToolAlreadyRegistered, fctx.With(context.Background()))
	}
	return nil
}

func (p *Registry) Unregister(name string) {
	p.tools.Delete(name)
}

func (p *Registry) RegisterAlias(alias string, target string) {
	p.aliases.Store(alias, target)
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

func (p *Registry) FindByADKName(name string) (*Tool, bool) {
	var found *Tool
	p.tools.Range(func(_ string, tool *Tool) bool {
		if tool.ADKName() == name || tool.Name() == name {
			found = tool
			return false
		}
		return true
	})
	return found, found != nil
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
		source := tool.Source
		if source == "" {
			source = "native"
		}
		if source == "native" && tool.CallableName != "" && tool.CallableName != tool.Definition.Name {
			source = "mcp"
		}
		tools = append(tools, CatalogueTool{
			ID:                   tool.Definition.Name,
			CallableName:         tool.ADKName(),
			Name:                 tool.Definition.Title,
			Description:          tool.Definition.Description,
			Source:               source,
			Available:            true,
			RequiresConfirmation: tool.Definition.RequiresConfirmation,
			RequiresWorkspace:    tool.Definition.RequiresWorkspace,
			ToolsetOnly:          tool.Definition.ToolsetOnly,
		})
		return true
	})
	p.aliases.Range(func(alias string, target string) bool {
		tool, ok := p.tools.Load(target)
		if !ok {
			return true
		}
		source := tool.Source
		if source == "" {
			source = "native"
		}
		tools = append(tools, CatalogueTool{
			ID:                   alias,
			CallableName:         tool.ADKName(),
			Name:                 tool.Definition.Title,
			Description:          tool.Definition.Description,
			Source:               source,
			Available:            true,
			RequiresConfirmation: tool.Definition.RequiresConfirmation,
			RequiresWorkspace:    tool.Definition.RequiresWorkspace,
			ToolsetOnly:          tool.Definition.ToolsetOnly,
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

// SearchCatalogue returns canonical tool identities for agent capability
// discovery. Aliases and management provenance remain available through the UI
// catalogue, but are intentionally excluded from this progressive-disclosure
// surface.
func (p *Registry) SearchCatalogue(query string, maxResults int) []CatalogueTool {
	if maxResults <= 0 || maxResults > 20 {
		maxResults = 8
	}

	parsedQuery := searchscore.New(query)
	type scoredTool struct {
		tool  CatalogueTool
		score int
	}

	results := make([]scoredTool, 0)
	p.tools.Range(func(_ string, tool *Tool) bool {
		if tool.Definition.ToolsetOnly {
			return true
		}
		name := tool.Definition.Title
		if strings.TrimSpace(name) == "" {
			name = tool.Definition.Name
		}
		item := CatalogueTool{
			ID:                   tool.Definition.Name,
			CallableName:         tool.ADKName(),
			Name:                 name,
			Description:          tool.Definition.Description,
			Available:            true,
			RequiresConfirmation: tool.Definition.RequiresConfirmation,
			RequiresWorkspace:    tool.Definition.RequiresWorkspace,
		}
		score := toolLexicalScore(parsedQuery, item)
		if parsedQuery.Empty() || score > 0 {
			results = append(results, scoredTool{tool: item, score: score})
		}
		return true
	})

	slices.SortFunc(results, func(a, b scoredTool) int {
		if a.score != b.score {
			return b.score - a.score
		}
		if a.tool.Name != b.tool.Name {
			return strings.Compare(a.tool.Name, b.tool.Name)
		}
		return strings.Compare(a.tool.ID, b.tool.ID)
	})
	if len(results) > maxResults {
		results = results[:maxResults]
	}

	tools := make([]CatalogueTool, 0, len(results))
	for _, result := range results {
		tools = append(tools, result.tool)
	}
	return tools
}

// ValidateStandaloneTool rejects capabilities that only make sense when their
// declared Toolset and its instruction are present.
func (p *Registry) ValidateStandaloneTool(name string) error {
	tool, ok := p.loadTool(name)
	if !ok || !tool.Definition.ToolsetOnly {
		return nil
	}
	return &ToolsetOnlyError{
		ToolName: tool.Definition.Name,
		Toolsets: slices.Clone(tool.Definition.Toolsets),
	}
}

func toolLexicalScore(query searchscore.Query, tool CatalogueTool) int {
	return query.Score(
		searchscore.Field{Text: tool.ID, ExactMatch: 120, SubstringMatch: 50, TermMatch: 10},
		searchscore.Field{Text: tool.Name, ExactMatch: 100, SubstringMatch: 40, TermMatch: 8},
		searchscore.Field{Text: tool.Description, SubstringMatch: 20, TermMatch: 4},
	)
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

func (p *Registry) CanonicalName(name string) string {
	return p.resolveName(name)
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
