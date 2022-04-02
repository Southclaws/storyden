package notification

import (
	"context"
	"time"
)

type Repository interface {
	Subscribe(ctx context.Context, userID string, refersType NotificationType, refersTo string) (*Subscription, error)
	Unsubscribe(ctx context.Context, userID, subID string) (int, error)

	GetSubscriptionsForUser(ctx context.Context, userID string) ([]Subscription, error)
	GetSubscriptionsForItem(ctx context.Context, refersType NotificationType, refersTo string) ([]Subscription, error)
	GetNotifications(ctx context.Context, userID string, read bool, after time.Time) ([]Notification, error)

	Notify(ctx context.Context, refersType NotificationType, refersTo string, title, desc, link string) (int, error)
	SetReadState(ctx context.Context, userID string, notificationID string, read bool) (*Notification, error)
	Delete(ctx context.Context, userID, notificationID string) (*Notification, error)
}
