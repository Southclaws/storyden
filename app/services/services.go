package services

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/account"
	"github.com/Southclaws/storyden/app/services/account_suspension"
	"github.com/Southclaws/storyden/app/services/asset_manager"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/app/services/avatar"
	"github.com/Southclaws/storyden/app/services/avatar_gen"
	"github.com/Southclaws/storyden/app/services/category"
	"github.com/Southclaws/storyden/app/services/cluster"
	"github.com/Southclaws/storyden/app/services/cluster/cluster_visibility"
	"github.com/Southclaws/storyden/app/services/clustertree"
	"github.com/Southclaws/storyden/app/services/collection"
	"github.com/Southclaws/storyden/app/services/hydrator"
	"github.com/Southclaws/storyden/app/services/hydrator/fetcher"
	"github.com/Southclaws/storyden/app/services/icon"
	"github.com/Southclaws/storyden/app/services/item_crud"
	"github.com/Southclaws/storyden/app/services/item_tree"
	"github.com/Southclaws/storyden/app/services/onboarding"
	"github.com/Southclaws/storyden/app/services/react"
	"github.com/Southclaws/storyden/app/services/reply"
	"github.com/Southclaws/storyden/app/services/search"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/app/services/semdex/datagraph_searcher"
	"github.com/Southclaws/storyden/app/services/semdex/indexer"
	"github.com/Southclaws/storyden/app/services/thread"
	"github.com/Southclaws/storyden/app/services/thread_mark"
	"github.com/Southclaws/storyden/app/services/url"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

func Build() fx.Option {
	return fx.Options(
		icon.Build(),
		onboarding.Build(),
		account.Build(),
		account_suspension.Build(),
		authentication.Build(),
		category.Build(),
		thread.Build(),
		reply.Build(),
		react.Build(),
		search.Build(),
		avatar.Build(),
		asset_manager.Build(),
		thread_mark.Build(),
		collection.Build(),
		url.Build(),
		hydrator.Build(),
		fetcher.Build(),
		buildSemdex(),
		fx.Provide(avatar_gen.New),
		fx.Provide(cluster.New, clustertree.New, cluster_visibility.New),
		fx.Provide(item_crud.New, item_tree.New),
		fx.Provide(datagraph_searcher.New),
	)
}

func buildSemdex() fx.Option {
	return fx.Provide(func(lc fx.Lifecycle, wc *weaviate.Client) (semdex.Service, error) {
		if wc == nil {
			if wc == nil {
				return semdex.Empty{}, nil
			}
		}

		return indexer.New(lc, wc)
	})
}
