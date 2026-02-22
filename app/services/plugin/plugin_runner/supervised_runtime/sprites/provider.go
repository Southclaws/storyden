package sprites

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/superfly/sprites-go"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/resources/plugin/plugin_reader"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/plugin_logger"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/supervised_runtime"
)

const spritesHTTPTimeout = 2 * time.Minute

type provider struct {
	client *sprites.Client

	pluginLogger *plugin_logger.Writer
	pluginReader *plugin_reader.Reader

	dataPath string

	maxRestartAttempts    int
	maxBackoffDuration    time.Duration
	runtimeCrashThreshold time.Duration
	runtimeCrashBackoff   time.Duration
}

func NewProvider(
	apiKey string,
	pluginLogger *plugin_logger.Writer,
	pluginReader *plugin_reader.Reader,
	dataPath string,
	maxRestartAttempts int,
	maxBackoffDuration time.Duration,
	runtimeCrashThreshold time.Duration,
	runtimeCrashBackoff time.Duration,
) supervised_runtime.Provider {
	client := sprites.New(
		apiKey,
		sprites.WithHTTPClient(&http.Client{
			Timeout: spritesHTTPTimeout,
		}),
	)

	return &provider{
		client:                client,
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
		p.client,
		p.pluginLogger,
		p.pluginReader,
		p.dataPath,
		p.maxRestartAttempts,
		p.maxBackoffDuration,
		p.runtimeCrashThreshold,
		p.runtimeCrashBackoff,
	)
}
