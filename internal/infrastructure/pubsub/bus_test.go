package pubsub_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"go.uber.org/fx"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestEventBus_SingleSubscriber(t *testing.T) {
	integration.Test(t, &config.Config{
		// QueueType: "amqp",
		// AmqpURL: "amqp://guest:guest@localhost:5672/",
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

			_, err := pubsub.Subscribe(ctx, bus, "test_service", func(ctx context.Context, event *EventTest) error {
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
		}))
	}))
}

func TestEventBus_MultipleSubscribers(t *testing.T) {
	integration.Test(t, &config.Config{
		// QueueType: "amqp",
		// AmqpURL: "amqp://guest:guest@localhost:5672/",
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

			r.NoError(sub1.Close())

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

func TestPublishNamed_ReceivesEventsViaGoChannel(t *testing.T) {
	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		bus *pubsub.Bus,
	) {
		lc.Append(fx.StartHook(func(ctx context.Context) {
			r := require.New(t)
			a := assert.New(t)

			topicName := "tests.named.channels.single"
			received := make(chan message.EventThreadPublished, 1)

			sub, err := pubsub.SubscribeNamed(ctx, bus, topicName, "named_channel_handler", func(ctx context.Context, payload json.RawMessage) error {
				var event message.EventThreadPublished
				if err := json.Unmarshal(payload, &event); err != nil {
					return err
				}
				received <- event
				return nil
			})
			r.NoError(err)
			t.Cleanup(func() {
				r.NoError(sub.Close())
			})

			wantID := post.ID(xid.New())
			err = bus.PublishNamed(ctx, topicName, message.EventThreadPublished{ID: wantID})
			r.NoError(err)

			select {
			case event := <-received:
				a.Equal(wantID, event.ID)
			case <-time.After(2 * time.Second):
				t.Fatal("timed out waiting for gochannel event")
			}
		}))
	}))
}

func TestSubscribeNamed_ReceivesEventsViaAMQP(t *testing.T) {
	integration.Test(t, &config.Config{
		// QueueType: "amqp",
		// AmqpURL: "amqp://guest:guest@localhost:5672/",
	}, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		bus *pubsub.Bus,
	) {
		lc.Append(fx.StartHook(func(ctx context.Context) {
			r := require.New(t)
			a := assert.New(t)

			topicName := "tests.named.amqp.single"
			received := make(chan message.EventThreadPublished, 1)

			sub, err := pubsub.SubscribeNamed(ctx, bus, topicName, "dynamic_single_handler", func(ctx context.Context, payload json.RawMessage) error {
				var event message.EventThreadPublished
				if err := json.Unmarshal(payload, &event); err != nil {
					return err
				}
				received <- event
				return nil
			})
			r.NoError(err)
			t.Cleanup(func() {
				r.NoError(sub.Close())
			})

			wantID := post.ID(xid.New())

			err = bus.PublishNamed(ctx, topicName, message.EventThreadPublished{ID: wantID})
			r.NoError(err)

			select {
			case event := <-received:
				a.Equal(wantID, event.ID)
			case <-time.After(5 * time.Second):
				t.Fatal("timed out waiting for dynamic event")
			}
		}))
	}))
}

func TestSubscribeNamed_FanoutToIndependentHandlers(t *testing.T) {
	integration.Test(t, &config.Config{
		// QueueType: "amqp",
		// AmqpURL: "amqp://guest:guest@localhost:5672/",
	}, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		bus *pubsub.Bus,
	) {
		lc.Append(fx.StartHook(func(ctx context.Context) {
			r := require.New(t)
			a := assert.New(t)

			type dynamicFanoutEvent struct {
				Value string `json:"value"`
			}

			topicName := "tests.named.amqp.fanout"
			recvOne := make(chan dynamicFanoutEvent, 1)
			recvTwo := make(chan dynamicFanoutEvent, 1)

			handler := func(target chan<- dynamicFanoutEvent) func(context.Context, json.RawMessage) error {
				return func(ctx context.Context, payload json.RawMessage) error {
					var event dynamicFanoutEvent
					if err := json.Unmarshal(payload, &event); err != nil {
						return err
					}
					target <- event
					return nil
				}
			}

			subOne, err := pubsub.SubscribeNamed(ctx, bus, topicName, "dynamic_fanout_one", handler(recvOne))
			r.NoError(err)
			t.Cleanup(func() {
				r.NoError(subOne.Close())
			})

			subTwo, err := pubsub.SubscribeNamed(ctx, bus, topicName, "dynamic_fanout_two", handler(recvTwo))
			r.NoError(err)
			t.Cleanup(func() {
				r.NoError(subTwo.Close())
			})

			err = bus.PublishNamed(ctx, topicName, dynamicFanoutEvent{Value: "fanout payload"})
			r.NoError(err)

			deadline := time.After(5 * time.Second)
			var one, two dynamicFanoutEvent
			for receivedOne, receivedTwo := false, false; !(receivedOne && receivedTwo); {
				select {
				case <-deadline:
					t.Fatal("timed out waiting for fanout events")
				case evt := <-recvOne:
					one = evt
					receivedOne = true
				case evt := <-recvTwo:
					two = evt
					receivedTwo = true
				}
			}

			a.Equal("fanout payload", one.Value)
			a.Equal("fanout payload", two.Value)
		}))
	}))
}

func TestSubscribeNamed_DuplicateHandlerNameErrors(t *testing.T) {
	integration.Test(t, &config.Config{
		// QueueType: "amqp",
		// AmqpURL: "amqp://guest:guest@localhost:5672/",
	}, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		bus *pubsub.Bus,
	) {
		lc.Append(fx.StartHook(func(ctx context.Context) {
			r := require.New(t)

			const handlerName = "dynamic_duplicate"
			topicName := "tests.named.amqp.duplicate_check"

			sub, err := pubsub.SubscribeNamed(ctx, bus, topicName, handlerName, func(context.Context, json.RawMessage) error {
				return nil
			})
			r.NoError(err)
			t.Cleanup(func() {
				r.NoError(sub.Close())
			})

			_, err = pubsub.SubscribeNamed(ctx, bus, topicName, handlerName, func(context.Context, json.RawMessage) error {
				return nil
			})
			r.Error(err)
		}))
	}))
}

func TestPublishNamed_DeliversToTypedSubscriber(t *testing.T) {
	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		bus *pubsub.Bus,
	) {
		lc.Append(fx.StartHook(func(ctx context.Context) {
			r := require.New(t)
			a := assert.New(t)

			recv := make(chan message.EventThreadPublished, 1)

			sub, err := pubsub.Subscribe(ctx, bus, "named_to_typed_handler", func(ctx context.Context, event *message.EventThreadPublished) error {
				recv <- *event
				return nil
			})
			r.NoError(err)
			t.Cleanup(func() {
				r.NoError(sub.Close())
			})

			topicName := reflect.TypeOf(message.EventThreadPublished{}).String()
			wantID := post.ID(xid.New())

			err = bus.PublishNamed(ctx, topicName, message.EventThreadPublished{ID: wantID})
			r.NoError(err)

			select {
			case evt := <-recv:
				a.Equal(wantID, evt.ID)
			case <-time.After(2 * time.Second):
				t.Fatal("timed out waiting for typed named event")
			}
		}))
	}))
}

func TestPublishersAndSubscribers_MixedNamedAndTyped(t *testing.T) {
	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		bus *pubsub.Bus,
	) {
		lc.Append(fx.StartHook(func(ctx context.Context) {
			r := require.New(t)
			a := assert.New(t)

			topicName := reflect.TypeOf(message.EventThreadPublished{}).String()
			typedRecv := make(chan post.ID, 2)
			dynamicRecv := make(chan post.ID, 2)

			typedSub, err := pubsub.Subscribe(ctx, bus, "mixed_typed_handler", func(ctx context.Context, event *message.EventThreadPublished) error {
				typedRecv <- event.ID
				return nil
			})
			r.NoError(err)
			t.Cleanup(func() {
				r.NoError(typedSub.Close())
			})

			dynamicSub, err := pubsub.SubscribeNamed(ctx, bus, topicName, "mixed_dynamic_handler", func(ctx context.Context, payload json.RawMessage) error {
				var event message.EventThreadPublished
				if err := json.Unmarshal(payload, &event); err != nil {
					return err
				}
				dynamicRecv <- event.ID
				return nil
			})
			r.NoError(err)
			t.Cleanup(func() {
				r.NoError(dynamicSub.Close())
			})

			waitID := func(ch <-chan post.ID, timeout time.Duration) post.ID {
				select {
				case id := <-ch:
					return id
				case <-time.After(timeout):
					t.Fatal("timed out waiting for mixed subscriber event")
				}
				return post.ID{}
			}

			wantAuto := post.ID(xid.New())
			r.NoError(bus.MustPublish(ctx, message.EventThreadPublished{ID: wantAuto}))
			a.Equal(wantAuto, waitID(typedRecv, 2*time.Second))
			a.Equal(wantAuto, waitID(dynamicRecv, 2*time.Second))

			wantNamed := post.ID(xid.New())
			r.NoError(bus.PublishNamed(ctx, topicName, message.EventThreadPublished{ID: wantNamed}))
			a.Equal(wantNamed, waitID(typedRecv, 2*time.Second))
			a.Equal(wantNamed, waitID(dynamicRecv, 2*time.Second))
		}))
	}))
}
