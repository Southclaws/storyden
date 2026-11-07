package notification_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/account/notification/push_subscription"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/notification/notify"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/infrastructure/webpush"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
	"github.com/Southclaws/storyden/tests"
)

func webPushConfig(t *testing.T) *config.Config {
	t.Helper()

	r := require.New(t)

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	r.NoError(err)

	private, err := key.Bytes()
	r.NoError(err)

	public, err := key.PublicKey.Bytes()
	r.NoError(err)

	return &config.Config{
		WebPushEnabled:          true,
		WebPushMaxSubscriptions: 5,
		WebPushVAPIDPrivateKey:  base64.RawURLEncoding.EncodeToString(private),
		WebPushVAPIDPublicKey:   base64.RawURLEncoding.EncodeToString(public),
		WebPushVAPIDSubject:     "mailto:admin@example.com",
	}
}

func TestPushSubscriptionAPI(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)

	var (
		root     context.Context
		cl       *openapi.ClientWithResponses
		sh       *e2e.SessionHelper
		accounts *account_writer.Writer
	)

	cfg := webPushConfig(t)
	integration.Test(t, cfg,
		e2e.Setup(),
		fx.Populate(&root, &cl, &sh, &accounts),
	)

	ctx, owner := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
	ownerSession := sh.WithSession(ctx)
	otherCtx, _ := e2e.WithAccount(root, accounts, seed.Account_003_Baldur)
	otherSession := sh.WithSession(otherCtx)

	input := openapi.PushSubscriptionUpsertJSONRequestBody{
		Endpoint: "https://push.example/" + owner.ID.String(),
	}
	input.Keys.P256dh = webPushConfig(t).WebPushVAPIDPublicKey
	input.Keys.Auth = "AAAAAAAAAAAAAAAAAAAAAA"

	info := tests.AssertRequest(
		cl.GetInfoWithResponse(root),
	)(t, http.StatusOK)
	r.NotNil(info.JSON200.WebPushPublicKey)
	a.Equal(cfg.WebPushVAPIDPublicKey, *info.JSON200.WebPushPublicKey)

	tests.AssertRequest(
		cl.PushSubscriptionUpsertWithResponse(root, input),
	)(t, http.StatusUnauthorized)
	created := tests.AssertRequest(
		cl.PushSubscriptionUpsertWithResponse(root, input, ownerSession),
	)(t, http.StatusOK)
	r.NotNil(created.JSON200)
	a.Equal(input.Endpoint, created.JSON200.Endpoint)

	input.Keys.Auth = "AQEBAQEBAQEBAQEBAQEBAQ"
	updated := tests.AssertRequest(
		cl.PushSubscriptionUpsertWithResponse(root, input, ownerSession),
	)(t, http.StatusOK)
	a.Equal(created.JSON200.Id, updated.JSON200.Id)

	tests.AssertRequest(
		cl.PushSubscriptionUpsertWithResponse(root, input, otherSession),
	)(t, http.StatusConflict)
	tests.AssertRequest(
		cl.PushSubscriptionDeleteWithResponse(root, created.JSON200.Id, otherSession),
	)(t, http.StatusNotFound)

	input.Endpoint = "https://127.0.0.1/internal"
	tests.AssertRequest(
		cl.PushSubscriptionUpsertWithResponse(root, input, ownerSession),
	)(t, http.StatusBadRequest)
	tests.AssertRequest(
		cl.PushSubscriptionDeleteWithResponse(root, created.JSON200.Id, ownerSession),
	)(t, http.StatusNoContent)
	tests.AssertRequest(
		cl.PushSubscriptionDeleteWithResponse(root, created.JSON200.Id, ownerSession),
	)(t, http.StatusNotFound)
}

func TestPushSubscriptionEndpoint(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)

	var (
		root          context.Context
		db            *ent.Client
		accounts      *account_writer.Writer
		subscriptions *push_subscription.Repository
	)

	integration.Test(t, &config.Config{},
		fx.Populate(&root, &db, &accounts, &subscriptions),
	)

	_, owner := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
	endpoint := "https://push.example/" + owner.ID.String() + "/token%2Fvalue?key=a%2Bb"
	input := webpush.SubscriptionInput{
		Endpoint: endpoint,
		P256DH:   webPushConfig(t).WebPushVAPIDPublicKey,
		Auth:     "AAAAAAAAAAAAAAAAAAAAAA",
	}

	subscription, err := subscriptions.Upsert(root, owner.ID, input)
	r.NoError(err)
	a.Equal(endpoint, subscription.Endpoint.String())

	stored, err := subscriptions.List(root, owner.ID)
	r.NoError(err)
	r.Len(stored, 1)
	a.Equal(endpoint, stored[0].Endpoint.String())
	a.Equal("https", stored[0].Endpoint.Scheme)
	a.Equal("push.example", stored[0].Endpoint.Hostname())

	input.Endpoint = "http://push.example/token"
	_, err = subscriptions.Upsert(root, owner.ID, input)
	r.Error(err)

	err = db.PushSubscription.UpdateOneID(subscription.ID).
		SetEndpoint("invalid endpoint").
		Exec(root)
	r.NoError(err)

	_, err = subscriptions.List(root, owner.ID)
	r.Error(err)
}

func TestPushDisabled(t *testing.T) {
	t.Parallel()

	var (
		root          context.Context
		cl            *openapi.ClientWithResponses
		sh            *e2e.SessionHelper
		accounts      *account_writer.Writer
		subscriptions *push_subscription.Repository
	)

	integration.Test(t, &config.Config{
		WebPushVAPIDPublicKey:  "invalid",
		WebPushVAPIDPrivateKey: "invalid",
		WebPushVAPIDSubject:    "invalid",
	},
		e2e.Setup(),
		fx.Populate(&root, &cl, &sh, &accounts, &subscriptions),
	)

	ctx, owner := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
	ownerSession := sh.WithSession(ctx)
	info := tests.AssertRequest(
		cl.GetInfoWithResponse(root),
	)(t, http.StatusOK)
	assert.Nil(t, info.JSON200.WebPushPublicKey)

	input := openapi.PushSubscriptionUpsertJSONRequestBody{
		Endpoint: "https://push.example/" + owner.ID.String(),
	}
	input.Keys.P256dh = webPushConfig(t).WebPushVAPIDPublicKey
	input.Keys.Auth = "AAAAAAAAAAAAAAAAAAAAAA"
	tests.AssertRequest(
		cl.PushSubscriptionUpsertWithResponse(root, input, ownerSession),
	)(t, http.StatusForbidden)

	existing, err := subscriptions.Upsert(root, owner.ID, webpush.SubscriptionInput{
		Endpoint: input.Endpoint,
		P256DH:   input.Keys.P256dh,
		Auth:     input.Keys.Auth,
	})
	require.NoError(t, err)
	tests.AssertRequest(
		cl.PushSubscriptionDeleteWithResponse(root, existing.ID.String(), ownerSession),
	)(t, http.StatusNoContent)
}

func TestNotificationCreatedEvent(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)

	var (
		root     context.Context
		bus      *pubsub.Bus
		notifier *notify.Notifier
		accounts *account_writer.Writer
		db       *ent.Client
	)

	integration.Test(t, nil, fx.Populate(&root, &bus, &notifier, &accounts, &db))

	_, owner := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
	created := make(chan rpc.EventNotificationCreated, 1)

	sub, err := pubsub.Subscribe(root, bus, "test.notification.created", func(_ context.Context, event *rpc.EventNotificationCreated) error {
		if event.AccountID == owner.ID {
			created <- *event
		}
		return nil
	})
	r.NoError(err)
	t.Cleanup(func() {
		r.NoError(sub.Close())
	})

	err = notifier.Send(root,
		owner.ID,
		opt.NewEmpty[account.AccountID](),
		notification.EventFollow,
		nil,
	)
	r.NoError(err)

	select {
	case event := <-created:
		row, err := db.Notification.Get(root, event.ID)
		r.NoError(err, "event must reference a persisted notification")
		a.Equal(owner.ID.String(), row.OwnerAccountID.String())
		a.Equal("follow", row.EventType)

	case <-time.After(5 * time.Second):
		t.Fatal("notification creation event was not delivered")
	}
}
