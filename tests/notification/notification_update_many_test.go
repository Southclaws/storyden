package notification_test

import (
	"context"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/account/notification/notify_writer"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestNotificationUpdateMany(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		nw *notify_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			userCtx, userAcc := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			userSession := sh.WithSession(userCtx)

			not1, err := nw.Notification(root, userAcc.ID, notification.EventThreadReply, opt.NewEmpty[datagraph.ItemRef](), opt.New(userAcc.ID))
			r.NoError(err)
			r.Equal(false, not1.Read)

			not2, err := nw.Notification(root, userAcc.ID, notification.EventPostLike, opt.NewEmpty[datagraph.ItemRef](), opt.New(userAcc.ID))
			r.NoError(err)
			r.Equal(false, not2.Read)

			not3, err := nw.Notification(root, userAcc.ID, notification.EventFollow, opt.NewEmpty[datagraph.ItemRef](), opt.New(userAcc.ID))
			r.NoError(err)
			r.Equal(false, not3.Read)

			notlist, err := cl.NotificationListWithResponse(root, &openapi.NotificationListParams{}, userSession)
			tests.Ok(t, err, notlist)
			r.Len(notlist.JSON200.Notifications, 3)

			for _, n := range notlist.JSON200.Notifications {
				a.Equal(openapi.Unread, n.Status)
			}

			statusRead := openapi.Read
			updateResp, err := cl.NotificationUpdateManyWithResponse(root, openapi.NotificationListUpdate{
				Notifications: []openapi.NotificationMutation{
					{
						Id:     not1.ID.String(),
						Status: &statusRead,
					},
					{
						Id:     not2.ID.String(),
						Status: &statusRead,
					},
				},
			}, userSession)
			tests.Ok(t, err, updateResp)

			r.Len(updateResp.JSON200.Notifications, 2)
			for _, n := range updateResp.JSON200.Notifications {
				a.Equal(openapi.Read, n.Status)
			}

			notlistAfter, err := cl.NotificationListWithResponse(root, &openapi.NotificationListParams{}, userSession)
			tests.Ok(t, err, notlistAfter)
			r.Len(notlistAfter.JSON200.Notifications, 3)

			readCount := 0
			unreadCount := 0
			for _, n := range notlistAfter.JSON200.Notifications {
				if n.Status == openapi.Read {
					readCount++
				} else {
					unreadCount++
				}
			}

			a.Equal(2, readCount)
			a.Equal(1, unreadCount)
		}))
	}))
}
