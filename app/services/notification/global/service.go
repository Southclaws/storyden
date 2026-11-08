package global

import (
	"context"
	"errors"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/post_querier"
	"github.com/Southclaws/storyden/app/services/notification/notify"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type Service struct {
	accounts *account_querier.Querier
	posts    *post_querier.Querier
	notifier *notify.Notifier
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(New),
		fx.Invoke(func(ctx context.Context, lc fx.Lifecycle, bus *pubsub.Bus, service *Service) {
			lc.Append(fx.StartHook(func() error {
				_, err := pubsub.Subscribe(ctx, bus, "notification.global.thread", service.threadPublished)
				if err != nil {
					return err
				}

				_, err = pubsub.Subscribe(ctx, bus, "notification.global.reply", service.replyPublished)
				if err != nil {
					return err
				}

				_, err = pubsub.Subscribe(ctx, bus, "notification.global.reaction", service.postReacted)
				return err
			}))
		}),
	)
}

func New(accounts *account_querier.Querier, posts *post_querier.Querier, notifier *notify.Notifier) *Service {
	return &Service{accounts: accounts, posts: posts, notifier: notifier}
}

func (s *Service) threadPublished(ctx context.Context, event *rpc.EventThreadPublished) error {
	return s.send(ctx, notification.EventGlobalThreadPublished, opt.NewEmpty[account.AccountID](), event.ID, nil)
}

func (s *Service) replyPublished(ctx context.Context, event *rpc.EventThreadReplyCreated) error {
	personal := map[account.AccountID][]notification.Event{
		event.ThreadAuthorID: {notification.EventThreadReply},
	}
	if author, ok := event.ReplyToAuthorID.Get(); ok {
		personal[author] = append(personal[author], notification.EventReplyToReply)
	}

	return s.send(ctx, notification.EventGlobalReplyPublished, opt.New(event.ReplyAuthorID), event.ReplyID, personal)
}

func (s *Service) postReacted(ctx context.Context, event *rpc.EventPostReacted) error {
	return s.send(ctx, notification.EventGlobalPostReacted, opt.New(event.AccountID), event.PostID, nil)
}

func (s *Service) send(
	ctx context.Context,
	event notification.Event,
	source opt.Optional[account.AccountID],
	id post.ID,
	personal map[account.AccountID][]notification.Event,
) error {
	p, err := s.posts.GetPublished(ctx, id)
	if err != nil {
		if ftag.Get(err) == ftag.NotFound {
			return nil
		}
		return fault.Wrap(err, fctx.With(ctx))
	}

	actor := source.Or(p.AuthorID)
	subscribers, err := s.accounts.ListNotificationSubscribers(ctx, event.String())
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	var errs []error
	for _, subscriber := range subscribers {
		if subscriber.ID == actor {
			continue
		}

		if receivesPersonalNotification(subscriber, personal[subscriber.ID]) {
			continue
		}

		err := s.notifier.Send(ctx, subscriber.ID, opt.New(actor), event, &datagraph.Ref{
			ID:   xid.ID(id),
			Kind: datagraph.KindPost,
		})
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func receivesPersonalNotification(subscriber account.NotificationSubscriber, events []notification.Event) bool {
	for _, event := range events {
		if notification.SettingsForEvent(subscriber.Settings.Events, event).Inbox {
			return true
		}
	}

	return false
}
