package notify_job

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/notification/notify"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(newNotifyConsumer),
		fx.Invoke(runNotifyConsumer),

		fx.Provide(notify.New),
	)
}
