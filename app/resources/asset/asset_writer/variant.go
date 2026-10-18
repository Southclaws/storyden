package asset_writer

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/internal/ent"
)

// AddVariant is for use with background job 2-phase commit transaction API.
func AddVariant(ctx context.Context, db *ent.Client, source, variant *asset.Asset) error {
	return db.Asset.Create().
		SetID(variant.Name.GetID()).
		SetFilename(variant.Name.String()).
		SetAccountID(source.OwnerID).
		SetParentAssetID(source.ID).
		SetSize(variant.Size).
		SetMimeType(variant.MIME.String()).
		SetMetadata(variant.Metadata.RawMessage).
		Exec(ctx)
}
