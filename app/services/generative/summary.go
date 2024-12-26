package generative

import (
	"context"
	"html/template"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

var SummarisePrompt = template.Must(template.New("").Parse(`
Write a short few paragraphs that are somewhat engaging but remaining relatively neutral in tone in the style of a wikipedia introduction based on the specified content. Focus on providing unique insights and interesting details while keeping the tone conversational and approachable. Imagine this will be read by someone browsing a directory or knowledgebase.

Be aware that the input to this may include broken HTML and other artifacts from the web and due to the nature of web scraping, there may be parts that do not make sense.

- Ignore any HTML tags, malformed content, or text that does not contribute meaningfully to the main topic.
- Based on the clear and coherent sections of the input, write short but engaging paragraphs. If the input lacks meaningful context, produce a neutral placeholder.
- If the input content is too fragmented or lacks sufficient context to produce a coherent response, produce a neutral placeholder.
- Do not describe the appearance of the input (e.g., broken HTML or artifacts). Instead, infer the main idea or purpose and expand on it creatively.
- If key parts of the content are missing or ambiguous, use creativity to fill gaps while maintaining relevance to the topic.

Output Format: Provide the output as a correctly formatted HTML document, but do not include any markdown tags around the output, you are free to use basic HTML formatting tags for emphasis, lists and headings. However, do not include the content title as a <h1> tag at the top. Start with a paragraph block immediately.

Content:

{{ .Content }}
`))

func (g *generator) Summarise(ctx context.Context, content datagraph.Content) (string, error) {
	template := strings.Builder{}
	err := SummarisePrompt.Execute(&template, map[string]any{
		"Content": content.Plaintext(),
	})
	if err != nil {
		return "", fault.Wrap(err, fctx.With(ctx))
	}

	result, err := g.prompter.Prompt(ctx, template.String())
	if err != nil {
		return "", fault.Wrap(err, fctx.With(ctx))
	}

	return result.Answer, nil
}
