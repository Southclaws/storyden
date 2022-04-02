package notification

import (
	"time"

	"github.com/Southclaws/storyden/api/src/infra/db"
)

type NotificationType string

const (
	NotificationTypeForumPostResponse NotificationType = NotificationType(db.NotificationTypeFORUMPOSTRESPONSE)
)

type Notification struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Link         string        `json:"link"`
	Read         bool          `json:"read"`
	CreatedAt    time.Time     `json:"createdAt"`
	Subscription *Subscription `json:"subscription"`
}

type Subscription struct {
	ID         string           `json:"id"`
	RefersType NotificationType `json:"refersType"`
	RefersTo   string           `json:"refersTo"`
	CreatedAt  time.Time        `json:"createdAt"`
	UpdatedAt  time.Time        `json:"updatedAt"`
	DeletedAt  *time.Time       `json:"deletedAt"`
}

func SubFromModel(m *db.SubscriptionModel) *Subscription {
	return &Subscription{
		ID:         m.ID,
		RefersType: NotificationType(m.RefersType),
		RefersTo:   m.RefersTo,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		DeletedAt:  m.InnerSubscription.DeletedAt,
	}
}

func FromModel(m *db.NotificationModel) *Notification {
	var sub *Subscription
	if m.RelationsNotification.Subscription != nil {
		sub = SubFromModel(m.RelationsNotification.Subscription)
	}

	return &Notification{
		ID:           m.InnerNotification.ID,
		Title:        m.InnerNotification.Title,
		Description:  m.InnerNotification.Description,
		Link:         m.InnerNotification.Link,
		Read:         m.InnerNotification.Read,
		CreatedAt:    m.InnerNotification.CreatedAt,
		Subscription: sub,
	}
}
