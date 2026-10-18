package background_job

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/background_job/job_dispatcher"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			New,
			func(runner *Runner) job_dispatcher.Dispatcher {
				return runner
			},
		),
		fx.Invoke(func(*Runner) {}),
	)
}
