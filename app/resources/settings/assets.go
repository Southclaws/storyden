package settings

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/internal/ent"
	ent_asset "github.com/Southclaws/storyden/internal/ent/asset"
)

type SettingsWithAssets struct {
	*Settings
	BannerAsset opt.Optional[asset.Asset]
}

func (d *SettingsRepository) GetWithAssets(ctx context.Context) (*SettingsWithAssets, error) {
	settings, err := d.Get(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := &SettingsWithAssets{Settings: settings}
	if id, ok := settings.BannerAssetID.Get(); ok {
		row, err := d.db.Asset.Query().Where(ent_asset.ID(id)).WithParent().Only(ctx)
		if err != nil {
			if !ent.IsNotFound(err) {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
		} else {
			result.BannerAsset = opt.New(*asset.Map(row))
		}
	}

	return result, nil
}
