package webpush_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/account/notification/push_subscription"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/notification/notify"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	entnotification "github.com/Southclaws/storyden/internal/ent/notification"
	"github.com/Southclaws/storyden/internal/infrastructure/webpush"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

type pushTransport struct {
	mu       sync.Mutex
	calls    []string
	statuses map[string]int
	bodies   map[string][]byte
}

func (s *pushTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	endpoint := request.URL.String()
	s.calls = append(s.calls, endpoint)
	s.bodies[endpoint] = body

	status := s.statuses[endpoint]
	if status == 0 {
		status = http.StatusCreated
	}

	return &http.Response{
		StatusCode: status,
		Body:       http.NoBody,
	}, nil
}

func TestNotificationPushDelivery(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)

	var (
		root          context.Context
		db            *ent.Client
		notifier      *notify.Notifier
		accounts      *account_writer.Writer
		subscriptions *push_subscription.Repository
	)

	base := "https://push.example/" + xid.New().String()
	transport := &pushTransport{
		statuses: map[string]int{
			base + "/expired":   http.StatusGone,
			base + "/temporary": http.StatusTooManyRequests,
		},
		bodies: make(map[string][]byte),
	}
	cfg, client := webpush.NewTestClient(t, transport)
	input, receiver := webpush.NewTestSubscription(t, base)
	parsed, err := webpush.ParseSubscription(input)
	r.NoError(err)

	integration.Test(t, &cfg,
		e2e.Setup(),
		fx.Replace(client),
		fx.Populate(&root, &db, &notifier, &accounts, &subscriptions),
	)

	_, owner := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
	_, other := e2e.WithAccount(root, accounts, seed.Account_003_Baldur)
	endpoints := []string{
		base + "/expired",
		base + "/temporary",
		base + "/working",
	}

	subscriptionIDs := make(map[string]string)
	for _, endpoint := range endpoints {
		input.Endpoint = endpoint
		subscription, err := subscriptions.Upsert(root, owner.ID, input)
		r.NoError(err)
		subscriptionIDs[endpoint] = subscription.ID.String()
	}

	input.Endpoint = base + "/other"
	_, err = subscriptions.Upsert(root, other.ID, input)
	r.NoError(err)

	post, err := db.Post.Create().
		SetUpdatedAt(time.Now()).
		SetAuthorID(xid.ID(other.ID)).
		SetTitle("A little harbour").
		SetSlug("a-little-harbour").
		SetBody("<p>" + strings.Repeat("A very large body. ", 10000) + "</p>").
		SetShort("A very large body.").
		SetLastReplyAt(time.Now()).
		Save(root)
	r.NoError(err)

	err = notifier.Send(root,
		owner.ID,
		opt.New(other.ID),
		notification.EventThreadReply,
		&datagraph.Ref{ID: post.ID, Kind: datagraph.KindPost},
	)
	r.NoError(err)

	r.Eventually(func() bool {
		transport.mu.Lock()
		delivered := len(transport.calls) == 3
		transport.mu.Unlock()

		remaining, err := subscriptions.List(root, owner.ID)
		return delivered && err == nil && len(remaining) == 2
	}, 5*time.Second, 10*time.Millisecond)

	transport.mu.Lock()
	a.ElementsMatch(endpoints, transport.calls)
	transport.mu.Unlock()

	storedNotification, err := db.Notification.Query().
		Where(entnotification.OwnerAccountIDEQ(xid.ID(owner.ID))).
		Only(root)
	r.NoError(err)

	for _, endpoint := range endpoints {
		transport.mu.Lock()
		body := transport.bodies[endpoint]
		transport.mu.Unlock()

		plaintext, err := webpush.DecryptTestPayload(t, receiver, parsed.Auth, body)
		r.NoError(err)
		data := bytes.TrimSuffix(plaintext, []byte{0x02})
		var payload struct {
			AccountID      string               `json:"account_id"`
			SubscriptionID string               `json:"subscription_id"`
			Notification   openapi.Notification `json:"notification"`
		}
		r.NoError(json.Unmarshal(data, &payload))
		a.Equal(owner.ID.String(), payload.AccountID)
		a.Equal(subscriptionIDs[endpoint], payload.SubscriptionID)
		a.Equal(storedNotification.ID.String(), payload.Notification.Id)
		a.Equal(openapi.NotificationEvent("thread_reply"), payload.Notification.Event)
		a.Equal(openapi.NotificationStatusUnread, payload.Notification.Status)
		r.NotNil(payload.Notification.Source)
		a.Equal(other.ID.String(), payload.Notification.Source.Id)
		a.Equal(other.Handle, payload.Notification.Source.Handle)
		r.NotNil(payload.Notification.Item)
		a.Equal(post.ID.String(), payload.Notification.Item.Id)
		a.Equal("A little harbour", payload.Notification.Item.Title)
		a.Less(len(data), 3993)

		var raw struct {
			Notification struct {
				Item   map[string]any `json:"item"`
				Source map[string]any `json:"source"`
			} `json:"notification"`
		}
		r.NoError(json.Unmarshal(data, &raw))
		a.NotContains(raw.Notification.Item, "body")
		a.NotContains(raw.Notification.Item, "author")
		a.NotContains(raw.Notification.Source, "roles")
	}

	remaining, err := subscriptions.List(root, owner.ID)
	r.NoError(err)
	r.Len(remaining, 2)

	for _, subscription := range remaining {
		a.NotEqual(base+"/expired", subscription.Endpoint.String())
		a.Equal("https", subscription.Endpoint.Scheme)
		a.Equal("push.example", subscription.Endpoint.Hostname())
	}
}

func TestNotificationPushPreferences(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)

	var (
		root          context.Context
		db            *ent.Client
		notifier      *notify.Notifier
		accounts      *account_writer.Writer
		subscriptions *push_subscription.Repository
	)

	transport := &pushTransport{
		statuses: make(map[string]int),
		bodies:   make(map[string][]byte),
	}
	cfg, client := webpush.NewTestClient(t, transport)
	input, _ := webpush.NewTestSubscription(t, "https://push.example/"+xid.New().String())

	integration.Test(t, &cfg,
		e2e.Setup(),
		fx.Replace(client),
		fx.Populate(&root, &db, &notifier, &accounts, &subscriptions),
	)

	_, owner := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
	_, err := subscriptions.Upsert(root, owner.ID, input)
	r.NoError(err)

	deliveryCount := func() int {
		transport.mu.Lock()
		defer transport.mu.Unlock()

		return len(transport.calls)
	}

	send := func(event notification.Event, expectedInbox int) {
		t.Helper()

		err := notifier.Send(root, owner.ID, opt.NewEmpty[account.AccountID](), event, nil)
		r.NoError(err)

		r.Eventually(func() bool {
			count, err := db.Notification.Query().
				Where(entnotification.OwnerAccountIDEQ(xid.ID(owner.ID))).
				Count(root)
			return err == nil && count == expectedInbox
		}, 5*time.Second, 10*time.Millisecond)
	}

	setPreferences := func(enabled bool, events map[string]account.NotificationEventSettings) {
		t.Helper()

		_, err := accounts.SetNotificationSettings(root, owner.ID, account.NotificationPreferences{
			Events: events,
			Push:   account.PushNotificationPreferences{Enabled: enabled},
		})
		r.NoError(err)
	}

	send(notification.EventThreadReply, 1)
	r.Eventually(func() bool {
		return deliveryCount() == 1
	}, 5*time.Second, 10*time.Millisecond)

	setPreferences(true, map[string]account.NotificationEventSettings{"thread_reply": {Inbox: true, Push: false}})
	send(notification.EventThreadReply, 2)
	a.Never(func() bool {
		return deliveryCount() != 1
	}, 200*time.Millisecond, 10*time.Millisecond)

	send(notification.EventReplyToReply, 3)
	r.Eventually(func() bool {
		return deliveryCount() == 2
	}, 5*time.Second, 10*time.Millisecond)

	setPreferences(true, map[string]account.NotificationEventSettings{"thread_reply": {Inbox: false, Push: false}})
	send(notification.EventThreadReply, 3)
	a.Never(func() bool {
		count, err := db.Notification.Query().
			Where(entnotification.OwnerAccountIDEQ(xid.ID(owner.ID))).
			Count(root)
		return err != nil || count != 3 || deliveryCount() != 2
	}, 200*time.Millisecond, 10*time.Millisecond)

	setPreferences(false, map[string]account.NotificationEventSettings{"thread_reply": {Inbox: true, Push: true}})
	send(notification.EventThreadReply, 4)
	a.Never(func() bool {
		return deliveryCount() != 2
	}, 200*time.Millisecond, 10*time.Millisecond)

	setPreferences(true, map[string]account.NotificationEventSettings{"thread_reply": {Inbox: true, Push: true}})
	send(notification.EventThreadReply, 5)
	r.Eventually(func() bool {
		return deliveryCount() == 3
	}, 5*time.Second, 10*time.Millisecond)

	send(notification.EventPostLike, 5)
	a.Never(func() bool {
		count, err := db.Notification.Query().
			Where(entnotification.OwnerAccountIDEQ(xid.ID(owner.ID))).
			Count(root)
		return err != nil || count != 5 || deliveryCount() != 3
	}, 200*time.Millisecond, 10*time.Millisecond)

	setPreferences(true, map[string]account.NotificationEventSettings{"post_like": {Inbox: true, Push: true}})
	send(notification.EventPostLike, 6)
	r.Eventually(func() bool {
		return deliveryCount() == 4
	}, 5*time.Second, 10*time.Millisecond)

	remaining, err := subscriptions.List(root, owner.ID)
	r.NoError(err)
	a.Len(remaining, 1)

	inbox, err := db.Notification.Query().
		Where(entnotification.OwnerAccountIDEQ(xid.ID(owner.ID))).
		All(root)
	r.NoError(err)
	r.Len(inbox, 6)

	for _, item := range inbox {
		a.False(item.Read)
	}
}
