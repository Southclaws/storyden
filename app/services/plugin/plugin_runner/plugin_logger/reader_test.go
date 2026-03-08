package plugin_logger

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/internal/config"
)

func TestStreamPluginLogs_StaysOpenWhileIdleThenStreamsNewLines(t *testing.T) {
	tmpDir := t.TempDir()
	pluginID := plugin.InstallationID(xid.New())

	logDir := getPluginLogDirectory(tmpDir, pluginID)
	require.NoError(t, os.MkdirAll(logDir, 0o755))

	currentLogPath := getOutputPath(tmpDir, pluginID)
	require.NoError(t, os.WriteFile(currentLogPath, []byte("line 1\n"), 0o644))

	reader := newReader(config.Config{PluginDataPath: tmpDir})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := reader.StreamPluginLogs(ctx, pluginID)
	require.NoError(t, err)

	received := make(chan string, 16)
	go func() {
		for line := range stream.Lines {
			received <- line
		}
	}()

	assertEventuallyContains(t, received, "line 1")

	// Simulates opening logs while plugin is inactive: no writes for a while.
	time.Sleep(350 * time.Millisecond)

	f, err := os.OpenFile(currentLogPath, os.O_APPEND|os.O_WRONLY, 0o644)
	require.NoError(t, err)
	_, err = f.WriteString("line 2\n")
	require.NoError(t, err)
	require.NoError(t, f.Close())

	assertEventuallyContains(t, received, "line 2")
}

func TestStreamPluginLogs_FollowsOutputAfterRotation(t *testing.T) {
	tmpDir := t.TempDir()
	pluginID := plugin.InstallationID(xid.New())

	logDir := getPluginLogDirectory(tmpDir, pluginID)
	require.NoError(t, os.MkdirAll(logDir, 0o755))

	currentLogPath := getOutputPath(tmpDir, pluginID)
	require.NoError(t, os.WriteFile(currentLogPath, []byte("line 1\n"), 0o644))

	reader := newReader(config.Config{PluginDataPath: tmpDir})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := reader.StreamPluginLogs(ctx, pluginID)
	require.NoError(t, err)

	received := make(chan string, 32)
	go func() {
		for line := range stream.Lines {
			received <- line
		}
	}()

	assertEventuallyContains(t, received, "line 1")

	f1, err := os.OpenFile(currentLogPath, os.O_APPEND|os.O_WRONLY, 0o644)
	require.NoError(t, err)
	_, err = f1.WriteString("line 2\n")
	require.NoError(t, err)
	require.NoError(t, f1.Close())
	assertEventuallyContains(t, received, "line 2")

	// Simulate rotate-on-stop: output.log -> timestamped file then new output.log.
	rotatedPath := filepath.Join(logDir, "output-2026-02-21T11-50-22.965.log")
	require.NoError(t, os.Rename(currentLogPath, rotatedPath))
	require.NoError(t, os.WriteFile(currentLogPath, []byte("line 3\n"), 0o644))

	assertEventuallyContains(t, received, "line 3")
}

func TestStreamPluginLogs_WaitsForFirstLogFile(t *testing.T) {
	tmpDir := t.TempDir()
	pluginID := plugin.InstallationID(xid.New())

	// Intentionally do not create logs directory yet.
	reader := newReader(config.Config{PluginDataPath: tmpDir})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := reader.StreamPluginLogs(ctx, pluginID)
	require.NoError(t, err)

	received := make(chan string, 8)
	go func() {
		for line := range stream.Lines {
			received <- line
		}
	}()

	time.Sleep(300 * time.Millisecond)

	logDir := getPluginLogDirectory(tmpDir, pluginID)
	require.NoError(t, os.MkdirAll(logDir, 0o755))
	require.NoError(t, os.WriteFile(getOutputPath(tmpDir, pluginID), []byte("first boot line\n"), 0o644))

	assertEventuallyContains(t, received, "first boot line")
}

func assertEventuallyContains(t *testing.T, ch <-chan string, want string) {
	t.Helper()

	timeout := time.NewTimer(2 * time.Second)
	defer timeout.Stop()

	for {
		select {
		case got := <-ch:
			if got == want {
				return
			}
		case <-timeout.C:
			t.Fatalf("timed out waiting for line: %s", want)
		}
	}
}
