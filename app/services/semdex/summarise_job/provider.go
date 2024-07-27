package summarise_job

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/internal/pubsub/queue"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			queue.New[mq.SummariseNode],
			// queue.New[mq.SummarisePost], // TODO
			// queue.New[mq.SummariseProfile], // TODO
		),

		fx.Provide(newSummariseConsumer),
		fx.Invoke(runSummariseConsumer),
	)
}
