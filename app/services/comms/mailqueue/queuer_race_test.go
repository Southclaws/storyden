package mailqueue_test

import (
	"context"
	"errors"
	"net/mail"
	"sync"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/email_queue"
	"github.com/Southclaws/storyden/app/resources/email_queue/email_queue_repo"
	"github.com/Southclaws/storyden/app/resources/message"
	mailqueue "github.com/Southclaws/storyden/app/services/comms/mailqueue"
	"github.com/Southclaws/storyden/internal/infrastructure/mailer"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/integration"
)

type scriptedSender struct {
	mu       sync.Mutex
	results  []error
	messages []mailer.Message
}

func (s *scriptedSender) Send(_ context.Context, msg mailer.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.messages = append(s.messages, msg)

	if len(s.results) == 0 {
		return nil
	}

	result := s.results[0]
	s.results = s.results[1:]
	return result
}

func (s *scriptedSender) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.messages)
}

func TestRetryNow_DoesNotMakeRowPollerEligibleBeforeImmediateAttempt(t *testing.T) {
	t.Parallel()

	sender := &scriptedSender{
		results: []error{
			errors.New("first retry send fails after provider may have accepted it"),
			nil,
		},
	}

	integration.Test(t, nil,
		fx.Decorate(func(mailer.Sender) mailer.Sender { return sender }),
		fx.Invoke(func(
			lc fx.Lifecycle,
			ctx context.Context,
			repo *email_queue_repo.Repository,
			bus *pubsub.Bus,
			_ *mailqueue.Queuer,
		) {
			lc.Append(fx.StartHook(func(context.Context) {
				r := require.New(t)

				id := xid.New()
				now := time.Now()

				msg := mailer.Message{
					Address: mail.Address{
						Address: "race-" + xid.New().String() + "@storyden.org",
					},
					Name:    "Race Recipient",
					Subject: "Retry Race",
					Content: mailer.Content{
						Plain: "plain body",
					},
				}

				rec, err := repo.CreateOrGetByID(ctx, id, msg)
				r.NoError(err)

				err = repo.MarkFailed(ctx, rec.ID, []*email_queue.Attempt{
					{
						Timestamp: now.Add(-time.Minute),
						Status:    email_queue.StatusFailed,
						Error:     opt.New("seed failure"),
					},
				}, now.Add(time.Hour))
				r.NoError(err)

				_, err = repo.RetryNow(ctx, rec.ID, now.Add(time.Minute))
				r.NoError(err)

				_, claimed, err := repo.ClaimNext(ctx, now)
				r.NoError(err)
				r.False(claimed, "manual retry should not make the row immediately poller-eligible")

				err = bus.SendCommand(ctx, &message.CommandAttemptQueuedEmail{ID: id})
				r.NoError(err)

				r.Eventually(func() bool {
					return sender.Count() == 1
				}, time.Second, 10*time.Millisecond, "manual retry should trigger only one provider send")

				r.Eventually(func() bool {
					final, err := repo.Get(ctx, rec.ID)
					if err != nil {
						return false
					}

					return final.Status == email_queue.StatusFailed && len(final.Attempts) == 2
				}, time.Second, 10*time.Millisecond, "seed failure plus one new failed retry attempt")
			}))
		}),
	)
}
