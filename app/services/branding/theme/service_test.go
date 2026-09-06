package theme

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
