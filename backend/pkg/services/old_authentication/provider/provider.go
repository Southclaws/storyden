package provider

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/backend/pkg/services/authentication/provider/basic"
	"github.com/Southclaws/storyden/backend/pkg/services/authentication/provider/magiclink"
	"github.com/Southclaws/storyden/backend/pkg/services/authentication/provider/oauth"
)

func Build() fx.Option {
	return fx.Options(
		// Basic auth is just simple username/password.
		fx.Provide(basic.NewBasicAuth),

		// Magic links are passwordless and use a provided communication method
		// to send the user a link that logs them in.
		magiclink.Build(),

		// OAuth is for integration with other services like login with twitter.
		oauth.Build(),
	)
}
