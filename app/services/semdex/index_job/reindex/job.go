package reindex

import (
	"context"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

func runReindexer(
	ctx context.Context,
	lc fx.Lifecycle,
	l *zap.Logger,

	qnode pubsub.Topic[mq.IndexNode],
	qpost pubsub.Topic[mq.IndexPost],
	re *reindexer,
) {
	if re == nil {
		return
	}

	lc.Append(fx.StartHook(func(_ context.Context) error {
		err := re.reindexAll(ctx)
		if err != nil {
			return err
		}

		go func() {
			for range time.NewTicker(time.Hour).C {
				err := re.reindexAll(ctx)
				if err != nil {
					l.Error("failed to run reindex job", zap.Error(err))
				}
			}
		}()

		return nil
	}))
}
