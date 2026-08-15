package pluginbuilder

import (
	"context"
	"iter"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	adksession "google.golang.org/adk/v2/session"
)

func TestResolveInstallationUsesBoundTarget(t *testing.T) {
	id := xid.New()
	ctx := newPluginBuilderTestContext(map[string]any{
		pluginBuildTargetStateKey: pluginBuildTarget{
			Mode:           "existing",
			InstallationID: id.String(),
			ManifestID:     "welcome-plugin",
		},
	})

	agent := &Agent{}
	got, found, action, err := agent.resolveInstallation(ctx, "welcome-plugin")
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, id.String(), got.String())
	require.Equal(t, "updated", action)
}

func TestResolveInstallationRejectsManifestMismatch(t *testing.T) {
	ctx := newPluginBuilderTestContext(map[string]any{
		pluginBuildTargetStateKey: pluginBuildTarget{
			Mode:           "existing",
			InstallationID: xid.New().String(),
			ManifestID:     "welcome-plugin",
		},
	})

	agent := &Agent{}
	_, _, _, err := agent.resolveInstallation(ctx, "other-plugin")
	require.ErrorContains(t, err, "start a new chat")
}

func TestResolveInstallationCreatesWhenUnbound(t *testing.T) {
	ctx := newPluginBuilderTestContext(nil)

	agent := &Agent{}
	got, found, action, err := agent.resolveInstallation(ctx, "welcome-plugin")
	require.NoError(t, err)
	require.False(t, found)
	require.Zero(t, got)
	require.Empty(t, action)
}

func TestSetPluginBuildTargetPersistsToState(t *testing.T) {
	ctx := newPluginBuilderTestContext(nil)
	id := xid.New().String()

	agent := &Agent{}
	err := agent.setPluginBuildTarget(ctx, pluginBuildTarget{
		InstallationID: id,
		ManifestID:     "welcome-plugin",
	})
	require.NoError(t, err)

	target, ok, err := pluginBuildTargetFromContext(ctx)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, pluginBuildTargetModeNew, target.Mode)
	require.Equal(t, id, target.InstallationID)
	require.Equal(t, "welcome-plugin", target.ManifestID)
}

func TestSetPluginBuildTargetAllowsNewPluginBeforeInstall(t *testing.T) {
	ctx := newPluginBuilderTestContext(nil)

	agent := &Agent{}
	err := agent.setPluginBuildTarget(ctx, pluginBuildTarget{
		Mode:       pluginBuildTargetModeNew,
		ManifestID: "welcome-plugin",
	})
	require.NoError(t, err)

	target, ok, err := pluginBuildTargetFromContext(ctx)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, pluginBuildTargetModeNew, target.Mode)
	require.Empty(t, target.InstallationID)
	require.Equal(t, "welcome-plugin", target.ManifestID)
}

type pluginBuilderTestContext struct {
	context.Context
	state pluginBuilderTestState
}

func newPluginBuilderTestContext(values map[string]any) *pluginBuilderTestContext {
	state := pluginBuilderTestState{}
	for key, value := range values {
		state[key] = value
	}
	return &pluginBuilderTestContext{
		Context: context.Background(),
		state:   state,
	}
}

func (c *pluginBuilderTestContext) State() adksession.State {
	return c.state
}

type pluginBuilderTestState map[string]any

func (s pluginBuilderTestState) Get(key string) (any, error) {
	value, ok := s[key]
	if !ok {
		return nil, adksession.ErrStateKeyNotExist
	}
	return value, nil
}

func (s pluginBuilderTestState) Set(key string, value any) error {
	s[key] = value
	return nil
}

func (s pluginBuilderTestState) All() iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		for key, value := range s {
			if !yield(key, value) {
				return
			}
		}
	}
}
