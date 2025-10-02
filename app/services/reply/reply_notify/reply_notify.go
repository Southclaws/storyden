package reply_notify

import (
	"context"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/services/notification/notify"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

func Build() fx.Option {
	return fx.Invoke(newReplyNotifier)
}

func newReplyNotifier(
	ctx context.Context,
	lc fx.Lifecycle,
	bus *pubsub.Bus,
	notifier *notify.Notifier,
) {
	lc.Append(fx.StartHook(func(hctx context.Context) error {
		_, err := pubsub.Subscribe(hctx, bus, "reply_notify.reply_created", func(ctx context.Context, evt *message.EventThreadReplyCreated) error {
			return notifier.Send(ctx,
				evt.ThreadAuthorID,
				opt.New(evt.ReplyAuthorID),
				notification.EventThreadReply,
				opt.New(datagraph.Ref{
					ID:   xid.ID(evt.ThreadID),
					Kind: datagraph.KindPost,
				}),
			)
		})
		return err
	}))
}
