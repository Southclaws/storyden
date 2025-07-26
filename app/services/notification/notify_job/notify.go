package notify_job

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/account/notification/notify_writer"
	"github.com/Southclaws/storyden/app/resources/datagraph"
)

type notifyConsumer struct {
	notifyWriter *notify_writer.Writer
}

func newNotifyConsumer(
	notifyWriter *notify_writer.Writer,
) *notifyConsumer {
	return &notifyConsumer{
		notifyWriter: notifyWriter,
	}
}

func (s *notifyConsumer) notify(ctx context.Context,
	targetID account.AccountID,
	sourceID opt.Optional[account.AccountID],
	event notification.Event,
	item *datagraph.Ref,
) error {
	itemref := opt.Map(opt.NewPtr(item), func(i datagraph.Ref) datagraph.ItemRef {
		return &i
	})

	_, err := s.notifyWriter.Notification(ctx, targetID, event, itemref, sourceID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
