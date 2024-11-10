package weaviate_semdexer

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

func (s *weaviateRefIndex) Search(ctx context.Context, q string) (datagraph.RefList, error) {
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

	arg := s.wc.GraphQL().
		HybridArgumentBuilder().
		WithAlpha(0.25).
		WithFusionType(graphql.RelativeScore).
		WithQuery(q)

	result, err := mergeErrors(s.wc.GraphQL().Get().
		WithClassName(s.cn.String()).
		WithFields(fields...).
		WithHybrid(arg).
		WithAutocut(2).
		WithLimit(30).
		Do(context.Background()))
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

	results, err := dt.MapErr(classData, mapToNodeReference)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return results, nil
}
