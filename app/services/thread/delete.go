package thread

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
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
	thr, err := s.threadQuerier.Get(ctx, id, pagination.Parameters{}, nil)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.authoriseThreadDelete(ctx, acc, thr); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.cache.Invalidate(ctx, xid.ID(id)); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	err = s.threadWriter.Delete(ctx, id)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to delete thread"))
	}

	s.bus.Publish(ctx, &rpc.EventThreadDeleted{
		ID: thr.ID,
	})

	return nil
}

func (s *service) authoriseThreadDelete(ctx context.Context, acc *account.AccountWithEdges, thr *thread.Thread) error {
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
