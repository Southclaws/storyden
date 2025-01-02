package pinecone_semdexer

import (
	"context"
	"sort"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/services/search/searcher"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/infrastructure/vector/pinecone"
)

func (s *pineconeSemdexer) Search(ctx context.Context, q string, p pagination.Parameters, opts searcher.Options) (*pagination.Result[datagraph.Item], error) {
	refs, err := s.SearchRefs(ctx, q, p, opts)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	items, err := s.hydrator.Hydrate(ctx, refs.Items...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := pagination.NewPageResult(p, refs.Results, items)
	return &result, nil
}

func (s *pineconeSemdexer) SearchRefs(ctx context.Context, q string, p pagination.Parameters, opts searcher.Options) (*pagination.Result[*datagraph.Ref], error) {
	objects, err := s.searchObjects(ctx, q, p, opts)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	results := objects.ToRefs()

	deduped := dedupeChunks(results)

	filtered := filterChunks(deduped)

	pagedResult := pagination.NewPageResult(p, len(results), filtered)

	return &pagedResult, nil
}

func (s *pineconeSemdexer) SearchChunks(ctx context.Context, q string, p pagination.Parameters, opts searcher.Options) ([]*semdex.Chunk, error) {
	objects, err := s.searchObjects(ctx, q, p, opts)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return objects.ToChunks(), nil
}

func (s *pineconeSemdexer) searchObjects(ctx context.Context, q string, p pagination.Parameters, opts searcher.Options) (Objects, error) {
	vec, err := s.ef(ctx, q)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	filterMap := map[string]any{}

	opts.Kinds.Call(func(kind []datagraph.Kind) {
		filterMap["datagraph_type"] = map[string]any{
			"$in": dt.Map(kind, func(k datagraph.Kind) any { return k.String() }),
		}
	})

	filter, err := structpb.NewStruct(filterMap)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	response, err := s.index.QueryByVectorValues(ctx, &pinecone.QueryByVectorValuesRequest{
		Vector:          vec,
		TopK:            uint32(p.Limit()),
		MetadataFilter:  filter,
		IncludeValues:   false,
		IncludeMetadata: true,
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return mapScoredVectors(response.Matches)
}

func filterChunks(results []*datagraph.Ref) []*datagraph.Ref {
	filtered := dt.Filter(results, func(r *datagraph.Ref) bool {
		return r.Relevance > 0.5
	})

	return filtered
}

func dedupeChunks(results []*datagraph.Ref) []*datagraph.Ref {
	groupedByID := lo.GroupBy(results, func(r *datagraph.Ref) xid.ID { return r.ID })

	// for each grouped result, compute the average score and flatten
	// the list of results into a single result per ID
	// this is a naive approach to deduplication

	list := lo.Values(groupedByID)

	deduped := dt.Reduce(list, func(acc []*datagraph.Ref, curr []*datagraph.Ref) []*datagraph.Ref {
		first := curr[0]
		score := []float64{}

		for _, r := range curr {
			score = append(score, r.Relevance)
		}

		next := &datagraph.Ref{
			ID:        first.ID,
			Kind:      first.Kind,
			Relevance: maxFloat64(score...),
		}

		return append(acc, next)
	}, []*datagraph.Ref{})

	sort.Sort(datagraph.RefList(deduped))

	return deduped
}

// max of all input floats
func maxFloat64(a ...float64) float64 {
	max := a[0]
	for _, n := range a {
		if n > max {
			max = n
		}
	}
	return max
}
