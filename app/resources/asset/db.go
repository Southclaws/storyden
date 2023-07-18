package asset

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
)

type database struct {
	db *ent.Client
}

func New(db *ent.Client) Repository {
	return &database{db}
}

func (d *database) Add(ctx context.Context,
	accountID account.AccountID,
	id, url, mt string,
	width, height int,
) (*Asset, error) {
	asset, err := d.db.Asset.Get(ctx, id)
	if err != nil && !ent.IsNotFound(err) {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if asset == nil {
		asset, err = d.db.Asset.
			Create().
			SetID(id).
			SetURL(url).
			SetWidth(width).
			SetHeight(height).
			SetMimetype(mt).
			SetAccountID(xid.ID(accountID)).
			Save(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
		}
	}

	return FromModel(asset), nil
}

func (d *database) Get(ctx context.Context, id string) (*Asset, error) {
	asset, err := d.db.Asset.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return FromModel(asset), nil
}

func (d *database) Remove(ctx context.Context, accountID account.AccountID, id AssetID) error {
	q := d.db.Asset.
		DeleteOneID(string(id))

	if err := q.Exec(ctx); err != nil {
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return nil
}
