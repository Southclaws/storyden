package toolsets

import (
	"context"
	"io"
	"iter"
	"log/slog"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	adkagent "google.golang.org/adk/v2/agent"
	adksession "google.golang.org/adk/v2/session"
	adktool "google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
	"google.golang.org/genai"

	robotresource "github.com/Southclaws/storyden/app/resources/robot"
	robottools "github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	"github.com/Southclaws/storyden/app/services/semdex/robot/toolsets/system_documents"
	"github.com/Southclaws/storyden/app/services/semdex/robot/toolsets/system_memory"
	"github.com/Southclaws/storyden/app/services/semdex/robot/toolsets/system_memory_management"
	"github.com/Southclaws/storyden/app/services/semdex/robot/toolsets/system_moderation"
	"github.com/Southclaws/storyden/app/services/semdex/robot/workspacestate"
	"github.com/Southclaws/storyden/lib/mcp"
)

func TestSystemDocumentsToolsetOwnsNavigationInstruction(t *testing.T) {
	definitions := systemDefinitions()
	for _, definition := range definitions {
		if definition.ID != system_documents.ID {
			continue
		}
		assert.ElementsMatch(t, []string{"document_close", "document_get", "document_list", "document_search"}, definition.ToolNames)
		require.NotNil(t, definition.InstructionProvider)
		instruction, err := definition.InstructionProvider(&toolsetReadonlyContext{Context: context.Background(), state: toolsetReadonlyState{}})
		require.NoError(t, err)
		assert.Contains(t, instruction, "document_get")
		assert.Contains(t, instruction, "document_search")
		return
	}
	t.Fatalf("%s is not registered", system_documents.ID)
}

func TestSystemModerationToolsetExcludesContentPurge(t *testing.T) {
	for _, definition := range systemDefinitions() {
		if definition.ID != system_moderation.ID {
			continue
		}
		assert.ElementsMatch(t, []string{
			"member_reinstate",
			"member_suspend",
			"report_create",
			"report_get",
			"report_list",
			"report_update",
		}, definition.ToolNames)
		assert.NotContains(t, definition.ToolNames, "content_purge")
		return
	}
	t.Fatalf("%s is not registered", system_moderation.ID)
}

func TestSystemMemoryToolsetsSeparateUseFromManagement(t *testing.T) {
	foundMemory := false
	foundManagement := false
	for _, definition := range systemDefinitions() {
		switch definition.ID {
		case system_memory.ID:
			foundMemory = true
			assert.ElementsMatch(t, []string{"memory_create", "memory_search"}, definition.ToolNames)
			assert.Equal(t, "Knowledge graph memory", definition.Name)
			assert.Contains(t, definition.Instruction, "Make memory decisions autonomously")
			assert.Contains(t, definition.Instruction, "Do not search or write memory reflexively")
			assert.Contains(t, definition.Instruction, "A simple query or automation needs no memory search")
			assert.Contains(t, definition.Instruction, "Every supplied text term and graph field in one call is ANDed")
			assert.Contains(t, definition.Instruction, "Do not search merely to prove that no duplicate exists")
			assert.Contains(t, definition.Instruction, "Duplicates are acceptable and can be consolidated later")
			assert.Contains(t, definition.Instruction, "always supply subject, predicate, and object together")
			assert.Contains(t, definition.Instruction, "(barney, is_searching_for, content_by_southclaws)")
			assert.NotContains(t, definition.Instruction, "memory_open")
			assert.NotContains(t, definition.Instruction, "memory_update")
		case system_memory_management.ID:
			foundManagement = true
			assert.ElementsMatch(t, []string{"memory_list", "memory_move", "memory_open", "memory_update"}, definition.ToolNames)
			assert.Equal(t, "Knowledge graph management", definition.Name)
			assert.Contains(t, definition.Instruction, "deliberate inspection, correction, and consolidation")
			assert.Contains(t, definition.Instruction, "scheduled maintenance Robot")
		}
	}
	assert.True(t, foundMemory, "%s is not registered", system_memory.ID)
	assert.True(t, foundManagement, "%s is not registered", system_memory_management.ID)
}

func TestBuildDeduplicatesToolsSharedByMultipleToolsets(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	toolRegistry := robottools.NewRegistry(logger)
	registerSystemTestTool(t, toolRegistry, mcp.GetThreadGetTool())
	registerSystemTestTool(t, toolRegistry, mcp.GetThreadSearchTool())

	registry := NewRegistry(logger, nil, toolRegistry)
	require.NoError(t, registry.Register(Definition{
		ID:          "system.discussions",
		Name:        "Discussion management",
		Instruction: "Manage discussions.",
		ToolNames:   []string{"thread_get", "thread_search"},
	}))
	require.NoError(t, registry.Register(Definition{
		ID:          "system.search",
		Name:        "Search",
		Instruction: "Search content.",
		ToolNames:   []string{"thread_search"},
	}))

	built, missing, err := registry.Build(context.Background(), []string{
		"system.search",
		"system.discussions",
	})
	require.NoError(t, err)
	assert.Empty(t, missing)
	require.Len(t, built, 2, "each Toolset must retain its own prompt even when all of its tools overlap")

	var names []string
	for _, toolset := range built {
		resolved, err := toolset.Tools(nil)
		require.NoError(t, err)
		for _, tool := range resolved {
			names = append(names, tool.Name())
		}
	}
	assert.Equal(t, []string{"thread_get", "thread_search"}, names)
}

func TestWorkspaceToolsetIsHiddenUntilSessionHasWorkspace(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	toolRegistry := robottools.NewRegistry(logger)
	builderCalls := 0
	require.NoError(t, toolRegistry.Register(&robottools.Tool{
		Definition: &mcp.ToolDefinition{
			Name:              "workspace_only",
			Title:             "Workspace Only",
			Description:       "Perform work inside the active workspace.",
			RequiresWorkspace: true,
		},
		Builder: func(context.Context) (adktool.Tool, error) {
			builderCalls++
			return functiontool.New(functiontool.Config{
				Name:        "workspace_only",
				Description: "Perform work inside the active workspace.",
			}, func(adkagent.Context, struct{}) (map[string]any, error) {
				return map[string]any{"ok": true}, nil
			})
		},
	}))

	registry := NewRegistry(logger, nil, toolRegistry)
	require.NoError(t, registry.Register(Definition{
		ID:        "system.workspace_test",
		Name:      "Workspace test",
		ToolNames: []string{"workspace_only"},
	}))

	definition, err := registry.Get(context.Background(), "system.workspace_test")
	require.NoError(t, err)
	assert.True(t, definition.RequiresWorkspace)

	built, missing, err := registry.Build(context.Background(), []string{"system.workspace_test"})
	require.NoError(t, err)
	assert.Empty(t, missing)
	require.Len(t, built, 1)

	withoutWorkspace := &toolsetReadonlyContext{Context: context.Background(), state: toolsetReadonlyState{}}
	available, err := built[0].Tools(withoutWorkspace)
	require.NoError(t, err)
	assert.Empty(t, available)
	assert.Equal(t, 0, builderCalls)

	withWorkspace := &toolsetReadonlyContext{Context: context.Background(), state: mountedWorkspaceState()}
	available, err = built[0].Tools(withWorkspace)
	require.NoError(t, err)
	require.Len(t, available, 1)
	assert.Equal(t, "workspace_only", available[0].Name())
	assert.Equal(t, 1, builderCalls)
}

func registerSystemTestTool(t *testing.T, registry *robottools.Registry, definition *mcp.ToolDefinition) {
	t.Helper()
	require.NoError(t, registry.Register(&robottools.Tool{
		Definition: definition,
		Builder: func(context.Context) (adktool.Tool, error) {
			return functiontool.New(functiontool.Config{
				Name:        definition.Name,
				Description: definition.Description,
				InputSchema: definition.InputSchema,
			}, func(adkagent.Context, struct{}) (map[string]any, error) {
				return map[string]any{"ok": true}, nil
			})
		},
	}))
}

type toolsetReadonlyContext struct {
	context.Context
	state adksession.ReadonlyState
}

func (c *toolsetReadonlyContext) UserContent() *genai.Content             { return nil }
func (c *toolsetReadonlyContext) InvocationID() string                    { return "test-invocation" }
func (c *toolsetReadonlyContext) AgentName() string                       { return "test-agent" }
func (c *toolsetReadonlyContext) ReadonlyState() adksession.ReadonlyState { return c.state }
func (c *toolsetReadonlyContext) UserID() string                          { return "test-user" }
func (c *toolsetReadonlyContext) AppName() string                         { return "test-app" }
func (c *toolsetReadonlyContext) SessionID() string                       { return "test-session" }
func (c *toolsetReadonlyContext) Branch() string                          { return "" }

type toolsetReadonlyState map[string]any

func (s toolsetReadonlyState) Get(key string) (any, error) {
	value, ok := s[key]
	if !ok {
		return nil, adksession.ErrStateKeyNotExist
	}
	return value, nil
}

func (s toolsetReadonlyState) All() iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		for key, value := range s {
			if !yield(key, value) {
				return
			}
		}
	}
}

func mountedWorkspaceState() toolsetReadonlyState {
	mount := &robotresource.WorkspaceMount{
		WorkspaceID:         robotresource.WorkspaceID(xid.New()),
		WorkspaceInstanceID: robotresource.WorkspaceInstanceID(xid.New()),
		Provider:            robotresource.WorkspaceProviderLocal,
	}
	return toolsetReadonlyState{
		workspacestate.WorkspaceStateKey: workspacestate.MountToState(mount),
	}
}
