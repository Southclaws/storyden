package asker

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/infrastructure/ai"
)

// defaultAsker uses whatever prompter is available and performs RAG prompting.
type defaultAsker struct {
	searcher semdex.Searcher
	prompter ai.Prompter
}

func (a *defaultAsker) Ask(ctx context.Context, q string, parent opt.Optional[xid.ID]) (semdex.AskResponseIterator, error) {
	t, err := buildContextPrompt(ctx, a.searcher, q)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	iter, err := a.prompter.PromptStream(ctx, t)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return streamExtractor(iter), nil
}
