package asset_download

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/internal/mime"
)

func TestSelectVariant(t *testing.T) {
	t.Parallel()

	source := &asset.Asset{MIME: mime.New("image/png")}
	variants := make(map[asset.Variant]*asset.Asset)
	for _, preset := range []asset.Variant{
		asset.Small,
		asset.Medium,
		asset.Large,
		asset.ThumbnailSmall,
		asset.ThumbnailMedium,
		asset.ThumbnailLarge,
	} {
		mediaType := mime.New("image/webp")
		metadata, err := asset.NewMetadata(mediaType, asset.ImageMetadata{Preset: preset})
		require.NoError(t, err)

		variants[preset] = &asset.Asset{MIME: mediaType, Metadata: metadata}
	}

	for _, tc := range []struct {
		name      string
		target    asset.Variant
		available []asset.Variant
		expected  *asset.Asset
	}{
		{
			name:     "no variants returns original",
			target:   asset.Large,
			expected: source,
		},
		{
			name:      "medium request skips small fallback",
			target:    asset.Medium,
			available: []asset.Variant{asset.Small},
			expected:  source,
		},
		{
			name:      "large request skips all smaller fallbacks",
			target:    asset.Large,
			available: []asset.Variant{asset.Small, asset.Medium},
			expected:  source,
		},
		{
			name:      "exact fitted variant wins",
			target:    asset.Medium,
			available: []asset.Variant{asset.Small, asset.Large, asset.Medium},
			expected:  variants[asset.Medium],
		},
		{
			name:      "smallest sufficient fitted fallback wins",
			target:    asset.Small,
			available: []asset.Variant{asset.Large, asset.Medium},
			expected:  variants[asset.Medium],
		},
		{
			name:      "thumbnail keeps largest smaller fallback",
			target:    asset.ThumbnailLarge,
			available: []asset.Variant{asset.ThumbnailMedium, asset.ThumbnailSmall},
			expected:  variants[asset.ThumbnailMedium],
		},
		{
			name:      "thumbnail keeps smallest sufficient fallback",
			target:    asset.ThumbnailSmall,
			available: []asset.Variant{asset.ThumbnailLarge, asset.ThumbnailMedium},
			expected:  variants[asset.ThumbnailMedium],
		},
		{
			name:      "thumbnail does not select fitted fallback",
			target:    asset.ThumbnailLarge,
			available: []asset.Variant{asset.Large},
			expected:  source,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			preset, ok := asset.Preset(tc.target)
			require.True(t, ok)

			var available []*asset.Asset
			for _, variant := range tc.available {
				available = append(available, variants[variant])
			}

			selected, err := selectVariant(source, available, preset)
			require.NoError(t, err)
			require.Same(t, tc.expected, selected)
		})
	}
}
