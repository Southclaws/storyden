package services

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/account/account_suspension"
	"github.com/Southclaws/storyden/app/services/account/register"
	"github.com/Southclaws/storyden/app/services/asset"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/app/services/avatar"
	"github.com/Southclaws/storyden/app/services/avatar_gen"
	"github.com/Southclaws/storyden/app/services/category"
	"github.com/Southclaws/storyden/app/services/collection"
	"github.com/Southclaws/storyden/app/services/icon"
	"github.com/Southclaws/storyden/app/services/library/node_mutate"
	"github.com/Southclaws/storyden/app/services/library/node_read"
	"github.com/Southclaws/storyden/app/services/library/node_visibility"
	"github.com/Southclaws/storyden/app/services/library/nodetree"
	"github.com/Southclaws/storyden/app/services/like/post_liker"
	"github.com/Southclaws/storyden/app/services/link"
	"github.com/Southclaws/storyden/app/services/mention/mention_job"
	"github.com/Southclaws/storyden/app/services/notification/notify_job"
	"github.com/Southclaws/storyden/app/services/onboarding"
	"github.com/Southclaws/storyden/app/services/profile/following"
	"github.com/Southclaws/storyden/app/services/react_manager"
	"github.com/Southclaws/storyden/app/services/reply"
	"github.com/Southclaws/storyden/app/services/search"
	"github.com/Southclaws/storyden/app/services/semdex/index_job"
	"github.com/Southclaws/storyden/app/services/semdex/semdexer"
	"github.com/Southclaws/storyden/app/services/semdex/summarise_job"
	"github.com/Southclaws/storyden/app/services/thread"
	"github.com/Southclaws/storyden/app/services/thread_mark"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(register.New),
		icon.Build(),
		onboarding.Build(),
		account_suspension.Build(),
		authentication.Build(),
		category.Build(),
		thread.Build(),
		reply.Build(),
		post_liker.Build(),
		react_manager.Build(),
		search.Build(),
		avatar.Build(),
		asset.Build(),
		thread_mark.Build(),
		collection.Build(),
		link.Build(),
		notify_job.Build(),
		mention_job.Build(),
		semdexer.Build(),
		index_job.Build(),
		summarise_job.Build(),
		fx.Provide(avatar_gen.New),
		fx.Provide(node_read.New, node_mutate.New, nodetree.New, node_visibility.New),
		fx.Provide(following.New),
	)
}
