package notification

// func TestThreadNotifications(t *testing.T) {
// 	t.Parallel()

// 	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
// 		lc fx.Lifecycle,
// 		root context.Context,
// 		cl *openapi.ClientWithResponses,
// 		cj *session.Jar,
// 		aw account_writer.Writer,
// 	) {
// 		lc.Append(fx.StartHook(func() {
// 			// NOTE: This test must wait for the entire app to be ready, since
// 			// the pubsub consumers use start hooks to subscribe to topics.
// 			// TODO: Figure out a better way to wait for the full app to start.
// 			r := require.New(t)
// 			a := assert.New(t)

// 			ctx1, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
// 			ctx2, acc2 := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
// 			user1session := sh.WithSession(ctx1)
// 			user2session := sh.WithSession(ctx2)

// 			cat1create, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{Admin: false, Colour: "#fe4efd", Description: "category testing", Name: "Category " + uuid.NewString()}, user1session)
// 			tests.Ok(t, err, cat1create)

// 			thread1create, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{Body: "<p>thread</p>", Category: cat1create.JSON200.Id, Visibility: openapi.Published, Title: "Thread testing"}, user1session)
// 			tests.Ok(t, err, thread1create)

// 			reply1, err := cl.ReplyCreateWithResponse(root, thread1create.JSON200.Slug, openapi.ReplyInitialProps{Body: "<p>reply</p>"}, user2session)
// 			tests.Ok(t, err, reply1)

// 			// TODO: Helper to wait for all messages to flush.
// 			time.Sleep(time.Second)

// 			notlist1, err := cl.NotificationListWithResponse(root, &openapi.NotificationListParams{}, user1session)
// 			tests.Ok(t, err, notlist1)

// 			r.Len(notlist1.JSON200.Notifications, 1)
// 			not1 := notlist1.JSON200.Notifications[0]
// 			a.Equal(not1.Event, "thread_reply")
// 			a.Equal(not1.Item.Kind, openapi.DatagraphItemKindPost)
// 			a.Equal(not1.Item.Id, thread1create.JSON200.Id)
// 			a.Equal(not1.Status, openapi.Unread)
// 			a.Equal(not1.Source.Id, acc2.ID.String())
// 		}))
// 	}))
// }
