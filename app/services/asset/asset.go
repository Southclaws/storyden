package asset

import (
	"github.com/Southclaws/storyden/app/services/asset/analyse"
	"github.com/Southclaws/storyden/app/services/asset/analyse_job"
	"github.com/Southclaws/storyden/app/services/asset/asset_download"
	"github.com/Southclaws/storyden/app/services/asset/asset_upload"
	"go.uber.org/fx"
)

func Build() fx.Option {
	return fx.Options(
		analyse_job.Build(),
		fx.Provide(
			analyse.New,
			asset_upload.New,
			asset_download.New,
		),
	)
}
