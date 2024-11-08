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
	"github.com/samber/lo"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

var SuggestTagsPrompt = template.Must(template.New("").Parse(`Analyze the provided content of \"{name}\" and generate relevant tags. Tags are either single words or multiple words separated only by a hyphen, no spaces.

It's very important that only tags that are relevant to the content are returned, any tags of low confidence MUST be omitted. Do not generate tags that are too vague or tags that are too specific and cannot easily be used in other contexts for other types of content. Generally avoid tags that are singular and not plural that too closely match phrases or words in the content.

Suggested Tags: Suggest any existing tags from the list above that best describe this content, prioritizing tags that closely match the main themes, ideas, or entities mentioned.
New Tags: If there are no suitable matches or if additional tags could enhance discoverability, suggest up to three new tags not found in the existing list. Ensure the tags are relevant, concise, and enhance content discoverability.

Output Format: Provide only a list of tags separated by commas with no additional text or symbols. Suggested tags should come first, followed by new tags, if any.

Consider the following list of existing tags on the platform for context and prioritization: {{ .AvailableTags }}

Content:

{{ .Content }}
`))

func (s *weaviateRefIndex) SuggestTags(ctx context.Context, content datagraph.Content, available tag_ref.Names) (tag_ref.Names, error) {
	// cap the available tags at 50, we don't to blow out the prompt size limit.
	sliced := lo.Splice(available, 50)

	template := strings.Builder{}
	err := SuggestTagsPrompt.Execute(&template, map[string]any{
		"AvailableTags": sliced,
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
