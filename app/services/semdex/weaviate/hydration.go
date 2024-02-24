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

func (h *withHydration) Search(ctx context.Context, query string) (datagraph.NodeReferenceList, error) {
	rs, err := h.wc.Search(ctx, query)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	results := datagraph.NodeReferenceList{}

	// NOTE: Should probably be parallelised at some point...
	for _, v := range rs {
		r, err := h.rh.Hydrate(ctx, v)
		if err != nil {
			h.l.Warn("failed to hydrate search result",
				zap.String("datagraph_kind", v.Kind.String()),
				zap.String("result_id", v.ID.String()),
				zap.Error(err))
			continue
		}

		results = append(results, r)
	}

	return results, nil
}

func (h *withHydration) Recommend(ctx context.Context, object datagraph.Indexable) (datagraph.NodeReferenceList, error) {
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

func (h *withHydration) Hydrate(ctx context.Context, results datagraph.NodeReferenceList) (datagraph.NodeReferenceList, error) {
	hydrated := datagraph.NodeReferenceList{}
	// NOTE: Should probably be parallelised at some point...
	for _, v := range results {
		r, err := h.rh.Hydrate(ctx, v)
		if err != nil {
			h.l.Warn("failed to hydrate search result",
				zap.String("datagraph_kind", v.Kind.String()),
				zap.String("result_id", v.ID.String()),
				zap.Error(err))
			continue
		}

		hydrated = append(hydrated, r)
	}

	return hydrated, nil
}
