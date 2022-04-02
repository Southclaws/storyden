package infra

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/api/src/infra/db"
	"github.com/Southclaws/storyden/api/src/infra/logger"
	"github.com/Southclaws/storyden/api/src/infra/mailer"
	"github.com/Southclaws/storyden/api/src/infra/pubsub"
)

func Build() fx.Option {
	return fx.Options(
		logger.Build(),
		mailer.Build(),
		pubsub.Build(),
		db.Build(),
	)
}
