package collection_read

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/app/resources/collection"
	"github.com/Southclaws/storyden/app/services/semdex"
)

type CollectionReader struct {
	fx.In

	Logger *zap.Logger
	Repo   collection.Repository
	Semdex semdex.Retriever
}

func (r *CollectionReader) IndexCollection(ctx context.Context, id collection.CollectionID) (*collection.Collection, error) {
	col, err := r.Repo.Get(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// ids := dt.Map(col.Items, func(i *collection.CollectionItem) xid.ID { return i.Item.GetID() })

	// vector, err := r.Semdex.GetVectorFor(ctx, ids...)
	// if err != nil {
	// 	r.Logger.Warn("failed to get vector for collection items", zap.Error(err))
	// }

	// TODO: dispatch indexing request for this vector

	return col, nil
}
