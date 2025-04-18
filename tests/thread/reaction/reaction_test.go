package reaction_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
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

func TestReactions(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			ctx1, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			ctx2, acc2 := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			session1 := sh.WithSession(ctx1)
			session2 := sh.WithSession(ctx2)

			cat1name := "Category " + uuid.NewString()

			cat1create, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{Admin: false, Colour: "#fe4efd", Description: "category testing", Name: cat1name}, session1)
			tests.Ok(t, err, cat1create)

			t.Run("react to thread", func(t *testing.T) {
				// acc1 creates a thread
				thread1create, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{Body: "<p>this is a thread</p>", Category: cat1create.JSON200.Id, Visibility: openapi.Published, Title: "Thread testing"}, session1)
				tests.Ok(t, err, thread1create)
				threadID := thread1create.JSON200.Id

				// acc2 reacts to it
				react1create, err := cl.PostReactAddWithResponse(root, threadID, openapi.PostReactAddJSONRequestBody{Emoji: "👻"}, session2)
				tests.Ok(t, err, react1create)
				a.Equal("👻", react1create.JSON200.Emoji)
				a.Equal(acc2.ID.String(), react1create.JSON200.Author.Id)

				thread1get, err := cl.ThreadGetWithResponse(root, threadID, nil, session1)
				tests.Ok(t, err, thread1get)
				r.Len(thread1get.JSON200.Reacts, 1)
				r.Equal("👻", thread1get.JSON200.Reacts[0].Emoji)
				r.Equal(acc2.ID.String(), thread1get.JSON200.Reacts[0].Author.Id)
			})

			t.Run("delete thread react", func(t *testing.T) {
				// acc1 creates a thread
				thread1create, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{Body: "<p>this is a thread</p>", Category: cat1create.JSON200.Id, Visibility: openapi.Published, Title: "Thread testing"}, session1)
				tests.Ok(t, err, thread1create)
				threadID := thread1create.JSON200.Id

				// acc2 reacts to it
				react1create, err := cl.PostReactAddWithResponse(root, threadID, openapi.PostReactAddJSONRequestBody{Emoji: "👻"}, session2)
				tests.Ok(t, err, react1create)

				// acc2 deletes the reaction
				react1delete, err := cl.PostReactRemoveWithResponse(root, threadID, react1create.JSON200.Id, session2)
				tests.Ok(t, err, react1delete)

				thread1get, err := cl.ThreadGetWithResponse(root, threadID, nil, session1)
				tests.Ok(t, err, thread1get)
				r.Len(thread1get.JSON200.Reacts, 0)
			})

			t.Run("react to reply", func(t *testing.T) {
				// acc1 creates a thread with 1 reply
				thread1create, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{Body: "<p>this is a thread</p>", Category: cat1create.JSON200.Id, Visibility: openapi.Published, Title: "Thread testing"}, session1)
				tests.Ok(t, err, thread1create)
				reply1create, err := cl.ReplyCreateWithResponse(root, thread1create.JSON200.Id, openapi.ReplyInitialProps{Body: "<p>this is a reply</p>"}, session1)
				tests.Ok(t, err, reply1create)

				// acc2 reacts to the reply
				react1create, err := cl.PostReactAddWithResponse(root, reply1create.JSON200.Id, openapi.PostReactAddJSONRequestBody{Emoji: "👻"}, session2)
				tests.Ok(t, err, react1create)
				a.Equal("👻", react1create.JSON200.Emoji)
				a.Equal(acc2.ID.String(), react1create.JSON200.Author.Id)

				thread1get, err := cl.ThreadGetWithResponse(root, thread1create.JSON200.Id, nil, session1)
				tests.Ok(t, err, thread1get)
				r.Len(thread1get.JSON200.Reacts, 0)
				r.Len(thread1get.JSON200.Replies.Replies, 1)
				reply := thread1get.JSON200.Replies.Replies[0]
				r.Len(reply.Reacts, 1)
				r.Equal("👻", reply.Reacts[0].Emoji)
				r.Equal(acc2.ID.String(), reply.Reacts[0].Author.Id)
			})

			t.Run("delete reply react", func(t *testing.T) {
				// acc1 creates a thread
				thread1create, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{Body: "<p>this is a thread</p>", Category: cat1create.JSON200.Id, Visibility: openapi.Published, Title: "Thread testing"}, session1)
				tests.Ok(t, err, thread1create)
				reply1create, err := cl.ReplyCreateWithResponse(root, thread1create.JSON200.Id, openapi.ReplyInitialProps{Body: "<p>this is a reply</p>"}, session1)
				tests.Ok(t, err, reply1create)

				// acc2 reacts to it
				react1create, err := cl.PostReactAddWithResponse(root, reply1create.JSON200.Id, openapi.PostReactAddJSONRequestBody{Emoji: "👻"}, session2)
				tests.Ok(t, err, react1create)

				// acc2 deletes the reaction
				react1delete, err := cl.PostReactRemoveWithResponse(root, reply1create.JSON200.Id, react1create.JSON200.Id, session2)
				tests.Ok(t, err, react1delete)

				thread1get, err := cl.ThreadGetWithResponse(root, thread1create.JSON200.Id, nil, session1)
				tests.Ok(t, err, thread1get)
				r.Len(thread1get.JSON200.Reacts, 0)
				r.Len(thread1get.JSON200.Replies.Replies, 1)
				reply := thread1get.JSON200.Replies.Replies[0]
				r.Len(reply.Reacts, 0)
			})

			t.Run("idempotent_reactions", func(t *testing.T) {
				thread1create, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{Body: "<p>this is a thread</p>", Category: cat1create.JSON200.Id, Visibility: openapi.Published, Title: "Thread testing"}, session1)
				tests.Ok(t, err, thread1create)
				threadID := thread1create.JSON200.Id

				r1, err := cl.PostReactAddWithResponse(root, threadID, openapi.PostReactAddJSONRequestBody{Emoji: "👻"}, session2)
				tests.Ok(t, err, r1)
				r2, err := cl.PostReactAddWithResponse(root, threadID, openapi.PostReactAddJSONRequestBody{Emoji: "👻"}, session2)
				tests.Ok(t, err, r2)
				r3, err := cl.PostReactAddWithResponse(root, threadID, openapi.PostReactAddJSONRequestBody{Emoji: "👻"}, session2)
				tests.Ok(t, err, r3)

				r4, err := cl.PostReactAddWithResponse(root, threadID, openapi.PostReactAddJSONRequestBody{Emoji: "🥶"}, session2)
				tests.Ok(t, err, r4)

				thread1get, err := cl.ThreadGetWithResponse(root, thread1create.JSON200.Slug, nil)
				tests.Ok(t, err, thread1get)

				r.Len(thread1get.JSON200.Reacts, 2, "2 reacts because 👻 is ignored after the first react, reactions are unique by (post, account, emoji) constraint")
			})
		}))
	}))
}
