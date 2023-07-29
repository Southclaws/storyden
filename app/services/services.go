package services

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/account"
	"github.com/Southclaws/storyden/app/services/asset"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/app/services/avatar"
	"github.com/Southclaws/storyden/app/services/avatar_gen"
	"github.com/Southclaws/storyden/app/services/react"
	"github.com/Southclaws/storyden/app/services/reply"
	"github.com/Southclaws/storyden/app/services/search"
	"github.com/Southclaws/storyden/app/services/thread"
	"github.com/Southclaws/storyden/app/services/thread_mark"
)

func Build() fx.Option {
	return fx.Options(
		account.Build(),
		authentication.Build(),
		thread.Build(),
		reply.Build(),
		react.Build(),
		search.Build(),
		avatar.Build(),
		asset.Build(),
		thread_mark.Build(),
		fx.Provide(avatar_gen.New),
	)
}
