package thread

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/kv"
)

var ErrNoPermission = fault.New("unauthenticated user cannot view unpublished threads", ftag.With(ftag.PermissionDenied))

func (s *service) Get(
	ctx context.Context,
	threadID post.ID,
	pageParams pagination.Parameters,
) (*thread.Thread, error) {
	accountID := session.GetOptAccountID(ctx)

	ctx, span := s.ins.Instrument(ctx,
		kv.String("account_id", accountID.String()),
	)
	defer span.End()

	thr, err := s.threadQuerier.Get(ctx, threadID, pageParams, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to get thread"))
	}

	if err := authoriseThreadRead(ctx, &thr.Post); err != nil {
		return nil, err
	}

	thr.Replies.Items = s.filterRepliesByVisibility(ctx, thr.Replies.Items, accountID)

	return thr, nil
}

func authoriseThreadRead(ctx context.Context, thr *post.Post) error {
	accountID := session.GetOptAccountID(ctx)
	if thr.Visibility != visibility.VisibilityPublished {
		aid, ok := accountID.Get()
		if !ok {
			return fault.Wrap(ErrNoPermission, fctx.With(ctx))
		}

		if err := session.Authorise(ctx, func() error {
			if thr.Author.ID == aid {
				return nil
			}

			return ErrNoPermission
		}, rbac.PermissionManagePosts); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	return nil
}

func (s *service) filterRepliesByVisibility(ctx context.Context, replies []*reply.Reply, accountID opt.Optional[account.AccountID]) []*reply.Reply {
	if accountID.Ok() && session.Authorise(ctx, nil, rbac.PermissionManagePosts) == nil {
		return replies
	}

	return dt.Filter(replies, func(r *reply.Reply) bool {
		if r.Visibility == visibility.VisibilityPublished {
			return true
		}

		if r.Visibility == visibility.VisibilityReview {
			if aid, ok := accountID.Get(); ok {
				return r.Author.ID == aid
			}
		}

		return false
	})
}
