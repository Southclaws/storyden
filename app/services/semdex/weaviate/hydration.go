package weaviate

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/app/services/semdex/result_hydrator"
)

type withHydration struct {
	l  *zap.Logger
	wc semdex.Semdexer
	rh *result_hydrator.Hydrator
}

func (h *withHydration) Index(ctx context.Context, object datagraph.Indexable) error {
	return h.wc.Index(ctx, object)
}

func (h *withHydration) Search(ctx context.Context, query string) ([]*semdex.Result, error) {
	rs, err := h.wc.Search(ctx, query)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	results := []*semdex.Result{}

	// NOTE: Should probably be parallelised at some point...
	for _, v := range rs {
		r, err := h.rh.Hydrate(ctx, v)
		if err != nil {
			h.l.Warn("failed to hydrate search result",
				zap.String("datagraph_kind", v.Type.String()),
				zap.String("result_id", v.Id.String()),
				zap.Error(err))
			continue
		}

		results = append(results, r)
	}

	return results, nil
}

func (h *withHydration) Recommend(ctx context.Context, object datagraph.Indexable) ([]*semdex.Result, error) {
	rs, err := h.wc.Recommend(ctx, object)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	results, err := h.Hydrate(ctx, rs)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return results, nil
}

func (h *withHydration) Hydrate(ctx context.Context, results []*semdex.Result) ([]*semdex.Result, error) {
	// NOTE: Should probably be parallelised at some point...
	for i, v := range results {
		r, err := h.rh.Hydrate(ctx, v)
		if err != nil {
			h.l.Warn("failed to hydrate search result",
				zap.String("datagraph_kind", v.Type.String()),
				zap.String("result_id", v.Id.String()),
				zap.Error(err))
			continue
		}

		results[i] = r
	}

	return results, nil
}
