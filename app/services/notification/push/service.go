package push

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification/push_subscription"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/infrastructure/webpush"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type Service struct {
	repository *push_subscription.Repository
	client     *webpush.Client
	logger     *slog.Logger
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(New),
		fx.Invoke(func(
			ctx context.Context,
			lc fx.Lifecycle,
			bus *pubsub.Bus,
			service *Service,
		) {
			if !service.client.Enabled() {
				return
			}

			lc.Append(fx.StartHook(func() error {
				_, err := pubsub.Subscribe(ctx, bus, "notification.push", service.deliver)
				return err
			}))
		}),
	)
}

func New(
	logger *slog.Logger,
	repository *push_subscription.Repository,
	client *webpush.Client,
) *Service {
	return &Service{
		repository: repository,
		client:     client,
		logger:     logger,
	}
}

func (s *Service) Register(
	ctx context.Context,
	owner account.AccountID,
	input webpush.SubscriptionInput,
) (*push_subscription.Subscription, error) {
	if !s.client.Enabled() {
		return nil, fault.Wrap(webpush.ErrDisabled, ftag.With(ftag.PermissionDenied))
	}

	return s.repository.Upsert(ctx, owner, input)
}

func (s *Service) deliver(ctx context.Context, event *rpc.EventNotificationCreated) error {
	subscriptions, err := s.repository.List(ctx, event.AccountID)
	if err != nil {
		s.logger.Error("failed to load push subscriptions",
			slog.String("notification_id", event.ID.String()),
			slog.String("error", err.Error()),
		)
		return nil
	}

	payload := []byte(event.ID.String())

	for _, subscription := range subscriptions {
		err := s.client.Send(ctx, subscription.Subscription, payload)
		if errors.Is(err, webpush.ErrExpired) {
			err = s.repository.Delete(ctx, event.AccountID, subscription.ID)
		}

		if err != nil {
			s.logger.Warn("push delivery failed",
				slog.String("notification_id", event.ID.String()),
				slog.String("subscription_id", subscription.ID.String()),
				slog.String("error", err.Error()),
			)
		}
	}
	return nil
}
