package toolsets_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/robot/robot_toolset"
	robottools "github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	robottoolsets "github.com/Southclaws/storyden/app/services/semdex/robot/toolsets"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestRegistrySearchesAcrossSystemPluginAndCustomToolsets(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		repo *robot_toolset.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			author, err := db.Account.Create().
				SetHandle("registry-owner").
				SetName("Registry Owner").
				Save(ctx)
			require.NoError(t, err)

			custom, err := repo.Create(
				ctx,
				account.AccountID(author.ID),
				"Plugin review",
				"Review plugin manifests and source",
				"Inspect before changing files.",
				[]string{"plugin_manifest_validate"},
				map[string]any{},
			)
			require.NoError(t, err)

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			registry := robottoolsets.NewRegistry(logger, repo, robottools.NewRegistry(logger))
			require.NoError(t, registry.Register(robottoolsets.Definition{
				ID:          "plugin.example.builder",
				Name:        "Plugin builder",
				Description: "Build and package plugins",
				Instruction: "Work inside the active workspace.",
				ToolNames:   []string{"plugin_build", "plugin_build"},
				Source:      robottoolsets.SourcePlugin,
				SourceRef:   "plugin-installation-1",
			}))
			require.NoError(t, registry.Register(robottoolsets.Definition{
				ID:          "system.research",
				Name:        "Web research",
				Description: "Search and verify web sources",
				Source:      robottoolsets.SourceSystem,
			}))

			results, err := registry.Search(ctx, "plugin", 10)
			require.NoError(t, err)
			require.Len(t, results, 2)
			assert.Equal(t, "plugin.example.builder", results[0].ID)
			assert.Equal(t, custom.ID.String(), results[1].ID)
			assert.Equal(t, robottoolsets.SourceCustom, results[1].Source)
			assert.True(t, results[1].Editable)

			plugin, err := registry.Get(ctx, "plugin.example.builder")
			require.NoError(t, err)
			assert.Equal(t, []string{"plugin_build"}, plugin.ToolNames)
			assert.False(t, plugin.Editable)

			require.NoError(t, registry.Validate(ctx, []string{"system.research", custom.ID.String()}))
			require.Error(t, registry.Validate(ctx, []string{"missing.toolset"}))
			require.Error(t, registry.Register(robottoolsets.Definition{ID: "system.research", Name: "Duplicate"}))

			registry.UnregisterSource("plugin-installation-1")
			_, err = registry.Get(ctx, "plugin.example.builder")
			require.Error(t, err)
		}))
	}))
}

func TestRegistryComposesDirectToolsWithToolsets(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	toolRegistry := robottools.NewRegistry(logger)
	toolRegistry.RegisterAlias("thread_read", "thread_get")
	registry := robottoolsets.NewRegistry(logger, nil, toolRegistry)
	require.NoError(t, registry.Register(robottoolsets.Definition{
		ID:        "system.thread_reader",
		Name:      "Thread reader",
		ToolNames: []string{"thread_get", "thread_list"},
	}))

	filtered, err := registry.RemoveContainedTools(
		ctx,
		[]string{"tag_list", "thread_get", "member_search"},
		[]string{"system.thread_reader"},
	)
	require.NoError(t, err)
	assert.Equal(t, []string{"tag_list", "member_search"}, filtered)

	err = registry.ValidateDirectTools(ctx, []string{"thread_get"}, []string{"system.thread_reader"})
	require.Error(t, err)
	assert.EqualError(t, err, `tool "thread_get" cannot be added directly because Toolset "Thread reader" (system.thread_reader) already contains that tool`)

	require.NoError(t, registry.ValidateDirectTools(ctx, []string{"tag_list"}, []string{"system.thread_reader"}))

	err = registry.ValidateDirectTools(ctx, []string{"thread_read"}, []string{"system.thread_reader"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), `tool "thread_read"`)
}
