package webpush_test

import (
	"context"
	"io"
	"net/http"
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
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/notification/notify"
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

	for _, endpoint := range endpoints {
		input.Endpoint = endpoint
		_, err := subscriptions.Upsert(root, owner.ID, input)
		r.NoError(err)
	}

	input.Endpoint = base + "/other"
	_, err = subscriptions.Upsert(root, other.ID, input)
	r.NoError(err)

	err = notifier.Send(root,
		owner.ID,
		opt.NewEmpty[account.AccountID](),
		notification.EventFollow,
		nil,
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
		a.Equal(storedNotification.ID.String()+"\x02", string(plaintext))
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
