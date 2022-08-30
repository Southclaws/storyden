package notification

import (
	"context"
	"time"

	"github.com/Southclaws/storyden/app/resources/account"
)

type Repository interface {
	Subscribe(ctx context.Context, userID account.AccountID, refersType NotificationType, refersTo string) (*Subscription, error)
	Unsubscribe(ctx context.Context, userID account.AccountID, subID SubscriptionID) (int, error)

	GetSubscriptionsForUser(ctx context.Context, userID account.AccountID) ([]Subscription, error)
	GetSubscriptionsForItem(ctx context.Context, refersType NotificationType, refersTo string) ([]Subscription, error)
	GetNotifications(ctx context.Context, userID account.AccountID, read bool, after time.Time) ([]Notification, error)

	Notify(ctx context.Context, refersType NotificationType, refersTo string, title, desc, link string) (int, error)
	SetReadState(ctx context.Context, userID account.AccountID, notificationID NotificationID, read bool) (*Notification, error)
	Delete(ctx context.Context, userID account.AccountID, notificationID NotificationID) (*Notification, error)
}
