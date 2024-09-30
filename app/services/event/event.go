package event

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/event/event_job"
	"github.com/Southclaws/storyden/app/services/event/event_management"
	"github.com/Southclaws/storyden/app/services/event/event_participation"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(event_management.New),
		fx.Provide(event_participation.New),
		event_job.Build(),
	)
}
