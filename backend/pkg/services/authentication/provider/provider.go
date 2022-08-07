package provider

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/backend/pkg/services/authentication/provider/magiclink"
	"github.com/Southclaws/storyden/backend/pkg/services/authentication/provider/oauth"
	"github.com/Southclaws/storyden/backend/pkg/services/authentication/provider/password"
)

func Build() fx.Option {
	return fx.Options(
		// Password auth is just simple username/password.
		fx.Provide(password.NewBasicAuth),

		// Magic links are passwordless and use a provided communication method
		// to send the user a link that logs them in.
		magiclink.Build(),

		// OAuth is for integration with other services like login with twitter.
		oauth.Build(),
	)
}
