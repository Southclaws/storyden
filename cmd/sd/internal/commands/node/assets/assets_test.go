package assets

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func TestFindAsset(t *testing.T) {
	assets := openapi.AssetList{
		{
			Id:       "asset-id",
			Filename: "cover.png",
			Path:     "/api/assets/stored-cover.png",
		},
	}

	tests := []string{
		"asset-id",
		"cover.png",
		"/api/assets/stored-cover.png",
		"stored-cover.png",
	}

	for _, selector := range tests {
		t.Run(selector, func(t *testing.T) {
			r := require.New(t)

			asset, err := findAsset(assets, selector)
			r.NoError(err)
			r.Equal("asset-id", asset.Id)
		})
	}
}

func TestFindAssetMissing(t *testing.T) {
	r := require.New(t)

	_, err := findAsset(nil, "missing")
	r.ErrorContains(err, "asset not attached to node: missing")
}

func TestAssetFilename(t *testing.T) {
	r := require.New(t)

	r.Equal("stored-cover.png", assetFilename(openapi.Asset{
		Filename: "cover.png",
		Path:     "/api/assets/stored-cover.png",
	}))
	r.Equal("cover.png", assetFilename(openapi.Asset{
		Filename: "cover.png",
	}))
}
