package event_job

import (
	"go.uber.org/fx"
)

func Build() fx.Option {
	return fx.Options(
		// TODO: Implement event notification consumers for:
		// - EventScheduledEventCreated, EventScheduledEventUpdated, EventScheduledEventDeleted
		// - EventParticipantJoined, EventParticipantLeft, etc.
		// This will handle notifying participants of event changes
	)
}
