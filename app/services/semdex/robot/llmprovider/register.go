package llmprovider

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/robot/llm_provider"
	"github.com/Southclaws/storyden/app/services/semdex/robot/llmprovider/anthropic"
	"github.com/Southclaws/storyden/app/services/semdex/robot/llmprovider/mockllm"
	"github.com/Southclaws/storyden/app/services/semdex/robot/llmprovider/openai"
	"github.com/Southclaws/storyden/app/services/semdex/robot/model_media"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(model_media.New, openai.New, anthropic.New),
		fx.Invoke(func(lc fx.Lifecycle, factory *llm_provider.Factory, openAI *openai.OpenAI, claude *anthropic.Anthropic) {
			factory.Put(openAI)
			factory.Put(claude)
			factory.Put(mockllm.Mock{})
		}),
	)
}
