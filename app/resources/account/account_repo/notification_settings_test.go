package account_repo_test

import (
	"context"
	"sync"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/notificationeventsettings"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

func TestNotificationSettingsStorage(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)
	var (
		root    context.Context
		db      *ent.Client
		writer  *account_writer.Writer
		querier *account_querier.Querier
	)

	integration.Test(t, nil, fx.Populate(&root, &db, &writer, &querier))
	settings := account.NotificationPreferences{
		Events: map[string]account.NotificationEventSettings{
			"thread_reply":         {Inbox: false, Push: false},
			"post_like":            {Inbox: true, Push: false},
			"node_version_applied": {Inbox: true, Push: true},
		},
		Push: account.PushNotificationPreferences{Enabled: false},
	}
	_, owner := e2e.WithAccount(root, writer, seed.Account_001_Odin)
	loaded, err := querier.GetNotificationSettings(root, owner.ID)
	r.NoError(err)
	a.Empty(loaded.Events)
	rows := db.NotificationEventSettings.Query().
		Where(notificationeventsettings.AccountIDEQ(xid.ID(owner.ID)))
	count, err := rows.Clone().Count(root)
	r.NoError(err)
	a.Zero(count)

	loaded, err = writer.SetNotificationSettings(root, owner.ID, settings)
	r.NoError(err)
	a.Equal(settings, loaded)

	invalid := account.NotificationPreferences{
		Events: map[string]account.NotificationEventSettings{
			"": {Inbox: true, Push: true},
		},
		Push: account.PushNotificationPreferences{Enabled: true},
	}
	_, err = writer.SetNotificationSettings(root, owner.ID, invalid)
	r.Error(err)
	unchanged, err := querier.GetNotificationSettings(root, owner.ID)
	r.NoError(err)
	a.Equal(settings, unchanged)

	row, err := db.NotificationEventSettings.Create().
		SetAccountID(xid.ID(owner.ID)).
		SetEvent("reply_to_reply").
		SetInboxEnabled(false).
		SetPushEnabled(true).
		Save(root)
	r.NoError(err)
	loaded, err = querier.GetNotificationSettings(root, owner.ID)
	r.NoError(err)
	a.Equal(account.NotificationEventSettings{}, loaded.Events[row.Event])

	reset := account.NotificationPreferences{Push: account.PushNotificationPreferences{Enabled: true}}
	_, err = writer.SetNotificationSettings(root, owner.ID, reset)
	r.NoError(err)
	loaded, err = querier.GetNotificationSettings(root, owner.ID)
	r.NoError(err)
	a.True(loaded.Push.Enabled)
	a.Empty(loaded.Events)
	count, err = rows.Clone().Count(root)
	r.NoError(err)
	a.Zero(count)
}

func TestNotificationSettingsConcurrentReplacement(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)
	var (
		root    context.Context
		db      *ent.Client
		writer  *account_writer.Writer
		querier *account_querier.Querier
	)

	integration.Test(t, nil, fx.Populate(&root, &db, &writer, &querier))
	_, owner := e2e.WithAccount(root, writer, seed.Account_001_Odin)
	choices := []account.NotificationPreferences{
		{
			Events: map[string]account.NotificationEventSettings{
				"thread_reply":   {Inbox: true, Push: false},
				"reply_to_reply": {Inbox: false, Push: false},
			},
			Push: account.PushNotificationPreferences{Enabled: false},
		},
		{
			Events: map[string]account.NotificationEventSettings{
				"post_like": {Inbox: true, Push: true},
				"follow":    {Inbox: false, Push: false},
			},
			Push: account.PushNotificationPreferences{Enabled: true},
		},
	}
	var group sync.WaitGroup
	errors := make(chan error, 8)
	for i := range 8 {
		group.Go(func() {
			_, err := writer.SetNotificationSettings(root, owner.ID, choices[i%len(choices)])
			errors <- err
		})
	}
	group.Wait()
	close(errors)
	for err := range errors {
		r.NoError(err)
	}

	loaded, err := querier.GetNotificationSettings(root, owner.ID)
	r.NoError(err)
	a.Contains(choices, loaded)
	count, err := db.NotificationEventSettings.Query().
		Where(notificationeventsettings.AccountIDEQ(xid.ID(owner.ID))).
		Count(root)
	r.NoError(err)
	a.Equal(2, count)
}
