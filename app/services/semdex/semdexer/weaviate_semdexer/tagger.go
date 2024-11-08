package weaviate_semdexer

import (
	"context"
	"strings"
	"text/template"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

var SuggestTagsPrompt = template.Must(template.New("").Parse(`Analyze the provided content of \"{name}\" and generate relevant tags. Tags are either single words or multiple words separated only by a hyphen, no spaces. Consider the following list of existing tags on the platform for context and prioritization: {{ .AvailableTags }}

Suggested Tags: Suggest any existing tags from the list above that best describe this content, prioritizing tags that closely match the main themes, ideas, or entities mentioned.
New Tags: If there are no suitable matches or if additional tags could enhance discoverability, suggest up to three new tags not found in the existing list. Ensure the tags are relevant, concise, and enhance content discoverability.

Output Format: Provide only a list of tags separated by commas with no additional text or symbols. Suggested tags should come first, followed by new tags, if any.

Content:

{{ .Content }}
`))

func (s *weaviateRefIndex) SuggestTags(ctx context.Context, content datagraph.Content, available tag_ref.Names) (tag_ref.Names, error) {
	template := strings.Builder{}
	err := SuggestTagsPrompt.Execute(&template, map[string]any{
		"AvailableTags": available,
		"Content":       content.Plaintext(),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	prompt := strings.ReplaceAll(template.String(), "\n", `\n`)

	gs := graphql.NewGenerativeSearch().SingleResult(prompt)

	r, err := mergeErrors(s.wc.GraphQL().
		Get().
		WithClassName(s.cn.String()).
		WithLimit(1).
		WithGenerativeSearch(gs).
		Do(ctx))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	wr, err := mapResponseObjects(r.Data)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	object, err := s.getFirstResult(wr)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	strings := strings.Split(object.Additional.Generate.SingleResult, ", ")

	tags := dt.Map(strings, func(s string) tag_ref.Name {
		return tag_ref.NewName(s)
	})

	return tags, nil
}
