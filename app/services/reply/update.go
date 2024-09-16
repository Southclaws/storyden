package reply

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

func (s *service) Update(ctx context.Context, threadID post.ID, partial Partial) (*reply.Reply, error) {
	aid, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := s.accountQuery.GetByID(ctx, aid)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	p, err := s.post_repo.Get(ctx, threadID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := acc.Roles.Permissions().Authorise(ctx, func() error {
		if p.Author.ID != aid {
			return fault.Wrap(rbac.ErrPermissions,
				fctx.With(ctx),
				fmsg.WithDesc("not owner", "You are not the owner of the post and do not have the Manage Posts permission."))
		}
		return nil
	}, rbac.PermissionManagePosts); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	opts := partial.Opts()

	p, err = s.post_repo.Update(ctx, threadID, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.indexQueue.Publish(ctx, mq.IndexPost{
		ID: p.ID,
	}); err != nil {
		s.l.Error("failed to publish index post message", zap.Error(err))
	}

	s.fetcher.HydrateContentURLs(ctx, p)

	return p, nil
}
