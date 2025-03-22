package weaviate_semdexer

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/filters"
)

func (w *weaviateSemdexer) Delete(ctx context.Context, id xid.ID) (int, error) {
	delete := w.wc.Batch().
		ObjectsBatchDeleter().
		WithClassName(w.cn.String()).
		WithWhere(
			filters.Where().
				WithPath([]string{"datagraph_id"}).
				WithOperator(filters.Equal).
				WithValueString(id.String()),
		)

	r, err := delete.Do(ctx)
	if err != nil {
		return 0, fault.Wrap(err, fctx.With(ctx))
	}

	return int(r.Results.Successful), nil
}
