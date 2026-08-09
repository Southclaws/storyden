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
Write a brief, concise summary of the following content in 2-3 short sentences, never exceeding one paragraph. This summary serves as a starting point for curators who will write their own descriptions later - it should just provide enough context to remind them what this content is about.

Keep the tone neutral and factual, similar to a Wikipedia introduction. Focus on what the subject IS, not why it's interesting or what makes it special.

The input may include a URL and Original Title which should help you understand the context, followed by the actual content. The content may also include broken HTML and other web scraping artifacts.

Guidelines:
- Maximum 2-3 sentences, always in a single paragraph
- Use URL and Original Title as context clues but don't reference them directly in the summary
- Focus on factual description of what the subject is
- Avoid promotional or marketing language
- Ignore HTML tags, malformed content, or irrelevant fragments
- If content is too fragmented, produce a simple neutral description based on URL/title
- Do NOT be creative or fill in gaps - stick to what's clearly stated
- Do NOT use phrases like "this article discusses" or "this page is about"

Output Format: Plain HTML paragraph tag(s) only. Start immediately with <p> tag. No headings, no multiple paragraphs, no markdown.

<example>
URL: https://selfh.st/
Original Title: Selfh.st - Modern Self-Hosting Made Easy

Content:
Selfh.st is a modern self-hosting platform that makes it easy to deploy applications...

Good output: <p>Selfh.st is a self-hosting platform for deploying applications. It provides tools for managing and running containerized services.</p>

Bad output: <p>Selfh.st is an innovative and exciting new platform that revolutionizes the way we think about self-hosting.</p><p>With its modern approach and user-friendly interface, it makes deployment easier than ever before.</p><p>Whether you're a beginner or an expert, you'll find something to love.</p>
</example>

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
