package notification

import "github.com/Southclaws/storyden/app/resources/account"

var defaultEventSettings = map[string]bool{
	EventGlobalThreadPublished.String(): false,
	EventGlobalReplyPublished.String():  false,
	EventGlobalPostReacted.String():     false,
	EventThreadReply.String():           true,
	EventReplyToReply.String():          true,
	EventPostLike.String():              false,
	EventFollow.String():                false,
	EventProfileMention.String():        true,
	EventEventHostAdded.String():        true,
	EventMemberAttendingEvent.String():  false,
	EventMemberDeclinedEvent.String():   false,
	EventAttendeeRemoved.String():       true,
	EventReportSubmitted.String():       true,
	EventReportUpdated.String():         true,
	EventWarningIssued.String():         true,
	EventNodeVersionCreated.String():    true,
	EventNodeVersionApplied.String():    true,
	EventNodeVersionDeleted.String():    true,
	EventTrailRunAttention.String():     true,
}

func SettingsForEvent(overrides map[string]account.NotificationEventSettings, event Event) account.NotificationEventSettings {
	if settings, ok := overrides[event.String()]; ok {
		settings.Push = settings.Inbox && settings.Push
		return settings
	}

	enabled := defaultEventSettings[event.String()]
	return account.NotificationEventSettings{Inbox: enabled, Push: enabled}
}

func PushEnabled(overrides map[string]account.NotificationEventSettings, event Event) bool {
	settings := SettingsForEvent(overrides, event)
	if settings.Inbox && settings.Push {
		return true
	}

	if event == EventThreadReply || event == EventReplyToReply {
		global := SettingsForEvent(overrides, EventGlobalReplyPublished)
		return global.Inbox && global.Push
	}

	return false
}

func ResolveEventSettings(overrides map[string]account.NotificationEventSettings) map[string]account.NotificationEventSettings {
	settings := make(map[string]account.NotificationEventSettings, len(defaultEventSettings))
	for event, enabled := range defaultEventSettings {
		settings[event] = account.NotificationEventSettings{Inbox: enabled, Push: enabled}
	}

	for event, override := range overrides {
		override.Push = override.Inbox && override.Push
		settings[event] = override
	}

	return settings
}
