package asker

import (
	"context"
	"html/template"
	"log/slog"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/question"
	"github.com/Southclaws/storyden/app/services/search/searcher"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/ai"
)

func New(
	cfg config.Config,
	logger *slog.Logger,
	searcher semdex.Searcher,
	prompter ai.Prompter,
	questions *question.Repository,
) (semdex.Asker, error) {
	asker, err := newAsker(cfg, searcher, prompter)
	if err != nil {
		return nil, err
	}

	return newCachedAsker(
		logger,
		asker,
		questions,
	)
}

func newAsker(cfg config.Config, searcher semdex.Searcher, prompter ai.Prompter) (semdex.Asker, error) {
	if cfg.SemdexProvider != "" && cfg.LanguageModelProvider == "" {
		return nil, fault.New("semdex requires a language model provider to be enabled")
	}

	switch cfg.AskerProvider {
	case "perplexity":
		// NOTE: While Perplexity looks like it could satisfy the language model
		// provider interface, it does not provide an embedding func, it's only
		// functional for chat-like interactions so it's only an Asker for now.
		// This means that if you wish to use Perplexity, you must also provide
		// a language model provider such as OpenAI along with an API key. Keep
		// this in mind when considering the cost of your Storyden installation.
		return newPerplexityAsker(cfg, searcher)

	default:
		return &defaultAsker{
			searcher: searcher,
			prompter: prompter,
		}, nil
	}
}

var AnswerPrompt = template.Must(template.New("").Parse(`
You are an expert assistant. Answer the user's question using the provided "Additional sources" as a primary reference. You MUST incorporate these sources into your answer and cite them in the "Sources" section. 
All sources must be cited in the "Sources" section without modification. Include URLs and descriptions as provided in the "Additional sources."
You MUST include AT LEAST ONE reference to the sources below in your answer IN ADDITION to other sources you may have been provided by a system prompt.
Combine the "Additional sources" with other sources from your system, but ensure at least one of the "Additional sources" is referenced in the response.

Additional sources that you MUST use in your answer:

{{- range .Context }}
- URL: {{ .URL.String }}
  Key points: {{ .Content }}
{{- end }}

Question: {{ .Question }}

Answer:
1. Provide your answer here in clear and concise paragraphs.
2. Use information from the sources above to support your answer, but do not include citations inline.
3. Include a "Sources" section with the sources listed.

Sources:
- <url> (<kind>) <short description of why this source was used>
`))

const maxContextForRAG = 10

func buildContextPrompt(ctx context.Context, s semdex.Searcher, q string) (string, error) {
	chunks, err := s.SearchChunks(ctx, q, pagination.NewPageParams(1, 200), searcher.Options{})
	if err != nil {
		return "", fault.Wrap(err, fctx.With(ctx))
	}

	if len(chunks) == 0 {
		return "", fault.New("no context found for question", fctx.With(ctx), ftag.With(ftag.NotFound))
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
		return "", fault.Wrap(err, fctx.With(ctx))
	}

	return t.String(), nil
}
