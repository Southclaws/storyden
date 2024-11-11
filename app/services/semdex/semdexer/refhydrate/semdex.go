package refhydrate

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/services/search/searcher"
	"github.com/Southclaws/storyden/app/services/semdex"
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

func (h *HydratedSemdexer) Delete(ctx context.Context, id xid.ID) error {
	return h.RefSemdex.Delete(ctx, id)
}

func (h *HydratedSemdexer) Search(ctx context.Context, query string, p pagination.Parameters, opts searcher.Options) (*pagination.Result[datagraph.Item], error) {
	rs, err := h.RefSemdex.Search(ctx, query, p, opts)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	hydrated, err := h.Hydrator.Hydrate(ctx, rs.Items...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	results := pagination.ConvertPageResult[*datagraph.Ref, datagraph.Item](*rs, hydrated)

	return &results, nil
}

func (h *HydratedSemdexer) Recommend(ctx context.Context, object datagraph.Item) (datagraph.ItemList, error) {
	rs, err := h.RefSemdex.Recommend(ctx, object)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return h.Hydrator.Hydrate(ctx, rs...)
}

func (h *HydratedSemdexer) SuggestTags(ctx context.Context, content datagraph.Content, available tag_ref.Names) (tag_ref.Names, error) {
	return h.RefSemdex.SuggestTags(ctx, content, available)
}

func (h *HydratedSemdexer) GetMany(ctx context.Context, limit uint, ids ...xid.ID) (datagraph.RefList, error) {
	return h.RefSemdex.GetMany(ctx, limit, ids...)
}

func (h *HydratedSemdexer) ScoreRelevance(ctx context.Context, object datagraph.Item, idx ...xid.ID) (map[xid.ID]float64, error) {
	return h.RefSemdex.ScoreRelevance(ctx, object, idx...)
}

func (h *HydratedSemdexer) Summarise(ctx context.Context, object datagraph.Item) (string, error) {
	return h.RefSemdex.Summarise(ctx, object)
}
