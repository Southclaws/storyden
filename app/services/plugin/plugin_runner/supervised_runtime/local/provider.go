package local

import (
	"context"
	"log/slog"
	"net/url"
	"time"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/resources/plugin/plugin_reader"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/plugin_logger"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/supervised_runtime"
)

type provider struct {
	pluginLogger *plugin_logger.Writer
	pluginReader *plugin_reader.Reader

	dataPath string

	maxRestartAttempts    int
	maxBackoffDuration    time.Duration
	runtimeCrashThreshold time.Duration
	runtimeCrashBackoff   time.Duration
}

func NewProvider(
	pluginLogger *plugin_logger.Writer,
	pluginReader *plugin_reader.Reader,
	dataPath string,
	maxRestartAttempts int,
	maxBackoffDuration time.Duration,
	runtimeCrashThreshold time.Duration,
	runtimeCrashBackoff time.Duration,
) supervised_runtime.Provider {
	return &provider{
		pluginLogger:          pluginLogger,
		pluginReader:          pluginReader,
		dataPath:              dataPath,
		maxRestartAttempts:    maxRestartAttempts,
		maxBackoffDuration:    maxBackoffDuration,
		runtimeCrashThreshold: runtimeCrashThreshold,
		runtimeCrashBackoff:   runtimeCrashBackoff,
	}
}

func (p *provider) New(
	id plugin.InstallationID,
	bin []byte,
	manifest *plugin.Validated,
	parentLogger *slog.Logger,
	serverURL url.URL,
	parentCtx context.Context,
) supervised_runtime.Runtime {
	return newRuntime(
		id,
		serverURL,
		parentCtx,
		bin,
		manifest,
		parentLogger,
		p.pluginLogger,
		p.pluginReader,
		p.dataPath,
		p.maxRestartAttempts,
		p.maxBackoffDuration,
		p.runtimeCrashThreshold,
		p.runtimeCrashBackoff,
	)
}
