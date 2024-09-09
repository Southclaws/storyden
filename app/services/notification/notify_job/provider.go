package notify_job

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/services/notification/notify"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub/queue"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			queue.New[mq.Notification],
		),

		fx.Provide(newNotifyConsumer),
		fx.Invoke(runNotifyConsumer),

		fx.Provide(notify.New),
	)
}
