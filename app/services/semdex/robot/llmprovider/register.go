package llmprovider

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/robot/llm_provider"
	"github.com/Southclaws/storyden/app/services/semdex/robot/llmprovider/anthropic"
	"github.com/Southclaws/storyden/app/services/semdex/robot/llmprovider/mockllm"
	"github.com/Southclaws/storyden/app/services/semdex/robot/llmprovider/openai"
)

func Build() fx.Option {
	return fx.Invoke(func(lc fx.Lifecycle, factory *llm_provider.Factory) {
		factory.Put(&openai.OpenAI{})
		factory.Put(&anthropic.Anthropic{})
		factory.Put(mockllm.Mock{})
	})
}
