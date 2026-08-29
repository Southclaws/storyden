package toolsets_test

import (
	"context"
	"errors"
	"iter"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	adkagent "google.golang.org/adk/v2/agent"
	adksession "google.golang.org/adk/v2/session"
	adktool "google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/robot/robot_toolset"
	robottools "github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	robottoolsets "github.com/Southclaws/storyden/app/services/semdex/robot/toolsets"
	"github.com/Southclaws/storyden/app/services/semdex/robot/toolsets/system_memory_management"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/lib/mcp"
)

func TestMemoryManagementToolsetIsDiscoverable(t *testing.T) {
	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		registry *robottoolsets.Registry,
	) {
		lc.Append(fx.StartHook(func() {
			results, err := registry.Search(ctx, "knowledge graph management", 20)
			require.NoError(t, err)

			for _, result := range results {
				if result.ID == system_memory_management.ID {
					assert.Equal(t, "Knowledge graph management", result.Name)
					assert.ElementsMatch(t, []string{"memory_list", "memory_move", "memory_open", "memory_update"}, result.ToolNames)
					return
				}
			}

			t.Fatalf("%s was not returned by Toolset search", system_memory_management.ID)
		}))
	}))
}

func TestDiscoveryContinuesWhenRegisteredCapabilitiesFailToInitialize(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		logger *slog.Logger,
		repo *robot_toolset.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			toolRegistry := robottools.NewRegistry(logger)
			require.NoError(t, toolRegistry.Register(&robottools.Tool{
				Definition: &mcp.ToolDefinition{
					Name:        "plugin__broken__search",
					Title:       "Broken plugin search",
					Description: "Search through a plugin that is currently broken.",
				},
				Source: "plugin",
				Builder: func(context.Context) (adktool.Tool, error) {
					return nil, errors.New("plugin process is unavailable")
				},
			}))
			require.NoError(t, toolRegistry.Register(&robottools.Tool{
				Definition: &mcp.ToolDefinition{
					Name:        "healthy_search",
					Title:       "Healthy search",
					Description: "Search a healthy source.",
				},
				Builder: func(context.Context) (adktool.Tool, error) {
					return functiontool.New(functiontool.Config{
						Name:        "healthy_search",
						Description: "Search a healthy source.",
					}, func(adkagent.Context, struct{}) (map[string]any, error) {
						return map[string]any{"ok": true}, nil
					})
				},
			}))

			registry := robottoolsets.NewRegistry(logger, repo, toolRegistry)
			require.NoError(t, registry.Register(robottoolsets.Definition{
				ID:          "plugin.broken.search",
				Name:        "Broken search plugin",
				Description: "Capabilities from a broken plugin.",
				ToolNames:   []string{"plugin__broken__search"},
				Source:      robottoolsets.SourcePlugin,
				SourceRef:   "plugin-installation-broken",
			}))
			require.NoError(t, registry.Register(robottoolsets.Definition{
				ID:          "plugin.broken.builder",
				Name:        "Broken builder",
				Description: "A Toolset whose builder fails.",
				Source:      robottoolsets.SourcePlugin,
				SourceRef:   "plugin-installation-broken",
				Builder: func(context.Context) (adktool.Toolset, error) {
					return nil, errors.New("plugin workspace failed to initialize")
				},
			}))

			discovery, err := registry.Discovery(ctx)
			require.NoError(t, err)
			available, err := discovery.Tools(&discoveryReadonlyContext{Context: ctx, state: discoveryReadonlyState{}})
			require.NoError(t, err)

			names := make([]string, 0, len(available))
			for _, item := range available {
				names = append(names, item.Name())
			}
			assert.Contains(t, names, "healthy_search")
			assert.Contains(t, names, "tool_runtime_issues")
			assert.NotContains(t, names, "plugin__broken__search")
		}))
	}))
}

type discoveryReadonlyContext struct {
	context.Context
	state adksession.ReadonlyState
}

func (c *discoveryReadonlyContext) UserContent() *genai.Content             { return nil }
func (c *discoveryReadonlyContext) InvocationID() string                    { return "test-invocation" }
func (c *discoveryReadonlyContext) AgentName() string                       { return "test-agent" }
func (c *discoveryReadonlyContext) ReadonlyState() adksession.ReadonlyState { return c.state }
func (c *discoveryReadonlyContext) UserID() string                          { return "test-user" }
func (c *discoveryReadonlyContext) AppName() string                         { return "test-app" }
func (c *discoveryReadonlyContext) SessionID() string                       { return "test-session" }
func (c *discoveryReadonlyContext) Branch() string                          { return "" }

type discoveryReadonlyState map[string]any

func (s discoveryReadonlyState) Get(key string) (any, error) {
	value, ok := s[key]
	if !ok {
		return nil, adksession.ErrStateKeyNotExist
	}
	return value, nil
}

func (s discoveryReadonlyState) All() iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		for key, value := range s {
			if !yield(key, value) {
				return
			}
		}
	}
}
