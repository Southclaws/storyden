package branding

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/branding/banner"
	"github.com/Southclaws/storyden/app/services/branding/icon"
	"github.com/Southclaws/storyden/app/services/branding/theme"
)

func Build() fx.Option {
	return fx.Provide(icon.New, banner.New, theme.New)
}
