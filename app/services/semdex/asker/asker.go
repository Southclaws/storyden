package asker

import (
	"context"
	"html/template"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/services/search/searcher"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/ai"
)

type Asker struct {
	searcher semdex.Searcher
	prompter ai.Prompter
}

func New(cfg config.Config, searcher semdex.Searcher, prompter ai.Prompter) (*Asker, error) {
	if cfg.SemdexProvider != "" && cfg.LanguageModelProvider == "" {
		return nil, fault.New("semdex requires a language model provider to be enabled")
	}

	return &Asker{
		searcher: searcher,
		prompter: prompter,
	}, nil
}

var AnswerPrompt = template.Must(template.New("").Parse(`
You are an expert assistant. Answer the user's question accurately and concisely using the provided sources. Cite the sources in a separate list at the end of your answer. 
Ensure that the source URLs (in "sdr" format) are kept exactly as they appear, without modification or breaking them across lines.

Sources:
{{- range .Context }}
- URL: {{ .URL.String }}
  Kind: {{ .Kind }}
  Content: {{ .Content }}
{{- end }}

Question: {{ .Question }}

Answer:
1. Provide your answer here in clear and concise paragraphs.
2. Use information from the sources above to support your answer, but do not include citations inline.
3. Include a "References" section with the source URLs listed, like this:

References:
- (the url to the source): (Short description of the source content)
`))

const maxContextForRAG = 10

func (a *Asker) Ask(ctx context.Context, q string) (chan string, chan error) {
	chunks, err := a.searcher.SearchChunks(ctx, q, pagination.NewPageParams(1, 200), searcher.Options{})
	if err != nil {
		ech := make(chan error, 1)
		ech <- fault.Wrap(err, fctx.With(ctx))
		return nil, ech
	}

	if len(chunks) == 0 {
		ech := make(chan error, 1)
		ech <- fault.New("no context found for question", fctx.With(ctx), ftag.With(ftag.NotFound))
		return nil, ech
	}

	if len(chunks) > maxContextForRAG {
		chunks = chunks[:maxContextForRAG]
	}

	t := strings.Builder{}
	err = AnswerPrompt.Execute(&t, map[string]any{
		"Context":  chunks,
		"Question": q,
	})
	if err != nil {
		ech := make(chan error, 1)
		ech <- fault.Wrap(err, fctx.With(ctx))
		return nil, ech
	}

	chch, ech := a.prompter.PromptStream(ctx, t.String())

	return chch, ech
}
