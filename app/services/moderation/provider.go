package moderation

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/moderation/content_policy"
	"github.com/Southclaws/storyden/app/services/moderation/spam"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(spam.New),
		fx.Provide(content_policy.New),
	)
}
