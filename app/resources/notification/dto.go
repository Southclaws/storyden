package notification

import (
	"time"

	"github.com/Southclaws/dt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/infrastructure/db/model"
	"github.com/Southclaws/storyden/internal/utils"
)

type (
	NotificationID xid.ID
	SubscriptionID xid.ID
)

type NotificationType string

type Notification struct {
	ID           NotificationID `json:"id"`
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	Link         string         `json:"link"`
	Read         bool           `json:"read"`
	CreatedAt    time.Time      `json:"createdAt"`
	Subscription Subscription   `json:"subscription"`
}

type Subscription struct {
	ID         SubscriptionID   `json:"id"`
	RefersType NotificationType `json:"refersType"`
	RefersTo   string           `json:"refersTo"`
	CreatedAt  time.Time        `json:"createdAt"`
}

func SubFromModel(m *model.Subscription) *Subscription {
	return &Subscription{
		ID:         SubscriptionID(m.ID),
		RefersType: NotificationType(m.RefersType),
		RefersTo:   m.RefersTo,
		CreatedAt:  m.CreatedAt,
	}
}

func FromModel(m *model.Notification) *Notification {
	return &Notification{
		ID:           NotificationID(m.ID),
		Title:        m.Title,
		Description:  m.Description,
		Link:         m.Link,
		Read:         m.Read,
		CreatedAt:    m.CreatedAt,
		Subscription: utils.Deref(SubFromModel(m.Edges.Subscription)),
	}
}

func FromModelMany(m []*model.Notification) []Notification {
	return dt.Map(m, func(t *model.Notification) Notification {
		return utils.Deref(FromModel(t))
	})
}

func SubFromModelMany(m []*model.Subscription) []Subscription {
	return dt.Map(m, func(t *model.Subscription) Subscription {
		return utils.Deref(SubFromModel(t))
	})
}
