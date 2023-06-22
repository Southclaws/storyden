package notification

import (
	"time"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/ent"
)

type NotificationID xid.ID

type NotificationType string

type Notification struct {
	ID          NotificationID
	Title       string
	Description string
	Link        string
	Read        bool
	CreatedAt   time.Time
}

func FromModel(m *ent.Notification) *Notification {
	return &Notification{
		ID:          NotificationID(m.ID),
		Title:       m.Title,
		Description: m.Description,
		Link:        m.Link,
		Read:        m.Read,
		CreatedAt:   m.CreatedAt,
	}
}
