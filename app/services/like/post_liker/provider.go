package post_liker

import (
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub/queue"
	"go.uber.org/fx"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			queue.New[mq.LikePost],
			queue.New[mq.UnlikePost],
		),
		fx.Provide(New),
	)
}
