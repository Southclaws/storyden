package asset_writer

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/internal/ent"
	ent_asset "github.com/Southclaws/storyden/internal/ent/asset"
	"github.com/Southclaws/storyden/internal/mime"
)

type Writer struct {
	db *ent.Client
}

func New(db *ent.Client) *Writer {
	return &Writer{db}
}

type Option func(*ent.AssetCreate)

func WithMetadata(metadata asset.Metadata) Option {
	return func(create *ent.AssetCreate) {
		create.SetMetadata(metadata)
	}
}

func (w *Writer) Add(ctx context.Context,
	accountID xid.ID,
	filename asset.Filename,
	size int,
	mt mime.Type,
	opts ...Option,
) (*asset.Asset, error) {
	create := w.db.Asset.
		Create().
		SetID(filename.GetID()).
		SetFilename(filename.String()).
		SetSize(size).
		SetMimeType(mt.String()).
		SetAccountID(xid.ID(accountID))

	for _, opt := range opts {
		opt(create)
	}

	r, err := create.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return asset.Map(r), nil
}

func (w *Writer) AddVersion(ctx context.Context,
	accountID xid.ID,
	filename asset.Filename,
	size int,
	mt mime.Type,
	parent asset.AssetID,
) (*asset.Asset, error) {
	r, err := w.db.Asset.
		Create().
		SetID(filename.GetID()).
		SetParentAssetID(parent).
		SetFilename(filename.String()).
		SetSize(size).
		SetMimeType(mt.String()).
		SetAccountID(xid.ID(accountID)).
		Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err = w.db.Asset.Query().Where(ent_asset.ID(r.ID)).WithParent().First(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return asset.Map(r), nil
}

func (w *Writer) Remove(ctx context.Context, accountID xid.ID, id asset.Filename) error {
	q := w.db.Asset.
		Delete().Where(
		ent_asset.Filename(id.String()),
	)

	if _, err := q.Exec(ctx); err != nil {
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return nil
}
