package resources

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/pkg/resources/account"
	"github.com/Southclaws/storyden/pkg/resources/authentication"
	"github.com/Southclaws/storyden/pkg/resources/category"
	"github.com/Southclaws/storyden/pkg/resources/notification"
	"github.com/Southclaws/storyden/pkg/resources/post"
	"github.com/Southclaws/storyden/pkg/resources/react"
	"github.com/Southclaws/storyden/pkg/resources/tag"
	"github.com/Southclaws/storyden/pkg/resources/thread"
)

func Build() fx.Option {
	return fx.Provide(
		account.New,
		authentication.New,
		category.New,
		post.New,
		tag.New,
		thread.New,
		react.New,
		notification.New,
	)
}
