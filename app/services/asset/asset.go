package asset

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/asset/asset_download"
	"github.com/Southclaws/storyden/app/services/asset/asset_upload"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			asset_upload.New,
			asset_download.New,
		),
	)
}
