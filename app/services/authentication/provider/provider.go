package provider

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth"
	"github.com/Southclaws/storyden/app/services/authentication/provider/password"
)

func Build() fx.Option {
	return fx.Options(
		// Password auth is just simple username/password.
		fx.Provide(password.NewBasicAuth),

		// Magic links are passwordless and use a provided communication method
		// to send the user a link that logs them in.
		// magiclink.Build(), // TODO

		// OAuth is for integration with other services like login with twitter.
		oauth.Build(),
	)
}
