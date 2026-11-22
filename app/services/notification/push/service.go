package push

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/account/notification/notify_querier"
	"github.com/Southclaws/storyden/app/resources/account/notification/push_subscription"
	"github.com/Southclaws/storyden/internal/infrastructure/webpush"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type SendFunc func(ctx context.Context, owner account.AccountID, subscription push_subscription.Subscription, item *notification.Notification) error

type Service struct {
	repository    *push_subscription.Repository
	accounts      *account_querier.Querier
	notifications *notify_querier.Querier
	client        *webpush.Client
	logger        *slog.Logger
}

func Build() fx.Option {
	return fx.Provide(New)
}

func New(
	logger *slog.Logger,
	repository *push_subscription.Repository,
	accounts *account_querier.Querier,
	notifications *notify_querier.Querier,
	client *webpush.Client,
) *Service {
	return &Service{
		repository:    repository,
		accounts:      accounts,
		notifications: notifications,
		client:        client,
		logger:        logger,
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

func (s *Service) Deliver(ctx context.Context, event *rpc.EventNotificationCreated, send SendFunc) error {
	settings, err := s.accounts.GetNotificationSettings(ctx, event.AccountID)
	if err != nil {
		s.logger.Error("failed to load push preferences",
			slog.String("notification_id", event.ID.String()),
			slog.String("error", err.Error()),
		)
		return nil
	}

	if !settings.Push.Enabled {
		return nil
	}

	item, err := s.notifications.Get(ctx, event.AccountID, event.ID)
	if err != nil {
		s.logger.Error("failed to load push notification",
			slog.String("notification_id", event.ID.String()),
			slog.String("error", err.Error()),
		)
		return nil
	}

	if !notification.PushEnabled(settings.Events, item.Event) {
		return nil
	}

	subscriptions, err := s.repository.List(ctx, event.AccountID)
	if err != nil {
		s.logger.Error("failed to load push subscriptions",
			slog.String("notification_id", event.ID.String()),
			slog.String("error", err.Error()),
		)
		return nil
	}

	for _, subscription := range subscriptions {
		err := send(ctx, event.AccountID, subscription, item)
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
