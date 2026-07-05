package pluginbuilder

import (
	"context"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"

	pluginresource "github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/plugin_logger"
)

func TestReadPluginLogsReturnsTail(t *testing.T) {
	id := xid.New()
	ctx := newPluginBuilderTestContext(map[string]any{
		pluginBuildTargetStateKey: pluginBuildTarget{
			InstallationID: id.String(),
			ManifestID:     "welcome-plugin",
		},
	})
	agent := &Agent{
		logs: fakePluginLogReader{
			lines: []string{"line 1", "line 2", "line 3"},
		},
	}

	result, err := agent.ReadPluginLogs(ctx, PluginLogsInput{
		MaxLines:   2,
		WaitMillis: 10,
	})
	require.NoError(t, err)
	require.Equal(t, id.String(), result.InstallationID)
	require.True(t, result.Truncated)
	require.Equal(t, []string{"line 2", "line 3"}, result.Lines)
}

func TestReadPluginLogsShapesWaitingForConfigurationState(t *testing.T) {
	id := xid.New()
	ctx := newPluginBuilderTestContext(map[string]any{
		pluginBuildTargetStateKey: pluginBuildTarget{
			InstallationID: id.String(),
			ManifestID:     "discord-plugin",
		},
	})
	agent := &Agent{
		logs: fakePluginLogReader{
			lines: []string{`time=2026-07-05T12:37:29.345+01:00 level=INFO msg="plugin is waiting for configuration (missing discord_token or discord_channel_id)"`},
		},
	}

	result, err := agent.ReadPluginLogs(ctx, PluginLogsInput{WaitMillis: 10})
	require.NoError(t, err)
	require.Equal(t, "plugin started and is waiting for user configuration", result.Message)
	require.Contains(t, result.NextAction, "Stop inspecting internals")
	require.Contains(t, result.NextAction, "settings they need to provide")
}

func TestReadPluginLogsRequiresBoundInstallation(t *testing.T) {
	agent := &Agent{logs: fakePluginLogReader{}}

	_, err := agent.ReadPluginLogs(newPluginBuilderTestContext(nil), PluginLogsInput{})
	require.ErrorContains(t, err, "no plugin installation is bound")
}

type fakePluginLogReader struct {
	lines []string
}

func (f fakePluginLogReader) StreamPluginLogs(ctx context.Context, pluginID pluginresource.InstallationID) (*plugin_logger.LogStream, error) {
	lines := make(chan string, len(f.lines))
	done := make(chan struct{})

	go func() {
		defer close(lines)
		defer close(done)
		for _, line := range f.lines {
			select {
			case <-ctx.Done():
				return
			case lines <- line:
			}
		}
	}()

	return &plugin_logger.LogStream{
		Lines: lines,
		Done:  done,
	}, nil
}
