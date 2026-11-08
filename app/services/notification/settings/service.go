package settings

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/notification"
)

type Service struct {
	reader *account_querier.Querier
	writer *account_writer.Writer
}

func New(reader *account_querier.Querier, writer *account_writer.Writer) *Service {
	return &Service{reader: reader, writer: writer}
}

func (s *Service) Get(ctx context.Context, id account.AccountID) (account.NotificationPreferences, error) {
	return s.reader.GetNotificationSettings(ctx, id)
}

func (s *Service) Update(ctx context.Context, id account.AccountID, settings account.NotificationPreferences) (account.NotificationPreferences, error) {
	for event, choice := range settings.Events {
		if _, err := notification.NewEvent(event); err != nil {
			return account.NotificationPreferences{}, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}

		if choice.Push && !choice.Inbox {
			return account.NotificationPreferences{}, fault.New("push notifications require inbox notifications", ftag.With(ftag.InvalidArgument))
		}
	}

	return s.writer.SetNotificationSettings(ctx, id, settings)
}
