package weaviate_semdexer

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
)

func (w *weaviateSemdexer) Delete(ctx context.Context, id xid.ID) error {
	wid := GetWeaviateID(id)

	err := w.wc.Data().Deleter().
		WithClassName(string(w.cn)).
		WithID(wid).
		Do(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
