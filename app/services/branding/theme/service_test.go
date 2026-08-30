package theme

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/settings"
)

func TestClassifyThemeAssetFilename(t *testing.T) {
	t.Parallel()

	stylesheet, err := classify("community.CSS")
	require.NoError(t, err)
	assert.Equal(t, AssetKindStylesheet, stylesheet)

	script, err := classify("enhancements.js")
	require.NoError(t, err)
	assert.Equal(t, AssetKindScript, script)

	for _, invalid := range []string{"../theme.css", "folder/theme.css", `folder\\theme.js`, ".", "theme.css.exe", ""} {
		invalid := invalid
		t.Run(invalid, func(t *testing.T) {
			_, err := classify(invalid)
			assert.Error(t, err)
		})
	}
}

func TestRevisionPreservesKindAndOrder(t *testing.T) {
	t.Parallel()

	stylesheets := []Asset{{ID: "css-a"}, {ID: "css-b"}}
	scripts := []Asset{{ID: "js-a"}}

	assert.Equal(t, "003adffd272ebc71c6fdf76ce249c7a49ccce3b6aa2d90b41d46b2c86c7e31c5", revisionFor(stylesheets, scripts))
	assert.NotEqual(t, revisionFor(stylesheets, scripts), revisionFor([]Asset{{ID: "css-b"}, {ID: "css-a"}}, scripts))
	assert.Empty(t, revisionFor(nil, nil))
}

func TestStoredManifestAssetsUsePublicImmutablePath(t *testing.T) {
	t.Parallel()

	items := mapStoredAssets([]settings.ThemeAsset{{
		ID: "asset-id", Filename: "theme.css", MIMEType: "text/css", Size: 42, Integrity: "sha256-dGVzdA==",
	}}, AssetKindStylesheet)

	require.Len(t, items, 1)
	assert.Equal(t, "/api/info/theme/assets/theme.css", items[0].Path)
	assert.Equal(t, "text/css", items[0].MIMEType)
	assert.Equal(t, AssetKindStylesheet, items[0].Kind)
}
