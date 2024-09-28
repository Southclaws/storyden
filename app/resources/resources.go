package resources

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/account/notification/notify_querier"
	"github.com/Southclaws/storyden/app/resources/account/notification/notify_writer"
	"github.com/Southclaws/storyden/app/resources/account/role/role_assign"
	"github.com/Southclaws/storyden/app/resources/account/role/role_querier"
	"github.com/Southclaws/storyden/app/resources/account/role/role_writer"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/collection"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_children"
	"github.com/Southclaws/storyden/app/resources/library/node_search"
	"github.com/Southclaws/storyden/app/resources/library/node_traversal"
	"github.com/Southclaws/storyden/app/resources/like/like_querier"
	"github.com/Southclaws/storyden/app/resources/like/like_writer"
	"github.com/Southclaws/storyden/app/resources/link/link_querier"
	"github.com/Southclaws/storyden/app/resources/link/link_writer"
	"github.com/Southclaws/storyden/app/resources/mailtemplate"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/resources/post/post_search"
	"github.com/Southclaws/storyden/app/resources/post/post_writer"
	"github.com/Southclaws/storyden/app/resources/post/reaction"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/profile/follow_querier"
	"github.com/Southclaws/storyden/app/resources/profile/follow_writer"
	"github.com/Southclaws/storyden/app/resources/profile/profile_search"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/resources/tag"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			settings.New,
			account_querier.New,
			account_writer.New,
			role_assign.New,
			role_querier.New,
			role_writer.New,
			asset.New,
			authentication.New,
			category.New,
			notify_querier.New,
			notify_writer.New,
			reply.New,
			tag.New,
			thread.New,
			reaction.New,
			like_querier.New,
			like_writer.New,
			post_search.New,
			post_writer.New,
			collection.New,
			library.New,
			node_traversal.New,
			node_children.New,
			node_search.New,
			link_querier.New,
			link_writer.New,
			profile_search.New,
			follow_writer.New,
			follow_querier.New,
			mailtemplate.New,
		),
	)
}
