package notification

import (
	"context"
	"time"

	"github.com/Southclaws/storyden/api/src/infra/db"
	"github.com/Southclaws/storyden/api/src/infra/db/model"
	"github.com/Southclaws/storyden/api/src/resources/post"
)

type DB struct {
	db *model.Client
}

func New(db *model.Client) Repository {
	return &DB{db}
}

func (d *DB) Subscribe(ctx context.Context, userID string, refersType NotificationType, refersTo string) (*Subscription, error) {
	sub, err := d.db.Subscription.CreateOne(
		db.Subscription.User.Link(db.User.ID.Equals(userID)),
		db.Subscription.RefersType.Set(db.NotificationType(refersType)),
		db.Subscription.RefersTo.Set(refersTo),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return SubFromModel(sub), nil
}

func (d *DB) Unsubscribe(ctx context.Context, userID, subID string) (int, error) {
	updateSubscription := d.db.Subscription.
		FindMany(
			db.Subscription.ID.Equals(subID),
			db.Subscription.UserID.Equals(userID),
		).
		With(
			db.Subscription.Notifications.Fetch(),
		).
		Update(
			db.Subscription.DeletedAt.Set(time.Now()),
		).
		Tx()

	deleteNotifications := d.db.Notification.
		FindMany(
			db.Notification.SubscriptionID.Equals(subID),
		).
		Delete().
		Tx()

	err := d.db.Prisma.TX.Transaction(
		updateSubscription,
		deleteNotifications,
	).Exec(ctx)
	if err != nil {
		return 0, err
	}

	return deleteNotifications.Result().Count, nil
}

func (d *DB) GetSubscriptionsForUser(ctx context.Context, userID string) ([]Subscription, error) {
	subs, err := d.db.Subscription.
		FindMany(
			db.Subscription.UserID.Equals(userID),
			db.Subscription.DeletedAt.IsNull(),
		).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	result := []Subscription{}
	for _, s := range subs {
		result = append(result, *SubFromModel(&s))
	}
	return result, nil
}

func (d *DB) GetSubscriptionsForItem(ctx context.Context, refersType NotificationType, refersTo string) ([]Subscription, error) {
	subs, err := d.db.Subscription.
		FindMany(
			db.Subscription.RefersType.Equals(db.NotificationType(refersType)),
			db.Subscription.RefersTo.Equals(refersTo),
			db.Subscription.DeletedAt.IsNull(),
		).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	result := []Subscription{}
	for _, s := range subs {
		result = append(result, *SubFromModel(&s))
	}
	return result, nil
}

func (d *DB) GetNotifications(ctx context.Context, userID string, read bool, after time.Time) ([]Notification, error) {
	filters := []db.NotificationWhereParam{
		db.Notification.Subscription.Where(
			db.Subscription.UserID.Equals(userID),
		),
	}

	// if read is false (default), only return unread notifications, otherwise, all.
	if !read {
		filters = append(filters, db.Notification.Read.Equals(false))
	}

	if !after.IsZero() {
		filters = append(filters, db.Notification.CreatedAt.After(after))
	}

	notifs, err := d.db.Notification.
		FindMany(filters...).
		With(db.Notification.Subscription.Fetch()).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	result := []Notification{}
	for _, n := range notifs {
		result = append(result, *FromModel(&n))
	}

	return result, nil
}

func (d *DB) Notify(ctx context.Context, refersType NotificationType, refersTo string, title, desc, link string) (int, error) {
	// NOTE: This is extremely inefficient for large forums!
	// TODO: Figure out a better way to do this. There are two options:
	//       1. A message queue, just to defer the DB ops. This would
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
			CreateOne(
				db.Notification.Title.Set(title),
				db.Notification.Description.Set(desc),
				db.Notification.Link.Set(link),
				db.Notification.Read.Set(false),
				db.Notification.Subscription.Link(
					db.Subscription.ID.Equals(sub.ID),
				),
			).
			Exec(ctx)
		if err != nil {
			return 0, err
		}
	}

	return len(subs), nil
}

func (d *DB) SetReadState(ctx context.Context, userID string, notificationID string, read bool) (*Notification, error) {
	ok, err := d.userHasRightsForNotification(ctx, userID, notificationID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, post.ErrUnauthorised
	}

	notif, err := d.db.Notification.
		FindUnique(
			db.Notification.ID.Equals(notificationID),
		).
		Update(
			db.Notification.Read.Set(read),
		).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return FromModel(notif), nil
}

// TODO: Cache these. Or do more clever queries.
func (d *DB) userHasRightsForNotification(ctx context.Context, userID, notificationID string) (bool, error) {
	n, err := d.db.Notification.
		FindUnique(db.Notification.ID.Equals(notificationID)).
		With(db.Notification.Subscription.Fetch().With(db.Subscription.User.Fetch())).
		Exec(ctx)
	if err != nil {
		return false, err
	}
	return n.Subscription().User().ID == userID, nil
}

func (d *DB) Delete(ctx context.Context, userID, notificationID string) (*Notification, error) {
	ok, err := d.userHasRightsForNotification(ctx, userID, notificationID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, post.ErrUnauthorised
	}

	n, err := d.db.Notification.
		FindUnique(db.Notification.ID.Equals(notificationID)).
		Delete().
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return FromModel(n), nil
}
