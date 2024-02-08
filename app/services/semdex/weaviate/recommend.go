package weaviate

import (
	"context"
	"encoding/json"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
	"go.uber.org/multierr"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/services/semdex"
)

func (w *weaviateSemdexer) Recommend(ctx context.Context, object datagraph.Indexable) ([]*semdex.Result, error) {
	wid := GetWeaviateID(object.GetID())

	result, err := w.wc.Data().ObjectsGetter().
		WithClassName(TestClassName).
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
		WithDistance(0.5)

	fields := []graphql.Field{
		{Name: "datagraph_id"},
		{Name: "datagraph_type"},
	}

	recommendations, err := w.wc.GraphQL().Get().
		WithClassName(TestClassName).
		WithFields(fields...).
		WithNearVector(withNearVector).
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

	classData := parsed.Get[TestClassName]

	results, err := dt.MapErr(classData, func(v WeaviateObject) (*semdex.Result, error) {
		id, err := xid.FromString(v.DatagraphID)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		dk, err := datagraph.NewKind(v.DatagraphType)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		return &semdex.Result{
			Id:   id,
			Type: dk,
			Name: v.Name,
		}, nil
	})
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
