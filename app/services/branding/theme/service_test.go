package theme

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/internal/mime"
)

func TestClassifyThemeAssetFilename(t *testing.T) {
	t.Parallel()

	stylesheet, err := classify("community.CSS")
	require.NoError(t, err)
	assert.Equal(t, asset.ThemeAssetKindStylesheet, stylesheet)

	script, err := classify("enhancements.js")
	require.NoError(t, err)
	assert.Equal(t, asset.ThemeAssetKindScript, script)

	for _, invalid := range []string{"../theme.css", "folder/theme.css", `folder\\theme.js`, ".", "theme.css.exe", ""} {
		invalid := invalid
		t.Run(invalid, func(t *testing.T) {
			_, err := classify(invalid)
			assert.Error(t, err)
		})
	}
}

func TestThemeAssetKind(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name      string
		mediaType string
		raw       string
		kind      asset.ThemeAssetKind
		valid     bool
	}{
		{
			name:      "stylesheet",
			mediaType: "text/css",
			raw:       `{"purpose":"theme","theme_kind":"stylesheet"}`,
			kind:      asset.ThemeAssetKindStylesheet,
			valid:     true,
		},
		{
			name:      "script",
			mediaType: "application/javascript",
			raw:       `{"purpose":"theme","theme_kind":"script"}`,
			kind:      asset.ThemeAssetKindScript,
			valid:     true,
		},
		{
			name:      "generic stylesheet",
			mediaType: "text/css",
			raw:       `{}`,
		},
		{
			name:      "wrong MIME class",
			mediaType: "application/pdf",
			raw:       `{"purpose":"theme","theme_kind":"stylesheet"}`,
		},
		{
			name:      "invalid kind",
			mediaType: "text/css",
			raw:       `{"purpose":"theme","theme_kind":"other"}`,
		},
		{
			name:      "invalid metadata",
			mediaType: "text/css",
			raw:       `{"purpose":"theme","theme_kind":123}`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			source := &asset.Asset{
				MIME:     mime.New(test.mediaType),
				Metadata: asset.Metadata{RawMessage: json.RawMessage(test.raw)},
			}

			kind, valid := themeAssetKind(source)
			require.Equal(t, test.valid, valid)
			require.Equal(t, test.kind, kind)
		})
	}
}
