package infrastructure

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/backend/internal/infrastructure/db"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/logger"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/securecookie"
)

func Build() fx.Option {
	return fx.Options(
		logger.Build(),
		// mailer.Build(),
		db.Build(),
		fx.Provide(securecookie.New),
	)
}
