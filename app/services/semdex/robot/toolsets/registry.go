package toolsets

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sort"
	"strings"
	"sync"

	adkagent "google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/robot/robot_toolset"
	"github.com/Southclaws/storyden/app/services/semdex/robot/searchscore"
	robottools "github.com/Southclaws/storyden/app/services/semdex/robot/tools"
)

type Source string

const (
	SourceSystem Source = "system"
	SourceCustom Source = "custom"
	SourcePlugin Source = "plugin"
)

type Builder func(context.Context) (tool.Toolset, error)

type Definition struct {
	ID                  string
	Name                string
	Description         string
	RequiresWorkspace   bool
	Instruction         string
	InstructionProvider func(adkagent.ReadonlyContext) (string, error)
	ToolNames           []string
	Source              Source
	SourceRef           string
	Editable            bool
	UsageCount          int
	Author              *account.Account
	Metadata            map[string]any
	Builder             Builder
}

type DirectToolConflictError struct {
	ToolName string
	Toolset  Definition
}

func (e *DirectToolConflictError) Error() string {
	return fmt.Sprintf(
		"tool %q cannot be added directly because Toolset %q (%s) already contains that tool",
		e.ToolName,
		e.Toolset.Name,
		e.Toolset.ID,
	)
}

type Registry struct {
	logger *slog.Logger
	repo   *robot_toolset.Repository
	tools  *robottools.Registry

	mu      sync.RWMutex
	dynamic map[string]Definition
}

func NewRegistry(logger *slog.Logger, repo *robot_toolset.Repository, tools *robottools.Registry) *Registry {
	return &Registry{
		logger:  logger,
		repo:    repo,
		tools:   tools,
		dynamic: make(map[string]Definition),
	}
}

func (r *Registry) Register(def Definition) error {
	def.ID = strings.TrimSpace(def.ID)
	def.Name = strings.TrimSpace(def.Name)
	if def.ID == "" || def.Name == "" {
		return fmt.Errorf("toolset ID and name are required")
	}
	if def.Source == "" {
		def.Source = SourceSystem
	}
	def.Editable = false
	def.ToolNames = uniqueStrings(def.ToolNames)

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.dynamic[def.ID]; exists {
		return fmt.Errorf("toolset %q is already registered", def.ID)
	}
	r.dynamic[def.ID] = def
	r.logger.Debug("registered Robot Toolset", slog.String("toolset_id", def.ID), slog.String("name", def.Name))
	return nil
}

func (r *Registry) Unregister(id string) {
	r.mu.Lock()
	delete(r.dynamic, id)
	r.mu.Unlock()
}

func (r *Registry) UnregisterSource(sourceRef string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, def := range r.dynamic {
		if def.SourceRef == sourceRef {
			delete(r.dynamic, id)
		}
	}
}

func (r *Registry) Get(ctx context.Context, id string) (Definition, error) {
	r.mu.RLock()
	def, ok := r.dynamic[id]
	r.mu.RUnlock()
	if ok {
		return r.withWorkspaceRequirement(ctx, def), nil
	}

	parsed, err := robot_toolset.ParseID(id)
	if err != nil {
		return Definition{}, fmt.Errorf("unknown toolset %q", id)
	}
	stored, err := r.repo.Get(ctx, parsed)
	if err != nil {
		return Definition{}, err
	}
	usage, err := r.repo.UsageCount(ctx, parsed)
	if err != nil {
		return Definition{}, err
	}
	return r.withWorkspaceRequirement(ctx, fromStored(stored, usage)), nil
}

func (r *Registry) List(ctx context.Context) ([]Definition, error) {
	r.mu.RLock()
	definitions := make([]Definition, 0, len(r.dynamic))
	for _, def := range r.dynamic {
		definitions = append(definitions, r.withWorkspaceRequirement(ctx, def))
	}
	r.mu.RUnlock()

	stored, err := r.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range stored {
		usage, err := r.repo.UsageCount(ctx, item.ID)
		if err != nil {
			return nil, err
		}
		definitions = append(definitions, r.withWorkspaceRequirement(ctx, fromStored(item, usage)))
	}

	sort.Slice(definitions, func(i, j int) bool {
		if definitions[i].Source != definitions[j].Source {
			return definitions[i].Source < definitions[j].Source
		}
		if definitions[i].Name != definitions[j].Name {
			return definitions[i].Name < definitions[j].Name
		}
		return definitions[i].ID < definitions[j].ID
	})
	return definitions, nil
}

func (r *Registry) Search(ctx context.Context, query string, maxResults int) ([]Definition, error) {
	all, err := r.List(ctx)
	if err != nil {
		return nil, err
	}
	if maxResults <= 0 || maxResults > 20 {
		maxResults = 8
	}
	parsedQuery := searchscore.New(query)
	if parsedQuery.Empty() {
		if len(all) > maxResults {
			all = all[:maxResults]
		}
		return all, nil
	}

	type scored struct {
		def   Definition
		score int
	}
	results := make([]scored, 0, len(all))
	for _, def := range all {
		score := lexicalScore(parsedQuery, def)
		if score > 0 {
			results = append(results, scored{def: def, score: score})
		}
	}
	sort.SliceStable(results, func(i, j int) bool {
		if results[i].score != results[j].score {
			return results[i].score > results[j].score
		}
		return results[i].def.Name < results[j].def.Name
	})
	if len(results) > maxResults {
		results = results[:maxResults]
	}
	definitions := make([]Definition, 0, len(results))
	for _, result := range results {
		definitions = append(definitions, result.def)
	}
	return definitions, nil
}

func (r *Registry) Validate(ctx context.Context, refs []string) error {
	for _, ref := range uniqueStrings(refs) {
		if _, err := r.Get(ctx, ref); err != nil {
			return err
		}
	}
	return nil
}

// RemoveContainedTools applies Toolset precedence to a Robot's direct tool
// assignment. Tool order is preserved so configuration responses remain stable.
func (r *Registry) RemoveContainedTools(ctx context.Context, toolNames, refs []string) ([]string, error) {
	contained := make(map[string]struct{})
	for _, ref := range uniqueStrings(refs) {
		def, err := r.Get(ctx, ref)
		if err != nil {
			return nil, err
		}
		for _, name := range def.ToolNames {
			contained[r.tools.CanonicalName(name)] = struct{}{}
		}
	}

	result := make([]string, 0, len(toolNames))
	for _, name := range toolNames {
		if _, found := contained[r.tools.CanonicalName(name)]; found {
			continue
		}
		result = append(result, name)
	}
	return result, nil
}

// ValidateDirectTools rejects direct tools already provided by an assigned
// Toolset and returns an error that names the owning Toolset.
func (r *Registry) ValidateDirectTools(ctx context.Context, toolNames, refs []string) error {
	definitions := make([]Definition, 0, len(refs))
	for _, ref := range uniqueStrings(refs) {
		def, err := r.Get(ctx, ref)
		if err != nil {
			return err
		}
		definitions = append(definitions, def)
	}

	for _, name := range toolNames {
		for _, def := range definitions {
			if slices.ContainsFunc(def.ToolNames, func(toolsetTool string) bool {
				return r.tools.CanonicalName(toolsetTool) == r.tools.CanonicalName(name)
			}) {
				return &DirectToolConflictError{ToolName: name, Toolset: def}
			}
		}
	}
	return nil
}

func (r *Registry) ListToolsets(ctx context.Context) ([]robottools.ToolsetInfo, error) {
	definitions, err := r.List(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]robottools.ToolsetInfo, 0, len(definitions))
	for _, def := range definitions {
		items = append(items, toolsetInfo(def))
	}
	return items, nil
}

func (r *Registry) SearchToolsets(ctx context.Context, query string, maxResults int) ([]robottools.ToolsetInfo, error) {
	definitions, err := r.Search(ctx, query, maxResults)
	if err != nil {
		return nil, err
	}
	items := make([]robottools.ToolsetInfo, 0, len(definitions))
	for _, def := range definitions {
		items = append(items, toolsetInfo(def))
	}
	return items, nil
}

func (r *Registry) GetToolset(ctx context.Context, id string) (robottools.ToolsetInfo, error) {
	def, err := r.Get(ctx, id)
	if err != nil {
		return robottools.ToolsetInfo{}, err
	}
	if def.InstructionProvider != nil {
		readonly, ok := ctx.(adkagent.ReadonlyContext)
		if !ok {
			return robottools.ToolsetInfo{}, fmt.Errorf("Toolset %q requires an active Robot conversation to resolve its instruction", id)
		}
		instruction, err := def.InstructionProvider(readonly)
		if err != nil {
			return robottools.ToolsetInfo{}, err
		}
		def.Instruction = instruction
	}
	return toolsetInfo(def), nil
}

func toolsetInfo(def Definition) robottools.ToolsetInfo {
	return robottools.ToolsetInfo{
		ID:                def.ID,
		Name:              def.Name,
		Description:       def.Description,
		Instruction:       def.Instruction,
		Tools:             slices.Clone(def.ToolNames),
		Editable:          def.Editable,
		RequiresWorkspace: def.RequiresWorkspace,
	}
}

func (r *Registry) withWorkspaceRequirement(ctx context.Context, def Definition) Definition {
	if def.RequiresWorkspace {
		return def
	}
	for _, name := range def.ToolNames {
		registered, err := r.tools.GetTool(ctx, name)
		if err == nil && registered.Definition.RequiresWorkspace {
			def.RequiresWorkspace = true
			break
		}
	}
	return def
}

func (r *Registry) resolveTools(ctx context.Context, def Definition) ([]tool.Tool, []string, error) {
	registered, missing := r.tools.GetToolsWithMissing(ctx, def.ToolNames...)
	resolved, err := registered.ToADKTools(ctx)
	return resolved, missing, err
}

func fromStored(stored *robot_toolset.Toolset, usage int) Definition {
	return Definition{
		ID:          stored.ID.String(),
		Name:        stored.Name,
		Description: stored.Description,
		Instruction: stored.Instruction,
		ToolNames:   append([]string(nil), stored.Tools...),
		Source:      SourceCustom,
		SourceRef:   stored.Author.ID.String(),
		Editable:    true,
		UsageCount:  usage,
		Author:      &stored.Author,
		Metadata:    stored.Metadata,
	}
}

func lexicalScore(query searchscore.Query, def Definition) int {
	fields := []searchscore.Field{
		{Text: def.Name, ExactMatch: 100, SubstringMatch: 40, TermMatch: 8},
		{Text: def.Description, SubstringMatch: 20, TermMatch: 4},
		{Text: def.Instruction, SubstringMatch: 5},
	}
	for _, toolName := range def.ToolNames {
		fields = append(fields, searchscore.Field{Text: toolName, SubstringMatch: 10})
	}
	return query.Score(fields...)
}

func uniqueStrings(values []string) []string {
	values = slices.Clone(values)
	slices.Sort(values)
	return slices.Compact(values)
}
