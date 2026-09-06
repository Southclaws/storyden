package bindings

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Southclaws/storyden/app/resources/settings"
)

func TestSerialiseThemeAssetUsesPublicImmutablePath(t *testing.T) {
	t.Parallel()

	got := serialiseThemeAsset(settings.ThemeAsset{
		ID:        "asset-id",
		Filename:  "theme.css",
		MIMEType:  "text/css",
		Size:      42,
		Integrity: "sha256-dGVzdA==",
	})

	assert.Equal(t, "/api/info/theme/assets/theme.css", got.Path)
	assert.Equal(t, "text/css", string(got.MimeType))
}
