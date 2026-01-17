// Package infrastructure simply provides all the plumbing packages to the DI.
package infrastructure

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/infrastructure/ai"
	"github.com/Southclaws/storyden/internal/infrastructure/cache"
	"github.com/Southclaws/storyden/internal/infrastructure/db"
	"github.com/Southclaws/storyden/internal/infrastructure/endec/jwt"
	"github.com/Southclaws/storyden/internal/infrastructure/frontend"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation"
	"github.com/Southclaws/storyden/internal/infrastructure/logger"
	"github.com/Southclaws/storyden/internal/infrastructure/mailer"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/infrastructure/rate"
	"github.com/Southclaws/storyden/internal/infrastructure/redis"
	"github.com/Southclaws/storyden/internal/infrastructure/sms"
	"github.com/Southclaws/storyden/internal/infrastructure/vector/pinecone"
	"github.com/Southclaws/storyden/internal/infrastructure/weaviate"
	"github.com/Southclaws/storyden/internal/infrastructure/webauthn"
)

func Build() fx.Option {
	return fx.Options(
		logger.Build(),
		instrumentation.Build(),
		db.Build(),
		redis.Build(),
		cache.Build(),
		fx.Provide(rate.NewFactory),
		mailer.Build(),
		sms.Build(),
		fx.Provide(webauthn.New),
		object.Build(),
		frontend.Build(),
		weaviate.Build(),
		pinecone.Build(),
		fx.Provide(ai.New),
		jwt.Build(),
		pubsub.Build(),
	)
}
