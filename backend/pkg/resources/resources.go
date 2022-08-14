package resources

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/backend/pkg/resources/account"
	"github.com/Southclaws/storyden/backend/pkg/resources/authentication"
)

func Build() fx.Option {
	return fx.Provide(
		account.New,
		authentication.New,
		// category.New,
		// post.New,
		// tag.New,
		// thread.New,
		// react.New,
		// notification.New,
	)
}
