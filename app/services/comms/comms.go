package comms

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/comms/mailqueue"
	"github.com/Southclaws/storyden/app/services/comms/mailtemplate"
)

func Build() fx.Option {
	return fx.Options(
		mailqueue.Build(),
		fx.Provide(mailtemplate.New),
	)
}
