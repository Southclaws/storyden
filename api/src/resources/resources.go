package resources

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/api/src/resources/user"
)

func Build() fx.Option {
	return fx.Provide(
		// category.New,
		// post.New,
		// tag.New,
		// thread.New,
		user.New,
		// react.New,
		// notification.New,
	)
}
