package resources

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/authentication"
	"github.com/Southclaws/storyden/app/resources/category"
	"github.com/Southclaws/storyden/app/resources/cluster"
	"github.com/Southclaws/storyden/app/resources/cluster_children"
	"github.com/Southclaws/storyden/app/resources/cluster_traversal"
	"github.com/Southclaws/storyden/app/resources/collection"
	"github.com/Southclaws/storyden/app/resources/item"
	"github.com/Southclaws/storyden/app/resources/item_search"
	"github.com/Southclaws/storyden/app/resources/link"
	"github.com/Southclaws/storyden/app/resources/link_graph"
	"github.com/Southclaws/storyden/app/resources/notification"
	"github.com/Southclaws/storyden/app/resources/post_search"
	"github.com/Southclaws/storyden/app/resources/profile_search"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/app/resources/reply"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/resources/tag"
	"github.com/Southclaws/storyden/app/resources/thread"
)

func Build() fx.Option {
	return fx.Options(
		rbac.Build(),
		fx.Provide(
			settings.New,
			account.New,
			asset.New,
			authentication.New,
			category.New,
			reply.New,
			tag.New,
			thread.New,
			react.New,
			notification.New,
			post_search.New,
			collection.New,
			cluster.New,
			cluster_traversal.New,
			cluster_children.New,
			item.New,
			item_search.New,
			link.New,
			link_graph.New,
			profile_search.New,
		),
	)
}
