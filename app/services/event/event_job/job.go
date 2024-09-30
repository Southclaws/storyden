package event_job

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

func runEventUpdateConsumer(
	ctx context.Context,
	lc fx.Lifecycle,
	l *zap.Logger,

	queue pubsub.Topic[mq.CreateEvent],

	ic *eventUpdateConsumer,
) {
	lc.Append(fx.StartHook(func(_ context.Context) error {
		channel, err := queue.Subscribe(ctx)
		if err != nil {
			panic(err)
		}

		go func() {
			for msg := range channel {
				// TODO:
				// - notify participants of certain changes to events
				//   - event has been cancelled
				//   - event has been rescheduled
				//   - event capacity changed
				//   - event location changed
				// - notify participants status and role changes
				//   - accepted after requesting
				//   - promoted to host
				//   - demoted to attendee
				//   - removed from event
				msg.Ack()
			}
		}()

		return nil
	}))
}
