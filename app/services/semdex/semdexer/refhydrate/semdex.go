package refhydrate

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/datagraph/semdex"
)

var _ semdex.Semdexer = &HydratedSemdexer{}

// HydratedSemdexer implements the Semdexer interface for semantic indexing. It
// wraps the weaviate ref index which works on non-hydrated lower level refs.
type HydratedSemdexer struct {
	RefSemdex semdex.RefSemdexer
	Hydrator  *Hydrator
}

func (h *HydratedSemdexer) Index(ctx context.Context, object datagraph.Item) error {
	return h.RefSemdex.Index(ctx, object)
}

func (h *HydratedSemdexer) Search(ctx context.Context, query string) (datagraph.ItemList, error) {
	rs, err := h.RefSemdex.Search(ctx, query)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return h.Hydrator.Hydrate(ctx, rs...)
}

func (h *HydratedSemdexer) Recommend(ctx context.Context, object datagraph.Item) (datagraph.ItemList, error) {
	rs, err := h.RefSemdex.Recommend(ctx, object)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return h.Hydrator.Hydrate(ctx, rs...)
}

func (h *HydratedSemdexer) GetAll(ctx context.Context) (datagraph.RefList, error) {
	return h.RefSemdex.GetAll(ctx)
}

func (h *HydratedSemdexer) ScoreRelevance(ctx context.Context, object datagraph.Item, idx ...xid.ID) (map[xid.ID]float64, error) {
	return h.RefSemdex.ScoreRelevance(ctx, object, idx...)
}

func (h *HydratedSemdexer) Summarise(ctx context.Context, object datagraph.Item) (string, error) {
	return h.RefSemdex.Summarise(ctx, object)
}
