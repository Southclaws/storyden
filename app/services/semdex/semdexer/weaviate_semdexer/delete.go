package weaviate_semdexer

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
)

func (w *weaviateSemdexer) Delete(ctx context.Context, id xid.ID) error {
	delete := w.wc.Batch().
		ObjectsBatchDeleter().
		WithWhere(
			filters.Where().
				WithPath([]string{"datagraph_id"}).
				WithOperator(filters.Equal).
				WithValueString(id.String()),
		)

	_, err := delete.Do(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
