package weaviate

import (
	"context"
	"encoding/json"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
	"go.uber.org/multierr"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

func (w *weaviateSemdexer) Recommend(ctx context.Context, object datagraph.Indexable) (datagraph.NodeReferenceList, error) {
	wid := GetWeaviateID(object.GetID())

	result, err := w.wc.Data().ObjectsGetter().
		WithClassName(w.cn.String()).
		WithVector().
		WithID(wid).
		Do(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	wobj := result[0]

	// TODO: Compute vector between account owner and object.

	withNearVector := w.wc.GraphQL().NearVectorArgBuilder().
		WithVector(wobj.Vector).
		WithCertainty(0.5)

	fields := []graphql.Field{
		{Name: "datagraph_id"},
		{Name: "datagraph_type"},
	}

	recommendations, err := w.wc.GraphQL().Get().
		WithClassName(w.cn.String()).
		WithFields(fields...).
		WithNearVector(withNearVector).
		WithAutocut(2).
		WithLimit(10).
		Do(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if len(recommendations.Errors) > 0 {
		return nil, fault.Wrap(gqlerror(recommendations.Errors), fctx.With(ctx))
	}

	j, err := json.Marshal(recommendations.Data)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	parsed := WeaviateResponse{}
	err = json.Unmarshal(j, &parsed)
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

	return results, nil
}

func gqlerror(gqe []*models.GraphQLError) error {
	return fault.Wrap(multierr.Combine(dt.Map(gqe, func(e *models.GraphQLError) error {
		return fault.New(e.Message)
	})...))
}
