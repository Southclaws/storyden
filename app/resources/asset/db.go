package asset

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/asset"
)

type database struct {
	db *ent.Client
}

func New(db *ent.Client) Repository {
	return &database{db}
}

func (d *database) Add(ctx context.Context,
	accountID account.AccountID,
	filename Filename,
	size int,
) (*Asset, error) {
	asset, err := d.db.Asset.
		Create().
		SetID(filename.GetID()).
		SetFilename(filename.name).
		SetSize(size).
		SetAccountID(xid.ID(accountID)).
		Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return FromModel(asset), nil
}

func (d *database) Get(ctx context.Context, id Filename) (*Asset, error) {
	asset, err := d.db.Asset.Query().Where(
		asset.Filename(id.name),
	).First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return FromModel(asset), nil
}

func (d *database) GetByID(ctx context.Context, id AssetID) (*Asset, error) {
	asset, err := d.db.Asset.Query().Where(
		asset.ID(id),
	).First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return FromModel(asset), nil
}

func (d *database) Remove(ctx context.Context, accountID account.AccountID, id Filename) error {
	q := d.db.Asset.
		Delete().Where(
		asset.Filename(id.name),
	)

	if _, err := q.Exec(ctx); err != nil {
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return nil
}
