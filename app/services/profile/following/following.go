package following

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/profile/follow_writer"
	"github.com/Southclaws/storyden/app/services/notification/notify"
)

type FollowManager struct {
	followWriter *follow_writer.Writer
	notifier     *notify.Notifier
}

func New(followWriter *follow_writer.Writer, notifier *notify.Notifier) *FollowManager {
	return &FollowManager{followWriter: followWriter, notifier: notifier}
}

func (f *FollowManager) Follow(ctx context.Context, follower, following account.AccountID) error {
	err := f.followWriter.Follow(ctx, follower, following)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	f.notifier.Send(ctx, following, opt.New(follower), notification.EventFollow, nil)

	return nil
}

func (f *FollowManager) Unfollow(ctx context.Context, follower, following account.AccountID) error {
	err := f.followWriter.Unfollow(ctx, follower, following)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
