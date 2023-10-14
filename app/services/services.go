package services

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/account"
	"github.com/Southclaws/storyden/app/services/asset"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/app/services/avatar"
	"github.com/Southclaws/storyden/app/services/avatar_gen"
	"github.com/Southclaws/storyden/app/services/category"
	"github.com/Southclaws/storyden/app/services/cluster"
	"github.com/Southclaws/storyden/app/services/clustertree"
	"github.com/Southclaws/storyden/app/services/collection"
	"github.com/Southclaws/storyden/app/services/icon"
	"github.com/Southclaws/storyden/app/services/item_crud"
	"github.com/Southclaws/storyden/app/services/onboarding"
	"github.com/Southclaws/storyden/app/services/react"
	"github.com/Southclaws/storyden/app/services/reply"
	"github.com/Southclaws/storyden/app/services/search"
	"github.com/Southclaws/storyden/app/services/thread"
	"github.com/Southclaws/storyden/app/services/thread_mark"
	"github.com/Southclaws/storyden/app/services/thread_url"
	"github.com/Southclaws/storyden/app/services/url"
)

func Build() fx.Option {
	return fx.Options(
		icon.Build(),
		onboarding.Build(),
		account.Build(),
		authentication.Build(),
		category.Build(),
		thread.Build(),
		reply.Build(),
		react.Build(),
		search.Build(),
		avatar.Build(),
		asset.Build(),
		thread_mark.Build(),
		collection.Build(),
		url.Build(),
		thread_url.Build(),
		fx.Provide(avatar_gen.New),
		fx.Provide(cluster.New, clustertree.New),
		fx.Provide(item_crud.New),
	)
}
