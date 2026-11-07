package notification

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/notification/notify"
	"github.com/Southclaws/storyden/app/services/notification/push"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(notify.New),
		push.Build(),
	)
}
