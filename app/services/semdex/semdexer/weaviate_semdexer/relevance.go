package weaviate_semdexer

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

func (w *weaviateRefIndex) ScoreRelevance(ctx context.Context, object datagraph.Item, ids ...xid.ID) (map[xid.ID]float64, error) {
	if len(ids) == 0 {
		return map[xid.ID]float64{}, nil
	}

	sids := dt.Map(ids, func(id xid.ID) string { return id.String() })

	near := graphql.NearObjectArgumentBuilder{}
	near.WithID(GetWeaviateID(object.GetID())).WithDistance(10)

	where := filters.Where().
		WithPath([]string{"datagraph_id"}).
		WithOperator(filters.ContainsAny).
		WithValueString(sids...)

	fields := []graphql.Field{
		{Name: "datagraph_id"},
		{Name: "datagraph_type"},
		{Name: "_additional", Fields: []graphql.Field{{Name: "distance"}}},
	}

	r, err := mergeErrors(w.wc.GraphQL().Get().
		WithNearObject(&near).
		WithClassName(w.cn.String()).
		WithFields(fields...).
		WithWhere(where).
		Do(ctx))
	if err != nil {
		// janky error handling because weaviate client doesn't use sentinals
		if err.Error() == "vector not found" {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	parsed, err := mapResponseObjects(r.Data)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	classData := parsed.Get[w.cn.String()]

	classData = dt.Filter(classData, func(o WeaviateObject) bool {
		return o.DatagraphID != object.GetID().String()
	})

	results, err := dt.MapErr(classData, mapToNodeReference)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	out := dt.Reduce(results, func(acc map[xid.ID]float64, ref *datagraph.Ref) map[xid.ID]float64 {
		acc[ref.ID] = ref.Relevance
		return acc
	}, map[xid.ID]float64{})

	return out, nil
}
