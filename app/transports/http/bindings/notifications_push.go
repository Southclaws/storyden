package bindings

import (
	"context"
	"encoding/json"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/account/notification/push_subscription"
	"github.com/Southclaws/storyden/app/services/notification/push"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/infrastructure/webpush"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type notificationPushSender struct {
	client *webpush.Client
}

func subscribeNotificationPush(
	ctx context.Context,
	lc fx.Lifecycle,
	bus *pubsub.Bus,
	client *webpush.Client,
	service *push.Service,
) {
	if !client.Enabled() {
		return
	}

	sender := &notificationPushSender{client: client}

	lc.Append(fx.StartHook(func() error {
		_, err := pubsub.Subscribe(ctx, bus, "notification.push", func(ctx context.Context, event *rpc.EventNotificationCreated) error {
			return service.Deliver(ctx, event, sender.Send)
		})
		return err
	}))
}

func (s *notificationPushSender) Send(ctx context.Context, owner account.AccountID, subscription push_subscription.Subscription, item *notification.Notification) error {
	payload, err := encodeNotificationPushPayload(owner, subscription.ID.String(), item)
	if err != nil {
		return err
	}

	return s.client.Send(ctx, subscription.Subscription, payload)
}

type notificationPushPayload struct {
	AccountID      openapi.Identifier   `json:"account_id"`
	SubscriptionID openapi.Identifier   `json:"subscription_id"`
	Notification   openapi.Notification `json:"notification"`
}

func encodeNotificationPushPayload(owner account.AccountID, subscriptionID string, item *notification.Notification) ([]byte, error) {
	value := notificationPushPayload{
		AccountID:      owner.String(),
		SubscriptionID: subscriptionID,
		Notification:   serialiseNotification(item),
	}

	var displayText []*string
	if value.Notification.Item != nil {
		displayText = append(displayText, &value.Notification.Item.Title)
	}

	if value.Notification.Source != nil {
		displayText = append(displayText, &value.Notification.Source.Name)
	}

	for {
		data, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}

		if len(data) <= webpush.MaxPayloadSize {
			return data, nil
		}

		var longest *string
		for _, text := range displayText {
			if longest == nil || len(*text) > len(*longest) {
				longest = text
			}
		}

		if longest == nil || len(*longest) == 0 {
			return nil, webpush.ErrPayloadTooLarge
		}

		runes := []rune(*longest)
		*longest = string(runes[:len(runes)/2])
	}
}
