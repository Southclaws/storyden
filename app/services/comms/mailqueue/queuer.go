package mailqueue

import (
	"context"
	"log/slog"
	"net/mail"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/services/comms/mailtemplate"
	"github.com/Southclaws/storyden/internal/infrastructure/mailer"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub/queue"
	"github.com/Southclaws/storyden/internal/infrastructure/rate"
)

var (
	EmailRateLimit       = 6
	EmailRateLimitPeriod = time.Hour
	EmailRateLimitReset  = time.Minute * 10
)

type Queuer struct {
	templates *mailtemplate.Builder
	limiter   rate.Limiter
	queue     pubsub.Topic[mq.Email]
	sender    mailer.Sender
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(queue.New[mq.Email]),
		fx.Provide(func(
			ctx context.Context,
			lc fx.Lifecycle,
			logger *slog.Logger,

			templates *mailtemplate.Builder,
			ratelimit *rate.LimiterFactory,
			queue pubsub.Topic[mq.Email],
			sender mailer.Sender,
		) *Queuer {
			q := &Queuer{
				templates: templates,
				limiter:   ratelimit.NewLimiter(EmailRateLimit, EmailRateLimitPeriod, EmailRateLimitReset),
				queue:     queue,
				sender:    sender,
			}

			lc.Append(fx.StartHook(func(_ context.Context) error {
				channel, err := queue.Subscribe(ctx)
				if err != nil {
					return err
				}

				go func() {
					for msg := range channel {
						if err := q.sender.Send(ctx, msg.Payload.Message); err != nil {
							logger.Error("failed to send email", slog.String("error", err.Error()))
						}

						msg.Ack()
					}
				}()

				return nil
			}))

			return q
		}),
	)
}

func (q *Queuer) Queue(ctx context.Context, address mail.Address, name string, subject string, intros []string, actions []mailtemplate.Action) error {
	err := q.limiter.Check(ctx, address.Address, 1)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	content, err := q.templates.Build(ctx, name, intros, actions)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	message, err := mailer.NewMessage(address, name, subject, *content)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return q.queue.Publish(ctx, mq.Email{
		Message: *message,
	})
}
