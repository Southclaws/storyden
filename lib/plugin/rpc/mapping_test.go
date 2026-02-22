package rpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManifestToMapUsesJSONFieldNames(t *testing.T) {
	manifest := Manifest{
		Author:         "you",
		Command:        "./plugin",
		Description:    "desc",
		EventsConsumed: []Event{EventEventThreadPublished},
		ID:             "my-plugin",
		Name:           "My Plugin",
		Version:        "0.0.1",
	}

	m := manifest.ToMap()

	events, ok := m["events_consumed"]
	require.True(t, ok, "expected events_consumed key in map")
	assert.NotContains(t, m, "EventsConsumed")

	eventList, ok := events.([]any)
	require.True(t, ok, "expected events_consumed to be []any")
	require.Len(t, eventList, 1)
	assert.Equal(t, string(EventEventThreadPublished), eventList[0])
}

func TestManifestFromMapDecodesEventsConsumed(t *testing.T) {
	raw := map[string]any{
		"id":              "my-plugin",
		"name":            "My Plugin",
		"description":     "desc",
		"version":         "0.0.1",
		"author":          "you",
		"command":         "./plugin",
		"events_consumed": []any{"EventThreadPublished", "EventThreadUpdated"},
	}

	manifest, err := ManifestFromMap(raw)
	require.NoError(t, err)
	require.NotNil(t, manifest)
	assert.Equal(t, []Event{
		EventEventThreadPublished,
		EventEventThreadUpdated,
	}, manifest.EventsConsumed)
}

func TestManifestFromMapValidatesManifest(t *testing.T) {
	raw := map[string]any{
		"id":              "bad id",
		"name":            "My Plugin",
		"description":     "desc",
		"version":         "0.0.1",
		"author":          "you",
		"command":         "./plugin",
		"events_consumed": []any{},
	}

	manifest, err := ManifestFromMap(raw)
	require.Error(t, err)
	assert.Nil(t, manifest)
}

func TestManifestFromMapRejectsInvalidEventsConsumed(t *testing.T) {
	raw := map[string]any{
		"id":              "my-plugin",
		"name":            "My Plugin",
		"description":     "desc",
		"version":         "0.0.1",
		"author":          "you",
		"command":         "./plugin",
		"events_consumed": []any{"EventThreadPublished", "EventThreadPublishedddsada"},
	}

	manifest, err := ManifestFromMap(raw)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid events_consumed value")
	assert.Nil(t, manifest)
}
