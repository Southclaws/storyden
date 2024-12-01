package branding

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/branding/banner"
	"github.com/Southclaws/storyden/app/services/branding/icon"
)

func Build() fx.Option {
	return fx.Provide(icon.New, banner.New)
}
