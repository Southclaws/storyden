package notification

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model"
	model_account "github.com/Southclaws/storyden/backend/internal/infrastructure/db/model/account"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model/notification"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model/subscription"
	"github.com/Southclaws/storyden/backend/pkg/resources/account"
	"github.com/Southclaws/storyden/backend/pkg/resources/post"
)

type database struct {
	db *model.Client
}

func New(db *model.Client) Repository {
	return &database{db}
}

func (d *database) Subscribe(ctx context.Context, accountID account.AccountID, refersType NotificationType, refersTo string) (*Subscription, error) {
	sub, err := d.db.Subscription.Create().
		SetAccountID(uuid.UUID(accountID)).
		SetRefersType(string(refersType)).
		SetRefersTo(refersTo).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return SubFromModel(sub), nil
}

func (d *database) Unsubscribe(ctx context.Context, accountID account.AccountID, subID SubscriptionID) (int, error) {
	i, err := d.db.Subscription.Delete().
		Where(
			subscription.IDEQ(uuid.UUID(subID)),
			subscription.HasAccountWith(
				model_account.IDEQ(uuid.UUID(accountID)),
			),
		).Exec(ctx)
	if err != nil {
		return 0, err
	}

	return i, nil
}

func (d *database) GetSubscriptionsForUser(ctx context.Context, accountID account.AccountID) ([]Subscription, error) {
	subs, err := d.db.Subscription.Query().
		Where(
			subscription.HasAccountWith(
				model_account.IDEQ(uuid.UUID(accountID)),
			),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return SubFromModelMany(subs), nil
}

func (d *database) GetSubscriptionsForItem(ctx context.Context, refersType NotificationType, refersTo string) ([]Subscription, error) {
	subs, err := d.db.Subscription.Query().
		Where(
			subscription.RefersTypeEQ(string(refersType)),
			subscription.RefersToEQ(refersTo),
			subscription.DeleteTimeIsNil(),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return SubFromModelMany(subs), nil
}

func (d *database) GetNotifications(ctx context.Context, accountID account.AccountID, read bool, after time.Time) ([]Notification, error) {
	q := d.db.Notification.Query().
		Where(
			notification.HasSubscriptionWith(
				subscription.HasAccountWith(
					model_account.IDEQ(uuid.UUID(accountID)),
				),
			),
		).
		WithSubscription()

	// if read is false (default), only return unread notifications, otherwise, all.
	if !read {
		q.Where(notification.ReadEQ(false))
	}

	if !after.IsZero() {
		q.Where(notification.CreateTimeGT(after))
	}

	notifs, err := q.All(ctx)
	if err != nil {
		return nil, err
	}

	return FromModelMany(notifs), nil
}

func (d *database) Notify(ctx context.Context, refersType NotificationType, refersTo string, title, desc, link string) (int, error) {
	// NOTE: This is extremely inefficient for large forums!
	// TODO: Figure out a better way to do this. There are two options:
	//       1. A message queue, just to defer the database ops. This would
	//          effectively be the same code, just spread over time.
	//       2. A new table that stores subscription notifiation events that
	//          aren't associated with specific users. Then, when a user queries
	//          their notifications list, this table is checked for items they
	//          are subscribed to and notifications are generated.
	subs, err := d.GetSubscriptionsForItem(ctx, refersType, refersTo)
	if err != nil {
		return 0, err
	}

	for _, sub := range subs {
		_, err := d.db.Notification.
			Create().
			SetTitle(title).
			SetDescription(desc).
			SetLink(link).
			SetRead(false).
			SetSubscriptionID(uuid.UUID(sub.ID)).
			Save(ctx)
		if err != nil {
			return 0, err
		}
	}

	return len(subs), nil
}

func (d *database) SetReadState(ctx context.Context, accountID account.AccountID, notificationID NotificationID, read bool) (*Notification, error) {
	ok, err := d.userHasRightsForNotification(ctx, accountID, notificationID)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, post.ErrUnauthorised
	}

	notif, err := d.db.Notification.
		UpdateOneID(uuid.UUID(notificationID)).
		SetRead(read).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return FromModel(notif), nil
}

// TODO: Cache these. Or do more clever queries.
func (d *database) userHasRightsForNotification(ctx context.Context, accountID account.AccountID, notificationID NotificationID) (bool, error) {
	n, err := d.db.Notification.Query().
		Where(notification.IDEQ(uuid.UUID(notificationID))).
		WithSubscription(func(sq *model.SubscriptionQuery) {
			sq.WithAccount()
		}).
		Only(ctx)
	if err != nil {
		return false, err
	}

	return account.AccountID(n.Edges.Subscription.Edges.Account.ID) == accountID, nil
}

func (d *database) Delete(ctx context.Context, accountID account.AccountID, notificationID NotificationID) (*Notification, error) {
	ok, err := d.userHasRightsForNotification(ctx, accountID, notificationID)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, post.ErrUnauthorised
	}

	n, err := d.db.Notification.Get(ctx, uuid.UUID(notificationID))
	if err != nil {
		return nil, err
	}

	err = d.db.Notification.DeleteOne(n).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return FromModel(n), nil
}
