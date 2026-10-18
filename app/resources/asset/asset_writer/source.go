package asset_writer

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/internal/ent"
)

func AddSource(ctx context.Context, db *ent.Client, source *asset.Asset) error {
	return db.Asset.Create().
		SetID(source.ID).
		SetFilename(source.Name.String()).
		SetAccountID(source.OwnerID).
		SetSize(source.Size).
		SetMimeType(source.MIME.String()).
		Exec(ctx)
}
