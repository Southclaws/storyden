package services

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/account/account_auth"
	"github.com/Southclaws/storyden/app/services/account/account_email"
	"github.com/Southclaws/storyden/app/services/account/account_suspension"
	"github.com/Southclaws/storyden/app/services/account/register"
	"github.com/Southclaws/storyden/app/services/asset"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/app/services/avatar"
	"github.com/Southclaws/storyden/app/services/avatar_gen"
	"github.com/Southclaws/storyden/app/services/category"
	"github.com/Southclaws/storyden/app/services/collection"
	"github.com/Southclaws/storyden/app/services/comms"
	"github.com/Southclaws/storyden/app/services/event"
	"github.com/Southclaws/storyden/app/services/icon"
	"github.com/Southclaws/storyden/app/services/library"
	"github.com/Southclaws/storyden/app/services/like/post_liker"
	"github.com/Southclaws/storyden/app/services/link"
	"github.com/Southclaws/storyden/app/services/mention/mention_job"
	"github.com/Southclaws/storyden/app/services/moderation"
	"github.com/Southclaws/storyden/app/services/notification/notify_job"
	"github.com/Southclaws/storyden/app/services/onboarding"
	"github.com/Southclaws/storyden/app/services/profile/following"
	"github.com/Southclaws/storyden/app/services/react_manager"
	"github.com/Southclaws/storyden/app/services/reply"
	"github.com/Southclaws/storyden/app/services/search"
	"github.com/Southclaws/storyden/app/services/semdex/index_job"
	"github.com/Southclaws/storyden/app/services/semdex/semdexer"
	"github.com/Southclaws/storyden/app/services/system/instance_info"
	"github.com/Southclaws/storyden/app/services/tag/autotagger"
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
		library.Build(),
		comms.Build(),
		link.Build(),
		notify_job.Build(),
		mention_job.Build(),
		semdexer.Build(),
		index_job.Build(),
		event.Build(),
		moderation.Build(),
		fx.Provide(avatar_gen.New),
		fx.Provide(following.New),
		fx.Provide(autotagger.New),
		fx.Provide(instance_info.New),
		fx.Provide(account_auth.New, account_email.New),
	)
}
