package notification_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/account/notification/notify_querier"
	"github.com/Southclaws/storyden/app/resources/account/notification/notify_writer"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestNotificationProjection(t *testing.T) {
	t.Parallel()
	r := require.New(t)
	a := assert.New(t)
	var (
		root     context.Context
		db       *ent.Client
		client   *openapi.ClientWithResponses
		sessions *e2e.SessionHelper
		accounts *account_writer.Writer
		writer   *notify_writer.Writer
		reader   *notify_querier.Querier
	)
	integration.Test(t, nil, e2e.Setup(), fx.Populate(&root, &db, &client, &sessions, &accounts, &writer, &reader))
	ownerCtx, owner := e2e.WithAccount(root, accounts, seed.Account_001_Odin)
	_, other := e2e.WithAccount(root, accounts, seed.Account_003_Baldur)
	session := sessions.WithSession(ownerCtx)
	largeContent := strings.Repeat("Long body text. ", 10000)
	post, err := db.Post.Create().
		SetUpdatedAt(time.Now()).
		SetAuthorID(xid.ID(other.ID)).
		SetTitle("A little harbour").
		SetSlug("a-little-harbour").
		SetBody("<p>" + largeContent + "</p>").
		SetShort("Long body text.").
		SetLastReplyAt(time.Now()).
		Save(root)
	r.NoError(err)
	reply, err := db.Post.Create().
		SetUpdatedAt(time.Now()).
		SetAuthorID(xid.ID(other.ID)).
		SetRootPostID(post.ID).
		SetBody("<p>A reply.</p>").
		SetShort("A reply.").
		SetLastReplyAt(time.Now()).
		Save(root)
	r.NoError(err)
	node, err := db.Node.Create().
		SetUpdatedAt(time.Now()).
		SetOwnerID(xid.ID(other.ID)).
		SetName("Field notes").
		SetSlug("field-notes-" + xid.New().String()).
		SetContent(largeContent).
		Save(root)
	r.NoError(err)

	for _, item := range []*datagraph.Ref{
		{ID: post.ID, Kind: datagraph.KindPost},
		{ID: reply.ID, Kind: datagraph.KindPost},
		{ID: node.ID, Kind: datagraph.KindNode},
		{ID: xid.ID(other.ID), Kind: datagraph.KindProfile},
	} {
		_, err := writer.Notification(root, owner.ID, notification.EventThreadReply,
			opt.New[datagraph.ItemRef](item), opt.New(other.ID), opt.NewEmpty[string]())
		r.NoError(err)
	}
	archived, err := writer.Notification(root, owner.ID, notification.EventTrailRunAttention,
		opt.NewEmpty[datagraph.ItemRef](), opt.NewEmpty[account.AccountID](), opt.New("run-one"))
	r.NoError(err)
	_, err = writer.SetRead(root, owner.ID, archived.ID, true)
	r.NoError(err)

	_, err = reader.Get(root, other.ID, archived.ID)
	r.Error(err)

	unread := tests.AssertRequest(client.NotificationListWithResponse(root, &openapi.NotificationListParams{
		Status: opt.New(openapi.NotificationStatusQuery{openapi.NotificationStatusUnread}).Ptr(),
	}, session))(t, http.StatusOK)
	r.NotNil(unread.JSON200)
	r.Len(unread.JSON200.Notifications, 4)
	a.Less(len(unread.Body), 4096)
	byItem := make(map[string]openapi.NotificationItemReference)
	for _, item := range unread.JSON200.Notifications {
		a.Equal(openapi.NotificationStatusUnread, item.Status)
		r.NotNil(item.Item)
		r.NotNil(item.Source)
		a.Equal(other.ID.String(), item.Source.Id)
		a.Equal(other.Handle, item.Source.Handle)
		byItem[item.Item.Id] = *item.Item
	}
	a.Equal("A little harbour", byItem[post.ID.String()].Title)
	a.Equal("A little harbour", byItem[reply.ID.String()].Title)
	a.Equal("Field notes", byItem[node.ID.String()].Title)
	a.Equal(node.Slug, *byItem[node.ID.String()].Slug)
	a.Equal(other.Name, byItem[other.ID.String()].Title)
	a.Equal(other.Handle, *byItem[other.ID.String()].Slug)

	var raw struct {
		Notifications []struct {
			Item   map[string]any `json:"item"`
			Source map[string]any `json:"source"`
		} `json:"notifications"`
	}
	r.NoError(json.Unmarshal(unread.Body, &raw))
	for _, item := range raw.Notifications {
		a.NotContains(item.Item, "body")
		a.NotContains(item.Item, "author")
		a.NotContains(item.Source, "roles")
	}

	read := tests.AssertRequest(client.NotificationListWithResponse(root, &openapi.NotificationListParams{
		Status: opt.New(openapi.NotificationStatusQuery{openapi.NotificationStatusRead}).Ptr(),
	}, session))(t, http.StatusOK)
	r.Len(read.JSON200.Notifications, 1)
	a.Equal(archived.ID.String(), read.JSON200.Notifications[0].Id)
	a.Equal("run-one", *read.JSON200.Notifications[0].Target)
	a.Nil(read.JSON200.Notifications[0].Item)
	a.Nil(read.JSON200.Notifications[0].Source)

	r.NoError(db.Post.DeleteOneID(reply.ID).Exec(root))
	r.NoError(db.Post.DeleteOneID(post.ID).Exec(root))
	remaining, err := reader.ListNotifications(root, owner.ID, false)
	r.NoError(err)
	a.Len(remaining, 2)
}
