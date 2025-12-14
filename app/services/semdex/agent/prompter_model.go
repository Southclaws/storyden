package agent

import (
	"context"
	"iter"

	"github.com/Southclaws/storyden/internal/infrastructure/ai"
	"google.golang.org/adk/model"
)

type prompterModel struct {
	prompter ai.Prompter
}

func newPrompterModel(p ai.Prompter) model.LLM {
	if p == nil {
		return nil
	}
	return &prompterModel{prompter: p}
}

func (p *prompterModel) Name() string {
	return "storyden_prompter"
}

func (p *prompterModel) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	return p.prompter.GenerateContent(ctx, req, stream)
}
