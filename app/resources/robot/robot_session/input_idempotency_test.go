package robot_session_test

import (
	"context"
	"sync"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/internal/ent"
	ent_robot_session_input "github.com/Southclaws/storyden/internal/ent/robotsessioninput"
	ent_robot_session_message "github.com/Southclaws/storyden/internal/ent/robotsessionmessage"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestEnqueueInputIsIdempotentAcrossConcurrentDeliveries(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		repo *robot_session.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			owner, err := db.Account.Create().SetHandle("concurrent-input-" + xid.New().String()).SetName("Concurrent Input Owner").Save(ctx)
			require.NoError(t, err)
			ownerID := account.AccountID(owner.ID)
			sessionID := robot.SessionID(xid.New())
			require.NoError(t, createSession(ctx, repo, sessionID, ownerID))

			inputID := robot.InputID(xid.New())
			const deliveries = 8
			start := make(chan struct{})
			errs := make(chan error, deliveries)
			var wg sync.WaitGroup
			for range deliveries {
				wg.Add(1)
				go func() {
					defer wg.Done()
					<-start
					errs <- enqueueVisibleInput(ctx, repo, sessionID, ownerID, inputID, "one durable input")
				}()
			}
			close(start)
			wg.Wait()
			close(errs)
			for err := range errs {
				require.NoError(t, err)
			}

			inputCount, err := db.RobotSessionInput.Query().
				Where(ent_robot_session_input.IDEQ(xid.ID(inputID))).
				Count(ctx)
			require.NoError(t, err)
			assert.Equal(t, 1, inputCount)

			messageCount, err := db.RobotSessionMessage.Query().
				Where(ent_robot_session_message.IDEQ(xid.ID(inputID))).
				Count(ctx)
			require.NoError(t, err)
			assert.Equal(t, 1, messageCount)
		}))
	}))
}
