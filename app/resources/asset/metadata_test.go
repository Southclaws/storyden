package asset_test

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/internal/mime"
)

func TestImageMetadata(t *testing.T) {
	t.Parallel()

	source := asset.Asset{
		MIME: mime.New("image/png"),
		Metadata: asset.Metadata{RawMessage: json.RawMessage(
			`{"width":960,"height":480,"image_preset":"m","image_recipe":1}`,
		)},
	}

	metadata, err := source.ImageMetadata()
	require.NoError(t, err)
	require.Equal(t, &asset.ImageMetadata{
		Width:  960,
		Height: 480,
		Preset: asset.Medium,
		Recipe: 1,
	}, metadata)

	encoded, err := json.Marshal(metadata)
	require.NoError(t, err)
	require.JSONEq(t, string(source.Metadata.RawMessage), string(encoded))
}

func TestImageMetadataUsesMIMEType(t *testing.T) {
	t.Parallel()

	metadata, err := asset.NewMetadata(mime.New("application/pdf"), struct {
		Width string `json:"width"`
		Pages int    `json:"pages"`
	}{Width: "custom value", Pages: 12})
	require.NoError(t, err)

	source := asset.Asset{MIME: mime.New("application/pdf"), Metadata: metadata}
	image, err := source.ImageMetadata()
	require.NoError(t, err)
	require.Nil(t, image)
	require.JSONEq(t, `{"width":"custom value","pages":12}`, string(metadata.RawMessage))
}

func TestImageMetadataRejectsInvalidFieldTypes(t *testing.T) {
	t.Parallel()

	for _, raw := range []string{
		`{"width":"960"}`,
		`{"height":1.5}`,
		`{"image_recipe":"1"}`,
	} {
		t.Run(raw, func(t *testing.T) {
			source := asset.Asset{
				MIME:     mime.New("image/png"),
				Metadata: asset.Metadata{RawMessage: json.RawMessage(raw)},
			}
			_, err := source.ImageMetadata()
			require.Error(t, err)

			_, err = asset.NewMetadata(source.MIME, json.RawMessage(raw))
			require.Error(t, err)
		})
	}
}

func TestMetadataMergePreservesAdditionalFields(t *testing.T) {
	t.Parallel()

	metadata := asset.Metadata{RawMessage: json.RawMessage(
		`{"width":1600,"height":800,"image_preset":"m","image_recipe":1,"counter":18446744073709551615}`,
	)}
	updated, err := metadata.Merge(asset.ImageMetadata{Width: 320, Height: 160})
	require.NoError(t, err)

	var stored struct {
		asset.ImageMetadata
		Counter uint64 `json:"counter"`
	}
	require.NoError(t, updated.Decode(&stored))
	require.Equal(t, 320, stored.Width)
	require.Equal(t, 160, stored.Height)
	require.Equal(t, asset.Medium, stored.Preset)
	require.Equal(t, 1, stored.Recipe)
	require.Equal(t, uint64(math.MaxUint64), stored.Counter)
}
