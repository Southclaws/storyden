package thread

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

func (s *service) Delete(ctx context.Context, id post.ID) error {
	aid, err := session.GetAccountID(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := s.accountQuery.GetByID(ctx, aid)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	// TODO: Minimal reader interface for thread.
	thr, err := s.thread_repo.Get(ctx, id, pagination.Parameters{}, nil)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.authoriseThreadDelete(ctx, acc, thr); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	err = s.thread_repo.Delete(ctx, id)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to delete thread"))
	}

	if err := s.deleteQueue.Publish(ctx, mq.DeleteThread{
		ID: thr.ID,
	}); err != nil {
		s.l.Error("failed to publish index post message", zap.Error(err))
	}

	return nil
}

func (s *service) authoriseThreadDelete(ctx context.Context, acc *account.Account, thr *thread.Thread) error {
	return acc.Roles.Permissions().Authorise(ctx, func() error {
		if thr.Author.ID != acc.ID {
			return fault.Wrap(rbac.ErrPermissions,
				fctx.With(ctx),
				fmsg.WithDesc("not author", "You are not the author of the thread and do not have the Manage Posts permission."),
			)
		}
		return nil
	}, rbac.PermissionManagePosts)
}
