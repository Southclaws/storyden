package weaviate_semdexer

import (
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

var SummarisePrompt = template.Must(template.New("").Parse(`
Write a short few paragraphs that are somewhat engaging but remaining relatively neutral in tone in the style of a wikipedia introduction about \"{name}\". Focus on providing unique insights and interesting details while keeping the tone conversational and approachable. Imagine this will be read by someone browsing a directory or knowledgebase.

Be aware that the input to this may include broken HTML and other artifacts from the web and due to the nature of web scraping, there may be parts that do not make sense.

- Ignore any HTML tags, malformed content, or text that does not contribute meaningfully to the main topic.
- Based on the clear and coherent sections of the input, write short but engaging paragraphs. If the input lacks meaningful context, produce a neutral placeholder.
- If the input content is too fragmented or lacks sufficient context to produce a coherent response, produce a neutral placeholder.
- Do not describe the appearance of the input (e.g., broken HTML or artifacts). Instead, infer the main idea or purpose and expand on it creatively.
- If key parts of the content are missing or ambiguous, use creativity to fill gaps while maintaining relevance to the topic.

Output Format: Provide the output as a correctly formatted HTML document, you are free to use basic HTML formatting tags for emphasis, lists and headings. However, do not include the content title as a <h1> tag at the top. Start with a paragraph block immediately.

Content:

{content}
`))

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

		template := strings.Builder{}
		err := SummarisePrompt.Execute(&template, map[string]any{})
		if err != nil {
			return "", fault.Wrap(err, fctx.With(ctx))
		}

		prompt := strings.ReplaceAll(template.String(), "\n", `\n`)

		summaryPrompt := fmt.Sprintf(`generate(singleResult: {
		prompt: """
			%s
			"""
		})`, prompt)

		fields = append(fields, graphql.Field{
			Name: "_additional",
			Fields: []graphql.Field{
				{Name: summaryPrompt, Fields: []graphql.Field{
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
