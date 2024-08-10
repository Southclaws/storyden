package weaviate_semdexer

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

func (s *weaviateRefIndex) Summarise(ctx context.Context, object datagraph.Item) (string, error) {
	fields := []graphql.Field{
		{Name: "datagraph_id"},
		{Name: "datagraph_type"},
		{Name: "name"},
		{Name: "content"},
	}

	// Switch summariser strategy based on the Weaviate class.
	// Local inference uses sum-transformers, remote inference uses openai.
	// TODO: Express this switcher in a better way at the top-level config.
	if s.cn.String() == "ContentText2vecTransformers" {
		fields = append(fields, graphql.Field{
			Name: "_additional",
			Fields: []graphql.Field{
				{Name: "summary(properties: [\"content\"])", Fields: []graphql.Field{
					{Name: "property"},
					{Name: "result"},
				}},
			},
		})
	} else if s.cn.String() == "ContentOpenAI" {
		fields = append(fields, graphql.Field{
			Name: "_additional",
			Fields: []graphql.Field{
				{Name: `generate(singleResult: {
					prompt: """
						Describe the following as a short summary: {content}
					"""
				})`, Fields: []graphql.Field{
					{Name: "singleResult"},
					{Name: "error"},
				}},
			},
		})
	} else {
		// No summariser available
		// TODO: return an error maybe?
		return "", nil
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
		return "", fault.Newf("expected exactly one result, got %d", len(classData))
	}

	if s.cn.String() == "ContentText2vecTransformers" {
		if classData[0].Additional.Summary == nil || len(classData[0].Additional.Summary) != 1 {
			return "", fault.New("summary not found in response")
		}

		return classData[0].Additional.Summary[0].Result, nil

	} else if s.cn.String() == "ContentOpenAI" {
		if classData[0].Additional.Generate.Error != "" {
			return "", fault.New(classData[0].Additional.Generate.Error)
		}

		return classData[0].Additional.Generate.SingleResult, nil
	}

	// TODO: handle this edge case
	return "", nil
}
