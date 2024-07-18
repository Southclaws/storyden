package weaviate

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

func (s *weaviateSemdexer) Summarise(ctx context.Context, object datagraph.Indexable) (string, error) {
	fields := []graphql.Field{
		{Name: "datagraph_id"},
		{Name: "datagraph_type"},
		{Name: "name"},
		{Name: "content"},
		{Name: "_additional", Fields: []graphql.Field{
			{Name: "summary(properties: [\"content\"])", Fields: []graphql.Field{
				{Name: "property"},
				{Name: "result"},
			}},
		}},
	}

	where := filters.Where().
		WithPath([]string{"datagraph_id"}).
		WithOperator(filters.ContainsAny).
		WithValueString(object.GetID().String())

	result, err := mergeErrors(s.wc.GraphQL().Get().
		WithClassName(s.cn.String()).
		WithFields(fields...).
		WithWhere(where).
		Do(context.Background()))
	if err != nil {
		return "", fault.Wrap(err, fctx.With(ctx))
	}

	parsed, err := mapResponseObjects(result.Data)
	if err != nil {
		return "", err
	}

	classData := parsed.Get[s.cn.String()]

	if len(classData) != 1 {
		return "", fault.New("expected exactly one result")
	}

	if classData[0].Additional.Summary == nil || len(classData[0].Additional.Summary) != 1 {
		return "", fault.New("summary not found in response")
	}

	summary := classData[0].Additional.Summary[0].Result

	return summary, nil
}
