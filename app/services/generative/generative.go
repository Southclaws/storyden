package generative

import (
	"context"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/robot/llm_provider"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
)

type Titler interface {
	SuggestTitle(ctx context.Context, content datagraph.Content) ([]string, error)
}

type Tagger interface {
	SuggestTags(ctx context.Context, content datagraph.Content, available tag_ref.Names) (tag_ref.Names, error)
}

type Summariser interface {
	Summarise(ctx context.Context, content datagraph.Content) (string, error)
}

var (
	_ Titler     = &generator{}
	_ Tagger     = &generator{}
	_ Summariser = &generator{}
)

type generator struct {
	models *llm_provider.Factory
}

func newGenerator(models *llm_provider.Factory) *generator {
	return &generator{models: models}
}

func Build() fx.Option {
	return fx.Provide(
		fx.Annotate(
			newGenerator,
			fx.As(new(Titler)),
			fx.As(new(Tagger)),
			fx.As(new(Summariser)),
		),
	)
}
