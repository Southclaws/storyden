package notification

import (
	"context"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/notification"
)

type database struct {
	db *ent.Client
}

func New(db *ent.Client) Repository {
	return &database{db}
}

func (d *database) List(ctx context.Context, accountID account.AccountID, read bool, after time.Time) ([]*Notification, error) {
	q := d.db.Notification.Query().
		Where(
		// notification.HasSubscriptionWith(
		// 	subscription.HasAccountWith(
		// 		model_account.IDEQ(xid.ID(accountID)),
		// 	),
		// ),
		)

	// if read is false (default), only return unread notifications, otherwise, all.
	if !read {
		q.Where(notification.ReadEQ(false))
	}

	if !after.IsZero() {
		q.Where(notification.CreatedAtGT(after))
	}

	notifs, err := q.All(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}

		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return dt.Map(notifs, FromModel), nil
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
	// subs, err := d.db.Subscription.Query().
	// 	Where(
	// 		subscription.RefersTypeEQ(string(refersType)),
	// 		subscription.RefersToEQ(refersTo),
	// 	).
	// 	All(ctx)
	// if err != nil {
	// 	return 0, fault.Wrap(err, fctx.With(ctx))
	// }

	// for _, sub := range subs {
	// 	_, err := d.db.Notification.
	// 		Create().
	// 		SetTitle(title).
	// 		SetDescription(desc).
	// 		SetLink(link).
	// 		SetRead(false).
	// 		Save(ctx)
	// 	if err != nil {
	// 		return 0, fault.Wrap(err, fctx.With(ctx))
	// 	}
	// }

	return 0, nil
}

func (d *database) SetReadState(ctx context.Context, accountID account.AccountID, notificationID NotificationID, read bool) (*Notification, error) {
	notif, err := d.db.Notification.
		UpdateOneID(xid.ID(notificationID)).
		SetRead(read).
		Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return FromModel(notif), nil
}

func (d *database) Delete(ctx context.Context, accountID account.AccountID, notificationID NotificationID) (*Notification, error) {
	n, err := d.db.Notification.Get(ctx, xid.ID(notificationID))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = d.db.Notification.DeleteOne(n).Exec(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return FromModel(n), nil
}
