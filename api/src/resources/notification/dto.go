package notification

import (
	"time"

	"4d63.com/optional"
	"github.com/Southclaws/storyden/api/src/infra/db/model"
	"github.com/Southclaws/storyden/api/src/utils"
	"github.com/google/uuid"
)

type (
	NotificationID uuid.UUID
	SubscriptionID uuid.UUID
)

type NotificationType string

const (
	NotificationTypeForumPostResponse NotificationType = NotificationType(model.RefersTyp)
)

type Notification struct {
	ID           NotificationID `json:"id"`
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	Link         string         `json:"link"`
	Read         bool           `json:"read"`
	CreatedAt    time.Time      `json:"createdAt"`
	Subscription *Subscription  `json:"subscription"`
}

type Subscription struct {
	ID         SubscriptionID               `json:"id"`
	RefersType NotificationType             `json:"refersType"`
	RefersTo   string                       `json:"refersTo"`
	CreatedAt  time.Time                    `json:"createdAt"`
	UpdatedAt  time.Time                    `json:"updatedAt"`
	DeletedAt  optional.Optional[time.Time] `json:"deletedAt"`
}

func SubFromModel(m *model.Subscription) *Subscription {
	return &Subscription{
		ID:         SubscriptionID(m.ID),
		RefersType: NotificationType(m.RefersType),
		RefersTo:   m.RefersTo,
		CreatedAt:  m.CreateTime,
		UpdatedAt:  m.UpdateTime,
		DeletedAt:  utils.OptionalZero(m.DeleteTime),
	}
}

func FromModel(m *model.Notification) *Notification {
	var sub *Subscription
	if m.Subscription != nil {
		sub = SubFromModel(m.Subscription)
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
