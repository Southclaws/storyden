package asset_querier

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/internal/ent"
	ent_asset "github.com/Southclaws/storyden/internal/ent/asset"
)

type Querier struct {
	db *ent.Client
}

func New(db *ent.Client) *Querier {
	return &Querier{db}
}

func (q *Querier) Get(ctx context.Context, id asset.Filename) (*asset.Asset, error) {
	r, err := q.db.Asset.Query().Where(
		ent_asset.Filename(id.String()),
	).First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return asset.Map(r), nil
}

func (q *Querier) GetByID(ctx context.Context, id asset.AssetID) (*asset.Asset, error) {
	r, err := q.db.Asset.Query().Where(
		ent_asset.ID(id),
	).First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return asset.Map(r), nil
}
