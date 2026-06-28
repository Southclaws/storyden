package pluginbuilder

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rs/xid"
	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	pluginresource "github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/plugin_logger"
)

const (
	defaultLogMaxLines   = 100
	defaultLogWaitMillis = 750
	maxLogMaxLines       = 500
	maxLogWaitMillis     = 3000
)

type pluginLogReader interface {
	StreamPluginLogs(ctx context.Context, pluginID pluginresource.InstallationID) (*plugin_logger.LogStream, error)
}

type PluginLogsInput struct {
	MaxLines   int `json:"max_lines,omitempty" jsonschema:"Maximum recent log lines to return"`
	WaitMillis int `json:"wait_millis,omitempty" jsonschema:"How long to wait for current log output before returning"`
}

type PluginLogsResult struct {
	InstallationID string   `json:"installation_id"`
	Lines          []string `json:"lines"`
	Truncated      bool     `json:"truncated"`
	Message        string   `json:"message"`
}

func (a *Agent) addLogTools(add toolAdder) error {
	return add(functiontool.New(functiontool.Config{
		Name: "plugin_logs_read",
		Description: `Read recent runtime logs for an installed supervised plugin.

Use this whenever the user asks to "check the logs", "see what happened at
runtime", or debug whether an installed plugin reacted to an event.

This tool reads runtime output only. It does not inspect source code. Do not use
plugin_go_symbol_search as a substitute for this tool when the user asks about
logs. It reads logs for the plugin bound to the current workspace.`,
	}, func(ctx adktool.Context, args PluginLogsInput) (PluginLogsResult, error) {
		return a.ReadPluginLogs(ctx, args)
	}))
}

func (a *Agent) ReadPluginLogs(ctx context.Context, in PluginLogsInput) (PluginLogsResult, error) {
	target, ok, err := pluginBuildTargetFromContext(ctx)
	if err != nil {
		return PluginLogsResult{}, err
	}
	if !ok || strings.TrimSpace(target.InstallationID) == "" {
		return PluginLogsResult{}, errors.New("no plugin installation is bound to this workspace; install the plugin first")
	}

	rawID := strings.TrimSpace(target.InstallationID)
	parsed, err := xid.FromString(rawID)
	if err != nil {
		return PluginLogsResult{}, fmt.Errorf("invalid bound plugin installation: %w", err)
	}

	maxLines := in.MaxLines
	if maxLines <= 0 {
		maxLines = defaultLogMaxLines
	}
	if maxLines > maxLogMaxLines {
		maxLines = maxLogMaxLines
	}

	waitMillis := in.WaitMillis
	if waitMillis <= 0 {
		waitMillis = defaultLogWaitMillis
	}
	if waitMillis > maxLogWaitMillis {
		waitMillis = maxLogWaitMillis
	}

	readCtx, cancel := context.WithTimeout(ctx, time.Duration(waitMillis)*time.Millisecond)
	defer cancel()

	stream, err := a.logs.StreamPluginLogs(readCtx, pluginresource.InstallationID(parsed))
	if err != nil {
		return PluginLogsResult{}, err
	}

	lines := []string{}
	truncated := false
	for {
		select {
		case line, ok := <-stream.Lines:
			if !ok {
				return pluginLogsResult(rawID, lines, truncated), nil
			}
			lines = append(lines, line)
			if len(lines) > maxLines {
				copy(lines, lines[1:])
				lines = lines[:maxLines]
				truncated = true
			}
		case <-readCtx.Done():
			return pluginLogsResult(rawID, lines, truncated), nil
		}
	}
}

func pluginLogsResult(installationID string, lines []string, truncated bool) PluginLogsResult {
	message := "no log lines were available before the read timed out"
	if len(lines) > 0 {
		message = fmt.Sprintf("read %d recent log lines", len(lines))
		if truncated {
			message = fmt.Sprintf("read the last %d log lines; earlier lines were omitted", len(lines))
		}
	}

	return PluginLogsResult{
		InstallationID: installationID,
		Lines:          lines,
		Truncated:      truncated,
		Message:        message,
	}
}
