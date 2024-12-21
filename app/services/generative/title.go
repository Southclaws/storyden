package generative

import (
	"context"
	"html/template"
	"strings"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

var TitlePrompt = template.Must(template.New("").Parse(`
Generate 1 to 3 concise and genuine title suggestions for the following content. Titles should feel natural, not overly stylized or pushy, and suitable for diverse content types. Limit each title to 60 characters. Respond with only the titles, separated by newlines.

Content:

{{ .Content }}
`))

func (g *generator) SuggestTitle(ctx context.Context, content datagraph.Content) ([]string, error) {
	template := strings.Builder{}
	err := TitlePrompt.Execute(&template, map[string]any{
		"Content": content.Plaintext(),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result, err := g.prompter.Prompt(ctx, template.String())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	titles := strings.Split(result.Answer, "\n")

	trimmed := dt.Map(titles, strings.TrimSpace)

	filtered := dt.Filter(trimmed, func(s string) bool { return len(s) > 0 })

	return filtered, nil
}
