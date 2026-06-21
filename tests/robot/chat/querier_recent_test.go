package chat_test

import (
	"context"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

func TestGetWithMessageFiltersReturnsNewestNotOldest(t *testing.T) {
	t.Parallel()
	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		sessionRepo *robot_session.Repository,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			ownerCtx, owner := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			_ = ownerCtx

			sessionID := robot.SessionID(xid.New())
			accountID := account.AccountID(owner.ID)

			_, err := sessionRepo.Create(root, sessionID, "test session", accountID, nil)
			require.NoError(t, err)

			// Append 5 messages with distinct, sequential invocation IDs so
			// we can identify them after retrieval.
			invocations := []string{"inv-A", "inv-B", "inv-C", "inv-D", "inv-E"}
			for _, inv := range invocations {
				err := sessionRepo.AppendMessage(
					root,
					sessionID,
					inv,
					opt.New(accountID),
					opt.NewEmpty[robot.Actor](),
					map[string]any{
						"author":  "user",
						"content": map[string]any{"parts": []any{map[string]any{"text": inv}}},
					},
				)
				require.NoError(t, err)
				// Small sleep so created_at timestamps differ on fast hardware.
				time.Sleep(2 * time.Millisecond)
			}

			t.Run("num_recent_events_returns_newest_messages", func(t *testing.T) {
				a := assert.New(t)

				// Ask for the 2 most recent events.
				sess, err := sessionRepo.GetWithMessageFilters(
					root,
					sessionID,
					opt.New(accountID),
					2,
					time.Time{},
				)
				require.NoError(t, err)
				require.Len(t, sess.Messages, 2, "expected exactly 2 messages back")

				gotInvocations := make([]string, len(sess.Messages))
				for i, m := range sess.Messages {
					gotInvocations[i] = m.InvocationID
				}

				a.Equal([]string{"inv-D", "inv-E"}, gotInvocations)
			})
		}))
	}))
}
