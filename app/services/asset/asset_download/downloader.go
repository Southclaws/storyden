package asset_download

import (
	"context"
	"io"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
)

type Downloader struct {
	l    *zap.Logger
	rbac rbac.AccessManager

	assets  asset.Repository
	objects object.Storer
}

func New(
	l *zap.Logger,
	rbac rbac.AccessManager,

	assets asset.Repository,
	objects object.Storer,
) *Downloader {
	return &Downloader{
		l:       l.With(zap.String("service", "asset")),
		rbac:    rbac,
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
