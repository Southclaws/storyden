package notification

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/notification/global"
	"github.com/Southclaws/storyden/app/services/notification/notify"
	"github.com/Southclaws/storyden/app/services/notification/push"
	"github.com/Southclaws/storyden/app/services/notification/settings"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(notify.New),
		fx.Provide(settings.New),
		global.Build(),
		push.Build(),
	)
}
