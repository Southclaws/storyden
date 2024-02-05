package weaviate

import (
	"context"
	"encoding/json"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/services/semdex"
)

type WeaviateObject struct {
	DatagraphID   string `json:"datagraph_id"`
	DatagraphType string `json:"datagraph_type"`
	Name          string `json:"name"`
	Content       string `json:"content"`
}

type WeaviateContent map[string][]WeaviateObject

type WeaviateResponse struct {
	Get WeaviateContent
}

func (s *weaviateSemdexer) Search(ctx context.Context, q string) ([]*semdex.Result, error) {
	fields := []graphql.Field{
		{Name: "datagraph_id"},
		{Name: "datagraph_type"},
		{Name: "name"},
		{Name: "content"},
	}

	arg := s.wc.GraphQL().
		HybridArgumentBuilder().
		WithAlpha(0.25).
		WithFusionType(graphql.RelativeScore).
		WithQuery(q)

	result, err := s.wc.GraphQL().Get().
		WithClassName(TestClassName).
		WithFields(fields...).
		WithHybrid(arg).
		WithLimit(30).
		Do(context.Background())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	j, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	parsed := WeaviateResponse{}
	err = json.Unmarshal(j, &parsed)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	classData, ok := parsed.Get[TestClassName]
	if !ok {
		return nil, fault.New("weaviate response did not contain expected class data")
	}

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
