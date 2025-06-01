package generative

import (
	"context"
	"html/template"
	"strings"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/internal/infrastructure/ai"
)

var SuggestTagsPrompt = template.Must(template.New("").Parse(`Analyze the provided content and generate up to three relevant tags. Tags are either single words or multiple words separated only by a hyphen, no spaces.

It's very important that only tags that are relevant to the content are returned, any tags of low confidence MUST be omitted. Do not generate tags that are too vague or tags that are too specific and cannot easily be used in other contexts for other types of content. Generally avoid tags that are singular and not plural that too closely match phrases or words in the content.

Suggested Tags: Suggest any existing tags from the list above that best describe this content, prioritizing tags that closely match the main themes, ideas, or entities mentioned.
New Tags: If there are no suitable matches or if additional tags could enhance discoverability, suggest up to three new tags not found in the existing list. Ensure the tags are relevant, concise, and enhance content discoverability.

Consider the following list of existing tags on the platform for context and prioritization: {{ .AvailableTags }}

Content:

{{ .Content }}
`))

type SuggestTagsResultSchema struct {
	Tags []string `json:"tags" jsonschema:"title=Tags,description=List of suggested tags,items=string"`
}

func (g *generator) SuggestTags(ctx context.Context, content datagraph.Content, available tag_ref.Names) (tag_ref.Names, error) {
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

	result, err := ai.PromptObject(ctx, g.prompter, "Suggest tags for content", template.String(), SuggestTagsResultSchema{})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	tags := dt.Map(result.Tags, func(s string) tag_ref.Name {
		return tag_ref.NewName(s)
	})

	return tags, nil
}
