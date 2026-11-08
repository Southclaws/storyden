package account_repo

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
	entaccount "github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/notificationeventsettings"
)

func (r *Repository) GetNotificationSettings(ctx context.Context, id account.AccountID) (account.NotificationPreferences, error) {
	row, err := r.db.Account.Query().
		Where(entaccount.IDEQ(xid.ID(id))).
		Select(entaccount.FieldPushNotificationsEnabled).
		WithNotificationEventSettings().
		Only(ctx)
	if err != nil {
		return account.NotificationPreferences{}, fault.Wrap(err, fctx.With(ctx))
	}

	return account.MapNotificationSettings(row.PushNotificationsEnabled, row.Edges.NotificationEventSettings), nil
}

func (r *Repository) SetNotificationSettings(ctx context.Context, id account.AccountID, settings account.NotificationPreferences) (account.NotificationPreferences, error) {
	err := ent.WithTx(ctx, r.db, func(tx *ent.Tx) error {
		if err := tx.Account.UpdateOneID(xid.ID(id)).
			SetPushNotificationsEnabled(settings.Push.Enabled).
			Exec(ctx); err != nil {
			return err
		}

		return replaceNotificationEventSettings(ctx, tx, xid.ID(id), settings.Events)
	})
	if err != nil {
		return account.NotificationPreferences{}, fault.Wrap(err, fctx.With(ctx))
	}

	return r.GetNotificationSettings(ctx, id)
}

func (r *Repository) ListNotificationSubscribers(ctx context.Context, event string) ([]account.NotificationSubscriber, error) {
	rows, err := r.db.Account.Query().
		Where(
			entaccount.DeletedAtIsNil(),
			entaccount.HasNotificationEventSettingsWith(
				notificationeventsettings.EventEQ(event),
				notificationeventsettings.InboxEnabledEQ(true),
			),
		).
		Select(entaccount.FieldID, entaccount.FieldPushNotificationsEnabled).
		WithNotificationEventSettings().
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	subscribers := make([]account.NotificationSubscriber, 0, len(rows))
	for _, row := range rows {
		subscribers = append(subscribers, account.NotificationSubscriber{
			ID:       account.AccountID(row.ID),
			Settings: account.MapNotificationSettings(row.PushNotificationsEnabled, row.Edges.NotificationEventSettings),
		})
	}

	return subscribers, nil
}

func replaceNotificationEventSettings(
	ctx context.Context,
	tx *ent.Tx,
	accountID xid.ID,
	events map[string]account.NotificationEventSettings,
) error {
	_, err := tx.NotificationEventSettings.Delete().
		Where(notificationeventsettings.AccountIDEQ(accountID)).
		Exec(ctx)
	if err != nil {
		return err
	}

	if len(events) == 0 {
		return nil
	}

	rows := make([]*ent.NotificationEventSettingsCreate, 0, len(events))
	for event, settings := range events {
		row := tx.NotificationEventSettings.Create().
			SetAccountID(accountID).
			SetEvent(event).
			SetInboxEnabled(settings.Inbox).
			SetPushEnabled(settings.Push)
		rows = append(rows, row)
	}

	return tx.NotificationEventSettings.CreateBulk(rows...).Exec(ctx)
}
