package mention_job

import (
	"context"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/services/notification/notify"
)

type mentionConsumer struct {
	notifySender *notify.Notifier
}

func newMentionConsumer(
	notifySender *notify.Notifier,
) *mentionConsumer {
	return &mentionConsumer{
		notifySender: notifySender,
	}
}

func (s *mentionConsumer) mention(ctx context.Context, by account.AccountID, source datagraph.Ref, item datagraph.Ref) error {
	switch item.Kind {
	case datagraph.KindProfile:
		s.notifySender.Send(ctx, account.AccountID(item.ID), opt.New(by), notification.EventProfileMention, &source)

		// TODO: Store mention in db
	}

	return nil
}
