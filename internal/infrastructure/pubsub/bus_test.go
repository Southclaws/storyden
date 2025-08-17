package pubsub_test

import (
	"context"
	"testing"

	"go.uber.org/fx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestEventBus_SingleSubscriber(t *testing.T) {
	integration.Test(t, &config.Config{
		// QueueType: "amqp",
		AmqpURL: "amqp://guest:guest@localhost:5672/",
	}, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		bus *pubsub.Bus,
	) {
		lc.Append(fx.StartHook(func(ctx context.Context) {
			r := require.New(t)
			a := assert.New(t)

			type EventTest struct {
				Value string
			}

			recv := make(chan EventTest)

			sub, err := pubsub.Subscribe(ctx, bus, "test_service", func(ctx context.Context, event *EventTest) error {
				recv <- *event
				return nil
			})
			r.NoError(err)

			err = bus.MustPublish(ctx, EventTest{
				Value: "Hello, World!",
			})
			r.NoError(err)

			received := <-recv
			a.Equal("Hello, World!", received.Value)

			sub.Close()
		}))
	}))
}

func TestEventBus_MultipleSubscribers(t *testing.T) {
	integration.Test(t, &config.Config{
		// QueueType: "amqp",
		AmqpURL: "amqp://guest:guest@localhost:5672/",
	}, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		bus *pubsub.Bus,
	) {
		lc.Append(fx.StartHook(func(ctx context.Context) {
			r := require.New(t)
			a := assert.New(t)

			type MultiEventTest struct {
				Value string
			}

			recv := make(chan MultiEventTest)

			sub1, err := pubsub.Subscribe(ctx, bus, "test_service_one", func(ctx context.Context, event *MultiEventTest) error {
				recv <- *event
				return nil
			})
			r.NoError(err)
			_, err = pubsub.Subscribe(ctx, bus, "test_service_two", func(ctx context.Context, event *MultiEventTest) error {
				recv <- *event
				return nil
			})
			r.NoError(err)

			err = bus.MustPublish(ctx, MultiEventTest{
				Value: "Hello, World!",
			})
			r.NoError(err)

			received1 := <-recv
			a.Equal("Hello, World!", received1.Value)

			received2 := <-recv
			a.Equal("Hello, World!", received2.Value)

			a.Equal(received1, received2)

			sub1.Close()

			err = bus.MustPublish(ctx, MultiEventTest{
				Value: "Message for only sub2",
			})
			r.NoError(err)

			received3 := <-recv
			a.Equal("Message for only sub2", received3.Value)

			// NOTE: This causes a flaky test because the router closes after
			// no more subscribers are left - in reality, this wouldn't happen
			// because there will always be some consumers hard-coded to run.
			// sub2.Close()

			// err = bus.MustPublish(ctx, MultiEventTest{
			// 	Value: "No more subscribers. No-op.",
			// })
			// r.NoError(err)
		}))
	}))
}
