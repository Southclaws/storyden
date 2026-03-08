package reply_notify

import (
	"context"
	"errors"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/services/notification/notify"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func Build() fx.Option {
	return fx.Invoke(func(
		ctx context.Context,
		lc fx.Lifecycle,
		bus *pubsub.Bus,
		notifier *notify.Notifier,
	) {
		consumer := func(hctx context.Context) error {
			_, err := pubsub.Subscribe(ctx, bus, "reply_notify.reply_created", func(ctx context.Context, evt *rpc.EventThreadReplyCreated) error {
				errs := []error{}

				if evt.ReplyAuthorID != evt.ThreadAuthorID {
					err := notifier.Send(ctx,
						evt.ThreadAuthorID,
						opt.New(evt.ReplyAuthorID),
						notification.EventThreadReply,
						&datagraph.Ref{
							ID:   xid.ID(evt.ThreadID),
							Kind: datagraph.KindPost,
						},
					)
					errs = append(errs, err)
				}

				if rtid, ok := evt.ReplyToAuthorID.Get(); ok && rtid != evt.ReplyAuthorID {

					err := notifier.Send(ctx,
						rtid,
						opt.New(evt.ReplyAuthorID),
						notification.EventReplyToReply,
						&datagraph.Ref{
							ID:   xid.ID(evt.ReplyID),
							Kind: datagraph.KindPost,
						},
					)
					errs = append(errs, err)
				}

				return errors.Join(errs...)
			})
			return err
		}

		lc.Append(fx.StartHook(consumer))
	})
}
