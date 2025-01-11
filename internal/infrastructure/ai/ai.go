package ai

import (
	"context"

	"github.com/Southclaws/storyden/internal/config"
)

type Result struct {
	Answer string
}

type Embedder func(ctx context.Context, text string) ([]float32, error)

type Prompter interface {
	Prompt(ctx context.Context, input string) (*Result, error)
	PromptStream(ctx context.Context, input string) (func(yield func(string, error) bool), error)
	EmbeddingFunc() func(ctx context.Context, text string) ([]float32, error)
}

func New(cfg config.Config) (Prompter, error) {
	switch cfg.LanguageModelProvider {
	case "openai":
		return newOpenAI(cfg)

	case "mock":
		return newMock()

	default:
		return &Disabled{}, nil
	}
}

type Disabled struct{}

func (d *Disabled) Prompt(ctx context.Context, input string) (*Result, error) {
	return nil, nil
}

func (d *Disabled) PromptStream(ctx context.Context, input string) (func(yield func(string, error) bool), error) {
	return nil, nil
}

func (d *Disabled) EmbeddingFunc() func(ctx context.Context, text string) ([]float32, error) {
	return nil
}
