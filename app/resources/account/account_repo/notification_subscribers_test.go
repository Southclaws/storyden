package account_repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/notification/settings"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestNotificationSubscribers(t *testing.T) {
	if tests.IsSharedPostgresDatabase() {
		t.Skip("global notification subscriptions affect parallel integration apps")
	}

	t.Parallel()

	r := require.New(t)
	a := assert.New(t)
	var (
		root        context.Context
		db          *ent.Client
		accounts    *account_writer.Writer
		querier     *account_querier.Querier
		preferences *settings.Service
	)

	integration.Test(t, nil, fx.Populate(&root, &db, &accounts, &querier, &preferences))

	_, enabled := e2e.WithAccount(root, accounts, seed.Account_003_Baldur)
	_, disabled := e2e.WithAccount(root, accounts, seed.Account_004_Loki)
	_, unset := e2e.WithAccount(root, accounts, seed.Account_006_Freyja)
	_, otherEvent := e2e.WithAccount(root, accounts, seed.Account_007_Freyr)

	setEvent := func(id account.AccountID, event notification.Event, inbox bool) {
		t.Helper()

		_, err := preferences.Update(root, id, account.NotificationPreferences{
			Events: map[string]account.NotificationEventSettings{
				event.String(): {Inbox: inbox, Push: false},
			},
			Push: account.PushNotificationPreferences{Enabled: true},
		})
		r.NoError(err)
	}

	setEvent(enabled.ID, notification.EventGlobalThreadPublished, true)
	setEvent(disabled.ID, notification.EventGlobalThreadPublished, false)
	setEvent(otherEvent.ID, notification.EventGlobalPostReacted, true)
	subscribers, err := querier.ListNotificationSubscribers(root, "global_thread_published")
	r.NoError(err)
	ids := make(map[account.AccountID]bool)
	for _, subscriber := range subscribers {
		ids[subscriber.ID] = true
	}

	a.True(ids[enabled.ID])
	a.False(ids[disabled.ID])
	a.False(ids[unset.ID])
	a.False(ids[otherEvent.ID])

	err = db.Account.UpdateOneID(xid.ID(enabled.ID)).SetDeletedAt(time.Now()).Exec(root)
	r.NoError(err)
	subscribers, err = querier.ListNotificationSubscribers(root, "global_thread_published")
	r.NoError(err)
	for _, subscriber := range subscribers {
		a.NotEqual(enabled.ID, subscriber.ID)
	}
}
