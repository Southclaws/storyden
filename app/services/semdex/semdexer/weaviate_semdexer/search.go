package weaviate_semdexer

import (
	"context"
	"fmt"
	"net/url"
	"sort"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/xid"
	"github.com/samber/lo"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/graphql"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/services/search/searcher"
	"github.com/Southclaws/storyden/app/services/semdex"
)

const (
	DefaultAutoCut     = 2
	VectorKeywordAlpha = 0.75
	RelevanceThreshold = 0.5
)

func (s *weaviateSemdexer) Search(ctx context.Context, q string, p pagination.Parameters, opts searcher.Options) (*pagination.Result[datagraph.Item], error) {
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

func (s *weaviateSemdexer) SearchRefs(ctx context.Context, q string, p pagination.Parameters, opts searcher.Options) (*pagination.Result[*datagraph.Ref], error) {
	objects, err := s.searchObjects(ctx, q, p, opts)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	results, err := dt.MapErr(objects, mapToNodeReference)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	filtered := filterChunks(results)

	deduped := dedupeChunks(filtered)

	pagedResult := pagination.NewPageResult(p, len(results), deduped)

	return &pagedResult, nil
}

func (s *weaviateSemdexer) SearchChunks(ctx context.Context, q string, p pagination.Parameters, opts searcher.Options) ([]*semdex.Chunk, error) {
	classData, err := s.searchObjects(ctx, q, p, opts)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return mapObjectsToChunks(classData)
}

func (s *weaviateSemdexer) searchObjects(ctx context.Context, q string, p pagination.Parameters, opts searcher.Options) ([]WeaviateObject, error) {
	fields := []graphql.Field{
		{Name: "datagraph_id"},
		{Name: "datagraph_type"},
		{Name: "name"},
		{Name: "content"},
		{Name: "_additional", Fields: []graphql.Field{
			{Name: "score"},
			{Name: "explainScore"},
		}},
	}

	autocut := DefaultAutoCut
	if p.PageZeroIndexed() > 0 {
		autocut = 0
	}

	arg := s.wc.GraphQL().
		HybridArgumentBuilder().
		WithAlpha(VectorKeywordAlpha).
		WithFusionType(graphql.RelativeScore).
		WithQuery(q)

	query := s.wc.GraphQL().Get().
		WithClassName(s.cn.String()).
		WithFields(fields...).
		WithHybrid(arg).
		WithAutocut(autocut).
		WithOffset(p.Offset()).
		WithLimit(p.Limit())

	countQuery := s.wc.GraphQL().Aggregate().
		WithClassName(s.cn.String()).
		WithFields(graphql.Field{
			Name:   "datagraph_id",
			Fields: []graphql.Field{{Name: "count"}},
		})

	if ks, ok := opts.Kinds.Get(); ok {
		kinds := dt.Map(ks, func(k datagraph.Kind) string { return k.String() })

		filter := filters.Where().
			WithPath([]string{"datagraph_type"}).
			WithOperator(filters.ContainsAny).
			WithValueString(kinds...)

		query.WithWhere(filter)
		countQuery.WithWhere(filter)
	}

	result, err := mergeErrors(query.Do(context.Background()))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	parsed, err := mapResponseObjects(result.Data)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	classData, ok := parsed.Get[s.cn.String()]
	if !ok {
		return nil, fault.New("weaviate response did not contain expected class data")
	}

	return classData, nil
}

// TODO: GroupBy on the datagraph_id
func (s *weaviateSemdexer) countObjects(ctx context.Context, countQuery graphql.AggregateBuilder) (int, error) {
	r, err := mergeErrors(countQuery.Do(ctx))
	if err != nil {
		return 0, fault.Wrap(err, fctx.With(ctx))
	}

	type AggregateResponse struct {
		Aggregate map[ /* class name */ string][]struct {
			Field struct {
				Count int `mapstructure:"count"`
			} `mapstructure:"datagraph_id"`
		}
	}

	var agg AggregateResponse
	err = mapstructure.Decode(r.Data, &agg)
	if err != nil {
		return 0, fault.Wrap(err, fctx.With(ctx))
	}

	classes := agg.Aggregate[s.cn.String()]
	if len(classes) < 1 {
		return 0, fault.New("no class data in aggregate count query")
	}

	count := classes[0].Field.Count

	return count, nil
}

func filterChunks(results []*datagraph.Ref) []*datagraph.Ref {
	filtered := dt.Filter(results, func(r *datagraph.Ref) bool {
		return r.Relevance > RelevanceThreshold
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
		score := 0.0

		for _, r := range curr {
			score += r.Relevance
		}

		next := &datagraph.Ref{
			ID:        first.ID,
			Kind:      first.Kind,
			Relevance: score / float64(len(curr)),
		}

		return append(acc, next)
	}, []*datagraph.Ref{})

	sort.Sort(datagraph.RefList(deduped))

	return deduped
}

func mapObjectToChunk(o WeaviateObject) (*semdex.Chunk, error) {
	id, err := xid.FromString(o.DatagraphID)
	if err != nil {
		return nil, err
	}

	kind, err := datagraph.NewKind(o.DatagraphType)
	if err != nil {
		return nil, err
	}

	sdr, err := url.Parse(fmt.Sprintf("%s:%s/%s", datagraph.RefScheme, kind, id.String()))
	if err != nil {
		return nil, err
	}

	return &semdex.Chunk{
		ID:      id,
		Kind:    kind,
		URL:     *sdr,
		Content: o.Content,
	}, nil
}

func mapObjectsToChunks(objects []WeaviateObject) ([]*semdex.Chunk, error) {
	return dt.MapErr(objects, mapObjectToChunk)
}
