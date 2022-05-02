package resources

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/api/src/resources/authentication"
	"github.com/Southclaws/storyden/api/src/resources/user"
)

func Build() fx.Option {
	return fx.Provide(
		authentication.New,
		user.New,
		// category.New,
		// post.New,
		// tag.New,
		// thread.New,
		// react.New,
		// notification.New,
	)
}
