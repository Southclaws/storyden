package notification_test

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
	"github.com/Southclaws/storyden/tests"
)

func TestGlobalNotifications(t *testing.T) {
	if tests.IsSharedPostgresDatabase() {
		t.Skip("global notification subscriptions affect parallel integration apps")
	}

	t.Parallel()

	r := require.New(t)
	a := assert.New(t)
	var (
		root     context.Context
		client   *openapi.ClientWithResponses
		sessions *e2e.SessionHelper
		accounts *account_writer.Writer
		bus      *pubsub.Bus
	)

	integration.Test(t, nil,
		e2e.Setup(),
		fx.Populate(&root, &client, &sessions, &accounts, &bus),
	)

	actorCtx, actor := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
	watcherCtx, watcher := e2e.WithAccount(root, accounts, seed.Account_003_Baldur)
	quietCtx, _ := e2e.WithAccount(root, accounts, seed.Account_004_Loki)
	actorSession := sessions.WithSession(actorCtx)
	watcherSession := sessions.WithSession(watcherCtx)
	quietSession := sessions.WithSession(quietCtx)
	choices := openapi.NotificationSettingsUpdateJSONRequestBody{
		Events: map[string]openapi.NotificationEventSettings{
			"global_thread_published": {Inbox: true, Push: true},
			"global_reply_published":  {Inbox: true, Push: true},
			"global_post_reacted":     {Inbox: true, Push: true},
		},
		Push: openapi.PushNotificationSettings{Enabled: true},
	}
	for _, session := range []openapi.RequestEditorFn{actorSession, watcherSession} {
		tests.AssertRequest(client.NotificationSettingsUpdateWithResponse(root, choices, session))(t, http.StatusOK)
	}

	var mutex sync.Mutex
	created := make(map[xid.ID]bool)
	_, err := pubsub.Subscribe(root, bus, "test.global.notifications", func(_ context.Context, event *rpc.EventNotificationCreated) error {
		if event.AccountID == watcher.ID {
			mutex.Lock()
			created[event.ID] = true
			mutex.Unlock()
		}

		return nil
	})
	r.NoError(err)

	verifyNotification := func(event, postID string, expectedEvents int) {
		t.Helper()

		r.Eventually(func() bool {
			mutex.Lock()
			defer mutex.Unlock()

			return len(created) == expectedEvents
		}, 5*time.Second, 10*time.Millisecond)

		inbox := tests.AssertRequest(client.NotificationListWithResponse(root, nil, watcherSession))(t, http.StatusOK)
		r.NotNil(inbox.JSON200)
		r.Len(inbox.JSON200.Notifications, expectedEvents)
		found := false
		for _, item := range inbox.JSON200.Notifications {
			if string(item.Event) != event {
				continue
			}

			found = true
			r.NotNil(item.Item)
			r.NotNil(item.Source)
			post, err := item.Item.AsDatagraphItemPost()
			r.NoError(err)
			a.Equal(postID, post.Ref.Id)
			a.Equal(actor.ID.String(), item.Source.Id)
			a.Equal(openapi.Unread, item.Status)

			id, err := xid.FromString(item.Id)
			r.NoError(err)
			mutex.Lock()
			published := created[id]
			mutex.Unlock()
			a.True(published)
		}
		a.True(found)
	}

	thread := tests.AssertRequest(client.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
		Title:      "Global notification test",
		Body:       opt.New("<p>Thread notification test.</p>").Ptr(),
		Visibility: opt.New(openapi.VisibilityPublished).Ptr(),
	}, actorSession))(t, http.StatusOK)
	verifyNotification("global_thread_published", thread.JSON200.Id, 1)

	reply := tests.AssertRequest(client.ReplyCreateWithResponse(root, thread.JSON200.Slug, openapi.ReplyInitialProps{
		Body: "<p>Reply notification test.</p>",
	}, actorSession))(t, http.StatusOK)
	verifyNotification("global_reply_published", reply.JSON200.Id, 2)

	react := openapi.PostReactAddJSONRequestBody{Emoji: "👍"}
	tests.AssertRequest(client.PostReactAddWithResponse(root, thread.JSON200.Id, react, actorSession))(t, http.StatusOK)
	verifyNotification("global_post_reacted", thread.JSON200.Id, 3)

	expected := []openapi.NotificationEvent{
		"global_thread_published",
		"global_reply_published",
		"global_post_reacted",
	}

	for _, session := range []openapi.RequestEditorFn{actorSession, quietSession} {
		inbox := tests.AssertRequest(client.NotificationListWithResponse(root, nil, session))(t, http.StatusOK)
		r.NotNil(inbox.JSON200)
		for _, item := range inbox.JSON200.Notifications {
			a.NotContains(expected, item.Event)
		}
	}

	tests.AssertRequest(client.PostReactAddWithResponse(root, thread.JSON200.Id, react, actorSession))(t, http.StatusOK)
	for _, visibility := range []openapi.Visibility{openapi.VisibilityDraft, openapi.VisibilityReview, openapi.VisibilityUnlisted} {
		private := tests.AssertRequest(client.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
			Title:      "Private thread",
			Body:       opt.New("<p>Private content.</p>").Ptr(),
			Visibility: opt.New(visibility).Ptr(),
		}, actorSession))(t, http.StatusOK)
		tests.AssertRequest(client.ReplyCreateWithResponse(root, private.JSON200.Slug, openapi.ReplyInitialProps{
			Body: "<p>Private reply.</p>",
		}, actorSession))(t, http.StatusOK)
		tests.AssertRequest(client.PostReactAddWithResponse(root, private.JSON200.Id, react, actorSession))(t, http.StatusOK)
	}

	a.Never(func() bool {
		mutex.Lock()
		published := len(created)
		mutex.Unlock()
		inbox, err := client.NotificationListWithResponse(root, nil, watcherSession)
		if err != nil || inbox.JSON200 == nil {
			return true
		}

		return len(inbox.JSON200.Notifications) != 3 || published != 3
	}, 300*time.Millisecond, 10*time.Millisecond)
}
