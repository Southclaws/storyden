package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractEventsIncludesFrontendMetadata(t *testing.T) {
	events, err := extractEventsFromSchema(filepath.Join("..", "..", "..", "api", "common", "events.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if len(events) == 0 {
		t.Fatal("expected events from the canonical schema")
	}

	first := events[0]
	if first.TypeName != "EventThreadPublished" {
		t.Fatalf("unexpected first event name: %q", first.TypeName)
	}
	if first.Label != "Thread published" {
		t.Fatalf("unexpected first event label: %q", first.Label)
	}
	if first.Description != "Emitted when a thread is visible as published, either on create or after a visibility change." {
		t.Fatalf("unexpected first event description: %q", first.Description)
	}
}

func TestGenerateWebEventNamesIncludesEscapedCatalogue(t *testing.T) {
	generated, err := generateWebEventNames([]Event{
		{
			TypeName:    "EventThreadPublished",
			HandlerName: "ThreadPublished",
			Label:       "Thread published",
			Description: `Emitted after a "published" transition.`,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	output := string(generated)
	for _, expected := range []string{
		`export const pluginEventNames = [`,
		`export const events = [`,
		`name: "EventThreadPublished"`,
		`label: "Thread published"`,
		`description: "Emitted after a \"published\" transition."`,
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("generated output missing %q", expected)
		}
	}
}
