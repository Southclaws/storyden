package robot

import (
	"fmt"

	"github.com/rs/xid"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/asset"
)

const ImageAssetIDMetadataKey = "storyden.asset.image_id"

func NewImageAssetPart(id asset.AssetID) *genai.Part {
	return &genai.Part{
		PartMetadata: map[string]any{
			ImageAssetIDMetadataKey: id.String(),
		},
	}
}

func ImageAssetIDs(contents ...*genai.Content) ([]asset.AssetID, error) {
	ids := []asset.AssetID{}

	for _, content := range contents {
		if content == nil {
			continue
		}
		for _, part := range content.Parts {
			id, ok, err := imageAssetID(part)
			if err != nil {
				return nil, err
			}
			if !ok {
				continue
			}
			ids = append(ids, id)
		}
	}

	return ids, nil
}

func imageAssetID(part *genai.Part) (asset.AssetID, bool, error) {
	if part == nil {
		return asset.AssetID{}, false, nil
	}
	if part.FileData != nil || part.InlineData != nil {
		return asset.AssetID{}, false, fmt.Errorf("model media must reference a Storyden asset ID")
	}

	raw, ok := part.PartMetadata[ImageAssetIDMetadataKey]
	if !ok {
		return asset.AssetID{}, false, nil
	}
	value, ok := raw.(string)
	if !ok {
		return asset.AssetID{}, false, fmt.Errorf("image asset ID metadata must be a string")
	}
	id, err := xid.FromString(value)
	if err != nil {
		return asset.AssetID{}, false, fmt.Errorf("invalid image asset ID %q: %w", value, err)
	}
	return asset.AssetID(id), true, nil
}

func UniqueImageAssetIDs(contents ...*genai.Content) ([]asset.AssetID, error) {
	ids, err := ImageAssetIDs(contents...)
	if err != nil {
		return nil, err
	}

	unique := make([]asset.AssetID, 0, len(ids))
	seen := make(map[asset.AssetID]struct{}, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		unique = append(unique, id)
	}
	return unique, nil
}
