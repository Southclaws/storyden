package reply

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"

	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

func (s *service) Delete(ctx context.Context, postID post.ID) error {
	aid, err := session.GetAccountID(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := s.accountQuery.GetByID(ctx, aid)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	p, err := s.post_repo.Get(ctx, postID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := acc.Roles.Permissions().Authorise(ctx, func() error {
		if p.Author.ID != aid {
			return fault.Wrap(rbac.ErrPermissions,
				fctx.With(ctx),
				fmsg.WithDesc("not owner", "You are not the owner of the post and do not have the Manage Posts permission."))
		}
		return nil
	}, rbac.PermissionManagePosts); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	err = s.post_repo.Delete(ctx, postID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	s.bus.Publish(ctx, &message.EventThreadReplyDeleted{
		ThreadID: p.RootPostID,
		ReplyID:  p.ID,
	})

	return nil
}
