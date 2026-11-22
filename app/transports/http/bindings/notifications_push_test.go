package bindings

import (
	"encoding/json"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/infrastructure/webpush"
)

func TestPayloadFitsPushDelivery(t *testing.T) {
	owner := account.AccountID(xid.New())
	subscription := xid.New().String()
	title := strings.Repeat("界<&", 10000)
	item := &notification.Notification{
		ID:    xid.New(),
		Event: notification.EventThreadReply,
		Item: opt.New(notification.Item{
			ID:    xid.New(),
			Kind:  datagraph.KindPost,
			Title: title,
		}),
		Source: opt.New(notification.Source{
			ID:     account.AccountID(xid.New()),
			Handle: "odin",
			Name:   "Odin",
		}),
	}

	data, err := encodeNotificationPushPayload(owner, subscription, item)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(data), webpush.MaxPayloadSize)
	assert.True(t, utf8.Valid(data))

	var decoded struct {
		AccountID      string               `json:"account_id"`
		SubscriptionID string               `json:"subscription_id"`
		Notification   openapi.Notification `json:"notification"`
	}
	require.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, owner.String(), decoded.AccountID)
	assert.Equal(t, subscription, decoded.SubscriptionID)
	assert.Equal(t, item.Item.OrZero().ID.String(), decoded.Notification.Item.Id)
	assert.NotEmpty(t, decoded.Notification.Item.Title)
	assert.True(t, strings.HasPrefix(title, decoded.Notification.Item.Title))
	assert.Equal(t, "Odin", decoded.Notification.Source.Name)
	assert.Equal(t, title, item.Item.OrZero().Title)
}
