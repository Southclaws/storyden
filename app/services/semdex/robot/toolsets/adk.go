package toolsets

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
	"google.golang.org/genai"

	robottools "github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	"github.com/Southclaws/storyden/app/services/semdex/robot/workspacestate"
	"github.com/Southclaws/storyden/lib/mcp"
)

const (
	activeToolsetsStateKey = "active_toolsets"
	activeToolsStateKey    = "active_tools"
	runtimeIssuesToolName  = "tool_runtime_issues"
)

type promptedToolset struct {
	definition Definition
	delegate   tool.Toolset
	builder    Builder
	buildCtx   context.Context
	registered robottools.Tools
	tools      []tool.Tool
	excluded   map[string]struct{}
}

func (t *promptedToolset) Name() string { return t.definition.ID }

func (t *promptedToolset) Tools(ctx agent.ReadonlyContext) ([]tool.Tool, error) {
	if t.definition.RequiresWorkspace && !workspacestate.Available(ctx.ReadonlyState()) {
		return nil, nil
	}
	resolved := t.tools
	if len(t.registered) > 0 {
		var err error
		resolved, err = t.registered.ToADKTools(t.buildCtx)
		if err != nil {
			return nil, err
		}
	}
	delegate := t.delegate
	if delegate == nil && t.builder != nil {
		var err error
		delegate, err = t.builder(t.buildCtx)
		if err != nil {
			return nil, fmt.Errorf("build Toolset %q: %w", t.definition.Name, err)
		}
	}
	if delegate != nil {
		var err error
		resolved, err = delegate.Tools(ctx)
		if err != nil {
			return nil, err
		}
	}
	if len(t.excluded) == 0 {
		return resolved, nil
	}

	filtered := make([]tool.Tool, 0, len(resolved))
	for _, item := range resolved {
		if _, duplicate := t.excluded[item.Name()]; duplicate {
			continue
		}
		filtered = append(filtered, item)
	}
	return filtered, nil
}

func (t *promptedToolset) ProcessRequest(ctx agent.Context, request *model.LLMRequest) error {
	if t.definition.RequiresWorkspace && !workspacestate.Available(ctx.ReadonlyState()) {
		appendInstruction(request, workspaceRequiredInstruction("Toolset", t.definition.Name, t.definition.ID))
		return nil
	}
	instruction, err := formatToolsetInstruction(ctx, t.definition)
	if err != nil {
		return err
	}
	appendInstruction(request, instruction)
	return nil
}

func (r *Registry) Build(ctx context.Context, refs []string) ([]tool.Toolset, []string, error) {
	toolsets := make([]tool.Toolset, 0, len(refs))
	var missing []string
	claimedTools := make(map[string]struct{})
	for _, ref := range uniqueStrings(refs) {
		def, err := r.Get(ctx, ref)
		if err != nil {
			return nil, nil, err
		}
		if def.Builder != nil {
			var delegate tool.Toolset
			if !def.RequiresWorkspace {
				delegate, err = def.Builder(ctx)
				if err != nil {
					return nil, nil, fmt.Errorf("build Toolset %q: %w", def.Name, err)
				}
			}
			excluded := make(map[string]struct{})
			for _, toolName := range def.ToolNames {
				if _, duplicate := claimedTools[toolName]; duplicate {
					excluded[toolName] = struct{}{}
					continue
				}
				claimedTools[toolName] = struct{}{}
			}
			toolsets = append(toolsets, &promptedToolset{
				definition: def,
				delegate:   delegate,
				builder:    def.Builder,
				buildCtx:   ctx,
				excluded:   excluded,
			})
			continue
		}
		if !def.RequiresWorkspace {
			resolved, unavailable, err := r.resolveTools(ctx, def)
			if err != nil {
				return nil, nil, fmt.Errorf("build Toolset %q: %w", def.Name, err)
			}
			missing = append(missing, unavailable...)
			filtered := make([]tool.Tool, 0, len(resolved))
			for _, item := range resolved {
				if _, duplicate := claimedTools[item.Name()]; duplicate {
					continue
				}
				claimedTools[item.Name()] = struct{}{}
				filtered = append(filtered, item)
			}
			toolsets = append(toolsets, &promptedToolset{definition: def, tools: filtered})
			continue
		}
		registered, unavailable := r.tools.GetToolsWithMissing(ctx, def.ToolNames...)
		missing = append(missing, unavailable...)
		filtered := make(robottools.Tools, 0, len(registered))
		for _, item := range registered {
			if _, duplicate := claimedTools[item.ADKName()]; duplicate {
				continue
			}
			claimedTools[item.ADKName()] = struct{}{}
			filtered = append(filtered, item)
		}
		toolsets = append(toolsets, &promptedToolset{definition: def, registered: filtered, buildCtx: ctx})
	}
	return toolsets, uniqueStrings(missing), nil
}

type discoveryToolset struct {
	registry *Registry
	buildCtx context.Context
	baseRefs []string
	tools    []tool.Tool

	mu         sync.RWMutex
	toolOwners map[string]map[string]struct{}
	managed    map[string]struct{}
	issues     []runtimeIssue
}

type runtimeIssue struct {
	CapabilityType  string `json:"capability_type"`
	ID              string `json:"id"`
	Name            string `json:"name"`
	Error           string `json:"error"`
	SuggestedAction string `json:"suggested_action"`
}

type runtimeIssuesOutput struct {
	Issues     []runtimeIssue `json:"issues"`
	NextAction string         `json:"next_action"`
}

type SearchInput struct {
	Query      string `json:"query" jsonschema:"What capability or task you need tools for"`
	MaxResults int    `json:"max_results,omitempty" jsonschema:"Maximum Toolsets to return, from 1 to 20"`
}

type CatalogueItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SearchOutput struct {
	Toolsets   []CatalogueItem `json:"toolsets"`
	NextAction string          `json:"next_action"`
}

type LoadInput struct {
	Toolsets []string `json:"toolsets" jsonschema:"Toolset IDs inspected with toolset_get"`
}

type LoadOutput struct {
	Loaded     []CatalogueItem `json:"loaded"`
	NextAction string          `json:"next_action"`
}

func (r *Registry) Discovery(ctx context.Context, baseRefs ...string) (tool.Toolset, error) {
	discovery := &discoveryToolset{
		registry:   r,
		buildCtx:   ctx,
		baseRefs:   uniqueStrings(baseRefs),
		toolOwners: make(map[string]map[string]struct{}),
		managed:    make(map[string]struct{}),
	}

	search, err := functiontool.New(functiontool.Config{
		Name:        "toolset_search",
		Description: "Search reusable Toolsets by the capability or task you need. Use this before loading tools that are not already available.",
	}, func(ctx agent.Context, input SearchInput) (SearchOutput, error) {
		definitions, err := r.Search(ctx, input.Query, input.MaxResults)
		if err != nil {
			return SearchOutput{}, err
		}
		return SearchOutput{
			Toolsets:   catalogueItems(definitions),
			NextAction: "Use toolset_get with a returned ID to inspect its tools and instruction before loading or assigning it.",
		}, nil
	})
	if err != nil {
		return nil, err
	}

	load, err := functiontool.New(functiontool.Config{
		Name:        "toolset_load",
		Description: "Activate one or more Toolsets for this conversation after inspecting their contents and runtime preconditions with toolset_get. Workspace-dependent Toolsets require an active workspace. Their tools and specialised prompt guidance become available on the next step.",
	}, func(ctx agent.Context, input LoadInput) (LoadOutput, error) {
		refs := uniqueStrings(input.Toolsets)
		if len(refs) == 0 {
			return LoadOutput{}, fmt.Errorf("at least one Toolset ID is required")
		}
		definitions := make([]Definition, 0, len(refs))
		for _, ref := range refs {
			def, err := r.Get(ctx, ref)
			if err != nil {
				return LoadOutput{}, err
			}
			if def.RequiresWorkspace && !workspacestate.Available(ctx.ReadonlyState()) {
				return LoadOutput{}, fmt.Errorf("%s", workspaceRequiredInstruction("Toolset", def.Name, def.ID))
			}
			if issue, ok := discovery.issueFor("toolset", def.ID); ok {
				return LoadOutput{}, runtimeIssueError(issue)
			}
			definitions = append(definitions, def)
		}
		active := append(activeToolsetRefs(ctx.ReadonlyState()), refs...)
		active = uniqueStrings(active)
		if err := ctx.State().Set(activeToolsetsStateKey, active); err != nil {
			return LoadOutput{}, err
		}
		return LoadOutput{
			Loaded:     catalogueItems(definitions),
			NextAction: "Continue the task using the newly available tools. Do not call toolset_load again for these Toolsets.",
		}, nil
	})
	if err != nil {
		return nil, err
	}

	loadTools, err := functiontool.New(functiontool.Config{
		Name:        "tool_load",
		Description: "Activate one or more individual tools for this conversation after inspecting their schemas and runtime preconditions with tool_get. Workspace-dependent tools require an active workspace. The tools become available on the next model step.",
	}, func(ctx agent.Context, input mcp.ToolToolLoadInput) (mcp.ToolToolLoadOutput, error) {
		refs := uniqueStrings(input.Tools)
		if len(refs) == 0 {
			return mcp.ToolToolLoadOutput{}, fmt.Errorf("at least one tool ID is required")
		}

		loaded := make([]mcp.RobotToolCatalogueItemYaml, 0, len(refs))
		canonical := make([]string, 0, len(refs))
		for _, ref := range refs {
			selected, err := r.tools.GetTool(ctx, ref)
			if err != nil {
				return mcp.ToolToolLoadOutput{}, err
			}
			if selected.Definition.RequiresWorkspace && !workspacestate.Available(ctx.ReadonlyState()) {
				name := selected.Definition.Title
				if strings.TrimSpace(name) == "" {
					name = selected.Definition.Name
				}
				return mcp.ToolToolLoadOutput{}, fmt.Errorf("%s", workspaceRequiredInstruction("Tool", name, selected.Definition.Name))
			}
			if issue, ok := discovery.issueFor("tool", selected.Definition.Name); ok {
				return mcp.ToolToolLoadOutput{}, runtimeIssueError(issue)
			}
			name := selected.Definition.Title
			if strings.TrimSpace(name) == "" {
				name = selected.Definition.Name
			}
			canonical = append(canonical, selected.Definition.Name)
			loaded = append(loaded, mcp.RobotToolCatalogueItemYaml{
				Id: selected.Definition.Name, Name: name, Description: selected.Definition.Description,
			})
		}

		active := append(activeToolRefs(ctx.ReadonlyState()), canonical...)
		active = uniqueStrings(active)
		if err := ctx.State().Set(activeToolsStateKey, active); err != nil {
			return mcp.ToolToolLoadOutput{}, err
		}
		return mcp.ToolToolLoadOutput{
			Loaded:     loaded,
			NextAction: "Continue the task using the newly available tools. Do not call tool_load again for these tools.",
		}, nil
	})
	if err != nil {
		return nil, err
	}

	issues, err := functiontool.New(functiontool.Config{
		Name:        runtimeIssuesToolName,
		Description: "Inspect registered Robot capabilities that could not be initialized for this conversation. This tool is only shown when runtime issues exist.",
	}, func(ctx agent.Context, _ struct{}) (runtimeIssuesOutput, error) {
		return runtimeIssuesOutput{
			Issues:     discovery.runtimeIssues(),
			NextAction: "Continue with healthy capabilities. Tell an administrator which capability failed and follow its suggested action; a broken plugin can be disabled while it is repaired.",
		}, nil
	})
	if err != nil {
		return nil, err
	}

	discovery.tools = []tool.Tool{search, load, loadTools, issues}
	return discovery, nil
}

func (t *discoveryToolset) Name() string { return "toolset_discovery" }

func (t *discoveryToolset) Tools(ctx agent.ReadonlyContext) ([]tool.Tool, error) {
	available := slices.Clone(t.tools)
	definitions, err := t.registry.List(ctx)
	if err != nil {
		return nil, err
	}
	owners := make(map[string]map[string]struct{})
	managed := make(map[string]struct{})
	issues := make([]runtimeIssue, 0)
	issueKeys := make(map[string]struct{})
	addIssue := func(issue runtimeIssue) {
		key := issue.CapabilityType + "\x00" + issue.ID + "\x00" + issue.Error
		if _, exists := issueKeys[key]; exists {
			return
		}
		issueKeys[key] = struct{}{}
		issues = append(issues, issue)
	}
	seen := make(map[string]struct{}, len(available))
	for _, item := range available {
		seen[item.Name()] = struct{}{}
	}
	for _, def := range definitions {
		if def.RequiresWorkspace && !workspacestate.Available(ctx.ReadonlyState()) {
			continue
		}
		if def.Builder != nil {
			built, err := def.Builder(t.buildCtx)
			if err != nil {
				addIssue(runtimeIssueForToolset(def, fmt.Errorf("initialize Toolset: %w", err)))
				continue
			}
			if built == nil {
				addIssue(runtimeIssueForToolset(def, fmt.Errorf("initialize Toolset: builder returned no Toolset")))
				continue
			}
			resolved, err := built.Tools(ctx)
			if err != nil {
				addIssue(runtimeIssueForToolset(def, fmt.Errorf("list Toolset tools: %w", err)))
				continue
			}
			for _, item := range resolved {
				managed[item.Name()] = struct{}{}
				addToolOwner(owners, item.Name(), def.ID)
				if _, exists := seen[item.Name()]; exists {
					continue
				}
				seen[item.Name()] = struct{}{}
				available = append(available, item)
			}
		} else {
			registered, missing := t.registry.tools.GetToolsWithMissing(t.buildCtx, def.ToolNames...)
			for _, missingID := range missing {
				addIssue(runtimeIssue{
					CapabilityType:  "tool",
					ID:              missingID,
					Name:            missingID,
					Error:           fmt.Sprintf("Toolset %q (%s) references a tool that is no longer registered", def.Name, def.ID),
					SuggestedAction: unavailableReferenceAction(def),
				})
			}
			for _, item := range registered {
				addToolOwner(owners, item.ADKName(), def.ID)
			}
		}
	}

	registered, err := t.registry.tools.GetTools(t.buildCtx)
	if err != nil {
		return nil, err
	}
	for _, item := range registered {
		if item.Definition.RequiresWorkspace && !workspacestate.Available(ctx.ReadonlyState()) {
			continue
		}
		if item.Builder == nil {
			addIssue(runtimeIssueForTool(item, fmt.Errorf("tool has no runtime builder")))
			continue
		}
		resolved, err := item.Builder(t.buildCtx)
		if err != nil {
			addIssue(runtimeIssueForTool(item, fmt.Errorf("initialize tool: %w", err)))
			continue
		}
		if resolved == nil {
			addIssue(runtimeIssueForTool(item, fmt.Errorf("initialize tool: builder returned no tool")))
			continue
		}
		if resolved.Name() != item.ADKName() {
			addIssue(runtimeIssueForTool(item, fmt.Errorf("runtime name %q does not match registered callable name %q", resolved.Name(), item.ADKName())))
			continue
		}
		managed[resolved.Name()] = struct{}{}
		if _, exists := seen[resolved.Name()]; exists {
			continue
		}
		seen[resolved.Name()] = struct{}{}
		available = append(available, resolved)
	}
	slices.SortFunc(issues, func(a, b runtimeIssue) int {
		if a.CapabilityType != b.CapabilityType {
			return strings.Compare(a.CapabilityType, b.CapabilityType)
		}
		return strings.Compare(a.ID, b.ID)
	})
	t.mu.Lock()
	t.toolOwners = owners
	t.managed = managed
	t.issues = issues
	t.mu.Unlock()
	return available, nil
}

func (t *discoveryToolset) ProcessRequest(ctx agent.Context, request *model.LLMRequest) error {
	appendInstruction(request, discoveryInstruction())

	refs := append(slices.Clone(t.baseRefs), activeToolsetRefs(ctx.ReadonlyState())...)
	refs = uniqueStrings(refs)
	active := make(map[string]struct{}, len(refs))
	for _, ref := range refs {
		def, err := t.registry.Get(ctx, ref)
		if err != nil {
			t.addRuntimeIssue(runtimeIssue{
				CapabilityType:  "toolset",
				ID:              ref,
				Name:            ref,
				Error:           "configured Toolset is no longer registered",
				SuggestedAction: "Continue without this Toolset and update the Robot configuration when convenient.",
			})
			continue
		}
		if def.RequiresWorkspace && !workspacestate.Available(ctx.ReadonlyState()) {
			appendInstruction(request, workspaceRequiredInstruction("Toolset", def.Name, def.ID))
			continue
		}
		if _, unavailable := t.issueFor("toolset", def.ID); unavailable {
			continue
		}
		instruction, err := formatToolsetInstruction(ctx, def)
		if err != nil {
			t.addRuntimeIssue(runtimeIssueForToolset(def, fmt.Errorf("resolve Toolset instruction: %w", err)))
			continue
		}
		active[ref] = struct{}{}
		appendInstruction(request, instruction)
	}
	activeTools := make(map[string]struct{})
	for _, ref := range activeToolRefs(ctx.ReadonlyState()) {
		selected, err := t.registry.tools.GetTool(ctx, ref)
		if err != nil {
			t.addRuntimeIssue(runtimeIssue{
				CapabilityType:  "tool",
				ID:              ref,
				Name:            ref,
				Error:           "loaded tool is no longer registered",
				SuggestedAction: "Continue without this tool and update the Robot configuration when convenient.",
			})
			continue
		}
		if _, unavailable := t.issueFor("tool", selected.Definition.Name); unavailable {
			continue
		}
		activeTools[selected.ADKName()] = struct{}{}
	}

	t.mu.RLock()
	for toolName := range t.managed {
		_, allowed := activeTools[toolName]
		owners := t.toolOwners[toolName]
		for owner := range owners {
			if _, ok := active[owner]; ok {
				allowed = true
				break
			}
		}
		if !allowed {
			delete(request.Tools, toolName)
		}
	}
	t.mu.RUnlock()

	issues := t.runtimeIssues()
	if len(issues) == 0 {
		delete(request.Tools, runtimeIssuesToolName)
	} else {
		appendInstruction(request, fmt.Sprintf(
			"%d registered Robot capabilities are unavailable in this conversation. Continue with healthy capabilities. Use %s to inspect the failures and recovery actions before explaining them to the user.",
			len(issues),
			runtimeIssuesToolName,
		))
	}
	return nil
}

func addToolOwner(owners map[string]map[string]struct{}, toolName, toolsetID string) {
	if owners[toolName] == nil {
		owners[toolName] = make(map[string]struct{})
	}
	owners[toolName][toolsetID] = struct{}{}
}

func (t *discoveryToolset) runtimeIssues() []runtimeIssue {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return slices.Clone(t.issues)
}

func (t *discoveryToolset) issueFor(capabilityType, id string) (runtimeIssue, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	for _, issue := range t.issues {
		if issue.CapabilityType == capabilityType && issue.ID == id {
			return issue, true
		}
	}
	return runtimeIssue{}, false
}

func (t *discoveryToolset) addRuntimeIssue(issue runtimeIssue) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, existing := range t.issues {
		if existing.CapabilityType == issue.CapabilityType && existing.ID == issue.ID && existing.Error == issue.Error {
			return
		}
	}
	t.issues = append(t.issues, issue)
	slices.SortFunc(t.issues, func(a, b runtimeIssue) int {
		if a.CapabilityType != b.CapabilityType {
			return strings.Compare(a.CapabilityType, b.CapabilityType)
		}
		return strings.Compare(a.ID, b.ID)
	})
}

func runtimeIssueForToolset(def Definition, err error) runtimeIssue {
	return runtimeIssue{
		CapabilityType:  "toolset",
		ID:              def.ID,
		Name:            def.Name,
		Error:           err.Error(),
		SuggestedAction: unavailableReferenceAction(def),
	}
}

func runtimeIssueForTool(item *robottools.Tool, err error) runtimeIssue {
	name := item.Definition.Title
	if strings.TrimSpace(name) == "" {
		name = item.Definition.Name
	}
	action := "Continue without this tool and report the failure to an administrator."
	if item.Source == "plugin" {
		action = "Disable or repair the plugin that provides this tool, then retry in a new turn."
	} else if item.Source == "mcp" {
		action = "Disable or repair the MCP server that provides this tool, then retry in a new turn."
	}
	return runtimeIssue{
		CapabilityType:  "tool",
		ID:              item.Definition.Name,
		Name:            name,
		Error:           err.Error(),
		SuggestedAction: action,
	}
}

func unavailableReferenceAction(def Definition) string {
	switch def.Source {
	case SourcePlugin:
		if def.SourceRef != "" {
			return fmt.Sprintf("Disable or repair plugin installation %s, then retry in a new turn.", def.SourceRef)
		}
		return "Disable or repair the plugin that provides this Toolset, then retry in a new turn."
	case SourceCustom:
		return "Update the custom Toolset to remove or replace unavailable tools, then retry in a new turn."
	default:
		return "Continue without this capability and report the failure to an administrator."
	}
}

func runtimeIssueError(issue runtimeIssue) error {
	return fmt.Errorf("%s %q (%s) is currently unavailable: %s. %s", issue.CapabilityType, issue.Name, issue.ID, issue.Error, issue.SuggestedAction)
}

func activeToolsetRefs(state interface{ Get(string) (any, error) }) []string {
	return stringStateValues(state, activeToolsetsStateKey)
}

func activeToolRefs(state interface{ Get(string) (any, error) }) []string {
	return stringStateValues(state, activeToolsStateKey)
}

func stringStateValues(state interface{ Get(string) (any, error) }, key string) []string {
	value, err := state.Get(key)
	if err != nil {
		return nil
	}
	switch refs := value.(type) {
	case []string:
		return slices.Clone(refs)
	case []any:
		out := make([]string, 0, len(refs))
		for _, value := range refs {
			if ref, ok := value.(string); ok {
				out = append(out, ref)
			}
		}
		return out
	default:
		return nil
	}
}

func discoveryInstruction() string {
	return "Additional capabilities use progressive disclosure. Search Toolsets when a task needs a coherent capability bundle, then inspect the selected Toolset with toolset_get before loading it. Search individual tools for a narrow task, inspect the selected schema with tool_get, then activate it with tool_load. Inspection identifies capabilities that require an active workspace; do not load those capabilities until a workspace is mounted. Load only what the current task needs."
}

func workspaceRequiredInstruction(kind, name, id string) string {
	return fmt.Sprintf("%s %q (%s) requires an active Robot workspace. Start or continue this conversation with a workspace, then try again.", kind, name, id)
}

func formatToolsetInstruction(ctx agent.ReadonlyContext, def Definition) (string, error) {
	instruction := def.Instruction
	if def.InstructionProvider != nil {
		var err error
		instruction, err = def.InstructionProvider(ctx)
		if err != nil {
			return "", err
		}
	}
	if strings.TrimSpace(instruction) == "" {
		return "Active Toolset: " + def.Name + " — " + def.Description, nil
	}
	return "Active Toolset: " + def.Name + "\n" + instruction, nil
}

func catalogueItems(definitions []Definition) []CatalogueItem {
	items := make([]CatalogueItem, 0, len(definitions))
	for _, def := range definitions {
		items = append(items, CatalogueItem{
			ID:          def.ID,
			Name:        def.Name,
			Description: def.Description,
		})
	}
	return items
}

func appendInstruction(request *model.LLMRequest, instruction string) {
	if strings.TrimSpace(instruction) == "" {
		return
	}
	if request.Config == nil {
		request.Config = &genai.GenerateContentConfig{}
	}
	if request.Config.SystemInstruction == nil {
		request.Config.SystemInstruction = genai.NewContentFromText(instruction, genai.RoleUser)
		return
	}
	request.Config.SystemInstruction.Parts = append(request.Config.SystemInstruction.Parts, genai.NewPartFromText("\n\n"+instruction))
}
