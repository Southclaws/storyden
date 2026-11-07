package notification_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestNotificationSettings(t *testing.T) {
	t.Parallel()

	r := require.New(t)
	a := assert.New(t)

	var (
		root     context.Context
		client   *openapi.ClientWithResponses
		sessions *e2e.SessionHelper
		accounts *account_writer.Writer
	)

	integration.Test(t, nil,
		e2e.Setup(),
		fx.Populate(&root, &client, &sessions, &accounts),
	)

	ownerCtx, _ := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
	ownerSession := sessions.WithSession(ownerCtx)
	otherCtx, _ := e2e.WithAccount(root, accounts, seed.Account_003_Baldur)
	otherSession := sessions.WithSession(otherCtx)

	initial := tests.AssertRequest(client.AccountGetWithResponse(root, ownerSession))(t, http.StatusOK)
	r.NotNil(initial.JSON200.NotificationSettings)
	a.True(initial.JSON200.NotificationSettings.Push.Enabled)
	a.Equal(openapi.NotificationEventSettings{Inbox: true, Push: true}, initial.JSON200.NotificationSettings.Events["thread_reply"])
	a.Equal(openapi.NotificationEventSettings{}, initial.JSON200.NotificationSettings.Events["post_like"])
	a.Equal(openapi.NotificationEventSettings{}, initial.JSON200.NotificationSettings.Events["follow"])

	events := map[string]openapi.NotificationEventSettings{
		"thread_reply":         {Inbox: false, Push: false},
		"reply_to_reply":       {Inbox: true, Push: false},
		"node_version_applied": {Inbox: true, Push: true},
	}
	preferences := &openapi.NotificationSettings{
		Events: events,
		Push:   openapi.PushNotificationSettings{Enabled: true},
	}
	input := openapi.AccountUpdateJSONRequestBody{NotificationSettings: preferences}

	tests.AssertRequest(client.AccountUpdateWithResponse(root, input))(t, http.StatusUnauthorized)
	updated := tests.AssertRequest(client.AccountUpdateWithResponse(root, input, ownerSession))(t, http.StatusOK)
	for event, settings := range events {
		a.Equal(settings, updated.JSON200.NotificationSettings.Events[event])
	}
	a.Equal(initial.JSON200.NotificationSettings.Events["follow"], updated.JSON200.NotificationSettings.Events["follow"])

	loaded := tests.AssertRequest(client.AccountGetWithResponse(root, ownerSession))(t, http.StatusOK)
	a.Equal(updated.JSON200.NotificationSettings, loaded.JSON200.NotificationSettings)

	other := tests.AssertRequest(client.AccountGetWithResponse(root, otherSession))(t, http.StatusOK)
	a.True(other.JSON200.NotificationSettings.Push.Enabled)
	a.Equal(initial.JSON200.NotificationSettings, other.JSON200.NotificationSettings)

	name := "New name"
	tests.AssertRequest(client.AccountUpdateWithResponse(root,
		openapi.AccountUpdateJSONRequestBody{Name: &name}, ownerSession,
	))(t, http.StatusOK)
	loaded = tests.AssertRequest(client.AccountGetWithResponse(root, ownerSession))(t, http.StatusOK)
	a.Equal(updated.JSON200.NotificationSettings, loaded.JSON200.NotificationSettings)

	invalidEvents := map[string]openapi.NotificationEventSettings{"replies": {Inbox: false, Push: false}}
	invalid := openapi.AccountUpdateJSONRequestBody{
		NotificationSettings: &openapi.NotificationSettings{
			Events: invalidEvents,
			Push:   openapi.PushNotificationSettings{Enabled: true},
		},
	}
	tests.AssertRequest(client.AccountUpdateWithResponse(root, invalid, ownerSession))(t, http.StatusBadRequest)
	loaded = tests.AssertRequest(client.AccountGetWithResponse(root, ownerSession))(t, http.StatusOK)
	a.Equal(updated.JSON200.NotificationSettings, loaded.JSON200.NotificationSettings)

	invalid.NotificationSettings.Events = map[string]openapi.NotificationEventSettings{
		"thread_reply": {Inbox: false, Push: true},
	}
	tests.AssertRequest(client.AccountUpdateWithResponse(root, invalid, ownerSession))(t, http.StatusBadRequest)
	loaded = tests.AssertRequest(client.AccountGetWithResponse(root, ownerSession))(t, http.StatusOK)
	a.Equal(updated.JSON200.NotificationSettings, loaded.JSON200.NotificationSettings)

	preferences.Push.Enabled = false
	tests.AssertRequest(client.AccountUpdateWithResponse(root, input, ownerSession))(t, http.StatusOK)
	loaded = tests.AssertRequest(client.AccountGetWithResponse(root, ownerSession))(t, http.StatusOK)
	a.False(loaded.JSON200.NotificationSettings.Push.Enabled)
	a.Equal(updated.JSON200.NotificationSettings.Events, loaded.JSON200.NotificationSettings.Events)

	reset := openapi.AccountUpdateJSONRequestBody{
		NotificationSettings: &openapi.NotificationSettings{
			Events: map[string]openapi.NotificationEventSettings{},
			Push:   openapi.PushNotificationSettings{Enabled: true},
		},
	}
	tests.AssertRequest(client.AccountUpdateWithResponse(root, reset, ownerSession))(t, http.StatusOK)
	loaded = tests.AssertRequest(client.AccountGetWithResponse(root, ownerSession))(t, http.StatusOK)
	a.True(loaded.JSON200.NotificationSettings.Push.Enabled)
	a.Equal(initial.JSON200.NotificationSettings, loaded.JSON200.NotificationSettings)
}
