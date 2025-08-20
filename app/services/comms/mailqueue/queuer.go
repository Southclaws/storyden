package mailqueue

import (
	"context"
	"log/slog"
	"net/mail"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/services/comms/mailtemplate"
	"github.com/Southclaws/storyden/internal/infrastructure/mailer"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
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
	bus       *pubsub.Bus
	sender    mailer.Sender
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(func(
			ctx context.Context,
			lc fx.Lifecycle,
			logger *slog.Logger,

			templates *mailtemplate.Builder,
			ratelimit *rate.LimiterFactory,
			bus *pubsub.Bus,
			sender mailer.Sender,
		) *Queuer {
			q := &Queuer{
				templates: templates,
				limiter:   ratelimit.NewLimiter(EmailRateLimit, EmailRateLimitPeriod, EmailRateLimitReset),
				bus:       bus,
				sender:    sender,
			}

			lc.Append(fx.StartHook(func(hctx context.Context) error {
				_, err := pubsub.SubscribeCommand(hctx, bus, "mailqueue.send_email", func(ctx context.Context, cmd *message.CommandSendEmail) error {
					if err := sender.Send(ctx, cmd.Message); err != nil {
						logger.Error("failed to send email", slog.String("error", err.Error()))
						return err
					}
					return nil
				})

				return err
			}))

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
		Message: *msg,
	}); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
