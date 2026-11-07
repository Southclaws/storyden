package notification_test

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/notification/push_subscription"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/infrastructure/webpush"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestPushSubscriptionEviction(t *testing.T) {
	for _, limit := range []int{5, 2, 1} {
		t.Run(strconv.Itoa(limit), func(t *testing.T) {
			t.Parallel()

			r := require.New(t)
			a := assert.New(t)

			var (
				root          context.Context
				cl            *openapi.ClientWithResponses
				sh            *e2e.SessionHelper
				db            *ent.Client
				accounts      *account_writer.Writer
				subscriptions *push_subscription.Repository
			)

			cfg := webPushConfig(t)
			cfg.WebPushMaxSubscriptions = limit
			integration.Test(t, cfg,
				e2e.Setup(),
				fx.Populate(&root, &cl, &sh, &db, &accounts, &subscriptions),
			)

			ctx, owner := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
			ownerSession := sh.WithSession(ctx)
			_, other := e2e.WithAccount(root, accounts, seed.Account_003_Baldur)
			base := "https://push.example/" + owner.ID.String() + "/"
			input := openapi.PushSubscriptionUpsertJSONRequestBody{
				Endpoint: base + "0",
			}
			input.Keys.P256dh = cfg.WebPushVAPIDPublicKey
			input.Keys.Auth = "AAAAAAAAAAAAAAAAAAAAAA"

			otherSubscription, err := subscriptions.Upsert(root, other.ID, webpush.SubscriptionInput{
				Endpoint: "https://push.example/" + other.ID.String(),
				P256DH:   input.Keys.P256dh,
				Auth:     input.Keys.Auth,
			})
			r.NoError(err)

			before, err := db.Account.Get(root, xid.ID(owner.ID))
			r.NoError(err)

			ids := make([]string, 0, limit)
			for i := range limit {
				input.Endpoint = base + strconv.Itoa(i)
				created := tests.AssertRequest(
					cl.PushSubscriptionUpsertWithResponse(root, input, ownerSession),
				)(t, http.StatusOK)
				r.NotNil(created.JSON200)
				ids = append(ids, created.JSON200.Id)
			}

			input.Endpoint = base + "0"
			input.Keys.Auth = "AQEBAQEBAQEBAQEBAQEBAQ"
			updated := tests.AssertRequest(
				cl.PushSubscriptionUpsertWithResponse(root, input, ownerSession),
			)(t, http.StatusOK)
			r.NotNil(updated.JSON200)
			a.Equal(ids[0], updated.JSON200.Id)

			stored, err := subscriptions.List(root, owner.ID)
			r.NoError(err)
			r.Len(stored, limit)

			input.Endpoint = base + "new"
			created := tests.AssertRequest(
				cl.PushSubscriptionUpsertWithResponse(root, input, ownerSession),
			)(t, http.StatusOK)
			r.NotNil(created.JSON200)

			stored, err = subscriptions.List(root, owner.ID)
			r.NoError(err)
			r.Len(stored, limit)

			actual := make([]string, 0, len(stored))
			for _, subscription := range stored {
				actual = append(actual, subscription.ID.String())
			}

			expected := append(ids[1:], created.JSON200.Id)
			a.ElementsMatch(expected, actual)
			a.NotContains(actual, ids[0])

			otherStored, err := subscriptions.List(root, other.ID)
			r.NoError(err)
			r.Len(otherStored, 1)
			a.Equal(otherSubscription.ID, otherStored[0].ID)

			after, err := db.Account.Get(root, xid.ID(owner.ID))
			r.NoError(err)
			a.Equal(before.UpdatedAt, after.UpdatedAt)
		})
	}
}

func TestPushSubscriptionConcurrentRegistration(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)

	var (
		root          context.Context
		db            *ent.Client
		accounts      *account_writer.Writer
		subscriptions *push_subscription.Repository
	)

	cfg := webPushConfig(t)
	integration.Test(t, cfg,
		fx.Populate(&root, &db, &accounts, &subscriptions),
	)

	_, owner := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
	otherRepository, err := push_subscription.New(db, *cfg)
	r.NoError(err)

	repositories := []*push_subscription.Repository{subscriptions, otherRepository}
	const registrations = 16
	start := make(chan struct{})
	errors := make(chan error, registrations)

	for i := range registrations {
		go func() {
			<-start

			repository := repositories[i%len(repositories)]
			_, err := repository.Upsert(root, owner.ID, webpush.SubscriptionInput{
				Endpoint: "https://push.example/" + owner.ID.String() + "/" + strconv.Itoa(i),
				P256DH:   cfg.WebPushVAPIDPublicKey,
				Auth:     "AAAAAAAAAAAAAAAAAAAAAA",
			})
			errors <- err
		}()
	}

	close(start)
	for range registrations {
		a.NoError(<-errors)
	}

	stored, err := subscriptions.List(root, owner.ID)
	r.NoError(err)
	a.Len(stored, 5)
}
