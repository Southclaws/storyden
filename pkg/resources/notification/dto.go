package notification

import (
	"time"

	"4d63.com/optional"
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/internal/infrastructure/db/model"
	"github.com/Southclaws/storyden/internal/utils"
)

type (
	NotificationID uuid.UUID
	SubscriptionID uuid.UUID
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
	return &Notification{
		ID:           NotificationID(m.ID),
		Title:        m.Title,
		Description:  m.Description,
		Link:         m.Link,
		Read:         m.Read,
		CreatedAt:    m.CreateTime,
		Subscription: utils.Deref(SubFromModel(m.Edges.Subscription), 0),
	}
}

func FromModelMany(m []*model.Notification) []Notification {
	return lo.Map(m, func(t *model.Notification, i int) Notification { return utils.Deref(FromModel(t), 0) })
}

func SubFromModelMany(m []*model.Subscription) []Subscription {
	return lo.Map(m, func(t *model.Subscription, i int) Subscription { return utils.Deref(SubFromModel(t), 0) })
}
