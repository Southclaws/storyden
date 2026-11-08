package account

import "github.com/Southclaws/storyden/internal/ent"

type NotificationPreferences struct {
	Events map[string]NotificationEventSettings
	Push   PushNotificationPreferences
}

type NotificationEventSettings struct {
	Inbox bool
	Push  bool
}

type PushNotificationPreferences struct {
	Enabled bool
}

type NotificationSubscriber struct {
	ID       AccountID
	Settings NotificationPreferences
}

func MapNotificationSettings(pushEnabled bool, rows []*ent.NotificationEventSettings) NotificationPreferences {
	settings := NotificationPreferences{
		Events: make(map[string]NotificationEventSettings, len(rows)),
		Push:   PushNotificationPreferences{Enabled: pushEnabled},
	}
	for _, row := range rows {
		settings.Events[row.Event] = NotificationEventSettings{
			Inbox: row.InboxEnabled,
			Push:  row.InboxEnabled && row.PushEnabled,
		}
	}

	return settings
}
