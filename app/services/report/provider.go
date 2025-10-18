package report

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/report/member_report"
	"github.com/Southclaws/storyden/app/services/report/report_manager"
	"github.com/Southclaws/storyden/app/services/report/report_notify"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			member_report.New,
			report_manager.New,
		),
		report_notify.Build(),
	)
}
