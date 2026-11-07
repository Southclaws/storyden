package notification

import "github.com/Southclaws/storyden/internal/ent/schema"

var defaultEventSettings = map[string]bool{
	EventThreadReply.String():          true,
	EventReplyToReply.String():         true,
	EventPostLike.String():             false,
	EventFollow.String():               false,
	EventProfileMention.String():       true,
	EventEventHostAdded.String():       true,
	EventMemberAttendingEvent.String(): false,
	EventMemberDeclinedEvent.String():  false,
	EventAttendeeRemoved.String():      true,
	EventReportSubmitted.String():      true,
	EventReportUpdated.String():        true,
	EventWarningIssued.String():        true,
	EventNodeVersionCreated.String():   true,
	EventNodeVersionApplied.String():   true,
	EventNodeVersionDeleted.String():   true,
	EventTrailRunAttention.String():    true,
}

func SettingsForEvent(overrides map[string]schema.NotificationEventSettings, event Event) schema.NotificationEventSettings {
	if settings, ok := overrides[event.String()]; ok {
		return settings
	}

	enabled := defaultEventSettings[event.String()]
	return schema.NotificationEventSettings{Inbox: enabled, Push: enabled}
}

func ResolveEventSettings(overrides map[string]schema.NotificationEventSettings) map[string]schema.NotificationEventSettings {
	settings := make(map[string]schema.NotificationEventSettings, len(defaultEventSettings))
	for event, enabled := range defaultEventSettings {
		settings[event] = schema.NotificationEventSettings{Inbox: enabled, Push: enabled}
	}

	for event, override := range overrides {
		settings[event] = override
	}

	return settings
}
