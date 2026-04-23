package mailqueue

import (
	"context"
	"log/slog"
	"net/mail"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/email_queue"
	"github.com/Southclaws/storyden/app/resources/email_queue/email_queue_repo"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/services/comms/mailtemplate"
	"github.com/Southclaws/storyden/internal/infrastructure/mailer"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/infrastructure/rate"
)

const (
	// EmailRateLimit is the number of emails a single recipient can trigger
	// within the rolling EmailRateLimitPeriod window.
	EmailRateLimit = 6
	// EmailRateLimitPeriod is the rolling window used by the per-recipient
	// rate limiter on email-triggering actions.
	EmailRateLimitPeriod = time.Hour
	// EmailRateLimitReset is the bucket reset interval used by the limiter.
	EmailRateLimitReset = 10 * time.Minute

	// pollInterval controls how often the background worker wakes up to drain
	// queued retry work from the durable email queue.
	pollInterval = time.Minute
	// retryDelay is the delay before an automatically retryable email is made
	// available to the poller again after a failed send attempt.
	retryDelay = time.Minute
	// maxRetries caps automatic retries after the initial eager send attempt.
	maxRetries = 5
	// terminalRetryDelay pushes exhausted failures far enough into the future
	// that they stop being auto-retried and require manual intervention.
	terminalRetryDelay = 100 * 365 * 24 * time.Hour
)

type Queuer struct {
	templates *mailtemplate.Builder
	limiter   rate.Limiter
	queue     *email_queue_repo.Repository
	bus       *pubsub.Bus
	sender    mailer.Sender
	logger    *slog.Logger
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(func(
			ctx context.Context,
			lc fx.Lifecycle,
			logger *slog.Logger,
			templates *mailtemplate.Builder,
			ratelimit *rate.LimiterFactory,
			queue *email_queue_repo.Repository,
			bus *pubsub.Bus,
			sender mailer.Sender,
		) *Queuer {
			q := &Queuer{
				templates: templates,
				limiter:   ratelimit.NewLimiter(EmailRateLimit, EmailRateLimitPeriod, EmailRateLimitReset),
				queue:     queue,
				bus:       bus,
				sender:    sender,
				logger:    logger,
			}

			if sender != nil {
				lc.Append(fx.StartHook(func(context.Context) error {
					if _, err := pubsub.SubscribeCommand(ctx, bus, "mailqueue.send_email", q.handleCommand); err != nil {
						return err
					}
					if _, err := pubsub.SubscribeCommand(ctx, bus, "mailqueue.attempt_email", q.handleAttemptCommand); err != nil {
						return err
					}
					go q.processLoop(ctx)
					return nil
				}))
			}

			return q
		}),
	)
}

func (q *Queuer) Queue(ctx context.Context, address mail.Address, name string, subject string, intros []string, actions []mailtemplate.Action) error {
	if q.sender == nil {
		return fault.New("email sending is not enabled")
	}

	err := q.limiter.Check(ctx, address.Address, 1)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	content, err := q.templates.Build(ctx, name, intros, actions)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	msg, err := mailer.NewMessage(address, name, subject, *content)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := q.bus.SendCommand(ctx, &message.CommandSendEmail{
		ID:      xid.New(),
		Message: *msg,
	}); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (q *Queuer) handleCommand(ctx context.Context, cmd *message.CommandSendEmail) error {
	rec, err := q.queue.CreateOrGetByID(ctx, cmd.ID, cmd.Message)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return q.attemptExisting(ctx, rec, email_queue.StatusPending)
}

func (q *Queuer) RetryNow(ctx context.Context, id email_queue.ID) (*email_queue.Email, error) {
	if q.sender == nil {
		return nil, fault.New("email sending is not enabled")
	}

	rec, err := q.queue.RetryNow(ctx, id, time.Now())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := q.bus.SendCommand(ctx, &message.CommandAttemptQueuedEmail{ID: xid.ID(id)}); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return rec, nil
}

func (q *Queuer) handleAttemptCommand(ctx context.Context, cmd *message.CommandAttemptQueuedEmail) error {
	rec, err := q.queue.Get(ctx, email_queue.ID(cmd.ID))
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return q.attemptExisting(ctx, rec, email_queue.StatusPending, email_queue.StatusFailed)
}

func (q *Queuer) attemptExisting(ctx context.Context, rec *email_queue.Email, allowedStatuses ...email_queue.Status) error {
	isAllowed := func(status email_queue.Status) bool {
		for _, allowed := range allowedStatuses {
			if status == allowed {
				return true
			}
		}
		return false
	}

	switch rec.Status {
	case email_queue.StatusSent:
		return nil

	case email_queue.StatusFailed:
		if !isAllowed(email_queue.StatusFailed) {
			return nil
		}

	case email_queue.StatusPending:
		if !isAllowed(email_queue.StatusPending) {
			return fault.New("queued email is not eligible for this attempt")
		}

	case email_queue.StatusProcessing:
		return fault.New("queued email is already processing")

	default:
		return fault.New("unknown queued email status")
	}

	claimed, err := q.queue.Claim(ctx, rec.ID, allowedStatuses...)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	if !claimed {
		rec, err = q.queue.Get(ctx, rec.ID)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		switch rec.Status {
		case email_queue.StatusSent:
			return nil
		case email_queue.StatusFailed:
			if isAllowed(email_queue.StatusFailed) {
				return nil
			}
			return nil
		case email_queue.StatusPending:
			if isAllowed(email_queue.StatusPending) {
				return nil
			}
			return fault.New("queued email is not eligible for this attempt")
		case email_queue.StatusProcessing:
			return fault.New("queued email is already processing")
		default:
			return fault.New("failed to claim queued email")
		}
	}

	return q.deliverClaimed(ctx, rec, time.Now())
}

func (q *Queuer) processLoop(ctx context.Context) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			processedCount := 0
			for {
				processed, err := q.processNext(ctx)
				if err != nil {
					q.logger.Error("failed to process email queue", slog.String("error", err.Error()))
					break
				}
				if !processed {
					break
				}
				processedCount++
			}
			if processedCount > 0 {
				q.logger.Debug("processed queued email retry batch", slog.Int("emails", processedCount))
			}
		}
	}
}

func (q *Queuer) processNext(ctx context.Context) (bool, error) {
	now := time.Now()

	rec, claimed, err := q.queue.ClaimNext(ctx, now)
	if err != nil {
		return false, fault.Wrap(err, fctx.With(ctx))
	}
	if !claimed {
		return false, nil
	}

	if err := q.deliverClaimed(ctx, rec, now); err != nil {
		return true, fault.Wrap(err, fctx.With(ctx))
	}

	return true, nil
}

func (q *Queuer) deliverClaimed(ctx context.Context, rec *email_queue.Email, now time.Time) error {
	attempt := &email_queue.Attempt{
		Timestamp: now,
	}

	msg, err := rec.Message()
	if err != nil {
		attempt.Status = email_queue.StatusFailed
		attempt.Error = opt.New(err.Error())
		return q.markFailedAttempt(ctx, rec, attempt, now)
	}

	if sendErr := q.sender.Send(ctx, *msg); sendErr != nil {
		attempt.Status = email_queue.StatusFailed
		attempt.Error = opt.New(sendErr.Error())
		return q.markFailedAttempt(ctx, rec, attempt, now)
	}

	attempt.Status = email_queue.StatusSent
	processedAt := time.Now()

	if err := q.queue.MarkSent(ctx, rec.ID, append(rec.Attempts, attempt), processedAt); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (q *Queuer) markFailedAttempt(ctx context.Context, rec *email_queue.Email, attempt *email_queue.Attempt, now time.Time) error {
	attempts := append(rec.Attempts, attempt)
	availableAt := now.Add(retryDelay)

	if len(rec.Attempts) >= maxRetries {
		availableAt = now.Add(terminalRetryDelay)
		q.logger.Warn(
			"email queue retries exhausted",
			slog.String("email_id", rec.ID.String()),
			slog.String("recipient_address", rec.RecipientAddress),
			slog.Int("attempts", len(attempts)),
			slog.Int("max_retries", maxRetries),
		)
	}

	if err := q.queue.MarkFailed(ctx, rec.ID, attempts, availableAt); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
