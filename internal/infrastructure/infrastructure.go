// Package infrastructure simply provides all the plumbing packages to the DI.
package infrastructure

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/db"
	"github.com/Southclaws/storyden/internal/frontend"
	"github.com/Southclaws/storyden/internal/logger"
	"github.com/Southclaws/storyden/internal/object"
	"github.com/Southclaws/storyden/internal/securecookie"
	"github.com/Southclaws/storyden/internal/sms"
	"github.com/Southclaws/storyden/internal/weaviate"
	"github.com/Southclaws/storyden/internal/webauthn"
)

func Build() fx.Option {
	return fx.Options(
		logger.Build(),
		db.Build(),
		securecookie.Build(),
		sms.Build(),
		fx.Provide(webauthn.New),
		object.Build(),
		frontend.Build(),
		weaviate.Build(),
	)
}
