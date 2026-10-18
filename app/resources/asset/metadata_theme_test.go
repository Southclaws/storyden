package asset_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/internal/mime"
)

func TestThemeMetadata(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		mediaType string
		kind      asset.ThemeAssetKind
		raw       string
	}{
		{
			mediaType: "text/css",
			kind:      asset.ThemeAssetKindStylesheet,
			raw:       `{"purpose":"theme","theme_kind":"stylesheet","theme_integrity":"sha256-example"}`,
		},
		{
			mediaType: "application/javascript",
			kind:      asset.ThemeAssetKindScript,
			raw:       `{"purpose":"theme","theme_kind":"script","theme_integrity":"sha256-example"}`,
		},
		{
			mediaType: "text/javascript",
			kind:      asset.ThemeAssetKindScript,
			raw:       `{"purpose":"theme","theme_kind":"script","theme_integrity":"sha256-example"}`,
		},
	} {
		t.Run(test.mediaType, func(t *testing.T) {
			source := asset.Asset{
				MIME:     mime.New(test.mediaType),
				Metadata: asset.Metadata{RawMessage: json.RawMessage(test.raw)},
			}

			metadata, err := source.ThemeMetadata()
			require.NoError(t, err)
			require.Equal(t, &asset.ThemeMetadata{
				Purpose:   "theme",
				Kind:      test.kind,
				Integrity: "sha256-example",
			}, metadata)

			encoded, err := asset.NewMetadata(source.MIME, metadata)
			require.NoError(t, err)
			require.JSONEq(t, test.raw, string(encoded.RawMessage))

			image, err := source.ImageMetadata()
			require.NoError(t, err)
			require.Nil(t, image)
		})
	}
}

func TestThemeMetadataUsesMIMEType(t *testing.T) {
	t.Parallel()

	for _, mediaType := range []string{"image/png", "application/pdf", "text/plain", "application/zip"} {
		t.Run(mediaType, func(t *testing.T) {
			source := asset.Asset{
				MIME: mime.New(mediaType),
				Metadata: asset.Metadata{RawMessage: json.RawMessage(
					`{"purpose":"theme","theme_kind":"stylesheet","theme_integrity":"sha256-example"}`,
				)},
			}

			metadata, err := source.ThemeMetadata()
			require.NoError(t, err)
			require.Nil(t, metadata)
		})
	}
}

func TestThemeMetadataRejectsInvalidFieldTypes(t *testing.T) {
	t.Parallel()

	for _, raw := range []string{
		`{"purpose":123}`,
		`{"theme_kind":123}`,
		`{"theme_integrity":123}`,
	} {
		t.Run(raw, func(t *testing.T) {
			source := asset.Asset{
				MIME:     mime.New("text/css"),
				Metadata: asset.Metadata{RawMessage: json.RawMessage(raw)},
			}

			_, err := source.ThemeMetadata()
			require.Error(t, err)

			_, err = asset.NewMetadata(source.MIME, json.RawMessage(raw))
			require.Error(t, err)
		})
	}
}
