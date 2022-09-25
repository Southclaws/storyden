package infrastructure

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/infrastructure/db"
	"github.com/Southclaws/storyden/internal/infrastructure/http"
	"github.com/Southclaws/storyden/internal/infrastructure/logger"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
	"github.com/Southclaws/storyden/internal/infrastructure/securecookie"
	"github.com/Southclaws/storyden/internal/infrastructure/webauthn"
)

func Build() fx.Option {
	return fx.Options(
		logger.Build(),
		// mailer.Build(),
		db.Build(),
		http.Build(),
		fx.Provide(securecookie.New),
		fx.Provide(webauthn.New),
		fx.Provide(object.NewS3Storer),
	)
}
