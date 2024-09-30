package event_job

import (
	"github.com/Southclaws/storyden/app/services/notification/notify"
)

type eventUpdateConsumer struct {
	notifySender *notify.Notifier
}

func newEventUpdateConsumer(
	notifySender *notify.Notifier,
) *eventUpdateConsumer {
	return &eventUpdateConsumer{
		notifySender: notifySender,
	}
}
