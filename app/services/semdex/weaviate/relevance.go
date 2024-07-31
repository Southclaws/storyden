package weaviate

import (
	"context"
	"encoding/json"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/rs/xid"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

func (w *weaviateSemdexer) ScoreRelevance(ctx context.Context, object datagraph.Indexable, ids ...xid.ID) (map[xid.ID]float64, error) {
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

	out := dt.Reduce(results, func(acc map[xid.ID]float64, ref *datagraph.NodeReference) map[xid.ID]float64 {
		acc[ref.ID] = ref.Score
		return acc
	}, map[xid.ID]float64{})

	return out, nil
}

func mapResponseObjects(raw map[string]models.JSONObject) (*WeaviateResponse, error) {
	j, err := json.Marshal(raw)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	parsed := WeaviateResponse{}
	err = json.Unmarshal(j, &parsed)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &parsed, nil
}
