package analyse_job

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub/queue"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			queue.New[mq.AnalyseAsset],
			queue.New[mq.DownloadAsset],
		),

		fx.Provide(newAnalyseConsumer),
		fx.Invoke(runAnalyseConsumer),
	)
}
