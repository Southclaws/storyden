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

func TestPushNotificationPreferences(t *testing.T) {
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
	r.NotNil(initial.JSON200.NotificationSettings.Push.Events)
	a.Empty(*initial.JSON200.NotificationSettings.Push.Events)

	events := map[string]bool{
		"thread_reply":         false,
		"reply_to_reply":       true,
		"node_version_applied": false,
	}
	preferences := &openapi.NotificationSettings{
		Push: openapi.PushNotificationSettings{
			Enabled: true,
			Events:  &events,
		},
	}
	input := openapi.AccountUpdateJSONRequestBody{NotificationSettings: preferences}

	tests.AssertRequest(client.AccountUpdateWithResponse(root, input))(t, http.StatusUnauthorized)
	updated := tests.AssertRequest(client.AccountUpdateWithResponse(root, input, ownerSession))(t, http.StatusOK)
	a.Equal(preferences, updated.JSON200.NotificationSettings)

	loaded := tests.AssertRequest(client.AccountGetWithResponse(root, ownerSession))(t, http.StatusOK)
	a.Equal(preferences, loaded.JSON200.NotificationSettings)

	other := tests.AssertRequest(client.AccountGetWithResponse(root, otherSession))(t, http.StatusOK)
	a.True(other.JSON200.NotificationSettings.Push.Enabled)
	a.Empty(*other.JSON200.NotificationSettings.Push.Events)

	name := "New name"
	tests.AssertRequest(client.AccountUpdateWithResponse(root,
		openapi.AccountUpdateJSONRequestBody{Name: &name}, ownerSession,
	))(t, http.StatusOK)
	loaded = tests.AssertRequest(client.AccountGetWithResponse(root, ownerSession))(t, http.StatusOK)
	a.Equal(preferences, loaded.JSON200.NotificationSettings)

	invalidEvents := map[string]bool{"replies": false}
	invalid := openapi.AccountUpdateJSONRequestBody{
		NotificationSettings: &openapi.NotificationSettings{
			Push: openapi.PushNotificationSettings{
				Enabled: true,
				Events:  &invalidEvents,
			},
		},
	}
	tests.AssertRequest(client.AccountUpdateWithResponse(root, invalid, ownerSession))(t, http.StatusBadRequest)
	loaded = tests.AssertRequest(client.AccountGetWithResponse(root, ownerSession))(t, http.StatusOK)
	a.Equal(preferences, loaded.JSON200.NotificationSettings)

	preferences.Push.Enabled = false
	tests.AssertRequest(client.AccountUpdateWithResponse(root, input, ownerSession))(t, http.StatusOK)
	loaded = tests.AssertRequest(client.AccountGetWithResponse(root, ownerSession))(t, http.StatusOK)
	a.Equal(preferences, loaded.JSON200.NotificationSettings)

	reset := openapi.AccountUpdateJSONRequestBody{
		NotificationSettings: &openapi.NotificationSettings{
			Push: openapi.PushNotificationSettings{Enabled: true},
		},
	}
	tests.AssertRequest(client.AccountUpdateWithResponse(root, reset, ownerSession))(t, http.StatusOK)
	loaded = tests.AssertRequest(client.AccountGetWithResponse(root, ownerSession))(t, http.StatusOK)
	a.True(loaded.JSON200.NotificationSettings.Push.Enabled)
	a.Empty(*loaded.JSON200.NotificationSettings.Push.Events)
}
