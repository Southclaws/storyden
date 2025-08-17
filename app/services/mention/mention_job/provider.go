package mention_job

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/mention/mentioner"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(newMentionConsumer),
		fx.Invoke(runMentionConsumer),

		fx.Provide(mentioner.New),
	)
}
