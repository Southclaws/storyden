package analyse_job

import (
	"go.uber.org/fx"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(newAnalyseConsumer),
	)
}
