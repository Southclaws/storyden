package index_job

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/services/semdex/index_job/reindex"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub/queue"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			queue.New[mq.IndexNode],
			queue.New[mq.IndexThread],
			queue.New[mq.IndexReply],
			queue.New[mq.IndexProfile],
		),

		fx.Provide(newIndexConsumer),
		fx.Invoke(runIndexConsumer),

		reindex.Build(),
	)
}
