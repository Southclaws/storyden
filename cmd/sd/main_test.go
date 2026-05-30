package main

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/carapace-sh/carapace"

	"github.com/Southclaws/storyden/cmd/sd/internal/cli"
)

func TestNewLoggerUsesConfiguredErrorStream(t *testing.T) {
	var stderr bytes.Buffer

	logger := newLogger(cli.Streams{Err: &stderr})
	logger.Error("meow?", slog.String("sound", "mrrp"))

	output := stderr.String()
	if !bytes.Contains([]byte(output), []byte("meow?")) {
		t.Fatalf("expected log output to contain message, got %q", output)
	}
	if !bytes.Contains([]byte(output), []byte("sound")) {
		t.Fatalf("expected log output to contain attribute key, got %q", output)
	}
	if !bytes.Contains([]byte(output), []byte("mrrp")) {
		t.Fatalf("expected log output to contain attribute value, got %q", output)
	}
}

func TestCarapace(t *testing.T) {
	carapace.Test(t)
}
