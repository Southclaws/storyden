package notification_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
)

func TestGlobalRepliesIncludePushForPersonalReplies(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	overrides := map[string]account.NotificationEventSettings{
		"thread_reply":           {Inbox: true, Push: false},
		"reply_to_reply":         {Inbox: true, Push: false},
		"global_reply_published": {Inbox: true, Push: true},
	}

	a.True(notification.PushEnabled(overrides, notification.EventThreadReply))
	a.True(notification.PushEnabled(overrides, notification.EventReplyToReply))
	a.False(notification.PushEnabled(overrides, notification.EventPostLike))

	overrides["global_reply_published"] = account.NotificationEventSettings{Inbox: true, Push: false}
	a.False(notification.PushEnabled(overrides, notification.EventThreadReply))
}

func TestNotificationEventOverrides(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	overrides := map[string]account.NotificationEventSettings{
		"thread_reply": {Inbox: true, Push: false},
		"post_like":    {Inbox: true, Push: true},
	}

	resolved := notification.ResolveEventSettings(overrides)
	a.Equal(account.NotificationEventSettings{Inbox: true, Push: false}, resolved["thread_reply"])
	a.Equal(account.NotificationEventSettings{Inbox: true, Push: true}, resolved["post_like"])
	a.Equal(account.NotificationEventSettings{Inbox: true, Push: true}, resolved["reply_to_reply"])

	resolved["thread_reply"] = account.NotificationEventSettings{}
	a.True(overrides["thread_reply"].Inbox)
	a.True(notification.ResolveEventSettings(nil)["thread_reply"].Inbox)
}

func TestPushRequiresInboxForStoredOverrides(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	overrides := map[string]account.NotificationEventSettings{
		"thread_reply":           {Inbox: false, Push: true},
		"global_reply_published": {Inbox: false, Push: true},
	}

	a.Equal(account.NotificationEventSettings{}, notification.SettingsForEvent(overrides, notification.EventThreadReply))
	a.Equal(account.NotificationEventSettings{}, notification.ResolveEventSettings(overrides)["thread_reply"])
	a.False(notification.PushEnabled(overrides, notification.EventThreadReply))
	a.True(overrides["thread_reply"].Push)
}
