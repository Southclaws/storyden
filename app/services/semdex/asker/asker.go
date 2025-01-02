package asker

import (
	"context"
	"html/template"
	"strings"
	"time"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/question"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/search/searcher"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/ai"
)

type Asker struct {
	logger    *zap.Logger
	searcher  semdex.Searcher
	prompter  ai.Prompter
	questions *question.Repository
}

func New(cfg config.Config, logger *zap.Logger, searcher semdex.Searcher, prompter ai.Prompter, questions *question.Repository) (*Asker, error) {
	if cfg.SemdexProvider != "" && cfg.LanguageModelProvider == "" {
		return nil, fault.New("semdex requires a language model provider to be enabled")
	}

	return &Asker{
		logger:    logger,
		searcher:  searcher,
		prompter:  prompter,
		questions: questions,
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

func (a *Asker) Ask(ctx context.Context, q string) (func(yield func(string, error) bool), error) {
	cached, err := a.questions.GetByQuerySlug(ctx, q)
	if err == nil {
		return a.cachedResult(ctx, cached)
	}

	return a.livePrompt(ctx, q)
}

func (a *Asker) cachedResult(ctx context.Context, q *question.Question) (func(yield func(string, error) bool), error) {
	md, err := htmltomarkdown.ConvertNode(q.Result.HTMLTree())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	chunks := strings.SplitAfter(string(md), " ")

	return func(yield func(string, error) bool) {
		for _, ch := range chunks {
			select {
			case <-ctx.Done():
				return

			default:
				if !yield(ch, nil) {
					return
				}
			}
			time.Sleep(time.Millisecond * 10)
		}
	}, nil
}

func (a *Asker) livePrompt(ctx context.Context, q string) (func(yield func(string, error) bool), error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	t, err := a.buildPrompt(ctx, q)
	if err != nil {
		return nil, err
	}

	chch, ech := a.prompter.PromptStream(ctx, t.String())

	return func(yield func(string, error) bool) {
		acc := []string{}

		cleanup := func() {
			err := a.cacheResult(ctx, accountID, q, acc)
			if err != nil {
				a.logger.Error("failed to cache result", zap.Error(err))
			}
		}

		for {
			select {
			case <-ctx.Done():
				return

			case chunk, ok := <-chch:
				if !ok {
					cleanup()
					return
				}

				acc = append(acc, chunk)
				if !yield(chunk, nil) {
					cleanup()
					return
				}

			case err := <-ech:
				if err != nil {
					yield("", err)
					return
				}
			}
		}
	}, nil
}

func (a *Asker) buildPrompt(ctx context.Context, q string) (*strings.Builder, error) {
	chunks, err := a.searcher.SearchChunks(ctx, q, pagination.NewPageParams(1, 200), searcher.Options{})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if len(chunks) == 0 {
		return nil, fault.New("no context found for question", fctx.With(ctx), ftag.With(ftag.NotFound))
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
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &t, nil
}

func (a *Asker) cacheResult(ctx context.Context, accountID account.AccountID, q string, chunks []string) error {
	result := strings.Join(chunks, "")

	acc, err := datagraph.NewRichTextFromMarkdown(result)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	_, err = a.questions.Store(ctx, accountID, q, acc)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
