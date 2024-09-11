package mention_job

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/services/mention/mentioner"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub/queue"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			queue.New[mq.Mention],
		),

		fx.Provide(newMentionConsumer),
		fx.Invoke(runMentionConsumer),

		fx.Provide(mentioner.New),
	)
}
