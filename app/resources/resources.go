package resources

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/collection"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_children"
	"github.com/Southclaws/storyden/app/resources/library/node_search"
	"github.com/Southclaws/storyden/app/resources/library/node_traversal"
	"github.com/Southclaws/storyden/app/resources/link/link_querier"
	"github.com/Southclaws/storyden/app/resources/link/link_writer"
	"github.com/Southclaws/storyden/app/resources/mailtemplate"
	"github.com/Southclaws/storyden/app/resources/notification"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/resources/post/post_search"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/profile/profile_search"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/resources/tag"
)

func Build() fx.Option {
	return fx.Options(
		rbac.Build(),
		fx.Provide(
			settings.New,
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
			library.New,
			node_traversal.New,
			node_children.New,
			node_search.New,
			link_querier.New,
			link_writer.New,
			profile_search.New,
			mailtemplate.New,
		),
	)
}
