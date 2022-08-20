package resources

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/backend/pkg/resources/account"
	"github.com/Southclaws/storyden/backend/pkg/resources/authentication"
	"github.com/Southclaws/storyden/backend/pkg/resources/category"
	"github.com/Southclaws/storyden/backend/pkg/resources/notification"
	"github.com/Southclaws/storyden/backend/pkg/resources/post"
	"github.com/Southclaws/storyden/backend/pkg/resources/react"
	"github.com/Southclaws/storyden/backend/pkg/resources/tag"
	"github.com/Southclaws/storyden/backend/pkg/resources/thread"
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
