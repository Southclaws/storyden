package services

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/account"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/app/services/avatar"
	"github.com/Southclaws/storyden/app/services/post"
	"github.com/Southclaws/storyden/app/services/rbac"
	"github.com/Southclaws/storyden/app/services/thread"
)

func Build() fx.Option {
	return fx.Options(
		account.Build(),
		authentication.Build(),
		rbac.Build(),
		thread.Build(),
		post.Build(),
		avatar.Build(),
	)
}
