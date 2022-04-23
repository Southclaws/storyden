package notification

import (
	"context"
	"time"

	"github.com/Southclaws/storyden/api/src/resources/user"
)

type Repository interface {
	Subscribe(ctx context.Context, userID user.UserID, refersType NotificationType, refersTo string) (*Subscription, error)
	Unsubscribe(ctx context.Context, userID user.UserID, subID SubscriptionID) (int, error)

	GetSubscriptionsForUser(ctx context.Context, userID user.UserID) ([]Subscription, error)
	GetSubscriptionsForItem(ctx context.Context, refersType NotificationType, refersTo string) ([]Subscription, error)
	GetNotifications(ctx context.Context, userID user.UserID, read bool, after time.Time) ([]Notification, error)

	Notify(ctx context.Context, refersType NotificationType, refersTo string, title, desc, link string) (int, error)
	SetReadState(ctx context.Context, userID user.UserID, notificationID NotificationID, read bool) (*Notification, error)
	Delete(ctx context.Context, userID user.UserID, notificationID NotificationID) (*Notification, error)
}
