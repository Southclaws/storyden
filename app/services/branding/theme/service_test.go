package theme

import (
	"testing"

	"github.com/Southclaws/opt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseManifest(t *testing.T) {
	t.Run("defaults when metadata is missing", func(t *testing.T) {
		out, err := ParseManifest(opt.NewEmpty[map[string]any]())
		require.NoError(t, err)
		assert.Empty(t, out.CSS)
		assert.Empty(t, out.Scripts)
	})

	t.Run("normalises plain filenames into theme asset paths", func(t *testing.T) {
		out, err := ParseManifest(opt.New(map[string]any{
			"theme": map[string]any{
				"css":     []any{"my-theme.css", "/api/info/theme/assets/existing-css"},
				"scripts": []any{"theme-hooks.js"},
			},
		}))
		require.NoError(t, err)

		assert.Equal(t,
			[]string{
				"/api/info/theme/assets/my-theme-css",
				"/api/info/theme/assets/existing-css",
			},
			out.CSS,
		)
		assert.Equal(t, []string{"/api/info/theme/assets/theme-hooks-js"}, out.Scripts)
	})

	t.Run("deduplicates while preserving order", func(t *testing.T) {
		out, err := ParseManifest(opt.New(map[string]any{
			"theme": map[string]any{
				"css": []any{
					"/api/info/theme/assets/one",
					"/api/info/theme/assets/one",
					"/api/info/theme/assets/two",
				},
			},
		}))
		require.NoError(t, err)
		assert.Equal(t, []string{
			"/api/info/theme/assets/one",
			"/api/info/theme/assets/two",
		}, out.CSS)
	})
}

func TestClassifyAssetKind(t *testing.T) {
	t.Run("accepts stylesheet and script extensions", func(t *testing.T) {
		kind, err := classifyAssetKind("custom.css")
		require.NoError(t, err)
		assert.Equal(t, AssetKindStylesheet, kind)

		kind, err = classifyAssetKind("behaviour.js")
		require.NoError(t, err)
		assert.Equal(t, AssetKindScript, kind)
	})

	t.Run("rejects other extensions", func(t *testing.T) {
		_, err := classifyAssetKind("bad.png")
		require.Error(t, err)
	})
}

func TestValidateDetectedMIME(t *testing.T) {
	t.Run("allows css mimes for stylesheets", func(t *testing.T) {
		require.NoError(t, validateDetectedMIME(AssetKindStylesheet, "text/css"))
		require.NoError(t, validateDetectedMIME(AssetKindStylesheet, "text/plain; charset=utf-8"))
	})

	t.Run("allows js mimes for scripts", func(t *testing.T) {
		require.NoError(t, validateDetectedMIME(AssetKindScript, "application/javascript"))
		require.NoError(t, validateDetectedMIME(AssetKindScript, "text/javascript"))
		require.NoError(t, validateDetectedMIME(AssetKindScript, "text/plain"))
	})

	t.Run("rejects non-text mime", func(t *testing.T) {
		require.Error(t, validateDetectedMIME(AssetKindStylesheet, "image/png"))
		require.Error(t, validateDetectedMIME(AssetKindScript, "application/zip"))
	})
}
