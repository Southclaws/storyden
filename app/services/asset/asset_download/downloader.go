package asset_download

import (
	"context"
	"io"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/asset/asset_querier"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
)

type Downloader struct {
	assets  *asset_querier.Querier
	objects object.Storer
}

func New(
	assets *asset_querier.Querier,
	objects object.Storer,
) *Downloader {
	return &Downloader{
		assets:  assets,
		objects: objects,
	}
}

func (d *Downloader) Get(ctx context.Context, id asset.Filename) (*asset.Asset, io.Reader, error) {
	a, err := d.assets.Get(ctx, id)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	path := asset.BuildAssetPath(a.Name)
	ctx = fctx.WithMeta(ctx, "path", path, "asset_id", id.String())

	r, size, err := d.objects.Read(ctx, path)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	a.Size = int(size)

	return a, r, nil
}

func (d *Downloader) GetByID(ctx context.Context, id asset.AssetID) (*asset.Asset, io.Reader, error) {
	a, err := d.assets.GetByID(ctx, id)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	path := asset.BuildAssetPath(a.Name)
	ctx = fctx.WithMeta(ctx, "path", path, "asset_id", id.String())

	r, size, err := d.objects.Read(ctx, path)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	a.Size = int(size)

	return a, r, nil
}
