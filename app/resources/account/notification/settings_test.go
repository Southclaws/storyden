package notification_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/internal/ent/schema"
)

func TestNotificationEventOverrides(t *testing.T) {
	t.Parallel()

	a := assert.New(t)
	overrides := map[string]schema.NotificationEventSettings{
		"thread_reply": {Inbox: true, Push: false},
		"post_like":    {Inbox: true, Push: true},
	}

	resolved := notification.ResolveEventSettings(overrides)
	a.Equal(schema.NotificationEventSettings{Inbox: true, Push: false}, resolved["thread_reply"])
	a.Equal(schema.NotificationEventSettings{Inbox: true, Push: true}, resolved["post_like"])
	a.Equal(schema.NotificationEventSettings{Inbox: true, Push: true}, resolved["reply_to_reply"])

	resolved["thread_reply"] = schema.NotificationEventSettings{}
	a.True(overrides["thread_reply"].Inbox)
	a.True(notification.ResolveEventSettings(nil)["thread_reply"].Inbox)
}
