package asset_writer

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/asset/asset_image"
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

type createOptions struct {
	metadata     any
	imagePresets int64
}

type Option func(*createOptions)

func WithMetadata(metadata any) Option {
	return func(options *createOptions) {
		options.metadata = metadata
	}
}

func WithImagePresets(presets int64) Option {
	return func(options *createOptions) {
		options.imagePresets = presets
	}
}

func (w *Writer) Add(
	ctx context.Context,
	accountID xid.ID,
	filename asset.Filename,
	size int,
	mediaType mime.Type,
	opts ...Option,
) (*asset.Asset, error) {
	return w.add(ctx, accountID, filename, size, mediaType, opt.NewEmpty[asset.AssetID](), opts...)
}

func (w *Writer) AddVersion(
	ctx context.Context,
	accountID xid.ID,
	filename asset.Filename,
	size int,
	mediaType mime.Type,
	parent asset.AssetID,
	opts ...Option,
) (*asset.Asset, error) {
	stored, err := w.add(ctx, accountID, filename, size, mediaType, opt.New(parent), opts...)
	if err != nil {
		return nil, err
	}

	row, err := w.db.Asset.Query().Where(ent_asset.ID(stored.ID)).WithParent().First(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return asset.Map(row), nil
}

func (w *Writer) add(
	ctx context.Context,
	accountID xid.ID,
	filename asset.Filename,
	size int,
	mediaType mime.Type,
	parent opt.Optional[asset.AssetID],
	opts ...Option,
) (*asset.Asset, error) {
	var options createOptions
	for _, option := range opts {
		option(&options)
	}

	metadata, err := asset.NewMetadata(mediaType, options.metadata)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	var row *ent.Asset
	err = ent.WithTx(ctx, w.db, func(tx *ent.Tx) error {
		create := tx.Asset.Create().
			SetID(filename.GetID()).
			SetFilename(filename.String()).
			SetSize(size).
			SetMimeType(mediaType.String()).
			SetAccountID(accountID)

		if parentID, ok := parent.Get(); ok {
			create.SetParentAssetID(parentID)
		}
		if options.metadata != nil {
			create.SetMetadata(metadata.RawMessage)
		}

		var err error
		row, err = create.Save(ctx)
		if err != nil {
			return err
		}
		if options.imagePresets == 0 {
			return nil
		}

		return asset_image.Request(ctx, tx.Client(), row.ID, options.imagePresets)
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return asset.Map(row), nil
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
