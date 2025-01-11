package asker

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/infrastructure/ai"
)

// defaultAsker uses whatever prompter is available and performs RAG prompting.
type defaultAsker struct {
	searcher semdex.Searcher
	prompter ai.Prompter
}

func (a *defaultAsker) Ask(ctx context.Context, q string) (chan string, chan error) {
	t, err := buildContextPrompt(ctx, a.searcher, q)
	if err != nil {
		ech := make(chan error, 1)
		ech <- fault.Wrap(err, fctx.With(ctx))
		return nil, ech
	}

	chch, ech := a.prompter.PromptStream(ctx, t)

	return chch, ech
}
