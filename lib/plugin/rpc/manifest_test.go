package rpc

import (
	"strings"
	"testing"

	"github.com/Southclaws/fault/fmsg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateName(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		require.NoError(t, ValidateName("My Plugin"))
	})

	t.Run("empty", func(t *testing.T) {
		err := ValidateName("")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidPluginName)
	})

	t.Run("too long", func(t *testing.T) {
		err := ValidateName(strings.Repeat("a", 101))
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrPluginNameTooLong)
	})
}

func TestManifestValidate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		m := Manifest{
			ID:             "my-plugin",
			Name:           "My Plugin",
			Description:    "desc",
			Version:        "0.0.1",
			Author:         "you",
			Command:        "./plugin",
			EventsConsumed: []Event{EventEventThreadPublished},
		}

		require.NoError(t, m.Validate())
	})

	t.Run("invalid fields", func(t *testing.T) {
		m := Manifest{
			ID:             "bad id",
			Name:           "",
			Description:    "desc",
			Version:        "0.0.1",
			Author:         "bad author",
			Command:        "./plugin",
			EventsConsumed: []Event{"EventThreadPublishedddsada"},
		}

		err := m.Validate()
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid plugin ID")
		require.ErrorContains(t, err, "invalid plugin author")
		require.ErrorContains(t, err, "invalid plugin name")
		require.ErrorContains(t, err, "invalid events_consumed value")
	})
}

func TestParseManifest(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		input := []byte(`{
			"id":"my-plugin",
			"name":"My Plugin",
			"description":"desc",
			"version":"0.0.1",
			"author":"you",
			"command":"./plugin",
			"events_consumed":["EventThreadPublished"]
		}`)

		manifest, err := ParseManifest(input)
		require.NoError(t, err)
		require.NotNil(t, manifest)
		assert.Equal(t, []Event{EventEventThreadPublished}, manifest.EventsConsumed)
	})

	t.Run("invalid event", func(t *testing.T) {
		input := []byte(`{
			"id":"my-plugin",
			"name":"My Plugin",
			"description":"desc",
			"version":"0.0.1",
			"author":"you",
			"command":"./plugin",
			"events_consumed":["EventThreadPublishedddsada"]
		}`)

		manifest, err := ParseManifest(input)
		require.Error(t, err)
		require.Nil(t, manifest)
		require.ErrorContains(t, err, "invalid events_consumed value")
		require.Contains(t, fmsg.GetIssue(err), `Field "events_consumed" contains unknown event`)
	})
}
